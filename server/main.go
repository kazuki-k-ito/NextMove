package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	gamepb "server/pkg/grpc"

	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

	coresdk "agones.dev/agones/pkg/sdk"
	sdk "agones.dev/agones/sdks/go"
	"google.golang.org/grpc"
)

type gameServer struct {
	gamepb.UnimplementedGameServiceServer
	characterList *CharacterList
}

func (s *gameServer) Move(stream gamepb.GameService_MoveServer) error {
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return stream.SendAndClose(&emptypb.Empty{})
		}
		if err != nil {
			return err
		}
		character := req.GetCharacter()
		log.Printf(
			"userID: %v, positionX:%v, positionY:%v, positionZ:%v, rotationY: %v, timestamp:%v",
			character.GetUserID(),
			character.GetPositionX(),
			character.GetPositionY(),
			character.GetPositionZ(),
			character.GetRotationY(),
			character.GetTimestamp(),
		)
		s.characterList.UpdateCharacter(character)
	}
}

func (s *gameServer) MoveServerStream(
	req *gamepb.MoveServerStreamRequest,
	stream gamepb.GameService_MoveServerStreamServer,
) error {
	resCount := 3
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&gamepb.MoveServerStreamResponse{
			Characters: s.characterList.GetPbCharactersExceptSelf(req.GetUserID()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func NewGameServer() *gameServer {
	return &gameServer{
		gamepb.UnimplementedGameServiceServer{},
		&CharacterList{},
	}
}

func main() {

	log.Print("Creating SDK instance")
	sdk, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf("Could not connect to sdk: %v", err)
	}

	log.Print("Starting Health Ping")
	ctx, cancel := context.WithCancel(context.Background())
	go doHealth(sdk, ctx)

	var gs *coresdk.GameServer
	gs, err = sdk.GameServer()
	if err != nil {
		log.Fatalf("Could not get gameserver port details: %s", err)
	}

	var p string
	if gs.Status.Ports != nil && len(gs.Status.Ports) > 0 {
		p = strconv.FormatInt(int64(gs.Status.Ports[0].Port), 10)
	} else {
		p = "7437"
	}
	port := &p
	listener, err := net.Listen("tcp", ":"+*port)
	log.Printf("Listen tcp port: %v", ":"+*port)
	if err != nil {
		panic(err)
	}

	debugListener := func() {
		ln, _ := net.Listen("tcp", ":28081")
		log.Println("Start debug tcp listener: 28081")
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Unable to accept incoming TCP connection: %v", err)
			}
			go tcpHandleConnection(conn, sdk, cancel)
		}
	}

	go debugListener()

	s := grpc.NewServer()

	gamepb.RegisterGameServiceServer(s, NewGameServer())

	// for grpcurl
	reflection.Register(s)

	go func() {
		log.Printf("Start gRPC server. port: %v", *port)
		s.Serve(listener)
	}()

	log.Print("Marking this server as ready")
	err = sdk.Ready()
	if err != nil {
		log.Fatalf("Could not send ready message")
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

func doHealth(sdk *sdk.SDK, ctx context.Context) {
	tick := time.Tick(2 * time.Second)
	log.Printf("Start Health Ping")
	for {
		err := sdk.Health()
		if err != nil {
			log.Fatalf("Could not send health ping, %v", err)
		}
		select {
		case <-ctx.Done():
			log.Print("Stopped health pings")
			return
		case <-tick:
		}
	}
}

func tcpHandleConnection(conn net.Conn, s *sdk.SDK, cancel context.CancelFunc) {
	log.Printf("TCP Client %s connected", conn.RemoteAddr().String())
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		tcpHandleCommand(conn, scanner.Text(), s, cancel)
	}
	log.Printf("TCP Client %s disconnected", conn.RemoteAddr().String())
}

func tcpHandleCommand(conn net.Conn, txt string, s *sdk.SDK, cancel context.CancelFunc) {
	log.Printf("TCP txt: %v", txt)

	response, addACK, err := handleResponse(txt, s, cancel)
	if err != nil {
		response = "ERROR: " + response + "\n"
	} else if addACK {
		response = "ACK TCP: " + response + "\n"
	}

	tcpRespond(conn, response)

	if response == "EXIT" {
		exit(s)
	}
}

func handleResponse(txt string, s *sdk.SDK, cancel context.CancelFunc) (response string, addACK bool, responseError error) {
	parts := strings.Split(strings.TrimSpace(txt), " ")
	response = txt
	addACK = true
	responseError = nil

	switch parts[0] {
	// shuts down the gameserver
	case "EXIT":
		// handle elsewhere, as we respond before exiting
		return
	case "GAMESERVER":
		response = gameServerName(s)
		addACK = false
	case "WATCH":
		watchGameServerEvents(s)
	}
	return
}

func gameServerName(s *sdk.SDK) string {
	var gs *coresdk.GameServer
	gs, err := s.GameServer()
	if err != nil {
		log.Fatalf("Could not retrieve GameServer: %v", err)
	}
	var j []byte
	j, err = json.Marshal(gs)
	if err != nil {
		log.Fatalf("error mashalling GameServer to JSON: %v", err)
	}
	log.Printf("GameServer: %s \n", string(j))
	return "NAME: " + gs.ObjectMeta.Name + "\n"
}

func watchGameServerEvents(s *sdk.SDK) {
	err := s.WatchGameServer(func(gs *coresdk.GameServer) {
		j, err := json.Marshal(gs)
		if err != nil {
			log.Fatalf("error mashalling GameServer to JSON: %v", err)
		}
		log.Printf("GameServer Event: %s \n", string(j))
	})
	if err != nil {
		log.Fatalf("Could not watch Game Server events, %v", err)
	}
}

func tcpRespond(conn net.Conn, txt string) {
	log.Printf("Responding to TCP with %q", txt)
	if _, err := conn.Write([]byte(txt + "\n")); err != nil {
		log.Fatalf("Could not write to TCP stream: %v", err)
	}
}

func exit(s *sdk.SDK) {
	log.Printf("Received EXIT command. Exiting.")
	// This tells Agones to shutdown this Game Server
	shutdownErr := s.Shutdown()
	if shutdownErr != nil {
		log.Printf("Could not shutdown")
	}
	// The process will exit when Agones removes the pod and the
	// container receives the SIGTERM signal
}

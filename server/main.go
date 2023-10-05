package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	gamepb "server/pkg/grpc"

	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"

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
			"userID: %v, positionX:%v, positionY:%v, positionZ:%v, rotationZ: %v, timestamp:%v",
			character.GetUserID(),
			character.GetPositionX(),
			character.GetPositionY(),
			character.GetPositionZ(),
			character.GetRotationZ(),
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
	port := 28080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	gamepb.RegisterGameServiceServer(s, NewGameServer())

	reflection.Register(s)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}

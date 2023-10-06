using System;
using System.Text;
using Cysharp.Net.Http;
using Game;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Grpc.Net.Client;
using UnityEngine;

public class RequestRunner : MonoBehaviour
{
    private AsyncClientStreamingCall<MoveRequest, Empty> _moveCall;
    private GameService.GameServiceClient _client;
    private string _userId;
    private void Awake()
    {
        var channel = GrpcChannel.ForAddress(
            "http://localhost:7155",
            new GrpcChannelOptions
            {
                HttpHandler = new YetAnotherHttpHandler
                {
                    Http2Only = true,
                },
                DisposeHttpClient = true,
            });
        _client = new GameService.GameServiceClient(channel);
        _moveCall = _client.Move();
        _userId = GenerateRandomString(16);
    }

    private string GenerateRandomString(int length)
    {
        var stringBuilder = new StringBuilder(length);
        var rand = new System.Random();

        for (var i = 0; i < length; i++)
        {
            var c = (char)('a' + rand.Next(26)); // Generating a random letter between 'a' to 'z'
            stringBuilder.Append(c);
        }

        return stringBuilder.ToString();
    }

    void Update()
    {
        if (Input.GetKeyDown(KeyCode.M))
        {
            MoveRequest();
        }
    }

    private async void MoveRequest()
    {
        var t = transform;
        var position = t.position;
        var timestamp = (ulong) (System.DateTime.UtcNow - new System.DateTime(1970, 1, 1)).TotalMilliseconds;
        await _moveCall.RequestStream.WriteAsync(new MoveRequest
        {
            Character = new Character
            {
                UserID = _userId,
                Timestamp = timestamp,
                PositionX = position.x,
                PositionY = position.y,
                PositionZ = position.z,
                RotationY = t.rotation.eulerAngles.y,
            }
        });
        
        Debug.Log($"Success Move Request (timestamp:{timestamp} rotation:{t.rotation.eulerAngles})");
    }
}

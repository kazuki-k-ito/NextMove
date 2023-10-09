using System;
using System.Text;
using Agones;
using Cysharp.Net.Http;
using Game;
using Google.Protobuf.WellKnownTypes;
using Grpc.Core;
using Grpc.Net.Client;
using UnityEngine;
using System.Collections.Generic;

public class RequestRunner : MonoBehaviour
{
    public GameObject otherCharacterPrefab;
    
    private AsyncClientStreamingCall<MoveRequest, Empty> _moveCall;
    private AsyncServerStreamingCall<MoveServerStreamResponse> _moveServerStreamCall;
    private GameService.GameServiceClient _client;
    private string _userId;
    private AgonesSdk _agones;
    private float _elapsedTime = 0f;
    private Dictionary<string, GameObject> _otherCharacters;
    private void Awake()
    {
        _otherCharacters = new Dictionary<string, GameObject>();
        var channel = GrpcChannel.ForAddress(
            "http://localhost:7654",
            new GrpcChannelOptions
            {
                HttpHandler = new YetAnotherHttpHandler
                {
                    Http2Only = true,
                },
                DisposeHttpClient = true,
            });
        _userId = GenerateRandomString(16);
        _client = new GameService.GameServiceClient(channel);
        _moveCall = _client.Move();
        MoveRequest();
        _moveServerStreamCall = _client.MoveServerStream(new MoveServerStreamRequest{UserID = _userId});
        MoveServerStream();
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
        _elapsedTime += Time.deltaTime;
        if (_elapsedTime >= 1.0f)
        {
            MoveRequest();
            _moveServerStreamCall = _client.MoveServerStream(new MoveServerStreamRequest{UserID = _userId});
            _elapsedTime = 0.0f;
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

    private async void MoveServerStream()
    {
        while (await _moveServerStreamCall.ResponseStream.MoveNext())
        {
            var response = _moveServerStreamCall.ResponseStream.Current;
            
            if (response == null || response.Characters == null) continue;
            foreach(var character in response.Characters)
            {
                // オブジェクトを生成
                if (!_otherCharacters.ContainsKey(character.UserID))
                {
                    _otherCharacters.Add(
                        character.UserID,
                        Instantiate(
                            otherCharacterPrefab,
                            new Vector3(character.PositionX, character.PositionY, character.PositionZ),
                            Quaternion.Euler(new Vector3(0, character.RotationY, 0))
                            )
                        );
                    Debug.Log("生成");
                }
                else
                {
                    var c = _otherCharacters[character.UserID];
                    c.transform.position = new Vector3(character.PositionX, character.PositionY, character.PositionZ);
                    c.transform.rotation = Quaternion.Euler(new Vector3(0, character.RotationY, 0));
                    Debug.Log("更新");
                }
                Debug.Log($"Receive Position. UserID:{character.UserID} Position:{character.PositionX},{character.PositionY},{character.PositionZ} Rotation:{character.RotationY}");
            }
        }
    }

    private void OnApplicationQuit()
    {
        _moveCall.Dispose();
        _moveServerStreamCall.Dispose();
    }
}

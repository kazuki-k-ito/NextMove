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
using UnityEngine.UIElements;

public class RequestRunner : MonoBehaviour
{
    class CharacterInfo
    {
        public Vector3 Pos;
        public Quaternion Rot;
        public ulong Time;
        public float rate;
    }
    
    public GameObject otherCharacterPrefab;
    
    private AsyncClientStreamingCall<MoveRequest, Empty> _moveCall;
    private AsyncServerStreamingCall<MoveServerStreamResponse> _moveServerStreamCall;
    private GameService.GameServiceClient _client;
    private string _userId;
    private AgonesSdk _agones;
    private float _elapsedTime = 0f;
    private Dictionary<string, GameObject> _otherCharacters;
    private Dictionary<string, List<CharacterInfo>> _otherCharacterInfo;
    private void Awake()
    {
        _otherCharacters = new Dictionary<string, GameObject>();
        _otherCharacterInfo = new Dictionary<string, List<CharacterInfo>>();
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
        if (_elapsedTime >= 0.1f)
        {
            MoveRequest();
            _moveServerStreamCall = _client.MoveServerStream(new MoveServerStreamRequest{UserID = _userId});
            _elapsedTime = 0.0f;
        }
    }
    
    private ulong Milliseconds()
    {
        return (ulong) (System.DateTime.UtcNow - new System.DateTime(1970, 1, 1)).TotalMilliseconds;
    }

    private void FixedUpdate()
    {
        foreach (var characterInfo in _otherCharacterInfo)
        {
            var character = _otherCharacters[characterInfo.Key];
            var t = character.transform;
            if (Vector3.Distance(characterInfo.Value[1].Pos, t.position) <= 0.05f) continue;

            var cOld = characterInfo.Value[0];
            var cNew = characterInfo.Value[1];
            var now = Milliseconds();
            cNew.rate += 0.1f;
            var rate = cNew.rate;
            Debug.Log($"rate:{rate} now:{now - 1000} cOldTime:{cOld.Time} cNewTime:{cNew.Time + 1000}");
            Debug.Log($"cOldPos:{cOld.Pos} cNewPos:{cNew.Pos}");
            character.transform.rotation = Quaternion.Slerp(t.rotation, cNew.Rot, rate);
            character.transform.position = Vector3.Lerp( cOld.Pos, cNew.Pos , rate);
        }
    }

    private async void MoveRequest()
    {
        var t = transform;
        var position = t.position;
        var timestamp = Milliseconds();
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
                    var characterInfo = new CharacterInfo();
                    characterInfo.Pos = new Vector3(
                        character.PositionX, character.PositionY, character.PositionZ
                    );
                    characterInfo.Rot = Quaternion.Euler(new Vector3(0, character.RotationY, 0));
                    characterInfo.Time = Milliseconds();
                    characterInfo.rate = 0.0f;
                    _otherCharacterInfo.Add(character.UserID, new List<CharacterInfo>(){characterInfo, characterInfo});
                    Debug.Log("生成");
                }
                else
                {
                    var characterInfo = _otherCharacterInfo[character.UserID];
                    var cNew = new CharacterInfo();
                    cNew.Pos = new Vector3(
                        character.PositionX, character.PositionY, character.PositionZ
                    );
                    cNew.Rot = Quaternion.Euler(new Vector3(0, character.RotationY, 0));
                    cNew.Time = character.Timestamp;
                    cNew.rate = 0.0f;
                    
                    characterInfo.Add(cNew);
                    characterInfo.RemoveAt(0);
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

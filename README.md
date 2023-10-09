# Installation

## protobuf
### Go言語
参考) https://zenn.dev/hsaki/books/golang-grpc-starting/viewer/codegenerate

protobufのインストール
```
brew install protobuf
```

go moduleのインストール
```
cd server
go mod init server
go install google.golang.org/grpc
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
```

ディレクトリの作成
```
mkdir server/pkg/grpc
```

コード生成
```
protoc --go_out=../server/pkg/grpc --go_opt=paths=source_relative --go-grpc_out=../server/pkg/grpc --go-grpc_opt=paths=source_relative game.proto
```

grpccurlのインストール
```
brew install grpcurl
```

## C#
protobufのdllをインストール
```
curl -Lo protobuf.zip https://www.nuget.org/api/v2/package/Google.Protobuf/3.24.3
unzip protobuf.zip -d protobuf
mkdir -p unity/Assets/Plugins/Google.Protobuf
cp protobuf/lib/net45/Google.Protobuf.dll ../unity/NextMove/Assets/Plugins/Google.Protobuf
```

Grpc.Toolsから必要なコマンドを取り込む
```
curl -Lo grpctools.zip https://www.nuget.org/api/v2/package/Grpc.Tools/2.58.0
unzip grpctools.zip -d grpctools
chmod +x grpctools/tools/macosx_x64/grpc_csharp_plugin
```

## minikube
```
brew install minikube
eval $(minikube -p minikube docker-env)
minikube config set memory 4096
minikube config set cpus 4
minikube start
```

## Agones
https://agones.dev/site/docs/installation/install-agones/yaml/

インストール
```
kubectl apply --server-side -f https://raw.githubusercontent.com/googleforgames/agones/release-1.35.0/install/yaml/install.yaml
```

fleetの作成
```
kubectl apply -f agones/fleet.yaml
kubectl apply -f agones/fleetautoscaler.yaml
```

## Unity
2022.3.10f1

### YetAnotherHttpHandlerのインストール
依存パッケージ
```
https://github.com/Cysharp/YetAnotherHttpHandler/releases/tag/redist-20230728-01
```

本体
Add package from git URL...
```
https://github.com/Cysharp/YetAnotherHttpHandler.git?path=src/YetAnotherHttpHandler#v0.1.0
```

# Quickstart

## Server

### ビルドとデプロイ
```
scripts/build_and_deploy.sh
```

### port開放
```
scripts/port_forward.sh
```

## Client
シーンを複数クライアントで再生
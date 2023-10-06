#!/bin/sh -xe

TOP_DIR=$(git rev-parse --show-toplevel)

GRPC_PLUGIN=$TOP_DIR/unity/NextMove/Tools/grpc_csharp_plugin

protoc -I $TOP_DIR/schema/ \
  --csharp_out=$TOP_DIR/unity/NextMove/Assets/Scripts/Protos \
  --grpc_out=$TOP_DIR/unity/NextMove/Assets/Scripts/Protos \
  --plugin=protoc-gen-grpc=$GRPC_PLUGIN \
  $TOP_DIR/schema/game.proto

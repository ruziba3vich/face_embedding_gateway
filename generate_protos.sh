#!/bin/bash
# chmod +x generate_protos.sh

set -e

PROTO_DIR="./"

OUT_DIR="./"

mkdir -p $OUT_DIR

protoc -I=$PROTO_DIR \
  --go_out=$OUT_DIR \
  --go-grpc_out=$OUT_DIR \
  $PROTO_DIR/face_recognition_protos/embed.proto

echo "protos generated successfully"

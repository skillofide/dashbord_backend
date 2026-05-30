#!/usr/bin/env bash
# scripts/proto-gen.sh
# Generates Go code from .proto files using protoc + protoc-gen-go.
# NOTE: With the JSON codec approach, this script is optional —
# the hand-written Go files in proto/ already implement the service interfaces.
# Run this script ONLY if you want to switch back to binary protobuf encoding.
#
# Prerequisites:
#   brew install protobuf
#   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

set -euo pipefail

PROTO_DIR="$(dirname "$0")/../proto"
OUT_DIR="$PROTO_DIR"

PROTO_FILES=(
  "problem/v1/problem.proto"
  "submission/v1/submission.proto"
  "execution/v1/execution.proto"
  "progress/v1/progress.proto"
  "notification/v1/notification.proto"
)

echo "Generating proto files..."

for f in "${PROTO_FILES[@]}"; do
  echo "  → $f"
  protoc \
    --proto_path="$PROTO_DIR" \
    --go_out="$OUT_DIR" --go_opt=paths=source_relative \
    --go-grpc_out="$OUT_DIR" --go-grpc_opt=paths=source_relative \
    "$PROTO_DIR/$f"
done

echo "Done! Proto files generated in $OUT_DIR"

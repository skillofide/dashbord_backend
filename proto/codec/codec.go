// Package codec registers a JSON-based gRPC codec under the name "proto",
// replacing the default protobuf binary codec. This means our service
// messages can be plain Go structs with JSON tags — no protoc needed.
package codec

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

func init() {
	encoding.RegisterCodec(JSONCodec{})
}

// JSONCodec implements encoding.Codec using JSON marshaling.
// It registers itself as "proto" so gRPC uses it transparently
// on both client and server sides.
type JSONCodec struct{}

func (JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Name returns "proto" to override the default protobuf codec.
func (JSONCodec) Name() string { return "proto" }

package proto

import (
	"encoding/json"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func MarshalProtoSliceToJSON[T proto.Message](messages []T, marshaler *protojson.MarshalOptions) ([]byte, error) {
	if len(messages) == 0 {
		return []byte("[]"), nil
	}

	if marshaler == nil {
		marshaler = &protojson.MarshalOptions{EmitUnpopulated: true, UseEnumNumbers: true}
	}

	// Create a slice to hold the JSON representation
	jsonObjects := make([]json.RawMessage, len(messages))

	for i, msg := range messages {
		jsonBytes, err := marshaler.Marshal(msg)
		if err != nil {
			return nil, err
		}

		jsonObjects[i] = json.RawMessage(jsonBytes)
	}

	// Marshal the slice to JSON objects
	return json.Marshal(jsonObjects)
}

func UnmarshalJSONToProtoSlice[T proto.Message](jsonData []byte, newMessage func() T) ([]T, error) {
	if len(jsonData) == 0 || string(jsonData) == "[]" || string(jsonData) == "null" {
		return []T{}, nil
	}

	// First unmarshal to slice to json.RawMessage
	var rawMessages []json.RawMessage
	if err := json.Unmarshal(jsonData, &rawMessages); err != nil {
		return nil, err
	}

	unmarshaler := protojson.UnmarshalOptions{}

	// Create slice to hold the proto messages
	messages := make([]T, len(rawMessages))

	for i, rawMsg := range rawMessages {
		msg := newMessage()
		if err := unmarshaler.Unmarshal(rawMsg, msg); err != nil {
			return nil, err
		}
		messages[i] = msg
	}

	return messages, nil
}

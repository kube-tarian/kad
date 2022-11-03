package model

import "encoding/json"

type RequestPayload struct {
	SubAction string          `json:"sub_action"`
	Data      json.RawMessage `json:"data,omitempty"` // TODO: This will be enhanced along with plugin implementation
}

type ResponsePayload struct {
	Status  string          `json:"status"`
	Message json.RawMessage `json:"message,omitempty"` // TODO: This will be enhanced along with plugin implementation
}

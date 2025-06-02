package protocol

import "encoding/json"

const (
	TypeRegionRequest = "region_request" // client -> server
	TypeRegionCount   = "region_count"   // server -> client
)

type Envelope struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type RegionRequestPayload struct {
	Region string `json:"region"`
}

type RegionCountPayload struct {
	Count int `json:"count"`
}

func Marshal(msgType string, payload any) ([]byte, error) {
	innerBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	env := Envelope{
		Type:    msgType,
		Payload: innerBytes,
	}
	return json.Marshal(env)
}

func Unmarshal(data []byte) (*Envelope, error) {
	var env Envelope
	err := json.Unmarshal(data, &env)
	if err != nil {
		return nil, err
	}
	return &env, nil
}

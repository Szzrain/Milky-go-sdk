package Milky_go_sdk

import (
	"encoding/json"
	"testing"
)

func TestForwardElementMarshalJSONUsesSenderName(t *testing.T) {
	tests := []struct {
		name    string
		message OutgoingForwardedMessage
	}{
		{
			name: "sender name",
			message: OutgoingForwardedMessage{
				UserID:     10010,
				SenderName: "MilkyBot",
				Segments:   []IMessageElement{&TextElement{Text: "card"}},
			},
		},
		{
			name: "deprecated name fallback",
			message: OutgoingForwardedMessage{
				UserID:   10010,
				Name:     "MilkyBot",
				Segments: []IMessageElement{&TextElement{Text: "card"}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(&ForwardElement{Messages: []OutgoingForwardedMessage{test.message}})
			if err != nil {
				t.Fatalf("failed to marshal forward element: %v", err)
			}

			const expected = `{"type":"forward","data":{"messages":[{"user_id":10010,"sender_name":"MilkyBot","segments":[{"type":"text","data":{"text":"card"}}]}]}}`
			if string(data) != expected {
				t.Fatalf("forward element JSON = %s, want %s", data, expected)
			}
		})
	}
}

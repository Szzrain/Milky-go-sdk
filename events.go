package Milky_go_sdk

import (
	"encoding/json"
)

// Event provides a basic initial struct for all websocket events.
type Event struct {
	Type    string          `json:"event_type"`
	RawData json.RawMessage `json:"data"`
	// Struct contains one of the other types in this file.
	Struct interface{} `json:"-"`
}

type ReceiveMessage struct {
	PeerId       int64  `json:"peer_id"`
	MessageSeq   int64  `json:"message_seq"`
	SenderId     int64  `json:"sender_id"`
	Time         int64  `json:"time"`
	MessageScene string `json:"message_scene"` // "friend", "group", etc.

	Segments    []IMessageElement `json:"segments"`
	Group       *GroupInfo        `json:"group"`
	Friend      *Friend           `json:"friend"`
	GroupMember *GroupMemberInfo  `json:"group_member"`
}

type MessageRet struct {
	MessageSeq int64 `json:"message_seq"`
	Time       int64 `json:"time"`
}

func (r *ReceiveMessage) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["peer_id"], &r.PeerId); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["message_seq"], &r.MessageSeq); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["sender_id"], &r.SenderId); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["time"], &r.Time); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["message_scene"], &r.MessageScene); err != nil {
		return err
	}
	if rawGroup, ok := raw["group"]; ok {
		var group GroupInfo
		if err := json.Unmarshal(rawGroup, &group); err != nil {
			return err
		}
		r.Group = &group
	}
	if rawGroupMember, ok := raw["group_member"]; ok {
		var groupMember GroupMemberInfo
		if err := json.Unmarshal(rawGroupMember, &groupMember); err != nil {
			return err
		}
		r.GroupMember = &groupMember
	}
	if rawFriend, ok := raw["friend"]; ok {
		var friend Friend
		if err := json.Unmarshal(rawFriend, &friend); err != nil {
			return err
		}
		r.Friend = &friend
	}
	if rawSegments, ok := raw["segments"]; ok {
		var elements []RawMessageElement
		if err := json.Unmarshal(rawSegments, &elements); err != nil {
			return err
		}
		for _, element := range elements {
			switch element.Type {
			case string(Text):
				var textElement TextElement
				if err := json.Unmarshal(element.Data, &textElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &textElement)
			case string(At):
				var atElement AtElement
				if err := json.Unmarshal(element.Data, &atElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &atElement)
			case string(Image):
				var imageElement ImageElement
				if err := json.Unmarshal(element.Data, &imageElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &imageElement)
			case string(Record):
				var recordElement RecordElement
				if err := json.Unmarshal(element.Data, &recordElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &recordElement)
			case string(Face):
				var faceElement FaceElement
				if err := json.Unmarshal(element.Data, &faceElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &faceElement)
			case string(Reply):
				var replyElement ReplyElement
				if err := json.Unmarshal(element.Data, &replyElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &replyElement)
			case string(Forward):
				var forwardElement ForwardElement
				if err := json.Unmarshal(element.Data, &forwardElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &forwardElement)
			default:
				continue // Ignore unknown types.
			}
		}
	}
	return nil
}

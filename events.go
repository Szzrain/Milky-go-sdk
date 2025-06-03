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

type IncomingMessage struct {
	PeerId     int64             `json:"peer_id"`
	MessageSeq int64             `json:"message_seq"`
	SenderId   int64             `json:"sender_id"`
	Time       int64             `json:"time"`
	Segments   []IMessageElement `json:"segments"`
}

func (i *IncomingMessage) UnmarshalJSON(data []byte) error {
	raw := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["peer_id"], &i.PeerId); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["message_seq"], &i.MessageSeq); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["sender_id"], &i.SenderId); err != nil {
		return err
	}
	if err := json.Unmarshal(raw["time"], &i.Time); err != nil {
		return err
	}
	if rawSegments, ok := raw["segments"]; ok {
		var elements []RawMessageElement
		if err := json.Unmarshal(rawSegments, &elements); err != nil {
			return err
		}
		for _, element := range elements {
			switch element.Type {
			case string(Text):
				var textElement ReceiveTextElement
				if err := json.Unmarshal(element.Data, &textElement); err != nil {
					return err
				}
				i.Segments = append(i.Segments, &textElement)
			case string(At):
				var atElement ReceiveAtElement
				if err := json.Unmarshal(element.Data, &atElement); err != nil {
					return err
				}
				i.Segments = append(i.Segments, &atElement)
			case string(Image):
				var imageElement ImageElement
				if err := json.Unmarshal(element.Data, &imageElement); err != nil {
					return err
				}
				i.Segments = append(i.Segments, &imageElement)
			case string(Record):
				var recordElement RecordElement
				if err := json.Unmarshal(element.Data, &recordElement); err != nil {
					return err
				}
				i.Segments = append(i.Segments, &recordElement)
			case string(Face):
				var faceElement FaceElement
				if err := json.Unmarshal(element.Data, &faceElement); err != nil {
					return err
				}
				i.Segments = append(i.Segments, &faceElement)
			case string(Reply):
				var replyElement ReplyElement
				if err := json.Unmarshal(element.Data, &replyElement); err != nil {
					return err
				}
				i.Segments = append(i.Segments, &replyElement)
			default:
				continue // Ignore unknown types.
			}
		}

	}
	return nil
}

type (
	RawMessageElement struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	IMessageElement interface {
		Type() ElementType
	}

	ElementType string
)

// Copied from sealdice/sealdice-core
const (
	Text   ElementType = "text"        // 文本
	At                 = "mention"     // 艾特
	AtAll              = "mention_all" // 艾特
	Video              = "video"       // 文件
	Image              = "image"       // 图片
	Reply              = "reply"       // 回复
	Record             = "record"      // 语音
	Face               = "face"        // 表情
)

const maxFileSize = 1024 * 1024 * 50 // 50MB

type ReceiveTextElement struct {
	Text string `json:"text"`
}

func (t *ReceiveTextElement) Type() ElementType {
	return Text
}

type ReceiveAtElement struct {
	UserID string `json:"user_id"`
}

func (t *ReceiveAtElement) Type() ElementType {
	return At
}

type ReplyElement struct {
	MessageSeq string `json:"message_seq"` // 回复的目标消息ID
}

func (t *ReplyElement) Type() ElementType {
	return Reply
}

type ImageElement struct {
	ResourceID string `json:"resource_id"` // 图片资源ID
	TempURL    string `json:"temp_url"`
	Summary    string `json:"summary"`  // 图片摘要
	SubType    string `json:"sub_type"` // 图片子类型
}

func (l *ImageElement) Type() ElementType {
	return Image
}

type RecordElement struct {
	ResourceID string `json:"resource_id"` // 资源ID
	TempURL    string `json:"temp_url"`
	Duration   int32  `json:"duration"`
}

func (r *RecordElement) Type() ElementType {
	return Record
}

type FaceElement struct {
	FaceID string `json:"face_id"`
}

func (f *FaceElement) Type() ElementType {
	return Face
}

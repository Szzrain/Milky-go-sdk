package Milky_go_sdk

import "encoding/json"

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
	At     ElementType = "mention"     // 艾特
	AtAll  ElementType = "mention_all" // 艾特全体
	Video  ElementType = "video"       // 视频
	Image  ElementType = "image"       // 图片
	Reply  ElementType = "reply"       // 回复
	Record ElementType = "record"      // 语音
	Face   ElementType = "face"        // 表情
)

type GroupInfo struct {
	GroupId        int64  `json:"group_id"`
	Name           string `json:"name"`
	MemberCount    int32  `json:"member_count"`
	MaxMemberCount int32  `json:"max_member_count"`
}

type GroupMemberInfo struct {
	GroupId      int64  `json:"group_id"`
	UserId       int64  `json:"user_id"`
	Nickname     string `json:"nickname"`
	Card         string `json:"card"`
	Title        string `json:"title"`
	Sex          string `json:"sex"`
	Level        int32  `json:"level"`
	Role         string `json:"role"` // "owner", "admin", "member"
	JoinTime     int64  `json:"join_time"`
	LastSentTime int64  `json:"last_sent_time"`
}

const maxFileSize = 1024 * 1024 * 50 // 50MB

type ReceiveTextElement struct {
	Text string `json:"text"`
}

func (t *ReceiveTextElement) Type() ElementType {
	return Text
}

type ReceiveAtElement struct {
	UserID int64 `json:"user_id"`
}

func (t *ReceiveAtElement) Type() ElementType {
	return At
}

type ReplyElement struct {
	MessageSeq int64 `json:"message_seq"` // 回复的目标消息ID
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

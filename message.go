package Milky_go_sdk

import "encoding/json"

type (
	RawMessageElement struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	IMessageElement interface {
		Type() ElementType
		MarshalJSON() ([]byte, error)
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

type APIResponse struct {
	Status  string          `json:"status"`
	RetCode int             `json:"retcode"`
	Data    json.RawMessage `json:"data"`
}

type LoginInfo struct {
	UIN      int64  `json:"uin"`
	Nickname string `json:"nickname"`
}

type FriendCategory struct {
	CategoryID   int32  `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type Friend struct {
	UserID   int64           `json:"user_id"`
	QID      int64           `json:"qid,omitempty"`
	Nickname string          `json:"nickname"`
	Sex      string          `json:"sex"`
	Remark   string          `json:"remark"`
	Category *FriendCategory `json:"category"`
}

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

type GroupAnnouncement struct {
	GroupId        int64  `json:"group_id"`
	AnnouncementID string `json:"announcement_id"`
	UserID         int64  `json:"user_id"`
	Time           int64  `json:"time"`
	Content        string `json:"content"`
	ImageURL       string `json:"image_url,omitempty"`
}

type GroupFile struct {
	GroupId         int64  `json:"group_id"`
	FileID          string `json:"file_id"`                    // 文件ID
	FileName        string `json:"file_name"`                  // 文件名
	ParentFolderID  string `json:"parent_folder_id,omitempty"` // 父文件夹ID
	FileSize        int64  `json:"file_size"`                  // 文件大小
	UploadedTime    int64  `json:"uploaded_time"`              // 上传时间
	ExpireTime      int64  `json:"expire_time"`                // 过期时间
	UploaderID      int64  `json:"uploader_id"`                // 上传者ID
	DownloadedTimes int32  `json:"downloaded_times"`           // 下载次数
}

type GroupFolder struct {
	GroupId          int64  `json:"group_id"`
	FolderID         string `json:"folder_id"` // 文件夹ID
	ParentFolderID   string `json:"parent_folder_id,omitempty"`
	FolderName       string `json:"folder_name"`        // 文件夹名称
	CreatedTime      int64  `json:"created_time"`       // 创建时间
	LastModifiedTime int64  `json:"last_modified_time"` // 最后修改时间
	CreatorID        int64  `json:"creator_id"`
	FileCount        int32  `json:"file_count"`
}

type FriendRequest struct {
	RequestID   string `json:"request_id"`
	Time        int64  `json:"time"`
	IsFiltered  bool   `json:"is_filtered"`  // 是否被过滤
	InitiatorID int64  `json:"initiator_id"` // 发起者ID
	State       string `json:"state"`
	Comment     string `json:"comment,omitempty"`
	Via         string `json:"via,omitempty"`
}

type GroupRequest struct {
	RequestID   string `json:"request_id"`
	Time        int64  `json:"time"`
	IsFiltered  bool   `json:"is_filtered"`  // 是否被过滤
	InitiatorID int64  `json:"initiator_id"` // 发起者ID
	State       string `json:"state"`
	GroupID     int64  `json:"group_id"`
	OperatorID  int64  `json:"operator_id,omitempty"`
	RequestType string `json:"request_type"`
	Comment     string `json:"comment,omitempty"`
	InviteeID   int64  `json:"invitee_id,omitempty"`
}

type GroupInvitation struct {
	RequestID   string `json:"request_id"`
	Time        int64  `json:"time"`
	IsFiltered  bool   `json:"is_filtered"`  // 是否被过滤
	InitiatorID int64  `json:"initiator_id"` // 发起者ID
	State       string `json:"state"`
	GroupID     int64  `json:"group_id"`
}

const maxFileSize = 1024 * 1024 * 50 // 50MB

type ReceiveTextElement struct {
	Text string `json:"text"`
}

func (t *ReceiveTextElement) Type() ElementType {
	return Text
}

func (t *ReceiveTextElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type ElementType `json:"type"`
		Data struct {
			Text string `json:"text"` // 文本内容
		} `json:"data"`
	}{
		Type: t.Type(),
		Data: struct {
			Text string `json:"text"` // 文本内容
		}{
			Text: t.Text,
		},
	})
}

type ReceiveAtElement struct {
	UserID int64 `json:"user_id"`
}

func (t *ReceiveAtElement) Type() ElementType {
	return At
}

func (t *ReceiveAtElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type ElementType `json:"type"`
		Data struct {
			UserID int64 `json:"user_id"`
		} `json:"data"`
	}{
		Type: t.Type(),
		Data: struct {
			UserID int64 `json:"user_id"`
		}{
			UserID: t.UserID,
		},
	})
}

type ReceiveAtAllElement struct{}

func (t *ReceiveAtAllElement) Type() ElementType {
	return AtAll
}

func (t *ReceiveAtAllElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type ElementType `json:"type"`
		Data struct{}    `json:"data"`
	}{
		Type: t.Type(),
		Data: struct{}{},
	})
}

type ReplyElement struct {
	MessageSeq int64 `json:"message_seq"` // 回复的目标消息ID
	ClientSeq  int64 `json:"client_seq"`
}

func (t *ReplyElement) Type() ElementType {
	return Reply
}

func (t *ReplyElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type ElementType `json:"type"`
		Data struct {
			MessageSeq int64 `json:"message_seq"`
			ClientSeq  int64 `json:"client_seq,omitempty"`
		} `json:"data"`
	}{
		Type: t.Type(),
		Data: struct {
			MessageSeq int64 `json:"message_seq"`
			ClientSeq  int64 `json:"client_seq,omitempty"`
		}{
			MessageSeq: t.MessageSeq,
			ClientSeq:  t.ClientSeq,
		},
	})
}

type ImageElement struct {
	URI        string `json:"uri,omitempty"` // 图片URI
	ResourceID string `json:"resource_id"`   // 图片资源ID
	TempURL    string `json:"temp_url"`
	Summary    string `json:"summary"`  // 图片摘要
	SubType    string `json:"sub_type"` // 图片子类型
}

func (l *ImageElement) Type() ElementType {
	return Image
}

func (l *ImageElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type ElementType `json:"type"`
		Data struct {
			URI     string `json:"uri"`
			Summary string `json:"summary"`
			SubType string `json:"sub_type"`
		} `json:"data"`
	}{
		Type: l.Type(),
		Data: struct {
			URI     string `json:"uri"`
			Summary string `json:"summary"`
			SubType string `json:"sub_type"`
		}{
			URI:     l.URI,
			Summary: l.Summary,
			SubType: l.SubType,
		},
	})
}

type RecordElement struct {
	URI        string `json:"uri,omitempty"` // 语音URI
	ResourceID string `json:"resource_id"`   // 资源ID
	TempURL    string `json:"temp_url"`
	Duration   int32  `json:"duration"`
}

func (r *RecordElement) Type() ElementType {
	return Record
}

func (r *RecordElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type ElementType `json:"type"`
		Data struct {
			URI string `json:"uri"`
		} `json:"data"`
	}{
		Type: r.Type(),
		Data: struct {
			URI string `json:"uri"`
		}{
			URI: r.URI,
		},
	})
}

type FaceElement struct {
	FaceID string `json:"face_id"`
}

func (f *FaceElement) Type() ElementType {
	return Face
}

func (f *FaceElement) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Type ElementType `json:"type"`
		Data struct {
			FaceID string `json:"face_id"`
		} `json:"data"`
	}{
		Type: f.Type(),
		Data: struct {
			FaceID string `json:"face_id"`
		}{
			FaceID: f.FaceID,
		},
	})
}

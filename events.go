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
			case string(AtAll):
				var atAllElement AtAllElement
				if err := json.Unmarshal(element.Data, &atAllElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &atAllElement)
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
			case string(Video):
				var videoElement VideoElement
				if err := json.Unmarshal(element.Data, &videoElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &videoElement)
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
			case string(MarketFace):
				var marketFaceElement MarketFaceElement
				if err := json.Unmarshal(element.Data, &marketFaceElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &marketFaceElement)
			case string(LightApp):
				var lightAppElement LightAppElement
				if err := json.Unmarshal(element.Data, &lightAppElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &lightAppElement)
			case string(XML):
				var xmlElement XmlElement
				if err := json.Unmarshal(element.Data, &xmlElement); err != nil {
					return err
				}
				r.Segments = append(r.Segments, &xmlElement)
			default:
				continue // Ignore unknown types.
			}
		}
	}
	return nil
}

type APIResponse struct {
	Status  string          `json:"status"`
	RetCode int             `json:"retcode"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message,omitempty"` // 错误信息
}

type LoginInfo struct {
	UIN      int64  `json:"uin"`
	Nickname string `json:"nickname"`
}

type ImplInfo struct {
	ImplName          string `json:"impl_name"`    // 实现名称
	ImplVersion       string `json:"impl_version"` // 实现版本
	QQProtocolVersion string `json:"qq_protocol_version"`
	QQProtocolType    string `json:"qq_protocol_type"`
	MilkyVersion      string `json:"milky_version"` // Milky 协议版本
}

type UserProfile struct {
	Nickname string `json:"nickname"`
	Qid      string `json:"qid"`
	Age      int32  `json:"age"`
	Sex      string `json:"sex"`
	Remark   string `json:"remark"`
	Bio      string `json:"bio"`
	Level    int32  `json:"level"`
	Country  string `json:"country"`
	City     string `json:"city"`
	School   string `json:"school"`
}

type FriendCategory struct {
	CategoryID   int32  `json:"category_id"`
	CategoryName string `json:"category_name"`
}

type Friend struct {
	UserID   int64           `json:"user_id"`
	QID      string          `json:"qid,omitempty"`
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
	InitiatorUID string `json:"initiator_uid"`
	InitiatorID  int64  `json:"initiator_id"` // 发起者ID
	Comment      string `json:"comment"`
	Via          string `json:"via"`
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
	InvitationSeq int64 `json:"invitation_seq"`
	InitiatorID   int64 `json:"initiator_id"` // 发起者ID
	GroupID       int64 `json:"group_id"`
}

type MessageRecall struct {
	MessageScene  string `json:"message_scene"` // "friend", "group", etc.
	PeerID        int64  `json:"peer_id"`
	MessageSeq    int64  `json:"message_seq"`
	SenderID      int64  `json:"sender_id"`
	OperatorID    int64  `json:"operator_id"`
	DisplaySuffix string `json:"display_suffix"`
}

type FriendNudge struct {
	UserID              int64  `json:"user_id"`
	IsSelfSend          bool   `json:"is_self_send"`
	IsSelfReceive       bool   `json:"is_self_receive"`
	DisplayAction       string `json:"display_action"`
	DisplaySuffix       string `json:"display_suffix"`
	DisplayActionImgUrl string `json:"display_action_img_url"`
}

type BotOffline struct {
	Reason string `json:"reason"` // 原因
}

type GroupNudge struct {
	GroupID             int64  `json:"group_id"`    // 群号
	SenderID            int64  `json:"sender_id"`   // 发送者ID
	ReceiverID          int64  `json:"receiver_id"` // 接收者ID
	DisplayAction       string `json:"display_action"`
	DisplaySuffix       string `json:"display_suffix"`
	DisplayActionImgUrl string `json:"display_action_img_url"`
}

type GroupMessageReaction struct {
	GroupID    int64  `json:"group_id"` // 群号
	UserID     int64  `json:"user_id"`
	MessageSeq int64  `json:"message_seq"`
	FaceID     string `json:"face_id"`
	IsAdd      bool   `json:"is_add"`
}

type GroupMute struct {
	GroupID    int64 `json:"group_id"` // 群号
	UserID     int64 `json:"user_id"`
	OperatorID int64 `json:"operator_id"`
	Duration   int32 `json:"duration"` // 禁言时长，单位秒，0表示取消禁言
}

type GroupWholeMute struct {
	GroupID    int64 `json:"group_id"` // 群号
	OperatorID int64 `json:"operator_id"`
	IsMute     bool  `json:"is_mute"`
}

type GroupMemberIncrease struct {
	GroupID    int64 `json:"group_id"`              // 群号
	UserID     int64 `json:"user_id"`               // 加群用户ID
	OperatorID int64 `json:"operator_id,omitempty"` // 操作人ID
	InvitorID  int64 `json:"invitor_id,omitempty"`  // 邀请人ID
}

type GroupMemberDecrease struct {
	GroupID    int64 `json:"group_id"`    // 群号
	UserID     int64 `json:"user_id"`     // 退群用户ID
	OperatorID int64 `json:"operator_id"` // 操作人ID
}

package Milky_go_sdk

import "fmt"

type apiEndpoints struct {
	Gateway string
}

// NewAPIEndpoints 创建一个新的 apiEndpoints 实例
func newAPIEndpoints(gateway string) *apiEndpoints {
	return &apiEndpoints{Gateway: gateway}
}

// Endpoint 生成完整的 API 端点
func (e *apiEndpoints) Endpoint(path string) string {
	return fmt.Sprintf("%s/api/%s", e.Gateway, path)
}

const (
	EndpointGetLoginInfo               = "get_login_info"
	EndpointGetFriendList              = "get_friend_list"
	EndpointGetFriendInfo              = "get_friend_info"
	EndpointGetGroupList               = "get_group_list"
	EndpointGetGroupInfo               = "get_group_info"
	EndpointGetGroupMemberList         = "get_group_member_list"
	EndpointGetGroupMemberInfo         = "get_group_member_info"
	EndpointSendPrivateMessage         = "send_private_message"
	EndpointSendGroupMessage           = "send_group_message"
	EndpointGetMessage                 = "get_message"
	EndpointGetHistoryMessage          = "get_history_message"
	EndpointGetResourceTempURL         = "get_resource_temp_url"
	EndpointGetForwardedMessages       = "get_forwarded_messages"
	EndpointRecallPrivateMessage       = "recall_private_message"
	EndpointRecallGroupMessage         = "recall_group_message"
	EndpointSendFriendNudge            = "send_friend_nudge"
	EndpointGetSendProfileLike         = "send_profile_like"
	EndpointSetGroupName               = "set_group_name"
	EndpointSetGroupAvatar             = "set_group_avatar"
	EndpointSetGroupMemberCard         = "set_group_member_card"
	EndpointSetGroupMemberSpecialTitle = "set_group_member_special_title"
	EndpointSetGroupAdmin              = "set_group_admin"
	EndpointSetGroupMemberMute         = "set_group_member_mute"
	EndpointSetGroupWholeMute          = "set_group_whole_mute"
	EndpointKickGroupMember            = "kick_group_member"
	EndpointGetGroupAnnouncementList   = "get_group_announcement_list"
	EndpointSendGroupAnnouncement      = "send_group_announcement"
	EndpointDeleteGroupAnnouncement    = "delete_group_announcement"
	EndpointQuitGroup                  = "quit_group"
	EndpointSendGroupMessageReaction   = "send_group_message_reaction"
	EndpointSendGroupNudge             = "send_group_nudge"
	EndpointGetFriendRequests          = "get_friend_requests"
	EndpointGetGroupRequests           = "get_group_requests"
	EndpointGetGroupInvitations        = "get_group_invitations"
	EndpointAcceptRequest              = "accept_request"
	EndpointRejectRequest              = "reject_request"
	EndpointUploadPrivateFile          = "upload_private_file"
	EndpointUploadGroupFile            = "upload_group_file"
	EndpointGetPrivateFileDownloadURL  = "get_private_file_download_url"
	EndpointGetGroupFileDownloadURL    = "get_group_file_download_url"
	EndpointGetGroupFiles              = "get_group_files"
	EndpointMoveGroupFile              = "move_group_file"
	EndpointRenameGroupFile            = "rename_group_file"
	EndpointDeleteGroupFile            = "delete_group_file"
	EndpointCreateGroupFolder          = "create_group_folder"
	EndpointRenameGroupFolder          = "rename_group_folder"
	EndpointDeleteGroupFolder          = "delete_group_folder"
)

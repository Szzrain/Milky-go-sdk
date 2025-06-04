package Milky_go_sdk

var (
	gateway     = ""
	EndpointAPI = func() string {
		return gateway + "/api/"
	}
	EndpointGetLoginInfo = func() string {
		return EndpointAPI() + "get_login_info"
	}
	EndpointGetFriendList = func() string {
		return EndpointAPI() + "get_friend_list"
	}
	EndpointGetFriendInfo = func() string {
		return EndpointAPI() + "get_friend_info"
	}
	EndpointGetGroupList = func() string {
		return EndpointAPI() + "get_group_list"
	}
	EndpointGetGroupInfo = func() string {
		return EndpointAPI() + "get_group_info"
	}
	EndpointGetGroupMemberList = func() string {
		return EndpointAPI() + "get_group_member_list"
	}
	EndpointGetGroupMemberInfo = func() string {
		return EndpointAPI() + "get_group_member_info"
	}
	EndpointSendPrivateMessage = func() string {
		return EndpointAPI() + "send_private_message"
	}
	EndpointSendGroupMessage = func() string {
		return EndpointAPI() + "send_group_message"
	}
	EndpointGetMessage = func() string {
		return EndpointAPI() + "get_message"
	}
	EndpointGetHistoryMessage = func() string {
		return EndpointAPI() + "get_history_message"
	}
	EndpointGetResourceTempURL = func() string {
		return EndpointAPI() + "get_resource_temp_url"
	}
	EndpointGetForwardedMessages = func() string {
		return EndpointAPI() + "get_forwarded_messages"
	}
	EndpointRecallPrivateMessage = func() string {
		return EndpointAPI() + "recall_private_message"
	}
	EndpointRecallGroupMessage = func() string {
		return EndpointAPI() + "recall_group_message"
	}
	EndpointSendFriendNudge = func() string {
		return EndpointAPI() + "send_friend_nudge"
	}
	EndpointGetSendProfileLike = func() string {
		return EndpointAPI() + "send_profile_like"
	}
	EndpointSetGroupName = func() string {
		return EndpointAPI() + "set_group_name"
	}
	EndpointSetGroupAvatar = func() string {
		return EndpointAPI() + "set_group_avatar"
	}
	EndpointSetGroupMemberCard = func() string {
		return EndpointAPI() + "set_group_member_card"
	}
	EndpointSetGroupMemberSpecialTitle = func() string {
		return EndpointAPI() + "set_group_member_special_title"
	}
	EndpointSetGroupAdmin = func() string {
		return EndpointAPI() + "set_group_admin"
	}
	EndpointSetGroupMemberMute = func() string {
		return EndpointAPI() + "set_group_member_mute"
	}
	EndpointSetGroupWholeMute = func() string {
		return EndpointAPI() + "set_group_whole_mute"
	}
	EndpointKickGroupMember = func() string {
		return EndpointAPI() + "kick_group_member"
	}
	EndpointGetGroupAnnouncementList = func() string {
		return EndpointAPI() + "get_group_announcement_list"
	}
	EndpointSendGroupAnnouncement = func() string {
		return EndpointAPI() + "send_group_announcement"
	}
	EndpointDeleteGroupAnnouncement = func() string {
		return EndpointAPI() + "delete_group_announcement"
	}
	EndpointQuitGroup = func() string {
		return EndpointAPI() + "quit_group"
	}
	EndpointSendGroupMessageReaction = func() string {
		return EndpointAPI() + "send_group_message_reaction"
	}
	EndpointSendGroupNudge = func() string {
		return EndpointAPI() + "send_group_nudge"
	}
	EndpointGetFriendRequests = func() string {
		return EndpointAPI() + "get_friend_requests"
	}
	EndpointGetGroupRequests = func() string {
		return EndpointAPI() + "get_group_requests"
	}
	EndpointGetGroupInvitations = func() string {
		return EndpointAPI() + "get_group_invitations"
	}
	EndpointAcceptRequest = func() string {
		return EndpointAPI() + "accept_request"
	}
	EndpointRejectRequest = func() string {
		return EndpointAPI() + "reject_request"
	}
	EndpointUploadPrivateFile = func() string {
		return EndpointAPI() + "upload_private_file"
	}
	EndpointUploadGroupFile = func() string {
		return EndpointAPI() + "upload_group_file"
	}
	EndpointGetPrivateFileDownloadURL = func() string {
		return EndpointAPI() + "get_private_file_download_url"
	}
	EndpointGetGroupFileDownloadURL = func() string {
		return EndpointAPI() + "get_group_file_download_url"
	}
	EndpointGetGroupFiles = func() string {
		return EndpointAPI() + "get_group_files"
	}
	EndpointMoveGroupFile = func() string {
		return EndpointAPI() + "move_group_file"
	}
	EndpointRenameGroupFile = func() string {
		return EndpointAPI() + "rename_group_file"
	}
	EndpointDeleteGroupFile = func() string {
		return EndpointAPI() + "delete_group_file"
	}
	EndpointCreateGroupFolder = func() string {
		return EndpointAPI() + "create_group_folder"
	}
	EndpointRenameGroupFolder = func() string {
		return EndpointAPI() + "rename_group_folder"
	}
	EndpointDeleteGroupFolder = func() string {
		return EndpointAPI() + "delete_group_folder"
	}
)

// SetGateway sets the gateway URL for the API endpoints.
func SetGateway(gw string) {
	gateway = gw
}

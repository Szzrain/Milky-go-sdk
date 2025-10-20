# Milky Go SDK

Ref: 

[DiscordGo](https://github.com/bwmarrin/discordgo)

[Milky Doc](https://milky.ntqqrev.org/)

## Feature List

### event transport

- [x] WebSocket
- [ ] ~~WebHook~~ (Not planned)

### api

#### system

- [x] /get_login_info
- [x] /get_impl_info
- [x] /get_user_profile
- [x] /get_friend_list
- [x] /get_friend_info
- [x] /get_group_list
- [x] /get_group_info
- [x] /get_group_member_list
- [x] /get_group_member_info
- [x] /get_cookies
- [x] /get_csrf_token

#### message

- [x] /send_private_message
- [x] /send_group_message
- [x] /get_message
- [x] /get_history_messages
- [x] /get_resource_temp_url
- [x] /get_forwarded_messages
- [x] /recall_private_message
- [x] /recall_group_message
- [x] /mark_message_as_read

#### friend

- [x] /send_friend_nudge
- [x] /send_profile_like
- [x] /get_friend_requests
- [x] /accept_friend_request
- [x] /reject_friend_request

#### group

- [x] /set_group_name
- [x] /set_group_avatar
- [x] /set_group_member_card
- [x] /set_group_member_special_title
- [x] /set_group_member_admin
- [x] /set_group_member_mute
- [x] /set_group_whole_mute
- [x] /kick_group_member
- [x] /get_group_announcement_list
- [x] /send_group_announcement
- [x] /delete_group_announcement
- [x] /get_group_essence_messages
- [x] /set_group_essence_message
- [x] /quit_group
- [x] /send_group_message_reaction
- [x] /send_group_nudge
- [x] /get_group_notifications
- [x] /accept_group_request
- [x] /reject_group_request
- [x] /accept_group_invitation
- [x] /reject_group_invitation

#### file

- [x] /upload_private_file
- [x] /upload_group_file
- [x] /get_private_file_download_url
- [x] /get_group_file_download_url
- [x] /get_group_files
- [x] /move_group_file
- [x] /rename_group_file
- [x] /delete_group_file
- [x] /create_group_folder
- [x] /rename_group_folder
- [x] /delete_group_folder

### event

- [x] bot_offline
- [x] message_receive
- [x] message_recall
- [x] friend_request
- [x] group_join_request
- [x] group_invited_join_request
- [x] group_invitation
- [x] friend_nudge
- [x] friend_file_upload
- [x] group_admin_change
- [x] group_essence_message_change
- [x] group_member_increase
- [x] group_member_decrease
- [x] group_name_change
- [x] group_message_reaction
- [x] group_mute
- [x] group_whole_mute
- [x] group_nudge
- [x] group_file_upload

### segment

#### incoming

- [x] text
- [x] mention
- [x] mention_all
- [x] face
- [x] reply
- [x] image
- [x] record
- [x] video
- [x] forward
- [x] market_face
- [x] light_app
- [x] xml

#### outgoing

- [x] text
- [x] mention
- [x] mention_all
- [x] face
- [x] reply
- [x] image
- [x] record
- [x] video
- [x] forward

#### forward

- [x] text
<!-- - [ ] mention -->
<!-- - [ ] mention_all -->
- [x] face
- [x] reply
- [x] image
- [x] record
- [x] video
- [x] forward
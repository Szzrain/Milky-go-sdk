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
- [ ] /get_resource_temp_url
- [ ] /get_forwarded_messages
- [ ] /recall_private_message
- [ ] /recall_group_message
- [ ] /mark_message_as_read

#### friend

- [x] /send_friend_nudge
- [x] /send_profile_like
- [ ] /get_friend_requests
- [x] /accept_friend_request
- [ ] /reject_friend_request

#### group

- [ ] /set_group_name
- [ ] /set_group_avatar
- [x] /set_group_member_card
- [ ] /set_group_member_special_title
- [ ] /set_group_member_admin
- [ ] /set_group_member_mute
- [ ] /set_group_whole_mute
- [ ] /kick_group_member
- [ ] /get_group_announcement_list
- [ ] /send_group_announcement
- [ ] /delete_group_announcement
- [ ] /get_group_essence_messages
- [ ] /set_group_essence_message
- [x] /quit_group
- [ ] /send_group_message_reaction
- [x] /send_group_nudge
- [ ] /get_group_notifications
- [ ] /accept_group_request
- [ ] /reject_group_request
- [x] /accept_group_invitation
- [ ] /reject_group_invitation

#### file

- [ ] /upload_private_file
- [x] /upload_group_file
- [ ] /get_private_file_download_url
- [ ] /get_group_file_download_url
- [ ] /get_group_files
- [ ] /move_group_file
- [ ] /rename_group_file
- [ ] /delete_group_file
- [ ] /create_group_folder
- [ ] /rename_group_folder
- [ ] /delete_group_folder

### event

- [x] bot_offline
- [x] message_receive
- [x] message_recall
- [x] friend_request
- [ ] group_join_request
- [ ] group_invited_join_request
- [x] group_invitation
- [x] friend_nudge
- [ ] friend_file_upload
- [ ] group_admin_change
- [ ] group_essence_message_change
- [x] group_member_increase
- [x] group_member_decrease
- [ ] group_name_change
- [x] group_message_reaction
- [x] group_mute
- [x] group_whole_mute
- [x] group_nudge
- [ ] group_file_upload

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
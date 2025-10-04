package Milky_go_sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	// Marshal defines function used to encode JSON payloads
	Marshal func(v interface{}) ([]byte, error) = json.Marshal
	// Unmarshal defines function used to decode JSON payloads
	Unmarshal func(src []byte, v interface{}) error = json.Unmarshal
)

// All error constants
var (
	ErrJSONUnmarshal = errors.New("json unmarshal")
	ErrUnauthorized  = errors.New("HTTP request was unauthorized")
)

// RESTError stores error information about a request with a bad response code.
// Message is not always present, there are cases where api calls can fail
// without returning a json message.
type RESTError struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte

	Message *APIErrorMessage // Message may be nil.
}

// newRestError returns a new REST API error.
func newRestError(req *http.Request, resp *http.Response, body []byte) *RESTError {
	restErr := &RESTError{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
	}

	// Attempt to decode the error and assume no message was provided if it fails
	var msg *APIErrorMessage
	err := Unmarshal(body, &msg)
	if err == nil {
		restErr.Message = msg
	}

	return restErr
}

// Error returns a Rest API Error with its status code and body.
func (r RESTError) Error() string {
	return "HTTP " + r.Response.Status + ", " + string(r.ResponseBody)
}

// RequestConfig is an HTTP request configuration.
type RequestConfig struct {
	Request        *http.Request
	MaxRestRetries int
	Client         *http.Client
}

// newRequestConfig returns a new HTTP request configuration based on parameters in Session.
func newRequestConfig(s *Session, req *http.Request) *RequestConfig {
	return &RequestConfig{
		MaxRestRetries: s.MaxRestRetries,
		Client:         s.Client,
		Request:        req,
	}
}

// RequestOption is a function which mutates request configuration.
// It can be supplied as an argument to any REST method.
type RequestOption func(cfg *RequestConfig)

// WithClient changes the HTTP client used for the request.
func WithClient(client *http.Client) RequestOption {
	return func(cfg *RequestConfig) {
		if client != nil {
			cfg.Client = client
		}
	}
}

// WithRestRetries changes maximum amount of retries if request fails.
func WithRestRetries(max int) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.MaxRestRetries = max
	}
}

// WithHeader sets a header in the request.
func WithHeader(key, value string) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.Request.Header.Set(key, value)
	}
}

// WithContext changes context of the request.
func WithContext(ctx context.Context) RequestOption {
	return func(cfg *RequestConfig) {
		cfg.Request = cfg.Request.WithContext(ctx)
	}
}

// Request makes a (GET/POST/...) Requests to REST API with JSON data.
func (s *Session) Request(method string, pathStr string, data interface{}, options ...RequestOption) (response []byte, err error) {
	var body []byte
	if data != nil {
		body, err = Marshal(data)
		if err != nil {
			return
		}
	}
	return s.RequestBase(method, s.apiEndpoints.Endpoint(pathStr), "application/json", body, 0, options...)
}

// RequestBase makes a request
func (s *Session) RequestBase(method, urlStr, contentType string, b []byte, sequence int, options ...RequestOption) (response []byte, err error) {

	s.Logger.Debugf("API REQUEST %8s :: %s\n", method, urlStr)
	s.Logger.Debugf("API REQUEST  PAYLOAD :: [%s]\n", string(b))

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(b))
	if err != nil {
		return
	}

	if s.Token != "" {
		req.Header.Set("Authorization", "Bearer "+s.Token)
	}

	if b != nil {
		req.Header.Set("Content-Type", contentType)
	}

	req.Header.Set("User-Agent", s.UserAgent)

	cfg := newRequestConfig(s, req)
	for _, opt := range options {
		opt(cfg)
	}
	req = cfg.Request

	for k, v := range req.Header {
		s.Logger.Debugf("API REQUEST   HEADER :: [%s] = %+v\n", k, v)
	}

	resp, err := cfg.Client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		err2 := resp.Body.Close()
		if err2 != nil {
			s.Logger.Debugf("error closing resp body")
		}
	}()

	response, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	s.Logger.Debugf("API RESPONSE  STATUS :: %s\n", resp.Status)
	for k, v := range resp.Header {
		s.Logger.Debugf("API RESPONSE  HEADER :: [%s] = %+v\n", k, v)
	}
	s.Logger.Debugf("API RESPONSE    BODY :: [%s]\n\n\n", response)

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusCreated:
	case http.StatusNoContent:
	case http.StatusBadGateway:
		// Retry sending request if possible
		if sequence < cfg.MaxRestRetries {

			s.Logger.Infof("%s Failed (%s), Retrying...", urlStr, resp.Status)
			response, err = s.RequestBase(method, urlStr, contentType, b, sequence+1, options...)
		} else {
			err = fmt.Errorf("exceeded Max retries HTTP %s, %s", resp.Status, response)
		}
	case http.StatusUnauthorized:
		s.Logger.Warnf(ErrUnauthorized.Error())
		err = ErrUnauthorized
		fallthrough
	default: // Error condition
		err = newRestError(req, resp, response)
	}

	return
}

func unmarshal(data []byte, v interface{}) error {
	err := Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrJSONUnmarshal, err)
	}

	return nil
}

func handleAPIResponse(request []byte, apiResponse *APIResponse, data interface{}) error {
	if err := unmarshal(request, apiResponse); err != nil {
		return err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		return fmt.Errorf("API call failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil && data != nil {
		return unmarshal(apiResponse.Data, data)
	}
	return nil
}

func (s *Session) GetLoginInfo() (*LoginInfo, error) {
	var apiResponse APIResponse
	var loginInfo LoginInfo
	request, err := s.Request("POST", EndpointGetLoginInfo, "{}")
	if err != nil {
		return nil, err
	}
	if err = handleAPIResponse(request, &apiResponse, &loginInfo); err != nil {
		return nil, err
	}
	return &loginInfo, nil
}

func (s *Session) GetImplInfo() (*ImplInfo, error) {
	request, err := s.Request("POST", EndpointGetImplInfo, "{}")
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var implInfo ImplInfo
	if err = handleAPIResponse(request, &apiResponse, &implInfo); err != nil {
		return nil, err
	}
	return &implInfo, nil
}

func (s *Session) GetFriendList(noCache bool) ([]Friend, error) {
	request, err := s.Request("POST", EndpointGetFriendList, map[string]interface{}{
		"no_cache": noCache,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var friendList struct {
		Friends []Friend `json:"friends"`
	}
	if err = handleAPIResponse(request, &apiResponse, &friendList); err != nil {
		return nil, err
	}
	return friendList.Friends, nil
}

func (s *Session) GetFriendInfo(userID int64, noCache bool) (*Friend, error) {
	request, err := s.Request("POST", EndpointGetFriendInfo, map[string]interface{}{
		"user_id":  userID,
		"no_cache": noCache,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var friendInfo struct {
		Friend Friend `json:"friend"`
	}
	if err = handleAPIResponse(request, &apiResponse, &friendInfo); err != nil {
		return nil, err
	}
	return &friendInfo.Friend, nil
}

func (s *Session) GetGroupList(noCache bool) ([]GroupInfo, error) {
	request, err := s.Request("POST", EndpointGetGroupList, map[string]interface{}{
		"no_cache": noCache,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var groupList struct {
		Groups []GroupInfo `json:"groups"`
	}
	if err = handleAPIResponse(request, &apiResponse, &groupList); err != nil {
		return nil, err
	}
	return groupList.Groups, nil
}

func (s *Session) GetGroupInfo(groupID int64, noCache bool) (*GroupInfo, error) {
	request, err := s.Request("POST", EndpointGetGroupInfo, map[string]interface{}{
		"group_id": groupID,
		"no_cache": noCache,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var groupInfo struct {
		Group GroupInfo `json:"group"`
	}
	if err = handleAPIResponse(request, &apiResponse, &groupInfo); err != nil {
		return nil, err
	}
	return &groupInfo.Group, nil
}

func (s *Session) GetGroupMemberList(groupID int64, noCache bool) ([]GroupMemberInfo, error) {
	request, err := s.Request("POST", EndpointGetGroupMemberList, map[string]interface{}{
		"group_id": groupID,
		"no_cache": noCache,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var memberList struct {
		Members []GroupMemberInfo `json:"members"`
	}
	if err = handleAPIResponse(request, &apiResponse, &memberList); err != nil {
		return nil, err
	}
	return memberList.Members, nil
}

func (s *Session) GetGroupMemberInfo(groupID, userID int64, noCache bool) (*GroupMemberInfo, error) {
	request, err := s.Request("POST", EndpointGetGroupMemberInfo, map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"no_cache": noCache,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var memberInfo struct {
		Member GroupMemberInfo `json:"member"`
	}
	if err = handleAPIResponse(request, &apiResponse, &memberInfo); err != nil {
		return nil, err
	}
	return &memberInfo.Member, nil
}

func (s *Session) SendGroupMessage(groupID int64, message *[]IMessageElement) (*MessageRet, error) {
	request, err := s.Request("POST", EndpointSendGroupMessage, map[string]interface{}{
		"group_id": groupID,
		"message":  message,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var messageRet MessageRet
	if err = handleAPIResponse(request, &apiResponse, &messageRet); err != nil {
		s.Logger.Errorf("Failed to unmarshal send group message response: %v", err)
		return nil, err
	}
	return &messageRet, nil
}

func (s *Session) SendPrivateMessage(userID int64, message *[]IMessageElement) (*MessageRet, error) {
	request, err := s.Request("POST", EndpointSendPrivateMessage, map[string]interface{}{
		"user_id": userID,
		"message": message,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var messageRet MessageRet
	if err = handleAPIResponse(request, &apiResponse, &messageRet); err != nil {
		s.Logger.Errorf("Failed to unmarshal send private message response: %v", err)
		return nil, err
	}
	return &messageRet, nil
}

func (s *Session) AcceptGroupInvitation(groupID int64, invitationSeq string) error {
	request, err := s.Request("POST", EndpointAcceptGroupInvitation, map[string]interface{}{
		"group_id":       groupID,
		"invitation_seq": invitationSeq,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) AcceptFriendRequest(initiatorUid string, isFiltered bool) error {
	request, err := s.Request("POST", EndpointAcceptFriendRequest, map[string]interface{}{
		"is_filtered":   isFiltered,
		"initiator_uid": initiatorUid,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) GetMessage(messageScene string, peerID int64, messageSeq int64) (*ReceiveMessage, error) {
	request, err := s.Request("POST", EndpointGetMessage, map[string]interface{}{
		"message_scene": messageScene,
		"peer_id":       peerID,
		"message_seq":   messageSeq,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var receiveMessage ReceiveMessage
	if err = handleAPIResponse(request, &apiResponse, &receiveMessage); err != nil {
		return nil, err
	}
	return &receiveMessage, nil
}

func (s *Session) GetHistoryMessages(messageScene string, peerID int64, startMessageSeq int64, direction string, limit int32) (*[]ReceiveMessage, error) {
	request, err := s.Request("POST", EndpointGetHistoryMessages, map[string]interface{}{
		"message_scene":     messageScene,
		"peer_id":           peerID,
		"start_message_seq": startMessageSeq,
		"direction":         direction,
		"limit":             limit,
	})
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var historyMessages []ReceiveMessage
	if err = handleAPIResponse(request, &apiResponse, &historyMessages); err != nil {
		return nil, err
	}
	return &historyMessages, nil
}

func (s *Session) SendFriendNudge(userID int64, isSelf bool) error {
	request, err := s.Request("POST", EndpointSendFriendNudge, map[string]interface{}{
		"user_id": userID,
		"is_self": isSelf,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) QuitGroup(groupID int64) error {
	request, err := s.Request("POST", EndpointQuitGroup, map[string]interface{}{
		"group_id": groupID,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) SendGroupNudge(groupID int64, userID int64) error {
	request, err := s.Request("POST", EndpointSendGroupNudge, map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) SetGroupName(groupID int64, name string) error {
	request, err := s.Request("POST", EndpointSetGroupName, map[string]interface{}{
		"group_id": groupID,
		"name":     name,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) SetGroupMemberCard(groupID int64, userID int64, card string) error {
	request, err := s.Request("POST", EndpointSetGroupMemberCard, map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"card":     card,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) SetGroupMemberSpecialTitle(groupID int64, userID int64, specialTitle string) error {
	request, err := s.Request("POST", EndpointSetGroupMemberSpecialTitle, map[string]interface{}{
		"group_id":      groupID,
		"user_id":       userID,
		"special_title": specialTitle,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) SetGroupAdmin(groupID int64, userID int64, isSet bool) error {
	request, err := s.Request("POST", EndpointSetGroupAdmin, map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"is_admin": isSet,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) SetGroupMemberMute(groupID int64, userID int64, duration int64) error {
	request, err := s.Request("POST", EndpointSetGroupMemberMute, map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"duration": duration,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) SetGroupWholeMute(groupID int64, isMute bool) error {
	request, err := s.Request("POST", EndpointSetGroupWholeMute, map[string]interface{}{
		"group_id": groupID,
		"is_mute":  isMute,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) KickGroupMember(groupID int64, userID int64, rejectAddRequest bool) error {
	request, err := s.Request("POST", EndpointKickGroupMember, map[string]interface{}{
		"group_id":           groupID,
		"user_id":            userID,
		"reject_add_request": rejectAddRequest,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

func (s *Session) UploadGroupFile(groupID int64, fileURI string, fileName string, parentFolderID string) (string, error) {
	request, err := s.Request("POST", EndpointUploadGroupFile, map[string]interface{}{
		"group_id":         groupID,
		"file_uri":         fileURI,
		"file_name":        fileName,
		"parent_folder_id": parentFolderID,
	})
	if err != nil {
		return "", err
	}
	var apiResponse APIResponse
	var uploadFileResponse struct {
		FileID string `json:"file_id"`
	}
	if err = handleAPIResponse(request, &apiResponse, &uploadFileResponse); err != nil {
		return "", err
	}
	return uploadFileResponse.FileID, nil
}

func (s *Session) GetGroupFileDownloadURL(groupID int64, fileID string) (string, error) {
	request, err := s.Request("POST", EndpointGetGroupFileDownloadURL, map[string]interface{}{
		"group_id": groupID,
		"file_id":  fileID,
	})
	if err != nil {
		return "", err
	}
	var apiResponse APIResponse
	var downloadURLResponse struct {
		DownloadURL string `json:"download_url"`
	}
	if err = handleAPIResponse(request, &apiResponse, &downloadURLResponse); err != nil {
		return "", err
	}
	return downloadURLResponse.DownloadURL, nil
}

func (s *Session) DeleteGroupFile(groupID int64, fileID string) error {
	request, err := s.Request("POST", EndpointDeleteGroupFile, map[string]interface{}{
		"group_id": groupID,
		"file_id":  fileID,
	})
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	return handleAPIResponse(request, &apiResponse, nil)
}

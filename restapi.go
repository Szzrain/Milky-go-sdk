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

func (s *Session) GetLoginInfo() (*LoginInfo, error) {
	request, err := s.Request("POST", EndpointGetLoginInfo, "{}", WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var loginInfo LoginInfo
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal login info: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get login info failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get login info failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &loginInfo); err != nil {
			s.Logger.Errorf("Failed to unmarshal login info data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Login info data is nil")
		return nil, fmt.Errorf("login info data is nil")
	}
	return &loginInfo, nil
}

func (s *Session) GetImplInfo() (*ImplInfo, error) {
	request, err := s.Request("POST", EndpointGetImplInfo, "{}", WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var implInfo ImplInfo
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal impl info: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get impl info failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get impl info failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &implInfo); err != nil {
			s.Logger.Errorf("Failed to unmarshal impl info data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Impl info data is nil")
		return nil, fmt.Errorf("impl info data is nil")
	}
	return &implInfo, nil
}

func (s *Session) GetFriendList(noCache bool) (*[]Friend, error) {
	request, err := s.Request("POST", EndpointGetFriendList, map[string]interface{}{
		"no_cache": noCache,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var friendList []Friend
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal friend list: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get friend list failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get friend list failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &friendList); err != nil {
			s.Logger.Errorf("Failed to unmarshal friend list data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Friend list data is nil")
		return nil, fmt.Errorf("friend list data is nil")
	}
	return &friendList, nil
}

func (s *Session) GetFriendInfo(userID int64, noCache bool) (*Friend, error) {
	request, err := s.Request("POST", EndpointGetFriendInfo, map[string]interface{}{
		"user_id":  userID,
		"no_cache": noCache,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var friendInfo Friend
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal friend info: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get friend info failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get friend info failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &friendInfo); err != nil {
			s.Logger.Errorf("Failed to unmarshal friend info data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Friend info data is nil")
		return nil, fmt.Errorf("friend info data is nil")
	}
	return &friendInfo, nil
}

func (s *Session) GetGroupList(noCache bool) (*[]GroupInfo, error) {
	request, err := s.Request("POST", EndpointGetGroupList, map[string]interface{}{
		"no_cache": noCache,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var groupList []GroupInfo
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal group list: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get group list failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get group list failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &groupList); err != nil {
			s.Logger.Errorf("Failed to unmarshal group list data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Group list data is nil")
		return nil, fmt.Errorf("group list data is nil")
	}
	return &groupList, nil
}

func (s *Session) GetGroupInfo(groupID int64, noCache bool) (*GroupInfo, error) {
	request, err := s.Request("POST", EndpointGetGroupInfo, map[string]interface{}{
		"group_id": groupID,
		"no_cache": noCache,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var groupInfo GroupInfo
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal group info: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get group info failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get group info failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &groupInfo); err != nil {
			s.Logger.Errorf("Failed to unmarshal group info data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Group info data is nil")
		return nil, fmt.Errorf("group info data is nil")
	}
	return &groupInfo, nil
}

func (s *Session) GetGroupMemberList(groupID int64, noCache bool) (*[]GroupMemberInfo, error) {
	request, err := s.Request("POST", EndpointGetGroupMemberList, map[string]interface{}{
		"group_id": groupID,
		"no_cache": noCache,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var memberList []GroupMemberInfo
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal group member list: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get group member list failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get group member list failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &memberList); err != nil {
			s.Logger.Errorf("Failed to unmarshal group member list data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Group member list data is nil")
		return nil, fmt.Errorf("group member list data is nil")
	}
	return &memberList, nil
}

func (s *Session) GetGroupMemberInfo(groupID, userID int64, noCache bool) (*GroupMemberInfo, error) {
	request, err := s.Request("POST", EndpointGetGroupMemberInfo, map[string]interface{}{
		"group_id": groupID,
		"user_id":  userID,
		"no_cache": noCache,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var memberInfo GroupMemberInfo
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal group member info: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get group member info failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get group member info failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &memberInfo); err != nil {
			s.Logger.Errorf("Failed to unmarshal group member info data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Group member info data is nil")
		return nil, fmt.Errorf("group member info data is nil")
	}
	return &memberInfo, nil
}

func (s *Session) SendGroupMessage(groupID int64, message *[]IMessageElement) (*MessageRet, error) {
	request, err := s.Request("POST", EndpointSendGroupMessage, map[string]interface{}{
		"group_id": groupID,
		"message":  message,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal send group message response: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Send group message failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("send group message failed: %s", apiResponse.Message)
	}
	var messageRet MessageRet
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &messageRet); err != nil {
			s.Logger.Errorf("Failed to unmarshal message ret data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Send group message data is nil")
		return nil, fmt.Errorf("send group message data is nil")
	}
	return &messageRet, nil
}

func (s *Session) SendPrivateMessage(userID int64, message *[]IMessageElement) (*MessageRet, error) {
	request, err := s.Request("POST", EndpointSendPrivateMessage, map[string]interface{}{
		"user_id": userID,
		"message": message,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal send private message response: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Send private message failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("send private message failed: %s", apiResponse.Message)
	}
	var messageRet MessageRet
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &messageRet); err != nil {
			s.Logger.Errorf("Failed to unmarshal message ret data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Send private message data is nil")
		return nil, fmt.Errorf("send private message data is nil")
	}
	return &messageRet, nil
}

func (s *Session) GetMessage(messageScene string, peerID int64, messageSeq int64) (*ReceiveMessage, error) {
	request, err := s.Request("POST", EndpointGetMessage, map[string]interface{}{
		"message_scene": messageScene,
		"peer_id":       peerID,
		"message_seq":   messageSeq,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var receiveMessage ReceiveMessage
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal get message response: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get message failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get message failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &receiveMessage); err != nil {
			s.Logger.Errorf("Failed to unmarshal receive message data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Get message data is nil")
		return nil, fmt.Errorf("get message data is nil")
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
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return nil, err
	}
	var apiResponse APIResponse
	var historyMessages []ReceiveMessage
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal get history message response: %v", err)
		return nil, err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get history message failed: %s", apiResponse.Message)
		return nil, fmt.Errorf("get history message failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &historyMessages); err != nil {
			s.Logger.Errorf("Failed to unmarshal history messages data: %v", err)
			return nil, err
		}
	} else {
		s.Logger.Errorf("Get history message data is nil")
		return nil, fmt.Errorf("get history message data is nil")
	}
	return &historyMessages, nil
}

func (s *Session) SendFriendNudge(userID int64, isSelf bool) error {
	request, err := s.Request("POST", EndpointSendFriendNudge, map[string]interface{}{
		"user_id": userID,
		"is_self": isSelf,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal send friend nudge response: %v", err)
		return err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Send friend nudge failed: %s", apiResponse.Message)
		return fmt.Errorf("send friend nudge failed: %s", apiResponse.Message)
	}
	return nil
}

func (s *Session) QuitGroup(groupID int64) error {
	request, err := s.Request("POST", EndpointQuitGroup, map[string]interface{}{
		"group_id": groupID,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return err
	}
	var apiResponse APIResponse
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal quit group response: %v", err)
		return err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Quit group failed: %s", apiResponse.Message)
		return fmt.Errorf("quit group failed: %s", apiResponse.Message)
	}
	return nil
}

func (s *Session) UploadGroupFile(groupID int64, fileURI string, fileName string, parentFolderID string) (string, error) {
	request, err := s.Request("POST", EndpointUploadGroupFile, map[string]interface{}{
		"group_id":         groupID,
		"file_uri":         fileURI,
		"file_name":        fileName,
		"parent_folder_id": parentFolderID,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return "", err
	}
	var apiResponse APIResponse
	var uploadFileResponse struct {
		FileID string `json:"file_id"`
	}
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal upload group file response: %v", err)
		return "", err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Upload group file failed: %s", apiResponse.Message)
		return "", fmt.Errorf("upload group file failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &uploadFileResponse); err != nil {
			s.Logger.Errorf("Failed to unmarshal upload file response data: %v", err)
			return "", err
		}
	} else {
		s.Logger.Errorf("Upload group file data is nil")
		return "", fmt.Errorf("upload group file data is nil")
	}
	return uploadFileResponse.FileID, nil
}

func (s *Session) GetGroupFileDownloadURL(groupID int64, fileID string) (string, error) {
	request, err := s.Request("POST", EndpointGetGroupFileDownloadURL, map[string]interface{}{
		"group_id": groupID,
		"file_id":  fileID,
	}, WithHeader("Content-Type", "application/json"))
	if err != nil {
		return "", err
	}
	var apiResponse APIResponse
	var downloadURLResponse struct {
		DownloadURL string `json:"download_url"`
	}
	if err = unmarshal(request, &apiResponse); err != nil {
		s.Logger.Errorf("Failed to unmarshal get group file download URL response: %v", err)
		return "", err
	}
	if apiResponse.RetCode != 0 || apiResponse.Status != "ok" {
		s.Logger.Errorf("Get group file download URL failed: %s", apiResponse.Message)
		return "", fmt.Errorf("get group file download URL failed: %s", apiResponse.Message)
	}
	if apiResponse.Data != nil {
		if err = unmarshal(apiResponse.Data, &downloadURLResponse); err != nil {
			s.Logger.Errorf("Failed to unmarshal download URL response data: %v", err)
			return "", err
		}
	} else {
		s.Logger.Errorf("Get group file download URL data is nil")
		return "", fmt.Errorf("get group file download URL data is nil")
	}
	return downloadURLResponse.DownloadURL, nil
}

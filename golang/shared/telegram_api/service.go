package telegram_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

const (
	viewCountEndpoint    = "post_view_count"
	membersCountEndpoint = "subscribers"
)

type Service struct {
	url string
}

func New() *Service {
	return &Service{
		url: "http://127.0.0.1:5000",
	}
}

type RequestData struct {
	ChannelHandle string `json:"channel_handle"`
	MessageID     int    `json:"message_id"`
}

type GetMessageViewsResponse struct {
	ViewCount int    `json:"view_count"`
	Error     string `json:"error"`
}

type GetMembersCountResponse struct {
	SubscribersCount int    `json:"subscribers_count"`
	Error            string `json:"error"`
}

func (s *Service) GetMessageViews(chatHandle string, messageID int) (int, error) {
	// Marshal the request data to JSON
	requestBody, err := json.Marshal(RequestData{
		ChannelHandle: fmt.Sprintf("@%s", chatHandle),
		MessageID:     messageID,
	})
	if err != nil {
		fmt.Println("Error marshalling request data:", err)
		return 0, err
	}

	resp, err := http.Post(
		fmt.Sprintf("%s/%s", s.url, viewCountEndpoint),
		"application/json", bytes.NewBuffer(requestBody),
	)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return 0, err
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return 0, err
	}

	var responseData GetMessageViewsResponse
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return 0, err
	}

	if responseData.Error != "" {
		fmt.Println("Error:", responseData.Error)
		return 0, errors.New(responseData.Error)
	}

	return responseData.ViewCount, nil
}

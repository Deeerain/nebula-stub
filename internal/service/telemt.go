package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type TelemtService struct {
	baseUrl string
	client  *http.Client
}

func NewTelemtService(host string, port int) *TelemtService {
	return &TelemtService{
		baseUrl: fmt.Sprintf("http://%s:%v", host, port),
		client:  http.DefaultClient,
	}
}

func (s *TelemtService) GetConnections() (map[string]any, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/users", s.baseUrl), nil)
	if err != nil {
		log.Fatalln("failed to create request: ", err.Error())
	}
	req.Header.Add("Content-type", "application/json")

	response, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error: %w", err)
	}

	var data map[string]any
	if err = json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode body: %w", err)
	}

	switch response.StatusCode {
	case 200:
		return data, nil
	default:
		return nil, fmt.Errorf("Telemt error: %v", data["error"])
	}
}

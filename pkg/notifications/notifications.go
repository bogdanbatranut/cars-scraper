package notifications

import (
	"bytes"
	"carscraper/pkg/amconfig"
	"fmt"
	"net/http"
)

type httpService struct {
	client *http.Client
}

func NewNotificationsService(cfg *amconfig.IConfig) *NotificationsService {
	urlStr := "http://dev.auto-mall.ro:88/automall"

	return &NotificationsService{
		httpService: httpService{
			client: &http.Client{},
		},
		baseURL: urlStr,
	}
}

type NotificationsService struct {
	httpService httpService
	baseURL     string
}

func (s *NotificationsService) SendTestNofication() {
}

func (s *NotificationsService) SendOpenAdNotification(adID uint, payload string) error {
	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer([]byte(payload)))

	if err != nil {
		return err
	}
	headerValue := fmt.Sprintf("view, Open ad, http://dev.auto-mall.ro/ad/%d, clear=true", adID)

	req.Header.Set(
		"Actions", headerValue)

	req.Header.Set("Content-Type", "text/plain")

	resp, err := s.httpService.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification, status code: %d", resp.StatusCode)
	}

	return nil
}

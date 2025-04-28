package notifications

import (
	"bytes"
	"carscraper/pkg/adsdb"
	"carscraper/pkg/amconfig"
	"fmt"
	"net/http"
)

type httpService struct {
	client *http.Client
}

func NewNotificationsService(cfg amconfig.IConfig) *NotificationsService {
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

func (s *NotificationsService) SendMinPriceInCriteria(ad adsdb.Ad) error {
	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer([]byte(*ad.Title)))

	if err != nil {
		return err
	}
	headerValue := fmt.Sprintf("view, Open ad, http://dev.auto-mall.ro/ad/%d, clear=true", ad.ID)

	req.Header.Set("Priority", "urgent")
	req.Header.Set("Tags", "bangbang")
	req.Header.Set("Title", "New MINIMUM PRICE in CRITERIA")
	req.Header.Set("Actions", headerValue)

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

func (s *NotificationsService) SendNewMinPrice(ad adsdb.Ad) error {
	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer([]byte(*ad.Title)))

	if err != nil {
		return err
	}
	headerValue := fmt.Sprintf("view, Open ad, http://dev.auto-mall.ro/ad/%d, clear=true", ad.ID)

	req.Header.Set("Priority", "urgent")
	req.Header.Set("Tags", "bangbang")
	req.Header.Set("Title", "New MINIMUM PRICE for car")
	req.Header.Set("Actions", headerValue)

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

func (s *NotificationsService) SendMinPriceCreatedNotification(ad adsdb.Ad) error {
	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer([]byte(*ad.Title)))

	if err != nil {
		return err
	}
	headerValue := fmt.Sprintf("view, Open ad, http://dev.auto-mall.ro/ad/%d, clear=true", ad.ID)

	req.Header.Set("Priority", "urgent")
	req.Header.Set("Tags", "warning")
	req.Header.Set("Title", "New car with MINIMUM PRICE")
	req.Header.Set("Actions", headerValue)

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

func (s *NotificationsService) SendPriceDecreaseNotification(ad adsdb.Ad) error {

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer([]byte(*ad.Title)))

	if err != nil {
		return err
	}
	headerValue := fmt.Sprintf("view, Open ad, http://dev.auto-mall.ro/ad/%d, clear=true", ad.ID)

	req.Header.Set("Priority", "urgent")
	req.Header.Set("Tags", "warning, small_red_triangle_down")
	req.Header.Set("Title", "PRICE DECREASE")
	req.Header.Set("Actions", headerValue)

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

func (s *NotificationsService) SendPriceIncreaseNotification(ad adsdb.Ad) error {

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer([]byte(*ad.Title)))

	if err != nil {
		return err
	}
	headerValue := fmt.Sprintf("view, Open ad, http://dev.auto-mall.ro/ad/%d, clear=true", ad.ID)

	req.Header.Set("Priority", "urgent")
	req.Header.Set("Tags", "warning, small_red_triangle")
	req.Header.Set("Title", "PRICE INCREASE")
	req.Header.Set("Actions", headerValue)

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

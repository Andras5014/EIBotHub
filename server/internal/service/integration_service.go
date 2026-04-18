package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

const (
	WebhookEventTest                 = "webhook.test"
	WebhookEventSkillCreated         = "skill.created"
	WebhookEventDiscussionCreated    = "discussion.created"
	WebhookEventWikiUpdated          = "wiki.updated"
	WebhookEventVerificationReviewed = "verification.reviewed"
	WebhookEventRewardRedeemed       = "reward.redeemed"
)

var supportedWebhookEvents = map[string]struct{}{
	WebhookEventTest:                 {},
	WebhookEventSkillCreated:         {},
	WebhookEventDiscussionCreated:    {},
	WebhookEventWikiUpdated:          {},
	WebhookEventVerificationReviewed: {},
	WebhookEventRewardRedeemed:       {},
}

type IntegrationService struct {
	repo   *repository.IntegrationRepository
	client *http.Client
}

func NewIntegrationService(repo *repository.IntegrationRepository) *IntegrationService {
	return &IntegrationService{
		repo: repo,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

func (s *IntegrationService) Spec() *dto.OpenAPISpecResponse {
	return &dto.OpenAPISpecResponse{
		Version:  "v1",
		BaseURL:  "/api/v1",
		Overview: "开放社区为模型、数据集、技能、Wiki 与社区协作提供统一 API，并支持 Webhook 事件回调。",
		Endpoints: []dto.OpenAPIEndpointItem{
			{Method: http.MethodPost, Path: "/api/v1/auth/login", Summary: "登录并获取访问令牌", Auth: "公开"},
			{Method: http.MethodGet, Path: "/api/v1/models", Summary: "查询模型列表", Auth: "公开"},
			{Method: http.MethodGet, Path: "/api/v1/datasets", Summary: "查询数据集列表", Auth: "公开"},
			{Method: http.MethodGet, Path: "/api/v1/wiki/pages", Summary: "查询 Wiki 词条", Auth: "公开"},
			{Method: http.MethodPost, Path: "/api/v1/skills", Summary: "发布技能", Auth: "Bearer Token"},
			{Method: http.MethodPost, Path: "/api/v1/discussions", Summary: "发起讨论", Auth: "Bearer Token"},
			{Method: http.MethodGet, Path: "/api/v1/webhooks", Summary: "列出当前用户的 Webhook 订阅", Auth: "Bearer Token"},
			{Method: http.MethodPost, Path: "/api/v1/webhooks", Summary: "创建 Webhook 订阅", Auth: "Bearer Token"},
			{Method: http.MethodPost, Path: "/api/v1/webhooks/{id}/test", Summary: "发送测试回调", Auth: "Bearer Token"},
		},
		WebhookEvents: supportedEventList(),
		CurlExamples: []dto.CLIExampleItem{
			{
				Title: "登录获取令牌",
				Command: `curl -X POST http://127.0.0.1:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@example.com","password":"Demo123!"}'`,
			},
			{
				Title: "创建 Webhook",
				Command: `curl -X POST http://127.0.0.1:8080/api/v1/webhooks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"本地调试","target_url":"http://127.0.0.1:9000/hook","secret":"opencommunity-secret","events":["webhook.test","reward.redeemed"]}'`,
			},
			{
				Title: "测试 Webhook",
				Command: `curl -X POST http://127.0.0.1:8080/api/v1/webhooks/1/test \
  -H "Authorization: Bearer <token>"`,
			},
		},
	}
}

func (s *IntegrationService) ListWebhooks(userID uint) ([]dto.WebhookSubscriptionItem, error) {
	items, err := s.repo.ListWebhooksByUser(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.WebhookSubscriptionItem, 0, len(items))
	for _, item := range items {
		result = append(result, toWebhookItem(item))
	}
	return result, nil
}

func (s *IntegrationService) CreateWebhook(userID uint, input dto.WebhookSubscriptionRequest) (*dto.WebhookSubscriptionItem, error) {
	if _, err := url.ParseRequestURI(input.TargetURL); err != nil {
		return nil, support.NewError(http.StatusBadRequest, "invalid_target_url", "Webhook 地址格式不正确")
	}
	events, err := normalizeWebhookEvents(input.Events)
	if err != nil {
		return nil, err
	}
	item := &model.WebhookSubscription{
		UserID:    userID,
		Name:      input.Name,
		TargetURL: input.TargetURL,
		Secret:    input.Secret,
		Events:    support.JoinCSV(events),
	}
	if err := s.repo.CreateWebhook(item); err != nil {
		return nil, err
	}
	result := toWebhookItem(*item)
	return &result, nil
}

func (s *IntegrationService) Deliveries(userID, webhookID uint) ([]dto.WebhookDeliveryItem, error) {
	if _, err := s.repo.GetWebhookByIDAndUser(webhookID, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "webhook_not_found", "Webhook 订阅不存在")
		}
		return nil, err
	}
	items, err := s.repo.ListDeliveries(webhookID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.WebhookDeliveryItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.WebhookDeliveryItem{
			ID:           item.ID,
			Event:        item.Event,
			Status:       item.Status,
			ResponseCode: item.ResponseCode,
			ResponseBody: item.ResponseBody,
			CreatedAt:    item.CreatedAt,
		})
	}
	return result, nil
}

func (s *IntegrationService) Test(userID, webhookID uint) (*dto.WebhookDeliveryItem, error) {
	item, err := s.repo.GetWebhookByIDAndUser(webhookID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "webhook_not_found", "Webhook 订阅不存在")
		}
		return nil, err
	}
	delivery, err := s.deliver(item, WebhookEventTest, map[string]any{
		"user_id": userID,
		"message": "这是开放社区发出的测试事件。",
	})
	if err != nil {
		return nil, err
	}
	result := dto.WebhookDeliveryItem{
		ID:           delivery.ID,
		Event:        delivery.Event,
		Status:       delivery.Status,
		ResponseCode: delivery.ResponseCode,
		ResponseBody: delivery.ResponseBody,
		CreatedAt:    delivery.CreatedAt,
	}
	return &result, nil
}

func (s *IntegrationService) Emit(userID uint, event string, payload any) {
	if s == nil || s.repo == nil || userID == 0 {
		return
	}
	items, err := s.repo.ListWebhooksByUser(userID)
	if err != nil {
		return
	}
	for i := range items {
		if !wantsWebhookEvent(items[i].Events, event) {
			continue
		}
		_, _ = s.deliver(&items[i], event, payload)
	}
}

func (s *IntegrationService) deliver(item *model.WebhookSubscription, event string, payload any) (*model.WebhookDelivery, error) {
	body := map[string]any{
		"event":       event,
		"occurred_at": time.Now(),
		"data":        payload,
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, item.TargetURL, bytes.NewReader(raw))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-OpenCommunity-Event", event)
	req.Header.Set("X-OpenCommunity-Signature", signWebhookPayload(item.Secret, raw))

	delivery := &model.WebhookDelivery{
		WebhookID:   item.ID,
		Event:       event,
		Status:      "failed",
		RequestBody: string(raw),
	}

	now := time.Now()
	item.LastDeliveredAt = &now

	resp, err := s.client.Do(req)
	if err != nil {
		delivery.ResponseBody = trimWebhookBody(err.Error())
		item.LastStatus = "failed"
		item.LastError = trimWebhookBody(err.Error())
		item.LastResponseCode = 0
		if saveErr := s.persistDelivery(item, delivery); saveErr != nil {
			return nil, saveErr
		}
		return delivery, nil
	}
	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
	delivery.ResponseCode = resp.StatusCode
	delivery.ResponseBody = trimWebhookBody(string(responseBody))
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		delivery.Status = "success"
		item.LastStatus = "success"
		item.LastError = ""
	} else {
		item.LastStatus = "failed"
		item.LastError = trimWebhookBody(resp.Status)
	}
	item.LastResponseCode = resp.StatusCode
	if saveErr := s.persistDelivery(item, delivery); saveErr != nil {
		return nil, saveErr
	}
	return delivery, nil
}

func (s *IntegrationService) persistDelivery(item *model.WebhookSubscription, delivery *model.WebhookDelivery) error {
	if err := s.repo.AddDelivery(delivery); err != nil {
		return err
	}
	return s.repo.UpdateWebhook(item)
}

func toWebhookItem(item model.WebhookSubscription) dto.WebhookSubscriptionItem {
	return dto.WebhookSubscriptionItem{
		ID:               item.ID,
		Name:             item.Name,
		TargetURL:        item.TargetURL,
		Events:           support.SplitCSV(item.Events),
		LastStatus:       item.LastStatus,
		LastResponseCode: item.LastResponseCode,
		LastError:        item.LastError,
		LastDeliveredAt:  item.LastDeliveredAt,
		CreatedAt:        item.CreatedAt,
	}
}

func normalizeWebhookEvents(events []string) ([]string, error) {
	seen := make(map[string]struct{}, len(events))
	result := make([]string, 0, len(events))
	for _, value := range events {
		event := strings.TrimSpace(value)
		if event == "" {
			continue
		}
		if _, ok := supportedWebhookEvents[event]; !ok {
			return nil, support.NewError(http.StatusBadRequest, "invalid_webhook_event", "存在不支持的 Webhook 事件")
		}
		if _, ok := seen[event]; ok {
			continue
		}
		seen[event] = struct{}{}
		result = append(result, event)
	}
	if len(result) == 0 {
		return nil, support.NewError(http.StatusBadRequest, "invalid_webhook_event", "至少选择一个 Webhook 事件")
	}
	sort.Strings(result)
	return result, nil
}

func wantsWebhookEvent(rawEvents string, event string) bool {
	for _, item := range support.SplitCSV(rawEvents) {
		if item == event {
			return true
		}
	}
	return false
}

func supportedEventList() []string {
	items := make([]string, 0, len(supportedWebhookEvents))
	for event := range supportedWebhookEvents {
		items = append(items, event)
	}
	sort.Strings(items)
	return items
}

func signWebhookPayload(secret string, payload []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func trimWebhookBody(input string) string {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) > 512 {
		return trimmed[:512]
	}
	return trimmed
}

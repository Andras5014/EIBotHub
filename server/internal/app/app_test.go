package app

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/config"
	"github.com/Andras5014/EIBotHub/server/internal/model"
)

func TestAuthProfileFlow(t *testing.T) {
	app := newTestApp(t)

	registerBody := map[string]any{
		"username": "alice",
		"email":    "alice@example.com",
		"password": "Alice123!",
	}
	resp := performJSON(t, app, http.MethodPost, "/api/v1/auth/register", registerBody, "")
	require.Equal(t, http.StatusCreated, resp.Code)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	token := payload["data"].(map[string]any)["token"].(string)

	profileResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/profile", nil, token, "")
	require.Equal(t, http.StatusOK, profileResp.Code)

	updateBody := map[string]any{
		"username": "alice-updated",
		"bio":      "builds embodied AI workflows",
		"avatar":   "https://example.com/avatar.png",
	}
	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/users/me/profile", updateBody, token)
	require.Equal(t, http.StatusOK, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "alice-updated")
}

func TestHomeDashboardSeedData(t *testing.T) {
	app := newTestApp(t)

	homeResp := performRequest(t, app, http.MethodGet, "/api/v1/portal/home", nil, "", "")
	require.Equal(t, http.StatusOK, homeResp.Code)
	require.Contains(t, homeResp.Body.String(), "\"module_settings\"")
	require.Contains(t, homeResp.Body.String(), "\"hero_config\"")
	require.Contains(t, homeResp.Body.String(), "\"rankings_config\"")
	require.Contains(t, homeResp.Body.String(), "\"badge_label\"")
	require.Contains(t, homeResp.Body.String(), "Patrol Corridor Vision")
	require.Contains(t, homeResp.Body.String(), "Factory Safety Alert Set")
	require.Contains(t, homeResp.Body.String(), "夜间安防巡检模板")
	require.Contains(t, homeResp.Body.String(), "园区夜间安防巡逻案例")

	wikiResp := performRequest(t, app, http.MethodGet, "/api/v1/wiki/pages", nil, "", "")
	require.Equal(t, http.StatusOK, wikiResp.Code)
	require.Contains(t, wikiResp.Body.String(), "搬运任务模板参数说明")

	rankingResp := performRequest(t, app, http.MethodGet, "/api/v1/rankings/contributors", nil, "", "")
	require.Equal(t, http.StatusOK, rankingResp.Code)
	require.Contains(t, rankingResp.Body.String(), "integrator")
	require.Contains(t, rankingResp.Body.String(), "ops")
}

func TestSeedDisplayDataCoverage(t *testing.T) {
	app := newTestApp(t)

	var demo model.User
	require.NoError(t, app.DB.Where("email = ?", "demo@example.com").First(&demo).Error)

	assertCountAtLeast(t, "scene pages", app.DB.Model(&model.ScenePageConfig{}).Where("enabled = ?", true), 5)
	assertCountAtLeast(t, "announcements", app.DB.Model(&model.Announcement{}), 5)
	assertCountAtLeast(t, "published models", app.DB.Model(&model.ModelAsset{}).Where("status = ?", model.StatusPublished), 5)
	assertCountAtLeast(t, "published datasets", app.DB.Model(&model.Dataset{}).Where("status = ?", model.StatusPublished), 5)
	assertCountAtLeast(t, "published templates", app.DB.Model(&model.TaskTemplate{}).Where("status = ?", model.StatusPublished), 5)
	assertCountAtLeast(t, "published cases", app.DB.Model(&model.ApplicationCase{}).Where("status = ?", model.StatusPublished), 5)
	assertCountAtLeast(t, "published docs", app.DB.Model(&model.Document{}).Where("status = ?", model.StatusPublished), 5)
	assertCountAtLeast(t, "faqs", app.DB.Model(&model.FAQ{}), 5)
	assertCountAtLeast(t, "videos", app.DB.Model(&model.VideoTutorial{}).Where("active = ?", true), 5)
	assertCountAtLeast(t, "wiki pages", app.DB.Model(&model.WikiPage{}).Where("status = ?", model.StatusPublished), 5)
	assertCountAtLeast(t, "hot keywords", app.DB.Model(&model.SearchKeywordConfig{}).Where("keyword_type = ? AND enabled = ?", "hot", true), 5)
	assertCountAtLeast(t, "recommended keywords", app.DB.Model(&model.SearchKeywordConfig{}).Where("keyword_type = ? AND enabled = ?", "recommended", true), 5)
	assertCountAtLeast(t, "community comments", app.DB.Model(&model.ResourceComment{}), 5)
	assertCountAtLeast(t, "demo skills", app.DB.Model(&model.Skill{}).Where("owner_id = ? AND status = ?", demo.ID, model.StatusPublished), 5)
	assertCountAtLeast(t, "demo discussions", app.DB.Model(&model.Discussion{}).Where("user_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "demo follows", app.DB.Model(&model.Follow{}).Where("follower_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "demo followers", app.DB.Model(&model.Follow{}).Where("followed_user_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "demo favorites", app.DB.Model(&model.Favorite{}).Where("user_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "demo downloads", app.DB.Model(&model.DownloadRecord{}).Where("user_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "demo notifications", app.DB.Model(&model.Notification{}).Where("user_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "reward benefits", app.DB.Model(&model.RewardBenefit{}).Where("active = ?", true), 5)
	assertCountAtLeast(t, "demo redemptions", app.DB.Model(&model.RewardRedemption{}).Where("user_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "demo access requests", app.DB.Model(&model.DatasetAccessRequest{}).Where("user_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "demo conversations", app.DB.Model(&model.Conversation{}).
		Joins("join conversation_participants on conversation_participants.conversation_id = conversations.id").
		Where("conversation_participants.user_id = ?", demo.ID), 5)
	assertCountAtLeast(t, "demo workspaces", app.DB.Model(&model.Workspace{}).
		Joins("join workspace_members on workspace_members.workspace_id = workspaces.id").
		Where("workspace_members.user_id = ?", demo.ID), 5)
}

func TestPortalModuleAndHighlightAdminFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	modulesResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/portal/modules", nil, adminToken, "")
	require.Equal(t, http.StatusOK, modulesResp.Code)
	require.Contains(t, modulesResp.Body.String(), "\"sort_order\":10")

	updateModuleResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/portal/modules/community", map[string]any{
		"enabled":    true,
		"sort_order": 15,
	}, adminToken)
	require.Equal(t, http.StatusOK, updateModuleResp.Code)

	createHighlightResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/portal/highlights", map[string]any{
		"text":       "支持首页亮点由后台灵活维护",
		"sort_order": 5,
		"enabled":    true,
	}, adminToken)
	require.Equal(t, http.StatusCreated, createHighlightResp.Code)
	highlightID := extractID(t, createHighlightResp.Body.Bytes())

	updateHighlightResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/portal/highlights/"+highlightID, map[string]any{
		"text":       "支持首页亮点与模块顺序灵活维护",
		"sort_order": 1,
		"enabled":    true,
	}, adminToken)
	require.Equal(t, http.StatusOK, updateHighlightResp.Code)

	homeResp := performRequest(t, app, http.MethodGet, "/api/v1/portal/home", nil, "", "")
	require.Equal(t, http.StatusOK, homeResp.Code)
	require.Contains(t, homeResp.Body.String(), "支持首页亮点与模块顺序灵活维护")
	require.Contains(t, homeResp.Body.String(), "\"module_key\":\"community\"")
	require.Contains(t, homeResp.Body.String(), "\"sort_order\":15")

	deleteHighlightResp := performRequest(t, app, http.MethodDelete, "/api/v1/admin/portal/highlights/"+highlightID, nil, adminToken, "")
	require.Equal(t, http.StatusOK, deleteHighlightResp.Code)
}

func TestPortalRankingConfigFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	getResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/portal/rankings-config", nil, adminToken, "")
	require.Equal(t, http.StatusOK, getResp.Code)
	require.Contains(t, getResp.Body.String(), "贡献排行榜")

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/portal/rankings-config", map[string]any{
		"title":    "社区贡献榜",
		"subtitle": "展示近期贡献最活跃的开发者。",
		"limit":    3,
		"enabled":  false,
	}, adminToken)
	require.Equal(t, http.StatusOK, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "社区贡献榜")

	homeResp := performRequest(t, app, http.MethodGet, "/api/v1/portal/home", nil, "", "")
	require.Equal(t, http.StatusOK, homeResp.Code)
	require.Contains(t, homeResp.Body.String(), "\"title\":\"社区贡献榜\"")
	require.Contains(t, homeResp.Body.String(), "\"limit\":3")
	require.Contains(t, homeResp.Body.String(), "\"enabled\":false")
}

func TestPortalHeroConfigFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	getResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/portal/hero-config", nil, adminToken, "")
	require.Equal(t, http.StatusOK, getResp.Code)
	require.Contains(t, getResp.Body.String(), "EIBotHub具生训练")

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/portal/hero-config", map[string]any{
		"tagline":          "机器人开发导读区",
		"title":            "一个面向具身智能资源协作的平台",
		"description":      "通过统一的模型、数据、模板和文档入口，帮助团队更快完成验证与复用。",
		"primary_button":   "立即上传模型",
		"secondary_button": "发布数据集",
		"search_button":    "搜索社区资源",
	}, adminToken)
	require.Equal(t, http.StatusOK, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "机器人开发导读区")

	homeResp := performRequest(t, app, http.MethodGet, "/api/v1/portal/home", nil, "", "")
	require.Equal(t, http.StatusOK, homeResp.Code)
	require.Contains(t, homeResp.Body.String(), "\"tagline\":\"机器人开发导读区\"")
	require.Contains(t, homeResp.Body.String(), "\"primary_button\":\"立即上传模型\"")
}

func TestFilterOptionsEndpoint(t *testing.T) {
	app := newTestApp(t)

	resp := performRequest(t, app, http.MethodGet, "/api/v1/filter-options", nil, "", "")
	require.Equal(t, http.StatusOK, resp.Code)
	require.Contains(t, resp.Body.String(), "\"robot_types\"")
	require.Contains(t, resp.Body.String(), "\"dataset_scenes\"")
	require.Contains(t, resp.Body.String(), "巡逻巡检")
	require.Contains(t, resp.Body.String(), "搬运")
}

func TestAdminFilterOptionConfigFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	createResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/content/filter-options", map[string]any{
		"kind":       "robot_type",
		"value":      "协同控制",
		"sort_order": 1,
		"enabled":    true,
	}, adminToken)
	require.Equal(t, http.StatusCreated, createResp.Code)
	filterOptionID := extractID(t, createResp.Body.Bytes())

	listResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/content/filter-options?kind=robot_type", nil, adminToken, "")
	require.Equal(t, http.StatusOK, listResp.Code)
	require.Contains(t, listResp.Body.String(), "协同控制")

	publicResp := performRequest(t, app, http.MethodGet, "/api/v1/filter-options", nil, "", "")
	require.Equal(t, http.StatusOK, publicResp.Code)
	require.Contains(t, publicResp.Body.String(), "协同控制")

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/content/filter-options/"+filterOptionID, map[string]any{
		"kind":       "robot_type",
		"value":      "协同机器人",
		"sort_order": 2,
		"enabled":    true,
	}, adminToken)
	require.Equal(t, http.StatusOK, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "协同机器人")

	deleteResp := performRequest(t, app, http.MethodDelete, "/api/v1/admin/content/filter-options/"+filterOptionID, nil, adminToken, "")
	require.Equal(t, http.StatusOK, deleteResp.Code)

	logsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/operations", nil, adminToken, "")
	require.Equal(t, http.StatusOK, logsResp.Code)
	require.Contains(t, logsResp.Body.String(), "\"resource_type\":\"filter-option\"")
	require.Contains(t, logsResp.Body.String(), "协同机器人")

	publicAfterDelete := performRequest(t, app, http.MethodGet, "/api/v1/filter-options", nil, "", "")
	require.Equal(t, http.StatusOK, publicAfterDelete.Code)
	require.NotContains(t, publicAfterDelete.Body.String(), "协同机器人")
}

func TestModelReviewAndSearchFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	modelResp := performMultipart(t, app, http.MethodPost, "/api/v1/models", map[string]string{
		"name":         "Factory Vision Pack",
		"summary":      "Factory line MVP model",
		"description":  "Detects anomalies for factory line robots",
		"tags":         "factory,inspection",
		"robot_type":   "巡逻巡检",
		"input_spec":   "RGB",
		"output_spec":  "events",
		"license":      "MIT",
		"dependencies": "opencv,onnxruntime",
		"version":      "v1.0.0",
		"changelog":    "first drop",
	}, "file", "factory.txt", "factory model", userToken)
	require.Equal(t, http.StatusCreated, modelResp.Code)
	modelID := extractID(t, modelResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/models/"+modelID+"/submit", nil, userToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	reviewList := performRequest(t, app, http.MethodGet, "/api/v1/admin/reviews?type=models", nil, adminToken, "")
	require.Equal(t, http.StatusOK, reviewList.Code)
	require.Contains(t, reviewList.Body.String(), "Factory Vision Pack")

	decisionResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/models/"+modelID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "ready for publish",
	}, adminToken)
	require.Equal(t, http.StatusOK, decisionResp.Code)

	searchResp := performRequest(t, app, http.MethodGet, "/api/v1/search?q=Factory&type=model", nil, "", "")
	require.Equal(t, http.StatusOK, searchResp.Code)
	require.Contains(t, searchResp.Body.String(), "Factory Vision Pack")
}

func TestDatasetAgreementAndDownloadFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	datasetResp := performMultipart(t, app, http.MethodPost, "/api/v1/datasets", map[string]string{
		"name":           "Route Samples",
		"summary":        "Route dataset",
		"description":    "Dataset for route planning",
		"tags":           "route,navigation",
		"sample_count":   "20",
		"device":         "LiDAR",
		"scene":          "移动作业",
		"privacy":        "public",
		"agreement_text": "accept before download",
		"version":        "v1.0.0",
		"changelog":      "initial",
		"sample_preview": "sample one\nsample two",
	}, "file", "routes.zip", "dataset binary", userToken)
	require.Equal(t, http.StatusCreated, datasetResp.Code)
	datasetID := extractID(t, datasetResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/submit", nil, userToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	decisionResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/datasets/"+datasetID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "publish",
	}, adminToken)
	require.Equal(t, http.StatusOK, decisionResp.Code)

	downloadFail := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, userToken, "")
	require.Equal(t, http.StatusBadRequest, downloadFail.Code)

	agreementResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/agreements/confirm", nil, userToken, "")
	require.Equal(t, http.StatusOK, agreementResp.Code)

	downloadResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, userToken, "")
	require.Equal(t, http.StatusOK, downloadResp.Code)

	historyResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/downloads", nil, userToken, "")
	require.Equal(t, http.StatusOK, historyResp.Code)
	require.Contains(t, historyResp.Body.String(), "Route Samples")

	sceneResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets?scene=移动作业", nil, "", "")
	require.Equal(t, http.StatusOK, sceneResp.Code)
	require.Contains(t, sceneResp.Body.String(), "Route Samples")
}

func TestModelEditFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	createResp := performMultipart(t, app, http.MethodPost, "/api/v1/models", map[string]string{
		"name":         "Editable Vision Pack",
		"summary":      "Editable summary",
		"description":  "Editable description for model flow test",
		"tags":         "vision,test",
		"robot_type":   "巡逻巡检",
		"input_spec":   "RGB",
		"output_spec":  "events",
		"license":      "MIT",
		"dependencies": "opencv",
		"version":      "v1.0.0",
		"changelog":    "initial release",
	}, "file", "editable.txt", "editable model", userToken)
	require.Equal(t, http.StatusCreated, createResp.Code)
	modelID := extractID(t, createResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/models/"+modelID+"/submit", nil, userToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/models/"+modelID, map[string]any{
		"name":         "Editable Vision Pack V2",
		"summary":      "Updated summary",
		"description":  "Updated description for model flow test",
		"tags":         "vision,test,updated",
		"robot_type":   "搬运",
		"input_spec":   "RGBD",
		"output_spec":  "obstacle map",
		"license":      "Apache-2.0",
		"dependencies": "opencv,onnxruntime",
	}, userToken)
	require.Equal(t, http.StatusOK, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "Editable Vision Pack V2")
	require.Contains(t, updateResp.Body.String(), "\"status\":\"draft\"")

	detailResp := performRequest(t, app, http.MethodGet, "/api/v1/models/"+modelID, nil, userToken, "")
	require.Equal(t, http.StatusOK, detailResp.Code)
	require.Contains(t, detailResp.Body.String(), "updated")

	resubmitResp := performRequest(t, app, http.MethodPost, "/api/v1/models/"+modelID+"/submit", nil, userToken, "")
	require.Equal(t, http.StatusOK, resubmitResp.Code)

	detailAfterSubmit := performRequest(t, app, http.MethodGet, "/api/v1/models/"+modelID, nil, userToken, "")
	require.Equal(t, http.StatusOK, detailAfterSubmit.Code)
	require.Contains(t, detailAfterSubmit.Body.String(), "\"status\":\"pending\"")

	rejectResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/models/"+modelID+"/decision", map[string]any{
		"decision": "rejected",
		"comment":  "metadata needs refinement",
	}, adminToken)
	require.Equal(t, http.StatusOK, rejectResp.Code)

	detailAfterReject := performRequest(t, app, http.MethodGet, "/api/v1/models/"+modelID, nil, userToken, "")
	require.Equal(t, http.StatusOK, detailAfterReject.Code)
	require.Contains(t, detailAfterReject.Body.String(), "\"review_comment\":\"metadata needs refinement\"")

	uploadsResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/uploads", nil, userToken, "")
	require.Equal(t, http.StatusOK, uploadsResp.Code)
	require.Contains(t, uploadsResp.Body.String(), "\"review_comment\":\"metadata needs refinement\"")
}

func TestDatasetEditFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	createResp := performMultipart(t, app, http.MethodPost, "/api/v1/datasets", map[string]string{
		"name":           "Editable Dataset Pack",
		"summary":        "Editable dataset summary",
		"description":    "Editable dataset description for update flow",
		"tags":           "dataset,test",
		"sample_count":   "24",
		"device":         "RGB Camera",
		"scene":          "巡逻巡检",
		"privacy":        "public",
		"agreement_text": "accept before download",
		"version":        "v1.0.0",
		"changelog":      "initial release",
		"sample_preview": "old sample 1\nold sample 2",
	}, "file", "editable-dataset.zip", "dataset binary", userToken)
	require.Equal(t, http.StatusCreated, createResp.Code)
	datasetID := extractID(t, createResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/submit", nil, userToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/datasets/"+datasetID, map[string]any{
		"name":           "Editable Dataset Pack V2",
		"summary":        "Updated dataset summary",
		"description":    "Updated dataset description for update flow",
		"tags":           "dataset,test,updated",
		"sample_count":   48,
		"device":         "LiDAR",
		"scene":          "移动作业",
		"privacy":        "restricted",
		"agreement_text": "updated agreement text",
		"sample_preview": "new sample 1\nnew sample 2",
	}, userToken)
	require.Equal(t, http.StatusOK, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "Editable Dataset Pack V2")
	require.Contains(t, updateResp.Body.String(), "\"status\":\"draft\"")
	require.Contains(t, updateResp.Body.String(), "\"privacy\":\"restricted\"")

	detailResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/"+datasetID, nil, userToken, "")
	require.Equal(t, http.StatusOK, detailResp.Code)
	require.Contains(t, detailResp.Body.String(), "new sample 1")
	require.NotContains(t, detailResp.Body.String(), "old sample 1")

	resubmitResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/submit", nil, userToken, "")
	require.Equal(t, http.StatusOK, resubmitResp.Code)

	detailAfterSubmit := performRequest(t, app, http.MethodGet, "/api/v1/datasets/"+datasetID, nil, userToken, "")
	require.Equal(t, http.StatusOK, detailAfterSubmit.Code)
	require.Contains(t, detailAfterSubmit.Body.String(), "\"status\":\"pending\"")

	rejectResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/datasets/"+datasetID+"/decision", map[string]any{
		"decision": "rejected",
		"comment":  "agreement text is too generic",
	}, adminToken)
	require.Equal(t, http.StatusOK, rejectResp.Code)

	detailAfterReject := performRequest(t, app, http.MethodGet, "/api/v1/datasets/"+datasetID, nil, userToken, "")
	require.Equal(t, http.StatusOK, detailAfterReject.Code)
	require.Contains(t, detailAfterReject.Body.String(), "\"review_comment\":\"agreement text is too generic\"")

	uploadsResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/uploads", nil, userToken, "")
	require.Equal(t, http.StatusOK, uploadsResp.Code)
	require.Contains(t, uploadsResp.Body.String(), "\"review_comment\":\"agreement text is too generic\"")
}

func TestDatasetSamplePreviewAndDownloadPackages(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	samplesResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/1/samples", nil, "", "")
	require.Equal(t, http.StatusOK, samplesResp.Code)
	require.Contains(t, samplesResp.Body.String(), "\"sample_type\":\"image\"")
	require.Contains(t, samplesResp.Body.String(), "\"sample_type\":\"pointcloud\"")
	require.Contains(t, samplesResp.Body.String(), "/api/v1/files/download/")

	agreementResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/1/agreements/confirm", nil, userToken, "")
	require.Equal(t, http.StatusOK, agreementResp.Code)

	taskResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/1/download-packages", map[string]any{
		"parts": 4,
	}, userToken)
	require.Equal(t, http.StatusCreated, taskResp.Code)
	require.Contains(t, taskResp.Body.String(), "\"total_parts\":4")
	require.Contains(t, taskResp.Body.String(), "/api/v1/files/download/")

	taskListResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/1/download-packages", nil, userToken, "")
	require.Equal(t, http.StatusOK, taskListResp.Code)
	require.Contains(t, taskListResp.Body.String(), "\"status\":\"ready\"")
}

func TestRestrictedDatasetAccessApprovalFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")
	ownerToken := loginToken(t, app, "demo@example.com", "Demo123!")
	secondAdminToken := createAdminUserToken(t, app, "reviewer", "reviewer@example.com", "Reviewer123!")

	registerResp := performJSON(t, app, http.MethodPost, "/api/v1/auth/register", map[string]any{
		"username": "guest",
		"email":    "guest@example.com",
		"password": "Guest123!",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &payload))
	guestToken := payload["data"].(map[string]any)["token"].(string)

	datasetResp := performMultipart(t, app, http.MethodPost, "/api/v1/datasets", map[string]string{
		"name":           "Restricted Route Samples",
		"summary":        "Restricted dataset",
		"description":    "Restricted dataset for approval flow",
		"tags":           "route,restricted",
		"sample_count":   "12",
		"device":         "RGB Camera",
		"scene":          "移动作业",
		"privacy":        "restricted",
		"agreement_text": "need approval before download",
		"version":        "v1.0.0",
		"changelog":      "initial",
		"sample_preview": "alpha\nbeta",
	}, "file", "restricted.zip", "dataset binary", ownerToken)
	require.Equal(t, http.StatusCreated, datasetResp.Code)
	datasetID := extractID(t, datasetResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/submit", nil, ownerToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	reviewResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/datasets/"+datasetID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "publish",
	}, adminToken)
	require.Equal(t, http.StatusOK, reviewResp.Code)

	agreementResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/agreements/confirm", nil, guestToken, "")
	require.Equal(t, http.StatusOK, agreementResp.Code)

	downloadDenied := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, guestToken, "")
	require.Equal(t, http.StatusForbidden, downloadDenied.Code)
	require.Contains(t, downloadDenied.Body.String(), "dataset_access_request_required")

	createAccessResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/access-requests", map[string]any{
		"reason": "需要用于联调验证",
	}, guestToken)
	require.Equal(t, http.StatusCreated, createAccessResp.Code)
	accessRequestID := extractID(t, createAccessResp.Body.Bytes())

	myRequestResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/"+datasetID+"/access-requests/me", nil, guestToken, "")
	require.Equal(t, http.StatusOK, myRequestResp.Code)
	require.Contains(t, myRequestResp.Body.String(), "\"status\":\"pending\"")

	adminListResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/datasets/access-requests", nil, adminToken, "")
	require.Equal(t, http.StatusOK, adminListResp.Code)
	require.Contains(t, adminListResp.Body.String(), "guest")
	require.Contains(t, adminListResp.Body.String(), "\"dataset_privacy\":\"restricted\"")
	require.Contains(t, adminListResp.Body.String(), "Restricted Route Samples")
	require.Contains(t, adminListResp.Body.String(), "\"sla_hours\":48")

	filteredAdminListResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/datasets/access-requests?status=pending&privacy=restricted&q=guest", nil, adminToken, "")
	require.Equal(t, http.StatusOK, filteredAdminListResp.Code)
	require.Contains(t, filteredAdminListResp.Body.String(), "\"dataset_owner_name\":\"demo\"")
	ontrackAdminListResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/datasets/access-requests?sla_status=ontrack&privacy=restricted", nil, adminToken, "")
	require.Equal(t, http.StatusOK, ontrackAdminListResp.Code)
	require.Contains(t, ontrackAdminListResp.Body.String(), "Restricted Route Samples")

	ownerNotificationsResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/notifications", nil, ownerToken, "")
	require.Equal(t, http.StatusOK, ownerNotificationsResp.Code)
	require.Contains(t, ownerNotificationsResp.Body.String(), "Restricted Route Samples")
	require.Contains(t, ownerNotificationsResp.Body.String(), "SLA 48 小时")

	userNotificationsResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/notifications", nil, guestToken, "")
	require.Equal(t, http.StatusOK, userNotificationsResp.Code)
	require.Contains(t, userNotificationsResp.Body.String(), "访问申请已提交")
	require.Contains(t, userNotificationsResp.Body.String(), "SLA 48 小时")

	ownerFilteredResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/datasets/access-requests?owner_id=2&privacy=restricted", nil, adminToken, "")
	require.Equal(t, http.StatusOK, ownerFilteredResp.Code)
	require.Contains(t, ownerFilteredResp.Body.String(), "Restricted Route Samples")

	batchDecisionResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/batch-decision", map[string]any{
		"ids":      []int{mustAtoi(t, accessRequestID)},
		"decision": "approved",
		"comment":  "允许下载",
	}, adminToken)
	require.Equal(t, http.StatusOK, batchDecisionResp.Code)

	stagePendingResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/datasets/access-requests?status=pending&privacy=restricted&q=guest", nil, adminToken, "")
	require.Equal(t, http.StatusOK, stagePendingResp.Code)
	require.Contains(t, stagePendingResp.Body.String(), "\"approval_stage\":1")
	require.Contains(t, stagePendingResp.Body.String(), "\"required_approvals\":2")

	adminDecisionResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/"+accessRequestID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "允许下载",
	}, adminToken)
	require.Equal(t, http.StatusConflict, adminDecisionResp.Code)

	finalDecisionResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/"+accessRequestID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "终审通过",
	}, secondAdminToken)
	require.Equal(t, http.StatusOK, finalDecisionResp.Code)

	reviewedListResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/datasets/access-requests?status=approved&privacy=restricted&q=guest", nil, adminToken, "")
	require.Equal(t, http.StatusOK, reviewedListResp.Code)
	require.Contains(t, reviewedListResp.Body.String(), "\"status\":\"approved\"")

	downloadAllowed := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, guestToken, "")
	require.Equal(t, http.StatusOK, downloadAllowed.Code)

	packageResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download-packages", map[string]any{
		"parts": 2,
	}, guestToken)
	require.Equal(t, http.StatusCreated, packageResp.Code)
}

func TestDatasetAccessHistoryEndpoints(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")
	ownerToken := loginToken(t, app, "demo@example.com", "Demo123!")

	registerResp := performJSON(t, app, http.MethodPost, "/api/v1/auth/register", map[string]any{
		"username": "history-user",
		"email":    "history-user@example.com",
		"password": "History123!",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &payload))
	userToken := payload["data"].(map[string]any)["token"].(string)

	datasetResp := performMultipart(t, app, http.MethodPost, "/api/v1/datasets", map[string]string{
		"name":           "History Restricted Dataset",
		"summary":        "Restricted dataset for history flow",
		"description":    "Dataset used to verify user-side access request history",
		"tags":           "history,restricted",
		"sample_count":   "16",
		"device":         "RGB Camera",
		"scene":          "巡逻巡检",
		"privacy":        "restricted",
		"agreement_text": "need approval before download",
		"version":        "v1.0.0",
		"changelog":      "initial",
		"sample_preview": "preview sample",
	}, "file", "history-dataset.zip", "dataset binary", ownerToken)
	require.Equal(t, http.StatusCreated, datasetResp.Code)
	datasetID := extractID(t, datasetResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/submit", nil, ownerToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	reviewResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/datasets/"+datasetID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "publish for history flow",
	}, adminToken)
	require.Equal(t, http.StatusOK, reviewResp.Code)

	firstRequestResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/access-requests", map[string]any{
		"reason": "首次申请，材料还不完整",
	}, userToken)
	require.Equal(t, http.StatusCreated, firstRequestResp.Code)
	firstRequestID := extractID(t, firstRequestResp.Body.Bytes())

	firstReviewResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/"+firstRequestID+"/decision", map[string]any{
		"decision": "rejected",
		"comment":  "缺少使用场景说明",
	}, adminToken)
	require.Equal(t, http.StatusOK, firstReviewResp.Code)

	secondRequestResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/access-requests", map[string]any{
		"reason": "已补充巡检联调场景说明，申请再次审核",
	}, userToken)
	require.Equal(t, http.StatusCreated, secondRequestResp.Code)

	datasetHistoryResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/"+datasetID+"/access-requests/history", nil, userToken, "")
	require.Equal(t, http.StatusOK, datasetHistoryResp.Code)
	require.Contains(t, datasetHistoryResp.Body.String(), "首次申请，材料还不完整")
	require.Contains(t, datasetHistoryResp.Body.String(), "缺少使用场景说明")
	require.Contains(t, datasetHistoryResp.Body.String(), "已补充巡检联调场景说明，申请再次审核")

	myHistoryResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/dataset-access-requests", nil, userToken, "")
	require.Equal(t, http.StatusOK, myHistoryResp.Code)
	require.Contains(t, myHistoryResp.Body.String(), "History Restricted Dataset")
	require.Contains(t, myHistoryResp.Body.String(), "\"dataset_privacy\":\"restricted\"")
	require.Contains(t, myHistoryResp.Body.String(), "\"status\":\"rejected\"")
}

func TestDatasetAccessGrantConstraints(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")
	ownerToken := loginToken(t, app, "demo@example.com", "Demo123!")
	secondAdminToken := createAdminUserToken(t, app, "grant-reviewer", "grant-reviewer@example.com", "GrantReviewer123!")

	registerResp := performJSON(t, app, http.MethodPost, "/api/v1/auth/register", map[string]any{
		"username": "limit-user",
		"email":    "limit-user@example.com",
		"password": "Limit123!",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &payload))
	userToken := payload["data"].(map[string]any)["token"].(string)

	datasetResp := performMultipart(t, app, http.MethodPost, "/api/v1/datasets", map[string]string{
		"name":           "Grant Restricted Dataset",
		"summary":        "Restricted dataset for grant constraint flow",
		"description":    "Dataset used to verify approval validity and download limits",
		"tags":           "grant,restricted",
		"sample_count":   "18",
		"device":         "RGB Camera",
		"scene":          "移动作业",
		"privacy":        "restricted",
		"agreement_text": "need approval before download",
		"version":        "v1.0.0",
		"changelog":      "initial",
		"sample_preview": "preview sample",
	}, "file", "grant-dataset.zip", "dataset binary", ownerToken)
	require.Equal(t, http.StatusCreated, datasetResp.Code)
	datasetID := extractID(t, datasetResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/submit", nil, ownerToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	reviewResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/datasets/"+datasetID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "publish for grant flow",
	}, adminToken)
	require.Equal(t, http.StatusOK, reviewResp.Code)

	requestResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/access-requests", map[string]any{
		"reason": "申请数据集用于多轮联调验证",
	}, userToken)
	require.Equal(t, http.StatusCreated, requestResp.Code)
	requestID := extractID(t, requestResp.Body.Bytes())

	approveResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/"+requestID+"/decision", map[string]any{
		"decision":       "approved",
		"comment":        "初审通过",
		"valid_days":     0,
		"download_limit": 0,
	}, adminToken)
	require.Equal(t, http.StatusOK, approveResp.Code)

	finalApproveResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/"+requestID+"/decision", map[string]any{
		"decision":       "approved",
		"comment":        "限期开放两次下载",
		"valid_days":     7,
		"download_limit": 2,
	}, secondAdminToken)
	require.Equal(t, http.StatusOK, finalApproveResp.Code)

	myRequestResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/"+datasetID+"/access-requests/me", nil, userToken, "")
	require.Equal(t, http.StatusOK, myRequestResp.Code)
	require.Contains(t, myRequestResp.Body.String(), "\"approval_stage\":2")
	require.Contains(t, myRequestResp.Body.String(), "\"download_limit\":2")
	require.Contains(t, myRequestResp.Body.String(), "\"download_count\":0")
	require.Contains(t, myRequestResp.Body.String(), "\"authorization_active\":true")

	agreementResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/agreements/confirm", nil, userToken, "")
	require.Equal(t, http.StatusOK, agreementResp.Code)

	downloadResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, userToken, "")
	require.Equal(t, http.StatusOK, downloadResp.Code)

	packageResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download-packages", map[string]any{
		"parts": 2,
	}, userToken)
	require.Equal(t, http.StatusCreated, packageResp.Code)

	limitExceededResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, userToken, "")
	require.Equal(t, http.StatusForbidden, limitExceededResp.Code)
	require.Contains(t, limitExceededResp.Body.String(), "dataset_access_limit_exceeded")

	historyResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/"+datasetID+"/access-requests/history", nil, userToken, "")
	require.Equal(t, http.StatusOK, historyResp.Code)
	require.Contains(t, historyResp.Body.String(), "\"download_count\":2")
	require.Contains(t, historyResp.Body.String(), "\"authorization_active\":false")
}

func TestDatasetAccessGrantExpiry(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")
	ownerToken := loginToken(t, app, "demo@example.com", "Demo123!")
	secondAdminToken := createAdminUserToken(t, app, "expiry-reviewer", "expiry-reviewer@example.com", "ExpiryReviewer123!")

	registerResp := performJSON(t, app, http.MethodPost, "/api/v1/auth/register", map[string]any{
		"username": "expiry-user",
		"email":    "expiry-user@example.com",
		"password": "Expiry123!",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &payload))
	userToken := payload["data"].(map[string]any)["token"].(string)

	datasetResp := performMultipart(t, app, http.MethodPost, "/api/v1/datasets", map[string]string{
		"name":           "Expiry Restricted Dataset",
		"summary":        "Restricted dataset for expiry flow",
		"description":    "Dataset used to verify expired access grants are blocked",
		"tags":           "expiry,restricted",
		"sample_count":   "14",
		"device":         "LiDAR",
		"scene":          "巡逻巡检",
		"privacy":        "restricted",
		"agreement_text": "need approval before download",
		"version":        "v1.0.0",
		"changelog":      "initial",
		"sample_preview": "preview sample",
	}, "file", "expiry-dataset.zip", "dataset binary", ownerToken)
	require.Equal(t, http.StatusCreated, datasetResp.Code)
	datasetID := extractID(t, datasetResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/submit", nil, ownerToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	reviewResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/datasets/"+datasetID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "publish for expiry flow",
	}, adminToken)
	require.Equal(t, http.StatusOK, reviewResp.Code)

	requestResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/access-requests", map[string]any{
		"reason": "申请数据集用于到期测试",
	}, userToken)
	require.Equal(t, http.StatusCreated, requestResp.Code)
	requestID := extractID(t, requestResp.Body.Bytes())

	approveResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/"+requestID+"/decision", map[string]any{
		"decision":       "approved",
		"comment":        "初审通过",
		"valid_days":     0,
		"download_limit": 0,
	}, adminToken)
	require.Equal(t, http.StatusOK, approveResp.Code)

	finalApproveResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/"+requestID+"/decision", map[string]any{
		"decision":       "approved",
		"comment":        "临时授权",
		"valid_days":     1,
		"download_limit": 0,
	}, secondAdminToken)
	require.Equal(t, http.StatusOK, finalApproveResp.Code)

	require.NoError(t, app.DB.Model(&model.DatasetAccessRequest{}).
		Where("id = ?", mustAtoi(t, requestID)).
		Update("approval_expires_at", time.Now().Add(-2*time.Hour)).Error)

	agreementResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/agreements/confirm", nil, userToken, "")
	require.Equal(t, http.StatusOK, agreementResp.Code)

	expiredResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, userToken, "")
	require.Equal(t, http.StatusForbidden, expiredResp.Code)
	require.Contains(t, expiredResp.Body.String(), "dataset_access_expired")
}

func TestInternalDatasetApprovedVerificationDirectAccess(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")
	ownerToken := loginToken(t, app, "demo@example.com", "Demo123!")

	registerResp := performJSON(t, app, http.MethodPost, "/api/v1/auth/register", map[string]any{
		"username": "internal-user",
		"email":    "internal-user@example.com",
		"password": "Internal123!",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &payload))
	userToken := payload["data"].(map[string]any)["token"].(string)

	applyResp := performJSON(t, app, http.MethodPost, "/api/v1/developer-verifications", map[string]any{
		"verification_type": "personal",
		"real_name":         "Internal User",
		"organization":      "Open Community Lab",
		"materials":         "portfolio",
		"reason":            "申请访问内部数据集进行联调。",
	}, userToken)
	require.Equal(t, http.StatusCreated, applyResp.Code)
	verificationID := extractID(t, applyResp.Body.Bytes())

	reviewVerificationResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/verifications/"+verificationID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "认证通过",
	}, adminToken)
	require.Equal(t, http.StatusOK, reviewVerificationResp.Code)

	datasetResp := performMultipart(t, app, http.MethodPost, "/api/v1/datasets", map[string]string{
		"name":           "Internal Access Dataset",
		"summary":        "Internal dataset for permission matrix flow",
		"description":    "Dataset used to verify approved developers can access internal resources directly",
		"tags":           "internal,permission",
		"sample_count":   "12",
		"device":         "RGB Camera",
		"scene":          "移动作业",
		"privacy":        "internal",
		"agreement_text": "accept before download",
		"version":        "v1.0.0",
		"changelog":      "initial",
		"sample_preview": "preview sample",
	}, "file", "internal-dataset.zip", "dataset binary", ownerToken)
	require.Equal(t, http.StatusCreated, datasetResp.Code)
	datasetID := extractID(t, datasetResp.Body.Bytes())

	submitResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/submit", nil, ownerToken, "")
	require.Equal(t, http.StatusOK, submitResp.Code)

	publishResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/reviews/datasets/"+datasetID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "publish internal dataset",
	}, adminToken)
	require.Equal(t, http.StatusOK, publishResp.Code)

	agreementResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/agreements/confirm", nil, userToken, "")
	require.Equal(t, http.StatusOK, agreementResp.Code)

	downloadResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, userToken, "")
	require.Equal(t, http.StatusOK, downloadResp.Code)

	packageResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download-packages", map[string]any{
		"parts": 2,
	}, userToken)
	require.Equal(t, http.StatusCreated, packageResp.Code)

	createAccessResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/access-requests", map[string]any{
		"reason": "已经认证，理论上不需要申请",
	}, userToken)
	require.Equal(t, http.StatusConflict, createAccessResp.Code)
	require.Contains(t, createAccessResp.Body.String(), "dataset_internal_auto_access")

	myRequestResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/"+datasetID+"/access-requests/me", nil, userToken, "")
	require.Equal(t, http.StatusOK, myRequestResp.Code)
	require.Contains(t, myRequestResp.Body.String(), "\"data\":null")
}

func TestCreateModelValidation(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	resp := performMultipart(t, app, http.MethodPost, "/api/v1/models", map[string]string{
		"name": "x",
	}, "file", "invalid.txt", "invalid", userToken)

	require.Equal(t, http.StatusBadRequest, resp.Code)
	require.Contains(t, resp.Body.String(), "invalid_request")
}

func TestV1ModelEvaluationAndTemplateInteraction(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	evaluationResp := performJSON(t, app, http.MethodPost, "/api/v1/models/1/evaluations", map[string]any{
		"benchmark": "Night Shift Benchmark",
		"summary":   "夜间巡检表现稳定",
		"score":     91,
		"notes":     "弱光场景下精度保持良好",
	}, userToken)
	require.Equal(t, http.StatusCreated, evaluationResp.Code)

	ratingResp := performJSON(t, app, http.MethodPost, "/api/v1/task-templates/1/ratings", map[string]any{
		"score":    5,
		"feedback": "模板编排清晰",
	}, userToken)
	require.Equal(t, http.StatusOK, ratingResp.Code)

	commentResp := performJSON(t, app, http.MethodPost, "/api/v1/task-templates/1/comments", map[string]any{
		"content": "适合做巡检场景的首个版本模板。",
	}, userToken)
	require.Equal(t, http.StatusCreated, commentResp.Code)

	listResp := performRequest(t, app, http.MethodGet, "/api/v1/models/1/evaluations", nil, "", "")
	require.Equal(t, http.StatusOK, listResp.Code)
	require.Contains(t, listResp.Body.String(), "Night Shift Benchmark")

	templateComments := performRequest(t, app, http.MethodGet, "/api/v1/task-templates/1/comments", nil, "", "")
	require.Equal(t, http.StatusOK, templateComments.Code)
	require.Contains(t, templateComments.Body.String(), "适合做巡检场景的首个版本模板。")
}

func TestV1SkillsDiscussionsAndSearchHot(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	skillResp := performJSON(t, app, http.MethodPost, "/api/v1/skills", map[string]any{
		"name":         "仓储告警聚合技能",
		"summary":      "把多源告警聚合成统一事件流。",
		"description":  "用于搬运与巡检机器人联合场景下的告警聚合与转发。",
		"category":     "巡逻巡检",
		"scene":        "移动作业",
		"guide":        "配置告警源、绑定 webhook、接入任务模板。",
		"resource_ref": "Warehouse Carrier Net,Inspection Route Set",
	}, userToken)
	require.Equal(t, http.StatusCreated, skillResp.Code)
	skillID := extractID(t, skillResp.Body.Bytes())

	forkResp := performRequest(t, app, http.MethodPost, "/api/v1/skills/"+skillID+"/fork", nil, userToken, "")
	require.Equal(t, http.StatusCreated, forkResp.Code)

	discussionResp := performJSON(t, app, http.MethodPost, "/api/v1/discussions", map[string]any{
		"title":   "如何组合技能与模板？",
		"tag":     "最佳实践",
		"content": "当前项目希望把告警技能接入任务模板，想了解推荐的复用方式。",
	}, userToken)
	require.Equal(t, http.StatusCreated, discussionResp.Code)
	discussionID := extractID(t, discussionResp.Body.Bytes())

	discussionComment := performJSON(t, app, http.MethodPost, "/api/v1/discussions/"+discussionID+"/comments", map[string]any{
		"content": "建议先在模板里只引用稳定版技能。",
	}, userToken)
	require.Equal(t, http.StatusCreated, discussionComment.Code)

	followResp := performRequest(t, app, http.MethodPost, "/api/v1/community/users/1/follow", nil, userToken, "")
	require.Equal(t, http.StatusOK, followResp.Code)

	searchResp := performRequest(t, app, http.MethodGet, "/api/v1/search?q=巡检&type=dataset", nil, "", "")
	require.Equal(t, http.StatusOK, searchResp.Code)
	hotResp := performRequest(t, app, http.MethodGet, "/api/v1/search/hot", nil, "", "")
	require.Equal(t, http.StatusOK, hotResp.Code)
	require.Contains(t, hotResp.Body.String(), "巡检")
	recommendedResp := performRequest(t, app, http.MethodGet, "/api/v1/search/recommended", nil, "", "")
	require.Equal(t, http.StatusOK, recommendedResp.Code)
	require.Contains(t, recommendedResp.Body.String(), "夜间安防")

	statsResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/1/stats", nil, "", "")
	require.Equal(t, http.StatusOK, statsResp.Code)
	require.Contains(t, statsResp.Body.String(), "download_count")
}

func TestSearchKeywordAdminFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	listResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/portal/search-keywords", nil, adminToken, "")
	require.Equal(t, http.StatusOK, listResp.Code)
	require.Contains(t, listResp.Body.String(), "recommended")

	createResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/portal/search-keywords", map[string]any{
		"query":        "机器人接入",
		"keyword_type": "recommended",
		"sort_order":   3,
		"enabled":      true,
	}, adminToken)
	require.Equal(t, http.StatusCreated, createResp.Code)
	keywordID := extractID(t, createResp.Body.Bytes())

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/portal/search-keywords/"+keywordID, map[string]any{
		"query":        "机器人接入指南",
		"keyword_type": "hot",
		"sort_order":   1,
		"enabled":      true,
	}, adminToken)
	require.Equal(t, http.StatusOK, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "机器人接入指南")

	hotResp := performRequest(t, app, http.MethodGet, "/api/v1/search/hot", nil, "", "")
	require.Equal(t, http.StatusOK, hotResp.Code)
	require.Contains(t, hotResp.Body.String(), "机器人接入指南")

	deleteResp := performRequest(t, app, http.MethodDelete, "/api/v1/admin/portal/search-keywords/"+keywordID, nil, adminToken, "")
	require.Equal(t, http.StatusOK, deleteResp.Code)
}

func TestSearchResponseEnhancements(t *testing.T) {
	app := newTestApp(t)

	searchResp := performRequest(t, app, http.MethodGet, "/api/v1/search?q=巡检", nil, "", "")
	require.Equal(t, http.StatusOK, searchResp.Code)
	require.Contains(t, searchResp.Body.String(), "\"type_counts\"")
	require.Contains(t, searchResp.Body.String(), "\"related_items\"")
	require.Contains(t, searchResp.Body.String(), "\"same_type_items\"")
	require.Contains(t, searchResp.Body.String(), "\"suggested_queries\"")

	typeResp := performRequest(t, app, http.MethodGet, "/api/v1/search?q=Vision&type=model", nil, "", "")
	require.Equal(t, http.StatusOK, typeResp.Code)
	require.Contains(t, typeResp.Body.String(), "\"focus_type\":\"model\"")
	require.Contains(t, typeResp.Body.String(), "\"same_type_items\"")
}

func TestFeaturedResourceBadgeLabelFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	createResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/portal/featured-resources", map[string]any{
		"resource_type": "model",
		"resource_id":   1,
		"badge_label":   "运营精选",
		"sort_order":    1,
		"enabled":       true,
	}, adminToken)
	require.Equal(t, http.StatusCreated, createResp.Code)
	featuredID := extractID(t, createResp.Body.Bytes())

	listResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/portal/featured-resources", nil, adminToken, "")
	require.Equal(t, http.StatusOK, listResp.Code)
	require.Contains(t, listResp.Body.String(), "运营精选")

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/portal/featured-resources/"+featuredID, map[string]any{
		"resource_type": "model",
		"resource_id":   1,
		"badge_label":   "首页精选",
		"sort_order":    1,
		"enabled":       true,
	}, adminToken)
	require.Equal(t, http.StatusOK, updateResp.Code)
	require.Contains(t, updateResp.Body.String(), "首页精选")

	homeResp := performRequest(t, app, http.MethodGet, "/api/v1/portal/home", nil, "", "")
	require.Equal(t, http.StatusOK, homeResp.Code)
	require.Contains(t, homeResp.Body.String(), "\"badge_label\":\"首页精选\"")
}

func TestAdminModelRecommendTagFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	listResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/content/model-recommend-tags", nil, adminToken, "")
	require.Equal(t, http.StatusOK, listResp.Code)
	require.Contains(t, listResp.Body.String(), "Patrol Corridor Vision")
	modelID := extractFirstListIDNumber(t, listResp.Body.Bytes())

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/content/model-recommend-tags/"+strconv.Itoa(modelID), map[string]any{
		"recommend_tag": "模型精选",
	}, adminToken)
	require.Equal(t, http.StatusOK, updateResp.Code)

	modelResp := performRequest(t, app, http.MethodGet, "/api/v1/models/"+strconv.Itoa(modelID), nil, "", "")
	require.Equal(t, http.StatusOK, modelResp.Code)
	require.Contains(t, modelResp.Body.String(), "\"recommend_tag\":\"模型精选\"")
	require.Contains(t, modelResp.Body.String(), "\"badge_label\":\"模型精选\"")

	modelListResp := performRequest(t, app, http.MethodGet, "/api/v1/models", nil, "", "")
	require.Equal(t, http.StatusOK, modelListResp.Code)
	require.Contains(t, modelListResp.Body.String(), "\"badge_label\":\"模型精选\"")
}

func TestScenePageFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	homeResp := performRequest(t, app, http.MethodGet, "/api/v1/portal/home", nil, "", "")
	require.Equal(t, http.StatusOK, homeResp.Code)
	require.Contains(t, homeResp.Body.String(), "\"scene_pages\"")
	require.Contains(t, homeResp.Body.String(), "巡逻巡检")

	sceneListResp := performRequest(t, app, http.MethodGet, "/api/v1/scenes", nil, "", "")
	require.Equal(t, http.StatusOK, sceneListResp.Code)
	require.Contains(t, sceneListResp.Body.String(), "warehouse-ops")

	sceneDetailResp := performRequest(t, app, http.MethodGet, "/api/v1/scenes/patrol-inspection", nil, "", "")
	require.Equal(t, http.StatusOK, sceneDetailResp.Code)
	require.Contains(t, sceneDetailResp.Body.String(), "\"scene\"")
	require.Contains(t, sceneDetailResp.Body.String(), "巡逻巡检")

	adminSceneResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/portal/scenes", map[string]any{
		"slug":        "night-security",
		"name":        "夜间安防",
		"tagline":     "Night Security",
		"summary":     "面向夜间安防与弱光告警联动的场景专题页。",
		"description": "聚合夜间安防相关模型、模板和案例。",
		"sort_order":  30,
		"enabled":     true,
	}, adminToken)
	require.Equal(t, http.StatusCreated, adminSceneResp.Code)

	sceneDetailCreatedResp := performRequest(t, app, http.MethodGet, "/api/v1/scenes/night-security", nil, "", "")
	require.Equal(t, http.StatusOK, sceneDetailCreatedResp.Code)
	require.Contains(t, sceneDetailCreatedResp.Body.String(), "夜间安防")
}

func TestRoleBasedAdminAccess(t *testing.T) {
	app := newTestApp(t)
	operatorToken := createRoleUserToken(t, app, "operator-user", "operator-user@example.com", "Operator123!", model.RoleOperator)
	reviewerToken := createRoleUserToken(t, app, "reviewer-user", "reviewer-user@example.com", "Reviewer123!", model.RoleReviewer)

	operatorPortalResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/portal/modules", nil, operatorToken, "")
	require.Equal(t, http.StatusOK, operatorPortalResp.Code)

	operatorReviewsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/reviews?type=models", nil, operatorToken, "")
	require.Equal(t, http.StatusForbidden, operatorReviewsResp.Code)
	require.Contains(t, operatorReviewsResp.Body.String(), "insufficient permission")

	reviewerReviewsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/reviews?type=models", nil, reviewerToken, "")
	require.Equal(t, http.StatusOK, reviewerReviewsResp.Code)

	reviewerPortalResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/portal/modules", nil, reviewerToken, "")
	require.Equal(t, http.StatusForbidden, reviewerPortalResp.Code)
	require.Contains(t, reviewerPortalResp.Body.String(), "insufficient permission")
}

func TestSuperAdminAccessAndPermissions(t *testing.T) {
	app := newTestApp(t)
	superAdminToken := createRoleUserToken(t, app, "super-admin-user", "super-admin-user@example.com", "Super123!", model.RoleSuperAdmin)

	meResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me", nil, superAdminToken, "")
	require.Equal(t, http.StatusOK, meResp.Code)
	require.Contains(t, meResp.Body.String(), "\"role\":\"super_admin\"")
	require.Contains(t, meResp.Body.String(), model.PermissionRewardManage)
	require.Contains(t, meResp.Body.String(), model.PermissionWikiManage)

	rewardResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/rewards/overview", nil, superAdminToken, "")
	require.Equal(t, http.StatusOK, rewardResp.Code)

	wikiResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/wiki/pages", nil, superAdminToken, "")
	require.Equal(t, http.StatusOK, wikiResp.Code)
}

func TestApprovedVerificationPromotesDeveloperRole(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	registerResp := performJSON(t, app, http.MethodPost, "/api/v1/auth/register", map[string]any{
		"username": "developer-upgrade",
		"email":    "developer-upgrade@example.com",
		"password": "Developer123!",
	}, "")
	require.Equal(t, http.StatusCreated, registerResp.Code)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(registerResp.Body.Bytes(), &payload))
	userToken := payload["data"].(map[string]any)["token"].(string)

	applyResp := performJSON(t, app, http.MethodPost, "/api/v1/developer-verifications", map[string]any{
		"verification_type": "personal",
		"real_name":         "Developer Upgrade",
		"organization":      "EIBotHub Lab",
		"materials":         "portfolio",
		"reason":            "申请开发者身份以访问内部资源。",
	}, userToken)
	require.Equal(t, http.StatusCreated, applyResp.Code)
	verificationID := extractID(t, applyResp.Body.Bytes())

	reviewResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/verifications/"+verificationID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "认证通过",
	}, adminToken)
	require.Equal(t, http.StatusOK, reviewResp.Code)

	meResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me", nil, userToken, "")
	require.Equal(t, http.StatusOK, meResp.Code)
	require.Contains(t, meResp.Body.String(), "\"role\":\"developer\"")
}

func TestV1ModelAndDatasetCommentsRatings(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	modelRating := performJSON(t, app, http.MethodPost, "/api/v1/models/1/ratings", map[string]any{
		"score":    4,
		"feedback": "模型结构清晰，适合继续优化。",
	}, userToken)
	require.Equal(t, http.StatusOK, modelRating.Code)

	modelComment := performJSON(t, app, http.MethodPost, "/api/v1/models/1/comments", map[string]any{
		"content": "这个模型在仓储场景里比较容易起步。",
	}, userToken)
	require.Equal(t, http.StatusCreated, modelComment.Code)

	modelComments := performRequest(t, app, http.MethodGet, "/api/v1/models/1/comments", nil, "", "")
	require.Equal(t, http.StatusOK, modelComments.Code)
	require.Contains(t, modelComments.Body.String(), "这个模型在仓储场景里比较容易起步。")

	datasetRating := performJSON(t, app, http.MethodPost, "/api/v1/datasets/1/ratings", map[string]any{
		"score":    5,
		"feedback": "数据集适合作为入门样本。",
	}, userToken)
	require.Equal(t, http.StatusOK, datasetRating.Code)

	datasetComment := performJSON(t, app, http.MethodPost, "/api/v1/datasets/1/comments", map[string]any{
		"content": "预览样本足够用于快速判断数据质量。",
	}, userToken)
	require.Equal(t, http.StatusCreated, datasetComment.Code)

	datasetComments := performRequest(t, app, http.MethodGet, "/api/v1/datasets/1/comments", nil, "", "")
	require.Equal(t, http.StatusOK, datasetComments.Code)
	require.Contains(t, datasetComments.Body.String(), "预览样本足够用于快速判断数据质量。")
}

func TestV1AdminCommunityModeration(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	overviewResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/community/overview", nil, adminToken, "")
	require.Equal(t, http.StatusOK, overviewResp.Code)
	require.Contains(t, overviewResp.Body.String(), "skills")

	skillsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/community/skills", nil, adminToken, "")
	require.Equal(t, http.StatusOK, skillsResp.Code)
	require.Contains(t, skillsResp.Body.String(), "巡检告警上报技能")

	hideSkillResp := performRequest(t, app, http.MethodPost, "/api/v1/admin/community/skills/1/hide", nil, adminToken, "")
	require.Equal(t, http.StatusOK, hideSkillResp.Code)

	discussionsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/community/discussions", nil, adminToken, "")
	require.Equal(t, http.StatusOK, discussionsResp.Code)
	require.Contains(t, discussionsResp.Body.String(), "移动作业模板如何接入现有机器人")

	commentsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/community/comments", nil, adminToken, "")
	require.Equal(t, http.StatusOK, commentsResp.Code)
	require.Contains(t, commentsResp.Body.String(), "这个模板适合做首个巡检 MVP")
}

func TestV1SearchIncludesSkillsAndDiscussions(t *testing.T) {
	app := newTestApp(t)

	skillSearch := performRequest(t, app, http.MethodGet, "/api/v1/search?q=告警&type=skill", nil, "", "")
	require.Equal(t, http.StatusOK, skillSearch.Code)
	require.Contains(t, skillSearch.Body.String(), "/skills/")

	discussionSearch := performRequest(t, app, http.MethodGet, "/api/v1/search?q=模板&type=discussion", nil, "", "")
	require.Equal(t, http.StatusOK, discussionSearch.Code)
	require.Contains(t, discussionSearch.Body.String(), "/discussions/")
}

func TestV1UserContributionEndpoints(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	myContributions := performRequest(t, app, http.MethodGet, "/api/v1/users/me/contributions", nil, userToken, "")
	require.Equal(t, http.StatusOK, myContributions.Code)
	require.Contains(t, myContributions.Body.String(), "巡检告警上报技能")

	publicContributions := performRequest(t, app, http.MethodGet, "/api/v1/community/users/2/contributions", nil, "", "")
	require.Equal(t, http.StatusOK, publicContributions.Code)
	require.Contains(t, publicContributions.Body.String(), "移动作业模板如何接入现有机器人")
}

func TestV1FollowCreatesNotificationAndFollowersList(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	followResp := performRequest(t, app, http.MethodPost, "/api/v1/community/users/1/follow", nil, userToken, "")
	require.Equal(t, http.StatusOK, followResp.Code)

	followersResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/followers", nil, adminToken, "")
	require.Equal(t, http.StatusOK, followersResp.Code)
	require.Contains(t, followersResp.Body.String(), "demo")

	notificationsResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/notifications", nil, adminToken, "")
	require.Equal(t, http.StatusOK, notificationsResp.Code)
	require.Contains(t, notificationsResp.Body.String(), "你有新的关注者")
}

func TestV2DeveloperVerificationFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	applyResp := performJSON(t, app, http.MethodPost, "/api/v1/developer-verifications", map[string]any{
		"verification_type": "personal",
		"real_name":         "Demo Developer",
		"organization":      "Open Community Lab",
		"materials":         "GitHub profile and project links",
		"reason":            "希望以开发者身份分享技能和模型。",
	}, userToken)
	require.Equal(t, http.StatusCreated, applyResp.Code)
	verificationID := extractID(t, applyResp.Body.Bytes())

	myStatus := performRequest(t, app, http.MethodGet, "/api/v1/users/me/verification", nil, userToken, "")
	require.Equal(t, http.StatusOK, myStatus.Code)
	require.Contains(t, myStatus.Body.String(), "pending")

	listResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/verifications", nil, adminToken, "")
	require.Equal(t, http.StatusOK, listResp.Code)
	require.Contains(t, listResp.Body.String(), "Demo Developer")

	reviewResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/verifications/"+verificationID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "认证通过",
	}, adminToken)
	require.Equal(t, http.StatusOK, reviewResp.Code)

	publicStatus := performRequest(t, app, http.MethodGet, "/api/v1/community/users/2/verification", nil, "", "")
	require.Equal(t, http.StatusOK, publicStatus.Code)
	require.Contains(t, publicStatus.Body.String(), "approved")
}

func TestV2EnterpriseVerificationEndpoint(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	applyResp := performJSON(t, app, http.MethodPost, "/api/v1/enterprise-verifications", map[string]any{
		"real_name":    "Open Community Inc.",
		"organization": "Open Community Inc.",
		"materials":    "营业执照与官网地址",
		"reason":       "希望以企业身份发布模型与解决方案。",
	}, userToken)
	require.Equal(t, http.StatusCreated, applyResp.Code)
	require.Contains(t, applyResp.Body.String(), "\"verification_type\":\"enterprise\"")
}

func TestV2WikiFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	createResp := performJSON(t, app, http.MethodPost, "/api/v1/wiki/pages", map[string]any{
		"title":   "机器人告警排查 Wiki",
		"summary": "整理常见告警排查步骤。",
		"content": "先核对设备状态，再检查模型依赖和告警回传配置。",
		"comment": "创建词条",
	}, userToken)
	require.Equal(t, http.StatusCreated, createResp.Code)
	pageID := extractID(t, createResp.Body.Bytes())

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/wiki/pages/"+pageID, map[string]any{
		"title":   "机器人告警排查 Wiki",
		"summary": "整理常见告警排查步骤。",
		"content": "先核对设备状态，再检查模型依赖和告警回传配置，最后核验通知链路。",
		"comment": "补充通知链路",
	}, userToken)
	require.Equal(t, http.StatusOK, updateResp.Code)

	pageResp := performRequest(t, app, http.MethodGet, "/api/v1/wiki/pages/"+pageID, nil, "", "")
	require.Equal(t, http.StatusOK, pageResp.Code)
	require.Contains(t, pageResp.Body.String(), "通知链路")

	revisionsResp := performRequest(t, app, http.MethodGet, "/api/v1/wiki/pages/"+pageID+"/revisions", nil, "", "")
	require.Equal(t, http.StatusOK, revisionsResp.Code)
	require.Contains(t, revisionsResp.Body.String(), "创建词条")
	require.Contains(t, revisionsResp.Body.String(), "补充通知链路")
}

func TestV2AdminContentOperations(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	templateResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/content/templates", map[string]any{
		"name":         "夜间巡检模板",
		"summary":      "用于夜间巡检任务编排。",
		"description":  "包含夜间巡检路线、告警和回充步骤。",
		"category":     "巡逻巡检",
		"scene":        "夜间作业",
		"guide":        "绑定设备 -> 选择路线 -> 发布执行",
		"resource_ref": "Warehouse Carrier Net",
		"status":       "draft",
	}, adminToken)
	require.Equal(t, http.StatusCreated, templateResp.Code)
	templateID := extractID(t, templateResp.Body.Bytes())

	batchTemplateResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/content/templates/status", map[string]any{
		"ids":    []int{mustAtoi(t, templateID)},
		"status": "published",
	}, adminToken)
	require.Equal(t, http.StatusOK, batchTemplateResp.Code)

	faqResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/content/faqs", map[string]any{
		"question": "如何维护推荐位？",
		"answer":   "进入后台门户运营页，选择资源后加入推荐位。",
	}, adminToken)
	require.Equal(t, http.StatusCreated, faqResp.Code)
	faqID := extractID(t, faqResp.Body.Bytes())

	deleteFAQResp := performRequest(t, app, http.MethodDelete, "/api/v1/admin/content/faqs/"+faqID, nil, adminToken, "")
	require.Equal(t, http.StatusOK, deleteFAQResp.Code)

	videoResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/content/videos", map[string]any{
		"title":      "内容运营演示",
		"summary":    "演示如何维护视频教程目录。",
		"link":       "https://example.com/videos/admin-content",
		"category":   "后台运营",
		"sort_order": 8,
		"active":     true,
	}, adminToken)
	require.Equal(t, http.StatusCreated, videoResp.Code)
	videoID := extractID(t, videoResp.Body.Bytes())

	batchVideoResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/content/videos/status", map[string]any{
		"ids":    []int{mustAtoi(t, videoID)},
		"active": false,
	}, adminToken)
	require.Equal(t, http.StatusOK, batchVideoResp.Code)

	logsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/operations", nil, adminToken, "")
	require.Equal(t, http.StatusOK, logsResp.Code)
	require.Contains(t, logsResp.Body.String(), "批量更新模板状态")
	require.Contains(t, logsResp.Body.String(), "如何维护推荐位？")
}

func TestV2WikiGovernanceFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	createResp := performJSON(t, app, http.MethodPost, "/api/v1/wiki/pages", map[string]any{
		"title":   "可回滚 Wiki",
		"summary": "验证锁定和回滚。",
		"content": "初始内容",
		"comment": "初始创建",
	}, userToken)
	require.Equal(t, http.StatusCreated, createResp.Code)
	pageID := extractID(t, createResp.Body.Bytes())

	revisionsResp := performRequest(t, app, http.MethodGet, "/api/v1/wiki/pages/"+pageID+"/revisions", nil, "", "")
	require.Equal(t, http.StatusOK, revisionsResp.Code)
	baseRevisionID := extractFirstListIDNumber(t, revisionsResp.Body.Bytes())

	updateResp := performJSON(t, app, http.MethodPut, "/api/v1/wiki/pages/"+pageID, map[string]any{
		"title":   "可回滚 Wiki",
		"summary": "验证锁定和回滚。",
		"content": "第二版内容",
		"comment": "第二版",
	}, userToken)
	require.Equal(t, http.StatusOK, updateResp.Code)

	lockResp := performRequest(t, app, http.MethodPost, "/api/v1/admin/wiki/pages/"+pageID+"/lock", nil, adminToken, "")
	require.Equal(t, http.StatusOK, lockResp.Code)

	blockedUpdateResp := performJSON(t, app, http.MethodPut, "/api/v1/wiki/pages/"+pageID, map[string]any{
		"title":   "可回滚 Wiki",
		"summary": "验证锁定和回滚。",
		"content": "普通用户不应修改成功",
		"comment": "锁定后修改",
	}, userToken)
	require.Equal(t, http.StatusLocked, blockedUpdateResp.Code)
	require.Contains(t, blockedUpdateResp.Body.String(), "wiki_locked")

	adminPagesResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/wiki/pages", nil, adminToken, "")
	require.Equal(t, http.StatusOK, adminPagesResp.Code)
	require.Contains(t, adminPagesResp.Body.String(), "\"locked\":true")

	rollbackResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/wiki/pages/"+pageID+"/rollback", map[string]any{
		"revision_id": baseRevisionID,
		"comment":     "回滚到初始版本",
	}, adminToken)
	require.Equal(t, http.StatusOK, rollbackResp.Code)

	pageResp := performRequest(t, app, http.MethodGet, "/api/v1/wiki/pages/"+pageID, nil, "", "")
	require.Equal(t, http.StatusOK, pageResp.Code)
	require.Contains(t, pageResp.Body.String(), "初始内容")

	unlockResp := performRequest(t, app, http.MethodPost, "/api/v1/admin/wiki/pages/"+pageID+"/unlock", nil, adminToken, "")
	require.Equal(t, http.StatusOK, unlockResp.Code)

	finalUpdateResp := performJSON(t, app, http.MethodPut, "/api/v1/wiki/pages/"+pageID, map[string]any{
		"title":   "可回滚 Wiki",
		"summary": "验证锁定和回滚。",
		"content": "解锁后再次修改",
		"comment": "恢复编辑",
	}, userToken)
	require.Equal(t, http.StatusOK, finalUpdateResp.Code)
}

func TestV2RewardsEndpoints(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	pointsResp := performRequest(t, app, http.MethodGet, "/api/v1/rewards/points", nil, userToken, "")
	require.Equal(t, http.StatusOK, pointsResp.Code)
	require.Contains(t, pointsResp.Body.String(), "points")

	ledgerResp := performRequest(t, app, http.MethodGet, "/api/v1/rewards/ledger", nil, userToken, "")
	require.Equal(t, http.StatusOK, ledgerResp.Code)
	require.Contains(t, ledgerResp.Body.String(), "source_type")

	benefitsResp := performRequest(t, app, http.MethodGet, "/api/v1/rewards/benefits", nil, userToken, "")
	require.Equal(t, http.StatusOK, benefitsResp.Code)
	require.Contains(t, benefitsResp.Body.String(), "首页创作者标识")

	redeemResp := performJSON(t, app, http.MethodPost, "/api/v1/rewards/redeem", map[string]any{
		"benefit_id": 1,
	}, userToken)
	require.Equal(t, http.StatusCreated, redeemResp.Code)
	require.Contains(t, redeemResp.Body.String(), "首页创作者标识")

	redemptionsResp := performRequest(t, app, http.MethodGet, "/api/v1/rewards/redemptions", nil, userToken, "")
	require.Equal(t, http.StatusOK, redemptionsResp.Code)
	require.Contains(t, redemptionsResp.Body.String(), "首页创作者标识")

	pointsAfterRedeem := performRequest(t, app, http.MethodGet, "/api/v1/rewards/points", nil, userToken, "")
	require.Equal(t, http.StatusOK, pointsAfterRedeem.Code)
	require.Contains(t, pointsAfterRedeem.Body.String(), "\"points\":20")

	rankingResp := performRequest(t, app, http.MethodGet, "/api/v1/rankings/contributors", nil, "", "")
	require.Equal(t, http.StatusOK, rankingResp.Code)
	require.Contains(t, rankingResp.Body.String(), "demo")
}

func TestV2CollaborationFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")

	messageResp := performJSON(t, app, http.MethodPost, "/api/v1/messages", map[string]any{
		"recipient_user_id": 1,
		"content":           "你好，我想讨论一下模型联调。",
	}, userToken)
	require.Equal(t, http.StatusCreated, messageResp.Code)
	conversationID := extractDataNumberField(t, messageResp.Body.Bytes(), "conversation_id")

	adminNotificationsResp := performRequest(t, app, http.MethodGet, "/api/v1/users/me/notifications", nil, adminToken, "")
	require.Equal(t, http.StatusOK, adminNotificationsResp.Code)
	require.Contains(t, adminNotificationsResp.Body.String(), "\"Type\":\"message\"")
	require.Contains(t, adminNotificationsResp.Body.String(), "\"Link\":\"/messages?conversation="+strconv.Itoa(conversationID)+"\"")

	conversationsResp := performRequest(t, app, http.MethodGet, "/api/v1/messages/conversations", nil, userToken, "")
	require.Equal(t, http.StatusOK, conversationsResp.Code)
	require.Contains(t, conversationsResp.Body.String(), "direct")
	require.Contains(t, conversationsResp.Body.String(), strconv.Itoa(conversationID))

	replyResp := performJSON(t, app, http.MethodPost, "/api/v1/messages", map[string]any{
		"conversation_id": conversationID,
		"content":         "收到，稍后同步模板依赖。",
	}, userToken)
	require.Equal(t, http.StatusCreated, replyResp.Code)

	workspaceResp := performJSON(t, app, http.MethodPost, "/api/v1/workspaces", map[string]any{
		"name":    "V2 协作空间",
		"summary": "用于私信和协作空间联调。",
	}, userToken)
	require.Equal(t, http.StatusCreated, workspaceResp.Code)
	workspaceID := extractID(t, workspaceResp.Body.Bytes())

	memberResp := performJSON(t, app, http.MethodPost, "/api/v1/workspaces/"+workspaceID+"/members", map[string]any{
		"user_id": 1,
	}, userToken)
	require.Equal(t, http.StatusOK, memberResp.Code)

	workspaceMessage := performJSON(t, app, http.MethodPost, "/api/v1/workspaces/"+workspaceID+"/messages", map[string]any{
		"content": "请一起确认模板依赖。",
	}, userToken)
	require.Equal(t, http.StatusCreated, workspaceMessage.Code)

	detailResp := performRequest(t, app, http.MethodGet, "/api/v1/workspaces/"+workspaceID, nil, userToken, "")
	require.Equal(t, http.StatusOK, detailResp.Code)
	require.Contains(t, detailResp.Body.String(), "请一起确认模板依赖。")
}

func TestV2CollaborationModerationFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")
	memberToken := createRoleUserToken(t, app, "workspace-member", "workspace-member@example.com", "Member123!", model.RoleUser)

	directResp := performJSON(t, app, http.MethodPost, "/api/v1/messages", map[string]any{
		"recipient_user_id": 1,
		"content":           "治理测试：先创建一条私信。",
	}, userToken)
	require.Equal(t, http.StatusCreated, directResp.Code)
	conversationID := extractDataNumberField(t, directResp.Body.Bytes(), "conversation_id")

	blockConversationResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/community/conversations/"+strconv.Itoa(conversationID)+"/block", map[string]any{
		"reason": "命中敏感内容治理策略",
	}, adminToken)
	require.Equal(t, http.StatusOK, blockConversationResp.Code)

	conversationListResp := performRequest(t, app, http.MethodGet, "/api/v1/messages/conversations", nil, userToken, "")
	require.Equal(t, http.StatusOK, conversationListResp.Code)
	require.Contains(t, conversationListResp.Body.String(), "\"active\":false")
	require.Contains(t, conversationListResp.Body.String(), "命中敏感内容治理策略")

	blockedMessagesResp := performRequest(t, app, http.MethodGet, "/api/v1/messages/conversations/"+strconv.Itoa(conversationID), nil, userToken, "")
	require.Equal(t, http.StatusForbidden, blockedMessagesResp.Code)
	require.Contains(t, blockedMessagesResp.Body.String(), "conversation_blocked")

	blockedReplyResp := performJSON(t, app, http.MethodPost, "/api/v1/messages", map[string]any{
		"conversation_id": conversationID,
		"content":         "治理测试：封禁后不应继续回复。",
	}, userToken)
	require.Equal(t, http.StatusForbidden, blockedReplyResp.Code)
	require.Contains(t, blockedReplyResp.Body.String(), "conversation_blocked")

	unblockConversationResp := performRequest(t, app, http.MethodPost, "/api/v1/admin/community/conversations/"+strconv.Itoa(conversationID)+"/unblock", nil, adminToken, "")
	require.Equal(t, http.StatusOK, unblockConversationResp.Code)

	replyResp := performJSON(t, app, http.MethodPost, "/api/v1/messages", map[string]any{
		"conversation_id": conversationID,
		"content":         "治理测试：解封后恢复发送。",
	}, userToken)
	require.Equal(t, http.StatusCreated, replyResp.Code)

	var workspaceMember model.User
	require.NoError(t, app.DB.Where("email = ?", "workspace-member@example.com").First(&workspaceMember).Error)

	workspaceResp := performJSON(t, app, http.MethodPost, "/api/v1/workspaces", map[string]any{
		"name":    "治理测试协作空间",
		"summary": "验证协作空间封禁与成员移除。",
	}, userToken)
	require.Equal(t, http.StatusCreated, workspaceResp.Code)
	workspaceID := extractID(t, workspaceResp.Body.Bytes())

	addMemberResp := performJSON(t, app, http.MethodPost, "/api/v1/workspaces/"+workspaceID+"/members", map[string]any{
		"user_id": workspaceMember.ID,
	}, userToken)
	require.Equal(t, http.StatusOK, addMemberResp.Code)

	blockWorkspaceResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/community/workspaces/"+workspaceID+"/block", map[string]any{
		"reason": "空间内容违规，暂停访问",
	}, adminToken)
	require.Equal(t, http.StatusOK, blockWorkspaceResp.Code)

	memberSpacesResp := performRequest(t, app, http.MethodGet, "/api/v1/workspaces", nil, memberToken, "")
	require.Equal(t, http.StatusOK, memberSpacesResp.Code)
	require.Contains(t, memberSpacesResp.Body.String(), "\"active\":false")
	require.Contains(t, memberSpacesResp.Body.String(), "空间内容违规，暂停访问")

	blockedWorkspaceDetailResp := performRequest(t, app, http.MethodGet, "/api/v1/workspaces/"+workspaceID, nil, memberToken, "")
	require.Equal(t, http.StatusForbidden, blockedWorkspaceDetailResp.Code)
	require.Contains(t, blockedWorkspaceDetailResp.Body.String(), "workspace_blocked")

	blockedWorkspaceMessageResp := performJSON(t, app, http.MethodPost, "/api/v1/workspaces/"+workspaceID+"/messages", map[string]any{
		"content": "治理测试：封禁后不应继续发送。",
	}, memberToken)
	require.Equal(t, http.StatusForbidden, blockedWorkspaceMessageResp.Code)
	require.Contains(t, blockedWorkspaceMessageResp.Body.String(), "workspace_blocked")

	unblockWorkspaceResp := performRequest(t, app, http.MethodPost, "/api/v1/admin/community/workspaces/"+workspaceID+"/unblock", nil, adminToken, "")
	require.Equal(t, http.StatusOK, unblockWorkspaceResp.Code)

	unblockedWorkspaceDetailResp := performRequest(t, app, http.MethodGet, "/api/v1/workspaces/"+workspaceID, nil, memberToken, "")
	require.Equal(t, http.StatusOK, unblockedWorkspaceDetailResp.Code)
	require.Contains(t, unblockedWorkspaceDetailResp.Body.String(), "治理测试协作空间")

	removeMemberResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/community/workspaces/"+workspaceID+"/remove-member", map[string]any{
		"user_id": workspaceMember.ID,
	}, adminToken)
	require.Equal(t, http.StatusOK, removeMemberResp.Code)

	removedMemberDetailResp := performRequest(t, app, http.MethodGet, "/api/v1/workspaces/"+workspaceID, nil, memberToken, "")
	require.Equal(t, http.StatusForbidden, removedMemberDetailResp.Code)
	require.Contains(t, removedMemberDetailResp.Body.String(), "\"code\":\"forbidden\"")

	adminLogsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/operations", nil, adminToken, "")
	require.Equal(t, http.StatusOK, adminLogsResp.Code)
	require.Contains(t, adminLogsResp.Body.String(), "\"resource_type\":\"conversation\"")
	require.Contains(t, adminLogsResp.Body.String(), "\"resource_type\":\"workspace\"")
	require.Contains(t, adminLogsResp.Body.String(), "remove_member")
}

func TestV2OpenAPIAndWebhookFlow(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	specResp := performRequest(t, app, http.MethodGet, "/api/v1/openapi/spec", nil, "", "")
	require.Equal(t, http.StatusOK, specResp.Code)
	require.Contains(t, specResp.Body.String(), "/api/v1/webhooks")

	var (
		mu          sync.Mutex
		events      []string
		signatures  []string
		requestBody []string
	)
	receiver := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mu.Lock()
		events = append(events, r.Header.Get("X-OpenCommunity-Event"))
		signatures = append(signatures, r.Header.Get("X-OpenCommunity-Signature"))
		requestBody = append(requestBody, string(body))
		mu.Unlock()
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"received":true}`))
	}))
	defer receiver.Close()

	createWebhookResp := performJSON(t, app, http.MethodPost, "/api/v1/webhooks", map[string]any{
		"name":       "本地联调",
		"target_url": receiver.URL,
		"secret":     "opencommunity-secret",
		"events":     []string{"webhook.test", "reward.redeemed"},
	}, userToken)
	require.Equal(t, http.StatusCreated, createWebhookResp.Code)
	webhookID := extractID(t, createWebhookResp.Body.Bytes())

	listResp := performRequest(t, app, http.MethodGet, "/api/v1/webhooks", nil, userToken, "")
	require.Equal(t, http.StatusOK, listResp.Code)
	require.Contains(t, listResp.Body.String(), receiver.URL)

	testResp := performRequest(t, app, http.MethodPost, "/api/v1/webhooks/"+webhookID+"/test", nil, userToken, "")
	require.Equal(t, http.StatusCreated, testResp.Code)
	require.Contains(t, testResp.Body.String(), "webhook.test")

	redeemResp := performJSON(t, app, http.MethodPost, "/api/v1/rewards/redeem", map[string]any{
		"benefit_id": 1,
	}, userToken)
	require.Equal(t, http.StatusCreated, redeemResp.Code)

	deliveriesResp := performRequest(t, app, http.MethodGet, "/api/v1/webhooks/"+webhookID+"/deliveries", nil, userToken, "")
	require.Equal(t, http.StatusOK, deliveriesResp.Code)
	require.Contains(t, deliveriesResp.Body.String(), "reward.redeemed")

	mu.Lock()
	defer mu.Unlock()
	require.GreaterOrEqual(t, len(events), 2)
	require.Contains(t, events, "webhook.test")
	require.Contains(t, events, "reward.redeemed")
	require.NotEmpty(t, signatures[0])
	require.Contains(t, strings.Join(requestBody, "\n"), "benefit_name")
}

func TestV2AdminRewardOperations(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	overviewResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/rewards/overview", nil, adminToken, "")
	require.Equal(t, http.StatusOK, overviewResp.Code)
	require.Contains(t, overviewResp.Body.String(), "active_benefits")

	benefitsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/rewards/benefits", nil, adminToken, "")
	require.Equal(t, http.StatusOK, benefitsResp.Code)
	require.Contains(t, benefitsResp.Body.String(), "首页创作者标识")

	updateBenefitResp := performJSON(t, app, http.MethodPut, "/api/v1/admin/rewards/benefits/1", map[string]any{
		"name":        "首页旗舰创作者标识",
		"summary":     "兑换后可在首页与公开开发者页显示旗舰创作者标识。",
		"cost_points": 25,
		"active":      true,
	}, adminToken)
	require.Equal(t, http.StatusOK, updateBenefitResp.Code)
	require.Contains(t, updateBenefitResp.Body.String(), "首页旗舰创作者标识")

	adjustResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/rewards/adjustments", map[string]any{
		"user_id": 2,
		"points":  15,
		"remark":  "活动补贴",
	}, adminToken)
	require.Equal(t, http.StatusCreated, adjustResp.Code)
	require.Contains(t, adjustResp.Body.String(), "活动补贴")

	adjustmentsResp := performRequest(t, app, http.MethodGet, "/api/v1/admin/rewards/adjustments", nil, adminToken, "")
	require.Equal(t, http.StatusOK, adjustmentsResp.Code)
	require.Contains(t, adjustmentsResp.Body.String(), "活动补贴")

	pointsResp := performRequest(t, app, http.MethodGet, "/api/v1/rewards/points", nil, userToken, "")
	require.Equal(t, http.StatusOK, pointsResp.Code)
	require.Contains(t, pointsResp.Body.String(), "\"points\":55")
}

func newTestApp(t *testing.T) *App {
	t.Helper()
	tempDir := t.TempDir()
	app, err := NewForTest(config.Config{
		Port:       "0",
		DBPath:     filepath.Join(tempDir, "test.db"),
		StorageDir: filepath.Join(tempDir, "storage"),
		AppSecret:  "test-secret",
		SeedDemo:   true,
	})
	require.NoError(t, err)
	sqlDB, err := app.DB.DB()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})
	return app
}

func performJSON(t *testing.T, app *App, method, target string, body any, token string) *httptest.ResponseRecorder {
	t.Helper()
	payload, err := json.Marshal(body)
	require.NoError(t, err)
	return performRequest(t, app, method, target, bytes.NewReader(payload), token, "application/json")
}

func performMultipart(t *testing.T, app *App, method, target string, fields map[string]string, fileField, fileName, fileContent, token string) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for key, value := range fields {
		require.NoError(t, writer.WriteField(key, value))
	}
	fileWriter, err := writer.CreateFormFile(fileField, fileName)
	require.NoError(t, err)
	_, err = io.Copy(fileWriter, strings.NewReader(fileContent))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return performRequest(t, app, method, target, &body, token, writer.FormDataContentType())
}

func performRequest(t *testing.T, app *App, method, target string, body io.Reader, token, contentType string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	return resp
}

func loginToken(t *testing.T, app *App, email, password string) string {
	t.Helper()
	resp := performJSON(t, app, http.MethodPost, "/api/v1/auth/login", map[string]any{
		"email":    email,
		"password": password,
	}, "")
	require.Equal(t, http.StatusOK, resp.Code)
	var payload map[string]any
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &payload))
	return payload["data"].(map[string]any)["token"].(string)
}

func extractID(t *testing.T, raw []byte) string {
	t.Helper()
	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))
	id := payload["data"].(map[string]any)["id"]
	return strconv.FormatInt(int64(id.(float64)), 10)
}

func extractDataNumberField(t *testing.T, raw []byte, field string) int {
	t.Helper()
	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))
	value := payload["data"].(map[string]any)[field]
	return int(value.(float64))
}

func extractFirstListIDNumber(t *testing.T, raw []byte) int {
	t.Helper()
	var payload map[string]any
	require.NoError(t, json.Unmarshal(raw, &payload))
	data := payload["data"].([]any)
	first := data[0].(map[string]any)
	return int(first["id"].(float64))
}

func mustAtoi(t *testing.T, value string) int {
	t.Helper()
	result, err := strconv.Atoi(value)
	require.NoError(t, err)
	return result
}

func createAdminUserToken(t *testing.T, app *App, username, email, password string) string {
	t.Helper()
	return createRoleUserToken(t, app, username, email, password, model.RoleAdmin)
}

func createRoleUserToken(t *testing.T, app *App, username, email, password, role string) string {
	t.Helper()
	resp := performJSON(t, app, http.MethodPost, "/api/v1/auth/register", map[string]any{
		"username": username,
		"email":    email,
		"password": password,
	}, "")
	require.Equal(t, http.StatusCreated, resp.Code)
	require.NoError(t, app.DB.Model(&model.User{}).Where("email = ?", email).Update("role", role).Error)
	return loginToken(t, app, email, password)
}

func assertCountAtLeast(t *testing.T, label string, query *gorm.DB, min int64) {
	t.Helper()
	var count int64
	require.NoError(t, query.Count(&count).Error)
	require.GreaterOrEqualf(t, count, min, "%s count = %d, want >= %d", label, count, min)
}

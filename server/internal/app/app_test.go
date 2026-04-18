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

	"github.com/stretchr/testify/require"

	"github.com/Andras5014/EIBotHub/server/internal/config"
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
}

func TestDatasetSamplePreviewAndDownloadPackages(t *testing.T) {
	app := newTestApp(t)
	userToken := loginToken(t, app, "demo@example.com", "Demo123!")

	samplesResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/1/samples", nil, "", "")
	require.Equal(t, http.StatusOK, samplesResp.Code)
	require.Contains(t, samplesResp.Body.String(), "\"sample_type\":\"image\"")
	require.Contains(t, samplesResp.Body.String(), "\"sample_type\":\"pointcloud\"")

	agreementResp := performRequest(t, app, http.MethodPost, "/api/v1/datasets/1/agreements/confirm", nil, userToken, "")
	require.Equal(t, http.StatusOK, agreementResp.Code)

	taskResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/1/download-packages", map[string]any{
		"parts": 4,
	}, userToken)
	require.Equal(t, http.StatusCreated, taskResp.Code)
	require.Contains(t, taskResp.Body.String(), "\"total_parts\":4")
	require.Contains(t, taskResp.Body.String(), "download-packages")

	taskListResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/1/download-packages", nil, userToken, "")
	require.Equal(t, http.StatusOK, taskListResp.Code)
	require.Contains(t, taskListResp.Body.String(), "\"status\":\"ready\"")
}

func TestRestrictedDatasetAccessApprovalFlow(t *testing.T) {
	app := newTestApp(t)
	adminToken := loginToken(t, app, "admin@opencommunity.local", "Admin123!")
	ownerToken := loginToken(t, app, "demo@example.com", "Demo123!")

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

	adminDecisionResp := performJSON(t, app, http.MethodPost, "/api/v1/admin/datasets/access-requests/"+accessRequestID+"/decision", map[string]any{
		"decision": "approved",
		"comment":  "允许下载",
	}, adminToken)
	require.Equal(t, http.StatusOK, adminDecisionResp.Code)

	downloadAllowed := performRequest(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download", nil, guestToken, "")
	require.Equal(t, http.StatusOK, downloadAllowed.Code)

	packageResp := performJSON(t, app, http.MethodPost, "/api/v1/datasets/"+datasetID+"/download-packages", map[string]any{
		"parts": 2,
	}, guestToken)
	require.Equal(t, http.StatusCreated, packageResp.Code)
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
		"title":    "如何组合技能与模板？",
		"summary":  "想确认技能分享和模板引用的最佳实践。",
		"content":  "当前项目希望把告警技能接入任务模板，想了解推荐的复用方式。",
		"category": "最佳实践",
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

	statsResp := performRequest(t, app, http.MethodGet, "/api/v1/datasets/1/stats", nil, "", "")
	require.Equal(t, http.StatusOK, statsResp.Code)
	require.Contains(t, statsResp.Body.String(), "download_count")
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

	messageResp := performJSON(t, app, http.MethodPost, "/api/v1/messages", map[string]any{
		"recipient_user_id": 1,
		"content":           "你好，我想讨论一下模型联调。",
	}, userToken)
	require.Equal(t, http.StatusCreated, messageResp.Code)

	conversationsResp := performRequest(t, app, http.MethodGet, "/api/v1/messages/conversations", nil, userToken, "")
	require.Equal(t, http.StatusOK, conversationsResp.Code)
	require.Contains(t, conversationsResp.Body.String(), "direct")
	conversationID := extractFirstListIDNumber(t, conversationsResp.Body.Bytes())

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

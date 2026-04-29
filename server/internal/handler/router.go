package handler

import (
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/middleware"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/service"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type Handler struct {
	auth          *service.AuthService
	portal        *service.PortalService
	models        *service.ModelService
	datasets      *service.DatasetService
	catalog       *service.CatalogService
	users         *service.UserService
	search        *service.SearchService
	admin         *service.AdminService
	community     *service.CommunityService
	verification  *service.VerificationService
	wiki          *service.WikiService
	rewards       *service.RewardService
	integrations  *service.IntegrationService
	collaboration *service.CollaborationService
	files         *service.FileService
	storage       support.ObjectStorage
	tokens        *support.TokenManager
}

var (
	registerLabels = map[string]string{
		"Username": "用户名",
		"Email":    "邮箱",
		"Password": "密码",
	}
	loginLabels = map[string]string{
		"Email":    "邮箱",
		"Password": "密码",
	}
	profileLabels = map[string]string{
		"Username": "用户名",
		"Bio":      "个人简介",
		"Avatar":   "头像链接",
	}
	searchLabels = map[string]string{
		"Q": "关键词",
	}
	searchKeywordLabels = map[string]string{
		"Query":       "关键词",
		"KeywordType": "运营词类型",
		"SortOrder":   "排序值",
	}
	filterOptionLabels = map[string]string{
		"Kind":      "字典类型",
		"Value":     "展示值",
		"SortOrder": "排序值",
	}
	modelRecommendTagLabels = map[string]string{
		"RecommendTag": "模型推荐标签",
	}
	batchDatasetAccessDecisionLabels = map[string]string{
		"IDs":      "申请列表",
		"Decision": "审核决定",
		"Comment":  "审核说明",
	}
	modelLabels = map[string]string{
		"Name":         "模型名称",
		"Summary":      "摘要",
		"Description":  "描述",
		"Tags":         "标签",
		"RobotType":    "机器人类型",
		"InputSpec":    "输入规格",
		"OutputSpec":   "输出规格",
		"License":      "许可证",
		"Dependencies": "依赖",
		"Version":      "版本",
		"Changelog":    "版本说明",
	}
	modelUpdateLabels = map[string]string{
		"Name":         "模型名称",
		"Summary":      "摘要",
		"Description":  "描述",
		"Tags":         "标签",
		"RobotType":    "机器人类型",
		"InputSpec":    "输入规格",
		"OutputSpec":   "输出规格",
		"License":      "许可证",
		"Dependencies": "依赖",
	}
	modelVersionLabels = map[string]string{
		"Version":   "版本",
		"Changelog": "版本说明",
	}
	datasetLabels = map[string]string{
		"Name":          "数据集名称",
		"Summary":       "摘要",
		"Description":   "描述",
		"Tags":          "标签",
		"SampleCount":   "样本量",
		"Device":        "采集设备",
		"Scene":         "场景",
		"Privacy":       "权限等级",
		"AgreementText": "协议文本",
		"Version":       "版本",
		"Changelog":     "版本说明",
		"SamplePreview": "样本预览",
	}
	datasetUpdateLabels = map[string]string{
		"Name":          "数据集名称",
		"Summary":       "摘要",
		"Description":   "描述",
		"Tags":          "标签",
		"SampleCount":   "样本量",
		"Device":        "采集设备",
		"Scene":         "场景",
		"Privacy":       "权限等级",
		"AgreementText": "协议文本",
		"SamplePreview": "样本预览",
	}
	favoriteLabels = map[string]string{
		"ResourceType": "资源类型",
		"ResourceID":   "资源ID",
		"Title":        "标题",
	}
	reviewLabels = map[string]string{
		"Decision": "审核决定",
		"Comment":  "审核说明",
	}
	announcementLabels = map[string]string{
		"Title":   "标题",
		"Summary": "摘要",
		"Link":    "链接",
	}
	moduleSettingLabels = map[string]string{
		"Enabled":   "启用状态",
		"SortOrder": "排序值",
	}
	homeHighlightLabels = map[string]string{
		"Text":      "亮点文案",
		"SortOrder": "排序值",
	}
	homeHeroLabels = map[string]string{
		"Tagline":         "导读标签",
		"Title":           "主标题",
		"Description":     "导读说明",
		"PrimaryButton":   "主按钮文案",
		"SecondaryButton": "次按钮文案",
		"SearchButton":    "搜索按钮文案",
	}
	rankingConfigLabels = map[string]string{
		"Title":    "榜单标题",
		"Subtitle": "榜单副标题",
		"Limit":    "展示数量",
	}
	scenePageLabels = map[string]string{
		"Slug":        "场景标识",
		"Name":        "场景名称",
		"Tagline":     "场景标签",
		"Summary":     "场景摘要",
		"Description": "场景说明",
		"SortOrder":   "排序值",
	}
	featuredResourceLabels = map[string]string{
		"ResourceType": "资源类型",
		"ResourceID":   "资源",
		"BadgeLabel":   "推荐标签",
		"SortOrder":    "排序",
	}
	templateAdminLabels = map[string]string{
		"Name":        "模板名称",
		"Summary":     "摘要",
		"Description": "描述",
		"Category":    "分类",
		"Scene":       "场景",
		"Guide":       "指南",
		"ResourceRef": "关联资源",
		"Status":      "状态",
	}
	applicationCaseAdminLabels = map[string]string{
		"Title":      "案例标题",
		"Summary":    "摘要",
		"Category":   "分类",
		"Guide":      "部署指南",
		"CoverImage": "封面图",
		"Status":     "状态",
	}
	documentCategoryLabels = map[string]string{
		"Name":    "分类名称",
		"DocType": "文档类型",
	}
	documentAdminLabels = map[string]string{
		"CategoryID": "分类",
		"Title":      "标题",
		"Summary":    "摘要",
		"Content":    "正文",
		"DocType":    "文档类型",
		"Status":     "状态",
	}
	faqAdminLabels = map[string]string{
		"Question": "问题",
		"Answer":   "答案",
	}
	videoTutorialLabels = map[string]string{
		"Title":     "标题",
		"Summary":   "摘要",
		"Link":      "链接",
		"Category":  "分类",
		"SortOrder": "排序",
	}
	agreementTemplateLabels = map[string]string{
		"Name":      "模板名称",
		"Content":   "模板内容",
		"SortOrder": "排序",
	}
	privacyOptionLabels = map[string]string{
		"Code":        "权限代码",
		"Name":        "权限名称",
		"Description": "说明",
		"SortOrder":   "排序",
	}
	modelEvaluationLabels = map[string]string{
		"Benchmark": "基准名称",
		"Summary":   "评测摘要",
		"Score":     "评分",
		"Notes":     "评测说明",
	}
	ratingLabels = map[string]string{
		"Score":    "评分",
		"Feedback": "反馈",
	}
	commentLabels = map[string]string{
		"Content": "评论内容",
	}
	skillLabels = map[string]string{
		"Name":        "技能名称",
		"Summary":     "摘要",
		"Description": "描述",
		"Category":    "分类",
		"Scene":       "场景",
		"Guide":       "指南",
		"ResourceRef": "关联资源",
	}
	discussionLabels = map[string]string{
		"Title":   "标题",
		"Tag":     "标签",
		"Content": "正文",
	}
	verificationLabels = map[string]string{
		"VerificationType": "认证类型",
		"RealName":         "真实姓名",
		"Organization":     "组织名称",
		"Materials":        "材料说明",
		"Reason":           "申请理由",
	}
	verificationDecisionLabels = map[string]string{
		"Decision": "审核决定",
		"Comment":  "审核说明",
	}
	batchStatusLabels = map[string]string{
		"IDs":    "资源列表",
		"Status": "状态",
	}
	batchDeleteLabels = map[string]string{
		"IDs": "资源列表",
	}
	rewardRedeemLabels = map[string]string{
		"BenefitID": "权益ID",
	}
	webhookLabels = map[string]string{
		"Name":      "名称",
		"TargetURL": "回调地址",
		"Secret":    "签名密钥",
		"Events":    "事件列表",
	}
	adminRewardBenefitLabels = map[string]string{
		"Name":       "权益名称",
		"Summary":    "权益摘要",
		"CostPoints": "所需积分",
	}
	adminRewardAdjustmentLabels = map[string]string{
		"UserID": "用户",
		"Points": "积分变更值",
		"Remark": "变更说明",
	}
	wikiRollbackLabels = map[string]string{
		"RevisionID": "修订版本",
		"Comment":    "回滚说明",
	}
	messageLabels = map[string]string{
		"Content": "消息内容",
	}
	workspaceLabels = map[string]string{
		"Name":    "协作空间名称",
		"Summary": "协作空间摘要",
	}
)

func New(auth *service.AuthService, portal *service.PortalService, models *service.ModelService, datasets *service.DatasetService, catalog *service.CatalogService, users *service.UserService, search *service.SearchService, admin *service.AdminService, community *service.CommunityService, verification *service.VerificationService, wiki *service.WikiService, rewards *service.RewardService, integrations *service.IntegrationService, collaboration *service.CollaborationService, files *service.FileService, storage support.ObjectStorage, tokens *support.TokenManager) *Handler {
	return &Handler{
		auth:          auth,
		portal:        portal,
		models:        models,
		datasets:      datasets,
		catalog:       catalog,
		users:         users,
		search:        search,
		admin:         admin,
		community:     community,
		verification:  verification,
		wiki:          wiki,
		rewards:       rewards,
		integrations:  integrations,
		collaboration: collaboration,
		files:         files,
		storage:       storage,
		tokens:        tokens,
	}
}

func (h *Handler) Register(router *gin.Engine) {
	router.Use(gin.Recovery(), middleware.CORS())

	api := router.Group("/api/v1")
	{
		api.GET("/files/download/:token", h.downloadFileByToken)
		api.GET("/portal/home", h.getHome)
		api.GET("/scenes", h.listScenePages)
		api.GET("/scenes/:slug", h.scenePageDetail)
		api.GET("/search", h.searchAll)
		api.GET("/search/hot", h.hotQueries)
		api.GET("/search/recommended", h.recommendedQueries)
		api.POST("/auth/register", h.register)
		api.POST("/auth/login", h.login)

		api.GET("/models", h.listModels)
		api.GET("/models/:id", h.modelDetail)
		api.GET("/models/:id/evaluations", h.listModelEvaluations)
		api.GET("/models/:id/ratings", h.getModelRatings)
		api.GET("/models/:id/comments", h.getModelComments)
		api.GET("/datasets/options", h.datasetOptions)
		api.GET("/filter-options", h.filterOptions)
		api.GET("/datasets", h.listDatasets)
		api.GET("/datasets/:id", h.datasetDetail)
		api.GET("/datasets/:id/samples", h.datasetSamples)
		api.GET("/datasets/:id/stats", h.datasetStats)
		api.GET("/datasets/:id/ratings", h.getDatasetRatings)
		api.GET("/datasets/:id/comments", h.getDatasetComments)
		api.GET("/task-templates", h.listTemplates)
		api.GET("/task-templates/:id", h.templateDetail)
		api.GET("/task-templates/:id/ratings", h.getTemplateRatings)
		api.GET("/task-templates/:id/comments", h.getTemplateComments)
		api.GET("/application-cases", h.listApplicationCases)
		api.GET("/application-cases/:id", h.applicationCaseDetail)
		api.GET("/skills", h.listSkills)
		api.GET("/skills/:id", h.skillDetail)
		api.GET("/skills/:id/ratings", h.getSkillRatings)
		api.GET("/skills/:id/comments", h.getSkillComments)
		api.GET("/docs/categories", h.listDocCategories)
		api.GET("/docs", h.listDocs)
		api.GET("/docs/:id", h.docDetail)
		api.GET("/faqs", h.listFAQs)
		api.GET("/videos", h.listVideos)
		api.GET("/community/users/:id", h.publicUserProfile)
		api.GET("/community/users/:id/contributions", h.publicUserContributions)
		api.GET("/community/users/:id/follow-stats", h.publicUserFollowStats)
		api.GET("/community/users/:id/verification", h.publicUserVerification)
		api.GET("/wiki/pages", h.listWikiPages)
		api.GET("/wiki/pages/:id", h.wikiPageDetail)
		api.GET("/wiki/pages/:id/revisions", h.listWikiRevisions)
		api.GET("/openapi/spec", h.openAPISpec)
		api.GET("/rankings/contributors", h.contributorRankings)
		api.GET("/discussions", h.listDiscussions)
		api.GET("/discussions/:id", h.discussionDetail)
		api.GET("/discussions/:id/comments", h.getDiscussionComments)

		secured := api.Group("/")
		secured.Use(middleware.AuthRequired(h.tokens))
		{
			secured.POST("/auth/logout", h.logout)
			secured.GET("/users/me", h.me)
			secured.GET("/users/me/profile", h.profile)
			secured.GET("/users/me/contributions", h.myContributions)
			secured.GET("/users/me/follow-stats", h.myFollowStats)
			secured.GET("/users/me/verification", h.myVerification)
			secured.GET("/rewards/points", h.myRewardSummary)
			secured.GET("/rewards/ledger", h.myRewardLedger)
			secured.GET("/rewards/benefits", h.rewardBenefits)
			secured.GET("/rewards/redemptions", h.myRewardRedemptions)
			secured.POST("/rewards/redeem", h.redeemRewardBenefit)
			secured.GET("/webhooks", h.listWebhooks)
			secured.POST("/webhooks", h.createWebhook)
			secured.GET("/webhooks/:id/deliveries", h.listWebhookDeliveries)
			secured.POST("/webhooks/:id/test", h.testWebhook)
			secured.PUT("/users/me/profile", h.updateProfile)
			secured.GET("/users/me/uploads", h.myUploads)
			secured.GET("/users/me/dataset-access-requests", h.myDatasetAccessRequests)
			secured.GET("/users/me/favorites", h.myFavorites)
			secured.POST("/favorites/toggle", h.toggleFavorite)
			secured.GET("/users/me/downloads", h.myDownloads)
			secured.GET("/users/me/notifications", h.myNotifications)
			secured.POST("/users/me/notifications/read", h.readNotifications)
			secured.POST("/users/me/notifications/:id/read", h.readNotification)

			secured.POST("/models", h.createModel)
			secured.PUT("/models/:id", h.updateModel)
			secured.POST("/models/:id/versions", h.addModelVersion)
			secured.POST("/models/:id/submit", h.submitModel)
			secured.POST("/models/:id/download", h.downloadModel)
			secured.POST("/models/:id/evaluations", h.createModelEvaluation)
			secured.POST("/models/:id/ratings", h.rateModel)
			secured.POST("/models/:id/comments", h.commentModel)

			secured.POST("/datasets", h.createDataset)
			secured.PUT("/datasets/:id", h.updateDataset)
			secured.POST("/datasets/:id/versions", h.addDatasetVersion)
			secured.POST("/datasets/:id/submit", h.submitDataset)
			secured.POST("/datasets/:id/agreements/confirm", h.confirmDatasetAgreement)
			secured.POST("/datasets/:id/download", h.downloadDataset)
			secured.GET("/datasets/:id/download-packages", h.listDatasetDownloadPackages)
			secured.POST("/datasets/:id/download-packages", h.createDatasetDownloadPackage)
			secured.GET("/datasets/:id/access-requests/me", h.myDatasetAccessRequest)
			secured.GET("/datasets/:id/access-requests/history", h.myDatasetAccessHistory)
			secured.POST("/datasets/:id/access-requests", h.createDatasetAccessRequest)
			secured.POST("/datasets/:id/ratings", h.rateDataset)
			secured.POST("/datasets/:id/comments", h.commentDataset)

			secured.POST("/task-templates/:id/ratings", h.rateTemplate)
			secured.POST("/task-templates/:id/comments", h.commentTemplate)
			secured.POST("/skills", h.createSkill)
			secured.POST("/skills/:id/fork", h.forkSkill)
			secured.POST("/skills/:id/ratings", h.rateSkill)
			secured.POST("/skills/:id/comments", h.commentSkill)
			secured.POST("/discussions", h.createDiscussion)
			secured.POST("/discussions/:id/comments", h.commentDiscussion)
			secured.POST("/community/users/:id/follow", h.toggleFollow)
			secured.GET("/users/me/follows", h.myFollows)
			secured.GET("/users/me/followers", h.myFollowers)
			secured.POST("/developer-verifications", h.applyVerification)
			secured.POST("/enterprise-verifications", h.applyEnterpriseVerification)
			secured.POST("/wiki/pages", h.createWikiPage)
			secured.PUT("/wiki/pages/:id", h.updateWikiPage)
			secured.GET("/messages/conversations", h.listConversations)
			secured.GET("/messages/conversations/:id", h.listConversationMessages)
			secured.POST("/messages", h.sendMessage)
			secured.GET("/workspaces", h.listWorkspaces)
			secured.GET("/workspaces/:id", h.workspaceDetail)
			secured.POST("/workspaces", h.createWorkspace)
			secured.POST("/workspaces/:id/members", h.addWorkspaceMember)
			secured.POST("/workspaces/:id/messages", h.sendWorkspaceMessage)
		}

		adminGroup := api.Group("/admin")
		adminGroup.Use(middleware.AuthRequired(h.tokens))
		{
			dashboardGroup := adminGroup.Group("/")
			dashboardGroup.Use(middleware.PermissionRequired(model.PermissionDashboardView))
			{
				dashboardGroup.GET("/dashboard", h.adminDashboard)
			}

			operationGroup := adminGroup.Group("/")
			operationGroup.Use(middleware.PermissionRequired(model.PermissionOperationLogView))
			{
				operationGroup.GET("/operations", h.adminOperationLogs)
			}

			reviewGroup := adminGroup.Group("/")
			reviewGroup.Use(middleware.PermissionRequired(model.PermissionReviewManage))
			{
				reviewGroup.GET("/reviews", h.adminReviews)
				reviewGroup.POST("/reviews/:type/:id/decision", h.adminReviewDecision)
			}

			datasetAccessGroup := adminGroup.Group("/")
			datasetAccessGroup.Use(middleware.PermissionRequired(model.PermissionDatasetAccessReview))
			{
				datasetAccessGroup.GET("/datasets/access-requests", h.adminDatasetAccessRequests)
				datasetAccessGroup.POST("/datasets/access-requests/batch-decision", h.adminBatchReviewDatasetAccessRequests)
				datasetAccessGroup.POST("/datasets/access-requests/:id/decision", h.adminReviewDatasetAccessRequest)
			}

			verificationGroup := adminGroup.Group("/")
			verificationGroup.Use(middleware.PermissionRequired(model.PermissionVerificationReview))
			{
				verificationGroup.GET("/verifications", h.adminVerifications)
				verificationGroup.POST("/verifications/:id/decision", h.adminReviewVerification)
			}

			announcementGroup := adminGroup.Group("/")
			announcementGroup.Use(middleware.PermissionRequired(model.PermissionAnnouncementManage))
			{
				announcementGroup.GET("/announcements", h.adminAnnouncements)
				announcementGroup.POST("/announcements", h.adminCreateAnnouncement)
			}

			portalGroup := adminGroup.Group("/")
			portalGroup.Use(middleware.PermissionRequired(model.PermissionPortalManage))
			{
				portalGroup.GET("/portal/modules", h.adminPortalModules)
				portalGroup.PUT("/portal/modules/:key", h.adminUpdatePortalModule)
				portalGroup.GET("/portal/hero-config", h.adminHomeHeroConfig)
				portalGroup.PUT("/portal/hero-config", h.adminUpdateHomeHeroConfig)
				portalGroup.GET("/portal/highlights", h.adminHomeHighlights)
				portalGroup.POST("/portal/highlights", h.adminCreateHomeHighlight)
				portalGroup.PUT("/portal/highlights/:id", h.adminUpdateHomeHighlight)
				portalGroup.DELETE("/portal/highlights/:id", h.adminDeleteHomeHighlight)
				portalGroup.GET("/portal/scenes", h.adminScenePages)
				portalGroup.POST("/portal/scenes", h.adminCreateScenePage)
				portalGroup.PUT("/portal/scenes/:id", h.adminUpdateScenePage)
				portalGroup.DELETE("/portal/scenes/:id", h.adminDeleteScenePage)
				portalGroup.GET("/portal/rankings-config", h.adminRankingConfig)
				portalGroup.PUT("/portal/rankings-config", h.adminUpdateRankingConfig)
			}

			searchKeywordGroup := adminGroup.Group("/")
			searchKeywordGroup.Use(middleware.PermissionRequired(model.PermissionSearchKeywordManage))
			{
				searchKeywordGroup.GET("/portal/search-keywords", h.adminSearchKeywords)
				searchKeywordGroup.POST("/portal/search-keywords", h.adminCreateSearchKeyword)
				searchKeywordGroup.PUT("/portal/search-keywords/:id", h.adminUpdateSearchKeyword)
				searchKeywordGroup.DELETE("/portal/search-keywords/:id", h.adminDeleteSearchKeyword)
			}

			featuredGroup := adminGroup.Group("/")
			featuredGroup.Use(middleware.PermissionRequired(model.PermissionFeaturedResourceManage))
			{
				featuredGroup.GET("/portal/featured-resources", h.adminFeaturedResources)
				featuredGroup.POST("/portal/featured-resources", h.adminCreateFeaturedResource)
				featuredGroup.PUT("/portal/featured-resources/:id", h.adminUpdateFeaturedResource)
				featuredGroup.DELETE("/portal/featured-resources/:id", h.adminDeleteFeaturedResource)
			}

			templateGroup := adminGroup.Group("/")
			templateGroup.Use(middleware.PermissionRequired(model.PermissionTemplateManage))
			{
				templateGroup.GET("/content/templates", h.adminTemplates)
				templateGroup.POST("/content/templates", h.adminCreateTemplate)
				templateGroup.PUT("/content/templates/:id", h.adminUpdateTemplate)
				templateGroup.DELETE("/content/templates/:id", h.adminDeleteTemplate)
				templateGroup.POST("/content/templates/status", h.adminBatchTemplateStatus)
				templateGroup.POST("/content/templates/delete", h.adminBatchDeleteTemplates)
			}

			caseGroup := adminGroup.Group("/")
			caseGroup.Use(middleware.PermissionRequired(model.PermissionApplicationCaseManage))
			{
				caseGroup.GET("/content/application-cases", h.adminApplicationCases)
				caseGroup.POST("/content/application-cases", h.adminCreateApplicationCase)
				caseGroup.PUT("/content/application-cases/:id", h.adminUpdateApplicationCase)
				caseGroup.DELETE("/content/application-cases/:id", h.adminDeleteApplicationCase)
				caseGroup.POST("/content/application-cases/status", h.adminBatchApplicationCaseStatus)
				caseGroup.POST("/content/application-cases/delete", h.adminBatchDeleteApplicationCases)
			}

			docCategoryGroup := adminGroup.Group("/")
			docCategoryGroup.Use(middleware.PermissionRequired(model.PermissionDocumentCategoryManage))
			{
				docCategoryGroup.GET("/content/doc-categories", h.adminDocCategories)
				docCategoryGroup.POST("/content/doc-categories", h.adminCreateDocCategory)
				docCategoryGroup.PUT("/content/doc-categories/:id", h.adminUpdateDocCategory)
				docCategoryGroup.DELETE("/content/doc-categories/:id", h.adminDeleteDocCategory)
			}

			documentGroup := adminGroup.Group("/")
			documentGroup.Use(middleware.PermissionRequired(model.PermissionDocumentManage))
			{
				documentGroup.GET("/content/docs", h.adminDocuments)
				documentGroup.POST("/content/docs", h.adminCreateDocument)
				documentGroup.PUT("/content/docs/:id", h.adminUpdateDocument)
				documentGroup.DELETE("/content/docs/:id", h.adminDeleteDocument)
				documentGroup.POST("/content/docs/status", h.adminBatchDocumentStatus)
				documentGroup.POST("/content/docs/delete", h.adminBatchDeleteDocuments)
			}

			faqGroup := adminGroup.Group("/")
			faqGroup.Use(middleware.PermissionRequired(model.PermissionFAQManage))
			{
				faqGroup.GET("/content/faqs", h.adminFAQs)
				faqGroup.POST("/content/faqs", h.adminCreateFAQ)
				faqGroup.PUT("/content/faqs/:id", h.adminUpdateFAQ)
				faqGroup.DELETE("/content/faqs/:id", h.adminDeleteFAQ)
			}

			videoGroup := adminGroup.Group("/")
			videoGroup.Use(middleware.PermissionRequired(model.PermissionVideoManage))
			{
				videoGroup.GET("/content/videos", h.adminVideos)
				videoGroup.POST("/content/videos", h.adminCreateVideo)
				videoGroup.PUT("/content/videos/:id", h.adminUpdateVideo)
				videoGroup.DELETE("/content/videos/:id", h.adminDeleteVideo)
				videoGroup.POST("/content/videos/status", h.adminBatchVideoStatus)
				videoGroup.POST("/content/videos/delete", h.adminBatchDeleteVideos)
			}

			agreementGroup := adminGroup.Group("/")
			agreementGroup.Use(middleware.PermissionRequired(model.PermissionAgreementManage))
			{
				agreementGroup.GET("/content/agreement-templates", h.adminAgreementTemplates)
				agreementGroup.POST("/content/agreement-templates", h.adminCreateAgreementTemplate)
				agreementGroup.PUT("/content/agreement-templates/:id", h.adminUpdateAgreementTemplate)
				agreementGroup.DELETE("/content/agreement-templates/:id", h.adminDeleteAgreementTemplate)
			}

			privacyGroup := adminGroup.Group("/")
			privacyGroup.Use(middleware.PermissionRequired(model.PermissionPrivacyManage))
			{
				privacyGroup.GET("/content/privacy-options", h.adminPrivacyOptions)
				privacyGroup.POST("/content/privacy-options", h.adminCreatePrivacyOption)
				privacyGroup.PUT("/content/privacy-options/:id", h.adminUpdatePrivacyOption)
				privacyGroup.DELETE("/content/privacy-options/:id", h.adminDeletePrivacyOption)
			}

			filterOptionGroup := adminGroup.Group("/")
			filterOptionGroup.Use(middleware.PermissionRequired(model.PermissionFilterOptionManage))
			{
				filterOptionGroup.GET("/content/filter-options", h.adminFilterOptions)
				filterOptionGroup.POST("/content/filter-options", h.adminCreateFilterOption)
				filterOptionGroup.PUT("/content/filter-options/:id", h.adminUpdateFilterOption)
				filterOptionGroup.DELETE("/content/filter-options/:id", h.adminDeleteFilterOption)
			}

			modelRecommendTagGroup := adminGroup.Group("/")
			modelRecommendTagGroup.Use(middleware.PermissionRequired(model.PermissionModelRecommendTagManage))
			{
				modelRecommendTagGroup.GET("/content/model-recommend-tags", h.adminModelRecommendTags)
				modelRecommendTagGroup.PUT("/content/model-recommend-tags/:id", h.adminUpdateModelRecommendTag)
			}

			communityOverviewGroup := adminGroup.Group("/")
			communityOverviewGroup.Use(middleware.PermissionRequired(model.PermissionCommunityAccess))
			{
				communityOverviewGroup.GET("/community/overview", h.adminCommunityOverview)
			}

			communityContentGroup := adminGroup.Group("/")
			communityContentGroup.Use(middleware.PermissionRequired(model.PermissionCommunityContentModerate))
			{
				communityContentGroup.GET("/community/skills", h.adminCommunitySkills)
				communityContentGroup.POST("/community/skills/:id/hide", h.adminHideSkill)
				communityContentGroup.GET("/community/discussions", h.adminCommunityDiscussions)
				communityContentGroup.POST("/community/discussions/:id/remove", h.adminDeleteDiscussion)
				communityContentGroup.GET("/community/comments", h.adminCommunityComments)
				communityContentGroup.POST("/community/comments/:id/remove", h.adminDeleteComment)
			}

			conversationModerationGroup := adminGroup.Group("/")
			conversationModerationGroup.Use(middleware.PermissionRequired(model.PermissionConversationModerate))
			{
				conversationModerationGroup.GET("/community/conversations", h.adminConversations)
				conversationModerationGroup.POST("/community/conversations/:id/block", h.adminBlockConversation)
				conversationModerationGroup.POST("/community/conversations/:id/unblock", h.adminUnblockConversation)
			}

			workspaceModerationGroup := adminGroup.Group("/")
			workspaceModerationGroup.Use(middleware.PermissionRequired(model.PermissionWorkspaceModerate))
			{
				workspaceModerationGroup.GET("/community/workspaces", h.adminWorkspaces)
				workspaceModerationGroup.POST("/community/workspaces/:id/block", h.adminBlockWorkspace)
				workspaceModerationGroup.POST("/community/workspaces/:id/unblock", h.adminUnblockWorkspace)
				workspaceModerationGroup.POST("/community/workspaces/:id/remove-member", h.adminRemoveWorkspaceMember)
			}

			wikiAdminGroup := adminGroup.Group("/")
			wikiAdminGroup.Use(middleware.PermissionRequired(model.PermissionWikiManage))
			{
				wikiAdminGroup.GET("/wiki/pages", h.adminWikiPages)
				wikiAdminGroup.POST("/wiki/pages/:id/lock", h.adminLockWikiPage)
				wikiAdminGroup.POST("/wiki/pages/:id/unlock", h.adminUnlockWikiPage)
				wikiAdminGroup.POST("/wiki/pages/:id/rollback", h.adminRollbackWikiPage)
			}

			rewardGroup := adminGroup.Group("/")
			rewardGroup.Use(middleware.PermissionRequired(model.PermissionRewardManage))
			{
				rewardGroup.GET("/rewards/overview", h.adminRewardOverview)
				rewardGroup.GET("/rewards/benefits", h.adminRewardBenefits)
				rewardGroup.PUT("/rewards/benefits/:id", h.adminUpdateRewardBenefit)
				rewardGroup.GET("/rewards/adjustments", h.adminRewardAdjustments)
				rewardGroup.POST("/rewards/adjustments", h.adminAdjustRewardPoints)
			}
		}
	}

	h.registerFrontend(router)
}

func (h *Handler) register(c *gin.Context) {
	var input dto.RegisterRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, registerLabels)))
		return
	}
	data, err := h.auth.Register(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) login(c *gin.Context) {
	var input dto.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, loginLabels)))
		return
	}
	data, err := h.auth.Login(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) me(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.auth.Profile(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) logout(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.auth.Logout(claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "logged out", nil)
}

func (h *Handler) profile(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.users.Profile(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) updateProfile(c *gin.Context) {
	var input dto.ProfileUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, profileLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.users.UpdateProfile(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) getHome(c *gin.Context) {
	data, err := h.portal.Home()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) searchAll(c *gin.Context) {
	var query dto.SearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_query", support.ValidationMessage(err, searchLabels)))
		return
	}
	data, err := h.search.Search(query)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listModels(c *gin.Context) {
	var query dto.ResourceListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_query", err.Error()))
		return
	}
	data, err := h.models.List(query)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) modelDetail(c *gin.Context) {
	id := parseUintParam(c, "id")
	data, err := h.models.Detail(id, h.currentUserID(c))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createModel(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var form dto.ModelCreateRequest
	if err := c.ShouldBind(&form); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, modelLabels)))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "missing_file", "请上传模型文件"))
		return
	}

	filePath, fileName, err := h.storage.SaveUploadedFile("models", fileHeader)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	payload := service.ModelCreateInput{
		Name:         form.Name,
		Summary:      form.Summary,
		Description:  form.Description,
		Tags:         strings.Split(form.Tags, ","),
		RobotType:    form.RobotType,
		InputSpec:    form.InputSpec,
		OutputSpec:   form.OutputSpec,
		License:      form.License,
		Dependencies: strings.Split(form.Dependencies, ","),
		Version:      form.Version,
		Changelog:    form.Changelog,
		FilePath:     filePath,
		FileName:     fileName,
	}
	data, err := h.models.Create(claims.UserID, payload)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) updateModel(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.ModelUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, modelUpdateLabels)))
		return
	}

	data, err := h.models.Update(parseUintParam(c, "id"), claims.UserID, service.ModelUpdateInput{
		Name:         input.Name,
		Summary:      input.Summary,
		Description:  input.Description,
		Tags:         strings.Split(input.Tags, ","),
		RobotType:    input.RobotType,
		InputSpec:    input.InputSpec,
		OutputSpec:   input.OutputSpec,
		License:      input.License,
		Dependencies: strings.Split(input.Dependencies, ","),
	})
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) addModelVersion(c *gin.Context) {
	claims := middleware.MustClaims(c)
	id := parseUintParam(c, "id")
	var input dto.ModelVersionRequest
	if err := c.ShouldBind(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, modelVersionLabels)))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "missing_file", "请上传版本文件"))
		return
	}
	filePath, fileName, err := h.storage.SaveUploadedFile("models", fileHeader)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	data, err := h.models.AddVersion(id, claims.UserID, service.ModelCreateInput{
		Version:   input.Version,
		Changelog: input.Changelog,
		FilePath:  filePath,
		FileName:  fileName,
	})
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) submitModel(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.models.Submit(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "model submitted for review", nil)
}

func (h *Handler) downloadModel(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.models.Download(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "model download recorded", nil)
}

func (h *Handler) listDatasets(c *gin.Context) {
	var query dto.ResourceListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_query", err.Error()))
		return
	}
	data, err := h.datasets.List(query)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) datasetOptions(c *gin.Context) {
	data, err := h.catalog.DatasetOptions()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) filterOptions(c *gin.Context) {
	data, err := h.catalog.FilterOptions()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) datasetDetail(c *gin.Context) {
	id := parseUintParam(c, "id")
	data, err := h.datasets.Detail(id, h.currentUserID(c))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) datasetSamples(c *gin.Context) {
	data, err := h.datasets.Samples(parseUintParam(c, "id"), h.currentUserID(c))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createDataset(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.DatasetCreateRequest
	if err := c.ShouldBind(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, datasetLabels)))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "missing_file", "请上传数据集文件"))
		return
	}
	filePath, fileName, err := h.storage.SaveUploadedFile("datasets", fileHeader)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	payload := service.DatasetCreateInput{
		Name:          input.Name,
		Summary:       input.Summary,
		Description:   input.Description,
		Tags:          strings.Split(input.Tags, ","),
		SampleCount:   input.SampleCount,
		Device:        input.Device,
		Scene:         input.Scene,
		Privacy:       input.Privacy,
		AgreementText: input.AgreementText,
		Version:       input.Version,
		Changelog:     input.Changelog,
		FilePath:      filePath,
		FileName:      fileName,
		SamplePreview: splitLines(input.SamplePreview),
		Samples:       collectDatasetSampleFiles(h.storage, c),
	}
	data, err := h.datasets.Create(claims.UserID, payload)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) updateDataset(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.DatasetUpdateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, datasetUpdateLabels)))
		return
	}

	data, err := h.datasets.Update(parseUintParam(c, "id"), claims.UserID, service.DatasetUpdateInput{
		Name:          input.Name,
		Summary:       input.Summary,
		Description:   input.Description,
		Tags:          strings.Split(input.Tags, ","),
		SampleCount:   input.SampleCount,
		Device:        input.Device,
		Scene:         input.Scene,
		Privacy:       input.Privacy,
		AgreementText: input.AgreementText,
		SamplePreview: splitLines(input.SamplePreview),
	})
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) addDatasetVersion(c *gin.Context) {
	claims := middleware.MustClaims(c)
	id := parseUintParam(c, "id")
	var input dto.DatasetVersionRequest
	if err := c.ShouldBind(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, modelVersionLabels)))
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil || fileHeader == nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "missing_file", "请上传版本文件"))
		return
	}
	filePath, fileName, err := h.storage.SaveUploadedFile("datasets", fileHeader)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	data, err := h.datasets.AddVersion(id, claims.UserID, service.DatasetCreateInput{
		Version:   input.Version,
		Changelog: input.Changelog,
		FilePath:  filePath,
		FileName:  fileName,
	})
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) submitDataset(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.datasets.Submit(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "dataset submitted for review", nil)
}

func (h *Handler) confirmDatasetAgreement(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.datasets.ConfirmAgreement(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "agreement accepted", nil)
}

func (h *Handler) downloadDataset(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.datasets.Download(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "dataset download recorded", nil)
}

func (h *Handler) listDatasetDownloadPackages(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.datasets.DownloadPackageTasks(parseUintParam(c, "id"), claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createDatasetDownloadPackage(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.DownloadPackageCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil && err.Error() != "EOF" {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "下载任务参数不正确"))
		return
	}
	data, err := h.datasets.CreateDownloadPackageTask(parseUintParam(c, "id"), claims.UserID, input.Parts)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) myDatasetAccessRequest(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.datasets.MyAccessRequest(parseUintParam(c, "id"), claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myDatasetAccessHistory(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.datasets.MyAccessRequestHistory(parseUintParam(c, "id"), claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createDatasetAccessRequest(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.DatasetAccessRequestPayload
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "访问申请参数不正确"))
		return
	}
	data, err := h.datasets.CreateAccessRequest(parseUintParam(c, "id"), claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) myDatasetAccessRequests(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.datasets.MyAccessRequests(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listTemplates(c *gin.Context) {
	data, err := h.catalog.Templates()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) templateDetail(c *gin.Context) {
	data, err := h.catalog.Template(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listApplicationCases(c *gin.Context) {
	data, err := h.catalog.ApplicationCases()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) applicationCaseDetail(c *gin.Context) {
	data, err := h.catalog.ApplicationCase(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listDocCategories(c *gin.Context) {
	data, err := h.catalog.DocCategories(c.Query("doc_type"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listDocs(c *gin.Context) {
	data, err := h.catalog.Docs(c.Query("doc_type"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) docDetail(c *gin.Context) {
	data, err := h.catalog.Doc(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listFAQs(c *gin.Context) {
	data, err := h.catalog.FAQs()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) publicUserProfile(c *gin.Context) {
	data, err := h.users.PublicProfile(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myUploads(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.users.Uploads(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myFavorites(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.users.Favorites(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) toggleFavorite(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.FavoriteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, favoriteLabels)))
		return
	}
	if err := h.users.ToggleFavorite(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "favorite toggled", nil)
}

func (h *Handler) myDownloads(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.users.Downloads(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myNotifications(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.users.Notifications(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) readNotifications(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.users.MarkNotificationsRead(claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "notifications marked as read", nil)
}

func (h *Handler) readNotification(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.users.MarkNotificationRead(claims.UserID, parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "notification marked as read", nil)
}

func (h *Handler) adminDashboard(c *gin.Context) {
	data, err := h.admin.Dashboard()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminReviews(c *gin.Context) {
	data, err := h.admin.Reviews(c.Query("type"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminReviewDecision(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.ReviewDecisionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, reviewLabels)))
		return
	}
	if err := h.admin.Decide(c.Param("type"), parseUintParam(c, "id"), claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "review decision applied", nil)
}

func (h *Handler) adminAnnouncements(c *gin.Context) {
	data, err := h.admin.Announcements()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateAnnouncement(c *gin.Context) {
	var input dto.AnnouncementRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, announcementLabels)))
		return
	}
	data, err := h.admin.CreateAnnouncement(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminDatasetAccessRequests(c *gin.Context) {
	var query dto.DatasetAccessAdminQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "访问申请筛选参数不正确"))
		return
	}
	data, err := h.datasets.AdminAccessRequests(query)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminReviewDatasetAccessRequest(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.DatasetAccessDecisionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "访问申请审核参数不正确"))
		return
	}
	if err := h.datasets.ReviewAccessRequest(parseUintParam(c, "id"), claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "dataset access request reviewed", nil)
}

func (h *Handler) adminBatchReviewDatasetAccessRequests(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.BatchDatasetAccessDecisionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, batchDatasetAccessDecisionLabels)))
		return
	}
	if err := h.datasets.BatchReviewAccessRequests(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "dataset access requests reviewed", nil)
}

func parseUintParam(c *gin.Context, key string) uint {
	value, _ := strconv.ParseUint(c.Param(key), 10, 64)
	return uint(value)
}

func (h *Handler) currentUserID(c *gin.Context) uint {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return 0
	}
	claims, err := h.tokens.Parse(strings.TrimPrefix(authHeader, "Bearer "))
	if err != nil {
		return 0
	}
	return claims.UserID
}

func splitLines(input string) []string {
	lines := strings.Split(input, "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func collectDatasetSampleFiles(storage support.ObjectStorage, c *gin.Context) []service.DatasetSampleInput {
	form, err := c.MultipartForm()
	if err != nil || form == nil {
		return nil
	}
	files := form.File["sample_files"]
	result := make([]service.DatasetSampleInput, 0, len(files))
	for _, fileHeader := range files {
		filePath, fileName, saveErr := storage.SaveUploadedFile("dataset-samples", fileHeader)
		if saveErr != nil {
			continue
		}
		sampleType := detectDatasetSampleType(fileName)
		previewText := ""
		if sampleType == "video" {
			previewText = "视频样本已上传，可通过链接查看或下载。"
		}
		if sampleType == "pointcloud" {
			previewText = "点云样本已上传，当前提供占位预览与文件链接。"
		}
		if sampleType == "image" {
			previewText = "图片样本预览"
		}
		result = append(result, service.DatasetSampleInput{
			SampleType:  sampleType,
			Title:       fileName,
			PreviewText: previewText,
			FilePath:    filePath,
			FileName:    fileName,
		})
	}
	return result
}

func detectDatasetSampleType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg":
		return "image"
	case ".mp4", ".webm", ".ogg", ".mov":
		return "video"
	case ".ply", ".pcd", ".las", ".laz", ".obj":
		return "pointcloud"
	default:
		return "file"
	}
}

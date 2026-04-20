package app

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/config"
	"github.com/Andras5014/EIBotHub/server/internal/handler"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/service"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type App struct {
	cfg    config.Config
	Router *gin.Engine
	DB     *gorm.DB
}

func New() (*App, error) {
	return newApp(config.Load())
}

func NewForTest(cfg config.Config) (*App, error) {
	return newApp(cfg)
}

func (a *App) Run() error {
	return a.Router.Run(":" + a.cfg.Port)
}

func newApp(cfg config.Config) (*App, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.DBPath), 0o755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(cfg.StorageDir, 0o755); err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Announcement{},
		&model.HomeModuleSetting{},
		&model.HomeHighlight{},
		&model.HomeHeroConfig{},
		&model.RankingConfig{},
		&model.ScenePageConfig{},
		&model.FeaturedResource{},
		&model.ModelAsset{},
		&model.ModelVersion{},
		&model.Dataset{},
		&model.DatasetVersion{},
		&model.DatasetSample{},
		&model.DownloadPackageTask{},
		&model.DatasetAccessRequest{},
		&model.TaskTemplate{},
		&model.ApplicationCase{},
		&model.DocumentCategory{},
		&model.Document{},
		&model.FAQ{},
		&model.VideoTutorial{},
		&model.AgreementTemplate{},
		&model.DatasetPrivacyOption{},
		&model.Favorite{},
		&model.DownloadRecord{},
		&model.FileObject{},
		&model.Notification{},
		&model.AgreementRecord{},
		&model.ReviewLog{},
		&model.AdminOperationLog{},
		&model.ModelEvaluation{},
		&model.ResourceRating{},
		&model.ResourceComment{},
		&model.Skill{},
		&model.Discussion{},
		&model.Follow{},
		&model.SearchRecord{},
		&model.SearchKeywordConfig{},
		&model.FilterOptionConfig{},
		&model.DeveloperVerification{},
		&model.WikiPage{},
		&model.WikiRevision{},
		&model.RewardLedger{},
		&model.RewardBenefit{},
		&model.RewardRedemption{},
		&model.WebhookSubscription{},
		&model.WebhookDelivery{},
		&model.Conversation{},
		&model.ConversationParticipant{},
		&model.Message{},
		&model.Workspace{},
		&model.WorkspaceMember{},
	); err != nil {
		return nil, err
	}

	storage := support.NewLocalStorage(cfg.StorageDir)
	if cfg.SeedDemo {
		if err := repository.EnsureSeeded(db, func(tx *gorm.DB) error {
			return seed(tx, storage, cfg.AppSecret)
		}); err != nil {
			return nil, err
		}
		if err := ensureV1SeedData(db, storage, cfg.AppSecret); err != nil {
			return nil, err
		}
	}

	tokenManager := support.NewTokenManager(cfg.AppSecret)
	fileTokenManager := support.NewFileTokenManager(cfg.AppSecret)

	userRepo := repository.NewUserRepository(db)
	portalRepo := repository.NewPortalRepository(db)
	modelRepo := repository.NewModelRepository(db)
	datasetRepo := repository.NewDatasetRepository(db)
	contentRepo := repository.NewContentRepository(db)
	activityRepo := repository.NewUserActivityRepository(db)
	fileRepo := repository.NewFileObjectRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	operationRepo := repository.NewOperationLogRepository(db)
	communityRepo := repository.NewCommunityRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)
	wikiRepo := repository.NewWikiRepository(db)
	rewardRepo := repository.NewRewardRepository(db)
	integrationRepo := repository.NewIntegrationRepository(db)
	collaborationRepo := repository.NewCollaborationRepository(db)

	authService := service.NewAuthService(userRepo, tokenManager)
	fileService := service.NewFileService(fileRepo, storage, fileTokenManager)
	modelService := service.NewModelService(modelRepo, activityRepo, reviewRepo, fileService)
	datasetService := service.NewDatasetService(datasetRepo, activityRepo, userRepo, verificationRepo, reviewRepo, storage, fileService)
	catalogService := service.NewCatalogService(modelRepo, datasetRepo, contentRepo)
	portalService := service.NewPortalService(portalRepo, modelRepo, datasetRepo, contentRepo)
	userService := service.NewUserService(userRepo, activityRepo, modelService, datasetService, reviewRepo)
	searchService := service.NewSearchService(modelRepo, datasetRepo, contentRepo, userRepo, communityRepo)
	adminService := service.NewAdminService(userRepo, portalRepo, modelRepo, datasetRepo, contentRepo, reviewRepo, operationRepo)
	integrationService := service.NewIntegrationService(integrationRepo)
	rewardService := service.NewRewardService(rewardRepo, userRepo, activityRepo, integrationService)
	communityService := service.NewCommunityService(communityRepo, userRepo, modelRepo, datasetRepo, contentRepo, activityRepo, rewardService, integrationService)
	verificationService := service.NewVerificationService(verificationRepo, userRepo, activityRepo, integrationService)
	wikiService := service.NewWikiService(wikiRepo, userRepo, integrationService)
	collaborationService := service.NewCollaborationService(collaborationRepo, userRepo, activityRepo, operationRepo)

	router := gin.New()
	handler.New(authService, portalService, modelService, datasetService, catalogService, userService, searchService, adminService, communityService, verificationService, wikiService, rewardService, integrationService, collaborationService, fileService, storage, tokenManager).Register(router)

	return &App{
		cfg:    cfg,
		Router: router,
		DB:     db,
	}, nil
}

func seed(db *gorm.DB, storage *support.LocalStorage, secret string) error {
	tokenManager := support.NewTokenManager(secret)
	now := time.Now().Add(-6 * time.Hour)

	admin := model.User{
		Username:     "admin",
		Email:        "admin@opencommunity.local",
		PasswordHash: tokenManager.HashPassword("Admin123!"),
		Role:         model.RoleAdmin,
		Bio:          "Open community operator",
	}
	user := model.User{
		Username:     "demo",
		Email:        "demo@example.com",
		PasswordHash: tokenManager.HashPassword("Demo123!"),
		Role:         model.RoleUser,
		Bio:          "Embodied AI developer",
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}
	if err := db.Create(&user).Error; err != nil {
		return err
	}

	modelFile, err := storage.WriteSeedFile("models", "seed-warehouse.pt", "seed model content")
	if err != nil {
		return err
	}
	datasetFile, err := storage.WriteSeedFile("datasets", "seed-inspection.zip", "seed dataset content")
	if err != nil {
		return err
	}
	imageSamplePath, err := storage.WriteGeneratedFile("dataset-samples", "seed-inspection-preview.svg", `<svg xmlns="http://www.w3.org/2000/svg" width="640" height="360"><rect width="100%" height="100%" fill="#eaf3ff"/><text x="40" y="180" font-size="28" fill="#16324f">Inspection Image Preview</text></svg>`)
	if err != nil {
		return err
	}
	videoSamplePath, err := storage.WriteGeneratedFile("dataset-samples", "seed-inspection-video.mp4", "video placeholder")
	if err != nil {
		return err
	}
	pointCloudSamplePath, err := storage.WriteGeneratedFile("dataset-samples", "seed-inspection-pointcloud.ply", "ply\nformat ascii 1.0\ncomment placeholder point cloud")
	if err != nil {
		return err
	}

	announcements := []model.Announcement{
		{Title: "开放社区 MVP 上线", Summary: "欢迎体验模型、数据集、模板与文档中心。", Link: "/docs", Pinned: true, PublishedAt: now},
		{Title: "首批示例资源已发布", Summary: "包含仓储搬运模型、巡检数据集和移动作业模板。", Link: "/search", Pinned: false, PublishedAt: now.Add(30 * time.Minute)},
	}
	if err := db.Create(&announcements).Error; err != nil {
		return err
	}

	moduleSettings := []model.HomeModuleSetting{
		{ModuleKey: "hero", SortOrder: 10, Enabled: true},
		{ModuleKey: "models", SortOrder: 20, Enabled: true},
		{ModuleKey: "announcements", SortOrder: 30, Enabled: true},
		{ModuleKey: "scenes", SortOrder: 40, Enabled: true},
		{ModuleKey: "resources", SortOrder: 50, Enabled: true},
		{ModuleKey: "community", SortOrder: 60, Enabled: true},
	}
	if err := db.Create(&moduleSettings).Error; err != nil {
		return err
	}

	scenePages := []model.ScenePageConfig{
		{
			Slug:        "patrol-inspection",
			Name:        "巡逻巡检",
			Tagline:     "Patrol Inspection",
			Summary:     "围绕园区巡检、楼宇安防和设备告警联动的资源专题页。",
			Description: "聚合巡检场景常用模型、数据集、模板和案例，优先满足园区安防、楼宇巡检和弱光告警场景的落地需求。",
			SortOrder:   10,
			Enabled:     true,
		},
		{
			Slug:        "warehouse-ops",
			Name:        "仓储搬运",
			Tagline:     "Warehouse Operation",
			Summary:     "围绕仓储搬运、分拣交接和物料调度的场景化入口。",
			Description: "聚合仓储场景常用抓取、避障、交接模板和案例，帮助开发者更快定位搬运与分拣场景的组合资源。",
			SortOrder:   20,
			Enabled:     true,
		},
	}
	if err := db.Create(&scenePages).Error; err != nil {
		return err
	}

	highlights := []model.HomeHighlight{
		{Text: "统一管理模型、数据集和任务模板", SortOrder: 10, Enabled: true},
		{Text: "面向移动作业、巡逻巡检、搬运等场景", SortOrder: 20, Enabled: true},
		{Text: "提供文档、下载统计、审核和运营后台", SortOrder: 30, Enabled: true},
	}
	if err := db.Create(&highlights).Error; err != nil {
		return err
	}

	heroConfig := model.HomeHeroConfig{
		Tagline:         "EIBotHub具生训练",
		Title:           "围绕模型、数据集、模板与文档的开放社区",
		Description:     "开放社区聚合模型、数据集、任务模板、文档与具身应用案例，帮助开发者围绕机器人场景快速完成接入、训练、部署和复用。",
		PrimaryButton:   "上传模型",
		SecondaryButton: "上传数据集",
		SearchButton:    "进入全局搜索",
	}
	if err := db.Create(&heroConfig).Error; err != nil {
		return err
	}

	rankingConfig := model.RankingConfig{
		Title:    "贡献排行榜",
		Subtitle: "基于积分展示近期社区贡献活跃度。",
		Limit:    5,
		Enabled:  true,
	}
	if err := db.Create(&rankingConfig).Error; err != nil {
		return err
	}

	models := []model.ModelAsset{
		{
			Name:         "Warehouse Carrier Net",
			Summary:      "适用于移动搬运机器人的基础感知模型。",
			Description:  "用于仓储场景下的移动搬运机器人感知与避障。",
			Tags:         "warehouse,carrier,vision",
			RobotType:    "搬运",
			InputSpec:    "RGB frame",
			OutputSpec:   "obstacle map",
			License:      "Apache-2.0",
			Dependencies: "onnxruntime,opencv",
			Status:       model.StatusPublished,
			Downloads:    18,
			OwnerID:      user.ID,
			Versions: []model.ModelVersion{
				{Version: "v1.0.0", FilePath: modelFile, FileName: "seed-warehouse.pt", Changelog: "initial release"},
			},
		},
	}
	if err := db.Create(&models).Error; err != nil {
		return err
	}

	datasets := []model.Dataset{
		{
			Name:          "Inspection Route Set",
			Summary:       "巡逻巡检场景的基础标注数据集。",
			Description:   "包含室内外巡检线路、设备缺陷和文字告警样本。",
			Tags:          "inspection,route,vision",
			SampleCount:   1280,
			Device:        "RGB Camera",
			Scene:         "巡逻巡检",
			Privacy:       "public",
			AgreementText: "下载并使用本数据集前需同意开放社区使用协议。",
			Status:        model.StatusPublished,
			Downloads:     24,
			OwnerID:       user.ID,
			Versions: []model.DatasetVersion{
				{Version: "v1.0.0", FilePath: datasetFile, FileName: "seed-inspection.zip", Changelog: "initial release"},
			},
			Samples: []model.DatasetSample{
				{SampleType: "text", Title: "Sample 1", PreviewText: "缺陷点位：A-12，告警等级：中。"},
				{SampleType: "text", Title: "Sample 2", PreviewText: "巡检线路：楼层 2 -> 楼层 3 -> 机房入口。"},
				{SampleType: "image", Title: "巡检图片预览", PreviewText: "图片样本预览", FilePath: imageSamplePath, FileName: "seed-inspection-preview.svg"},
				{SampleType: "video", Title: "巡检视频占位样本", PreviewText: "视频样本已上传，可通过链接查看或下载。", FilePath: videoSamplePath, FileName: "seed-inspection-video.mp4"},
				{SampleType: "pointcloud", Title: "点云占位样本", PreviewText: "点云样本已上传，当前提供占位预览与文件链接。", FilePath: pointCloudSamplePath, FileName: "seed-inspection-pointcloud.ply"},
			},
		},
	}
	if err := db.Create(&datasets).Error; err != nil {
		return err
	}

	templates := []model.TaskTemplate{
		{
			Name:        "移动作业巡检模板",
			Summary:     "适用于常规室内巡逻巡检任务的标准模板。",
			Description: "包含起点、巡逻点、告警上报和回充流程。",
			Category:    "巡逻巡检",
			Scene:       "移动作业",
			Guide:       "1. 绑定设备 2. 选择路线 3. 配置告警 4. 发布执行",
			ResourceRef: "Warehouse Carrier Net,Inspection Route Set",
			UsageCount:  12,
			Status:      model.StatusPublished,
		},
	}
	if err := db.Create(&templates).Error; err != nil {
		return err
	}

	appCases := []model.ApplicationCase{
		{
			Title:    "仓储搬运联合调度案例",
			Summary:  "展示搬运机器人在仓储调度中的联合工作方式。",
			Category: "搬运",
			Guide:    "部署步骤：接入设备、导入模型、绑定任务模板、发布调度。",
			Status:   model.StatusPublished,
		},
	}
	if err := db.Create(&appCases).Error; err != nil {
		return err
	}

	featuredResources := []model.FeaturedResource{
		{ResourceType: "model", ResourceID: models[0].ID, BadgeLabel: "编辑推荐", SortOrder: 1, Enabled: true},
		{ResourceType: "dataset", ResourceID: datasets[0].ID, BadgeLabel: "精选数据", SortOrder: 1, Enabled: true},
		{ResourceType: "task-template", ResourceID: templates[0].ID, BadgeLabel: "标准模板", SortOrder: 1, Enabled: true},
		{ResourceType: "application-case", ResourceID: appCases[0].ID, BadgeLabel: "场景案例", SortOrder: 1, Enabled: true},
	}
	if err := db.Create(&featuredResources).Error; err != nil {
		return err
	}

	platformCategory := model.DocumentCategory{Name: "平台指南", DocType: "platform"}
	technicalCategory := model.DocumentCategory{Name: "接入文档", DocType: "technical"}
	if err := db.Create(&platformCategory).Error; err != nil {
		return err
	}
	if err := db.Create(&technicalCategory).Error; err != nil {
		return err
	}

	documents := []model.Document{
		{
			CategoryID: platformCategory.ID,
			Title:      "开放社区使用指南",
			Summary:    "快速了解门户、资源上传、审核和下载流程。",
			Content:    "使用流程：注册登录 -> 上传资源 -> 提交审核 -> 审核通过后展示与下载。",
			DocType:    "platform",
			Status:     model.StatusPublished,
		},
		{
			CategoryID: technicalCategory.ID,
			Title:      "机器人接入说明",
			Summary:    "说明机器人类型、任务模板和资源引用关系。",
			Content:    "接入步骤：准备设备参数 -> 选择模型与数据集 -> 导入模板 -> 调试运行。",
			DocType:    "technical",
			Status:     model.StatusPublished,
		},
	}
	if err := db.Create(&documents).Error; err != nil {
		return err
	}

	faqs := []model.FAQ{
		{Question: "如何上传模型？", Answer: "登录后进入模型上传页，填写元数据并上传模型文件即可。"},
		{Question: "数据集下载前为何需要确认协议？", Answer: "用于记录数据集下载协议接受情况，保障资源使用合规。"},
	}
	if err := db.Create(&faqs).Error; err != nil {
		return err
	}

	videos := []model.VideoTutorial{
		{
			Title:     "开放社区五分钟上手",
			Summary:   "快速了解首页、资源上传、审核和下载流程。",
			Link:      "https://example.com/videos/getting-started",
			Category:  "平台指南",
			SortOrder: 1,
			Active:    true,
		},
		{
			Title:     "机器人接入与模板绑定演示",
			Summary:   "演示如何把模型、数据集和模板串起来完成联调。",
			Link:      "https://example.com/videos/integration-demo",
			Category:  "接入实践",
			SortOrder: 2,
			Active:    true,
		},
	}
	if err := db.Create(&videos).Error; err != nil {
		return err
	}

	agreementTemplates := []model.AgreementTemplate{
		{
			Name:      "开放社区标准协议",
			Content:   "下载并使用本数据集前需同意开放社区使用协议，并在使用时保留数据来源说明。",
			SortOrder: 1,
			Active:    true,
		},
		{
			Name:      "科研试用协议",
			Content:   "仅限科研试用，不得用于商业化发布，成果需注明数据集来源。",
			SortOrder: 2,
			Active:    true,
		},
	}
	if err := db.Create(&agreementTemplates).Error; err != nil {
		return err
	}

	privacyOptions := []model.DatasetPrivacyOption{
		{Code: "public", Name: "公开", Description: "任何登录用户在确认协议后可下载。", SortOrder: 1, Active: true},
		{Code: "internal", Name: "内部", Description: "仅限受邀或指定成员访问。", SortOrder: 2, Active: true},
		{Code: "restricted", Name: "受限", Description: "需额外审批后才能下载。", SortOrder: 3, Active: true},
	}
	if err := db.Create(&privacyOptions).Error; err != nil {
		return err
	}

	notifications := []model.Notification{
		{UserID: user.ID, Type: "system", Title: "欢迎加入开放社区", Content: "你可以先浏览门户首页，再尝试上传模型或数据集。"},
		{UserID: user.ID, Type: "resource", Title: "示例模型已发布", Content: "你可以在模型仓库查看 Warehouse Carrier Net。"},
	}
	if err := db.Create(&notifications).Error; err != nil {
		return err
	}

	evaluations := []model.ModelEvaluation{
		{
			ModelID:   models[0].ID,
			UserID:    user.ID,
			Benchmark: "Warehouse Obstacle Benchmark",
			Summary:   "在仓储避障基准上表现稳定。",
			Score:     86.5,
			Notes:     "适合中等密度障碍场景。",
		},
	}
	if err := db.Create(&evaluations).Error; err != nil {
		return err
	}

	ratings := []model.ResourceRating{
		{ResourceType: "model", ResourceID: models[0].ID, UserID: user.ID, Score: 5, Feedback: "部署简单，适合 MVP。"},
		{ResourceType: "task-template", ResourceID: templates[0].ID, UserID: user.ID, Score: 4, Feedback: "模板结构清晰。"},
	}
	if err := db.Create(&ratings).Error; err != nil {
		return err
	}

	skills := []model.Skill{
		{
			Name:        "巡检告警上报技能",
			Summary:     "把巡检任务中的异常告警统一封装成可复用技能。",
			Description: "适用于移动作业和巡逻巡检流程中的告警上传与状态同步。",
			Category:    "巡逻巡检",
			Scene:       "移动作业",
			Guide:       "接入告警源 -> 配置 webhook -> 绑定任务模板 -> 验证状态同步。",
			ResourceRef: "Warehouse Carrier Net,Inspection Route Set",
			Status:      model.StatusPublished,
			OwnerID:     user.ID,
			UsageCount:  8,
		},
	}
	if err := db.Create(&skills).Error; err != nil {
		return err
	}

	discussions := []model.Discussion{
		{
			Title:    "移动作业模板如何接入现有机器人？",
			Summary:  "讨论模板资源和设备参数之间的最小接入集。",
			Content:  "想确认在已有机器人项目里，模板与模型的依赖绑定最少需要准备哪些参数。",
			Category: "接入讨论",
			UserID:   user.ID,
		},
	}
	if err := db.Create(&discussions).Error; err != nil {
		return err
	}

	comments := []model.ResourceComment{
		{ResourceType: "task-template", ResourceID: templates[0].ID, UserID: user.ID, Content: "这个模板适合做首个巡检 MVP。"},
		{ResourceType: "skill", ResourceID: skills[0].ID, UserID: user.ID, Content: "技能拆分后便于团队复用。"},
		{ResourceType: "discussion", ResourceID: discussions[0].ID, UserID: admin.ID, Content: "建议优先确认设备参数和回传协议。"},
	}
	if err := db.Create(&comments).Error; err != nil {
		return err
	}

	follows := []model.Follow{
		{FollowerID: admin.ID, FollowedUserID: user.ID},
	}
	if err := db.Create(&follows).Error; err != nil {
		return err
	}

	searchRecords := []model.SearchRecord{
		{Query: "巡检", SearchType: "dataset"},
		{Query: "巡检", SearchType: "task-template"},
		{Query: "搬运", SearchType: "model"},
		{Query: "部署", SearchType: "doc"},
	}
	if err := db.Create(&searchRecords).Error; err != nil {
		return err
	}

	searchKeywordConfigs := []model.SearchKeywordConfig{
		{Query: "巡检", KeywordType: "hot", SortOrder: 1, Enabled: true},
		{Query: "搬运", KeywordType: "hot", SortOrder: 2, Enabled: true},
		{Query: "夜间安防", KeywordType: "recommended", SortOrder: 1, Enabled: true},
		{Query: "仓储搬运", KeywordType: "recommended", SortOrder: 2, Enabled: true},
	}
	return db.Create(&searchKeywordConfigs).Error
}

func ensureV1SeedData(db *gorm.DB, storage *support.LocalStorage, secret string) error {
	var admin model.User
	if err := db.Where("email = ?", "admin@opencommunity.local").First(&admin).Error; err != nil {
		return nil
	}
	var demo model.User
	if err := db.Where("email = ?", "demo@example.com").First(&demo).Error; err != nil {
		return nil
	}
	tokenManager := support.NewTokenManager(secret)

	var baseModel model.ModelAsset
	_ = db.Order("id asc").First(&baseModel).Error
	var baseDataset model.Dataset
	_ = db.Order("id asc").First(&baseDataset).Error
	var baseTemplate model.TaskTemplate
	_ = db.Order("id asc").First(&baseTemplate).Error
	var baseAppCase model.ApplicationCase
	_ = db.Order("id asc").First(&baseAppCase).Error

	integrator, err := ensureSeedUser(db, tokenManager, "integrator", "integrator@example.com", "现场接入工程师")
	if err != nil {
		return err
	}
	opsUser, err := ensureSeedUser(db, tokenManager, "ops", "ops@example.com", "平台运营与质检")
	if err != nil {
		return err
	}
	reviewerUser, err := ensureSeedUser(db, tokenManager, "seed-reviewer", "seed-reviewer@example.com", "内容审核与质量把关")
	if err != nil {
		return err
	}
	maintainerUser, err := ensureSeedUser(db, tokenManager, "seed-maintainer", "seed-maintainer@example.com", "社区维护与模板整理")
	if err != nil {
		return err
	}
	analystUser, err := ensureSeedUser(db, tokenManager, "seed-analyst", "seed-analyst@example.com", "场景分析与数据标注")
	if err != nil {
		return err
	}
	creatorUser, err := ensureSeedUser(db, tokenManager, "seed-creator", "seed-creator@example.com", "内容策划与案例撰写")
	if err != nil {
		return err
	}
	qaUser, err := ensureSeedUser(db, tokenManager, "seed-qa", "seed-qa@example.com", "测试验证与回归检查")
	if err != nil {
		return err
	}

	var evaluationCount int64
	if err := db.Model(&model.ModelEvaluation{}).Count(&evaluationCount).Error; err != nil {
		return err
	}
	if evaluationCount == 0 && baseModel.ID > 0 {
		if err := db.Create(&model.ModelEvaluation{
			ModelID:   baseModel.ID,
			UserID:    demo.ID,
			Benchmark: "Warehouse Obstacle Benchmark",
			Summary:   "在仓储避障基准上表现稳定。",
			Score:     86.5,
			Notes:     "适合中等密度障碍场景。",
		}).Error; err != nil {
			return err
		}
	}

	var ratingCount int64
	if err := db.Model(&model.ResourceRating{}).Count(&ratingCount).Error; err != nil {
		return err
	}
	if ratingCount == 0 {
		ratings := []model.ResourceRating{}
		if baseModel.ID > 0 {
			ratings = append(ratings, model.ResourceRating{ResourceType: "model", ResourceID: baseModel.ID, UserID: demo.ID, Score: 5, Feedback: "部署简单，适合 MVP。"})
		}
		if baseTemplate.ID > 0 {
			ratings = append(ratings, model.ResourceRating{ResourceType: "task-template", ResourceID: baseTemplate.ID, UserID: demo.ID, Score: 4, Feedback: "模板结构清晰。"})
		}
		if len(ratings) > 0 {
			if err := db.Create(&ratings).Error; err != nil {
				return err
			}
		}
	}

	var skillCount int64
	if err := db.Model(&model.Skill{}).Count(&skillCount).Error; err != nil {
		return err
	}
	if skillCount == 0 {
		skill := model.Skill{
			Name:        "巡检告警上报技能",
			Summary:     "把巡检任务中的异常告警统一封装成可复用技能。",
			Description: "适用于移动作业和巡逻巡检流程中的告警上传与状态同步。",
			Category:    "巡逻巡检",
			Scene:       "移动作业",
			Guide:       "接入告警源 -> 配置 webhook -> 绑定任务模板 -> 验证状态同步。",
			ResourceRef: "Warehouse Carrier Net,Inspection Route Set",
			Status:      model.StatusPublished,
			OwnerID:     demo.ID,
			UsageCount:  8,
		}
		if err := db.Create(&skill).Error; err != nil {
			return err
		}
	}

	var discussionCount int64
	if err := db.Model(&model.Discussion{}).Count(&discussionCount).Error; err != nil {
		return err
	}
	if discussionCount == 0 {
		discussion := model.Discussion{
			Title:    "移动作业模板如何接入现有机器人？",
			Summary:  "讨论模板资源和设备参数之间的最小接入集。",
			Content:  "想确认在已有机器人项目里，模板与模型的依赖绑定最少需要准备哪些参数。",
			Category: "接入讨论",
			UserID:   demo.ID,
		}
		if err := db.Create(&discussion).Error; err != nil {
			return err
		}
	}

	var commentCount int64
	if err := db.Model(&model.ResourceComment{}).Count(&commentCount).Error; err != nil {
		return err
	}
	if commentCount == 0 {
		var skill model.Skill
		_ = db.Order("id asc").First(&skill).Error
		var discussion model.Discussion
		_ = db.Order("id asc").First(&discussion).Error
		comments := []model.ResourceComment{}
		if baseTemplate.ID > 0 {
			comments = append(comments, model.ResourceComment{ResourceType: "task-template", ResourceID: baseTemplate.ID, UserID: demo.ID, Content: "这个模板适合做首个巡检 MVP。"})
		}
		if skill.ID > 0 {
			comments = append(comments, model.ResourceComment{ResourceType: "skill", ResourceID: skill.ID, UserID: demo.ID, Content: "技能拆分后便于团队复用。"})
		}
		if discussion.ID > 0 {
			comments = append(comments, model.ResourceComment{ResourceType: "discussion", ResourceID: discussion.ID, UserID: admin.ID, Content: "建议优先确认设备参数和回传协议。"})
		}
		if len(comments) > 0 {
			if err := db.Create(&comments).Error; err != nil {
				return err
			}
		}
	}

	var followCount int64
	if err := db.Model(&model.Follow{}).Count(&followCount).Error; err != nil {
		return err
	}
	if followCount == 0 {
		if err := db.Create(&model.Follow{FollowerID: admin.ID, FollowedUserID: demo.ID}).Error; err != nil {
			return err
		}
	}

	var searchCount int64
	if err := db.Model(&model.SearchRecord{}).Count(&searchCount).Error; err != nil {
		return err
	}
	if searchCount == 0 {
		searchRecords := []model.SearchRecord{
			{Query: "巡检", SearchType: "dataset"},
			{Query: "巡检", SearchType: "task-template"},
			{Query: "搬运", SearchType: "model"},
			{Query: "部署", SearchType: "doc"},
		}
		if err := db.Create(&searchRecords).Error; err != nil {
			return err
		}
	}

	var searchKeywordConfigCount int64
	if err := db.Model(&model.SearchKeywordConfig{}).Count(&searchKeywordConfigCount).Error; err != nil {
		return err
	}
	if searchKeywordConfigCount == 0 {
		searchKeywordConfigs := []model.SearchKeywordConfig{
			{Query: "巡检", KeywordType: "hot", SortOrder: 1, Enabled: true},
			{Query: "搬运", KeywordType: "hot", SortOrder: 2, Enabled: true},
			{Query: "夜间安防", KeywordType: "recommended", SortOrder: 1, Enabled: true},
			{Query: "仓储搬运", KeywordType: "recommended", SortOrder: 2, Enabled: true},
		}
		if err := db.Create(&searchKeywordConfigs).Error; err != nil {
			return err
		}
	}

	var moduleSettingCount int64
	if err := db.Model(&model.HomeModuleSetting{}).Count(&moduleSettingCount).Error; err != nil {
		return err
	}
	if moduleSettingCount == 0 {
		moduleSettings := []model.HomeModuleSetting{
			{ModuleKey: "hero", SortOrder: 10, Enabled: true},
			{ModuleKey: "models", SortOrder: 20, Enabled: true},
			{ModuleKey: "announcements", SortOrder: 30, Enabled: true},
			{ModuleKey: "scenes", SortOrder: 40, Enabled: true},
			{ModuleKey: "resources", SortOrder: 50, Enabled: true},
			{ModuleKey: "community", SortOrder: 60, Enabled: true},
		}
		if err := db.Create(&moduleSettings).Error; err != nil {
			return err
		}
	}
	var highlightCount int64
	if err := db.Model(&model.HomeHighlight{}).Count(&highlightCount).Error; err != nil {
		return err
	}
	if highlightCount == 0 {
		highlights := []model.HomeHighlight{
			{Text: "统一管理模型、数据集和任务模板", SortOrder: 10, Enabled: true},
			{Text: "面向移动作业、巡逻巡检、搬运等场景", SortOrder: 20, Enabled: true},
			{Text: "提供文档、下载统计、审核和运营后台", SortOrder: 30, Enabled: true},
		}
		if err := db.Create(&highlights).Error; err != nil {
			return err
		}
	}
	var heroConfigCount int64
	if err := db.Model(&model.HomeHeroConfig{}).Count(&heroConfigCount).Error; err != nil {
		return err
	}
	if heroConfigCount == 0 {
		heroConfig := model.HomeHeroConfig{
			Tagline:         "OpenLoong 风格信息门户",
			Title:           "围绕模型、数据集、模板与文档的开放社区",
			Description:     "开放社区聚合模型、数据集、任务模板、文档与具身应用案例，帮助开发者围绕机器人场景快速完成接入、训练、部署和复用。",
			PrimaryButton:   "上传模型",
			SecondaryButton: "上传数据集",
			SearchButton:    "进入全局搜索",
		}
		if err := db.Create(&heroConfig).Error; err != nil {
			return err
		}
	}
	var rankingConfigCount int64
	if err := db.Model(&model.RankingConfig{}).Count(&rankingConfigCount).Error; err != nil {
		return err
	}
	if rankingConfigCount == 0 {
		rankingConfig := model.RankingConfig{
			Title:    "贡献排行榜",
			Subtitle: "基于积分展示近期社区贡献活跃度。",
			Limit:    5,
			Enabled:  true,
		}
		if err := db.Create(&rankingConfig).Error; err != nil {
			return err
		}
	}
	var scenePageCount int64
	if err := db.Model(&model.ScenePageConfig{}).Count(&scenePageCount).Error; err != nil {
		return err
	}
	if scenePageCount == 0 {
		scenePages := []model.ScenePageConfig{
			{
				Slug:        "patrol-inspection",
				Name:        "巡逻巡检",
				Tagline:     "Patrol Inspection",
				Summary:     "围绕园区巡检、楼宇安防和设备告警联动的资源专题页。",
				Description: "聚合巡检场景常用模型、数据集、模板和案例，优先满足园区安防、楼宇巡检和弱光告警场景的落地需求。",
				SortOrder:   10,
				Enabled:     true,
			},
			{
				Slug:        "warehouse-ops",
				Name:        "仓储搬运",
				Tagline:     "Warehouse Operation",
				Summary:     "围绕仓储搬运、分拣交接和物料调度的场景化入口。",
				Description: "聚合仓储场景常用抓取、避障、交接模板和案例，帮助开发者更快定位搬运与分拣场景的组合资源。",
				SortOrder:   20,
				Enabled:     true,
			},
		}
		if err := db.Create(&scenePages).Error; err != nil {
			return err
		}
	}

	curatedModels := make([]model.ModelAsset, 0, 4)
	for _, spec := range []seedModelSpec{
		{
			Name:         "Patrol Corridor Vision",
			Summary:      "面向园区与楼宇巡检走廊的轻量视觉识别模型。",
			Description:  "适配走廊和楼层巡检机器人，聚焦通道障碍、门禁状态与异常物体识别。",
			Tags:         "patrol,corridor,vision",
			RobotType:    "巡逻巡检",
			InputSpec:    "RGB frame",
			OutputSpec:   "event labels",
			License:      "Apache-2.0",
			Dependencies: "onnxruntime,opencv",
			Downloads:    16,
			FileName:     "seed-patrol-corridor.onnx",
			FileContent:  "seed patrol corridor model",
			Changelog:    "seed demo release",
		},
		{
			Name:         "Sorting Cell Grasp Policy",
			Summary:      "适用于小型分拣单元的抓取与放置策略模型。",
			Description:  "用于轻量搬运与分拣场景中的抓取动作决策、抓取失败恢复与节拍优化。",
			Tags:         "sorting,grasp,manipulation",
			RobotType:    "搬运",
			InputSpec:    "RGBD frame",
			OutputSpec:   "grasp action",
			License:      "MIT",
			Dependencies: "pytorch,opencv",
			Downloads:    14,
			FileName:     "seed-sorting-grasp.pt",
			FileContent:  "seed sorting grasp model",
			Changelog:    "seed demo release",
		},
		{
			Name:         "Night Security Audio Fusion",
			Summary:      "面向夜间安防与园区巡逻场景的音视频融合识别模型。",
			Description:  "用于夜间安防、低照巡逻和异响告警场景中的多传感联动识别。",
			Tags:         "night-security,audio,patrol",
			RobotType:    "巡逻巡检",
			InputSpec:    "RGB + audio stream",
			OutputSpec:   "alert events",
			License:      "Apache-2.0",
			Dependencies: "ffmpeg,onnxruntime",
			Downloads:    13,
			FileName:     "seed-night-security.onnx",
			FileContent:  "seed night security model",
			Changelog:    "seed demo release",
		},
		{
			Name:         "Factory Quality Defect Locator",
			Summary:      "适用于产线质检与缺陷复检场景的轻量定位模型。",
			Description:  "聚焦产线质检中的表面瑕疵、装配偏差和复检框选定位。",
			Tags:         "factory-quality,defect,vision",
			RobotType:    "移动作业",
			InputSpec:    "RGB frame",
			OutputSpec:   "defect boxes",
			License:      "MIT",
			Dependencies: "opencv,tensorrt",
			Downloads:    12,
			FileName:     "seed-factory-quality.engine",
			FileContent:  "seed factory quality model",
			Changelog:    "seed demo release",
		},
	} {
		item, err := ensureSeedModel(db, storage, demo.ID, spec)
		if err != nil {
			return err
		}
		curatedModels = append(curatedModels, *item)
	}

	curatedDatasets := make([]model.Dataset, 0, 4)
	for _, spec := range []seedDatasetSpec{
		{
			Name:          "Factory Safety Alert Set",
			Summary:       "产线安全巡检中的告警、围栏和人员靠近样本集。",
			Description:   "覆盖告警灯、围栏入侵、人员靠近与设备停机等典型安全巡检事件。",
			Tags:          "factory,safety,inspection",
			SampleCount:   860,
			Device:        "RGB Camera + Event Log",
			Scene:         "巡逻巡检",
			Privacy:       "public",
			AgreementText: "下载并使用本数据集前需同意开放社区使用协议。",
			Downloads:     19,
			FileName:      "seed-factory-safety.zip",
			FileContent:   "seed factory safety dataset",
			Changelog:     "seed demo release",
		},
		{
			Name:          "Forklift Aisle Mapping Pack",
			Summary:       "面向仓储搬运机器人的巷道定位与障碍标注数据集。",
			Description:   "提供叉车通道、货架间距、遮挡物与动态障碍的结构化标注样本。",
			Tags:          "forklift,mapping,warehouse",
			SampleCount:   1430,
			Device:        "LiDAR + RGBD",
			Scene:         "仓储搬运",
			Privacy:       "public",
			AgreementText: "下载并使用本数据集前需同意开放社区使用协议。",
			Downloads:     17,
			FileName:      "seed-forklift-mapping.zip",
			FileContent:   "seed forklift mapping dataset",
			Changelog:     "seed demo release",
		},
		{
			Name:          "Night Security Audio Event Pack",
			Summary:       "夜间安防场景中的异响、脚步和门禁异常事件数据集。",
			Description:   "覆盖夜间安防巡逻中的门禁异常、玻璃破碎、脚步靠近和环境异响等样本。",
			Tags:          "night-security,audio,alert",
			SampleCount:   920,
			Device:        "Microphone + RGB Camera",
			Scene:         "夜间安防",
			Privacy:       "public",
			AgreementText: "下载并使用本数据集前需同意开放社区使用协议。",
			Downloads:     15,
			FileName:      "seed-night-security-audio.zip",
			FileContent:   "seed night security dataset",
			Changelog:     "seed demo release",
		},
		{
			Name:          "Factory Quality Defect Set",
			Summary:       "面向产线质检与复检流程的缺陷检测数据集。",
			Description:   "覆盖划痕、装配偏差、螺丝缺失和表面污染等产线质检样本。",
			Tags:          "factory-quality,defect,inspection",
			SampleCount:   1180,
			Device:        "Industrial Camera",
			Scene:         "产线质检",
			Privacy:       "public",
			AgreementText: "下载并使用本数据集前需同意开放社区使用协议。",
			Downloads:     14,
			FileName:      "seed-factory-quality.zip",
			FileContent:   "seed factory quality dataset",
			Changelog:     "seed demo release",
		},
	} {
		item, err := ensureSeedDataset(db, storage, demo.ID, spec)
		if err != nil {
			return err
		}
		curatedDatasets = append(curatedDatasets, *item)
	}

	curatedTemplates := make([]model.TaskTemplate, 0, 4)
	for _, spec := range []seedTemplateSpec{
		{
			Name:        "夜间安防巡检模板",
			Summary:     "适用于园区夜间安防巡检任务的标准执行模板。",
			Description: "覆盖巡逻点、异常上报、视频回传与回充收尾流程。",
			Category:    "安防巡检",
			Scene:       "夜间巡检",
			Guide:       "1. 配置路线 2. 绑定告警规则 3. 设置录像回传 4. 发布执行",
			ResourceRef: "Patrol Corridor Vision,Factory Safety Alert Set",
			UsageCount:  18,
		},
		{
			Name:        "仓储搬运交接模板",
			Summary:     "适用于入库、出库与工位交接的搬运任务模板。",
			Description: "包含起始工位、交接确认、异常中断与回传闭环。",
			Category:    "搬运作业",
			Scene:       "仓储搬运",
			Guide:       "1. 选择搬运单元 2. 绑定抓取策略 3. 配置交接规则 4. 发布任务",
			ResourceRef: "Sorting Cell Grasp Policy,Forklift Aisle Mapping Pack",
			UsageCount:  15,
		},
		{
			Name:        "多模态告警联动模板",
			Summary:     "适用于多模态告警聚合和通知分发的任务模板。",
			Description: "用于夜间安防与多传感巡逻场景中的音视频告警聚合、分发和闭环记录。",
			Category:    "多模态告警",
			Scene:       "多模态告警",
			Guide:       "1. 绑定音视频输入 2. 配置事件等级 3. 设置通知链路 4. 发布执行",
			ResourceRef: "Night Security Audio Fusion,Night Security Audio Event Pack",
			UsageCount:  13,
		},
		{
			Name:        "产线质检复检模板",
			Summary:     "适用于产线质检与异常复检流程的标准模板。",
			Description: "覆盖缺陷初检、人工复核、复检回传与关闭工单闭环。",
			Category:    "产线质检",
			Scene:       "产线质检",
			Guide:       "1. 绑定质检模型 2. 配置复检工位 3. 设置回传字段 4. 发布执行",
			ResourceRef: "Factory Quality Defect Locator,Factory Quality Defect Set",
			UsageCount:  11,
		},
	} {
		item, err := ensureSeedTemplate(db, spec)
		if err != nil {
			return err
		}
		curatedTemplates = append(curatedTemplates, *item)
	}

	curatedCases := make([]model.ApplicationCase, 0, 4)
	for _, spec := range []seedApplicationCaseSpec{
		{
			Title:    "园区夜间安防巡逻案例",
			Summary:  "展示夜间巡逻、异常识别与告警联动的闭环流程。",
			Category: "安防",
			Guide:    "部署步骤：配置巡逻路线、绑定巡检模型、接入告警通道、联调回传。",
		},
		{
			Title:    "产线物料搬运闭环案例",
			Summary:  "展示料箱搬运、工位交接与异常回退的完整流程。",
			Category: "搬运",
			Guide:    "部署步骤：导入搬运模型、加载巷道数据集、绑定交接模板、发布任务。",
		},
		{
			Title:    "多模态告警联动案例",
			Summary:  "展示音视频事件聚合、告警派发和异常升级的闭环流程。",
			Category: "多模态告警",
			Guide:    "部署步骤：绑定音视频模型、导入事件数据集、配置告警模板、联调通知链路。",
		},
		{
			Title:    "产线质检复检案例",
			Summary:  "展示缺陷检测、人工复核和复检回传的完整质量闭环。",
			Category: "产线质检",
			Guide:    "部署步骤：绑定质检模型、导入缺陷数据、配置复检模板、发布质检流程。",
		},
	} {
		item, err := ensureSeedApplicationCase(db, spec)
		if err != nil {
			return err
		}
		curatedCases = append(curatedCases, *item)
	}

	for _, spec := range []seedAnnouncementSpec{
		{
			Title:   "社区榜单与 Wiki 已开放",
			Summary: "现在可以在首页查看贡献排行，并直接进入最新 Wiki 词条。",
			Link:    "/wiki",
			Pinned:  false,
		},
		{
			Title:   "新增多场景演示资源",
			Summary: "补充了夜间安防巡检、产线搬运等首页演示内容。",
			Link:    "/applications",
			Pinned:  false,
		},
		{
			Title:   "文档中心补充 5 组示例内容",
			Summary: "平台指南、接入文档、FAQ 和视频教程均已补充可演示的测试数据。",
			Link:    "/docs",
			Pinned:  false,
		},
		{
			Title:   "社区资料页已补充完整演示数据",
			Summary: "个人中心、公开主页、关注关系、收藏与下载记录都已具备 5 条以上样例。",
			Link:    "/me",
			Pinned:  false,
		},
		{
			Title:   "私信与协作空间增加演示样板",
			Summary: "已补齐多条私信会话与协作空间，便于测试治理、消息和成员流转。",
			Link:    "/messages",
			Pinned:  false,
		},
	} {
		if err := ensureSeedAnnouncement(db, spec); err != nil {
			return err
		}
	}

	for _, spec := range []seedWikiSpec{
		{
			Title:    "巡检任务接入 Wiki",
			Summary:  "记录移动作业和巡检模板接入的关键步骤。",
			Content:  "接入建议：先确认设备参数，再绑定模型、数据集和任务模板，最后联调告警链路。",
			EditorID: demo.ID,
			Comment:  "补充巡检任务接入 Wiki",
		},
		{
			Title:    "搬运任务模板参数说明",
			Summary:  "整理仓储搬运任务中常用的模板字段与参数含义。",
			Content:  "建议优先确认起始点、目标点、交接确认方式和异常回退策略，再绑定模型与数据集。",
			EditorID: integrator.ID,
			Comment:  "补充搬运模板说明",
		},
		{
			Title:    "园区巡检告警联调手册",
			Summary:  "整理园区巡检任务的告警通道、回传字段与联调顺序。",
			Content:  "联调建议：先通设备状态，再通告警回传，最后验证任务模板与通知链路的闭环。",
			EditorID: opsUser.ID,
			Comment:  "补充告警联调手册",
		},
		{
			Title:    "夜间安防值守排班说明",
			Summary:  "记录夜间安防巡逻的排班规则、告警等级和联动方式。",
			Content:  "建议将夜间安防任务按时段拆分，分别绑定音视频事件阈值、告警升级规则和回传人群组。",
			EditorID: reviewerUser.ID,
			Comment:  "补充夜间安防说明",
		},
		{
			Title:    "多模态告警事件字段约定",
			Summary:  "统一音频、视频和状态告警的字段命名与事件层级。",
			Content:  "建议统一 event_type、event_level、source_id、snapshot_url、review_status 等字段，保证多系统联动一致。",
			EditorID: analystUser.ID,
			Comment:  "补充告警字段约定",
		},
	} {
		if err := ensureSeedWikiPage(db, spec); err != nil {
			return err
		}
	}

	if baseModel.ID == 0 && len(curatedModels) > 0 {
		baseModel = curatedModels[0]
	}
	if baseDataset.ID == 0 && len(curatedDatasets) > 0 {
		baseDataset = curatedDatasets[0]
	}
	if baseTemplate.ID == 0 && len(curatedTemplates) > 0 {
		baseTemplate = curatedTemplates[0]
	}
	if baseAppCase.ID == 0 && len(curatedCases) > 0 {
		baseAppCase = curatedCases[0]
	}

	for _, item := range []struct {
		resourceType string
		resourceID   uint
		badgeLabel   string
		sortOrder    int
	}{
		{resourceType: "model", resourceID: baseModel.ID, badgeLabel: "编辑推荐", sortOrder: 1},
		{resourceType: "model", resourceID: curatedModels[0].ID, badgeLabel: "场景精选", sortOrder: 2},
		{resourceType: "model", resourceID: curatedModels[1].ID, badgeLabel: "抓取策略", sortOrder: 3},
		{resourceType: "model", resourceID: curatedModels[2].ID, badgeLabel: "夜巡融合", sortOrder: 4},
		{resourceType: "model", resourceID: curatedModels[3].ID, badgeLabel: "质检定位", sortOrder: 5},
		{resourceType: "dataset", resourceID: baseDataset.ID, badgeLabel: "精选数据", sortOrder: 1},
		{resourceType: "dataset", resourceID: curatedDatasets[0].ID, badgeLabel: "安全巡检", sortOrder: 2},
		{resourceType: "dataset", resourceID: curatedDatasets[1].ID, badgeLabel: "仓储地图", sortOrder: 3},
		{resourceType: "dataset", resourceID: curatedDatasets[2].ID, badgeLabel: "夜间音频", sortOrder: 4},
		{resourceType: "dataset", resourceID: curatedDatasets[3].ID, badgeLabel: "质检样本", sortOrder: 5},
		{resourceType: "task-template", resourceID: baseTemplate.ID, badgeLabel: "标准模板", sortOrder: 1},
		{resourceType: "task-template", resourceID: curatedTemplates[0].ID, badgeLabel: "夜巡方案", sortOrder: 2},
		{resourceType: "task-template", resourceID: curatedTemplates[1].ID, badgeLabel: "交接闭环", sortOrder: 3},
		{resourceType: "task-template", resourceID: curatedTemplates[2].ID, badgeLabel: "告警聚合", sortOrder: 4},
		{resourceType: "task-template", resourceID: curatedTemplates[3].ID, badgeLabel: "复检流程", sortOrder: 5},
		{resourceType: "application-case", resourceID: baseAppCase.ID, badgeLabel: "场景案例", sortOrder: 1},
		{resourceType: "application-case", resourceID: curatedCases[0].ID, badgeLabel: "园区方案", sortOrder: 2},
		{resourceType: "application-case", resourceID: curatedCases[1].ID, badgeLabel: "产线实践", sortOrder: 3},
		{resourceType: "application-case", resourceID: curatedCases[2].ID, badgeLabel: "告警联动", sortOrder: 4},
		{resourceType: "application-case", resourceID: curatedCases[3].ID, badgeLabel: "质检闭环", sortOrder: 5},
	} {
		if err := ensureSeedFeaturedResource(db, item.resourceType, item.resourceID, item.sortOrder, item.badgeLabel); err != nil {
			return err
		}
	}

	var featuredCount int64
	if err := db.Model(&model.FeaturedResource{}).Count(&featuredCount).Error; err != nil {
		return err
	}
	if featuredCount == 0 {
		var appCase model.ApplicationCase
		_ = db.Order("id asc").First(&appCase).Error
		featuredResources := []model.FeaturedResource{}
		if baseModel.ID > 0 {
			featuredResources = append(featuredResources, model.FeaturedResource{ResourceType: "model", ResourceID: baseModel.ID, BadgeLabel: "编辑推荐", SortOrder: 1, Enabled: true})
		}
		if baseDataset.ID > 0 {
			featuredResources = append(featuredResources, model.FeaturedResource{ResourceType: "dataset", ResourceID: baseDataset.ID, BadgeLabel: "精选数据", SortOrder: 1, Enabled: true})
		}
		if baseTemplate.ID > 0 {
			featuredResources = append(featuredResources, model.FeaturedResource{ResourceType: "task-template", ResourceID: baseTemplate.ID, BadgeLabel: "标准模板", SortOrder: 1, Enabled: true})
		}
		if appCase.ID > 0 {
			featuredResources = append(featuredResources, model.FeaturedResource{ResourceType: "application-case", ResourceID: appCase.ID, BadgeLabel: "场景案例", SortOrder: 1, Enabled: true})
		}
		if len(featuredResources) > 0 {
			if err := db.Create(&featuredResources).Error; err != nil {
				return err
			}
		}
	}

	for _, spec := range []seedScenePageSpec{
		{
			Slug:        "patrol-inspection",
			Name:        "巡逻巡检",
			Tagline:     "Patrol Inspection",
			Summary:     "聚合巡逻巡检场景下的模型、数据集、模板与案例。",
			Description: "面向楼宇、园区和工厂巡逻巡检任务，提供从识别、告警到执行模板的一站式资源入口。",
			SortOrder:   10,
		},
		{
			Slug:        "warehouse-ops",
			Name:        "仓储搬运",
			Tagline:     "Warehouse Operation",
			Summary:     "聚合仓储搬运、交接和巷道通行的资源组合。",
			Description: "围绕仓储搬运任务的定位、抓取、交接和异常回退流程，集中整理高频可复用资源。",
			SortOrder:   20,
		},
		{
			Slug:        "campus-guard",
			Name:        "园区值守",
			Tagline:     "Campus Guard",
			Summary:     "聚合园区值守、夜间巡逻和异常上报的专题资源。",
			Description: "帮助开发者快速定位园区值守场景下的巡逻模型、数据集、模板与案例配置。",
			SortOrder:   30,
		},
		{
			Slug:        "factory-quality",
			Name:        "产线质检",
			Tagline:     "Factory Quality",
			Summary:     "聚合产线质检、复检和异常回传相关资源。",
			Description: "围绕产线质检场景提供缺陷定位模型、质检数据集、复检模板和业务案例。",
			SortOrder:   40,
		},
		{
			Slug:        "multimodal-alert",
			Name:        "多模态告警",
			Tagline:     "Multimodal Alert",
			Summary:     "聚合多模态告警联动、事件聚合和通知分发资源。",
			Description: "适用于音视频、状态事件等多源信息融合告警场景，支持聚合与分发闭环测试。",
			SortOrder:   50,
		},
	} {
		if _, err := ensureSeedScenePage(db, spec); err != nil {
			return err
		}
	}

	var platformCategory model.DocumentCategory
	if err := db.Where("name = ?", "平台指南").First(&platformCategory).Error; err != nil {
		return err
	}
	var technicalCategory model.DocumentCategory
	if err := db.Where("name = ?", "接入文档").First(&technicalCategory).Error; err != nil {
		return err
	}

	for _, spec := range []seedDocumentSpec{
		{
			CategoryID: platformCategory.ID,
			Title:      "开放社区使用指南",
			Summary:    "快速了解门户、资源上传、审核和下载流程。",
			Content:    "使用流程：注册登录 -> 上传资源 -> 提交审核 -> 审核通过后展示与下载。",
			DocType:    "platform",
		},
		{
			CategoryID: technicalCategory.ID,
			Title:      "机器人接入说明",
			Summary:    "说明机器人类型、任务模板和资源引用关系。",
			Content:    "接入步骤：准备设备参数 -> 选择模型与数据集 -> 导入模板 -> 调试运行。",
			DocType:    "technical",
		},
		{
			CategoryID: platformCategory.ID,
			Title:      "数据集授权下载说明",
			Summary:    "说明协议确认、审批、授权有效期和下载次数限制。",
			Content:    "建议在下载前先核对权限级别、审批状态、有效期和剩余下载次数，再生成下载或分批任务。",
			DocType:    "platform",
		},
		{
			CategoryID: technicalCategory.ID,
			Title:      "搜索与筛选接入约定",
			Summary:    "说明搜索关键词、筛选字典和资源标签的接入规则。",
			Content:    "建议统一标签、场景、分类和推荐词配置，保持模型、数据集和模板在搜索结果中的一致性。",
			DocType:    "technical",
		},
		{
			CategoryID: platformCategory.ID,
			Title:      "协作空间治理说明",
			Summary:    "说明私信治理、协作空间封禁、成员移除和审计日志。",
			Content:    "后台可对私信和协作空间执行封禁、解封、成员移除等治理动作，并记录管理操作日志。",
			DocType:    "platform",
		},
	} {
		if _, err := ensureSeedDocument(db, spec); err != nil {
			return err
		}
	}

	for _, spec := range []seedFAQSpec{
		{Question: "如何上传模型？", Answer: "登录后进入模型上传页，填写元数据并上传模型文件即可。"},
		{Question: "数据集下载前为何需要确认协议？", Answer: "用于记录数据集下载协议接受情况，保障资源使用合规。"},
		{Question: "为什么有些数据集需要审批？", Answer: "受限或内部数据集会按权限矩阵执行审批、授权有效期和下载次数控制。"},
		{Question: "协作空间能做什么？", Answer: "可用于围绕模型、数据集和模板进行成员协作、消息同步和治理演示。"},
		{Question: "搜索推荐词在哪里维护？", Answer: "管理员可在门户运营页维护热门词、推荐词和筛选字典。"},
	} {
		if _, err := ensureSeedFAQ(db, spec); err != nil {
			return err
		}
	}

	for _, spec := range []seedVideoSpec{
		{
			Title:     "开放社区五分钟上手",
			Summary:   "快速了解首页、资源上传、审核和下载流程。",
			Link:      "https://example.com/videos/getting-started",
			Category:  "平台指南",
			SortOrder: 1,
		},
		{
			Title:     "机器人接入与模板绑定演示",
			Summary:   "演示如何把模型、数据集和模板串起来完成联调。",
			Link:      "https://example.com/videos/integration-demo",
			Category:  "接入实践",
			SortOrder: 2,
		},
		{
			Title:     "夜间安防告警联动演示",
			Summary:   "演示夜间安防场景中的音视频告警联动和通知闭环。",
			Link:      "https://example.com/videos/night-security-demo",
			Category:  "场景演示",
			SortOrder: 3,
		},
		{
			Title:     "产线质检复检流程演示",
			Summary:   "演示缺陷检测、人工复核和复检回传流程。",
			Link:      "https://example.com/videos/factory-quality-demo",
			Category:  "场景演示",
			SortOrder: 4,
		},
		{
			Title:     "社区治理后台操作演示",
			Summary:   "演示私信、协作空间、评论和讨论的治理链路。",
			Link:      "https://example.com/videos/community-governance",
			Category:  "后台运营",
			SortOrder: 5,
		},
	} {
		if _, err := ensureSeedVideo(db, spec); err != nil {
			return err
		}
	}

	for _, spec := range []seedSkillSpec{
		{
			Name:        "巡检告警上报技能",
			Summary:     "把巡检任务中的异常告警统一封装成可复用技能。",
			Description: "适用于移动作业和巡逻巡检流程中的告警上传与状态同步。",
			Category:    "巡逻巡检",
			Scene:       "巡逻巡检",
			Guide:       "接入告警源 -> 配置 webhook -> 绑定任务模板 -> 验证状态同步。",
			ResourceRef: "Warehouse Carrier Net,Inspection Route Set",
			OwnerID:     demo.ID,
			UsageCount:  8,
		},
		{
			Name:        "夜间视频巡逻编排技能",
			Summary:     "用于夜间安防任务中的视频采集、事件汇总和告警上报。",
			Description: "适用于夜间安防、低照巡逻和园区值守流程中的视频巡逻编排。",
			Category:    "夜间安防",
			Scene:       "夜间安防",
			Guide:       "配置巡逻点位 -> 绑定音视频模型 -> 设置事件阈值 -> 联调通知链路。",
			ResourceRef: "Night Security Audio Fusion,Night Security Audio Event Pack",
			OwnerID:     demo.ID,
			UsageCount:  7,
		},
		{
			Name:        "巷道搬运交接校验技能",
			Summary:     "用于仓储搬运交接前后的状态确认与失败回退。",
			Description: "适用于仓储搬运、入出库交接和工位交互场景中的状态校验。",
			Category:    "仓储搬运",
			Scene:       "仓储搬运",
			Guide:       "绑定工位 -> 设置交接字段 -> 配置失败回退规则 -> 发布执行。",
			ResourceRef: "Sorting Cell Grasp Policy,Forklift Aisle Mapping Pack",
			OwnerID:     demo.ID,
			UsageCount:  6,
		},
		{
			Name:        "产线缺陷复检协同技能",
			Summary:     "把缺陷初检、人工复核和复检回传整理成可复用步骤。",
			Description: "适用于产线质检、复检和异常复盘场景中的协同流程封装。",
			Category:    "产线质检",
			Scene:       "产线质检",
			Guide:       "绑定缺陷模型 -> 配置复检字段 -> 接入回传接口 -> 联调工单关闭。",
			ResourceRef: "Factory Quality Defect Locator,Factory Quality Defect Set",
			OwnerID:     demo.ID,
			UsageCount:  5,
		},
		{
			Name:        "多传感告警聚合技能",
			Summary:     "适用于多模态告警场景下的事件聚合与升级分发。",
			Description: "用于多传感巡逻、状态监控和音视频事件融合后的告警治理流程。",
			Category:    "多模态告警",
			Scene:       "多模态告警",
			Guide:       "接入多源事件 -> 配置聚合规则 -> 设置升级策略 -> 发布告警闭环。",
			ResourceRef: "Night Security Audio Fusion,多模态告警联动模板",
			OwnerID:     demo.ID,
			UsageCount:  9,
		},
	} {
		if _, err := ensureSeedSkill(db, spec); err != nil {
			return err
		}
	}

	for _, spec := range []seedDiscussionSpec{
		{
			Title:    "移动作业模板如何接入现有机器人？",
			Summary:  "讨论模板资源和设备参数之间的最小接入集。",
			Content:  "想确认在已有机器人项目里，模板与模型的依赖绑定最少需要准备哪些参数。",
			Category: "接入讨论",
			UserID:   demo.ID,
		},
		{
			Title:    "夜间安防巡检模型怎么做告警阈值校准？",
			Summary:  "讨论夜间安防场景里音视频联合阈值的校准策略。",
			Content:  "目前夜间场景里误报偏高，想确认音频和视频事件在联合告警时应该如何设置阈值分层。",
			Category: "夜间安防",
			UserID:   demo.ID,
		},
		{
			Title:    "仓储搬运交接环节如何设计失败回退？",
			Summary:  "讨论交接确认、状态回滚和失败重试的设计边界。",
			Content:  "在搬运交接过程中，遇到工位响应超时或识别失败时，大家一般怎么设计回退和重试策略？",
			Category: "仓储搬运",
			UserID:   demo.ID,
		},
		{
			Title:    "点云与视觉联合巡检是否值得默认开启？",
			Summary:  "讨论多传感巡检在精度和成本之间的取舍。",
			Content:  "当前在巡检车上可同时启用点云和视觉，想了解在演示环境里默认开启是否值得，还是先保留轻量组合更合适。",
			Category: "多模态告警",
			UserID:   demo.ID,
		},
		{
			Title:    "产线质检复检模板需要哪些最小输入？",
			Summary:  "讨论质检复检流程里最关键的字段设计。",
			Content:  "准备把产线质检做成模板，想确认缺陷类别、工位编号、复检状态和关闭原因是否应作为必填字段。",
			Category: "产线质检",
			UserID:   demo.ID,
		},
	} {
		if _, err := ensureSeedDiscussion(db, spec); err != nil {
			return err
		}
	}

	var seededSkills []model.Skill
	if err := db.Where("owner_id = ?", demo.ID).Order("updated_at desc").Find(&seededSkills).Error; err != nil {
		return err
	}
	var seededDiscussions []model.Discussion
	if err := db.Where("user_id = ?", demo.ID).Order("updated_at desc").Find(&seededDiscussions).Error; err != nil {
		return err
	}

	commentTargets := []struct {
		resourceType string
		resourceID   uint
		userID       uint
		content      string
	}{
		{resourceType: "model", resourceID: baseModel.ID, userID: admin.ID, content: "这个模型的部署链路比较顺，适合首页展示。"},
		{resourceType: "model", resourceID: curatedModels[2].ID, userID: reviewerUser.ID, content: "夜间安防场景下建议结合音频阈值一起调。"},
		{resourceType: "dataset", resourceID: baseDataset.ID, userID: integrator.ID, content: "巡检线路样本组织得比较规范，适合做快速验证。"},
		{resourceType: "dataset", resourceID: curatedDatasets[3].ID, userID: analystUser.ID, content: "产线质检样本的缺陷划分维度比较清晰。"},
		{resourceType: "task-template", resourceID: curatedTemplates[0].ID, userID: opsUser.ID, content: "夜间巡检模板的告警闭环配置很完整。"},
		{resourceType: "skill", resourceID: seededSkills[0].ID, userID: maintainerUser.ID, content: "这个技能拆分得很细，适合作为团队复用基线。"},
		{resourceType: "discussion", resourceID: seededDiscussions[0].ID, userID: admin.ID, content: "建议先把设备字段和模板参数一一映射，再做最小接入。"},
		{resourceType: "discussion", resourceID: seededDiscussions[3].ID, userID: qaUser.ID, content: "多传感默认开启前最好先压一轮性能和误报率。"},
	}
	for _, item := range commentTargets {
		if item.resourceID == 0 {
			continue
		}
		if err := ensureSeedComment(db, item.resourceType, item.resourceID, item.userID, item.content); err != nil {
			return err
		}
	}

	for _, item := range []struct {
		query       string
		keywordType string
		sortOrder   int
	}{
		{query: "巡检", keywordType: "hot", sortOrder: 1},
		{query: "搬运", keywordType: "hot", sortOrder: 2},
		{query: "夜间安防", keywordType: "hot", sortOrder: 3},
		{query: "产线质检", keywordType: "hot", sortOrder: 4},
		{query: "多模态告警", keywordType: "hot", sortOrder: 5},
		{query: "夜间安防", keywordType: "recommended", sortOrder: 1},
		{query: "仓储搬运", keywordType: "recommended", sortOrder: 2},
		{query: "巡逻巡检", keywordType: "recommended", sortOrder: 3},
		{query: "产线质检", keywordType: "recommended", sortOrder: 4},
		{query: "多模态告警", keywordType: "recommended", sortOrder: 5},
	} {
		if err := ensureSeedSearchKeyword(db, item.query, item.keywordType, item.sortOrder); err != nil {
			return err
		}
	}

	for _, item := range []struct {
		query      string
		searchType string
	}{
		{query: "巡检", searchType: "dataset"},
		{query: "巡检", searchType: "task-template"},
		{query: "搬运", searchType: "model"},
		{query: "部署", searchType: "doc"},
		{query: "夜间安防", searchType: "model"},
		{query: "夜间安防", searchType: "dataset"},
		{query: "产线质检", searchType: "dataset"},
		{query: "产线质检", searchType: "task-template"},
		{query: "多模态告警", searchType: "skill"},
		{query: "多模态告警", searchType: "discussion"},
	} {
		if err := ensureSeedSearchRecord(db, item.query, item.searchType); err != nil {
			return err
		}
	}

	displayResources := []struct {
		resourceType string
		resourceID   uint
		title        string
	}{
		{resourceType: "model", resourceID: baseModel.ID, title: baseModel.Name},
		{resourceType: "model", resourceID: curatedModels[0].ID, title: curatedModels[0].Name},
		{resourceType: "dataset", resourceID: baseDataset.ID, title: baseDataset.Name},
		{resourceType: "dataset", resourceID: curatedDatasets[0].ID, title: curatedDatasets[0].Name},
		{resourceType: "task-template", resourceID: baseTemplate.ID, title: baseTemplate.Name},
	}
	for _, item := range displayResources {
		if item.resourceID == 0 {
			continue
		}
		if err := ensureSeedFavorite(db, demo.ID, item.resourceType, item.resourceID, item.title); err != nil {
			return err
		}
		if err := ensureSeedDownload(db, demo.ID, item.resourceType, item.resourceID, item.title); err != nil {
			return err
		}
	}

	for _, item := range []struct {
		notifyType string
		title      string
		content    string
	}{
		{notifyType: "system", title: "欢迎加入开放社区", content: "你可以先浏览门户首页，再尝试上传模型或数据集。"},
		{notifyType: "resource", title: "示例模型已发布", content: "你可以在模型仓库查看 Warehouse Carrier Net。"},
		{notifyType: "review", title: "夜间安防模板已进入展示区", content: "夜间安防巡检模板已进入首页专题资源，可继续完善说明与案例。"},
		{notifyType: "community", title: "你的讨论收到新回复", content: "“仓储搬运交接环节如何设计失败回退？”已有新的社区回复。"},
		{notifyType: "reward", title: "积分与权益样例已补充", content: "你的个人中心已补齐积分账本、兑换记录和权益示例数据。"},
	} {
		if err := ensureSeedNotification(db, demo.ID, item.notifyType, item.title, item.content); err != nil {
			return err
		}
	}

	for _, followedUserID := range []uint{integrator.ID, opsUser.ID, reviewerUser.ID, maintainerUser.ID, creatorUser.ID} {
		if err := ensureSeedFollow(db, demo.ID, followedUserID); err != nil {
			return err
		}
	}
	for _, followerID := range []uint{admin.ID, integrator.ID, opsUser.ID, reviewerUser.ID, analystUser.ID} {
		if err := ensureSeedFollow(db, followerID, demo.ID); err != nil {
			return err
		}
	}

	for _, spec := range []seedRewardSpec{
		{UserID: demo.ID, SourceType: "seed_skill", Points: 30, Remark: "示例技能发布"},
		{UserID: demo.ID, SourceType: "seed_discussion", Points: 10, Remark: "示例讨论创建"},
		{UserID: demo.ID, SourceType: "seed_model_showcase", Points: 60, Remark: "模型专题样例补充"},
		{UserID: demo.ID, SourceType: "seed_dataset_showcase", Points: 50, Remark: "数据集专题样例补充"},
		{UserID: demo.ID, SourceType: "seed_template_showcase", Points: 45, Remark: "模板专题样例补充"},
		{UserID: demo.ID, SourceType: "seed_wiki_showcase", Points: 45, Remark: "Wiki 与资料页样例补充"},
		{UserID: admin.ID, SourceType: "seed_admin", Points: 5, Remark: "示例运营积分"},
		{UserID: reviewerUser.ID, SourceType: "seed_review", Points: 22, Remark: "审核演示积分"},
		{UserID: maintainerUser.ID, SourceType: "seed_maintain", Points: 19, Remark: "社区维护积分"},
		{UserID: analystUser.ID, SourceType: "seed_analysis", Points: 17, Remark: "数据分析积分"},
		{UserID: creatorUser.ID, SourceType: "seed_creator", Points: 21, Remark: "内容策划积分"},
		{UserID: qaUser.ID, SourceType: "seed_qa", Points: 15, Remark: "测试验证积分"},
	} {
		if err := ensureSeedReward(db, spec); err != nil {
			return err
		}
	}

	benefitSpecs := []struct {
		name       string
		summary    string
		costPoints int
	}{
		{name: "首页创作者标识", summary: "兑换后可在公开开发者页显示创作者标识。", costPoints: 20},
		{name: "社区推荐位申请", summary: "兑换后可获得一次社区首页推荐位申请资格。", costPoints: 40},
		{name: "协作空间扩容", summary: "兑换后可为一个协作空间增加更多成员名额。", costPoints: 60},
		{name: "专题页案例露出", summary: "兑换后可获得一次场景专题页案例展示机会。", costPoints: 45},
		{name: "数据集授权加急通道", summary: "兑换后可为一次数据集授权申请开启加急处理。", costPoints: 35},
	}
	seededBenefits := make([]model.RewardBenefit, 0, len(benefitSpecs))
	for _, spec := range benefitSpecs {
		item, err := ensureSeedRewardBenefit(db, spec.name, spec.summary, spec.costPoints)
		if err != nil {
			return err
		}
		seededBenefits = append(seededBenefits, *item)
	}
	for _, item := range seededBenefits {
		if err := ensureSeedRewardRedemption(db, demo.ID, item); err != nil {
			return err
		}
	}

	accessExpiryShort := time.Now().Add(7 * 24 * time.Hour)
	accessExpiryLong := time.Now().Add(30 * 24 * time.Hour)
	reviewedAt := time.Now().Add(-6 * time.Hour)
	for _, item := range []seedAccessRequestSpec{
		{
			DatasetID:         baseDataset.ID,
			UserID:            demo.ID,
			Reason:            "需要用于巡逻巡检流程联调与样例展示。",
			Status:            "approved",
			ReviewComment:     "已授权用于演示环境验证。",
			ApprovalStage:     1,
			RequiredApprovals: 1,
			ApprovalExpiresAt: &accessExpiryLong,
			DownloadLimit:     10,
			DownloadCount:     2,
			ReviewedBy:        &admin.ID,
			ReviewedAt:        &reviewedAt,
		},
		{
			DatasetID:         curatedDatasets[0].ID,
			UserID:            demo.ID,
			Reason:            "用于夜间安防与安全巡检样例页展示。",
			Status:            "approved",
			ReviewComment:     "已通过，可继续生成下载包。",
			ApprovalStage:     1,
			RequiredApprovals: 1,
			ApprovalExpiresAt: &accessExpiryShort,
			DownloadLimit:     6,
			DownloadCount:     1,
			ReviewedBy:        &admin.ID,
			ReviewedAt:        &reviewedAt,
		},
		{
			DatasetID:         curatedDatasets[1].ID,
			UserID:            demo.ID,
			Reason:            "用于仓储搬运巷道地图样例测试。",
			Status:            "pending",
			ReviewComment:     "",
			ApprovalStage:     0,
			RequiredApprovals: 1,
		},
		{
			DatasetID:         curatedDatasets[2].ID,
			UserID:            demo.ID,
			Reason:            "用于多模态告警音频事件测试。",
			Status:            "rejected",
			ReviewComment:     "请补充使用场景和保留时长说明。",
			ApprovalStage:     0,
			RequiredApprovals: 1,
			ReviewedBy:        &admin.ID,
			ReviewedAt:        &reviewedAt,
		},
		{
			DatasetID:         curatedDatasets[3].ID,
			UserID:            demo.ID,
			Reason:            "用于产线质检复检流程联调。",
			Status:            "pending",
			ReviewComment:     "",
			ApprovalStage:     1,
			RequiredApprovals: 2,
		},
	} {
		if item.DatasetID == 0 {
			continue
		}
		if err := ensureSeedAccessRequest(db, item); err != nil {
			return err
		}
	}

	for _, item := range []struct {
		peer    uint
		title   string
		content string
		sender  uint
	}{
		{peer: admin.ID, title: "管理员与开发者私信", content: "欢迎加入开放社区，如需协作可以直接在这里沟通。", sender: admin.ID},
		{peer: integrator.ID, title: "开发者与接入工程师私信", content: "巡检模板和模型已经同步，准备开始联调。", sender: integrator.ID},
		{peer: opsUser.ID, title: "开发者与运营同学私信", content: "首页展示位的文案和资源标签已经更新。", sender: opsUser.ID},
		{peer: reviewerUser.ID, title: "开发者与审核同学私信", content: "请再补一版使用说明，我这边可以继续审核。", sender: reviewerUser.ID},
		{peer: maintainerUser.ID, title: "开发者与社区维护私信", content: "技能页和讨论区的演示数据已经补充完毕。", sender: maintainerUser.ID},
	} {
		if _, err := ensureDirectConversation(db, demo.ID, item.peer, item.title, item.content, item.sender); err != nil {
			return err
		}
	}

	for _, item := range []struct {
		name     string
		summary  string
		ownerID  uint
		members  []uint
		senderID uint
		content  string
	}{
		{name: "巡检联调空间", summary: "用于巡逻巡检模型、数据集和模板的联合验证。", ownerID: demo.ID, members: []uint{admin.ID, integrator.ID}, senderID: demo.ID, content: "巡检联调资源已就绪，今天先验证异常上报链路。"},
		{name: "夜间安防值守空间", summary: "用于夜间安防音视频告警与通知闭环验证。", ownerID: demo.ID, members: []uint{reviewerUser.ID, opsUser.ID}, senderID: reviewerUser.ID, content: "夜间值守模板需要补一版弱光阈值配置。"},
		{name: "仓储搬运交接空间", summary: "用于仓储搬运交接、失败回退和回传联调。", ownerID: integrator.ID, members: []uint{demo.ID, maintainerUser.ID}, senderID: integrator.ID, content: "交接失败回退策略已经按新模板更新。"},
		{name: "产线质检复检空间", summary: "用于产线质检缺陷复检和工单关闭演示。", ownerID: analystUser.ID, members: []uint{demo.ID, qaUser.ID}, senderID: analystUser.ID, content: "质检复检模板已绑定到新的缺陷数据集。"},
		{name: "多模态告警治理空间", summary: "用于多模态告警聚合与治理策略验证。", ownerID: opsUser.ID, members: []uint{demo.ID, admin.ID, creatorUser.ID}, senderID: opsUser.ID, content: "告警升级规则和社区治理流程可以一起做联调。"},
	} {
		if _, err := ensureWorkspaceSeed(db, item.name, item.summary, item.ownerID, item.members, item.senderID, item.content); err != nil {
			return err
		}
	}

	var videoCount int64
	if err := db.Model(&model.VideoTutorial{}).Count(&videoCount).Error; err != nil {
		return err
	}
	if videoCount == 0 {
		videos := []model.VideoTutorial{
			{
				Title:     "开放社区五分钟上手",
				Summary:   "快速了解首页、资源上传、审核和下载流程。",
				Link:      "https://example.com/videos/getting-started",
				Category:  "平台指南",
				SortOrder: 1,
				Active:    true,
			},
			{
				Title:     "机器人接入与模板绑定演示",
				Summary:   "演示如何把模型、数据集和模板串起来完成联调。",
				Link:      "https://example.com/videos/integration-demo",
				Category:  "接入实践",
				SortOrder: 2,
				Active:    true,
			},
		}
		if err := db.Create(&videos).Error; err != nil {
			return err
		}
	}

	var agreementTemplateCount int64
	if err := db.Model(&model.AgreementTemplate{}).Count(&agreementTemplateCount).Error; err != nil {
		return err
	}
	if agreementTemplateCount == 0 {
		templates := []model.AgreementTemplate{
			{
				Name:      "开放社区标准协议",
				Content:   "下载并使用本数据集前需同意开放社区使用协议，并在使用时保留数据来源说明。",
				SortOrder: 1,
				Active:    true,
			},
			{
				Name:      "科研试用协议",
				Content:   "仅限科研试用，不得用于商业化发布，成果需注明数据集来源。",
				SortOrder: 2,
				Active:    true,
			},
		}
		if err := db.Create(&templates).Error; err != nil {
			return err
		}
	}

	var privacyOptionCount int64
	if err := db.Model(&model.DatasetPrivacyOption{}).Count(&privacyOptionCount).Error; err != nil {
		return err
	}
	if privacyOptionCount == 0 {
		options := []model.DatasetPrivacyOption{
			{Code: "public", Name: "公开", Description: "任何登录用户在确认协议后可下载。", SortOrder: 1, Active: true},
			{Code: "internal", Name: "内部", Description: "仅限受邀或指定成员访问。", SortOrder: 2, Active: true},
			{Code: "restricted", Name: "受限", Description: "需额外审批后才能下载。", SortOrder: 3, Active: true},
		}
		if err := db.Create(&options).Error; err != nil {
			return err
		}
	}

	if baseDataset.ID > 0 {
		var downloadCount int64
		if err := db.Model(&model.DownloadRecord{}).Where("resource_type = ? AND resource_id = ?", "dataset", baseDataset.ID).Count(&downloadCount).Error; err != nil {
			return err
		}
		if downloadCount == 0 {
			downloads := []model.DownloadRecord{
				{UserID: demo.ID, ResourceType: "dataset", ResourceID: baseDataset.ID, ResourceTitle: baseDataset.Name},
				{UserID: admin.ID, ResourceType: "dataset", ResourceID: baseDataset.ID, ResourceTitle: baseDataset.Name},
			}
			if err := db.Create(&downloads).Error; err != nil {
				return err
			}
		}

		var sampleMediaCount int64
		if err := db.Model(&model.DatasetSample{}).Where("dataset_id = ? AND sample_type IN ?", baseDataset.ID, []string{"image", "video", "pointcloud"}).Count(&sampleMediaCount).Error; err != nil {
			return err
		}
		if sampleMediaCount == 0 {
			imageSamplePath, err := storage.WriteGeneratedFile("dataset-samples", "seed-inspection-preview.svg", `<svg xmlns="http://www.w3.org/2000/svg" width="640" height="360"><rect width="100%" height="100%" fill="#eaf3ff"/><text x="40" y="180" font-size="28" fill="#16324f">Inspection Image Preview</text></svg>`)
			if err != nil {
				return err
			}
			videoSamplePath, err := storage.WriteGeneratedFile("dataset-samples", "seed-inspection-video.mp4", "video placeholder")
			if err != nil {
				return err
			}
			pointCloudSamplePath, err := storage.WriteGeneratedFile("dataset-samples", "seed-inspection-pointcloud.ply", "ply\nformat ascii 1.0\ncomment placeholder point cloud")
			if err != nil {
				return err
			}
			samples := []model.DatasetSample{
				{DatasetID: baseDataset.ID, SampleType: "image", Title: "巡检图片预览", PreviewText: "图片样本预览", FilePath: imageSamplePath, FileName: "seed-inspection-preview.svg"},
				{DatasetID: baseDataset.ID, SampleType: "video", Title: "巡检视频占位样本", PreviewText: "视频样本已上传，可通过链接查看或下载。", FilePath: videoSamplePath, FileName: "seed-inspection-video.mp4"},
				{DatasetID: baseDataset.ID, SampleType: "pointcloud", Title: "点云占位样本", PreviewText: "点云样本已上传，当前提供占位预览与文件链接。", FilePath: pointCloudSamplePath, FileName: "seed-inspection-pointcloud.ply"},
			}
			if err := db.Create(&samples).Error; err != nil {
				return err
			}
		}
	}

	var wikiCount int64
	if err := db.Model(&model.WikiPage{}).Count(&wikiCount).Error; err != nil {
		return err
	}
	if wikiCount == 0 {
		page := model.WikiPage{
			Title:    "巡检任务接入 Wiki",
			Summary:  "记录移动作业和巡检模板接入的关键步骤。",
			Content:  "接入建议：先确认设备参数，再绑定模型、数据集和任务模板，最后联调告警链路。",
			Status:   model.StatusPublished,
			EditorID: demo.ID,
		}
		if err := db.Create(&page).Error; err != nil {
			return err
		}
		if err := db.Create(&model.WikiRevision{
			PageID:   page.ID,
			EditorID: demo.ID,
			Title:    page.Title,
			Summary:  page.Summary,
			Content:  page.Content,
			Comment:  "初始词条",
		}).Error; err != nil {
			return err
		}
	}

	var rewardCount int64
	if err := db.Model(&model.RewardLedger{}).Count(&rewardCount).Error; err != nil {
		return err
	}
	if rewardCount == 0 {
		rewards := []model.RewardLedger{
			{UserID: demo.ID, SourceType: "seed_skill", Points: 30, Remark: "示例技能发布"},
			{UserID: demo.ID, SourceType: "seed_discussion", Points: 10, Remark: "示例讨论创建"},
			{UserID: admin.ID, SourceType: "seed_admin", Points: 5, Remark: "示例运营积分"},
		}
		if err := db.Create(&rewards).Error; err != nil {
			return err
		}
	}
	for _, spec := range []seedRewardSpec{
		{UserID: integrator.ID, SourceType: "seed_integration", Points: 28, Remark: "接入实践整理"},
		{UserID: opsUser.ID, SourceType: "seed_ops", Points: 18, Remark: "首页运营与质检"},
	} {
		if err := ensureSeedReward(db, spec); err != nil {
			return err
		}
	}

	var benefitCount int64
	if err := db.Model(&model.RewardBenefit{}).Count(&benefitCount).Error; err != nil {
		return err
	}
	if benefitCount == 0 {
		benefits := []model.RewardBenefit{
			{Name: "首页创作者标识", Summary: "兑换后可在公开开发者页显示创作者标识。", CostPoints: 20, Active: true},
			{Name: "社区推荐位申请", Summary: "兑换后可获得一次社区首页推荐位申请资格。", CostPoints: 40, Active: true},
			{Name: "协作空间扩容", Summary: "兑换后可为一个协作空间增加更多成员名额。", CostPoints: 60, Active: true},
		}
		if err := db.Create(&benefits).Error; err != nil {
			return err
		}
	}

	var conversationCount int64
	if err := db.Model(&model.Conversation{}).Count(&conversationCount).Error; err != nil {
		return err
	}
	if conversationCount == 0 {
		direct := model.Conversation{
			Kind:  "direct",
			Key:   repository.DirectKey(admin.ID, demo.ID),
			Title: "管理员与开发者私信",
		}
		if err := db.Create(&direct).Error; err != nil {
			return err
		}
		if err := db.Create(&[]model.ConversationParticipant{
			{ConversationID: direct.ID, UserID: admin.ID},
			{ConversationID: direct.ID, UserID: demo.ID},
		}).Error; err != nil {
			return err
		}
		if err := db.Create(&model.Message{
			ConversationID: direct.ID,
			SenderID:       admin.ID,
			Content:        "欢迎加入开放社区，如需协作可以直接在这里沟通。",
		}).Error; err != nil {
			return err
		}

		workspaceConversation := model.Conversation{
			Kind:  "workspace",
			Key:   "workspace:1",
			Title: "MVP 协作空间",
		}
		if err := db.Create(&workspaceConversation).Error; err != nil {
			return err
		}
		workspace := model.Workspace{
			Name:           "MVP 协作空间",
			Summary:        "用于模型、数据集和模板联调的协作空间。",
			OwnerID:        admin.ID,
			ConversationID: workspaceConversation.ID,
		}
		if err := db.Create(&workspace).Error; err != nil {
			return err
		}
		workspaceConversation.Key = repository.WorkspaceKey(workspace.ID)
		workspaceConversation.WorkspaceID = &workspace.ID
		if err := db.Save(&workspaceConversation).Error; err != nil {
			return err
		}
		if err := db.Create(&[]model.WorkspaceMember{
			{WorkspaceID: workspace.ID, UserID: admin.ID, Role: "owner"},
			{WorkspaceID: workspace.ID, UserID: demo.ID, Role: "member"},
		}).Error; err != nil {
			return err
		}
		if err := db.Create(&[]model.ConversationParticipant{
			{ConversationID: workspaceConversation.ID, UserID: admin.ID},
			{ConversationID: workspaceConversation.ID, UserID: demo.ID},
		}).Error; err != nil {
			return err
		}
		if err := db.Create(&model.Message{
			ConversationID: workspaceConversation.ID,
			SenderID:       demo.ID,
			Content:        "已准备好巡检模板和数据集，可以开始联调。",
		}).Error; err != nil {
			return err
		}
	}
	return nil
}

type seedModelSpec struct {
	Name         string
	Summary      string
	Description  string
	Tags         string
	RobotType    string
	InputSpec    string
	OutputSpec   string
	License      string
	Dependencies string
	Downloads    int64
	FileName     string
	FileContent  string
	Changelog    string
}

type seedDatasetSpec struct {
	Name          string
	Summary       string
	Description   string
	Tags          string
	SampleCount   int
	Device        string
	Scene         string
	Privacy       string
	AgreementText string
	Downloads     int64
	FileName      string
	FileContent   string
	Changelog     string
}

type seedTemplateSpec struct {
	Name        string
	Summary     string
	Description string
	Category    string
	Scene       string
	Guide       string
	ResourceRef string
	UsageCount  int64
}

type seedApplicationCaseSpec struct {
	Title    string
	Summary  string
	Category string
	Guide    string
}

type seedAnnouncementSpec struct {
	Title   string
	Summary string
	Link    string
	Pinned  bool
}

type seedWikiSpec struct {
	Title    string
	Summary  string
	Content  string
	EditorID uint
	Comment  string
}

type seedRewardSpec struct {
	UserID     uint
	SourceType string
	Points     int
	Remark     string
}

func ensureSeedUser(db *gorm.DB, tokenManager *support.TokenManager, username, email, bio string) (*model.User, error) {
	var item model.User
	err := db.Where("email = ?", email).First(&item).Error
	if err == nil {
		updates := map[string]any{}
		if item.Username == "" {
			updates["username"] = username
		}
		if item.Bio == "" {
			updates["bio"] = bio
		}
		if len(updates) > 0 {
			if err := db.Model(&item).Updates(updates).Error; err != nil {
				return nil, err
			}
			if err := db.Where("email = ?", email).First(&item).Error; err != nil {
				return nil, err
			}
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.User{
		Username:     username,
		Email:        email,
		PasswordHash: tokenManager.HashPassword("Demo123!"),
		Role:         model.RoleUser,
		Bio:          bio,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedModel(db *gorm.DB, storage *support.LocalStorage, ownerID uint, spec seedModelSpec) (*model.ModelAsset, error) {
	var item model.ModelAsset
	err := db.Where("name = ?", spec.Name).First(&item).Error
	if err == nil {
		if item.Downloads < spec.Downloads || item.Status != model.StatusPublished {
			if err := db.Model(&item).Updates(map[string]any{
				"downloads": spec.Downloads,
				"status":    model.StatusPublished,
			}).Error; err != nil {
				return nil, err
			}
			if err := db.Where("name = ?", spec.Name).First(&item).Error; err != nil {
				return nil, err
			}
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	filePath, err := storage.WriteSeedFile("models", spec.FileName, spec.FileContent)
	if err != nil {
		return nil, err
	}
	item = model.ModelAsset{
		Name:         spec.Name,
		Summary:      spec.Summary,
		Description:  spec.Description,
		Tags:         spec.Tags,
		RobotType:    spec.RobotType,
		InputSpec:    spec.InputSpec,
		OutputSpec:   spec.OutputSpec,
		License:      spec.License,
		Dependencies: spec.Dependencies,
		Status:       model.StatusPublished,
		Downloads:    spec.Downloads,
		OwnerID:      ownerID,
		Versions: []model.ModelVersion{
			{Version: "v1.0.0", FilePath: filePath, FileName: spec.FileName, Changelog: spec.Changelog},
		},
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedDataset(db *gorm.DB, storage *support.LocalStorage, ownerID uint, spec seedDatasetSpec) (*model.Dataset, error) {
	var item model.Dataset
	err := db.Where("name = ?", spec.Name).First(&item).Error
	if err == nil {
		if item.Downloads < spec.Downloads || item.Status != model.StatusPublished {
			if err := db.Model(&item).Updates(map[string]any{
				"downloads": spec.Downloads,
				"status":    model.StatusPublished,
			}).Error; err != nil {
				return nil, err
			}
			if err := db.Where("name = ?", spec.Name).First(&item).Error; err != nil {
				return nil, err
			}
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	filePath, err := storage.WriteSeedFile("datasets", spec.FileName, spec.FileContent)
	if err != nil {
		return nil, err
	}
	item = model.Dataset{
		Name:          spec.Name,
		Summary:       spec.Summary,
		Description:   spec.Description,
		Tags:          spec.Tags,
		SampleCount:   spec.SampleCount,
		Device:        spec.Device,
		Scene:         spec.Scene,
		Privacy:       spec.Privacy,
		AgreementText: spec.AgreementText,
		Status:        model.StatusPublished,
		Downloads:     spec.Downloads,
		OwnerID:       ownerID,
		Versions: []model.DatasetVersion{
			{Version: "v1.0.0", FilePath: filePath, FileName: spec.FileName, Changelog: spec.Changelog},
		},
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedTemplate(db *gorm.DB, spec seedTemplateSpec) (*model.TaskTemplate, error) {
	var item model.TaskTemplate
	err := db.Where("name = ?", spec.Name).First(&item).Error
	if err == nil {
		if item.UsageCount < spec.UsageCount || item.Status != model.StatusPublished {
			if err := db.Model(&item).Updates(map[string]any{
				"usage_count": spec.UsageCount,
				"status":      model.StatusPublished,
			}).Error; err != nil {
				return nil, err
			}
			if err := db.Where("name = ?", spec.Name).First(&item).Error; err != nil {
				return nil, err
			}
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.TaskTemplate{
		Name:        spec.Name,
		Summary:     spec.Summary,
		Description: spec.Description,
		Category:    spec.Category,
		Scene:       spec.Scene,
		Guide:       spec.Guide,
		ResourceRef: spec.ResourceRef,
		UsageCount:  spec.UsageCount,
		Status:      model.StatusPublished,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedApplicationCase(db *gorm.DB, spec seedApplicationCaseSpec) (*model.ApplicationCase, error) {
	var item model.ApplicationCase
	err := db.Where("title = ?", spec.Title).First(&item).Error
	if err == nil {
		if item.Status != model.StatusPublished {
			if err := db.Model(&item).Update("status", model.StatusPublished).Error; err != nil {
				return nil, err
			}
			if err := db.Where("title = ?", spec.Title).First(&item).Error; err != nil {
				return nil, err
			}
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.ApplicationCase{
		Title:    spec.Title,
		Summary:  spec.Summary,
		Category: spec.Category,
		Guide:    spec.Guide,
		Status:   model.StatusPublished,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedAnnouncement(db *gorm.DB, spec seedAnnouncementSpec) error {
	var item model.Announcement
	err := db.Where("title = ?", spec.Title).First(&item).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	item = model.Announcement{
		Title:       spec.Title,
		Summary:     spec.Summary,
		Link:        spec.Link,
		Pinned:      spec.Pinned,
		PublishedAt: time.Now(),
	}
	return db.Create(&item).Error
}

func ensureSeedWikiPage(db *gorm.DB, spec seedWikiSpec) error {
	var item model.WikiPage
	err := db.Where("title = ?", spec.Title).First(&item).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	item = model.WikiPage{
		Title:    spec.Title,
		Summary:  spec.Summary,
		Content:  spec.Content,
		Status:   model.StatusPublished,
		EditorID: spec.EditorID,
	}
	if err := db.Create(&item).Error; err != nil {
		return err
	}
	return db.Create(&model.WikiRevision{
		PageID:   item.ID,
		EditorID: spec.EditorID,
		Title:    spec.Title,
		Summary:  spec.Summary,
		Content:  spec.Content,
		Comment:  spec.Comment,
	}).Error
}

func ensureSeedReward(db *gorm.DB, spec seedRewardSpec) error {
	var item model.RewardLedger
	err := db.Where("user_id = ? AND source_type = ? AND remark = ?", spec.UserID, spec.SourceType, spec.Remark).First(&item).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return db.Create(&model.RewardLedger{
		UserID:     spec.UserID,
		SourceType: spec.SourceType,
		Points:     spec.Points,
		Remark:     spec.Remark,
	}).Error
}

func ensureSeedFeaturedResource(db *gorm.DB, resourceType string, resourceID uint, sortOrder int, badgeLabel string) error {
	if resourceID == 0 {
		return nil
	}

	var item model.FeaturedResource
	err := db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).First(&item).Error
	if err == nil {
		return db.Model(&item).Updates(map[string]any{
			"badge_label": badgeLabel,
			"sort_order":  sortOrder,
			"enabled":     true,
		}).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return db.Create(&model.FeaturedResource{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		BadgeLabel:   badgeLabel,
		SortOrder:    sortOrder,
		Enabled:      true,
	}).Error
}

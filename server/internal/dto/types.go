package dto

import "time"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type ProfileUpdateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Bio      string `json:"bio" binding:"max=512"`
	Avatar   string `json:"avatar" binding:"max=255"`
}

type ResourceListQuery struct {
	Q         string `form:"q"`
	Tags      string `form:"tags"`
	RobotType string `form:"robot_type"`
	Status    string `form:"status"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Sort      string `form:"sort"`
}

type SearchQuery struct {
	Q             string `form:"q" binding:"required"`
	Type          string `form:"type"`
	Tags          string `form:"tags"`
	RobotType     string `form:"robot_type"`
	Sort          string `form:"sort"`
	UpdatedWithin int    `form:"updated_within"`
	Page          int    `form:"page"`
	PageSize      int    `form:"page_size"`
}

type FavoriteRequest struct {
	ResourceType string `json:"resource_type" binding:"required"`
	ResourceID   uint   `json:"resource_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
}

type ReviewDecisionRequest struct {
	Decision string `json:"decision" binding:"required,oneof=approved rejected"`
	Comment  string `json:"comment" binding:"max=512"`
}

type AnnouncementRequest struct {
	Title   string `json:"title" binding:"required,max=160"`
	Summary string `json:"summary" binding:"required,max=512"`
	Link    string `json:"link" binding:"max=255"`
	Pinned  bool   `json:"pinned"`
}

type ModuleSettingRequest struct {
	Enabled bool `json:"enabled"`
}

type FeaturedResourceRequest struct {
	ResourceType string `json:"resource_type" binding:"required,oneof=model dataset task-template application-case"`
	ResourceID   uint   `json:"resource_id" binding:"required"`
	SortOrder    int    `json:"sort_order" binding:"gte=0,lte=1000"`
	Enabled      bool   `json:"enabled"`
}

type TaskTemplateAdminRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=160"`
	Summary     string `json:"summary" binding:"required,min=4,max=255"`
	Description string `json:"description" binding:"required,min=10,max=5000"`
	Category    string `json:"category" binding:"required,max=128"`
	Scene       string `json:"scene" binding:"required,max=128"`
	Guide       string `json:"guide" binding:"required,min=4,max=5000"`
	ResourceRef string `json:"resource_ref" binding:"max=255"`
	Status      string `json:"status" binding:"required,oneof=draft pending published rejected"`
}

type ApplicationCaseAdminRequest struct {
	Title      string `json:"title" binding:"required,min=2,max=160"`
	Summary    string `json:"summary" binding:"required,min=4,max=255"`
	Category   string `json:"category" binding:"required,max=128"`
	Guide      string `json:"guide" binding:"required,min=4,max=5000"`
	CoverImage string `json:"cover_image" binding:"max=255"`
	Status     string `json:"status" binding:"required,oneof=draft pending published rejected"`
}

type DocumentCategoryRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=120"`
	DocType string `json:"doc_type" binding:"required,oneof=platform technical"`
}

type DocumentAdminRequest struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Title      string `json:"title" binding:"required,min=2,max=160"`
	Summary    string `json:"summary" binding:"required,min=4,max=255"`
	Content    string `json:"content" binding:"required,min=10,max=10000"`
	DocType    string `json:"doc_type" binding:"required,oneof=platform technical"`
	Status     string `json:"status" binding:"required,oneof=draft pending published rejected"`
}

type FAQRequest struct {
	Question string `json:"question" binding:"required,min=2,max=255"`
	Answer   string `json:"answer" binding:"required,min=2,max=10000"`
}

type VideoTutorialRequest struct {
	Title     string `json:"title" binding:"required,min=2,max=160"`
	Summary   string `json:"summary" binding:"required,min=4,max=255"`
	Link      string `json:"link" binding:"required,max=512"`
	Category  string `json:"category" binding:"required,max=120"`
	SortOrder int    `json:"sort_order" binding:"gte=0,lte=1000"`
	Active    bool   `json:"active"`
}

type AgreementTemplateRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=160"`
	Content   string `json:"content" binding:"required,min=4,max=10000"`
	SortOrder int    `json:"sort_order" binding:"gte=0,lte=1000"`
	Active    bool   `json:"active"`
}

type DatasetPrivacyOptionRequest struct {
	Code        string `json:"code" binding:"required,min=2,max=64"`
	Name        string `json:"name" binding:"required,min=2,max=120"`
	Description string `json:"description" binding:"max=255"`
	SortOrder   int    `json:"sort_order" binding:"gte=0,lte=1000"`
	Active      bool   `json:"active"`
}

type BatchStatusRequest struct {
	IDs    []uint `json:"ids" binding:"required,min=1,max=100"`
	Status string `json:"status" binding:"required,oneof=draft pending published rejected"`
}

type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1,max=100"`
}

type ModelCreateRequest struct {
	Name         string `form:"name" binding:"required,min=2,max=160"`
	Summary      string `form:"summary" binding:"required,min=4,max=255"`
	Description  string `form:"description" binding:"required,min=10,max=5000"`
	Tags         string `form:"tags" binding:"max=255"`
	RobotType    string `form:"robot_type" binding:"required,max=128"`
	InputSpec    string `form:"input_spec" binding:"required,max=255"`
	OutputSpec   string `form:"output_spec" binding:"required,max=255"`
	License      string `form:"license" binding:"required,max=120"`
	Dependencies string `form:"dependencies" binding:"max=255"`
	Version      string `form:"version" binding:"required,max=64"`
	Changelog    string `form:"changelog" binding:"required,max=512"`
}

type ModelVersionRequest struct {
	Version   string `form:"version" binding:"required,max=64"`
	Changelog string `form:"changelog" binding:"required,max=512"`
}

type DatasetCreateRequest struct {
	Name          string `form:"name" binding:"required,min=2,max=160"`
	Summary       string `form:"summary" binding:"required,min=4,max=255"`
	Description   string `form:"description" binding:"required,min=10,max=5000"`
	Tags          string `form:"tags" binding:"max=255"`
	SampleCount   int    `form:"sample_count" binding:"required,gte=1,lte=100000000"`
	Device        string `form:"device" binding:"required,max=128"`
	Scene         string `form:"scene" binding:"required,max=128"`
	Privacy       string `form:"privacy" binding:"required,max=64"`
	AgreementText string `form:"agreement_text" binding:"required,min=6,max=5000"`
	Version       string `form:"version" binding:"required,max=64"`
	Changelog     string `form:"changelog" binding:"required,max=512"`
	SamplePreview string `form:"sample_preview" binding:"max=5000"`
}

type DatasetVersionRequest struct {
	Version   string `form:"version" binding:"required,max=64"`
	Changelog string `form:"changelog" binding:"required,max=512"`
}

type DownloadPackageCreateRequest struct {
	Parts int `json:"parts" binding:"omitempty,gte=1,lte=10"`
}

type DatasetAccessRequestPayload struct {
	Reason string `json:"reason" binding:"required,min=4,max=2000"`
}

type DatasetAccessDecisionRequest struct {
	Decision string `json:"decision" binding:"required,oneof=approved rejected"`
	Comment  string `json:"comment" binding:"max=512"`
}

type UserSummary struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  UserSummary `json:"user"`
}

type AnnouncementItem struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Link        string    `json:"link"`
	Pinned      bool      `json:"pinned"`
	PublishedAt time.Time `json:"published_at"`
}

type ModuleSettingItem struct {
	ID        uint      `json:"id"`
	ModuleKey string    `json:"module_key"`
	Label     string    `json:"label"`
	Enabled   bool      `json:"enabled"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FeaturedResourceItem struct {
	ID           uint      `json:"id"`
	ResourceType string    `json:"resource_type"`
	ResourceID   uint      `json:"resource_id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Route        string    `json:"route"`
	SortOrder    int       `json:"sort_order"`
	Enabled      bool      `json:"enabled"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ResourceCard struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Summary     string    `json:"summary"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type"`
	Tags        []string  `json:"tags"`
	RobotType   string    `json:"robot_type,omitempty"`
	Downloads   int64     `json:"downloads"`
	Status      string    `json:"status,omitempty"`
	Owner       string    `json:"owner,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SearchItem struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Tags      []string  `json:"tags,omitempty"`
	Route     string    `json:"route"`
	ScoreHint int64     `json:"score_hint"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ModelDetail struct {
	ResourceCard
	InputSpec    string        `json:"input_spec"`
	OutputSpec   string        `json:"output_spec"`
	License      string        `json:"license"`
	Dependencies []string      `json:"dependencies"`
	Versions     []FileVersion `json:"versions"`
	Favorited    bool          `json:"favorited"`
}

type DatasetDetail struct {
	ResourceCard
	SampleCount   int             `json:"sample_count"`
	Device        string          `json:"device"`
	Scene         string          `json:"scene"`
	Privacy       string          `json:"privacy"`
	AgreementText string          `json:"agreement_text"`
	Samples       []DatasetSample `json:"samples"`
	Versions      []FileVersion   `json:"versions"`
	Favorited     bool            `json:"favorited"`
}

type DatasetSample struct {
	ID          uint   `json:"id"`
	SampleType  string `json:"sample_type"`
	Title       string `json:"title"`
	PreviewText string `json:"preview_text"`
	PreviewURL  string `json:"preview_url,omitempty"`
	FileName    string `json:"file_name,omitempty"`
}

type FileVersion struct {
	ID        uint      `json:"id"`
	Version   string    `json:"version"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	Changelog string    `json:"changelog"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskTemplateItem struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Scene       string    `json:"scene"`
	Guide       string    `json:"guide"`
	ResourceRef []string  `json:"resource_ref"`
	UsageCount  int64     `json:"usage_count"`
	Status      string    `json:"status,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ApplicationCaseItem struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Category   string    `json:"category"`
	Guide      string    `json:"guide"`
	CoverImage string    `json:"cover_image,omitempty"`
	Status     string    `json:"status,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DocumentCategoryItem struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	DocType string `json:"doc_type"`
}

type DocumentItem struct {
	ID         uint      `json:"id"`
	CategoryID uint      `json:"category_id"`
	Category   string    `json:"category"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Content    string    `json:"content"`
	DocType    string    `json:"doc_type"`
	Status     string    `json:"status,omitempty"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type FAQItem struct {
	ID        uint      `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	UpdatedAt time.Time `json:"updated_at"`
}

type VideoTutorialItem struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Link      string    `json:"link"`
	Category  string    `json:"category"`
	SortOrder int       `json:"sort_order"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AgreementTemplateItem struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	SortOrder int       `json:"sort_order"`
	Active    bool      `json:"active"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DatasetPrivacyOptionItem struct {
	ID          uint      `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SortOrder   int       `json:"sort_order"`
	Active      bool      `json:"active"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DatasetOptionsResponse struct {
	AgreementTemplates []AgreementTemplateItem    `json:"agreement_templates"`
	PrivacyOptions     []DatasetPrivacyOptionItem `json:"privacy_options"`
}

type HomeResponse struct {
	PlatformIntro    string                `json:"platform_intro"`
	Highlights       []string              `json:"highlights"`
	Announcements    []AnnouncementItem    `json:"announcements"`
	HotModels        []ResourceCard        `json:"hot_models"`
	HotDatasets      []ResourceCard        `json:"hot_datasets"`
	TaskTemplates    []TaskTemplateItem    `json:"task_templates"`
	ApplicationCases []ApplicationCaseItem `json:"application_cases"`
	ModuleSettings   []ModuleSettingItem   `json:"module_settings"`
}

type DashboardResponse struct {
	Users             int64 `json:"users"`
	PublishedModels   int64 `json:"published_models"`
	PendingModels     int64 `json:"pending_models"`
	PublishedDatasets int64 `json:"published_datasets"`
	PendingDatasets   int64 `json:"pending_datasets"`
	Announcements     int64 `json:"announcements"`
}

type ReviewItem struct {
	ID        uint      `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Status    string    `json:"status"`
	Owner     string    `json:"owner"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ModelEvaluationRequest struct {
	Benchmark string  `json:"benchmark" binding:"required,max=160"`
	Summary   string  `json:"summary" binding:"required,min=4,max=255"`
	Score     float64 `json:"score" binding:"required,gte=0,lte=100"`
	Notes     string  `json:"notes" binding:"max=5000"`
}

type ModelEvaluationItem struct {
	ID        uint      `json:"id"`
	Benchmark string    `json:"benchmark"`
	Summary   string    `json:"summary"`
	Score     float64   `json:"score"`
	Notes     string    `json:"notes"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type RatingRequest struct {
	Score    int    `json:"score" binding:"required,gte=1,lte=5"`
	Feedback string `json:"feedback" binding:"max=512"`
}

type RatingItem struct {
	ID        uint      `json:"id"`
	Score     int       `json:"score"`
	Feedback  string    `json:"feedback"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type RatingSummary struct {
	Average float64      `json:"average"`
	Count   int64        `json:"count"`
	Items   []RatingItem `json:"items"`
}

type CommentRequest struct {
	Content  string `json:"content" binding:"required,min=2,max=5000"`
	ParentID *uint  `json:"parent_id"`
}

type CommentItem struct {
	ID        uint      `json:"id"`
	ParentID  *uint     `json:"parent_id,omitempty"`
	Content   string    `json:"content"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"`
	CreatedAt time.Time `json:"created_at"`
}

type DatasetStatsResponse struct {
	DatasetID     uint             `json:"dataset_id"`
	DownloadCount int64            `json:"download_count"`
	SampleCount   int              `json:"sample_count"`
	SampleTypeMix []StatValueItem  `json:"sample_type_mix"`
	DownloadTrend []TrendValueItem `json:"download_trend"`
}

type DownloadPackageTaskItem struct {
	ID         uint      `json:"id"`
	DatasetID  uint      `json:"dataset_id"`
	Status     string    `json:"status"`
	BundleURL  string    `json:"bundle_url,omitempty"`
	PartLinks  []string  `json:"part_links"`
	TotalParts int       `json:"total_parts"`
	CreatedAt  time.Time `json:"created_at"`
}

type DatasetAccessRequestItem struct {
	ID            uint       `json:"id"`
	DatasetID     uint       `json:"dataset_id"`
	UserID        uint       `json:"user_id"`
	UserName      string     `json:"user_name,omitempty"`
	Reason        string     `json:"reason"`
	Status        string     `json:"status"`
	ReviewComment string     `json:"review_comment"`
	ReviewedAt    *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

type StatValueItem struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type TrendValueItem struct {
	Label string `json:"label"`
	Value int64  `json:"value"`
}

type SkillCreateRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=160"`
	Summary     string `json:"summary" binding:"required,min=4,max=255"`
	Description string `json:"description" binding:"required,min=10,max=5000"`
	Category    string `json:"category" binding:"required,max=128"`
	Scene       string `json:"scene" binding:"required,max=128"`
	Guide       string `json:"guide" binding:"required,min=4,max=5000"`
	ResourceRef string `json:"resource_ref" binding:"max=255"`
}

type SkillItem struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Summary      string        `json:"summary"`
	Description  string        `json:"description"`
	Category     string        `json:"category"`
	Scene        string        `json:"scene"`
	Guide        string        `json:"guide"`
	ResourceRef  []string      `json:"resource_ref"`
	Status       string        `json:"status"`
	ForkedFromID *uint         `json:"forked_from_id,omitempty"`
	UsageCount   int64         `json:"usage_count"`
	OwnerID      uint          `json:"owner_id"`
	OwnerName    string        `json:"owner_name"`
	Rating       RatingSummary `json:"rating"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type DiscussionRequest struct {
	Title    string `json:"title" binding:"required,min=2,max=160"`
	Summary  string `json:"summary" binding:"required,min=4,max=255"`
	Content  string `json:"content" binding:"required,min=10,max=5000"`
	Category string `json:"category" binding:"required,max=128"`
}

type DiscussionItem struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	Category     string    `json:"category"`
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name"`
	CommentCount int64     `json:"comment_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type FollowItem struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"created_at"`
}

type FollowStats struct {
	Follows   int64 `json:"follows"`
	Followers int64 `json:"followers"`
}

type SearchHotItem struct {
	Query string `json:"query"`
	Count int64  `json:"count"`
}

type UserContributionPayload struct {
	Skills      []SkillItem      `json:"skills"`
	Discussions []DiscussionItem `json:"discussions"`
}

type AdminCommunityOverview struct {
	Skills      int64 `json:"skills"`
	Discussions int64 `json:"discussions"`
	Comments    int64 `json:"comments"`
}

type AdminSkillModerationItem struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Summary   string    `json:"summary"`
	Status    string    `json:"status"`
	OwnerID   uint      `json:"owner_id"`
	OwnerName string    `json:"owner_name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AdminDiscussionModerationItem struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name"`
	CommentCount int64     `json:"comment_count"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AdminCommentModerationItem struct {
	ID           uint      `json:"id"`
	ResourceType string    `json:"resource_type"`
	ResourceID   uint      `json:"resource_id"`
	UserID       uint      `json:"user_id"`
	UserName     string    `json:"user_name"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

type DeveloperVerificationRequest struct {
	VerificationType string `json:"verification_type" binding:"required,oneof=personal enterprise"`
	RealName         string `json:"real_name" binding:"required,min=2,max=120"`
	Organization     string `json:"organization" binding:"max=160"`
	Materials        string `json:"materials" binding:"max=5000"`
	Reason           string `json:"reason" binding:"required,min=6,max=5000"`
}

type EnterpriseVerificationRequest struct {
	RealName     string `json:"real_name" binding:"required,min=2,max=120"`
	Organization string `json:"organization" binding:"required,min=2,max=160"`
	Materials    string `json:"materials" binding:"max=5000"`
	Reason       string `json:"reason" binding:"required,min=6,max=5000"`
}

type VerificationStatusItem struct {
	ID               uint       `json:"id"`
	UserID           uint       `json:"user_id"`
	VerificationType string     `json:"verification_type"`
	RealName         string     `json:"real_name"`
	Organization     string     `json:"organization"`
	Materials        string     `json:"materials"`
	Reason           string     `json:"reason"`
	Status           string     `json:"status"`
	ReviewComment    string     `json:"review_comment"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type VerificationDecisionRequest struct {
	Decision string `json:"decision" binding:"required,oneof=approved rejected"`
	Comment  string `json:"comment" binding:"max=512"`
}

type AdminVerificationItem struct {
	ID               uint       `json:"id"`
	UserID           uint       `json:"user_id"`
	UserName         string     `json:"user_name"`
	VerificationType string     `json:"verification_type"`
	RealName         string     `json:"real_name"`
	Organization     string     `json:"organization"`
	Reason           string     `json:"reason"`
	Status           string     `json:"status"`
	ReviewComment    string     `json:"review_comment"`
	ReviewedAt       *time.Time `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type WikiPageRequest struct {
	Title   string `json:"title" binding:"required,min=2,max=160"`
	Summary string `json:"summary" binding:"required,min=4,max=255"`
	Content string `json:"content" binding:"required,min=2,max=10000"`
	Comment string `json:"comment" binding:"max=512"`
}

type WikiPageItem struct {
	ID            uint       `json:"id"`
	Title         string     `json:"title"`
	Summary       string     `json:"summary"`
	Content       string     `json:"content"`
	Status        string     `json:"status"`
	Locked        bool       `json:"locked"`
	LockedBy      *uint      `json:"locked_by,omitempty"`
	LockedByName  string     `json:"locked_by_name,omitempty"`
	LockedAt      *time.Time `json:"locked_at,omitempty"`
	RevisionCount int64      `json:"revision_count,omitempty"`
	EditorID      uint       `json:"editor_id"`
	EditorName    string     `json:"editor_name"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type AdminWikiRollbackRequest struct {
	RevisionID uint   `json:"revision_id" binding:"required"`
	Comment    string `json:"comment" binding:"max=512"`
}

type WikiRevisionItem struct {
	ID         uint      `json:"id"`
	PageID     uint      `json:"page_id"`
	EditorID   uint      `json:"editor_id"`
	EditorName string    `json:"editor_name"`
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Content    string    `json:"content"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

type RewardSummary struct {
	Points int64 `json:"points"`
}

type RewardLedgerItem struct {
	ID         uint      `json:"id"`
	SourceType string    `json:"source_type"`
	Points     int       `json:"points"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}

type ContributorRankingItem struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Points   int64  `json:"points"`
}

type RewardBenefitItem struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Summary    string `json:"summary"`
	CostPoints int    `json:"cost_points"`
	Active     bool   `json:"active"`
}

type AdminRewardOverview struct {
	Benefits       int64 `json:"benefits"`
	ActiveBenefits int64 `json:"active_benefits"`
	Redemptions    int64 `json:"redemptions"`
	LedgerEntries  int64 `json:"ledger_entries"`
	NetPoints      int64 `json:"net_points"`
}

type AdminRewardBenefitRequest struct {
	Name       string `json:"name" binding:"required,min=2,max=160"`
	Summary    string `json:"summary" binding:"required,min=4,max=255"`
	CostPoints int    `json:"cost_points" binding:"required,gte=1,lte=100000"`
	Active     bool   `json:"active"`
}

type RewardRedeemRequest struct {
	BenefitID uint `json:"benefit_id" binding:"required"`
}

type RewardRedemptionItem struct {
	ID          uint      `json:"id"`
	BenefitID   uint      `json:"benefit_id"`
	BenefitName string    `json:"benefit_name"`
	CostPoints  int       `json:"cost_points"`
	CreatedAt   time.Time `json:"created_at"`
}

type AdminRewardAdjustmentRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Points int    `json:"points" binding:"required,ne=0,gte=-100000,lte=100000"`
	Remark string `json:"remark" binding:"required,min=2,max=255"`
}

type AdminRewardAdjustmentItem struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"`
	Points    int       `json:"points"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminOperationLogItem struct {
	ID           uint      `json:"id"`
	AdminUserID  uint      `json:"admin_user_id"`
	AdminName    string    `json:"admin_name"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   uint      `json:"resource_id"`
	Summary      string    `json:"summary"`
	Detail       string    `json:"detail"`
	CreatedAt    time.Time `json:"created_at"`
}

type OpenAPIEndpointItem struct {
	Method  string `json:"method"`
	Path    string `json:"path"`
	Summary string `json:"summary"`
	Auth    string `json:"auth"`
}

type CLIExampleItem struct {
	Title   string `json:"title"`
	Command string `json:"command"`
}

type OpenAPISpecResponse struct {
	Version       string                `json:"version"`
	BaseURL       string                `json:"base_url"`
	Overview      string                `json:"overview"`
	Endpoints     []OpenAPIEndpointItem `json:"endpoints"`
	WebhookEvents []string              `json:"webhook_events"`
	CurlExamples  []CLIExampleItem      `json:"curl_examples"`
}

type WebhookSubscriptionRequest struct {
	Name      string   `json:"name" binding:"required,min=2,max=120"`
	TargetURL string   `json:"target_url" binding:"required,max=512"`
	Secret    string   `json:"secret" binding:"required,min=8,max=120"`
	Events    []string `json:"events" binding:"required,min=1,max=12"`
}

type WebhookSubscriptionItem struct {
	ID               uint       `json:"id"`
	Name             string     `json:"name"`
	TargetURL        string     `json:"target_url"`
	Events           []string   `json:"events"`
	LastStatus       string     `json:"last_status"`
	LastResponseCode int        `json:"last_response_code"`
	LastError        string     `json:"last_error"`
	LastDeliveredAt  *time.Time `json:"last_delivered_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

type WebhookDeliveryItem struct {
	ID           uint      `json:"id"`
	Event        string    `json:"event"`
	Status       string    `json:"status"`
	ResponseCode int       `json:"response_code"`
	ResponseBody string    `json:"response_body"`
	CreatedAt    time.Time `json:"created_at"`
}

type MessageSendRequest struct {
	ConversationID  uint   `json:"conversation_id"`
	RecipientUserID uint   `json:"recipient_user_id"`
	Content         string `json:"content" binding:"required,min=1,max=5000"`
}

type MessageItem struct {
	ID             uint      `json:"id"`
	ConversationID uint      `json:"conversation_id"`
	SenderID       uint      `json:"sender_id"`
	SenderName     string    `json:"sender_name"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

type ConversationItem struct {
	ID               uint      `json:"id"`
	Kind             string    `json:"kind"`
	Title            string    `json:"title"`
	WorkspaceID      *uint     `json:"workspace_id,omitempty"`
	ParticipantNames []string  `json:"participant_names"`
	LatestMessage    string    `json:"latest_message"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type WorkspaceCreateRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=160"`
	Summary string `json:"summary" binding:"required,min=4,max=255"`
}

type WorkspaceMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

type WorkspaceMessageRequest struct {
	Content string `json:"content" binding:"required,min=1,max=5000"`
}

type WorkspaceItem struct {
	ID             uint      `json:"id"`
	Name           string    `json:"name"`
	Summary        string    `json:"summary"`
	OwnerID        uint      `json:"owner_id"`
	ConversationID uint      `json:"conversation_id"`
	MemberCount    int64     `json:"member_count"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WorkspaceDetail struct {
	WorkspaceItem
	Messages []MessageItem `json:"messages"`
	Members  []FollowItem  `json:"members"`
}

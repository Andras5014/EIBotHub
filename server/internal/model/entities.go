package model

import "time"

const (
	RoleUser      = "user"
	RoleDeveloper = "developer"
	RoleOperator  = "operator"
	RoleReviewer  = "reviewer"
	RoleAdmin     = "admin"

	StatusDraft     = "draft"
	StatusPending   = "pending"
	StatusPublished = "published"
	StatusRejected  = "rejected"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:64;uniqueIndex;not null"`
	Email        string `gorm:"size:128;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string `gorm:"size:32;not null;default:user"`
	Bio          string `gorm:"size:512"`
	Avatar       string `gorm:"size:255"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Announcement struct {
	ID          uint      `gorm:"primaryKey"`
	Title       string    `gorm:"size:160;not null"`
	Summary     string    `gorm:"size:512;not null"`
	Link        string    `gorm:"size:255"`
	Pinned      bool      `gorm:"not null;default:false"`
	PublishedAt time.Time `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type HomeModuleSetting struct {
	ID        uint   `gorm:"primaryKey"`
	ModuleKey string `gorm:"size:64;uniqueIndex;not null"`
	SortOrder int    `gorm:"not null;default:0"`
	Enabled   bool   `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type HomeHighlight struct {
	ID        uint   `gorm:"primaryKey"`
	Text      string `gorm:"size:255;not null"`
	SortOrder int    `gorm:"not null;default:0"`
	Enabled   bool   `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type HomeHeroConfig struct {
	ID              uint   `gorm:"primaryKey"`
	Tagline         string `gorm:"size:120;not null"`
	Title           string `gorm:"size:255;not null"`
	Description     string `gorm:"size:1024;not null"`
	PrimaryButton   string `gorm:"size:64;not null"`
	SecondaryButton string `gorm:"size:64;not null"`
	SearchButton    string `gorm:"size:64;not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RankingConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:160;not null"`
	Subtitle  string `gorm:"size:255"`
	Limit     int    `gorm:"not null;default:5"`
	Enabled   bool   `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ScenePageConfig struct {
	ID          uint   `gorm:"primaryKey"`
	Slug        string `gorm:"size:64;uniqueIndex;not null"`
	Name        string `gorm:"size:120;not null"`
	Tagline     string `gorm:"size:120"`
	Summary     string `gorm:"size:255;not null"`
	Description string `gorm:"size:1024"`
	SortOrder   int    `gorm:"not null;default:0"`
	Enabled     bool   `gorm:"not null;default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FeaturedResource struct {
	ID           uint   `gorm:"primaryKey"`
	ResourceType string `gorm:"size:32;index;not null"`
	ResourceID   uint   `gorm:"index;not null"`
	BadgeLabel   string `gorm:"size:64"`
	SortOrder    int    `gorm:"not null;default:0"`
	Enabled      bool   `gorm:"not null;default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ModelAsset struct {
	ID           uint           `gorm:"primaryKey"`
	Name         string         `gorm:"size:160;index;not null"`
	Summary      string         `gorm:"size:255;not null"`
	Description  string         `gorm:"type:text;not null"`
	Tags         string         `gorm:"size:255"`
	RecommendTag string         `gorm:"size:64"`
	RobotType    string         `gorm:"size:128;index"`
	InputSpec    string         `gorm:"size:255"`
	OutputSpec   string         `gorm:"size:255"`
	License      string         `gorm:"size:120"`
	Dependencies string         `gorm:"size:255"`
	Status       string         `gorm:"size:32;index;not null;default:draft"`
	Downloads    int64          `gorm:"not null;default:0"`
	OwnerID      uint           `gorm:"index;not null"`
	Owner        User           `gorm:"foreignKey:OwnerID"`
	Versions     []ModelVersion `gorm:"foreignKey:ModelID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ModelVersion struct {
	ID         uint   `gorm:"primaryKey"`
	ModelID    uint   `gorm:"index;not null"`
	Version    string `gorm:"size:64;not null"`
	FilePath   string `gorm:"size:255"`
	FileName   string `gorm:"size:255"`
	Changelog  string `gorm:"size:512"`
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type Dataset struct {
	ID            uint             `gorm:"primaryKey"`
	Name          string           `gorm:"size:160;index;not null"`
	Summary       string           `gorm:"size:255;not null"`
	Description   string           `gorm:"type:text;not null"`
	Tags          string           `gorm:"size:255"`
	SampleCount   int              `gorm:"not null;default:0"`
	Device        string           `gorm:"size:128"`
	Scene         string           `gorm:"size:128;index"`
	Privacy       string           `gorm:"size:64"`
	AgreementText string           `gorm:"type:text"`
	Status        string           `gorm:"size:32;index;not null;default:draft"`
	Downloads     int64            `gorm:"not null;default:0"`
	OwnerID       uint             `gorm:"index;not null"`
	Owner         User             `gorm:"foreignKey:OwnerID"`
	Versions      []DatasetVersion `gorm:"foreignKey:DatasetID"`
	Samples       []DatasetSample  `gorm:"foreignKey:DatasetID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type DatasetVersion struct {
	ID         uint   `gorm:"primaryKey"`
	DatasetID  uint   `gorm:"index;not null"`
	Version    string `gorm:"size:64;not null"`
	FilePath   string `gorm:"size:255"`
	FileName   string `gorm:"size:255"`
	Changelog  string `gorm:"size:512"`
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type DatasetSample struct {
	ID          uint   `gorm:"primaryKey"`
	DatasetID   uint   `gorm:"index;not null"`
	SampleType  string `gorm:"size:64;not null"`
	Title       string `gorm:"size:128;not null"`
	PreviewText string `gorm:"size:512"`
	FilePath    string `gorm:"size:255"`
	FileName    string `gorm:"size:255"`
	CreatedAt   time.Time
}

type DownloadPackageTask struct {
	ID         uint   `gorm:"primaryKey"`
	DatasetID  uint   `gorm:"index;not null"`
	UserID     uint   `gorm:"index;not null"`
	Status     string `gorm:"size:32;index;not null;default:ready"`
	BundlePath string `gorm:"size:255"`
	PartLinks  string `gorm:"type:text"`
	TotalParts int    `gorm:"not null;default:0"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DatasetAccessRequest struct {
	ID                uint   `gorm:"primaryKey"`
	DatasetID         uint   `gorm:"index;not null"`
	UserID            uint   `gorm:"index;not null"`
	Reason            string `gorm:"type:text;not null"`
	Status            string `gorm:"size:32;index;not null;default:pending"`
	ReviewComment     string `gorm:"size:512"`
	ApprovalStage     int    `gorm:"not null;default:0"`
	RequiredApprovals int    `gorm:"not null;default:1"`
	ApproverIDs       string `gorm:"size:255"`
	ApprovalExpiresAt *time.Time
	DownloadLimit     int   `gorm:"not null;default:0"`
	DownloadCount     int   `gorm:"not null;default:0"`
	ReviewedBy        *uint `gorm:"index"`
	ReviewedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type TaskTemplate struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:160;index;not null"`
	Summary     string `gorm:"size:255;not null"`
	Description string `gorm:"type:text;not null"`
	Category    string `gorm:"size:128;index"`
	Scene       string `gorm:"size:128;index"`
	Guide       string `gorm:"type:text"`
	ResourceRef string `gorm:"size:255"`
	UsageCount  int64  `gorm:"not null;default:0"`
	Status      string `gorm:"size:32;index;not null;default:published"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ApplicationCase struct {
	ID         uint   `gorm:"primaryKey"`
	Title      string `gorm:"size:160;index;not null"`
	Summary    string `gorm:"size:255;not null"`
	Category   string `gorm:"size:128;index"`
	Guide      string `gorm:"type:text"`
	CoverImage string `gorm:"size:255"`
	Status     string `gorm:"size:32;index;not null;default:published"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type DocumentCategory struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:120;uniqueIndex;not null"`
	DocType   string `gorm:"size:32;index;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Document struct {
	ID         uint             `gorm:"primaryKey"`
	CategoryID uint             `gorm:"index;not null"`
	Category   DocumentCategory `gorm:"foreignKey:CategoryID"`
	Title      string           `gorm:"size:160;index;not null"`
	Summary    string           `gorm:"size:255;not null"`
	Content    string           `gorm:"type:text;not null"`
	DocType    string           `gorm:"size:32;index;not null"`
	Status     string           `gorm:"size:32;index;not null;default:published"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type FAQ struct {
	ID        uint   `gorm:"primaryKey"`
	Question  string `gorm:"size:255;not null"`
	Answer    string `gorm:"type:text;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VideoTutorial struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:160;not null"`
	Summary   string `gorm:"size:255;not null"`
	Link      string `gorm:"size:512;not null"`
	Category  string `gorm:"size:120;index"`
	SortOrder int    `gorm:"not null;default:0"`
	Active    bool   `gorm:"index;not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AgreementTemplate struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:160;not null"`
	Content   string `gorm:"type:text;not null"`
	SortOrder int    `gorm:"not null;default:0"`
	Active    bool   `gorm:"index;not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DatasetPrivacyOption struct {
	ID          uint   `gorm:"primaryKey"`
	Code        string `gorm:"size:64;uniqueIndex;not null"`
	Name        string `gorm:"size:120;not null"`
	Description string `gorm:"size:255"`
	SortOrder   int    `gorm:"not null;default:0"`
	Active      bool   `gorm:"index;not null;default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Favorite struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `gorm:"index;not null"`
	ResourceType  string `gorm:"size:32;index;not null"`
	ResourceID    uint   `gorm:"index;not null"`
	ResourceTitle string `gorm:"size:160;not null"`
	CreatedAt     time.Time
}

type DownloadRecord struct {
	ID            uint   `gorm:"primaryKey"`
	UserID        uint   `gorm:"index;not null"`
	ResourceType  string `gorm:"size:32;index;not null"`
	ResourceID    uint   `gorm:"index;not null"`
	ResourceTitle string `gorm:"size:160;not null"`
	CreatedAt     time.Time
}

type FileObject struct {
	ID            uint   `gorm:"primaryKey"`
	ObjectKey     string `gorm:"size:255;uniqueIndex;not null"`
	OriginalName  string `gorm:"size:255;not null"`
	MIMEType      string `gorm:"size:128"`
	SizeBytes     int64  `gorm:"not null;default:0"`
	StorageDriver string `gorm:"size:32;not null;default:local"`
	Scope         string `gorm:"size:64;index"`
	UploadedBy    uint   `gorm:"index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Notification struct {
	ID        uint   `gorm:"primaryKey"`
	UserID    uint   `gorm:"index;not null"`
	Type      string `gorm:"size:64;not null"`
	Title     string `gorm:"size:160;not null"`
	Content   string `gorm:"type:text;not null"`
	Link      string `gorm:"size:255"`
	Read      bool   `gorm:"index;not null;default:false"`
	CreatedAt time.Time
}

type AgreementRecord struct {
	ID         uint      `gorm:"primaryKey"`
	UserID     uint      `gorm:"index;not null"`
	DatasetID  uint      `gorm:"index;not null"`
	AcceptedAt time.Time `gorm:"not null"`
}

type ReviewLog struct {
	ID           uint   `gorm:"primaryKey"`
	ResourceType string `gorm:"size:32;index;not null"`
	ResourceID   uint   `gorm:"index;not null"`
	ReviewerID   uint   `gorm:"index;not null"`
	Decision     string `gorm:"size:32;not null"`
	Comment      string `gorm:"size:512"`
	CreatedAt    time.Time
}

type AdminOperationLog struct {
	ID           uint   `gorm:"primaryKey"`
	AdminUserID  uint   `gorm:"index;not null"`
	Action       string `gorm:"size:64;index;not null"`
	ResourceType string `gorm:"size:32;index;not null"`
	ResourceID   uint   `gorm:"index;not null"`
	Summary      string `gorm:"size:255;not null"`
	Detail       string `gorm:"type:text"`
	CreatedAt    time.Time
}

type ModelEvaluation struct {
	ID        uint    `gorm:"primaryKey"`
	ModelID   uint    `gorm:"index;not null"`
	UserID    uint    `gorm:"index;not null"`
	Benchmark string  `gorm:"size:160;not null"`
	Summary   string  `gorm:"size:255;not null"`
	Score     float64 `gorm:"not null"`
	Notes     string  `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ResourceRating struct {
	ID           uint   `gorm:"primaryKey"`
	ResourceType string `gorm:"size:32;index:idx_rating_target,priority:1;not null"`
	ResourceID   uint   `gorm:"index:idx_rating_target,priority:2;not null"`
	UserID       uint   `gorm:"index:idx_rating_user_target,priority:1;not null"`
	Score        int    `gorm:"not null"`
	Feedback     string `gorm:"size:512"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type ResourceComment struct {
	ID           uint   `gorm:"primaryKey"`
	ResourceType string `gorm:"size:32;index:idx_comment_target,priority:1;not null"`
	ResourceID   uint   `gorm:"index:idx_comment_target,priority:2;not null"`
	UserID       uint   `gorm:"index;not null"`
	ParentID     *uint  `gorm:"index"`
	Content      string `gorm:"type:text;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Skill struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:160;index;not null"`
	Summary      string `gorm:"size:255;not null"`
	Description  string `gorm:"type:text;not null"`
	Category     string `gorm:"size:128;index"`
	Scene        string `gorm:"size:128;index"`
	Guide        string `gorm:"type:text"`
	ResourceRef  string `gorm:"size:255"`
	Status       string `gorm:"size:32;index;not null;default:published"`
	ForkedFromID *uint  `gorm:"index"`
	UsageCount   int64  `gorm:"not null;default:0"`
	OwnerID      uint   `gorm:"index;not null"`
	Owner        User   `gorm:"foreignKey:OwnerID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Discussion struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:160;index;not null"`
	Summary   string `gorm:"size:255;not null"`
	Content   string `gorm:"type:text;not null"`
	Category  string `gorm:"size:128;index"`
	UserID    uint   `gorm:"index;not null"`
	User      User   `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Follow struct {
	ID             uint `gorm:"primaryKey"`
	FollowerID     uint `gorm:"uniqueIndex:idx_follow_pair,priority:1;not null"`
	FollowedUserID uint `gorm:"uniqueIndex:idx_follow_pair,priority:2;not null"`
	CreatedAt      time.Time
}

type SearchRecord struct {
	ID         uint   `gorm:"primaryKey"`
	Query      string `gorm:"size:160;index;not null"`
	SearchType string `gorm:"size:32;index"`
	CreatedAt  time.Time
}

type SearchKeywordConfig struct {
	ID          uint   `gorm:"primaryKey"`
	Query       string `gorm:"size:160;index;not null"`
	KeywordType string `gorm:"size:32;index;not null"`
	SortOrder   int    `gorm:"not null;default:0"`
	Enabled     bool   `gorm:"not null;default:true"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type FilterOptionConfig struct {
	ID        uint   `gorm:"primaryKey"`
	Kind      string `gorm:"size:64;index;not null"`
	Value     string `gorm:"size:160;index;not null"`
	SortOrder int    `gorm:"not null;default:0"`
	Enabled   bool   `gorm:"not null;default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeveloperVerification struct {
	ID               uint   `gorm:"primaryKey"`
	UserID           uint   `gorm:"index;not null"`
	User             User   `gorm:"foreignKey:UserID"`
	VerificationType string `gorm:"size:32;index;not null"`
	RealName         string `gorm:"size:120;not null"`
	Organization     string `gorm:"size:160"`
	Materials        string `gorm:"type:text"`
	Reason           string `gorm:"type:text;not null"`
	Status           string `gorm:"size:32;index;not null;default:pending"`
	ReviewComment    string `gorm:"size:512"`
	ReviewedBy       *uint  `gorm:"index"`
	ReviewedAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type WikiPage struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:160;index;not null"`
	Summary   string `gorm:"size:255;not null"`
	Content   string `gorm:"type:text;not null"`
	Status    string `gorm:"size:32;index;not null;default:published"`
	Locked    bool   `gorm:"index;not null;default:false"`
	LockedBy  *uint  `gorm:"index"`
	LockedAt  *time.Time
	EditorID  uint           `gorm:"index;not null"`
	Editor    User           `gorm:"foreignKey:EditorID"`
	Revisions []WikiRevision `gorm:"foreignKey:PageID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WikiRevision struct {
	ID        uint   `gorm:"primaryKey"`
	PageID    uint   `gorm:"index;not null"`
	EditorID  uint   `gorm:"index;not null"`
	Title     string `gorm:"size:160;not null"`
	Summary   string `gorm:"size:255;not null"`
	Content   string `gorm:"type:text;not null"`
	Comment   string `gorm:"size:512"`
	CreatedAt time.Time
}

type RewardLedger struct {
	ID         uint   `gorm:"primaryKey"`
	UserID     uint   `gorm:"index;not null"`
	User       User   `gorm:"foreignKey:UserID"`
	SourceType string `gorm:"size:64;index;not null"`
	Points     int    `gorm:"not null"`
	Remark     string `gorm:"size:255"`
	CreatedAt  time.Time
}

type RewardBenefit struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"size:160;not null"`
	Summary    string `gorm:"size:255;not null"`
	CostPoints int    `gorm:"not null"`
	Active     bool   `gorm:"index;not null;default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type RewardRedemption struct {
	ID         uint          `gorm:"primaryKey"`
	UserID     uint          `gorm:"index;not null"`
	BenefitID  uint          `gorm:"index;not null"`
	Benefit    RewardBenefit `gorm:"foreignKey:BenefitID"`
	CostPoints int           `gorm:"not null"`
	CreatedAt  time.Time
}

type WebhookSubscription struct {
	ID               uint   `gorm:"primaryKey"`
	UserID           uint   `gorm:"index;not null"`
	Name             string `gorm:"size:120;not null"`
	TargetURL        string `gorm:"size:512;not null"`
	Secret           string `gorm:"size:120;not null"`
	Events           string `gorm:"size:512;not null"`
	LastStatus       string `gorm:"size:32"`
	LastResponseCode int    `gorm:"not null;default:0"`
	LastError        string `gorm:"size:512"`
	LastDeliveredAt  *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type WebhookDelivery struct {
	ID           uint   `gorm:"primaryKey"`
	WebhookID    uint   `gorm:"index;not null"`
	Event        string `gorm:"size:120;index;not null"`
	Status       string `gorm:"size:32;not null"`
	ResponseCode int    `gorm:"not null;default:0"`
	ResponseBody string `gorm:"type:text"`
	RequestBody  string `gorm:"type:text"`
	CreatedAt    time.Time
}

type Conversation struct {
	ID            uint   `gorm:"primaryKey"`
	Kind          string `gorm:"size:32;index;not null"`
	Key           string `gorm:"size:160;uniqueIndex;not null"`
	WorkspaceID   *uint  `gorm:"index"`
	Title         string `gorm:"size:160"`
	Active        bool   `gorm:"index;not null;default:true"`
	BlockedReason string `gorm:"size:512"`
	BlockedBy     *uint  `gorm:"index"`
	BlockedAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type ConversationParticipant struct {
	ID             uint `gorm:"primaryKey"`
	ConversationID uint `gorm:"uniqueIndex:idx_conversation_user,priority:1;not null"`
	UserID         uint `gorm:"uniqueIndex:idx_conversation_user,priority:2;not null"`
	CreatedAt      time.Time
}

type Message struct {
	ID             uint   `gorm:"primaryKey"`
	ConversationID uint   `gorm:"index;not null"`
	SenderID       uint   `gorm:"index;not null"`
	Content        string `gorm:"type:text;not null"`
	CreatedAt      time.Time
}

type Workspace struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"size:160;index;not null"`
	Summary        string `gorm:"size:255;not null"`
	OwnerID        uint   `gorm:"index;not null"`
	ConversationID uint   `gorm:"index;not null"`
	Active         bool   `gorm:"index;not null;default:true"`
	BlockedReason  string `gorm:"size:512"`
	BlockedBy      *uint  `gorm:"index"`
	BlockedAt      *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type WorkspaceMember struct {
	ID          uint   `gorm:"primaryKey"`
	WorkspaceID uint   `gorm:"uniqueIndex:idx_workspace_user,priority:1;not null"`
	UserID      uint   `gorm:"uniqueIndex:idx_workspace_user,priority:2;not null"`
	Role        string `gorm:"size:32;not null;default:member"`
	CreatedAt   time.Time
}

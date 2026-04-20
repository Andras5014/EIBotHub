package app

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
)

type seedScenePageSpec struct {
	Slug        string
	Name        string
	Tagline     string
	Summary     string
	Description string
	SortOrder   int
}

type seedDocumentSpec struct {
	CategoryID uint
	Title      string
	Summary    string
	Content    string
	DocType    string
}

type seedFAQSpec struct {
	Question string
	Answer   string
}

type seedVideoSpec struct {
	Title     string
	Summary   string
	Link      string
	Category  string
	SortOrder int
}

type seedSkillSpec struct {
	Name        string
	Summary     string
	Description string
	Category    string
	Scene       string
	Guide       string
	ResourceRef string
	OwnerID     uint
	UsageCount  int64
}

type seedDiscussionSpec struct {
	Title    string
	Summary  string
	Content  string
	Category string
	UserID   uint
}

type seedAccessRequestSpec struct {
	DatasetID         uint
	UserID            uint
	Reason            string
	Status            string
	ReviewComment     string
	ApprovalStage     int
	RequiredApprovals int
	ApprovalExpiresAt *time.Time
	DownloadLimit     int
	DownloadCount     int
	ReviewedBy        *uint
	ReviewedAt        *time.Time
}

func ensureSeedScenePage(db *gorm.DB, spec seedScenePageSpec) (*model.ScenePageConfig, error) {
	var item model.ScenePageConfig
	err := db.Where("slug = ?", spec.Slug).First(&item).Error
	if err == nil {
		if err := db.Model(&item).Updates(map[string]any{
			"name":        spec.Name,
			"tagline":     spec.Tagline,
			"summary":     spec.Summary,
			"description": spec.Description,
			"sort_order":  spec.SortOrder,
			"enabled":     true,
		}).Error; err != nil {
			return nil, err
		}
		if err := db.Where("slug = ?", spec.Slug).First(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.ScenePageConfig{
		Slug:        spec.Slug,
		Name:        spec.Name,
		Tagline:     spec.Tagline,
		Summary:     spec.Summary,
		Description: spec.Description,
		SortOrder:   spec.SortOrder,
		Enabled:     true,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedDocument(db *gorm.DB, spec seedDocumentSpec) (*model.Document, error) {
	var item model.Document
	err := db.Where("title = ?", spec.Title).First(&item).Error
	if err == nil {
		if err := db.Model(&item).Updates(map[string]any{
			"category_id": spec.CategoryID,
			"summary":     spec.Summary,
			"content":     spec.Content,
			"doc_type":    spec.DocType,
			"status":      model.StatusPublished,
		}).Error; err != nil {
			return nil, err
		}
		if err := db.Where("title = ?", spec.Title).First(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.Document{
		CategoryID: spec.CategoryID,
		Title:      spec.Title,
		Summary:    spec.Summary,
		Content:    spec.Content,
		DocType:    spec.DocType,
		Status:     model.StatusPublished,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedFAQ(db *gorm.DB, spec seedFAQSpec) (*model.FAQ, error) {
	var item model.FAQ
	err := db.Where("question = ?", spec.Question).First(&item).Error
	if err == nil {
		if err := db.Model(&item).Update("answer", spec.Answer).Error; err != nil {
			return nil, err
		}
		if err := db.Where("question = ?", spec.Question).First(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.FAQ{
		Question: spec.Question,
		Answer:   spec.Answer,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedVideo(db *gorm.DB, spec seedVideoSpec) (*model.VideoTutorial, error) {
	var item model.VideoTutorial
	err := db.Where("title = ?", spec.Title).First(&item).Error
	if err == nil {
		if err := db.Model(&item).Updates(map[string]any{
			"summary":    spec.Summary,
			"link":       spec.Link,
			"category":   spec.Category,
			"sort_order": spec.SortOrder,
			"active":     true,
		}).Error; err != nil {
			return nil, err
		}
		if err := db.Where("title = ?", spec.Title).First(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.VideoTutorial{
		Title:     spec.Title,
		Summary:   spec.Summary,
		Link:      spec.Link,
		Category:  spec.Category,
		SortOrder: spec.SortOrder,
		Active:    true,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedSkill(db *gorm.DB, spec seedSkillSpec) (*model.Skill, error) {
	var item model.Skill
	err := db.Where("name = ?", spec.Name).First(&item).Error
	if err == nil {
		if err := db.Model(&item).Updates(map[string]any{
			"summary":      spec.Summary,
			"description":  spec.Description,
			"category":     spec.Category,
			"scene":        spec.Scene,
			"guide":        spec.Guide,
			"resource_ref": spec.ResourceRef,
			"owner_id":     spec.OwnerID,
			"usage_count":  spec.UsageCount,
			"status":       model.StatusPublished,
		}).Error; err != nil {
			return nil, err
		}
		if err := db.Where("name = ?", spec.Name).First(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.Skill{
		Name:        spec.Name,
		Summary:     spec.Summary,
		Description: spec.Description,
		Category:    spec.Category,
		Scene:       spec.Scene,
		Guide:       spec.Guide,
		ResourceRef: spec.ResourceRef,
		Status:      model.StatusPublished,
		OwnerID:     spec.OwnerID,
		UsageCount:  spec.UsageCount,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedDiscussion(db *gorm.DB, spec seedDiscussionSpec) (*model.Discussion, error) {
	var item model.Discussion
	err := db.Where("title = ?", spec.Title).First(&item).Error
	if err == nil {
		if err := db.Model(&item).Updates(map[string]any{
			"summary":  spec.Summary,
			"content":  spec.Content,
			"category": spec.Category,
			"user_id":  spec.UserID,
		}).Error; err != nil {
			return nil, err
		}
		if err := db.Where("title = ?", spec.Title).First(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item = model.Discussion{
		Title:    spec.Title,
		Summary:  spec.Summary,
		Content:  spec.Content,
		Category: spec.Category,
		UserID:   spec.UserID,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedComment(db *gorm.DB, resourceType string, resourceID, userID uint, content string) error {
	var item model.ResourceComment
	err := db.Where("resource_type = ? AND resource_id = ? AND user_id = ? AND content = ?", resourceType, resourceID, userID, content).First(&item).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(&model.ResourceComment{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		UserID:       userID,
		Content:      content,
	}).Error
}

func ensureSeedFollow(db *gorm.DB, followerID, followedUserID uint) error {
	if followerID == 0 || followedUserID == 0 || followerID == followedUserID {
		return nil
	}
	return db.Where("follower_id = ? AND followed_user_id = ?", followerID, followedUserID).
		FirstOrCreate(&model.Follow{FollowerID: followerID, FollowedUserID: followedUserID}).Error
}

func ensureSeedFavorite(db *gorm.DB, userID uint, resourceType string, resourceID uint, resourceTitle string) error {
	var item model.Favorite
	err := db.Where("user_id = ? AND resource_type = ? AND resource_id = ?", userID, resourceType, resourceID).First(&item).Error
	if err == nil {
		if item.ResourceTitle != resourceTitle {
			return db.Model(&item).Update("resource_title", resourceTitle).Error
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(&model.Favorite{
		UserID:        userID,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		ResourceTitle: resourceTitle,
	}).Error
}

func ensureSeedDownload(db *gorm.DB, userID uint, resourceType string, resourceID uint, resourceTitle string) error {
	var item model.DownloadRecord
	err := db.Where("user_id = ? AND resource_type = ? AND resource_id = ?", userID, resourceType, resourceID).First(&item).Error
	if err == nil {
		if item.ResourceTitle != resourceTitle {
			return db.Model(&item).Update("resource_title", resourceTitle).Error
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(&model.DownloadRecord{
		UserID:        userID,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		ResourceTitle: resourceTitle,
	}).Error
}

func ensureSeedNotification(db *gorm.DB, userID uint, notifyType, title, content string) error {
	var item model.Notification
	err := db.Where("user_id = ? AND type = ? AND title = ?", userID, notifyType, title).First(&item).Error
	if err == nil {
		if err := db.Model(&item).Updates(map[string]any{
			"content": content,
			"read":    false,
		}).Error; err != nil {
			return err
		}
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(&model.Notification{
		UserID:  userID,
		Type:    notifyType,
		Title:   title,
		Content: content,
	}).Error
}

func ensureSeedSearchKeyword(db *gorm.DB, query, keywordType string, sortOrder int) error {
	var item model.SearchKeywordConfig
	err := db.Where("query = ? AND keyword_type = ?", query, keywordType).First(&item).Error
	if err == nil {
		return db.Model(&item).Updates(map[string]any{
			"sort_order": sortOrder,
			"enabled":    true,
		}).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(&model.SearchKeywordConfig{
		Query:       query,
		KeywordType: keywordType,
		SortOrder:   sortOrder,
		Enabled:     true,
	}).Error
}

func ensureSeedSearchRecord(db *gorm.DB, query, searchType string) error {
	var item model.SearchRecord
	err := db.Where("query = ? AND search_type = ?", query, searchType).First(&item).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	return db.Create(&model.SearchRecord{
		Query:      query,
		SearchType: searchType,
	}).Error
}

func ensureSeedRewardBenefit(db *gorm.DB, name, summary string, costPoints int) (*model.RewardBenefit, error) {
	var item model.RewardBenefit
	err := db.Where("name = ?", name).First(&item).Error
	if err == nil {
		if err := db.Model(&item).Updates(map[string]any{
			"summary":     summary,
			"cost_points": costPoints,
			"active":      true,
		}).Error; err != nil {
			return nil, err
		}
		if err := db.Where("name = ?", name).First(&item).Error; err != nil {
			return nil, err
		}
		return &item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	item = model.RewardBenefit{
		Name:       name,
		Summary:    summary,
		CostPoints: costPoints,
		Active:     true,
	}
	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func ensureSeedRewardRedemption(db *gorm.DB, userID uint, benefit model.RewardBenefit) error {
	var item model.RewardRedemption
	err := db.Where("user_id = ? AND benefit_id = ?", userID, benefit.ID).First(&item).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if err := db.Create(&model.RewardRedemption{
		UserID:     userID,
		BenefitID:  benefit.ID,
		CostPoints: benefit.CostPoints,
	}).Error; err != nil {
		return err
	}
	return ensureSeedReward(db, seedRewardSpec{
		UserID:     userID,
		SourceType: "redeem",
		Points:     -benefit.CostPoints,
		Remark:     "兑换权益：" + benefit.Name,
	})
}

func ensureSeedAccessRequest(db *gorm.DB, spec seedAccessRequestSpec) error {
	var item model.DatasetAccessRequest
	err := db.Where("dataset_id = ? AND user_id = ? AND reason = ?", spec.DatasetID, spec.UserID, spec.Reason).First(&item).Error
	if err == nil {
		return db.Model(&item).Updates(map[string]any{
			"status":              spec.Status,
			"review_comment":      spec.ReviewComment,
			"approval_stage":      spec.ApprovalStage,
			"required_approvals":  spec.RequiredApprovals,
			"approval_expires_at": spec.ApprovalExpiresAt,
			"download_limit":      spec.DownloadLimit,
			"download_count":      spec.DownloadCount,
			"reviewed_by":         spec.ReviewedBy,
			"reviewed_at":         spec.ReviewedAt,
		}).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	item = model.DatasetAccessRequest{
		DatasetID:         spec.DatasetID,
		UserID:            spec.UserID,
		Reason:            spec.Reason,
		Status:            spec.Status,
		ReviewComment:     spec.ReviewComment,
		ApprovalStage:     spec.ApprovalStage,
		RequiredApprovals: spec.RequiredApprovals,
		ApprovalExpiresAt: spec.ApprovalExpiresAt,
		DownloadLimit:     spec.DownloadLimit,
		DownloadCount:     spec.DownloadCount,
		ReviewedBy:        spec.ReviewedBy,
		ReviewedAt:        spec.ReviewedAt,
	}
	return db.Create(&item).Error
}

func ensureDirectConversation(db *gorm.DB, userA, userB uint, title, content string, senderID uint) (*model.Conversation, error) {
	key := repository.DirectKey(userA, userB)
	var item model.Conversation
	err := db.Where("key = ?", key).First(&item).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		item = model.Conversation{
			Kind:  "direct",
			Key:   key,
			Title: title,
		}
		if err := db.Create(&item).Error; err != nil {
			return nil, err
		}
	}

	if item.Title != title {
		if err := db.Model(&item).Update("title", title).Error; err != nil {
			return nil, err
		}
	}
	for _, userID := range []uint{userA, userB} {
		if err := db.Where("conversation_id = ? AND user_id = ?", item.ID, userID).
			FirstOrCreate(&model.ConversationParticipant{ConversationID: item.ID, UserID: userID}).Error; err != nil {
			return nil, err
		}
	}

	var messageCount int64
	if err := db.Model(&model.Message{}).Where("conversation_id = ?", item.ID).Count(&messageCount).Error; err != nil {
		return nil, err
	}
	if messageCount == 0 {
		if err := db.Create(&model.Message{
			ConversationID: item.ID,
			SenderID:       senderID,
			Content:        content,
		}).Error; err != nil {
			return nil, err
		}
	}
	return &item, nil
}

func ensureWorkspaceSeed(db *gorm.DB, name, summary string, ownerID uint, memberIDs []uint, senderID uint, content string) (*model.Workspace, error) {
	var item model.Workspace
	err := db.Where("name = ?", name).First(&item).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		conversation := model.Conversation{
			Kind:  "workspace",
			Key:   "workspace:seed:" + time.Now().Format("20060102150405.000000000"),
			Title: name,
		}
		if err := db.Create(&conversation).Error; err != nil {
			return nil, err
		}
		item = model.Workspace{
			Name:           name,
			Summary:        summary,
			OwnerID:        ownerID,
			ConversationID: conversation.ID,
		}
		if err := db.Create(&item).Error; err != nil {
			return nil, err
		}
		conversation.Key = repository.WorkspaceKey(item.ID)
		conversation.WorkspaceID = &item.ID
		if err := db.Save(&conversation).Error; err != nil {
			return nil, err
		}
	} else {
		if err := db.Model(&item).Updates(map[string]any{
			"summary":  summary,
			"owner_id": ownerID,
			"active":   true,
		}).Error; err != nil {
			return nil, err
		}
	}

	memberRole := map[uint]string{ownerID: "owner"}
	for _, memberID := range memberIDs {
		if memberID != ownerID {
			memberRole[memberID] = "member"
		}
	}
	for userID, role := range memberRole {
		if err := db.Where("workspace_id = ? AND user_id = ?", item.ID, userID).
			FirstOrCreate(&model.WorkspaceMember{WorkspaceID: item.ID, UserID: userID, Role: role}).Error; err != nil {
			return nil, err
		}
		if err := db.Where("conversation_id = ? AND user_id = ?", item.ConversationID, userID).
			FirstOrCreate(&model.ConversationParticipant{ConversationID: item.ConversationID, UserID: userID}).Error; err != nil {
			return nil, err
		}
	}

	var messageCount int64
	if err := db.Model(&model.Message{}).Where("conversation_id = ?", item.ConversationID).Count(&messageCount).Error; err != nil {
		return nil, err
	}
	if messageCount == 0 {
		if err := db.Create(&model.Message{
			ConversationID: item.ConversationID,
			SenderID:       senderID,
			Content:        content,
		}).Error; err != nil {
			return nil, err
		}
	}
	return &item, nil
}

package service

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type CommunityService struct {
	community *repository.CommunityRepository
	users     *repository.UserRepository
	models    *repository.ModelRepository
	datasets  *repository.DatasetRepository
	content   *repository.ContentRepository
	activity  *repository.UserActivityRepository
	rewards   *RewardService
	hooks     *IntegrationService
}

func NewCommunityService(community *repository.CommunityRepository, users *repository.UserRepository, models *repository.ModelRepository, datasets *repository.DatasetRepository, content *repository.ContentRepository, activity *repository.UserActivityRepository, rewards *RewardService, hooks *IntegrationService) *CommunityService {
	return &CommunityService{
		community: community,
		users:     users,
		models:    models,
		datasets:  datasets,
		content:   content,
		activity:  activity,
		rewards:   rewards,
		hooks:     hooks,
	}
}

func (s *CommunityService) CreateModelEvaluation(userID, modelID uint, input dto.ModelEvaluationRequest) (*dto.ModelEvaluationItem, error) {
	if _, err := s.models.GetByID(modelID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "model_not_found", "模型不存在")
		}
		return nil, err
	}

	item := &model.ModelEvaluation{
		ModelID:   modelID,
		UserID:    userID,
		Benchmark: input.Benchmark,
		Summary:   input.Summary,
		Score:     input.Score,
		Notes:     input.Notes,
	}
	if err := s.community.AddModelEvaluation(item); err != nil {
		return nil, err
	}
	user, _ := s.users.FindByID(userID)
	s.rewards.Add(userID, "model_evaluation", 8, "提交模型评测")
	if modelItem, err := s.models.GetByID(modelID); err == nil {
		s.notifyUser(modelItem.OwnerID, userID, "model_evaluation", "模型收到新评测", user.Username+" 提交了新的模型评测")
	}
	result := &dto.ModelEvaluationItem{
		ID:        item.ID,
		Benchmark: item.Benchmark,
		Summary:   item.Summary,
		Score:     item.Score,
		Notes:     item.Notes,
		UserName:  user.Username,
		CreatedAt: item.CreatedAt,
	}
	return result, nil
}

func (s *CommunityService) ListModelEvaluations(modelID uint) ([]dto.ModelEvaluationItem, error) {
	items, err := s.community.ListModelEvaluations(modelID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ModelEvaluationItem, 0, len(items))
	for _, item := range items {
		user, _ := s.users.FindByID(item.UserID)
		result = append(result, dto.ModelEvaluationItem{
			ID:        item.ID,
			Benchmark: item.Benchmark,
			Summary:   item.Summary,
			Score:     item.Score,
			Notes:     item.Notes,
			UserName:  user.Username,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) UpsertRating(resourceType string, resourceID, userID uint, input dto.RatingRequest) (*dto.RatingSummary, error) {
	if err := s.community.UpsertRating(&model.ResourceRating{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		UserID:       userID,
		Score:        input.Score,
		Feedback:     input.Feedback,
	}); err != nil {
		return nil, err
	}
	if ownerID, ok := s.resourceOwner(resourceType, resourceID); ok {
		user, _ := s.users.FindByID(userID)
		s.notifyUser(ownerID, userID, "rating", "资源收到新评分", user.Username+" 提交了新的评分")
	}
	return s.RatingSummary(resourceType, resourceID)
}

func (s *CommunityService) RatingSummary(resourceType string, resourceID uint) (*dto.RatingSummary, error) {
	items, err := s.community.ListRatings(resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	result := &dto.RatingSummary{
		Count: int64(len(items)),
		Items: make([]dto.RatingItem, 0, len(items)),
	}
	var total int64
	for _, item := range items {
		user, _ := s.users.FindByID(item.UserID)
		total += int64(item.Score)
		result.Items = append(result.Items, dto.RatingItem{
			ID:        item.ID,
			Score:     item.Score,
			Feedback:  item.Feedback,
			UserName:  user.Username,
			CreatedAt: item.CreatedAt,
		})
	}
	if result.Count > 0 {
		result.Average = float64(total) / float64(result.Count)
	}
	return result, nil
}

func (s *CommunityService) AddComment(resourceType string, resourceID, userID uint, input dto.CommentRequest) (*dto.CommentItem, error) {
	item := &model.ResourceComment{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		UserID:       userID,
		ParentID:     input.ParentID,
		Content:      input.Content,
	}
	if err := s.community.AddComment(item); err != nil {
		return nil, err
	}
	user, _ := s.users.FindByID(userID)
	s.rewards.Add(userID, "comment", 2, "发表社区评论")
	if ownerID, ok := s.resourceOwner(resourceType, resourceID); ok {
		s.notifyUser(ownerID, userID, "comment", "资源收到新评论", user.Username+" 发表了新的评论")
	}
	return &dto.CommentItem{
		ID:        item.ID,
		ParentID:  item.ParentID,
		Content:   item.Content,
		UserID:    item.UserID,
		UserName:  user.Username,
		CreatedAt: item.CreatedAt,
	}, nil
}

func (s *CommunityService) ListComments(resourceType string, resourceID uint) ([]dto.CommentItem, error) {
	items, err := s.community.ListComments(resourceType, resourceID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.CommentItem, 0, len(items))
	for _, item := range items {
		user, _ := s.users.FindByID(item.UserID)
		result = append(result, dto.CommentItem{
			ID:        item.ID,
			ParentID:  item.ParentID,
			Content:   item.Content,
			UserID:    item.UserID,
			UserName:  user.Username,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) DatasetStats(datasetID uint) (*dto.DatasetStatsResponse, error) {
	dataset, err := s.datasets.GetByID(datasetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "dataset_not_found", "数据集不存在")
		}
		return nil, err
	}
	mix, err := s.community.DatasetSampleBreakdown(datasetID)
	if err != nil {
		return nil, err
	}
	trend, err := s.community.DatasetDownloadTrend(datasetID, 7)
	if err != nil {
		return nil, err
	}
	result := &dto.DatasetStatsResponse{
		DatasetID:     datasetID,
		DownloadCount: dataset.Downloads,
		SampleCount:   dataset.SampleCount,
		SampleTypeMix: make([]dto.StatValueItem, 0, len(mix)),
		DownloadTrend: make([]dto.TrendValueItem, 0, len(trend)),
	}
	for _, item := range mix {
		result.SampleTypeMix = append(result.SampleTypeMix, dto.StatValueItem{Label: item.Label, Value: item.Value})
	}
	for _, item := range trend {
		result.DownloadTrend = append(result.DownloadTrend, dto.TrendValueItem{Label: item.Label, Value: item.Value})
	}
	return result, nil
}

func (s *CommunityService) CreateSkill(userID uint, input dto.SkillCreateRequest) (*dto.SkillItem, error) {
	item := &model.Skill{
		Name:        input.Name,
		Summary:     input.Summary,
		Description: input.Description,
		Category:    input.Category,
		Scene:       input.Scene,
		Guide:       input.Guide,
		ResourceRef: input.ResourceRef,
		Status:      model.StatusPublished,
		OwnerID:     userID,
	}
	if err := s.community.CreateSkill(item); err != nil {
		return nil, err
	}
	s.rewards.Add(userID, "skill_create", 30, "发布技能")
	result, err := s.GetSkill(item.ID)
	if err != nil {
		return nil, err
	}
	if s.hooks != nil {
		s.hooks.Emit(userID, WebhookEventSkillCreated, map[string]any{
			"skill_id":     result.ID,
			"name":         result.Name,
			"category":     result.Category,
			"scene":        result.Scene,
			"owner_id":     result.OwnerID,
			"owner_name":   result.OwnerName,
			"resource_ref": result.ResourceRef,
		})
	}
	return result, nil
}

func (s *CommunityService) ForkSkill(skillID, userID uint) (*dto.SkillItem, error) {
	source, err := s.community.GetSkill(skillID)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "skill_not_found", "技能不存在")
	}
	item := &model.Skill{
		Name:         source.Name + " Fork",
		Summary:      source.Summary,
		Description:  source.Description,
		Category:     source.Category,
		Scene:        source.Scene,
		Guide:        source.Guide,
		ResourceRef:  source.ResourceRef,
		Status:       model.StatusPublished,
		ForkedFromID: &source.ID,
		OwnerID:      userID,
	}
	if err := s.community.CreateSkill(item); err != nil {
		return nil, err
	}
	return s.GetSkill(item.ID)
}

func (s *CommunityService) ListSkills() ([]dto.SkillItem, error) {
	items, err := s.community.ListSkills()
	if err != nil {
		return nil, err
	}
	result := make([]dto.SkillItem, 0, len(items))
	for _, item := range items {
		summary, err := s.RatingSummary("skill", item.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, toSkillItem(item, *summary))
	}
	return result, nil
}

func (s *CommunityService) GetSkill(id uint) (*dto.SkillItem, error) {
	item, err := s.community.GetSkill(id)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "skill_not_found", "技能不存在")
	}
	summary, err := s.RatingSummary("skill", id)
	if err != nil {
		return nil, err
	}
	result := toSkillItem(*item, *summary)
	return &result, nil
}

func (s *CommunityService) CreateDiscussion(userID uint, input dto.DiscussionRequest) (*dto.DiscussionItem, error) {
	item := &model.Discussion{
		Title:    input.Title,
		Summary:  input.Summary,
		Content:  input.Content,
		Category: input.Category,
		UserID:   userID,
	}
	if err := s.community.CreateDiscussion(item); err != nil {
		return nil, err
	}
	s.rewards.Add(userID, "discussion_create", 10, "发起讨论")
	result, err := s.GetDiscussion(item.ID)
	if err != nil {
		return nil, err
	}
	if s.hooks != nil {
		s.hooks.Emit(userID, WebhookEventDiscussionCreated, map[string]any{
			"discussion_id": result.ID,
			"title":         result.Title,
			"category":      result.Category,
			"user_id":       result.UserID,
			"user_name":     result.UserName,
		})
	}
	return result, nil
}

func (s *CommunityService) ListDiscussions() ([]dto.DiscussionItem, error) {
	items, err := s.community.ListDiscussions()
	if err != nil {
		return nil, err
	}
	result := make([]dto.DiscussionItem, 0, len(items))
	for _, item := range items {
		count, err := s.community.CountComments("discussion", item.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, dto.DiscussionItem{
			ID:           item.ID,
			Title:        item.Title,
			Summary:      item.Summary,
			Content:      item.Content,
			Category:     item.Category,
			UserID:       item.UserID,
			UserName:     item.User.Username,
			CommentCount: count,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) GetDiscussion(id uint) (*dto.DiscussionItem, error) {
	item, err := s.community.GetDiscussion(id)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "discussion_not_found", "讨论不存在")
	}
	count, err := s.community.CountComments("discussion", id)
	if err != nil {
		return nil, err
	}
	result := &dto.DiscussionItem{
		ID:           item.ID,
		Title:        item.Title,
		Summary:      item.Summary,
		Content:      item.Content,
		Category:     item.Category,
		UserID:       item.UserID,
		UserName:     item.User.Username,
		CommentCount: count,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
	return result, nil
}

func (s *CommunityService) ToggleFollow(followerID, followedUserID uint) error {
	if followerID == followedUserID {
		return support.NewError(http.StatusBadRequest, "invalid_follow", "不能关注自己")
	}
	if _, err := s.users.FindByID(followedUserID); err != nil {
		return support.NewError(http.StatusNotFound, "user_not_found", "用户不存在")
	}
	alreadyFollowing, err := s.community.IsFollowing(followerID, followedUserID)
	if err != nil {
		return err
	}
	if err := s.community.UpsertFollow(&model.Follow{
		FollowerID:     followerID,
		FollowedUserID: followedUserID,
	}); err != nil {
		return err
	}
	if !alreadyFollowing {
		follower, _ := s.users.FindByID(followerID)
		s.notifyUser(followedUserID, followerID, "follow", "你有新的关注者", follower.Username+" 关注了你")
	}
	return nil
}

func (s *CommunityService) MyFollows(followerID uint) ([]dto.FollowItem, error) {
	items, err := s.community.ListFollows(followerID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.FollowItem, 0, len(items))
	for _, item := range items {
		user, err := s.users.FindByID(item.FollowedUserID)
		if err != nil {
			return nil, err
		}
		result = append(result, dto.FollowItem{
			ID:        item.ID,
			UserID:    user.ID,
			UserName:  user.Username,
			Bio:       user.Bio,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) MyFollowers(userID uint) ([]dto.FollowItem, error) {
	items, err := s.community.ListFollowers(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.FollowItem, 0, len(items))
	for _, item := range items {
		user, err := s.users.FindByID(item.FollowerID)
		if err != nil {
			return nil, err
		}
		result = append(result, dto.FollowItem{
			ID:        item.ID,
			UserID:    user.ID,
			UserName:  user.Username,
			Bio:       user.Bio,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) FollowStats(userID uint) (*dto.FollowStats, error) {
	follows, err := s.community.CountFollows(userID)
	if err != nil {
		return nil, err
	}
	followers, err := s.community.CountFollowers(userID)
	if err != nil {
		return nil, err
	}
	return &dto.FollowStats{
		Follows:   follows,
		Followers: followers,
	}, nil
}

func (s *CommunityService) HotQueries() ([]dto.SearchHotItem, error) {
	items, err := s.community.HotQueries(8)
	if err != nil {
		return nil, err
	}
	result := make([]dto.SearchHotItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.SearchHotItem{Query: item.Query, Count: item.Count})
	}
	return result, nil
}

func (s *CommunityService) UserContributions(userID uint) (*dto.UserContributionPayload, error) {
	skills, err := s.community.ListSkillsByOwner(userID)
	if err != nil {
		return nil, err
	}
	discussions, err := s.community.ListDiscussionsByUser(userID)
	if err != nil {
		return nil, err
	}

	result := &dto.UserContributionPayload{
		Skills:      make([]dto.SkillItem, 0, len(skills)),
		Discussions: make([]dto.DiscussionItem, 0, len(discussions)),
	}
	for _, item := range skills {
		summary, err := s.RatingSummary("skill", item.ID)
		if err != nil {
			return nil, err
		}
		result.Skills = append(result.Skills, toSkillItem(item, *summary))
	}
	for _, item := range discussions {
		count, err := s.community.CountComments("discussion", item.ID)
		if err != nil {
			return nil, err
		}
		result.Discussions = append(result.Discussions, dto.DiscussionItem{
			ID:           item.ID,
			Title:        item.Title,
			Summary:      item.Summary,
			Content:      item.Content,
			Category:     item.Category,
			UserID:       item.UserID,
			UserName:     item.User.Username,
			CommentCount: count,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) AdminOverview() (*dto.AdminCommunityOverview, error) {
	skills, err := s.community.CountSkills()
	if err != nil {
		return nil, err
	}
	discussions, err := s.community.CountDiscussions()
	if err != nil {
		return nil, err
	}
	comments, err := s.community.CountAllComments()
	if err != nil {
		return nil, err
	}
	return &dto.AdminCommunityOverview{
		Skills:      skills,
		Discussions: discussions,
		Comments:    comments,
	}, nil
}

func (s *CommunityService) AdminSkills() ([]dto.AdminSkillModerationItem, error) {
	items, err := s.community.ListAllSkills()
	if err != nil {
		return nil, err
	}
	result := make([]dto.AdminSkillModerationItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.AdminSkillModerationItem{
			ID:        item.ID,
			Name:      item.Name,
			Summary:   item.Summary,
			Status:    item.Status,
			OwnerID:   item.OwnerID,
			OwnerName: item.Owner.Username,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) AdminHideSkill(id uint) error {
	return s.community.UpdateSkillStatus(id, model.StatusRejected)
}

func (s *CommunityService) AdminDiscussions() ([]dto.AdminDiscussionModerationItem, error) {
	items, err := s.community.ListDiscussions()
	if err != nil {
		return nil, err
	}
	result := make([]dto.AdminDiscussionModerationItem, 0, len(items))
	for _, item := range items {
		count, err := s.community.CountComments("discussion", item.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, dto.AdminDiscussionModerationItem{
			ID:           item.ID,
			Title:        item.Title,
			Category:     item.Category,
			UserID:       item.UserID,
			UserName:     item.User.Username,
			CommentCount: count,
			UpdatedAt:    item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) AdminDeleteDiscussion(id uint) error {
	return s.community.DeleteDiscussion(id)
}

func (s *CommunityService) AdminComments() ([]dto.AdminCommentModerationItem, error) {
	items, err := s.community.ListAllComments()
	if err != nil {
		return nil, err
	}
	result := make([]dto.AdminCommentModerationItem, 0, len(items))
	for _, item := range items {
		user, _ := s.users.FindByID(item.UserID)
		result = append(result, dto.AdminCommentModerationItem{
			ID:           item.ID,
			ResourceType: item.ResourceType,
			ResourceID:   item.ResourceID,
			UserID:       item.UserID,
			UserName:     user.Username,
			Content:      item.Content,
			CreatedAt:    item.CreatedAt,
		})
	}
	return result, nil
}

func (s *CommunityService) AdminDeleteComment(id uint) error {
	return s.community.DeleteComment(id)
}

func toSkillItem(item model.Skill, rating dto.RatingSummary) dto.SkillItem {
	return dto.SkillItem{
		ID:           item.ID,
		Name:         item.Name,
		Summary:      item.Summary,
		Description:  item.Description,
		Category:     item.Category,
		Scene:        item.Scene,
		Guide:        item.Guide,
		ResourceRef:  support.SplitCSV(item.ResourceRef),
		Status:       item.Status,
		ForkedFromID: item.ForkedFromID,
		UsageCount:   item.UsageCount,
		OwnerID:      item.OwnerID,
		OwnerName:    item.Owner.Username,
		Rating:       rating,
		UpdatedAt:    item.UpdatedAt,
	}
}

func (s *CommunityService) notifyUser(targetUserID, actorUserID uint, notificationType, title, content string) {
	if targetUserID == 0 || targetUserID == actorUserID || s.activity == nil {
		return
	}
	_ = s.activity.AddNotification(&model.Notification{
		UserID:  targetUserID,
		Type:    notificationType,
		Title:   title,
		Content: content,
	})
}

func (s *CommunityService) resourceOwner(resourceType string, resourceID uint) (uint, bool) {
	switch resourceType {
	case "model":
		item, err := s.models.GetByID(resourceID)
		if err == nil {
			return item.OwnerID, true
		}
	case "dataset":
		item, err := s.datasets.GetByID(resourceID)
		if err == nil {
			return item.OwnerID, true
		}
	case "skill":
		item, err := s.community.GetSkill(resourceID)
		if err == nil {
			return item.OwnerID, true
		}
	case "discussion":
		item, err := s.community.GetDiscussion(resourceID)
		if err == nil {
			return item.UserID, true
		}
	}
	return 0, false
}

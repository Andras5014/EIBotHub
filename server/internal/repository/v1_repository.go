package repository

import (
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type HotQuery struct {
	Query string
	Count int64
}

type AggregationValue struct {
	Label string
	Value int64
}

type CommunityRepository struct {
	db *gorm.DB
}

func NewCommunityRepository(db *gorm.DB) *CommunityRepository {
	return &CommunityRepository{db: db}
}

func (r *CommunityRepository) AddModelEvaluation(item *model.ModelEvaluation) error {
	return r.db.Create(item).Error
}

func (r *CommunityRepository) ListModelEvaluations(modelID uint) ([]model.ModelEvaluation, error) {
	var items []model.ModelEvaluation
	err := r.db.Where("model_id = ?", modelID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) UpsertRating(item *model.ResourceRating) error {
	var existing model.ResourceRating
	err := r.db.Where("resource_type = ? AND resource_id = ? AND user_id = ?", item.ResourceType, item.ResourceID, item.UserID).First(&existing).Error
	if err == nil {
		existing.Score = item.Score
		existing.Feedback = item.Feedback
		return r.db.Save(&existing).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(item).Error
}

func (r *CommunityRepository) ListRatings(resourceType string, resourceID uint) ([]model.ResourceRating, error) {
	var items []model.ResourceRating
	err := r.db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) AddComment(item *model.ResourceComment) error {
	return r.db.Create(item).Error
}

func (r *CommunityRepository) ListComments(resourceType string, resourceID uint) ([]model.ResourceComment, error) {
	var items []model.ResourceComment
	err := r.db.Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Order("created_at asc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) CreateSkill(item *model.Skill) error {
	return r.db.Create(item).Error
}

func (r *CommunityRepository) GetSkill(id uint) (*model.Skill, error) {
	var item model.Skill
	err := r.db.Preload("Owner").First(&item, id).Error
	return &item, err
}

func (r *CommunityRepository) ListSkills() ([]model.Skill, error) {
	var items []model.Skill
	err := r.db.Preload("Owner").Where("status = ?", model.StatusPublished).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) SearchSkills(query string, page, pageSize int) ([]model.Skill, int64, error) {
	db := r.db.Model(&model.Skill{}).Where("status = ?", model.StatusPublished)
	db = applyTextSearch(db, query, []string{"name", "summary", "description", "category", "scene"})

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Skill
	err := applyPagination(db.Order("usage_count desc, updated_at desc"), page, pageSize).
		Preload("Owner").
		Find(&items).Error
	return items, total, err
}

func (r *CommunityRepository) ListAllSkills() ([]model.Skill, error) {
	var items []model.Skill
	err := r.db.Preload("Owner").Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) ListSkillsByOwner(ownerID uint) ([]model.Skill, error) {
	var items []model.Skill
	err := r.db.Preload("Owner").Where("owner_id = ?", ownerID).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) UpdateSkillStatus(id uint, status string) error {
	return r.db.Model(&model.Skill{}).Where("id = ?", id).Update("status", status).Error
}

func (r *CommunityRepository) CreateDiscussion(item *model.Discussion) error {
	return r.db.Create(item).Error
}

func (r *CommunityRepository) GetDiscussion(id uint) (*model.Discussion, error) {
	var item model.Discussion
	err := r.db.Preload("User").First(&item, id).Error
	return &item, err
}

func (r *CommunityRepository) ListDiscussions() ([]model.Discussion, error) {
	var items []model.Discussion
	err := r.db.Preload("User").Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) SearchDiscussions(query string, page, pageSize int) ([]model.Discussion, int64, error) {
	db := r.db.Model(&model.Discussion{})
	db = applyTextSearch(db, query, []string{"title", "summary", "content", "category"})

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Discussion
	err := applyPagination(db.Order("updated_at desc"), page, pageSize).
		Preload("User").
		Find(&items).Error
	return items, total, err
}

func (r *CommunityRepository) ListDiscussionsByUser(userID uint) ([]model.Discussion, error) {
	var items []model.Discussion
	err := r.db.Preload("User").Where("user_id = ?", userID).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) DeleteDiscussion(id uint) error {
	return r.db.Delete(&model.Discussion{}, id).Error
}

func (r *CommunityRepository) UpsertFollow(item *model.Follow) error {
	var existing model.Follow
	err := r.db.Where("follower_id = ? AND followed_user_id = ?", item.FollowerID, item.FollowedUserID).First(&existing).Error
	if err == nil {
		return r.db.Delete(&existing).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(item).Error
}

func (r *CommunityRepository) IsFollowing(followerID, followedUserID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Follow{}).Where("follower_id = ? AND followed_user_id = ?", followerID, followedUserID).Count(&count).Error
	return count > 0, err
}

func (r *CommunityRepository) ListFollows(followerID uint) ([]model.Follow, error) {
	var items []model.Follow
	err := r.db.Where("follower_id = ?", followerID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) ListFollowers(followedUserID uint) ([]model.Follow, error) {
	var items []model.Follow
	err := r.db.Where("followed_user_id = ?", followedUserID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) ListAllComments() ([]model.ResourceComment, error) {
	var items []model.ResourceComment
	err := r.db.Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *CommunityRepository) DeleteComment(id uint) error {
	return r.db.Delete(&model.ResourceComment{}, id).Error
}

func (r *CommunityRepository) AddSearchRecord(item *model.SearchRecord) error {
	return r.db.Create(item).Error
}

func (r *CommunityRepository) HotQueries(limit int) ([]HotQuery, error) {
	var rows []HotQuery
	err := r.db.Model(&model.SearchRecord{}).
		Select("query, count(*) as count").
		Group("query").
		Order("count desc, max(created_at) desc").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

func (r *CommunityRepository) CountComments(resourceType string, resourceID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ResourceComment{}).Where("resource_type = ? AND resource_id = ?", resourceType, resourceID).Count(&count).Error
	return count, err
}

func (r *CommunityRepository) CountSkills() (int64, error) {
	var count int64
	return count, r.db.Model(&model.Skill{}).Count(&count).Error
}

func (r *CommunityRepository) CountFollows(followerID uint) (int64, error) {
	var count int64
	return count, r.db.Model(&model.Follow{}).Where("follower_id = ?", followerID).Count(&count).Error
}

func (r *CommunityRepository) CountFollowers(followedUserID uint) (int64, error) {
	var count int64
	return count, r.db.Model(&model.Follow{}).Where("followed_user_id = ?", followedUserID).Count(&count).Error
}

func (r *CommunityRepository) CountDiscussions() (int64, error) {
	var count int64
	return count, r.db.Model(&model.Discussion{}).Count(&count).Error
}

func (r *CommunityRepository) CountAllComments() (int64, error) {
	var count int64
	return count, r.db.Model(&model.ResourceComment{}).Count(&count).Error
}

func (r *CommunityRepository) DatasetSampleBreakdown(datasetID uint) ([]AggregationValue, error) {
	var rows []AggregationValue
	err := r.db.Model(&model.DatasetSample{}).
		Select("sample_type as label, count(*) as value").
		Where("dataset_id = ?", datasetID).
		Group("sample_type").
		Order("value desc").
		Scan(&rows).Error
	return rows, err
}

func (r *CommunityRepository) DatasetDownloadTrend(datasetID uint, days int) ([]AggregationValue, error) {
	from := time.Now().AddDate(0, 0, -(days - 1))
	var rows []AggregationValue
	err := r.db.Model(&model.DownloadRecord{}).
		Select("strftime('%Y-%m-%d', created_at) as label, count(*) as value").
		Where("resource_type = ? AND resource_id = ? AND created_at >= ?", "dataset", datasetID, from).
		Group("strftime('%Y-%m-%d', created_at)").
		Order("label asc").
		Scan(&rows).Error
	return rows, err
}

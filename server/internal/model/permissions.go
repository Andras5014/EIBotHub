package model

import "slices"

const (
	RoleSuperAdmin = "super_admin"
)

const (
	PermissionDashboardView            = "dashboard:view"
	PermissionOperationLogView         = "operation_log:view"
	PermissionPortalAccess             = "portal:access"
	PermissionPortalManage             = "portal:manage"
	PermissionAnnouncementManage       = "announcement:manage"
	PermissionSearchKeywordManage      = "search_keyword:manage"
	PermissionFeaturedResourceManage   = "featured_resource:manage"
	PermissionContentAccess            = "content:access"
	PermissionTemplateManage           = "template:manage"
	PermissionApplicationCaseManage    = "application_case:manage"
	PermissionDocumentCategoryManage   = "document_category:manage"
	PermissionDocumentManage           = "document:manage"
	PermissionFAQManage                = "faq:manage"
	PermissionVideoManage              = "video:manage"
	PermissionAgreementManage          = "agreement_template:manage"
	PermissionPrivacyManage            = "privacy_option:manage"
	PermissionFilterOptionManage       = "filter_option:manage"
	PermissionModelRecommendTagManage  = "model_recommend_tag:manage"
	PermissionReviewAccess             = "review:access"
	PermissionReviewManage             = "review:manage"
	PermissionDatasetAccessAccess      = "dataset_access:access"
	PermissionDatasetAccessReview      = "dataset_access:review"
	PermissionVerificationAccess       = "verification:access"
	PermissionVerificationReview       = "verification:review"
	PermissionCommunityAccess          = "community:access"
	PermissionCommunityContentModerate = "community_content:moderate"
	PermissionConversationModerate     = "conversation:moderate"
	PermissionWorkspaceModerate        = "workspace:moderate"
	PermissionWikiManage               = "wiki:manage"
	PermissionRewardAccess             = "reward:access"
	PermissionRewardManage             = "reward:manage"
)

var allPermissions = []string{
	PermissionDashboardView,
	PermissionOperationLogView,
	PermissionPortalAccess,
	PermissionPortalManage,
	PermissionAnnouncementManage,
	PermissionSearchKeywordManage,
	PermissionFeaturedResourceManage,
	PermissionContentAccess,
	PermissionTemplateManage,
	PermissionApplicationCaseManage,
	PermissionDocumentCategoryManage,
	PermissionDocumentManage,
	PermissionFAQManage,
	PermissionVideoManage,
	PermissionAgreementManage,
	PermissionPrivacyManage,
	PermissionFilterOptionManage,
	PermissionModelRecommendTagManage,
	PermissionReviewAccess,
	PermissionReviewManage,
	PermissionDatasetAccessAccess,
	PermissionDatasetAccessReview,
	PermissionVerificationAccess,
	PermissionVerificationReview,
	PermissionCommunityAccess,
	PermissionCommunityContentModerate,
	PermissionConversationModerate,
	PermissionWorkspaceModerate,
	PermissionWikiManage,
	PermissionRewardAccess,
	PermissionRewardManage,
}

var rolePermissionMatrix = map[string][]string{
	RoleUser:      {},
	RoleDeveloper: {},
	RoleOperator: {
		PermissionDashboardView,
		PermissionOperationLogView,
		PermissionPortalAccess,
		PermissionPortalManage,
		PermissionAnnouncementManage,
		PermissionSearchKeywordManage,
		PermissionFeaturedResourceManage,
		PermissionContentAccess,
		PermissionTemplateManage,
		PermissionApplicationCaseManage,
		PermissionDocumentCategoryManage,
		PermissionDocumentManage,
		PermissionFAQManage,
		PermissionVideoManage,
		PermissionAgreementManage,
		PermissionPrivacyManage,
		PermissionFilterOptionManage,
		PermissionModelRecommendTagManage,
		PermissionCommunityAccess,
		PermissionCommunityContentModerate,
		PermissionConversationModerate,
		PermissionWorkspaceModerate,
	},
	RoleReviewer: {
		PermissionDashboardView,
		PermissionOperationLogView,
		PermissionReviewAccess,
		PermissionReviewManage,
		PermissionDatasetAccessAccess,
		PermissionDatasetAccessReview,
		PermissionVerificationAccess,
		PermissionVerificationReview,
		PermissionCommunityAccess,
		PermissionCommunityContentModerate,
		PermissionConversationModerate,
		PermissionWorkspaceModerate,
	},
	RoleAdmin:      append([]string{}, allPermissions...),
	RoleSuperAdmin: append([]string{}, allPermissions...),
}

func RolePermissions(role string) []string {
	permissions, ok := rolePermissionMatrix[role]
	if !ok {
		return []string{}
	}
	return append([]string{}, permissions...)
}

func HasPermission(role, permission string) bool {
	if permission == "" {
		return true
	}
	return slices.Contains(RolePermissions(role), permission)
}

func HasAllPermissions(role string, permissions ...string) bool {
	for _, permission := range permissions {
		if !HasPermission(role, permission) {
			return false
		}
	}
	return true
}

func IsAdminRole(role string) bool {
	return role == RoleAdmin || role == RoleSuperAdmin
}

func HasBackofficeAccess(role string) bool {
	return len(RolePermissions(role)) > 0
}

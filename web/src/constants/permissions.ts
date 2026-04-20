export const ROLE_SUPER_ADMIN = 'super_admin';

export const PERMISSIONS = {
  dashboardView: 'dashboard:view',
  operationLogView: 'operation_log:view',
  portalAccess: 'portal:access',
  portalManage: 'portal:manage',
  announcementManage: 'announcement:manage',
  searchKeywordManage: 'search_keyword:manage',
  featuredResourceManage: 'featured_resource:manage',
  contentAccess: 'content:access',
  templateManage: 'template:manage',
  applicationCaseManage: 'application_case:manage',
  documentCategoryManage: 'document_category:manage',
  documentManage: 'document:manage',
  faqManage: 'faq:manage',
  videoManage: 'video:manage',
  agreementManage: 'agreement_template:manage',
  privacyManage: 'privacy_option:manage',
  filterOptionManage: 'filter_option:manage',
  modelRecommendTagManage: 'model_recommend_tag:manage',
  reviewAccess: 'review:access',
  reviewManage: 'review:manage',
  datasetAccessAccess: 'dataset_access:access',
  datasetAccessReview: 'dataset_access:review',
  verificationAccess: 'verification:access',
  verificationReview: 'verification:review',
  communityAccess: 'community:access',
  communityContentModerate: 'community_content:moderate',
  conversationModerate: 'conversation:moderate',
  workspaceModerate: 'workspace:moderate',
  wikiManage: 'wiki:manage',
  rewardAccess: 'reward:access',
  rewardManage: 'reward:manage',
} as const;

const allPermissions = Object.values(PERMISSIONS);

export const rolePermissionMatrix: Record<string, string[]> = {
  user: [],
  developer: [],
  operator: [
    PERMISSIONS.dashboardView,
    PERMISSIONS.operationLogView,
    PERMISSIONS.portalAccess,
    PERMISSIONS.portalManage,
    PERMISSIONS.announcementManage,
    PERMISSIONS.searchKeywordManage,
    PERMISSIONS.featuredResourceManage,
    PERMISSIONS.contentAccess,
    PERMISSIONS.templateManage,
    PERMISSIONS.applicationCaseManage,
    PERMISSIONS.documentCategoryManage,
    PERMISSIONS.documentManage,
    PERMISSIONS.faqManage,
    PERMISSIONS.videoManage,
    PERMISSIONS.agreementManage,
    PERMISSIONS.privacyManage,
    PERMISSIONS.filterOptionManage,
    PERMISSIONS.modelRecommendTagManage,
    PERMISSIONS.communityAccess,
    PERMISSIONS.communityContentModerate,
    PERMISSIONS.conversationModerate,
    PERMISSIONS.workspaceModerate,
  ],
  reviewer: [
    PERMISSIONS.dashboardView,
    PERMISSIONS.operationLogView,
    PERMISSIONS.reviewAccess,
    PERMISSIONS.reviewManage,
    PERMISSIONS.datasetAccessAccess,
    PERMISSIONS.datasetAccessReview,
    PERMISSIONS.verificationAccess,
    PERMISSIONS.verificationReview,
    PERMISSIONS.communityAccess,
    PERMISSIONS.communityContentModerate,
    PERMISSIONS.conversationModerate,
    PERMISSIONS.workspaceModerate,
  ],
  admin: [...allPermissions],
  [ROLE_SUPER_ADMIN]: [...allPermissions],
};

export function permissionsByRole(role?: string) {
  return rolePermissionMatrix[role ?? ''] ?? [];
}

export function resolvePermissions(user?: { role?: string; permissions?: string[] | null }) {
  if (user?.permissions?.length) {
    return [...user.permissions];
  }
  return permissionsByRole(user?.role);
}

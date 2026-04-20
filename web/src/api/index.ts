import type {
  AgreementTemplateItem,
  AnnouncementItem,
  ApplicationCaseItem,
  AuthPayload,
  BatchDatasetAccessDecisionPayload,
  DashboardPayload,
  DatasetDetail,
  DatasetAccessRequestItem,
  DatasetOptionsResponse,
  FilterOptionsResponse,
  DatasetPrivacyOptionItem,
  DownloadPackageTaskItem,
  DatasetSample,
  DatasetStatsResponse,
  DocumentCategoryItem,
  DocumentItem,
  FAQItem,
  FeaturedResourceItem,
  FavoriteRecord,
  FilterOptionConfigItem,
  FollowItem,
  FollowStats,
  HomePayload,
  HomeHighlightItem,
  HomeHeroConfigItem,
  RankingConfigItem,
  ModelEvaluationItem,
  ModelDetail,
  NotificationRecord,
  Paginated,
  RatingSummary,
  ResourceCard,
  ReviewItem,
  ScenePageDetailPayload,
  ScenePageItem,
  SearchHotItem,
  SearchResponse,
  SearchSuggestionItem,
  SearchKeywordConfigItem,
  SearchItem,
  SkillItem,
  TaskTemplateItem,
  UploadsPayload,
  UserContributionPayload,
  UserSummary,
  DownloadRecord,
  CommentItem,
  DiscussionItem,
  AdminCommunityOverview,
  AdminConversationModerationItem,
  AdminSkillModerationItem,
  AdminDiscussionModerationItem,
  AdminCommentModerationItem,
  AdminWorkspaceModerationItem,
  AdminModelRecommendTagItem,
  AdminVerificationItem,
  AdminOperationLogItem,
  AdminRewardOverview,
  AdminRewardAdjustmentItem,
  VerificationStatusItem,
  WikiPageItem,
  WikiRevisionItem,
  RewardSummary,
  RewardLedgerItem,
  ContributorRankingItem,
  RewardBenefitItem,
  RewardRedemptionItem,
  OpenAPISpecResponse,
  WebhookSubscriptionItem,
  WebhookDeliveryItem,
  MessageItem,
  ModuleSettingItem,
  ConversationItem,
  VideoTutorialItem,
  WorkspaceItem,
  WorkspaceDetail,
} from '@/types/api';
import { request, sendWithoutData } from './client';

export const api = {
  register: (payload: { username: string; email: string; password: string }) =>
    request<AuthPayload>({ url: '/auth/register', method: 'post', data: payload }),
  login: (payload: { email: string; password: string }) =>
    request<AuthPayload>({ url: '/auth/login', method: 'post', data: payload }),
  logout: () => sendWithoutData({ url: '/auth/logout', method: 'post' }),
  me: () => request<UserSummary>({ url: '/users/me', method: 'get' }),
  home: () => request<HomePayload>({ url: '/portal/home', method: 'get' }),
  listScenePages: () => request<ScenePageItem[]>({ url: '/scenes', method: 'get' }),
  getScenePage: (slug: string) => request<ScenePageDetailPayload>({ url: `/scenes/${slug}`, method: 'get' }),
  search: (params: Record<string, unknown>) =>
    request<SearchResponse>({ url: '/search', method: 'get', params }),
  hotQueries: () => request<SearchHotItem[]>({ url: '/search/hot', method: 'get' }),
  recommendedQueries: () => request<SearchSuggestionItem[]>({ url: '/search/recommended', method: 'get' }),
  listModels: (params: Record<string, unknown>) =>
    request<Paginated<ResourceCard>>({ url: '/models', method: 'get', params }),
  getModel: (id: string | number) => request<ModelDetail>({ url: `/models/${id}`, method: 'get' }),
  getModelEvaluations: (id: string | number) =>
    request<ModelEvaluationItem[]>({ url: `/models/${id}/evaluations`, method: 'get' }),
  createModelEvaluation: (id: string | number, payload: { benchmark: string; summary: string; score: number; notes: string }) =>
    request<ModelEvaluationItem>({ url: `/models/${id}/evaluations`, method: 'post', data: payload }),
  getModelRatings: (id: string | number) =>
    request<RatingSummary>({ url: `/models/${id}/ratings`, method: 'get' }),
  rateModel: (id: string | number, payload: { score: number; feedback: string }) =>
    request<RatingSummary>({ url: `/models/${id}/ratings`, method: 'post', data: payload }),
  getModelComments: (id: string | number) =>
    request<CommentItem[]>({ url: `/models/${id}/comments`, method: 'get' }),
  commentModel: (id: string | number, payload: { content: string; parent_id?: number }) =>
    request<CommentItem>({ url: `/models/${id}/comments`, method: 'post', data: payload }),
  createModel: (formData: FormData) => request<ModelDetail>({ url: '/models', method: 'post', data: formData }),
  updateModel: (
    id: string | number,
    payload: {
      name: string;
      summary: string;
      description: string;
      tags: string;
      robot_type: string;
      input_spec: string;
      output_spec: string;
      license: string;
      dependencies: string;
    },
  ) => request<ModelDetail>({ url: `/models/${id}`, method: 'put', data: payload }),
  submitModel: (id: string | number) => sendWithoutData({ url: `/models/${id}/submit`, method: 'post' }),
  addModelVersion: (id: string | number, formData: FormData) =>
    request<ModelDetail>({ url: `/models/${id}/versions`, method: 'post', data: formData }),
  downloadModel: (id: string | number) => sendWithoutData({ url: `/models/${id}/download`, method: 'post' }),
  listDatasets: (params: Record<string, unknown>) =>
    request<Paginated<ResourceCard>>({ url: '/datasets', method: 'get', params }),
  getDatasetOptions: () => request<DatasetOptionsResponse>({ url: '/datasets/options', method: 'get' }),
  getFilterOptions: () => request<FilterOptionsResponse>({ url: '/filter-options', method: 'get' }),
  getDataset: (id: string | number) => request<DatasetDetail>({ url: `/datasets/${id}`, method: 'get' }),
  getDatasetSamples: (id: string | number) => request<DatasetSample[]>({ url: `/datasets/${id}/samples`, method: 'get' }),
  getDatasetStats: (id: string | number) => request<DatasetStatsResponse>({ url: `/datasets/${id}/stats`, method: 'get' }),
  getDatasetRatings: (id: string | number) => request<RatingSummary>({ url: `/datasets/${id}/ratings`, method: 'get' }),
  rateDataset: (id: string | number, payload: { score: number; feedback: string }) =>
    request<RatingSummary>({ url: `/datasets/${id}/ratings`, method: 'post', data: payload }),
  getDatasetComments: (id: string | number) => request<CommentItem[]>({ url: `/datasets/${id}/comments`, method: 'get' }),
  commentDataset: (id: string | number, payload: { content: string; parent_id?: number }) =>
    request<CommentItem>({ url: `/datasets/${id}/comments`, method: 'post', data: payload }),
  createDataset: (formData: FormData) => request<DatasetDetail>({ url: '/datasets', method: 'post', data: formData }),
  updateDataset: (
    id: string | number,
    payload: {
      name: string;
      summary: string;
      description: string;
      tags: string;
      sample_count: number;
      device: string;
      scene: string;
      privacy: string;
      agreement_text: string;
      sample_preview: string;
    },
  ) => request<DatasetDetail>({ url: `/datasets/${id}`, method: 'put', data: payload }),
  addDatasetVersion: (id: string | number, formData: FormData) =>
    request<DatasetDetail>({ url: `/datasets/${id}/versions`, method: 'post', data: formData }),
  submitDataset: (id: string | number) => sendWithoutData({ url: `/datasets/${id}/submit`, method: 'post' }),
  confirmDatasetAgreement: (id: string | number) =>
    sendWithoutData({ url: `/datasets/${id}/agreements/confirm`, method: 'post' }),
  downloadDataset: (id: string | number) => sendWithoutData({ url: `/datasets/${id}/download`, method: 'post' }),
  getDatasetDownloadPackages: (id: string | number) =>
    request<DownloadPackageTaskItem[]>({ url: `/datasets/${id}/download-packages`, method: 'get' }),
  createDatasetDownloadPackage: (id: string | number, payload: { parts?: number }) =>
    request<DownloadPackageTaskItem>({ url: `/datasets/${id}/download-packages`, method: 'post', data: payload }),
  getMyDatasetAccessRequest: (id: string | number) =>
    request<DatasetAccessRequestItem | null>({ url: `/datasets/${id}/access-requests/me`, method: 'get' }),
  getMyDatasetAccessHistory: (id: string | number) =>
    request<DatasetAccessRequestItem[]>({ url: `/datasets/${id}/access-requests/history`, method: 'get' }),
  createDatasetAccessRequest: (id: string | number, payload: { reason: string }) =>
    request<DatasetAccessRequestItem>({ url: `/datasets/${id}/access-requests`, method: 'post', data: payload }),
  listTemplates: () => request<TaskTemplateItem[]>({ url: '/task-templates', method: 'get' }),
  getTemplate: (id: string | number) => request<TaskTemplateItem>({ url: `/task-templates/${id}`, method: 'get' }),
  getTemplateRatings: (id: string | number) => request<RatingSummary>({ url: `/task-templates/${id}/ratings`, method: 'get' }),
  rateTemplate: (id: string | number, payload: { score: number; feedback: string }) =>
    request<RatingSummary>({ url: `/task-templates/${id}/ratings`, method: 'post', data: payload }),
  getTemplateComments: (id: string | number) =>
    request<CommentItem[]>({ url: `/task-templates/${id}/comments`, method: 'get' }),
  commentTemplate: (id: string | number, payload: { content: string; parent_id?: number }) =>
    request<CommentItem>({ url: `/task-templates/${id}/comments`, method: 'post', data: payload }),
  listApplicationCases: () => request<ApplicationCaseItem[]>({ url: '/application-cases', method: 'get' }),
  getApplicationCase: (id: string | number) =>
    request<ApplicationCaseItem>({ url: `/application-cases/${id}`, method: 'get' }),
  listSkills: () => request<SkillItem[]>({ url: '/skills', method: 'get' }),
  getSkill: (id: string | number) => request<SkillItem>({ url: `/skills/${id}`, method: 'get' }),
  createSkill: (payload: { name: string; summary: string; description: string; category: string; scene: string; guide: string; resource_ref: string }) =>
    request<SkillItem>({ url: '/skills', method: 'post', data: payload }),
  forkSkill: (id: string | number) => request<SkillItem>({ url: `/skills/${id}/fork`, method: 'post' }),
  getSkillRatings: (id: string | number) => request<RatingSummary>({ url: `/skills/${id}/ratings`, method: 'get' }),
  rateSkill: (id: string | number, payload: { score: number; feedback: string }) =>
    request<RatingSummary>({ url: `/skills/${id}/ratings`, method: 'post', data: payload }),
  getSkillComments: (id: string | number) =>
    request<CommentItem[]>({ url: `/skills/${id}/comments`, method: 'get' }),
  commentSkill: (id: string | number, payload: { content: string; parent_id?: number }) =>
    request<CommentItem>({ url: `/skills/${id}/comments`, method: 'post', data: payload }),
  listDocCategories: (docType = '') =>
    request<DocumentCategoryItem[]>({ url: '/docs/categories', method: 'get', params: { doc_type: docType } }),
  listDocs: (docType = '') => request<DocumentItem[]>({ url: '/docs', method: 'get', params: { doc_type: docType } }),
  getDoc: (id: string | number) => request<DocumentItem>({ url: `/docs/${id}`, method: 'get' }),
  listFaqs: () => request<FAQItem[]>({ url: '/faqs', method: 'get' }),
  listVideos: () => request<VideoTutorialItem[]>({ url: '/videos', method: 'get' }),
  listWikiPages: () => request<WikiPageItem[]>({ url: '/wiki/pages', method: 'get' }),
  getWikiPage: (id: string | number) => request<WikiPageItem>({ url: `/wiki/pages/${id}`, method: 'get' }),
  createWikiPage: (payload: { title: string; summary: string; content: string; comment: string }) =>
    request<WikiPageItem>({ url: '/wiki/pages', method: 'post', data: payload }),
  updateWikiPage: (id: string | number, payload: { title: string; summary: string; content: string; comment: string }) =>
    request<WikiPageItem>({ url: `/wiki/pages/${id}`, method: 'put', data: payload }),
  getWikiRevisions: (id: string | number) => request<WikiRevisionItem[]>({ url: `/wiki/pages/${id}/revisions`, method: 'get' }),
  getPublicUser: (id: string | number) => request<UserSummary>({ url: `/community/users/${id}`, method: 'get' }),
  getPublicUserContributions: (id: string | number) =>
    request<UserContributionPayload>({ url: `/community/users/${id}/contributions`, method: 'get' }),
  getPublicUserFollowStats: (id: string | number) =>
    request<FollowStats>({ url: `/community/users/${id}/follow-stats`, method: 'get' }),
  toggleFollow: (id: string | number) => sendWithoutData({ url: `/community/users/${id}/follow`, method: 'post' }),
  myFollows: () => request<FollowItem[]>({ url: '/users/me/follows', method: 'get' }),
  myFollowers: () => request<FollowItem[]>({ url: '/users/me/followers', method: 'get' }),
  myFollowStats: () => request<FollowStats>({ url: '/users/me/follow-stats', method: 'get' }),
  myContributions: () => request<UserContributionPayload>({ url: '/users/me/contributions', method: 'get' }),
  myRewardSummary: () => request<RewardSummary>({ url: '/rewards/points', method: 'get' }),
  myRewardLedger: () => request<RewardLedgerItem[]>({ url: '/rewards/ledger', method: 'get' }),
  rewardBenefits: () => request<RewardBenefitItem[]>({ url: '/rewards/benefits', method: 'get' }),
  myRewardRedemptions: () => request<RewardRedemptionItem[]>({ url: '/rewards/redemptions', method: 'get' }),
  redeemRewardBenefit: (payload: { benefit_id: number }) =>
    request<RewardRedemptionItem>({ url: '/rewards/redeem', method: 'post', data: payload }),
  openAPISpec: () => request<OpenAPISpecResponse>({ url: '/openapi/spec', method: 'get' }),
  listWebhooks: () => request<WebhookSubscriptionItem[]>({ url: '/webhooks', method: 'get' }),
  createWebhook: (payload: { name: string; target_url: string; secret: string; events: string[] }) =>
    request<WebhookSubscriptionItem>({ url: '/webhooks', method: 'post', data: payload }),
  getWebhookDeliveries: (id: string | number) =>
    request<WebhookDeliveryItem[]>({ url: `/webhooks/${id}/deliveries`, method: 'get' }),
  testWebhook: (id: string | number) =>
    request<WebhookDeliveryItem>({ url: `/webhooks/${id}/test`, method: 'post' }),
  contributorRankings: () => request<ContributorRankingItem[]>({ url: '/rankings/contributors', method: 'get' }),
  listConversations: () => request<ConversationItem[]>({ url: '/messages/conversations', method: 'get' }),
  getConversationMessages: (id: string | number) => request<MessageItem[]>({ url: `/messages/conversations/${id}`, method: 'get' }),
  sendMessage: (payload: { conversation_id?: number; recipient_user_id?: number; content: string }) =>
    request<MessageItem>({ url: '/messages', method: 'post', data: payload }),
  listWorkspaces: () => request<WorkspaceItem[]>({ url: '/workspaces', method: 'get' }),
  getWorkspace: (id: string | number) => request<WorkspaceDetail>({ url: `/workspaces/${id}`, method: 'get' }),
  createWorkspace: (payload: { name: string; summary: string }) =>
    request<WorkspaceDetail>({ url: '/workspaces', method: 'post', data: payload }),
  addWorkspaceMember: (id: string | number, payload: { user_id: number }) =>
    sendWithoutData({ url: `/workspaces/${id}/members`, method: 'post', data: payload }),
  sendWorkspaceMessage: (id: string | number, payload: { content: string }) =>
    request<MessageItem>({ url: `/workspaces/${id}/messages`, method: 'post', data: payload }),
  listDiscussions: (params: { q?: string; tag?: string; sort?: 'hot' | 'latest' | 'comments'; limit?: number } = {}) =>
    request<DiscussionItem[]>({ url: '/discussions', method: 'get', params }),
  getDiscussion: (id: string | number) => request<DiscussionItem>({ url: `/discussions/${id}`, method: 'get' }),
  createDiscussion: (payload: { title: string; tag: string; content: string }) =>
    request<DiscussionItem>({ url: '/discussions', method: 'post', data: payload }),
  getDiscussionComments: (id: string | number) => request<CommentItem[]>({ url: `/discussions/${id}/comments`, method: 'get' }),
  commentDiscussion: (id: string | number, payload: { content: string; parent_id?: number }) =>
    request<CommentItem>({ url: `/discussions/${id}/comments`, method: 'post', data: payload }),
  getProfile: () => request<UserSummary>({ url: '/users/me/profile', method: 'get' }),
  updateProfile: (payload: { username: string; bio: string; avatar: string }) =>
    request<UserSummary>({ url: '/users/me/profile', method: 'put', data: payload }),
  getMyDatasetAccessRequests: () => request<DatasetAccessRequestItem[]>({ url: '/users/me/dataset-access-requests', method: 'get' }),
  getUploads: () => request<UploadsPayload>({ url: '/users/me/uploads', method: 'get' }),
  getFavorites: () => request<FavoriteRecord[]>({ url: '/users/me/favorites', method: 'get' }),
  toggleFavorite: (payload: { resource_type: string; resource_id: number; title: string }) =>
    sendWithoutData({ url: '/favorites/toggle', method: 'post', data: payload }),
  getDownloads: () => request<DownloadRecord[]>({ url: '/users/me/downloads', method: 'get' }),
  getNotifications: () => request<NotificationRecord[]>({ url: '/users/me/notifications', method: 'get' }),
  readNotifications: () => sendWithoutData({ url: '/users/me/notifications/read', method: 'post' }),
  applyVerification: (payload: { verification_type: 'personal' | 'enterprise'; real_name: string; organization: string; materials: string; reason: string }) =>
    request<VerificationStatusItem>({ url: '/developer-verifications', method: 'post', data: payload }),
  applyEnterpriseVerification: (payload: { real_name: string; organization: string; materials: string; reason: string }) =>
    request<VerificationStatusItem>({ url: '/enterprise-verifications', method: 'post', data: payload }),
  myVerification: () => request<VerificationStatusItem | null>({ url: '/users/me/verification', method: 'get' }),
  getPublicVerification: (id: string | number) => request<VerificationStatusItem | null>({ url: `/community/users/${id}/verification`, method: 'get' }),
  getAdminDashboard: () => request<DashboardPayload>({ url: '/admin/dashboard', method: 'get' }),
  getAdminReviews: (type: 'models' | 'datasets') =>
    request<ReviewItem[]>({ url: '/admin/reviews', method: 'get', params: { type } }),
  decideReview: (type: 'models' | 'datasets', id: number, payload: { decision: 'approved' | 'rejected'; comment: string }) =>
    sendWithoutData({ url: `/admin/reviews/${type}/${id}/decision`, method: 'post', data: payload }),
  getAdminDatasetAccessRequests: (params: Record<string, unknown> = {}) =>
    request<DatasetAccessRequestItem[]>({ url: '/admin/datasets/access-requests', method: 'get', params }),
  reviewAdminDatasetAccessRequest: (id: number, payload: { decision: 'approved' | 'rejected'; comment: string; valid_days: number; download_limit: number }) =>
    sendWithoutData({ url: `/admin/datasets/access-requests/${id}/decision`, method: 'post', data: payload }),
  batchReviewAdminDatasetAccessRequests: (payload: BatchDatasetAccessDecisionPayload) =>
    sendWithoutData({ url: '/admin/datasets/access-requests/batch-decision', method: 'post', data: payload }),
  getAdminAnnouncements: () => request<AnnouncementItem[]>({ url: '/admin/announcements', method: 'get' }),
  createAnnouncement: (payload: { title: string; summary: string; link: string; pinned: boolean }) =>
    request<AnnouncementItem>({ url: '/admin/announcements', method: 'post', data: payload }),
  getAdminPortalModules: () => request<ModuleSettingItem[]>({ url: '/admin/portal/modules', method: 'get' }),
  updateAdminPortalModule: (key: string, payload: { enabled: boolean; sort_order: number }) =>
    sendWithoutData({ url: `/admin/portal/modules/${key}`, method: 'put', data: payload }),
  getAdminHomeHeroConfig: () => request<HomeHeroConfigItem>({ url: '/admin/portal/hero-config', method: 'get' }),
  updateAdminHomeHeroConfig: (payload: {
    tagline: string;
    title: string;
    description: string;
    primary_button: string;
    secondary_button: string;
    search_button: string;
  }) => request<HomeHeroConfigItem>({ url: '/admin/portal/hero-config', method: 'put', data: payload }),
  getAdminHomeHighlights: () => request<HomeHighlightItem[]>({ url: '/admin/portal/highlights', method: 'get' }),
  createAdminHomeHighlight: (payload: { text: string; sort_order: number; enabled: boolean }) =>
    request<HomeHighlightItem>({ url: '/admin/portal/highlights', method: 'post', data: payload }),
  updateAdminHomeHighlight: (id: number, payload: { text: string; sort_order: number; enabled: boolean }) =>
    request<HomeHighlightItem>({ url: `/admin/portal/highlights/${id}`, method: 'put', data: payload }),
  deleteAdminHomeHighlight: (id: number) =>
    sendWithoutData({ url: `/admin/portal/highlights/${id}`, method: 'delete' }),
  getAdminScenePages: () => request<ScenePageItem[]>({ url: '/admin/portal/scenes', method: 'get' }),
  createAdminScenePage: (payload: { slug: string; name: string; tagline: string; summary: string; description: string; sort_order: number; enabled: boolean }) =>
    request<ScenePageItem>({ url: '/admin/portal/scenes', method: 'post', data: payload }),
  updateAdminScenePage: (id: number, payload: { slug: string; name: string; tagline: string; summary: string; description: string; sort_order: number; enabled: boolean }) =>
    request<ScenePageItem>({ url: `/admin/portal/scenes/${id}`, method: 'put', data: payload }),
  deleteAdminScenePage: (id: number) =>
    sendWithoutData({ url: `/admin/portal/scenes/${id}`, method: 'delete' }),
  getAdminRankingConfig: () => request<RankingConfigItem>({ url: '/admin/portal/rankings-config', method: 'get' }),
  updateAdminRankingConfig: (payload: { title: string; subtitle: string; limit: number; enabled: boolean }) =>
    request<RankingConfigItem>({ url: '/admin/portal/rankings-config', method: 'put', data: payload }),
  getAdminSearchKeywords: (keywordType = '') =>
    request<SearchKeywordConfigItem[]>({ url: '/admin/portal/search-keywords', method: 'get', params: { keyword_type: keywordType || undefined } }),
  createAdminSearchKeyword: (payload: { query: string; keyword_type: string; sort_order: number; enabled: boolean }) =>
    request<SearchKeywordConfigItem>({ url: '/admin/portal/search-keywords', method: 'post', data: payload }),
  updateAdminSearchKeyword: (id: number, payload: { query: string; keyword_type: string; sort_order: number; enabled: boolean }) =>
    request<SearchKeywordConfigItem>({ url: `/admin/portal/search-keywords/${id}`, method: 'put', data: payload }),
  deleteAdminSearchKeyword: (id: number) =>
    sendWithoutData({ url: `/admin/portal/search-keywords/${id}`, method: 'delete' }),
  getAdminFeaturedResources: () => request<FeaturedResourceItem[]>({ url: '/admin/portal/featured-resources', method: 'get' }),
  createAdminFeaturedResource: (payload: { resource_type: string; resource_id: number; badge_label: string; sort_order: number; enabled: boolean }) =>
    request<FeaturedResourceItem>({ url: '/admin/portal/featured-resources', method: 'post', data: payload }),
  updateAdminFeaturedResource: (id: number, payload: { resource_type: string; resource_id: number; badge_label: string; sort_order: number; enabled: boolean }) =>
    request<FeaturedResourceItem>({ url: `/admin/portal/featured-resources/${id}`, method: 'put', data: payload }),
  deleteAdminFeaturedResource: (id: number) =>
    sendWithoutData({ url: `/admin/portal/featured-resources/${id}`, method: 'delete' }),
  getAdminTemplates: () => request<TaskTemplateItem[]>({ url: '/admin/content/templates', method: 'get' }),
  createAdminTemplate: (payload: { name: string; summary: string; description: string; category: string; scene: string; guide: string; resource_ref: string; status: string }) =>
    request<TaskTemplateItem>({ url: '/admin/content/templates', method: 'post', data: payload }),
  updateAdminTemplate: (id: number, payload: { name: string; summary: string; description: string; category: string; scene: string; guide: string; resource_ref: string; status: string }) =>
    request<TaskTemplateItem>({ url: `/admin/content/templates/${id}`, method: 'put', data: payload }),
  deleteAdminTemplate: (id: number) =>
    sendWithoutData({ url: `/admin/content/templates/${id}`, method: 'delete' }),
  batchAdminTemplateStatus: (payload: { ids: number[]; status: string }) =>
    sendWithoutData({ url: '/admin/content/templates/status', method: 'post', data: payload }),
  batchDeleteAdminTemplates: (payload: { ids: number[] }) =>
    sendWithoutData({ url: '/admin/content/templates/delete', method: 'post', data: payload }),
  getAdminApplicationCases: () => request<ApplicationCaseItem[]>({ url: '/admin/content/application-cases', method: 'get' }),
  createAdminApplicationCase: (payload: { title: string; summary: string; category: string; guide: string; cover_image: string; status: string }) =>
    request<ApplicationCaseItem>({ url: '/admin/content/application-cases', method: 'post', data: payload }),
  updateAdminApplicationCase: (id: number, payload: { title: string; summary: string; category: string; guide: string; cover_image: string; status: string }) =>
    request<ApplicationCaseItem>({ url: `/admin/content/application-cases/${id}`, method: 'put', data: payload }),
  deleteAdminApplicationCase: (id: number) =>
    sendWithoutData({ url: `/admin/content/application-cases/${id}`, method: 'delete' }),
  batchAdminApplicationCaseStatus: (payload: { ids: number[]; status: string }) =>
    sendWithoutData({ url: '/admin/content/application-cases/status', method: 'post', data: payload }),
  batchDeleteAdminApplicationCases: (payload: { ids: number[] }) =>
    sendWithoutData({ url: '/admin/content/application-cases/delete', method: 'post', data: payload }),
  getAdminDocCategories: () => request<DocumentCategoryItem[]>({ url: '/admin/content/doc-categories', method: 'get' }),
  createAdminDocCategory: (payload: { name: string; doc_type: string }) =>
    request<DocumentCategoryItem>({ url: '/admin/content/doc-categories', method: 'post', data: payload }),
  updateAdminDocCategory: (id: number, payload: { name: string; doc_type: string }) =>
    request<DocumentCategoryItem>({ url: `/admin/content/doc-categories/${id}`, method: 'put', data: payload }),
  deleteAdminDocCategory: (id: number) =>
    sendWithoutData({ url: `/admin/content/doc-categories/${id}`, method: 'delete' }),
  getAdminDocuments: () => request<DocumentItem[]>({ url: '/admin/content/docs', method: 'get' }),
  createAdminDocument: (payload: { category_id: number; title: string; summary: string; content: string; doc_type: string; status: string }) =>
    request<DocumentItem>({ url: '/admin/content/docs', method: 'post', data: payload }),
  updateAdminDocument: (id: number, payload: { category_id: number; title: string; summary: string; content: string; doc_type: string; status: string }) =>
    request<DocumentItem>({ url: `/admin/content/docs/${id}`, method: 'put', data: payload }),
  deleteAdminDocument: (id: number) =>
    sendWithoutData({ url: `/admin/content/docs/${id}`, method: 'delete' }),
  batchAdminDocumentStatus: (payload: { ids: number[]; status: string }) =>
    sendWithoutData({ url: '/admin/content/docs/status', method: 'post', data: payload }),
  batchDeleteAdminDocuments: (payload: { ids: number[] }) =>
    sendWithoutData({ url: '/admin/content/docs/delete', method: 'post', data: payload }),
  getAdminFaqs: () => request<FAQItem[]>({ url: '/admin/content/faqs', method: 'get' }),
  createAdminFaq: (payload: { question: string; answer: string }) =>
    request<FAQItem>({ url: '/admin/content/faqs', method: 'post', data: payload }),
  updateAdminFaq: (id: number, payload: { question: string; answer: string }) =>
    request<FAQItem>({ url: `/admin/content/faqs/${id}`, method: 'put', data: payload }),
  deleteAdminFaq: (id: number) =>
    sendWithoutData({ url: `/admin/content/faqs/${id}`, method: 'delete' }),
  getAdminVideos: () => request<VideoTutorialItem[]>({ url: '/admin/content/videos', method: 'get' }),
  createAdminVideo: (payload: { title: string; summary: string; link: string; category: string; sort_order: number; active: boolean }) =>
    request<VideoTutorialItem>({ url: '/admin/content/videos', method: 'post', data: payload }),
  updateAdminVideo: (id: number, payload: { title: string; summary: string; link: string; category: string; sort_order: number; active: boolean }) =>
    request<VideoTutorialItem>({ url: `/admin/content/videos/${id}`, method: 'put', data: payload }),
  deleteAdminVideo: (id: number) =>
    sendWithoutData({ url: `/admin/content/videos/${id}`, method: 'delete' }),
  batchAdminVideoStatus: (payload: { ids: number[]; active: boolean }) =>
    sendWithoutData({ url: '/admin/content/videos/status', method: 'post', data: payload }),
  batchDeleteAdminVideos: (payload: { ids: number[] }) =>
    sendWithoutData({ url: '/admin/content/videos/delete', method: 'post', data: payload }),
  getAdminAgreementTemplates: () => request<AgreementTemplateItem[]>({ url: '/admin/content/agreement-templates', method: 'get' }),
  createAdminAgreementTemplate: (payload: { name: string; content: string; sort_order: number; active: boolean }) =>
    request<AgreementTemplateItem>({ url: '/admin/content/agreement-templates', method: 'post', data: payload }),
  updateAdminAgreementTemplate: (id: number, payload: { name: string; content: string; sort_order: number; active: boolean }) =>
    request<AgreementTemplateItem>({ url: `/admin/content/agreement-templates/${id}`, method: 'put', data: payload }),
  deleteAdminAgreementTemplate: (id: number) =>
    sendWithoutData({ url: `/admin/content/agreement-templates/${id}`, method: 'delete' }),
  getAdminPrivacyOptions: () => request<DatasetPrivacyOptionItem[]>({ url: '/admin/content/privacy-options', method: 'get' }),
  createAdminPrivacyOption: (payload: { code: string; name: string; description: string; sort_order: number; active: boolean }) =>
    request<DatasetPrivacyOptionItem>({ url: '/admin/content/privacy-options', method: 'post', data: payload }),
  updateAdminPrivacyOption: (id: number, payload: { code: string; name: string; description: string; sort_order: number; active: boolean }) =>
    request<DatasetPrivacyOptionItem>({ url: `/admin/content/privacy-options/${id}`, method: 'put', data: payload }),
  deleteAdminPrivacyOption: (id: number) =>
    sendWithoutData({ url: `/admin/content/privacy-options/${id}`, method: 'delete' }),
  getAdminFilterOptions: (kind = '') =>
    request<FilterOptionConfigItem[]>({ url: '/admin/content/filter-options', method: 'get', params: { kind: kind || undefined } }),
  createAdminFilterOption: (payload: { kind: string; value: string; sort_order: number; enabled: boolean }) =>
    request<FilterOptionConfigItem>({ url: '/admin/content/filter-options', method: 'post', data: payload }),
  updateAdminFilterOption: (id: number, payload: { kind: string; value: string; sort_order: number; enabled: boolean }) =>
    request<FilterOptionConfigItem>({ url: `/admin/content/filter-options/${id}`, method: 'put', data: payload }),
  deleteAdminFilterOption: (id: number) =>
    sendWithoutData({ url: `/admin/content/filter-options/${id}`, method: 'delete' }),
  getAdminModelRecommendTags: () => request<AdminModelRecommendTagItem[]>({ url: '/admin/content/model-recommend-tags', method: 'get' }),
  updateAdminModelRecommendTag: (id: number, payload: { recommend_tag: string }) =>
    sendWithoutData({ url: `/admin/content/model-recommend-tags/${id}`, method: 'put', data: payload }),
  getAdminOperationLogs: () => request<AdminOperationLogItem[]>({ url: '/admin/operations', method: 'get' }),
  getAdminCommunityOverview: () => request<AdminCommunityOverview>({ url: '/admin/community/overview', method: 'get' }),
  getAdminCommunitySkills: () => request<AdminSkillModerationItem[]>({ url: '/admin/community/skills', method: 'get' }),
  hideAdminSkill: (id: number) => sendWithoutData({ url: `/admin/community/skills/${id}/hide`, method: 'post' }),
  getAdminCommunityDiscussions: () => request<AdminDiscussionModerationItem[]>({ url: '/admin/community/discussions', method: 'get' }),
  removeAdminDiscussion: (id: number) => sendWithoutData({ url: `/admin/community/discussions/${id}/remove`, method: 'post' }),
  getAdminCommunityComments: () => request<AdminCommentModerationItem[]>({ url: '/admin/community/comments', method: 'get' }),
  removeAdminComment: (id: number) => sendWithoutData({ url: `/admin/community/comments/${id}/remove`, method: 'post' }),
  getAdminCommunityConversations: () =>
    request<AdminConversationModerationItem[]>({ url: '/admin/community/conversations', method: 'get' }),
  blockAdminConversation: (id: number, payload: { reason: string }) =>
    sendWithoutData({ url: `/admin/community/conversations/${id}/block`, method: 'post', data: payload }),
  unblockAdminConversation: (id: number) =>
    sendWithoutData({ url: `/admin/community/conversations/${id}/unblock`, method: 'post' }),
  getAdminCommunityWorkspaces: () =>
    request<AdminWorkspaceModerationItem[]>({ url: '/admin/community/workspaces', method: 'get' }),
  blockAdminWorkspace: (id: number, payload: { reason: string }) =>
    sendWithoutData({ url: `/admin/community/workspaces/${id}/block`, method: 'post', data: payload }),
  unblockAdminWorkspace: (id: number) =>
    sendWithoutData({ url: `/admin/community/workspaces/${id}/unblock`, method: 'post' }),
  removeAdminWorkspaceMember: (id: number, payload: { user_id: number }) =>
    sendWithoutData({ url: `/admin/community/workspaces/${id}/remove-member`, method: 'post', data: payload }),
  getAdminVerifications: () => request<AdminVerificationItem[]>({ url: '/admin/verifications', method: 'get' }),
  reviewAdminVerification: (id: number, payload: { decision: 'approved' | 'rejected'; comment: string }) =>
    sendWithoutData({ url: `/admin/verifications/${id}/decision`, method: 'post', data: payload }),
  getAdminRewardOverview: () => request<AdminRewardOverview>({ url: '/admin/rewards/overview', method: 'get' }),
  getAdminRewardBenefits: () => request<RewardBenefitItem[]>({ url: '/admin/rewards/benefits', method: 'get' }),
  updateAdminRewardBenefit: (id: number, payload: { name: string; summary: string; cost_points: number; active: boolean }) =>
    request<RewardBenefitItem>({ url: `/admin/rewards/benefits/${id}`, method: 'put', data: payload }),
  getAdminRewardAdjustments: () => request<AdminRewardAdjustmentItem[]>({ url: '/admin/rewards/adjustments', method: 'get' }),
  createAdminRewardAdjustment: (payload: { user_id: number; points: number; remark: string }) =>
    request<AdminRewardAdjustmentItem>({ url: '/admin/rewards/adjustments', method: 'post', data: payload }),
};

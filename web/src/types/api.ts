export interface Envelope<T> {
  success: boolean;
  message?: string;
  data: T;
  error?: {
    code: string;
    message: string;
  };
}

export interface UserSummary {
  id: number;
  username: string;
  email: string;
  role: string;
  permissions?: string[];
  bio: string;
  avatar: string;
}

export interface AuthPayload {
  token: string;
  user: UserSummary;
}

export interface AnnouncementItem {
  id: number;
  title: string;
  summary: string;
  link: string;
  pinned: boolean;
  published_at: string;
}

export interface ModuleSettingItem {
  id?: number;
  module_key: string;
  label: string;
  sort_order: number;
  enabled: boolean;
  updated_at?: string;
}

export interface HomeHighlightItem {
  id: number;
  text: string;
  sort_order: number;
  enabled: boolean;
  updated_at: string;
}

export interface HomeHeroConfigItem {
  id?: number;
  tagline: string;
  title: string;
  description: string;
  primary_button: string;
  secondary_button: string;
  search_button: string;
  updated_at?: string;
}

export interface RankingConfigItem {
  id?: number;
  title: string;
  subtitle: string;
  limit: number;
  enabled: boolean;
  updated_at?: string;
}

export interface ScenePageItem {
  id: number;
  slug: string;
  name: string;
  tagline: string;
  summary: string;
  description: string;
  sort_order: number;
  enabled: boolean;
  updated_at: string;
}

export interface FeaturedResourceItem {
  id: number;
  resource_type: string;
  resource_id: number;
  title: string;
  summary: string;
  route: string;
  badge_label: string;
  sort_order: number;
  enabled: boolean;
  updated_at: string;
}

export interface ResourceCard {
  id: number;
  name: string;
  summary: string;
  description?: string;
  type: string;
  tags: string[];
  robot_type?: string;
  badge_label?: string;
  downloads: number;
  status?: string;
  owner?: string;
  review_comment?: string;
  updated_at: string;
}

export interface FileVersion {
  id: number;
  version: string;
  file_name: string;
  file_url: string;
  changelog: string;
  created_at: string;
}

export interface ModelDetail extends ResourceCard {
  review_comment?: string;
  recommend_tag?: string;
  input_spec: string;
  output_spec: string;
  license: string;
  dependencies: string[];
  versions: FileVersion[];
  favorited: boolean;
}

export interface DatasetSample {
  id: number;
  sample_type: string;
  title: string;
  preview_text: string;
  preview_url?: string;
  file_name?: string;
}

export interface DownloadPackageTaskItem {
  id: number;
  dataset_id: number;
  status: string;
  bundle_url?: string;
  part_links: string[];
  total_parts: number;
  created_at: string;
}

export interface DatasetAccessRequestItem {
  id: number;
  dataset_id: number;
  dataset_name?: string;
  dataset_privacy?: string;
  dataset_owner_id?: number;
  dataset_owner_name?: string;
  user_id: number;
  user_name?: string;
  reason: string;
  status: string;
  review_comment: string;
  approval_stage: number;
  required_approvals: number;
  approval_expires_at?: string;
  download_limit: number;
  download_count: number;
  remaining_downloads?: number;
  is_expired: boolean;
  authorization_active: boolean;
  sla_hours: number;
  sla_deadline_at?: string;
  sla_overdue: boolean;
  sla_remaining_minutes: number;
  reviewed_at?: string;
  created_at: string;
}

export interface BatchDatasetAccessDecisionPayload {
  ids: number[];
  decision: 'approved' | 'rejected';
  comment: string;
  valid_days: number;
  download_limit: number;
}

export interface DatasetDetail extends ResourceCard {
  review_comment?: string;
  sample_count: number;
  device: string;
  scene: string;
  privacy: string;
  agreement_text: string;
  samples: DatasetSample[];
  versions: FileVersion[];
  favorited: boolean;
}

export interface TaskTemplateItem {
  id: number;
  name: string;
  summary: string;
  description: string;
  category: string;
  scene: string;
  guide: string;
  resource_ref: string[];
  usage_count: number;
  status?: string;
  updated_at: string;
}

export interface ApplicationCaseItem {
  id: number;
  title: string;
  summary: string;
  category: string;
  guide: string;
  cover_image?: string;
  status?: string;
  updated_at: string;
}

export interface DocumentCategoryItem {
  id: number;
  name: string;
  doc_type: string;
}

export interface DocumentItem {
  id: number;
  category_id: number;
  category: string;
  title: string;
  summary: string;
  content: string;
  doc_type: string;
  status?: string;
  updated_at: string;
}

export interface FAQItem {
  id: number;
  question: string;
  answer: string;
  updated_at: string;
}

export interface VideoTutorialItem {
  id: number;
  title: string;
  summary: string;
  link: string;
  category: string;
  sort_order: number;
  active: boolean;
  updated_at: string;
}

export interface AgreementTemplateItem {
  id: number;
  name: string;
  content: string;
  sort_order: number;
  active: boolean;
  updated_at: string;
}

export interface DatasetPrivacyOptionItem {
  id: number;
  code: string;
  name: string;
  description: string;
  sort_order: number;
  active: boolean;
  updated_at: string;
}

export interface DatasetOptionsResponse {
  agreement_templates: AgreementTemplateItem[];
  privacy_options: DatasetPrivacyOptionItem[];
}

export interface FilterOptionsResponse {
  tags: string[];
  model_tags: string[];
  dataset_tags: string[];
  robot_types: string[];
  dataset_scenes: string[];
  template_categories: string[];
  template_scenes: string[];
  application_case_categories: string[];
}

export interface SearchItem {
  id: number;
  type: string;
  title: string;
  summary: string;
  tags?: string[];
  route: string;
  score_hint: number;
  updated_at: string;
}

export interface SearchTypeCountItem {
  type: string;
  label: string;
  count: number;
}

export interface SearchResponse extends Paginated<SearchItem> {
  focus_type?: string;
  type_counts: SearchTypeCountItem[];
  same_type_items: SearchItem[];
  related_items: SearchItem[];
  suggested_queries: SearchSuggestionItem[];
}

export interface HomePayload {
  platform_intro: string;
  hero_config: HomeHeroConfigItem;
  highlights: string[];
  announcements: AnnouncementItem[];
  hot_models: ResourceCard[];
  hot_datasets: ResourceCard[];
  task_templates: TaskTemplateItem[];
  application_cases: ApplicationCaseItem[];
  scene_pages: ScenePageItem[];
  module_settings: ModuleSettingItem[];
  rankings_config: RankingConfigItem;
}

export interface ScenePageDetailPayload {
  scene: ScenePageItem;
  models: ResourceCard[];
  datasets: ResourceCard[];
  task_templates: TaskTemplateItem[];
  application_cases: ApplicationCaseItem[];
}

export interface Paginated<T> {
  items: T[];
  page: number;
  page_size: number;
  total: number;
}

export interface DashboardPayload {
  users: number;
  published_models: number;
  pending_models: number;
  published_datasets: number;
  pending_datasets: number;
  announcements: number;
}

export interface ReviewItem {
  id: number;
  type: string;
  title: string;
  summary: string;
  status: string;
  owner: string;
  updated_at: string;
}

export interface ModelEvaluationItem {
  id: number;
  benchmark: string;
  summary: string;
  score: number;
  notes: string;
  user_name: string;
  created_at: string;
}

export interface RatingItem {
  id: number;
  score: number;
  feedback: string;
  user_name: string;
  created_at: string;
}

export interface RatingSummary {
  average: number;
  count: number;
  items: RatingItem[];
}

export interface CommentItem {
  id: number;
  parent_id?: number;
  content: string;
  user_id: number;
  user_name: string;
  created_at: string;
}

export interface DatasetStatsResponse {
  dataset_id: number;
  download_count: number;
  sample_count: number;
  sample_type_mix: Array<{ label: string; value: number }>;
  download_trend: Array<{ label: string; value: number }>;
}

export interface SkillItem {
  id: number;
  name: string;
  summary: string;
  description: string;
  category: string;
  scene: string;
  guide: string;
  resource_ref: string[];
  status: string;
  forked_from_id?: number;
  usage_count: number;
  owner_id: number;
  owner_name: string;
  rating: RatingSummary;
  updated_at: string;
}

export interface DiscussionItem {
  id: number;
  title: string;
  summary: string;
  content: string;
  tag: string;
  user_id: number;
  user_name: string;
  comment_count: number;
  hot_score: number;
  created_at: string;
  updated_at: string;
}

export interface FollowItem {
  id: number;
  user_id: number;
  user_name: string;
  bio: string;
  created_at: string;
}

export interface FollowStats {
  follows: number;
  followers: number;
}

export interface SearchHotItem {
  query: string;
  count: number;
}

export interface SearchSuggestionItem {
  query: string;
}

export interface SearchKeywordConfigItem {
  id: number;
  query: string;
  keyword_type: string;
  sort_order: number;
  enabled: boolean;
  updated_at: string;
}

export interface FilterOptionConfigItem {
  id: number;
  kind: string;
  value: string;
  sort_order: number;
  enabled: boolean;
  updated_at: string;
}

export interface AdminModelRecommendTagItem {
  id: number;
  name: string;
  summary: string;
  recommend_tag: string;
  status: string;
  owner_name: string;
  updated_at: string;
}

export interface UserContributionPayload {
  skills: SkillItem[];
  discussions: DiscussionItem[];
}

export interface AdminCommunityOverview {
  skills: number;
  discussions: number;
  comments: number;
}

export interface AdminConversationModerationItem {
  id: number;
  kind: string;
  title: string;
  participant_names: string[];
  latest_message: string;
  active: boolean;
  blocked_reason: string;
  updated_at: string;
}

export interface AdminWorkspaceModerationItem {
  id: number;
  name: string;
  summary: string;
  owner_id: number;
  owner_name: string;
  member_count: number;
  members: FollowItem[];
  active: boolean;
  blocked_reason: string;
  updated_at: string;
}

export interface AdminSkillModerationItem {
  id: number;
  name: string;
  summary: string;
  status: string;
  owner_id: number;
  owner_name: string;
  updated_at: string;
}

export interface AdminDiscussionModerationItem {
  id: number;
  title: string;
  tag: string;
  user_id: number;
  user_name: string;
  comment_count: number;
  updated_at: string;
}

export interface AdminCommentModerationItem {
  id: number;
  resource_type: string;
  resource_id: number;
  user_id: number;
  user_name: string;
  content: string;
  created_at: string;
}

export interface VerificationStatusItem {
  id: number;
  user_id: number;
  verification_type: string;
  real_name: string;
  organization: string;
  materials: string;
  reason: string;
  status: string;
  review_comment: string;
  reviewed_at?: string;
  created_at: string;
}

export interface AdminVerificationItem {
  id: number;
  user_id: number;
  user_name: string;
  verification_type: string;
  real_name: string;
  organization: string;
  reason: string;
  status: string;
  review_comment: string;
  reviewed_at?: string;
  created_at: string;
}

export interface AdminRewardOverview {
  benefits: number;
  active_benefits: number;
  redemptions: number;
  ledger_entries: number;
  net_points: number;
}

export interface AdminRewardAdjustmentItem {
  id: number;
  user_id: number;
  user_name: string;
  points: number;
  remark: string;
  created_at: string;
}

export interface AdminOperationLogItem {
  id: number;
  admin_user_id: number;
  admin_name: string;
  action: string;
  resource_type: string;
  resource_id: number;
  summary: string;
  detail: string;
  created_at: string;
}

export interface WikiPageItem {
  id: number;
  title: string;
  summary: string;
  content: string;
  status: string;
  editor_id: number;
  editor_name: string;
  updated_at: string;
}

export interface WikiRevisionItem {
  id: number;
  page_id: number;
  editor_id: number;
  editor_name: string;
  title: string;
  summary: string;
  content: string;
  comment: string;
  created_at: string;
}

export interface RewardSummary {
  points: number;
}

export interface RewardLedgerItem {
  id: number;
  source_type: string;
  points: number;
  remark: string;
  created_at: string;
}

export interface ContributorRankingItem {
  user_id: number;
  user_name: string;
  points: number;
}

export interface RewardBenefitItem {
  id: number;
  name: string;
  summary: string;
  cost_points: number;
  active: boolean;
}

export interface RewardRedemptionItem {
  id: number;
  benefit_id: number;
  benefit_name: string;
  cost_points: number;
  created_at: string;
}

export interface OpenAPIEndpointItem {
  method: string;
  path: string;
  summary: string;
  auth: string;
}

export interface CLIExampleItem {
  title: string;
  command: string;
}

export interface OpenAPISpecResponse {
  version: string;
  base_url: string;
  overview: string;
  endpoints: OpenAPIEndpointItem[];
  webhook_events: string[];
  curl_examples: CLIExampleItem[];
}

export interface WebhookSubscriptionItem {
  id: number;
  name: string;
  target_url: string;
  events: string[];
  last_status: string;
  last_response_code: number;
  last_error: string;
  last_delivered_at?: string;
  created_at: string;
}

export interface WebhookDeliveryItem {
  id: number;
  event: string;
  status: string;
  response_code: number;
  response_body: string;
  created_at: string;
}

export interface MessageItem {
  id: number;
  conversation_id: number;
  sender_id: number;
  sender_name: string;
  content: string;
  created_at: string;
}

export interface ConversationItem {
  id: number;
  kind: string;
  title: string;
  workspace_id?: number;
  participant_names: string[];
  latest_message: string;
  active: boolean;
  blocked_reason: string;
  updated_at: string;
}

export interface WorkspaceItem {
  id: number;
  name: string;
  summary: string;
  owner_id: number;
  conversation_id: number;
  member_count: number;
  active: boolean;
  blocked_reason: string;
  updated_at: string;
}

export interface WorkspaceDetail extends WorkspaceItem {
  messages: MessageItem[];
  members: FollowItem[];
}

export interface FavoriteRecord {
  ID: number;
  UserID: number;
  ResourceType: string;
  ResourceID: number;
  ResourceTitle: string;
  CreatedAt: string;
}

export interface DownloadRecord {
  ID: number;
  UserID: number;
  ResourceType: string;
  ResourceID: number;
  ResourceTitle: string;
  CreatedAt: string;
}

export interface NotificationRecord {
  ID: number;
  UserID: number;
  Type: string;
  Title: string;
  Content: string;
  Link?: string;
  Read: boolean;
  CreatedAt: string;
}

export interface UploadsPayload {
  models: ResourceCard[];
  datasets: ResourceCard[];
}

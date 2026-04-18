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
  enabled: boolean;
  updated_at?: string;
}

export interface FeaturedResourceItem {
  id: number;
  resource_type: string;
  resource_id: number;
  title: string;
  summary: string;
  route: string;
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
  downloads: number;
  status?: string;
  owner?: string;
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
  user_id: number;
  user_name?: string;
  reason: string;
  status: string;
  review_comment: string;
  reviewed_at?: string;
  created_at: string;
}

export interface DatasetDetail extends ResourceCard {
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

export interface HomePayload {
  platform_intro: string;
  highlights: string[];
  announcements: AnnouncementItem[];
  hot_models: ResourceCard[];
  hot_datasets: ResourceCard[];
  task_templates: TaskTemplateItem[];
  application_cases: ApplicationCaseItem[];
  module_settings: ModuleSettingItem[];
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
  category: string;
  user_id: number;
  user_name: string;
  comment_count: number;
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

export interface UserContributionPayload {
  skills: SkillItem[];
  discussions: DiscussionItem[];
}

export interface AdminCommunityOverview {
  skills: number;
  discussions: number;
  comments: number;
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
  category: string;
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
  updated_at: string;
}

export interface WorkspaceItem {
  id: number;
  name: string;
  summary: string;
  owner_id: number;
  conversation_id: number;
  member_count: number;
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
  Read: boolean;
  CreatedAt: string;
}

export interface UploadsPayload {
  models: ResourceCard[];
  datasets: ResourceCard[];
}

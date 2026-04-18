import { createRouter, createWebHistory } from 'vue-router';

const HomeView = () => import('@/views/home/HomeView.vue');
const SearchView = () => import('@/views/home/SearchView.vue');
const LoginView = () => import('@/views/auth/LoginView.vue');
const ModelListView = () => import('@/views/models/ModelListView.vue');
const ModelCompareView = () => import('@/views/models/ModelCompareView.vue');
const ModelUploadView = () => import('@/views/models/ModelUploadView.vue');
const ModelDetailView = () => import('@/views/models/ModelDetailView.vue');
const DatasetListView = () => import('@/views/datasets/DatasetListView.vue');
const DatasetUploadView = () => import('@/views/datasets/DatasetUploadView.vue');
const DatasetDetailView = () => import('@/views/datasets/DatasetDetailView.vue');
const TemplateListView = () => import('@/views/templates/TemplateListView.vue');
const TemplateDetailView = () => import('@/views/templates/TemplateDetailView.vue');
const ApplicationCaseListView = () => import('@/views/templates/ApplicationCaseListView.vue');
const ApplicationCaseDetailView = () => import('@/views/templates/ApplicationCaseDetailView.vue');
const SkillListView = () => import('@/views/skills/SkillListView.vue');
const SkillCreateView = () => import('@/views/skills/SkillCreateView.vue');
const SkillDetailView = () => import('@/views/skills/SkillDetailView.vue');
const DiscussionListView = () => import('@/views/community/DiscussionListView.vue');
const DiscussionDetailView = () => import('@/views/community/DiscussionDetailView.vue');
const DocsView = () => import('@/views/docs/DocsView.vue');
const WikiView = () => import('@/views/docs/WikiView.vue');
const DeveloperHubView = () => import('@/views/docs/DeveloperHubView.vue');
const ProfileView = () => import('@/views/profile/ProfileView.vue');
const PublicUserView = () => import('@/views/community/PublicUserView.vue');
const MessagesView = () => import('@/views/community/MessagesView.vue');
const WorkspacesView = () => import('@/views/community/WorkspacesView.vue');
const AdminDashboardView = () => import('@/views/admin/AdminDashboardView.vue');
const AdminPortalView = () => import('@/views/admin/AdminPortalView.vue');
const AdminDatasetAccessView = () => import('@/views/admin/AdminDatasetAccessView.vue');
const AdminReviewsView = () => import('@/views/admin/AdminReviewsView.vue');
const AdminAnnouncementsView = () => import('@/views/admin/AdminAnnouncementsView.vue');
const AdminContentView = () => import('@/views/admin/AdminContentView.vue');
const AdminCommunityView = () => import('@/views/admin/AdminCommunityView.vue');
const AdminVerificationView = () => import('@/views/admin/AdminVerificationView.vue');
const AdminRewardsView = () => import('@/views/admin/AdminRewardsView.vue');

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/search', name: 'search', component: SearchView },
    { path: '/login', name: 'login', component: LoginView },
    { path: '/models', name: 'models', component: ModelListView },
    { path: '/models/compare', name: 'model-compare', component: ModelCompareView },
    { path: '/models/upload', name: 'model-upload', component: ModelUploadView, meta: { requiresAuth: true } },
    { path: '/models/:id', name: 'model-detail', component: ModelDetailView, props: true },
    { path: '/datasets', name: 'datasets', component: DatasetListView },
    { path: '/datasets/upload', name: 'dataset-upload', component: DatasetUploadView, meta: { requiresAuth: true } },
    { path: '/datasets/:id', name: 'dataset-detail', component: DatasetDetailView, props: true },
    { path: '/templates', name: 'templates', component: TemplateListView },
    { path: '/templates/:id', name: 'template-detail', component: TemplateDetailView, props: true },
    { path: '/applications', name: 'applications', component: ApplicationCaseListView },
    { path: '/applications/:id', name: 'application-detail', component: ApplicationCaseDetailView, props: true },
    { path: '/skills', name: 'skills', component: SkillListView },
    { path: '/skills/new', name: 'skill-new', component: SkillCreateView, meta: { requiresAuth: true } },
    { path: '/skills/:id', name: 'skill-detail', component: SkillDetailView, props: true },
    { path: '/discussions', name: 'discussions', component: DiscussionListView },
    { path: '/discussions/:id', name: 'discussion-detail', component: DiscussionDetailView, props: true },
    { path: '/docs/:id?', name: 'docs', component: DocsView, props: true },
    { path: '/wiki/:id?', name: 'wiki', component: WikiView, props: true },
    { path: '/developers', name: 'developers', component: DeveloperHubView },
    { path: '/me', name: 'me', component: ProfileView, meta: { requiresAuth: true } },
    { path: '/community/users/:id', name: 'public-user', component: PublicUserView, props: true },
    { path: '/messages', name: 'messages', component: MessagesView, meta: { requiresAuth: true } },
    { path: '/workspaces', name: 'workspaces', component: WorkspacesView, meta: { requiresAuth: true } },
    { path: '/admin', name: 'admin-dashboard', component: AdminDashboardView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/portal', name: 'admin-portal', component: AdminPortalView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/dataset-access', name: 'admin-dataset-access', component: AdminDatasetAccessView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/reviews', name: 'admin-reviews', component: AdminReviewsView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/announcements', name: 'admin-announcements', component: AdminAnnouncementsView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/content', name: 'admin-content', component: AdminContentView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/community', name: 'admin-community', component: AdminCommunityView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/verifications', name: 'admin-verifications', component: AdminVerificationView, meta: { requiresAuth: true, requiresAdmin: true } },
    { path: '/admin/rewards', name: 'admin-rewards', component: AdminRewardsView, meta: { requiresAuth: true, requiresAdmin: true } },
  ],
  scrollBehavior() {
    return { top: 0 };
  },
});

router.beforeEach((to) => {
  const token = localStorage.getItem('open-community-token');
  const user = localStorage.getItem('open-community-user');
  const parsedUser = user ? (JSON.parse(user) as { role?: string }) : null;

  if (to.meta.requiresAuth && !token) {
    return { name: 'login', query: { redirect: to.fullPath } };
  }
  if (to.meta.requiresAdmin && parsedUser?.role !== 'admin') {
    return { name: 'home' };
  }
  return true;
});

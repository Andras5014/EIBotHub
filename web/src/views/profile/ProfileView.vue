<template>
  <div class="page-shell">
    <div class="section-head page-topbar">
      <div>
        <div class="page-kicker">Personal Profile</div>
        <h1 class="section-title">个人中心</h1>
        <p class="section-subtitle">参考 GitHub 用户主页的信息层级，把身份、贡献、资源和管理操作集中到一张结构清晰的个人主页中。</p>
      </div>
      <a-space wrap class="page-topbar-actions">
        <a-badge :count="unreadNotificationCount" :offset="[8, 2]">
          <a-button @click="openNotificationCenter">通知中心</a-button>
        </a-badge>
        <a-button type="primary" @click="openProfileSettings()">个人设置</a-button>
      </a-space>
    </div>

    <div class="profile-layout">
      <div class="profile-aside">
        <UserProfileSidebar
          :name="profile.username || auth.user?.username || '未命名用户'"
          :handle="profile.username || auth.user?.username || 'user'"
          :bio="profile.bio"
          :avatar="profile.avatar"
          :badges="profileBadges"
          :stats="profileStats"
        >
          <template #actions>
            <a-space direction="vertical" style="width: 100%">
              <a-button type="primary" @click="openPublicProfile">查看公开主页</a-button>
              <a-button @click="openProfileSettings()">打开设置面板</a-button>
            </a-space>
          </template>

          <div class="meta-line">账号角色：{{ auth.user?.role === 'admin' ? '管理员' : '开发者' }}</div>
          <div class="meta-line">开发者认证：{{ verification ? statusLabel(verification.status) : '未申请' }}</div>
        </UserProfileSidebar>

      </div>

      <div class="profile-main">
        <section class="page-card hero-panel">
          <div class="hero-grid">
            <div class="hero-copy">
              <div class="hero-kicker">Dashboard Snapshot</div>
              <h2 class="hero-title">把“我是谁、我做了什么、我还要处理什么”放到同一视图里</h2>
              <p class="hero-subtitle">
                GitHub 用户主页的价值不只是看资料，而是把身份、代表性内容和最近活动压缩进同一条信息流。这里沿用这个思路来管理你的资源、通知和社区关系。
              </p>
            </div>
            <div class="hero-stats">
              <div class="hero-stat-card">
                <div class="hero-stat-label">积分</div>
                <div class="hero-stat-value">{{ rewardSummary.points }}</div>
              </div>
              <div class="hero-stat-card">
                <div class="hero-stat-label">技能</div>
                <div class="hero-stat-value">{{ contributions.skills.length }}</div>
              </div>
              <div class="hero-stat-card">
                <div class="hero-stat-label">讨论</div>
                <div class="hero-stat-value">{{ contributions.discussions.length }}</div>
              </div>
              <div class="hero-stat-card">
                <div class="hero-stat-label">粉丝</div>
                <div class="hero-stat-value">{{ followStats.followers }}</div>
              </div>
            </div>
          </div>
        </section>

        <section class="page-card main-card">
          <div class="section-head">
            <div>
              <h2 class="section-title">我的空间</h2>
              <p class="section-subtitle">保留当前所有管理能力，但按 GitHub 个人主页的浏览顺序重新整理内容层级。</p>
            </div>
          </div>

          <a-tabs v-model:activeKey="mainTab">
            <a-tab-pane key="overview" tab="概览">
              <div class="overview-grid">
                <a-card title="贡献概览" class="inner-card">
                  <a-row :gutter="[12, 12]">
                    <a-col :span="12"><a-statistic title="关注中" :value="followStats.follows" /></a-col>
                    <a-col :span="12"><a-statistic title="粉丝" :value="followStats.followers" /></a-col>
                    <a-col :span="12"><a-statistic title="已兑换权益" :value="rewardRedemptions.length" /></a-col>
                    <a-col :span="12"><a-statistic title="访问申请" :value="accessRequests.length" /></a-col>
                  </a-row>
                </a-card>

                <a-card title="最近技能" class="inner-card">
                  <a-empty v-if="!contributions.skills.length" description="还没有公开技能" />
                  <a-list v-else :data-source="contributions.skills.slice(0, 5)">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <a-list-item-meta :title="item.name" :description="item.summary" />
                      </a-list-item>
                    </template>
                  </a-list>
                </a-card>

                <a-card title="最近讨论" class="inner-card">
                  <a-empty v-if="!contributions.discussions.length" description="还没有公开讨论" />
                  <a-list v-else :data-source="contributions.discussions.slice(0, 5)">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <a-list-item-meta :title="item.title" :description="item.summary" />
                      </a-list-item>
                    </template>
                  </a-list>
                </a-card>
              </div>
            </a-tab-pane>

            <a-tab-pane key="resources" tab="资源与记录">
              <a-tabs size="small">
                <a-tab-pane key="models" tab="我的模型">
                  <a-space style="margin-bottom: 12px" wrap>
                    <span class="section-subtitle">状态筛选</span>
                    <a-select v-model:value="modelStatusFilter" style="width: 180px" :options="statusOptions" />
                  </a-space>
                  <a-empty v-if="!filteredModels.length" description="当前筛选条件下没有模型" />
                  <a-list v-else :data-source="filteredModels">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <a-list-item-meta>
                          <template #title>
                            <a-space wrap>
                              <RouterLink :to="`/models/${item.id}`">{{ item.name }}</RouterLink>
                              <a-tag :color="statusColor(item.status)">{{ statusLabel(item.status) }}</a-tag>
                              <span class="section-subtitle">下载 {{ item.downloads }}</span>
                            </a-space>
                          </template>
                          <template #description>
                            <div class="resource-summary">{{ item.summary || '暂无摘要' }}</div>
                            <div v-if="item.review_comment && item.status === 'rejected'" class="review-comment">驳回原因：{{ item.review_comment }}</div>
                          </template>
                        </a-list-item-meta>
                        <template #actions>
                          <RouterLink :to="`/models/${item.id}`">详情</RouterLink>
                          <RouterLink :to="`/models/${item.id}/edit`">编辑</RouterLink>
                          <a-button
                            v-if="canResubmit(item.status)"
                            type="link"
                            :loading="resubmittingKey === `model-${item.id}`"
                            @click="resubmitModel(item.id)"
                          >
                            重新提交审核
                          </a-button>
                        </template>
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="datasets" tab="我的数据集">
                  <a-space style="margin-bottom: 12px" wrap>
                    <span class="section-subtitle">状态筛选</span>
                    <a-select v-model:value="datasetStatusFilter" style="width: 180px" :options="statusOptions" />
                  </a-space>
                  <a-empty v-if="!filteredDatasets.length" description="当前筛选条件下没有数据集" />
                  <a-list v-else :data-source="filteredDatasets">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <a-list-item-meta>
                          <template #title>
                            <a-space wrap>
                              <RouterLink :to="`/datasets/${item.id}`">{{ item.name }}</RouterLink>
                              <a-tag :color="statusColor(item.status)">{{ statusLabel(item.status) }}</a-tag>
                              <span class="section-subtitle">下载 {{ item.downloads }}</span>
                            </a-space>
                          </template>
                          <template #description>
                            <div class="resource-summary">{{ item.summary || '暂无摘要' }}</div>
                            <div v-if="item.review_comment && item.status === 'rejected'" class="review-comment">驳回原因：{{ item.review_comment }}</div>
                          </template>
                        </a-list-item-meta>
                        <template #actions>
                          <RouterLink :to="`/datasets/${item.id}`">详情</RouterLink>
                          <RouterLink :to="`/datasets/${item.id}/edit`">编辑</RouterLink>
                          <a-button
                            v-if="canResubmit(item.status)"
                            type="link"
                            :loading="resubmittingKey === `dataset-${item.id}`"
                            @click="resubmitDataset(item.id)"
                          >
                            重新提交审核
                          </a-button>
                        </template>
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="favorites" tab="收藏记录">
                  <a-list :data-source="favorites">
                    <template #renderItem="{ item }">
                      <a-list-item>{{ item.ResourceTitle }} · {{ item.ResourceType }}</a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="downloads" tab="下载历史">
                  <a-list :data-source="downloads">
                    <template #renderItem="{ item }">
                      <a-list-item>{{ item.ResourceTitle }} · {{ item.ResourceType }}</a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="access-requests" tab="访问申请">
                  <a-space style="margin-bottom: 12px" wrap>
                    <span class="section-subtitle">状态筛选</span>
                    <a-select v-model:value="accessRequestStatusFilter" style="width: 180px" :options="accessRequestStatusOptions" />
                  </a-space>
                  <a-empty v-if="!filteredAccessRequests.length" description="当前筛选条件下没有访问申请记录" />
                  <a-list v-else :data-source="filteredAccessRequests">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <a-list-item-meta>
                          <template #title>
                            <a-space wrap>
                              <RouterLink :to="`/datasets/${item.dataset_id}`">{{ item.dataset_name || `数据集 ${item.dataset_id}` }}</RouterLink>
                              <a-tag v-if="item.dataset_privacy" color="purple">{{ item.dataset_privacy }}</a-tag>
                              <a-tag :color="statusColor(item.status)">{{ statusLabel(item.status) }}</a-tag>
                            </a-space>
                          </template>
                          <template #description>
                            <div class="resource-summary">{{ item.reason }}</div>
                            <div v-if="item.required_approvals > 1" class="section-subtitle">审批进度：{{ item.approval_stage }}/{{ item.required_approvals }}</div>
                            <div v-if="item.sla_hours > 0" class="section-subtitle">处理时限：{{ slaSummary(item) }}</div>
                            <div class="section-subtitle">申请时间：{{ formatDate(item.created_at) }}</div>
                            <div v-if="item.reviewed_at" class="section-subtitle">处理时间：{{ formatDate(item.reviewed_at) }}</div>
                            <div v-if="item.status === 'approved'" class="section-subtitle">授权信息：{{ accessGrantSummary(item) }}</div>
                            <div v-if="item.review_comment" class="review-comment">审核意见：{{ item.review_comment }}</div>
                          </template>
                        </a-list-item-meta>
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>
              </a-tabs>
            </a-tab-pane>

            <a-tab-pane key="community" tab="社区与通知">
              <a-tabs v-model:activeKey="communityTab" size="small">
                <a-tab-pane key="skills" tab="我的技能">
                  <a-list :data-source="contributions.skills">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <RouterLink :to="`/skills/${item.id}`">{{ item.name }}</RouterLink>
                        <span class="section-subtitle">{{ item.summary }}</span>
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="discussions" tab="我的讨论">
                  <a-list :data-source="contributions.discussions">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <RouterLink :to="`/discussions/${item.id}`">{{ item.title }}</RouterLink>
                        <span class="section-subtitle">{{ item.summary }}</span>
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="notifications">
                  <template #tab>
                    <a-badge :count="unreadNotificationCount" :offset="[8, 0]">
                      <span>通知中心</span>
                    </a-badge>
                  </template>
                  <a-space style="margin-bottom: 12px">
                    <a-button :disabled="!unreadNotificationCount" @click="markNotificationsRead">全部标记已读</a-button>
                    <span class="section-subtitle">未读 {{ unreadNotificationCount }} 条</span>
                  </a-space>
                  <a-empty v-if="!notifications.length" description="当前没有通知" />
                  <a-list v-else :data-source="notifications">
                    <template #renderItem="{ item }">
                      <a-list-item
                        :class="['notification-item', { unread: !item.Read, clickable: Boolean(item.Link) }]"
                        @click="openNotificationItem(item)"
                      >
                        <a-list-item-meta>
                          <template #title>
                            <div class="notification-title-row">
                              <div class="notification-title-main">
                                <span v-if="!item.Read" class="notification-dot" />
                                <span>{{ item.Title }}</span>
                              </div>
                              <a-tag :color="item.Read ? 'default' : 'red'">{{ item.Read ? '已读' : '未读' }}</a-tag>
                            </div>
                          </template>
                          <template #description>
                            <div class="notification-content">{{ item.Content }}</div>
                            <div class="notification-meta">
                              <span>{{ notificationTypeLabel(item.Type) }}</span>
                              <span>{{ formatDate(item.CreatedAt) }}</span>
                            </div>
                          </template>
                        </a-list-item-meta>
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="follows" tab="我的关注">
                  <a-list :data-source="follows">
                    <template #renderItem="{ item }">
                      <a-list-item>{{ item.user_name }} · {{ item.bio }}</a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="followers" tab="我的粉丝">
                  <a-list :data-source="followers">
                    <template #renderItem="{ item }">
                      <a-list-item>{{ item.user_name }} · {{ item.bio }}</a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="rankings" tab="排行榜">
                  <a-list :data-source="rankings">
                    <template #renderItem="{ item, index }">
                      <a-list-item>{{ index + 1 }}. {{ item.user_name }} · {{ item.points }}</a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>
              </a-tabs>
            </a-tab-pane>

            <a-tab-pane key="rewards" tab="积分与权益">
              <a-tabs size="small">
                <a-tab-pane key="rewards-ledger" tab="积分账本">
                  <a-list :data-source="rewardLedger">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <a-list-item-meta :title="`${item.points > 0 ? '+' : ''}${item.points} · ${item.remark}`" :description="item.source_type" />
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>

                <a-tab-pane key="benefits" tab="权益兑换">
                  <a-row :gutter="[16, 16]">
                    <a-col v-for="item in rewardBenefits" :key="item.id" :xs="24" :md="12">
                      <a-card class="redeem-card">
                        <div class="section-head redeem-head">
                          <div>
                            <div class="list-title">{{ item.name }}</div>
                            <div class="list-desc">{{ item.summary }}</div>
                          </div>
                          <a-tag color="gold">{{ item.cost_points }} 积分</a-tag>
                        </div>
                        <a-space>
                          <a-tag :color="rewardSummary.points >= item.cost_points ? 'green' : 'default'">
                            {{ rewardSummary.points >= item.cost_points ? '积分充足' : '积分不足' }}
                          </a-tag>
                          <a-button
                            type="primary"
                            :loading="redeemingId === item.id"
                            :disabled="rewardSummary.points < item.cost_points || !item.active"
                            @click="redeem(item.id)"
                          >
                            立即兑换
                          </a-button>
                        </a-space>
                      </a-card>
                    </a-col>
                  </a-row>
                </a-tab-pane>

                <a-tab-pane key="redemptions" tab="兑换记录">
                  <a-list :data-source="rewardRedemptions">
                    <template #renderItem="{ item }">
                      <a-list-item>
                        <a-list-item-meta
                          :title="`${item.benefit_name} · -${item.cost_points} 积分`"
                          :description="new Date(item.created_at).toLocaleString()"
                        />
                      </a-list-item>
                    </template>
                  </a-list>
                </a-tab-pane>
              </a-tabs>
            </a-tab-pane>
          </a-tabs>
        </section>
      </div>
    </div>

    <a-modal
      v-model:open="profileSettingsOpen"
      title="个人设置"
      :footer="null"
      width="760px"
      destroy-on-close
      @cancel="closeProfileSettings"
    >
      <a-tabs v-model:activeKey="profileSettingsTab">
        <a-tab-pane key="profile" tab="编辑资料">
          <a-form :model="profile" layout="vertical" @finish="saveProfile">
            <a-form-item label="用户名" name="username"><a-input v-model:value="profile.username" /></a-form-item>
            <a-form-item label="头像链接" name="avatar"><a-input v-model:value="profile.avatar" /></a-form-item>
            <a-form-item label="个人简介" name="bio"><a-textarea v-model:value="profile.bio" :rows="4" /></a-form-item>
            <a-form-item style="margin-bottom: 0">
              <a-button type="primary" html-type="submit">保存资料</a-button>
            </a-form-item>
          </a-form>
        </a-tab-pane>

        <a-tab-pane key="verification" tab="提交认证">
          <template v-if="verification">
            <a-space wrap>
              <a-tag color="cyan">{{ verification.verification_type }}</a-tag>
              <a-tag :color="verification.status === 'approved' ? 'green' : verification.status === 'rejected' ? 'red' : 'gold'">
                {{ statusLabel(verification.status) }}
              </a-tag>
            </a-space>
            <p class="section-subtitle modal-section-copy">
              {{ verification.real_name }} {{ verification.organization ? `· ${verification.organization}` : '' }}
            </p>
            <p class="section-subtitle">{{ verification.review_comment || '申请已提交，等待审核。' }}</p>
          </template>

          <a-form :model="verificationForm" layout="vertical" class="verification-form" @finish="submitVerification">
            <a-form-item label="认证类型">
              <a-radio-group v-model:value="verificationForm.verification_type">
                <a-radio-button value="personal">个人</a-radio-button>
                <a-radio-button value="enterprise">企业</a-radio-button>
              </a-radio-group>
            </a-form-item>
            <a-form-item label="真实姓名"><a-input v-model:value="verificationForm.real_name" /></a-form-item>
            <a-form-item label="组织名称"><a-input v-model:value="verificationForm.organization" /></a-form-item>
            <a-form-item label="材料说明"><a-textarea v-model:value="verificationForm.materials" :rows="2" /></a-form-item>
            <a-form-item label="申请理由"><a-textarea v-model:value="verificationForm.reason" :rows="3" /></a-form-item>
            <a-form-item style="margin-bottom: 0"><a-button type="primary" html-type="submit">提交认证申请</a-button></a-form-item>
          </a-form>
        </a-tab-pane>

        <a-tab-pane key="settings" tab="设置">
          <div class="settings-grid">
            <a-card class="inner-card settings-card" title="提醒">
              <p class="section-subtitle">当前还有 {{ unreadNotificationCount }} 条未读提醒。</p>
              <a-space wrap>
                <a-button @click="openNotificationCenter">查看通知中心</a-button>
                <a-button :disabled="!unreadNotificationCount" @click="markNotificationsRead">全部标记已读</a-button>
              </a-space>
            </a-card>
            <a-card class="inner-card settings-card" title="账号">
              <p class="section-subtitle">你当前的身份是 {{ auth.user?.role === 'admin' ? '管理员' : '开发者' }}。</p>
              <a-space wrap>
                <a-button type="primary" @click="openPublicProfile">查看公开主页</a-button>
                <a-button danger @click="logout">退出登录</a-button>
              </a-space>
            </a-card>
          </div>
        </a-tab-pane>
      </a-tabs>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { message } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import UserProfileSidebar from '@/components/UserProfileSidebar.vue';
import { useAuthStore } from '@/stores/auth';
import type {
  ContributorRankingItem,
  DatasetAccessRequestItem,
  DownloadRecord,
  FavoriteRecord,
  FollowItem,
  FollowStats,
  NotificationRecord,
  ResourceCard,
  RewardBenefitItem,
  RewardLedgerItem,
  RewardRedemptionItem,
  RewardSummary,
  UserContributionPayload,
  VerificationStatusItem,
} from '@/types/api';

const auth = useAuthStore();
const route = useRoute();
const router = useRouter();
const mainTab = ref('overview');
const communityTab = ref('skills');
const profileSettingsOpen = ref(false);
const profileSettingsTab = ref<'profile' | 'verification' | 'settings'>('profile');

const profile = reactive({
  username: '',
  avatar: '',
  bio: '',
});

const uploads = reactive<{ models: ResourceCard[]; datasets: ResourceCard[] }>({
  models: [],
  datasets: [],
});
const favorites = ref<FavoriteRecord[]>([]);
const downloads = ref<DownloadRecord[]>([]);
const accessRequests = ref<DatasetAccessRequestItem[]>([]);
const notifications = ref<NotificationRecord[]>([]);
const rewardSummary = ref<RewardSummary>({ points: 0 });
const rewardLedger = ref<RewardLedgerItem[]>([]);
const rewardBenefits = ref<RewardBenefitItem[]>([]);
const rewardRedemptions = ref<RewardRedemptionItem[]>([]);
const rankings = ref<ContributorRankingItem[]>([]);
const follows = ref<FollowItem[]>([]);
const followers = ref<FollowItem[]>([]);
const modelStatusFilter = ref('all');
const datasetStatusFilter = ref('all');
const accessRequestStatusFilter = ref('all');
const resubmittingKey = ref('');
const redeemingId = ref<number>();
const followStats = ref<FollowStats>({
  follows: 0,
  followers: 0,
});
const verification = ref<VerificationStatusItem | null>(null);
const contributions = reactive<UserContributionPayload>({
  skills: [],
  discussions: [],
});
const verificationForm = reactive({
  verification_type: 'personal' as 'personal' | 'enterprise',
  real_name: '',
  organization: '',
  materials: '',
  reason: '',
});

const statusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '草稿', value: 'draft' },
  { label: '待审核', value: 'pending' },
  { label: '已发布', value: 'published' },
  { label: '已驳回', value: 'rejected' },
];
const accessRequestStatusOptions = [
  { label: '全部状态', value: 'all' },
  { label: '待审核', value: 'pending' },
  { label: '已通过', value: 'approved' },
  { label: '已驳回', value: 'rejected' },
];

const filteredModels = computed(() => filterByStatus(uploads.models, modelStatusFilter.value));
const filteredDatasets = computed(() => filterByStatus(uploads.datasets, datasetStatusFilter.value));
const filteredAccessRequests = computed(() => filterByStatus(accessRequests.value, accessRequestStatusFilter.value));
const unreadNotificationCount = computed(() => notifications.value.filter((item) => !item.Read).length);
const profileBadges = computed(() => {
  const badges = [
    { text: auth.user?.role === 'admin' ? '管理员' : '开发者', color: auth.user?.role === 'admin' ? 'volcano' : 'blue' },
  ];
  if (verification.value) {
    badges.push({ text: verification.value.verification_type, color: 'cyan' });
    badges.push({
      text: verification.value.status === 'approved' ? '认证通过' : verification.value.status === 'rejected' ? '认证驳回' : '认证审核中',
      color: verification.value.status === 'approved' ? 'green' : verification.value.status === 'rejected' ? 'red' : 'gold',
    });
  }
  return badges;
});
const profileStats = computed(() => [
  { label: '积分', value: rewardSummary.value.points },
  { label: '技能', value: contributions.skills.length },
  { label: '讨论', value: contributions.discussions.length },
  { label: '粉丝', value: followStats.value.followers },
]);

async function load() {
  const profileData = await api.getProfile();
  profile.username = profileData.username;
  profile.avatar = profileData.avatar;
  profile.bio = profileData.bio;

  const uploadData = await api.getUploads();
  uploads.models = uploadData.models;
  uploads.datasets = uploadData.datasets;
  favorites.value = await api.getFavorites();
  downloads.value = await api.getDownloads();
  accessRequests.value = await api.getMyDatasetAccessRequests();
  notifications.value = await api.getNotifications();
  rewardSummary.value = await api.myRewardSummary();
  rewardLedger.value = await api.myRewardLedger();
  rewardBenefits.value = await api.rewardBenefits();
  rewardRedemptions.value = await api.myRewardRedemptions();
  rankings.value = await api.contributorRankings();
  follows.value = await api.myFollows();
  followers.value = await api.myFollowers();
  followStats.value = await api.myFollowStats();
  verification.value = await api.myVerification();
  const contributionData = await api.myContributions();
  contributions.skills = contributionData.skills;
  contributions.discussions = contributionData.discussions;
}

async function saveProfile() {
  const updated = await api.updateProfile(profile);
  auth.setAuth(auth.token, updated);
  message.success('资料已保存');
}

async function markNotificationsRead() {
  await api.readNotifications();
  await load();
  message.success('通知已标记已读');
}

async function redeem(benefitId: number) {
  redeemingId.value = benefitId;
  try {
    const item = await api.redeemRewardBenefit({ benefit_id: benefitId });
    await load();
    message.success(`已兑换：${item.benefit_name}`);
  } catch (error) {
    message.error(error instanceof Error ? error.message : '兑换失败');
  } finally {
    redeemingId.value = undefined;
  }
}

async function resubmitModel(id: number) {
  resubmittingKey.value = `model-${id}`;
  try {
    await api.submitModel(id);
    await load();
    message.success('模型已重新提交审核');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '模型重新提交失败');
  } finally {
    resubmittingKey.value = '';
  }
}

async function resubmitDataset(id: number) {
  resubmittingKey.value = `dataset-${id}`;
  try {
    await api.submitDataset(id);
    await load();
    message.success('数据集已重新提交审核');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '数据集重新提交失败');
  } finally {
    resubmittingKey.value = '';
  }
}

async function submitVerification() {
  try {
    verification.value =
      verificationForm.verification_type === 'enterprise'
        ? await api.applyEnterpriseVerification({
            real_name: verificationForm.real_name,
            organization: verificationForm.organization,
            materials: verificationForm.materials,
            reason: verificationForm.reason,
          })
        : await api.applyVerification(verificationForm);
    profileSettingsOpen.value = false;
    message.success('认证申请已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '认证申请失败');
  }
}

async function logout() {
  await auth.logout();
  router.push('/');
}

function openPublicProfile() {
  if (!auth.user?.id) return;
  router.push(`/community/users/${auth.user.id}`);
}

function openProfileSettings(tab: 'profile' | 'verification' | 'settings' = 'profile') {
  profileSettingsTab.value = tab;
  profileSettingsOpen.value = true;
}

function closeProfileSettings() {
  profileSettingsOpen.value = false;
  if (route.query.settings || route.query.settingsTab) {
    const nextQuery = { ...route.query };
    delete nextQuery.settings;
    delete nextQuery.settingsTab;
    void router.replace({ path: '/me', query: nextQuery });
  }
}

function openNotificationCenter() {
  mainTab.value = 'community';
  communityTab.value = 'notifications';
  void router.replace({ path: '/me', query: { tab: 'community', subtab: 'notifications' } });
}

function filterByStatus<T extends { status?: string }>(items: T[], status: string) {
  if (status === 'all') {
    return items;
  }
  return items.filter((item) => item.status === status);
}

function statusLabel(status?: string) {
  return {
    draft: '草稿',
    pending: '待审核',
    approved: '已通过',
    published: '已发布',
    rejected: '已驳回',
  }[status ?? ''] ?? status;
}

function statusColor(status?: string) {
  return {
    draft: 'default',
    pending: 'gold',
    approved: 'green',
    published: 'green',
    rejected: 'red',
  }[status ?? ''] ?? 'default';
}

function canResubmit(status?: string) {
  return status === 'draft' || status === 'rejected';
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('zh-CN');
}

function accessGrantSummary(item: DatasetAccessRequestItem) {
  const limitText = item.download_limit > 0 ? `剩余 ${item.remaining_downloads ?? 0} 次` : '不限次数';
  const expireText = item.approval_expires_at ? formatDate(item.approval_expires_at) : '长期有效';
  const statusText = item.authorization_active ? '授权有效' : item.is_expired ? '授权已过期' : '次数已用尽';
  return `${statusText} · ${limitText} · 有效期至 ${expireText}`;
}

function slaSummary(item: DatasetAccessRequestItem) {
  if (!item.sla_hours || !item.sla_deadline_at) {
    return '无 SLA 要求';
  }
  return item.sla_overdue
    ? `SLA ${item.sla_hours} 小时，已超时 ${formatMinutes(-item.sla_remaining_minutes)}`
    : `SLA ${item.sla_hours} 小时，剩余 ${formatMinutes(item.sla_remaining_minutes)}`;
}

function formatMinutes(value: number) {
  const minutes = Math.max(value, 0);
  const hours = Math.floor(minutes / 60);
  const rest = minutes % 60;
  if (hours > 0) {
    return `${hours} 小时 ${rest} 分`;
  }
  return `${rest} 分`;
}

function notificationTypeLabel(value: string) {
  return {
    follow: '关注提醒',
    rating: '评分提醒',
    comment: '评论提醒',
    resource: '资源提醒',
    system: '系统通知',
    dataset_access: '访问审批',
    dataset_access_review: '审批结果',
    reward: '积分提醒',
    reward_adjustment: '积分变更',
    verification: '认证提醒',
    message: '私信提醒',
    workspace: '协作提醒',
    model_evaluation: '评测提醒',
  }[value] ?? value;
}

function openNotificationItem(item: NotificationRecord) {
  if (item.Link) {
    void router.push(item.Link);
    return;
  }
  if (item.Type === 'message') {
    void router.push('/messages');
    return;
  }
}

watch(
  () => [route.query.tab, route.query.subtab],
  ([tab, subtab]) => {
    const nextMain = typeof tab === 'string' ? tab : 'overview';
    const nextCommunity = typeof subtab === 'string' ? subtab : 'skills';
    if (['overview', 'resources', 'community', 'rewards'].includes(nextMain)) {
      mainTab.value = nextMain;
    }
    if (['skills', 'discussions', 'notifications', 'follows', 'followers', 'rankings'].includes(nextCommunity)) {
      communityTab.value = nextCommunity;
    }
  },
  { immediate: true },
);

watch(
  () => [route.query.settings, route.query.settingsTab],
  ([settings, settingsTab]) => {
    if (settings !== '1') {
      return;
    }
    const nextTab =
      settingsTab === 'verification' || settingsTab === 'settings' || settingsTab === 'profile'
        ? settingsTab
        : 'profile';
    profileSettingsTab.value = nextTab;
    profileSettingsOpen.value = true;
  },
  { immediate: true },
);

onMounted(load);
</script>

<style scoped>
.page-topbar {
  margin-bottom: 20px;
}

.page-topbar-actions {
  flex: none;
}

.page-kicker {
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  margin-bottom: 8px;
}

.profile-layout {
  display: grid;
  gap: 20px;
  grid-template-columns: minmax(280px, 320px) minmax(0, 1fr);
}

.inner-card {
  border-radius: 18px;
}

.verification-form {
  margin-top: 12px;
}

.hero-panel,
.main-card {
  padding: 24px;
  border-radius: 24px;
}

.hero-panel {
  background:
    radial-gradient(circle at top right, rgba(22, 119, 255, 0.16), transparent 32%),
    linear-gradient(145deg, #f8fbff 0%, #ffffff 58%, #f2fbf7 100%);
}

.hero-grid {
  display: grid;
  gap: 20px;
  grid-template-columns: minmax(0, 1.3fr) minmax(280px, 360px);
  align-items: start;
}

.hero-kicker {
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.hero-title {
  margin: 12px 0 0;
  font-size: 30px;
  line-height: 1.2;
}

.hero-subtitle {
  margin-top: 12px;
  color: var(--text-secondary);
  line-height: 1.8;
}

.hero-stats {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.hero-stat-card {
  padding: 18px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(220, 230, 242, 0.82);
}

.hero-stat-label {
  color: var(--text-secondary);
  font-size: 12px;
}

.hero-stat-value {
  margin-top: 8px;
  color: var(--text-main);
  font-size: 28px;
  font-weight: 700;
}

.main-card {
  margin-top: 20px;
}

.overview-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.redeem-card {
  border-radius: 16px;
}

.redeem-head {
  margin-bottom: 16px;
}

.resource-summary {
  color: var(--text-secondary);
}

.review-comment {
  margin-top: 8px;
  color: #b42318;
}

.list-title {
  color: var(--text-main);
  font-weight: 700;
}

.list-desc {
  margin-top: 6px;
  color: var(--text-secondary);
}

.meta-line {
  margin-top: 4px;
}

.modal-section-copy {
  margin-top: 10px;
}

.settings-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.settings-card {
  height: 100%;
}

.notification-item {
  border: 1px solid rgba(220, 230, 242, 0.82);
  border-radius: 16px;
  margin-bottom: 12px;
  padding: 14px 16px;
  background: #fff;
  cursor: default;
}

.notification-item.unread {
  border-color: rgba(255, 77, 79, 0.28);
  background: linear-gradient(180deg, #fff8f8 0%, #ffffff 100%);
}

.notification-item.clickable {
  cursor: pointer;
}

.notification-item.clickable:hover {
  border-color: rgba(22, 119, 255, 0.28);
  background: linear-gradient(180deg, #f8fbff 0%, #ffffff 100%);
}

.notification-title-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.notification-title-main {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  color: var(--text-main);
  font-weight: 700;
}

.notification-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #ff4d4f;
  flex: none;
}

.notification-content {
  color: var(--text-main);
  line-height: 1.75;
}

.notification-meta {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  margin-top: 8px;
  color: var(--text-secondary);
  font-size: 12px;
}

@media (max-width: 1080px) {
  .profile-layout,
  .hero-grid,
  .overview-grid,
  .settings-grid {
    grid-template-columns: 1fr;
  }
}
</style>

<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">个人中心</h1>
          <p class="section-subtitle">管理个人资料、上传资源、收藏、下载历史、通知、关注关系、开发者身份和积分。</p>
        </div>
        <a-button danger @click="logout">退出登录</a-button>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="8">
          <a-card title="个人资料" class="inner-card">
            <a-form :model="profile" layout="vertical" @finish="saveProfile">
              <a-form-item label="用户名" name="username"><a-input v-model:value="profile.username" /></a-form-item>
              <a-form-item label="头像链接" name="avatar"><a-input v-model:value="profile.avatar" /></a-form-item>
              <a-form-item label="个人简介" name="bio"><a-textarea v-model:value="profile.bio" :rows="4" /></a-form-item>
              <a-form-item><a-button type="primary" html-type="submit">保存资料</a-button></a-form-item>
            </a-form>
          </a-card>

          <a-card title="贡献概览" class="inner-card" style="margin-top: 16px">
            <a-row :gutter="[12, 12]">
              <a-col :span="12"><a-statistic title="积分" :value="rewardSummary.points" /></a-col>
              <a-col :span="12"><a-statistic title="技能" :value="contributions.skills.length" /></a-col>
              <a-col :span="12"><a-statistic title="讨论" :value="contributions.discussions.length" /></a-col>
              <a-col :span="12"><a-statistic title="关注中" :value="followStats.follows" /></a-col>
              <a-col :span="12"><a-statistic title="粉丝" :value="followStats.followers" /></a-col>
              <a-col :span="12"><a-statistic title="已兑换权益" :value="rewardRedemptions.length" /></a-col>
            </a-row>
          </a-card>

          <a-card title="开发者认证" class="inner-card" style="margin-top: 16px">
            <template v-if="verification">
              <a-space wrap>
                <a-tag color="cyan">{{ verification.verification_type }}</a-tag>
                <a-tag :color="verification.status === 'approved' ? 'green' : verification.status === 'rejected' ? 'red' : 'gold'">
                  {{ verification.status }}
                </a-tag>
              </a-space>
              <p class="section-subtitle" style="margin-top: 10px">
                {{ verification.real_name }} {{ verification.organization ? `· ${verification.organization}` : '' }}
              </p>
              <p class="section-subtitle">{{ verification.review_comment || '申请已提交，等待审核。' }}</p>
            </template>

            <a-form :model="verificationForm" layout="vertical" @finish="submitVerification" style="margin-top: 12px">
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
              <a-form-item><a-button type="primary" html-type="submit">提交认证申请</a-button></a-form-item>
            </a-form>
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="16">
          <a-card title="我的资源" class="inner-card">
            <a-tabs>
              <a-tab-pane key="models" tab="我的模型">
                <ResourceCards :items="uploads.models" />
              </a-tab-pane>
              <a-tab-pane key="datasets" tab="我的数据集">
                <ResourceCards :items="uploads.datasets" />
              </a-tab-pane>
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
              <a-tab-pane key="notifications" tab="通知中心">
                <a-space style="margin-bottom: 12px">
                  <a-button @click="markNotificationsRead">全部标记已读</a-button>
                </a-space>
                <a-list :data-source="notifications">
                  <template #renderItem="{ item }">
                    <a-list-item>
                      <a-list-item-meta :title="item.Title" :description="item.Content" />
                    </a-list-item>
                  </template>
                </a-list>
              </a-tab-pane>
              <a-tab-pane key="rewards" tab="积分账本">
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
              <a-tab-pane key="rankings" tab="排行榜">
                <a-list :data-source="rankings">
                  <template #renderItem="{ item, index }">
                    <a-list-item>{{ index + 1 }}. {{ item.user_name }} · {{ item.points }}</a-list-item>
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
            </a-tabs>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink, useRouter } from 'vue-router';

import { api } from '@/api';
import ResourceCards from '@/components/ResourceCards.vue';
import { useAuthStore } from '@/stores/auth';
import type {
  ContributorRankingItem,
  DownloadRecord,
  FavoriteRecord,
  FollowItem,
  FollowStats,
  NotificationRecord,
  ResourceCard,
  RewardLedgerItem,
  RewardBenefitItem,
  RewardRedemptionItem,
  RewardSummary,
  UserContributionPayload,
  VerificationStatusItem,
} from '@/types/api';

const auth = useAuthStore();
const router = useRouter();

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
const notifications = ref<NotificationRecord[]>([]);
const rewardSummary = ref<RewardSummary>({ points: 0 });
const rewardLedger = ref<RewardLedgerItem[]>([]);
const rewardBenefits = ref<RewardBenefitItem[]>([]);
const rewardRedemptions = ref<RewardRedemptionItem[]>([]);
const rankings = ref<ContributorRankingItem[]>([]);
const follows = ref<FollowItem[]>([]);
const followers = ref<FollowItem[]>([]);
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
    message.success('认证申请已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '认证申请失败');
  }
}

async function logout() {
  await auth.logout();
  router.push('/');
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.inner-card {
  border-radius: 18px;
}

.redeem-card {
  border-radius: 16px;
}

.redeem-head {
  margin-bottom: 16px;
}
</style>

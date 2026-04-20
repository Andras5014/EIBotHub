<template>
  <div class="page-shell" v-if="user">
    <div class="profile-layout">
      <div class="profile-aside">
        <UserProfileSidebar
          :name="user.username"
          :handle="user.username"
          :bio="user.bio"
          :avatar="user.avatar"
          :badges="profileBadges"
          :stats="profileStats"
        >
          <template #actions>
            <a-space direction="vertical" style="width: 100%">
              <a-button type="primary" @click="follow">关注这个开发者</a-button>
              <a-button @click="messageUser">发送私信</a-button>
            </a-space>
          </template>

          <div class="meta-line">身份：社区开发者主页</div>
          <div class="meta-line">布局参考 GitHub 用户主页的资料侧栏与内容流结构。</div>
        </UserProfileSidebar>
      </div>

      <div class="profile-main">
        <section class="page-card hero-panel">
          <div class="section-head hero-head">
            <div>
              <div class="hero-kicker">Developer Profile</div>
              <h1 class="section-title hero-title">{{ user.username }} 的主页</h1>
              <p class="section-subtitle hero-subtitle">
                这里不再只是一张开发者卡片，而是像 GitHub 个人主页一样，把身份信息、代表性内容和最近参与放在同一条浏览路径里。
              </p>
            </div>
            <a-space wrap>
              <a-tag color="blue">技能 {{ contributions.skills.length }}</a-tag>
              <a-tag color="purple">讨论 {{ contributions.discussions.length }}</a-tag>
              <a-tag color="cyan">粉丝 {{ followStats.followers }}</a-tag>
            </a-space>
          </div>
        </section>

        <section class="page-card section-block">
          <a-tabs>
            <a-tab-pane key="overview" tab="概览">
              <div class="summary-grid">
                <div class="summary-card">
                  <div class="summary-label">技能贡献</div>
                  <div class="summary-value">{{ contributions.skills.length }}</div>
                </div>
                <div class="summary-card">
                  <div class="summary-label">讨论参与</div>
                  <div class="summary-value">{{ contributions.discussions.length }}</div>
                </div>
                <div class="summary-card">
                  <div class="summary-label">关注中</div>
                  <div class="summary-value">{{ followStats.follows }}</div>
                </div>
                <div class="summary-card">
                  <div class="summary-label">粉丝数</div>
                  <div class="summary-value">{{ followStats.followers }}</div>
                </div>
              </div>

              <div class="section-group">
                <div class="section-head">
                  <div>
                    <h2 class="section-title sub-title">置顶技能</h2>
                    <p class="section-subtitle">参考 GitHub pinned repositories 的方式，优先展示最能代表这个开发者能力边界的技能沉淀。</p>
                  </div>
                </div>
                <div v-if="pinnedSkills.length" class="content-grid">
                  <RouterLink v-for="item in pinnedSkills" :key="item.id" :to="`/skills/${item.id}`" class="content-card">
                    <div class="content-title">{{ item.name }}</div>
                    <div class="content-desc">{{ item.summary }}</div>
                    <a-space wrap>
                      <span class="pill-meta">{{ item.category }}</span>
                      <span class="pill-meta">{{ item.scene }}</span>
                      <span class="pill-meta">评分 {{ item.rating.average.toFixed(1) }}</span>
                    </a-space>
                  </RouterLink>
                </div>
                <a-empty v-else description="还没有公开技能内容" />
              </div>

              <div class="section-group">
                <div class="section-head">
                  <div>
                    <h2 class="section-title sub-title">最近讨论</h2>
                    <p class="section-subtitle">像 GitHub 活动流一样，保留最近公开讨论，方便快速判断当前关注主题。</p>
                  </div>
                </div>
                <a-list :data-source="recentDiscussions" class="activity-list">
                  <template #renderItem="{ item }">
                    <a-list-item>
                      <a-list-item-meta>
                        <template #title>
                          <a-space wrap>
                            <RouterLink :to="`/discussions/${item.id}`">{{ item.title }}</RouterLink>
                            <span class="pill-meta">{{ item.tag }}</span>
                            <span class="pill-meta">评论 {{ item.comment_count }}</span>
                          </a-space>
                        </template>
                        <template #description>{{ item.summary }}</template>
                      </a-list-item-meta>
                    </a-list-item>
                  </template>
                </a-list>
              </div>
            </a-tab-pane>

            <a-tab-pane key="skills" tab="技能">
              <a-empty v-if="!contributions.skills.length" description="还没有公开技能" />
              <div v-else class="content-grid">
                <RouterLink v-for="item in contributions.skills" :key="item.id" :to="`/skills/${item.id}`" class="content-card">
                  <div class="content-title">{{ item.name }}</div>
                  <div class="content-desc">{{ item.summary }}</div>
                  <a-space wrap>
                    <span class="pill-meta">{{ item.category }}</span>
                    <span class="pill-meta">{{ item.scene }}</span>
                    <span class="pill-meta">评分 {{ item.rating.average.toFixed(1) }}</span>
                  </a-space>
                </RouterLink>
              </div>
            </a-tab-pane>

            <a-tab-pane key="discussions" tab="讨论">
              <a-empty v-if="!contributions.discussions.length" description="还没有公开讨论" />
              <div v-else class="content-grid">
                <RouterLink v-for="item in contributions.discussions" :key="item.id" :to="`/discussions/${item.id}`" class="content-card">
                  <div class="content-title">{{ item.title }}</div>
                  <div class="content-desc">{{ item.summary }}</div>
                  <a-space wrap>
                    <span class="pill-meta">{{ item.tag }}</span>
                    <span class="pill-meta">评论 {{ item.comment_count }}</span>
                  </a-space>
                </RouterLink>
              </div>
            </a-tab-pane>
          </a-tabs>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import UserProfileSidebar from '@/components/UserProfileSidebar.vue';
import type { FollowStats, UserContributionPayload, UserSummary, VerificationStatusItem } from '@/types/api';

const route = useRoute();
const router = useRouter();
const user = ref<UserSummary>();
const contributions = ref<UserContributionPayload>({
  skills: [],
  discussions: [],
});
const followStats = ref<FollowStats>({
  follows: 0,
  followers: 0,
});
const verification = ref<VerificationStatusItem | null>(null);

const profileBadges = computed(() => {
  const badges = [{ text: '开发者', color: 'blue' }];
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
  { label: '技能', value: contributions.value.skills.length },
  { label: '讨论', value: contributions.value.discussions.length },
  { label: '关注中', value: followStats.value.follows },
  { label: '粉丝', value: followStats.value.followers },
]);

const pinnedSkills = computed(() => contributions.value.skills.slice(0, 5));
const recentDiscussions = computed(() => contributions.value.discussions.slice(0, 6));

async function load() {
  const id = String(route.params.id);
  const [profile, contributionData, stats, verificationStatus] = await Promise.all([
    api.getPublicUser(id),
    api.getPublicUserContributions(id),
    api.getPublicUserFollowStats(id),
    api.getPublicVerification(id),
  ]);
  user.value = profile;
  contributions.value = contributionData;
  followStats.value = stats;
  verification.value = verificationStatus;
}

watch(
  () => route.params.id,
  () => {
    void load();
  },
  { immediate: true },
);

async function follow() {
  if (!user.value) return;
  try {
    await api.toggleFollow(user.value.id);
    message.success('关注状态已更新');
    await load();
  } catch (error) {
    message.error(error instanceof Error ? error.message : '关注失败');
  }
}

function messageUser() {
  if (!user.value) return;
  router.push({ name: 'messages', query: { recipient: String(user.value.id) } });
}
</script>

<style scoped>
.profile-layout {
  display: grid;
  gap: 20px;
  grid-template-columns: minmax(280px, 320px) minmax(0, 1fr);
}

.hero-panel,
.section-block {
  padding: 24px;
  border-radius: 24px;
}

.hero-panel {
  background:
    radial-gradient(circle at top right, rgba(22, 119, 255, 0.16), transparent 32%),
    linear-gradient(145deg, #f8fbff 0%, #ffffff 58%, #f2fbf7 100%);
}

.hero-kicker {
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.hero-title {
  margin-top: 10px;
}

.hero-subtitle {
  max-width: 760px;
  line-height: 1.75;
}

.section-block {
  margin-top: 20px;
  padding-top: 12px;
}

.summary-grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.summary-card {
  padding: 18px;
  border-radius: 18px;
  background: var(--surface-soft);
  border: 1px solid rgba(220, 230, 242, 0.82);
}

.summary-label {
  color: var(--text-secondary);
  font-size: 12px;
}

.summary-value {
  margin-top: 8px;
  color: var(--text-main);
  font-size: 28px;
  font-weight: 700;
}

.section-group {
  margin-top: 28px;
}

.sub-title {
  font-size: 22px;
}

.content-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.content-card {
  display: block;
  padding: 18px;
  border-radius: 18px;
  border: 1px solid rgba(220, 230, 242, 0.82);
  background: linear-gradient(180deg, #fff 0%, #f9fbff 100%);
  transition: transform 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

.content-card:hover {
  transform: translateY(-2px);
  border-color: rgba(22, 119, 255, 0.34);
  box-shadow: 0 16px 32px rgba(22, 50, 79, 0.08);
}

.content-title {
  color: var(--text-main);
  font-size: 17px;
  font-weight: 700;
}

.content-desc {
  margin: 10px 0 14px;
  color: var(--text-secondary);
  line-height: 1.75;
  min-height: 48px;
}

.activity-list {
  border: 1px solid rgba(220, 230, 242, 0.82);
  border-radius: 18px;
  overflow: hidden;
}

.meta-line {
  margin-top: 4px;
}

@media (max-width: 1080px) {
  .profile-layout,
  .summary-grid,
  .content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>

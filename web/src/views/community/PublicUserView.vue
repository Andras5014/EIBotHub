<template>
  <div class="page-shell" v-if="user">
    <div class="page-card block">
      <a-tag color="blue">开发者</a-tag>
      <a-space v-if="verification" style="margin-left: 8px">
        <a-tag color="cyan">{{ verification.verification_type }}</a-tag>
        <a-tag :color="verification.status === 'approved' ? 'green' : verification.status === 'rejected' ? 'red' : 'gold'">
          {{ verification.status }}
        </a-tag>
      </a-space>
      <h1 class="section-title" style="margin-top: 10px">{{ user.username }}</h1>
      <p class="section-subtitle">{{ user.bio || '这个开发者暂未填写简介。' }}</p>
      <a-space style="margin-top: 16px">
        <a-button type="primary" @click="follow">关注这个开发者</a-button>
        <a-button @click="messageUser">发送私信</a-button>
      </a-space>

      <a-row :gutter="[12, 12]" style="margin-top: 16px">
        <a-col :xs="12" :lg="6"><a-statistic title="技能" :value="contributions.skills.length" /></a-col>
        <a-col :xs="12" :lg="6"><a-statistic title="讨论" :value="contributions.discussions.length" /></a-col>
        <a-col :xs="12" :lg="6"><a-statistic title="关注中" :value="followStats.follows" /></a-col>
        <a-col :xs="12" :lg="6"><a-statistic title="粉丝" :value="followStats.followers" /></a-col>
      </a-row>

      <a-row :gutter="[16, 16]" style="margin-top: 24px">
        <a-col :xs="24" :lg="12">
          <a-card title="TA 的技能">
            <a-list :data-source="contributions.skills">
              <template #renderItem="{ item }">
                <a-list-item>
                  <RouterLink :to="`/skills/${item.id}`">{{ item.name }}</RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="TA 的讨论">
            <a-list :data-source="contributions.discussions">
              <template #renderItem="{ item }">
                <a-list-item>
                  <RouterLink :to="`/discussions/${item.id}`">{{ item.title }}</RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
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

onMounted(async () => {
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
});

async function follow() {
  if (!user.value) return;
  try {
    await api.toggleFollow(user.value.id);
    message.success('关注状态已更新');
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
.block {
  padding: 24px;
}
</style>

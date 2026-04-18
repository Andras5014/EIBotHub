<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">社区治理</h1>
          <p class="section-subtitle">查看技能、讨论和评论，支持基础下架与删除操作。</p>
        </div>
      </div>

      <div class="grid-cards" v-if="overview">
        <a-card><a-statistic title="技能数" :value="overview.skills" /></a-card>
        <a-card><a-statistic title="讨论数" :value="overview.discussions" /></a-card>
        <a-card><a-statistic title="评论数" :value="overview.comments" /></a-card>
      </div>

      <a-row :gutter="[16, 16]" style="margin-top: 18px">
        <a-col :xs="24" :lg="10">
          <a-input v-model:value="keyword" placeholder="按标题、摘要、用户名或评论内容筛选" />
        </a-col>
        <a-col :xs="24" :lg="6">
          <a-select v-model:value="skillStatusFilter" style="width: 100%" :options="skillStatusOptions" />
        </a-col>
        <a-col :xs="24" :lg="8">
          <a-select v-model:value="commentTypeFilter" style="width: 100%" :options="commentTypeOptions" />
        </a-col>
      </a-row>

      <a-tabs style="margin-top: 18px">
        <a-tab-pane key="skills" tab="技能">
          <a-list :data-source="filteredSkills">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-space wrap>
                      <span>{{ item.name }}</span>
                      <a-tag :color="item.status === 'published' ? 'green' : 'red'">{{ item.status }}</a-tag>
                    </a-space>
                  </template>
                  <template #description>
                    <div>{{ item.summary }}</div>
                    <a-space wrap style="margin-top: 6px">
                      <RouterLink :to="`/community/users/${item.owner_id}`">查看作者</RouterLink>
                      <span>作者：{{ item.owner_name }}</span>
                      <span>更新时间：{{ formatDate(item.updated_at) }}</span>
                    </a-space>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button danger type="link" @click="hideSkill(item.id)">下架</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
        <a-tab-pane key="discussions" tab="讨论">
          <a-list :data-source="filteredDiscussions">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <RouterLink :to="`/discussions/${item.id}`">{{ item.title }}</RouterLink>
                  </template>
                  <template #description>
                    <a-space wrap>
                      <span class="pill-meta">{{ item.category }}</span>
                      <RouterLink :to="`/community/users/${item.user_id}`">查看作者</RouterLink>
                      <span>作者：{{ item.user_name }}</span>
                      <span>评论：{{ item.comment_count }}</span>
                      <span>更新时间：{{ formatDate(item.updated_at) }}</span>
                    </a-space>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button danger type="link" @click="removeDiscussion(item.id)">删除</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
        <a-tab-pane key="comments" tab="评论">
          <a-list :data-source="filteredComments">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-space wrap>
                      <RouterLink :to="resourceRoute(item.resource_type, item.resource_id)">
                        {{ resourceLabel(item.resource_type) }} #{{ item.resource_id }}
                      </RouterLink>
                      <RouterLink :to="`/community/users/${item.user_id}`">查看作者</RouterLink>
                      <span>作者：{{ item.user_name }}</span>
                    </a-space>
                  </template>
                  <template #description>
                    <div>{{ item.content }}</div>
                    <a-space wrap style="margin-top: 6px">
                      <span>创建时间：{{ formatDate(item.created_at) }}</span>
                    </a-space>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button danger type="link" @click="removeComment(item.id)">删除</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import type {
  AdminCommentModerationItem,
  AdminCommunityOverview,
  AdminDiscussionModerationItem,
  AdminSkillModerationItem,
} from '@/types/api';

const overview = ref<AdminCommunityOverview>();
const skills = ref<AdminSkillModerationItem[]>([]);
const discussions = ref<AdminDiscussionModerationItem[]>([]);
const comments = ref<AdminCommentModerationItem[]>([]);
const keyword = ref('');
const skillStatusFilter = ref<'all' | string>('all');
const commentTypeFilter = ref<'all' | string>('all');

const skillStatusOptions = [
  { label: '全部技能状态', value: 'all' },
  { label: '已发布', value: 'published' },
  { label: '已驳回', value: 'rejected' },
];

const commentTypeOptions = [
  { label: '全部评论类型', value: 'all' },
  { label: '模型', value: 'model' },
  { label: '数据集', value: 'dataset' },
  { label: '模板', value: 'task-template' },
  { label: '技能', value: 'skill' },
  { label: '讨论', value: 'discussion' },
];

const filteredSkills = computed(() =>
  skills.value.filter((item) => {
    const matchKeyword =
      !keyword.value.trim() ||
      `${item.name} ${item.summary} ${item.owner_name}`.includes(keyword.value.trim());
    const matchStatus = skillStatusFilter.value === 'all' || item.status === skillStatusFilter.value;
    return matchKeyword && matchStatus;
  }),
);

const filteredDiscussions = computed(() =>
  discussions.value.filter((item) =>
    !keyword.value.trim() || `${item.title} ${item.category} ${item.user_name}`.includes(keyword.value.trim()),
  ),
);

const filteredComments = computed(() =>
  comments.value.filter((item) => {
    const matchKeyword =
      !keyword.value.trim() ||
      `${item.user_name} ${item.content} ${item.resource_type}`.includes(keyword.value.trim());
    const matchType = commentTypeFilter.value === 'all' || item.resource_type === commentTypeFilter.value;
    return matchKeyword && matchType;
  }),
);

async function load() {
  const [communityOverview, skillItems, discussionItems, commentItems] = await Promise.all([
    api.getAdminCommunityOverview(),
    api.getAdminCommunitySkills(),
    api.getAdminCommunityDiscussions(),
    api.getAdminCommunityComments(),
  ]);
  overview.value = communityOverview;
  skills.value = skillItems;
  discussions.value = discussionItems;
  comments.value = commentItems;
}

async function hideSkill(id: number) {
  await api.hideAdminSkill(id);
  message.success('技能已下架');
  await load();
}

async function removeDiscussion(id: number) {
  await api.removeAdminDiscussion(id);
  message.success('讨论已删除');
  await load();
}

async function removeComment(id: number) {
  await api.removeAdminComment(id);
  message.success('评论已删除');
  await load();
}

function resourceRoute(resourceType: string, resourceID: number) {
  return {
    model: `/models/${resourceID}`,
    dataset: `/datasets/${resourceID}`,
    'task-template': `/templates/${resourceID}`,
    skill: `/skills/${resourceID}`,
    discussion: `/discussions/${resourceID}`,
  }[resourceType] ?? '/admin/community';
}

function resourceLabel(resourceType: string) {
  return {
    model: '模型',
    dataset: '数据集',
    'task-template': '模板',
    skill: '技能',
    discussion: '讨论',
  }[resourceType] ?? resourceType;
}

function formatDate(value: string) {
  return new Date(value).toLocaleString();
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

<template>
  <div class="page-shell" v-if="detail">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <a-space wrap>
            <span class="pill-meta">{{ detail.category }}</span>
            <span class="pill-meta">{{ detail.scene }}</span>
            <span class="pill-meta">作者：{{ detail.owner_name }}</span>
          </a-space>
          <h1 class="section-title" style="margin-top: 10px">{{ detail.name }}</h1>
          <p class="section-subtitle">{{ detail.summary }}</p>
        </div>
        <a-space>
          <a-button @click="contactAuthor">联系作者</a-button>
          <a-button @click="startWorkspace">进入协作</a-button>
          <a-button type="primary" @click="forkSkill">Fork 技能</a-button>
        </a-space>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="12">
          <a-card title="技能说明">
            <p>{{ detail.description }}</p>
            <a-divider />
            <p style="white-space: pre-wrap">{{ detail.guide }}</p>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="评分与互动">
            <a-statistic :value="ratings?.average ?? 0" :precision="1" title="平均分" />
            <p class="section-subtitle">共 {{ ratings?.count ?? 0 }} 条评分</p>
            <a-divider />
            <a-form :model="ratingForm" layout="vertical" @finish="submitRating">
              <a-form-item label="评分"><a-rate v-model:value="ratingForm.score" /></a-form-item>
              <a-form-item label="反馈"><a-textarea v-model:value="ratingForm.feedback" :rows="3" /></a-form-item>
              <a-form-item><a-button type="primary" html-type="submit">提交评分</a-button></a-form-item>
            </a-form>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="评论区" style="margin-top: 16px">
        <a-form :model="commentForm" layout="vertical" @finish="submitComment">
          <a-form-item label="评论内容"><a-textarea v-model:value="commentForm.content" :rows="3" /></a-form-item>
          <a-form-item><a-button type="primary" html-type="submit">发表评论</a-button></a-form-item>
        </a-form>
        <a-list :data-source="comments">
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta :title="item.user_name" :description="item.content" />
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import type { CommentItem, RatingSummary, SkillItem } from '@/types/api';

const route = useRoute();
const router = useRouter();
const detail = ref<SkillItem>();
const ratings = ref<RatingSummary>();
const comments = ref<CommentItem[]>([]);
const ratingForm = reactive({ score: 5, feedback: '' });
const commentForm = reactive({ content: '' });

async function load() {
  const id = String(route.params.id);
  const [skill, ratingSummary, commentItems] = await Promise.all([
    api.getSkill(id),
    api.getSkillRatings(id),
    api.getSkillComments(id),
  ]);
  detail.value = skill;
  ratings.value = ratingSummary;
  comments.value = commentItems;
}

async function submitRating() {
  if (!detail.value) return;
  try {
    ratings.value = await api.rateSkill(detail.value.id, ratingForm);
    ratingForm.feedback = '';
    message.success('评分已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '评分失败');
  }
}

async function submitComment() {
  if (!detail.value) return;
  try {
    await api.commentSkill(detail.value.id, commentForm);
    comments.value = await api.getSkillComments(detail.value.id);
    commentForm.content = '';
    message.success('评论已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '评论失败');
  }
}

async function forkSkill() {
  if (!detail.value) return;
  try {
    const next = await api.forkSkill(detail.value.id);
    message.success('Fork 成功');
    window.location.href = `/skills/${next.id}`;
  } catch (error) {
    message.error(error instanceof Error ? error.message : 'Fork 失败');
  }
}

function contactAuthor() {
  if (!detail.value) return;
  router.push({ name: 'messages', query: { recipient: String(detail.value.owner_id) } });
}

function startWorkspace() {
  if (!detail.value) return;
  router.push({
    name: 'workspaces',
    query: {
      name: `${detail.value.name} 协作空间`,
      summary: `围绕技能《${detail.value.name}》进行协作与联调。`,
    },
  });
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

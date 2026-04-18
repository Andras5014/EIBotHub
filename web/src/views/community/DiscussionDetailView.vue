<template>
  <div class="page-shell" v-if="detail">
    <div class="page-card block">
      <a-space wrap>
        <span class="pill-meta">{{ detail.category }}</span>
        <span class="pill-meta">{{ detail.user_name }}</span>
        <span class="pill-meta">评论 {{ detail.comment_count }}</span>
      </a-space>
      <h1 class="section-title" style="margin-top: 10px">{{ detail.title }}</h1>
      <p class="section-subtitle">{{ detail.summary }}</p>
      <a-space style="margin-top: 12px">
        <a-button @click="contactAuthor">联系作者</a-button>
        <a-button @click="startWorkspace">创建协作空间</a-button>
      </a-space>

      <a-card title="讨论正文" style="margin-top: 18px">
        <p style="white-space: pre-wrap">{{ detail.content }}</p>
      </a-card>

      <a-card title="评论回复" style="margin-top: 16px">
        <a-alert
          v-if="replyTo"
          type="info"
          show-icon
          style="margin-bottom: 16px"
          :message="`正在回复 ${replyTo.user_name}`"
          :description="replyTo.content"
        >
          <template #action>
            <a-button type="link" @click="clearReply">取消回复</a-button>
          </template>
        </a-alert>

        <a-form :model="form" layout="vertical" @finish="submit">
          <a-form-item label="评论内容">
            <a-textarea
              v-model:value="form.content"
              :rows="3"
              :disabled="!auth.isAuthenticated"
              :placeholder="replyTo ? '补充你的回复内容' : '发表你的看法'"
            />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" html-type="submit" :disabled="!auth.isAuthenticated">发表评论</a-button>
          </a-form-item>
        </a-form>

        <a-empty v-if="!commentThreads.length" description="还没有评论" />

        <div v-else class="comment-thread-list">
          <div v-for="item in commentThreads" :key="item.id" class="comment-card">
            <div class="comment-head">
              <div class="comment-author">{{ item.user_name }}</div>
              <div class="comment-date">{{ formatDate(item.created_at) }}</div>
            </div>
            <div class="comment-body">{{ item.content }}</div>
            <a-space>
              <a-button type="link" size="small" @click="setReply(item)">回复</a-button>
            </a-space>

            <div v-if="item.replies.length" class="reply-list">
              <div v-for="reply in item.replies" :key="reply.id" class="reply-card">
                <div class="comment-head">
                  <div class="comment-author">{{ reply.user_name }}</div>
                  <div class="comment-date">{{ formatDate(reply.created_at) }}</div>
                </div>
                <div class="comment-body">{{ reply.content }}</div>
                <a-button type="link" size="small" @click="setReply(reply)">继续回复</a-button>
              </div>
            </div>
          </div>
        </div>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import { useAuthStore } from '@/stores/auth';
import type { CommentItem, DiscussionItem } from '@/types/api';

type ThreadComment = CommentItem & { replies: CommentItem[] };

const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const detail = ref<DiscussionItem>();
const comments = ref<CommentItem[]>([]);
const replyTo = ref<CommentItem>();
const form = reactive({
  content: '',
  parent_id: undefined as number | undefined,
});

const commentThreads = computed<ThreadComment[]>(() => {
  const roots = comments.value.filter((item) => !item.parent_id);
  return roots.map((item) => ({
    ...item,
    replies: comments.value.filter((reply) => reply.parent_id === item.id),
  }));
});

async function load() {
  const id = String(route.params.id);
  const [discussion, discussionComments] = await Promise.all([
    api.getDiscussion(id),
    api.getDiscussionComments(id),
  ]);
  detail.value = discussion;
  comments.value = discussionComments;
}

async function submit() {
  if (!detail.value) return;
  if (!auth.isAuthenticated) {
    message.info('请先登录后再发表评论');
    return;
  }
  try {
    await api.commentDiscussion(detail.value.id, {
      content: form.content,
      parent_id: form.parent_id,
    });
    form.content = '';
    clearReply();
    await load();
    message.success('评论已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '评论失败');
  }
}

function setReply(item: CommentItem) {
  replyTo.value = item;
  form.parent_id = item.parent_id ?? item.id;
}

function clearReply() {
  replyTo.value = undefined;
  form.parent_id = undefined;
}

function contactAuthor() {
  if (!detail.value) return;
  router.push({ name: 'messages', query: { recipient: String(detail.value.user_id) } });
}

function startWorkspace() {
  if (!detail.value) return;
  router.push({
    name: 'workspaces',
    query: {
      name: `${detail.value.title} 协作空间`,
      summary: `围绕讨论《${detail.value.title}》进行协作。`,
    },
  });
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('zh-CN');
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.comment-thread-list {
  display: grid;
  gap: 12px;
}

.comment-card,
.reply-card {
  border-radius: 16px;
  border: 1px solid var(--line);
  background: var(--surface-soft);
  padding: 14px;
}

.reply-list {
  margin-top: 12px;
  display: grid;
  gap: 10px;
  padding-left: 18px;
  border-left: 2px solid rgba(22, 119, 255, 0.14);
}

.comment-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.comment-author {
  font-weight: 700;
}

.comment-date {
  color: var(--text-secondary);
  font-size: 12px;
}

.comment-body {
  white-space: pre-wrap;
  color: var(--text-main);
  margin-bottom: 8px;
}
</style>

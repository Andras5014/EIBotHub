<template>
  <div class="page-shell" v-if="detail">
    <div class="page-card block">
      <a-space wrap>
        <span class="pill-meta">{{ detail.tag }}</span>
        <span class="pill-meta">{{ detail.user_name }}</span>
        <span class="pill-meta">评论 {{ detail.comment_count }}</span>
      </a-space>
      <h1 class="section-title" style="margin-top: 10px">{{ detail.title }}</h1>
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
          <div v-for="item in commentThreads" :key="item.id" :class="['comment-card', { 'author-comment': isPostAuthor(item.user_id) }]">
            <div class="comment-author-row">
              <div class="comment-author">{{ item.user_name }}</div>
              <span v-if="isPostAuthor(item.user_id)" class="author-badge">贴主</span>
            </div>
            <div class="comment-meta">
              <div class="comment-relation">{{ relationLabel(item) }}</div>
              <div class="comment-date">{{ formatDate(item.created_at) }}</div>
            </div>
            <div class="comment-body">{{ item.content }}</div>
            <a-space class="comment-actions">
              <a-button type="link" size="small" @click="setReply(item)">回复</a-button>
            </a-space>

            <div v-if="item.replies.length" class="reply-list">
              <div v-for="reply in item.replies" :key="reply.id" :class="['reply-card', { 'author-comment': isPostAuthor(reply.user_id) }]">
                <div class="comment-author-row">
                  <div class="comment-author">{{ reply.user_name }}</div>
                  <span v-if="isPostAuthor(reply.user_id)" class="author-badge">贴主</span>
                </div>
                <div class="comment-meta">
                  <div class="comment-relation">{{ relationLabel(reply) }}</div>
                  <div class="comment-date">{{ formatDate(reply.created_at) }}</div>
                </div>
                <div class="comment-body">{{ reply.content }}</div>
                <a-button type="link" size="small" class="comment-actions" @click="setReply(reply)">继续回复</a-button>
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

const commentLookup = computed(() => new Map(comments.value.map((item) => [item.id, item])));

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

function isPostAuthor(userID: number) {
  return userID === detail.value?.user_id;
}

function relationLabel(item: CommentItem) {
  if (!item.parent_id) {
    return '回复帖子';
  }
  const targetName = commentLookup.value.get(item.parent_id)?.user_name ?? detail.value?.user_name ?? '贴主';
  return `回复 ${targetName}`;
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
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  grid-template-areas:
    "author meta"
    "body meta"
    "actions meta";
  gap: 4px 12px;
  align-items: start;
  border-radius: 16px;
  border: 1px solid var(--line);
  background: var(--surface-soft);
  padding: 14px;
}

.author-comment {
  border-color: rgba(22, 119, 255, 0.28);
  background: linear-gradient(180deg, #fafdff 0%, #f2f8ff 100%);
}

.reply-list {
  margin-top: 12px;
  display: grid;
  gap: 10px;
  padding-left: 18px;
  border-left: 2px solid rgba(22, 119, 255, 0.14);
}

.comment-author-row {
  grid-area: author;
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.comment-author {
  font-weight: 700;
}

.author-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px 8px;
  border-radius: 999px;
  background: rgba(22, 119, 255, 0.12);
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 700;
}

.comment-meta {
  grid-area: meta;
  display: grid;
  justify-items: end;
  gap: 2px;
}

.comment-relation {
  color: var(--text-secondary);
  font-size: 12px;
  white-space: nowrap;
}

.comment-date {
  color: var(--text-secondary);
  font-size: 12px;
}

.comment-body {
  grid-area: body;
  white-space: pre-wrap;
  color: var(--text-main);
  margin-bottom: 0;
}

.comment-actions {
  grid-area: actions;
}

@media (max-width: 720px) {
  .comment-card,
  .reply-card {
    grid-template-columns: 1fr;
    grid-template-areas:
      "author"
      "meta"
      "body"
      "actions";
  }

  .comment-meta {
    justify-items: start;
  }
}
</style>

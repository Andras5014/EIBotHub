<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">私信</h1>
          <p class="section-subtitle">查看已有会话，或直接向指定用户发送消息。</p>
        </div>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="8">
          <a-card title="发起私信" class="inner-card">
            <a-form :model="form" layout="vertical" @finish="sendDirect">
              <a-form-item label="接收用户">
                <UserSearchSelect v-model="form.recipient_user_id" placeholder="搜索用户名" />
              </a-form-item>
              <a-form-item label="内容"><a-textarea v-model:value="form.content" :rows="4" /></a-form-item>
              <a-form-item><a-button type="primary" html-type="submit">发送</a-button></a-form-item>
            </a-form>
          </a-card>
          <a-card title="会话列表" class="inner-card" style="margin-top: 16px">
            <a-list :data-source="conversations">
              <template #renderItem="{ item }">
                <a-list-item :class="{ active: selectedConversation?.id === item.id }" @click="selectConversation(item)" style="cursor: pointer">
                  <a-list-item-meta>
                    <template #title>
                      <a-space wrap>
                        <span>{{ conversationLabel(item) }}</span>
                        <a-tag :color="item.active ? 'green' : 'red'">{{ item.active ? '正常' : '已封禁' }}</a-tag>
                      </a-space>
                    </template>
                    <template #description>
                      <div class="participant-text">参与人：{{ participantSummary(item.participant_names) }}</div>
                      <div>{{ item.latest_message || '暂无消息' }}</div>
                      <div v-if="item.blocked_reason" class="blocked-text">封禁原因：{{ item.blocked_reason }}</div>
                    </template>
                  </a-list-item-meta>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="16">
          <a-card title="消息内容" class="inner-card">
            <template v-if="selectedConversation">
              <a-space wrap style="margin-bottom: 12px">
                <span class="pill-meta">会话类型：{{ selectedConversation.kind }}</span>
                <span class="pill-meta">参与人：{{ participantSummary(selectedConversation.participant_names) }}</span>
                <a-tag :color="selectedConversation.active ? 'green' : 'red'">
                  {{ selectedConversation.active ? '正常' : '已封禁' }}
                </a-tag>
              </a-space>
              <a-alert
                v-if="!selectedConversation.active"
                type="warning"
                show-icon
                :message="selectedConversation.blocked_reason || '当前会话已被管理员封禁，暂不可继续发送或查看消息。'"
                style="margin-bottom: 16px"
              />
              <a-list :data-source="messages">
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta :title="item.sender_name" :description="`${item.content} · ${formatDate(item.created_at)}`" />
                  </a-list-item>
                </template>
              </a-list>
              <a-divider />
              <a-form :model="replyForm" layout="vertical" @finish="sendReply">
                <a-form-item label="回复内容">
                  <a-textarea v-model:value="replyForm.content" :rows="3" :disabled="!selectedConversation.active" />
                </a-form-item>
                <a-form-item>
                  <a-button type="primary" html-type="submit" :disabled="!selectedConversation.active">回复</a-button>
                </a-form-item>
              </a-form>
            </template>
            <a-empty v-else description="请选择一个会话" />
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { message } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import UserSearchSelect from '@/components/UserSearchSelect.vue';
import type { ConversationItem, MessageItem } from '@/types/api';

const route = useRoute();
const router = useRouter();
const conversations = ref<ConversationItem[]>([]);
const selectedConversation = ref<ConversationItem>();
const messages = ref<MessageItem[]>([]);
const form = reactive({
  recipient_user_id: undefined as number | undefined,
  content: '',
});
const replyForm = reactive({
  content: '',
});

async function loadConversations() {
  conversations.value = await api.listConversations();
}

async function selectConversation(item: ConversationItem) {
  selectedConversation.value = item;
  if (!item.active) {
    messages.value = [];
    router.replace({ query: { ...route.query, conversation: String(item.id) } });
    return;
  }
  messages.value = await api.getConversationMessages(item.id);
  router.replace({ query: { ...route.query, conversation: String(item.id) } });
}

async function sendDirect() {
  if (!form.recipient_user_id) {
    message.error('请选择接收用户');
    return;
  }
  try {
    const created = await api.sendMessage(form);
    form.content = '';
    await loadConversations();
    const target = conversations.value.find((item) => item.id === created.conversation_id);
    if (target) {
      await selectConversation(target);
    }
    message.success('私信已发送');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '发送失败');
  }
}

async function sendReply() {
  if (!selectedConversation.value) return;
  if (!selectedConversation.value.active) {
    message.error('当前会话已被封禁');
    return;
  }
  try {
    await api.sendMessage({
      conversation_id: selectedConversation.value.id,
      content: replyForm.content,
    });
    replyForm.content = '';
    messages.value = await api.getConversationMessages(selectedConversation.value.id);
    await loadConversations();
    message.success('消息已发送');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '发送失败');
  }
}

watch(
  () => route.query.recipient,
  (recipient) => {
    if (recipient) {
      form.recipient_user_id = Number(recipient);
    }
  },
  { immediate: true },
);

watch(
  () => route.query.conversation,
  async (conversation) => {
    if (!conversation) return;
    const found = conversations.value.find((item) => item.id === Number(conversation));
    if (found) {
      await selectConversation(found);
    }
  },
);

onMounted(async () => {
  await loadConversations();
  const conversation = route.query.conversation;
  if (conversation) {
    const found = conversations.value.find((item) => item.id === Number(conversation));
    if (found) {
      await selectConversation(found);
      return;
    }
  }
  if (conversations.value.length > 0) {
    await selectConversation(conversations.value[0]);
  }
});

function formatDate(value: string) {
  return new Date(value).toLocaleString();
}

function conversationLabel(item: ConversationItem) {
  if (item.kind === 'direct' && item.participant_names.length > 0) {
    return item.participant_names.join(' / ');
  }
  if (item.title) {
    return item.title;
  }
  if (item.participant_names.length > 0) {
    return item.participant_names.join(' / ');
  }
  return '未命名会话';
}

function participantSummary(names: string[]) {
  return names.length > 0 ? names.join('、') : '暂无';
}
</script>

<style scoped>
.block {
  padding: 24px;
}

.inner-card {
  border-radius: 18px;
}

.active {
  background: var(--surface-soft);
  border-radius: 14px;
}

.blocked-text {
  margin-top: 4px;
  color: #cf1322;
}

.participant-text {
  color: var(--text-secondary);
}
</style>

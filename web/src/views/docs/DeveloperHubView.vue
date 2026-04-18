<template>
  <div class="page-shell">
    <div class="page-card hero-card">
      <div class="section-head">
        <div>
          <a-tag color="blue">开放生态</a-tag>
          <h1 class="section-title" style="margin-top: 12px">开发者接入</h1>
          <p class="section-subtitle">
            统一查看 OpenAPI 基础信息、可订阅事件、CLI 示例，并直接管理当前账号的 Webhook 回调。
          </p>
        </div>
        <a-space wrap>
          <a-tag color="cyan">Base URL {{ spec?.base_url ?? '/api/v1' }}</a-tag>
          <a-tag color="gold">Version {{ spec?.version ?? 'v1' }}</a-tag>
        </a-space>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="14">
          <a-card title="OpenAPI 概览" class="inner-card">
            <p class="section-subtitle">{{ spec?.overview }}</p>
            <a-table
              :data-source="spec?.endpoints ?? []"
              :pagination="false"
              row-key="path"
              size="small"
              :columns="endpointColumns"
            />
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="10">
          <a-card title="可订阅事件" class="inner-card">
            <a-space wrap>
              <a-tag v-for="item in spec?.webhook_events ?? []" :key="item" color="processing">{{ item }}</a-tag>
            </a-space>
          </a-card>
          <a-card title="CLI 示例" class="inner-card" style="margin-top: 16px">
            <div v-for="item in spec?.curl_examples ?? []" :key="item.title" class="code-block">
              <div class="list-title">{{ item.title }}</div>
              <pre>{{ item.command }}</pre>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <div class="page-card section-block" style="margin-top: 20px">
      <div class="section-head">
        <div>
          <h2 class="section-title">Webhook 管理</h2>
          <p class="section-subtitle">登录后可创建回调地址、选择事件，并直接发起测试投递。</p>
        </div>
      </div>

      <template v-if="auth.isAuthenticated">
        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :lg="9">
            <a-card title="新建订阅" class="inner-card">
              <a-form :model="createForm" layout="vertical" @finish="submitWebhook">
                <a-form-item label="名称" name="name">
                  <a-input v-model:value="createForm.name" placeholder="例如：本地联调" />
                </a-form-item>
                <a-form-item label="回调地址" name="target_url">
                  <a-input v-model:value="createForm.target_url" placeholder="http://127.0.0.1:9000/hook" />
                </a-form-item>
                <a-form-item label="签名密钥" name="secret">
                  <a-input-password v-model:value="createForm.secret" placeholder="至少 8 位" />
                </a-form-item>
                <a-form-item label="订阅事件" name="events">
                  <a-select v-model:value="createForm.events" mode="multiple" :options="eventOptions" placeholder="选择事件" />
                </a-form-item>
                <a-form-item>
                  <a-button type="primary" html-type="submit" :loading="creating">创建订阅</a-button>
                </a-form-item>
              </a-form>
            </a-card>
          </a-col>

          <a-col :xs="24" :lg="15">
            <a-card title="我的订阅" class="inner-card">
              <a-empty v-if="!webhooks.length" description="还没有创建 Webhook 订阅" />
              <a-row v-else :gutter="[16, 16]">
                <a-col v-for="item in webhooks" :key="item.id" :xs="24" :xl="12">
                  <a-card class="webhook-card" :class="{ active: selectedWebhook?.id === item.id }">
                    <div class="section-head" style="margin-bottom: 12px">
                      <div>
                        <div class="list-title">{{ item.name }}</div>
                        <div class="list-desc">{{ item.target_url }}</div>
                      </div>
                      <a-tag :color="item.last_status === 'success' ? 'green' : item.last_status === 'failed' ? 'red' : 'default'">
                        {{ item.last_status || '未投递' }}
                      </a-tag>
                    </div>
                    <a-space wrap style="margin-bottom: 12px">
                      <a-tag v-for="event in item.events" :key="event" color="processing">{{ event }}</a-tag>
                    </a-space>
                    <p class="section-subtitle">
                      最近响应：{{ item.last_response_code || '无' }}
                      <span v-if="item.last_error"> · {{ item.last_error }}</span>
                    </p>
                    <a-space wrap>
                      <a-button @click="selectWebhook(item)">查看记录</a-button>
                      <a-button type="primary" :loading="testingId === item.id" @click="triggerTest(item)">测试投递</a-button>
                    </a-space>
                  </a-card>
                </a-col>
              </a-row>
            </a-card>
          </a-col>
        </a-row>

        <a-card v-if="selectedWebhook" :title="`${selectedWebhook.name} · 最近投递`" class="inner-card" style="margin-top: 16px">
          <a-list :data-source="deliveries">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta
                  :title="`${item.event} · ${item.status} · ${item.response_code || '无响应码'}`"
                  :description="`${item.response_body || '无响应内容'} · ${new Date(item.created_at).toLocaleString()}`"
                />
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </template>

      <a-empty v-else description="登录后可管理当前账号的 Webhook 订阅">
        <RouterLink to="/login">
          <a-button type="primary">登录后继续</a-button>
        </RouterLink>
      </a-empty>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import { useAuthStore } from '@/stores/auth';
import type {
  OpenAPIEndpointItem,
  OpenAPISpecResponse,
  WebhookDeliveryItem,
  WebhookSubscriptionItem,
} from '@/types/api';

const auth = useAuthStore();
const spec = ref<OpenAPISpecResponse>();
const webhooks = ref<WebhookSubscriptionItem[]>([]);
const deliveries = ref<WebhookDeliveryItem[]>([]);
const selectedWebhook = ref<WebhookSubscriptionItem>();
const creating = ref(false);
const testingId = ref<number>();

const createForm = reactive({
  name: '',
  target_url: '',
  secret: '',
  events: [] as string[],
});

const endpointColumns = [
  { title: '方法', dataIndex: 'method', key: 'method', width: 88 },
  { title: '路径', dataIndex: 'path', key: 'path' },
  { title: '说明', dataIndex: 'summary', key: 'summary' },
  { title: '鉴权', dataIndex: 'auth', key: 'auth', width: 120 },
] as Array<{ title: string; dataIndex: keyof OpenAPIEndpointItem | 'auth'; key: string; width?: number }>;

const eventOptions = computed(() =>
  (spec.value?.webhook_events ?? []).map((item) => ({
    label: item,
    value: item,
  })),
);

async function loadSpec() {
  spec.value = await api.openAPISpec();
}

async function loadWebhooks() {
  if (!auth.isAuthenticated) {
    webhooks.value = [];
    deliveries.value = [];
    selectedWebhook.value = undefined;
    return;
  }
  webhooks.value = await api.listWebhooks();
  if (selectedWebhook.value) {
    const current = webhooks.value.find((item) => item.id === selectedWebhook.value?.id);
    if (current) {
      await selectWebhook(current);
      return;
    }
  }
  if (webhooks.value.length) {
    await selectWebhook(webhooks.value[0]);
  }
}

async function selectWebhook(item: WebhookSubscriptionItem) {
  selectedWebhook.value = item;
  deliveries.value = await api.getWebhookDeliveries(item.id);
}

async function submitWebhook() {
  creating.value = true;
  try {
    const created = await api.createWebhook(createForm);
    createForm.name = '';
    createForm.target_url = '';
    createForm.secret = '';
    createForm.events = [];
    await loadWebhooks();
    await selectWebhook(created);
    message.success('Webhook 订阅已创建');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '创建 Webhook 失败');
  } finally {
    creating.value = false;
  }
}

async function triggerTest(item: WebhookSubscriptionItem) {
  testingId.value = item.id;
  try {
    await api.testWebhook(item.id);
    await loadWebhooks();
    await selectWebhook(item);
    message.success('测试投递已发送');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '测试投递失败');
  } finally {
    testingId.value = undefined;
  }
}

onMounted(async () => {
  await loadSpec();
  await loadWebhooks();
});
</script>

<style scoped>
.hero-card,
.section-block {
  padding: 24px;
}

.inner-card {
  border-radius: 18px;
}

.code-block + .code-block {
  margin-top: 16px;
}

.code-block pre {
  margin: 10px 0 0;
  padding: 14px;
  overflow: auto;
  border-radius: 14px;
  background: #f5f7fb;
  color: #23314e;
  font-size: 12px;
  line-height: 1.6;
}

.webhook-card {
  border-radius: 16px;
}

.webhook-card.active {
  border-color: #1677ff;
  box-shadow: 0 10px 30px rgba(22, 119, 255, 0.08);
}
</style>

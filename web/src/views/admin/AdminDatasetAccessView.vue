<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">数据集访问审批</h1>
          <p class="section-subtitle">审核 `internal / restricted` 数据集的下载申请，通过后用户才能下载或生成分批下载任务。</p>
        </div>
      </div>

      <a-list :data-source="items">
        <template #renderItem="{ item }">
          <a-list-item class="request-item">
            <a-list-item-meta>
              <template #title>
                <a-space wrap>
                  <span>{{ item.user_name || `用户 ${item.user_id}` }}</span>
                  <a-tag color="blue">数据集 {{ item.dataset_id }}</a-tag>
                  <a-tag :color="statusColor(item.status)">{{ item.status }}</a-tag>
                </a-space>
              </template>
              <template #description>
                <div>申请理由：{{ item.reason }}</div>
                <div v-if="item.review_comment">审核意见：{{ item.review_comment }}</div>
                <div>申请时间：{{ formatDate(item.created_at) }}</div>
              </template>
            </a-list-item-meta>
            <template v-if="item.status === 'pending'" #actions>
              <a-button type="link" @click="review(item.id, 'approved')">通过</a-button>
              <a-button danger type="link" @click="review(item.id, 'rejected')">驳回</a-button>
            </template>
          </a-list-item>
        </template>
      </a-list>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { message } from 'ant-design-vue';

import { api } from '@/api';
import type { DatasetAccessRequestItem } from '@/types/api';

const items = ref<DatasetAccessRequestItem[]>([]);

async function load() {
  items.value = await api.getAdminDatasetAccessRequests();
}

async function review(id: number, decision: 'approved' | 'rejected') {
  await api.reviewAdminDatasetAccessRequest(id, {
    decision,
    comment: decision === 'approved' ? '访问申请通过' : '当前申请未通过，请补充说明后重试',
  });
  message.success('访问申请状态已更新');
  await load();
}

function statusColor(status: string) {
  return {
    pending: 'gold',
    approved: 'green',
    rejected: 'red',
  }[status] ?? 'default';
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

.request-item {
  border-radius: 16px;
  border: 1px solid var(--line);
  margin-bottom: 12px;
  background: var(--surface-soft);
}
</style>

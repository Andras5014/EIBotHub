<template>
  <div class="grid-cards">
    <a-card v-for="item in items" :key="item.id" :bordered="false" class="resource-card">
      <div class="card-top">
        <div>
          <h3 class="card-title">{{ item.name }}</h3>
          <p class="card-summary">{{ item.summary }}</p>
        </div>
        <a-space direction="vertical" align="end" :size="8">
          <a-tag v-if="item.badge_label" color="orange">{{ item.badge_label }}</a-tag>
          <a-tag color="blue">{{ typeLabel(item.type) }}</a-tag>
          <a-tag v-if="item.status" :color="statusColor(item.status)">{{ statusLabel(item.status) }}</a-tag>
        </a-space>
      </div>

      <a-space wrap size="small" style="margin-bottom: 12px">
        <span v-for="tag in item.tags" :key="tag" class="pill-meta">{{ tag }}</span>
        <span v-if="item.robot_type" class="pill-meta">{{ item.robot_type }}</span>
      </a-space>

      <div class="card-foot">
        <span class="meta-line">下载 {{ item.downloads }}</span>
        <span class="meta-line">更新 {{ formatDate(item.updated_at) }}</span>
      </div>

      <template #actions>
        <RouterLink :to="detailRoute(item)">查看详情</RouterLink>
        <RouterLink v-if="props.editable" :to="editRoute(item)">编辑资源</RouterLink>
      </template>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router';

import type { ResourceCard } from '@/types/api';

const props = defineProps<{
  items: ResourceCard[];
  editable?: boolean;
}>();

function typeLabel(type: string) {
  if (type === 'model') return '模型';
  if (type === 'dataset') return '数据集';
  return type;
}

function statusLabel(status?: string) {
  return {
    draft: '草稿',
    pending: '待审核',
    published: '已发布',
    rejected: '已驳回',
  }[status ?? ''] ?? status;
}

function statusColor(status?: string) {
  return {
    draft: 'default',
    pending: 'gold',
    published: 'green',
    rejected: 'red',
  }[status ?? ''] ?? 'default';
}

function detailRoute(item: ResourceCard) {
  return item.type === 'dataset' ? `/datasets/${item.id}` : `/models/${item.id}`;
}

function editRoute(item: ResourceCard) {
  return item.type === 'dataset' ? `/datasets/${item.id}/edit` : `/models/${item.id}/edit`;
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString();
}
</script>

<style scoped>
.resource-card {
  display: flex;
  height: 100%;
  flex-direction: column;
  border-radius: 22px;
  background: linear-gradient(180deg, #ffffff, #f8fbff);
  box-shadow: var(--shadow);
}

.resource-card :deep(.ant-card-body) {
  display: flex;
  flex: 1;
  flex-direction: column;
}

.card-top {
  display: flex;
  gap: 12px;
  justify-content: space-between;
  margin-bottom: 12px;
}

.card-title {
  margin: 0 0 6px;
  font-size: 18px;
}

.card-summary {
  margin: 0;
  min-height: 44px;
  color: var(--text-secondary);
}

.card-foot {
  margin-top: auto;
  display: flex;
  justify-content: space-between;
  color: var(--text-secondary);
  font-size: 12px;
}
</style>

<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">审核中心</h1>
          <p class="section-subtitle">聚焦模型和数据集的发布审核链路。</p>
        </div>
      </div>

      <a-radio-group v-model:value="reviewType" @change="load">
        <a-radio-button value="models">模型</a-radio-button>
        <a-radio-button value="datasets">数据集</a-radio-button>
      </a-radio-group>

      <a-list :data-source="items" style="margin-top: 18px">
        <template #renderItem="{ item }">
          <a-list-item class="review-item">
            <a-list-item-meta :title="item.title" :description="`${item.summary} · 提交人：${item.owner}`" />
            <template #actions>
              <a-button type="link" @click="decide(item.id, 'approved')">通过</a-button>
              <a-button danger type="link" @click="decide(item.id, 'rejected')">驳回</a-button>
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
import type { ReviewItem } from '@/types/api';

const reviewType = ref<'models' | 'datasets'>('models');
const items = ref<ReviewItem[]>([]);

async function load() {
  items.value = await api.getAdminReviews(reviewType.value);
}

async function decide(id: number, decision: 'approved' | 'rejected') {
  await api.decideReview(reviewType.value, id, { decision, comment: decision === 'approved' ? 'approved in admin panel' : 'rejected in admin panel' });
  message.success('审核状态已更新');
  await load();
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.review-item {
  border-radius: 16px;
  border: 1px solid var(--line);
  margin-bottom: 12px;
  padding: 8px 12px;
}
</style>

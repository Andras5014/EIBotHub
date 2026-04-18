<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">认证审核</h1>
          <p class="section-subtitle">审核开发者个人和企业认证申请，并将结果写回用户通知中心。</p>
        </div>
      </div>

      <a-list :data-source="items">
        <template #renderItem="{ item }">
          <a-list-item class="verify-item">
            <a-list-item-meta>
              <template #title>
                <a-space wrap>
                  <span>{{ item.user_name }}</span>
                  <a-tag color="blue">{{ item.verification_type }}</a-tag>
                  <a-tag :color="item.status === 'approved' ? 'green' : item.status === 'rejected' ? 'red' : 'gold'">
                    {{ item.status }}
                  </a-tag>
                </a-space>
              </template>
              <template #description>
                <div>真实姓名：{{ item.real_name }}</div>
                <div>组织：{{ item.organization || '未填写' }}</div>
                <div>申请理由：{{ item.reason }}</div>
              </template>
            </a-list-item-meta>
            <template #actions>
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
import type { AdminVerificationItem } from '@/types/api';

const items = ref<AdminVerificationItem[]>([]);

async function load() {
  items.value = await api.getAdminVerifications();
}

async function review(id: number, decision: 'approved' | 'rejected') {
  await api.reviewAdminVerification(id, {
    decision,
    comment: decision === 'approved' ? '认证通过' : '材料不足或不符合要求',
  });
  message.success('认证状态已更新');
  await load();
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.verify-item {
  border: 1px solid var(--line);
  border-radius: 18px;
  margin-bottom: 12px;
  background: var(--surface-soft);
}
</style>

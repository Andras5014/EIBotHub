<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">后台总览</h1>
          <p class="section-subtitle">MVP 先覆盖资源审核、公告运营和基础统计。</p>
        </div>
        <a-space>
          <RouterLink to="/admin/portal"><a-button type="primary">门户运营</a-button></RouterLink>
          <RouterLink to="/admin/reviews"><a-button>审核中心</a-button></RouterLink>
          <RouterLink to="/admin/dataset-access"><a-button>下载审批</a-button></RouterLink>
          <RouterLink to="/admin/content"><a-button>内容管理</a-button></RouterLink>
          <RouterLink to="/admin/announcements"><a-button type="primary">公告管理</a-button></RouterLink>
          <RouterLink to="/admin/community"><a-button>社区治理</a-button></RouterLink>
          <RouterLink to="/admin/verifications"><a-button>认证审核</a-button></RouterLink>
          <RouterLink to="/admin/rewards"><a-button>积分运营</a-button></RouterLink>
        </a-space>
      </div>

      <div class="grid-cards" v-if="dashboard">
        <a-card><a-statistic title="用户数" :value="dashboard.users" /></a-card>
        <a-card><a-statistic title="已发布模型" :value="dashboard.published_models" /></a-card>
        <a-card><a-statistic title="待审模型" :value="dashboard.pending_models" /></a-card>
        <a-card><a-statistic title="已发布数据集" :value="dashboard.published_datasets" /></a-card>
        <a-card><a-statistic title="待审数据集" :value="dashboard.pending_datasets" /></a-card>
        <a-card><a-statistic title="公告数" :value="dashboard.announcements" /></a-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import type { DashboardPayload } from '@/types/api';

const dashboard = ref<DashboardPayload>();

onMounted(async () => {
  dashboard.value = await api.getAdminDashboard();
});
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

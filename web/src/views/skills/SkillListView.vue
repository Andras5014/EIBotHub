<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">技能分享</h1>
          <p class="section-subtitle">V1 新增用户技能分享、Fork、评分与评论，便于任务流复用。</p>
        </div>
        <RouterLink to="/skills/new">
          <a-button type="primary">发布技能</a-button>
        </RouterLink>
      </div>

      <a-list :grid="listGrid" :data-source="items">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-card :title="item.name" class="skill-card">
              <p class="skill-summary">{{ item.summary }}</p>
              <a-space wrap style="margin-bottom: 10px">
                <span class="pill-meta">{{ item.category }}</span>
                <span class="pill-meta">{{ item.scene }}</span>
                <span class="pill-meta">评分 {{ item.rating.average.toFixed(1) }}</span>
              </a-space>
              <div class="section-subtitle skill-owner">作者：{{ item.owner_name }}</div>
              <template #actions>
                <RouterLink :to="`/skills/${item.id}`">查看详情</RouterLink>
              </template>
            </a-card>
          </a-list-item>
        </template>
      </a-list>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import type { SkillItem } from '@/types/api';

const items = ref<SkillItem[]>([]);
const listGrid = {
  gutter: 16,
  xs: 1,
  sm: 1,
  md: 2,
  lg: 3,
  xl: 3,
};

onMounted(async () => {
  items.value = await api.listSkills();
});
</script>

<style scoped>
.block {
  padding: 24px;
}

.skill-card {
  height: 100%;
  border-radius: 18px;
}

.skill-card :deep(.ant-card-body) {
  min-height: 178px;
}

.skill-summary {
  min-height: 44px;
  margin-bottom: 12px;
  color: var(--text-main);
  overflow-wrap: anywhere;
}

.skill-owner {
  overflow-wrap: anywhere;
}
</style>

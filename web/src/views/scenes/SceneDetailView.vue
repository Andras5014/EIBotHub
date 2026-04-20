<template>
  <div class="page-shell" v-if="detail">
    <section class="page-card hero-panel">
      <div>
        <a-tag color="cyan">{{ detail.scene.tagline || '场景专题' }}</a-tag>
        <h1 class="section-title hero-title">{{ detail.scene.name }}</h1>
        <p class="section-subtitle hero-subtitle">{{ detail.scene.summary }}</p>
        <p class="hero-description">{{ detail.scene.description }}</p>
      </div>
      <div class="hero-stats">
        <a-card class="hero-stat-card" :bordered="false"><a-statistic title="模型" :value="detail.models.length" /></a-card>
        <a-card class="hero-stat-card" :bordered="false"><a-statistic title="数据集" :value="detail.datasets.length" /></a-card>
        <a-card class="hero-stat-card" :bordered="false"><a-statistic title="模板" :value="detail.task_templates.length" /></a-card>
        <a-card class="hero-stat-card" :bordered="false"><a-statistic title="案例" :value="detail.application_cases.length" /></a-card>
      </div>
    </section>

    <section class="page-card block">
      <div class="section-head">
        <div>
          <h2 class="section-title">场景资源</h2>
          <p class="section-subtitle">按当前场景聚合模型、数据集、模板和案例。</p>
        </div>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="12">
          <a-card title="模型" class="inner-card">
            <ResourceCards :items="detail.models" />
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="数据集" class="inner-card">
            <ResourceCards :items="detail.datasets" />
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="任务模板" class="inner-card">
            <a-empty v-if="!detail.task_templates.length" description="当前场景下没有模板" />
            <a-list v-else :data-source="detail.task_templates">
              <template #renderItem="{ item }">
                <a-list-item>
                  <RouterLink :to="`/templates/${item.id}`">{{ item.name }}</RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="12">
          <a-card title="具身案例" class="inner-card">
            <a-empty v-if="!detail.application_cases.length" description="当前场景下没有案例" />
            <a-list v-else :data-source="detail.application_cases">
              <template #renderItem="{ item }">
                <a-list-item>
                  <RouterLink :to="`/applications/${item.id}`">{{ item.title }}</RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

import { api } from '@/api';
import ResourceCards from '@/components/ResourceCards.vue';
import type { ScenePageDetailPayload } from '@/types/api';

const route = useRoute();
const detail = ref<ScenePageDetailPayload>();

onMounted(async () => {
  detail.value = await api.getScenePage(String(route.params.slug));
});
</script>

<style scoped>
.hero-panel {
  padding: 28px;
  display: grid;
  gap: 20px;
  grid-template-columns: 1.4fr 1fr;
  margin-bottom: 20px;
}

.hero-title {
  margin-top: 12px;
}

.hero-subtitle,
.hero-description {
  line-height: 1.75;
}

.hero-description {
  color: var(--text-secondary);
}

.hero-stats {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.hero-stat-card {
  border-radius: 20px;
  background: linear-gradient(180deg, #fbfdff, #eef7fb);
}

.block {
  padding: 24px;
}

.inner-card {
  border-radius: 18px;
}

@media (max-width: 960px) {
  .hero-panel,
  .hero-stats {
    grid-template-columns: 1fr;
  }
}
</style>

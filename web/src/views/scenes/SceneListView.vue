<template>
  <div class="page-shell">
    <section class="page-card hero-panel">
      <div>
        <a-tag color="blue">Scene Hub</a-tag>
        <h1 class="section-title hero-title">场景专区</h1>
        <p class="section-subtitle hero-subtitle">
          把模型、数据集、模板和案例围绕具体场景重新组织，先看落地目标，再反查所需资源。
        </p>
      </div>
      <div class="hero-stats">
        <a-card class="hero-stat-card" :bordered="false">
          <a-statistic title="场景卡片" :value="items.length" />
        </a-card>
      </div>
    </section>

    <section class="page-card block">
      <div class="section-head">
        <div>
          <h2 class="section-title">已配置场景</h2>
          <p class="section-subtitle">这些场景由后台运营配置，可作为首页导流和场景化运营入口。</p>
        </div>
      </div>

      <a-empty v-if="!items.length" description="还没有配置场景页" />
      <div v-else class="scene-grid">
        <RouterLink v-for="item in items" :key="item.id" :to="`/scenes/${item.slug}`" class="scene-card">
          <div class="scene-kicker">{{ item.tagline || '场景专题' }}</div>
          <div class="scene-title">{{ item.name }}</div>
          <div class="scene-summary">{{ item.summary }}</div>
        </RouterLink>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import type { ScenePageItem } from '@/types/api';

const items = ref<ScenePageItem[]>([]);

onMounted(async () => {
  items.value = await api.listScenePages();
});
</script>

<style scoped>
.hero-panel {
  padding: 28px;
  display: grid;
  gap: 20px;
  grid-template-columns: 1.5fr 1fr;
  margin-bottom: 20px;
}

.hero-title {
  margin-top: 12px;
}

.hero-subtitle {
  max-width: 760px;
  line-height: 1.75;
}

.hero-stat-card {
  border-radius: 20px;
  background: linear-gradient(180deg, #fbfdff, #eef7fb);
}

.block {
  padding: 24px;
}

.scene-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.scene-card {
  display: block;
  padding: 20px;
  border-radius: 20px;
  border: 1px solid rgba(220, 230, 242, 0.82);
  background: linear-gradient(180deg, #fff 0%, #f9fbff 100%);
}

.scene-kicker {
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.scene-title {
  margin-top: 10px;
  color: var(--text-main);
  font-size: 22px;
  font-weight: 700;
}

.scene-summary {
  margin-top: 10px;
  color: var(--text-secondary);
  line-height: 1.7;
}

@media (max-width: 960px) {
  .hero-panel,
  .scene-grid {
    grid-template-columns: 1fr;
  }
}
</style>

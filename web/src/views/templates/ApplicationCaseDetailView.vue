<template>
  <div class="page-shell">
    <a-spin :spinning="loading">
      <a-result
        v-if="error"
        status="warning"
        title="案例详情加载失败"
        :sub-title="error"
      >
        <template #extra>
          <a-space>
            <RouterLink to="/applications">
              <a-button>返回案例库</a-button>
            </RouterLink>
            <a-button type="primary" @click="load">重试</a-button>
          </a-space>
        </template>
      </a-result>

      <template v-else-if="detail">
        <section class="page-card hero-panel">
          <div>
            <a-space wrap>
              <span class="pill-meta">{{ detail.category }}</span>
              <span class="pill-meta">{{ guideSteps.length }} 个部署步骤</span>
              <span class="pill-meta">更新于 {{ formatDate(detail.updated_at) }}</span>
            </a-space>
            <h1 class="section-title hero-title">{{ detail.title }}</h1>
            <p class="section-subtitle hero-subtitle">{{ detail.summary }}</p>
          </div>
          <a-space wrap>
            <RouterLink to="/applications">
              <a-button>返回列表</a-button>
            </RouterLink>
            <RouterLink :to="{ name: 'templates', query: { category: detail.category } }">
              <a-button>查看模板</a-button>
            </RouterLink>
            <RouterLink :to="{ name: 'search', query: { q: detail.title } }">
              <a-button type="primary">全局搜索资源</a-button>
            </RouterLink>
          </a-space>
        </section>

        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :lg="15">
            <a-card title="部署路径" class="content-card">
              <div v-if="guideSteps.length" class="step-list">
                <div v-for="(step, index) in guideSteps" :key="`${detail.id}-${index}`" class="step-item">
                  <span class="step-index">{{ index + 1 }}</span>
                  <div>
                    <div class="step-title">阶段 {{ index + 1 }}</div>
                    <div class="step-copy">{{ step }}</div>
                  </div>
                </div>
              </div>
              <a-empty v-else description="当前案例还没有配置详细部署步骤" />
            </a-card>

            <a-card title="案例解读" class="content-card" style="margin-top: 16px">
              <p class="content-copy">
                该案例聚焦于「{{ detail.category }}」方向，适合作为落地方案导读页。你可以先按步骤查看部署链路，再回到模板库筛选更具体的执行方案。
              </p>
            </a-card>
          </a-col>

          <a-col :xs="24" :lg="9">
            <a-card title="案例概览" class="content-card">
              <div class="overview-grid">
                <div class="overview-item">
                  <div class="overview-label">业务分类</div>
                  <div class="overview-value">{{ detail.category }}</div>
                </div>
                <div class="overview-item">
                  <div class="overview-label">部署阶段数</div>
                  <div class="overview-value">{{ guideSteps.length }}</div>
                </div>
                <div class="overview-item">
                  <div class="overview-label">最近更新</div>
                  <div class="overview-value">{{ formatDate(detail.updated_at) }}</div>
                </div>
              </div>
            </a-card>

            <a-card title="推荐动作" class="content-card" style="margin-top: 16px">
              <div class="action-list">
                <RouterLink :to="{ name: 'templates', query: { category: detail.category } }" class="action-link">
                  <span class="action-title">筛选同类任务模板</span>
                  <span class="action-copy">按当前业务分类回查已沉淀的流程模板</span>
                </RouterLink>
                <RouterLink :to="{ name: 'search', query: { q: detail.category } }" class="action-link">
                  <span class="action-title">搜索相关模型与数据集</span>
                  <span class="action-copy">查看该方向已有的资源组合与说明文档</span>
                </RouterLink>
              </div>
            </a-card>
          </a-col>
        </a-row>

        <a-card
          v-if="relatedTemplates.length"
          title="可直接复用的任务模板"
          class="content-card"
          style="margin-top: 16px"
        >
          <div class="related-grid">
            <a-card v-for="item in relatedTemplates" :key="item.id" hoverable class="related-card">
              <template #title>{{ item.name }}</template>
              <p class="card-summary">{{ item.summary }}</p>
              <a-space wrap>
                <span class="pill-meta">{{ item.category }}</span>
                <span class="pill-meta">{{ item.scene }}</span>
                <span class="pill-meta">使用 {{ item.usage_count }}</span>
              </a-space>
              <template #actions>
                <RouterLink :to="`/templates/${item.id}`">查看模板</RouterLink>
              </template>
            </a-card>
          </div>
        </a-card>
      </template>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

import { api } from '@/api';
import type { ApplicationCaseItem, TaskTemplateItem } from '@/types/api';

import { formatDate, splitGuideSteps } from './utils';

const route = useRoute();
const detail = ref<ApplicationCaseItem>();
const relatedTemplates = ref<TaskTemplateItem[]>([]);
const loading = ref(false);
const error = ref('');

const guideSteps = computed(() => splitGuideSteps(detail.value?.guide));

watch(
  () => route.params.id,
  () => {
    void load();
  },
  { immediate: true },
);

async function load() {
  loading.value = true;
  error.value = '';

  try {
    const id = String(route.params.id);
    const [caseDetail, templates] = await Promise.all([
      api.getApplicationCase(id),
      api.listTemplates(),
    ]);

    detail.value = caseDetail;
    relatedTemplates.value = templates
      .filter((item) => item.category === caseDetail.category)
      .sort((left, right) => right.usage_count - left.usage_count)
      .slice(0, 3);
  } catch (loadError) {
    error.value = loadError instanceof Error ? loadError.message : '案例详情获取失败';
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.hero-panel {
  padding: 28px;
  display: flex;
  justify-content: space-between;
  gap: 20px;
  align-items: flex-start;
  margin-bottom: 16px;
}

.hero-title {
  margin-top: 14px;
}

.hero-subtitle {
  max-width: 760px;
}

.content-card {
  border-radius: 22px;
}

.content-copy {
  margin: 0;
  color: var(--text-secondary);
  white-space: pre-wrap;
}

.step-list {
  display: grid;
  gap: 12px;
}

.step-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 14px;
  border-radius: 18px;
  background: var(--surface-soft);
  border: 1px solid var(--line);
}

.step-index {
  width: 28px;
  height: 28px;
  display: inline-grid;
  place-items: center;
  border-radius: 999px;
  background: rgba(19, 163, 127, 0.14);
  color: #0f7b61;
  font-weight: 700;
  flex: none;
}

.step-title {
  font-weight: 700;
  margin-bottom: 4px;
}

.step-copy {
  color: var(--text-secondary);
}

.overview-grid {
  display: grid;
  gap: 12px;
}

.overview-item {
  padding: 14px;
  border-radius: 18px;
  background: var(--surface-soft);
  border: 1px solid var(--line);
}

.overview-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.overview-value {
  font-size: 16px;
  font-weight: 700;
}

.action-list {
  display: grid;
  gap: 12px;
}

.action-link {
  display: grid;
  gap: 4px;
  padding: 14px;
  border-radius: 18px;
  background: var(--surface-soft);
  border: 1px solid var(--line);
}

.action-title {
  font-weight: 700;
}

.action-copy {
  color: var(--text-secondary);
}

.related-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.related-card {
  border-radius: 18px;
}

.card-summary {
  color: var(--text-secondary);
  min-height: 44px;
}

@media (max-width: 960px) {
  .hero-panel {
    flex-direction: column;
  }

  .related-grid {
    grid-template-columns: 1fr;
  }
}
</style>

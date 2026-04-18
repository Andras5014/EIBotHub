<template>
  <div class="page-shell">
    <section class="page-card hero-panel">
      <div>
        <a-tag color="cyan">具身案例</a-tag>
        <h1 class="section-title hero-title">从可复用案例反推部署路径</h1>
        <p class="section-subtitle hero-subtitle">
          用案例页承接模板、模型和数据集之间的关系，让用户先看到落地方式，再反查所需资源。
        </p>
      </div>
      <div class="hero-stats">
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="案例总数" :value="items.length" />
        </a-card>
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="覆盖分类" :value="categories.length" />
        </a-card>
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="部署阶段" :value="totalSteps" />
        </a-card>
      </div>
    </section>

    <section class="page-card block">
      <div class="section-head">
        <div>
          <h2 class="section-title">案例目录</h2>
          <p class="section-subtitle">按行业分类和部署步骤快速筛选，直接跳到可落地的实践方案。</p>
        </div>
        <span class="pill-meta">命中 {{ filteredItems.length }} 项</span>
      </div>

      <div class="toolbar">
        <a-input-search
          v-model:value="filters.q"
          class="toolbar-search"
          allow-clear
          placeholder="搜索案例标题、摘要或部署步骤"
        />
        <a-select v-model:value="filters.category" class="toolbar-select">
          <a-select-option value="">全部分类</a-select-option>
          <a-select-option v-for="item in categories" :key="item" :value="item">{{ item }}</a-select-option>
        </a-select>
        <a-select v-model:value="filters.sort" class="toolbar-select">
          <a-select-option value="latest">最近更新</a-select-option>
          <a-select-option value="steps">步骤数量</a-select-option>
          <a-select-option value="title">标题顺序</a-select-option>
        </a-select>
      </div>

      <a-spin :spinning="loading">
        <a-result
          v-if="error"
          status="warning"
          title="案例加载失败"
          :sub-title="error"
        >
          <template #extra>
            <a-button type="primary" @click="load">重新加载</a-button>
          </template>
        </a-result>

        <template v-else-if="filteredItems.length">
          <div class="case-grid">
            <a-card v-for="item in pagedItems" :key="item.id" hoverable class="case-card">
              <template #title>
                <div class="card-title-row">
                  <span class="card-title">{{ item.title }}</span>
                  <span class="card-date">{{ formatDate(item.updated_at) }}</span>
                </div>
              </template>

              <a-space wrap style="margin-bottom: 12px">
                <span class="pill-meta">{{ item.category }}</span>
                <span class="pill-meta">{{ splitGuideSteps(item.guide).length }} 个步骤</span>
              </a-space>

              <p class="card-summary">{{ item.summary }}</p>

              <div class="guide-panel">
                <div class="guide-title">部署路径预览</div>
                <div v-for="(step, index) in splitGuideSteps(item.guide).slice(0, 4)" :key="`${item.id}-${index}`" class="guide-step">
                  <span class="guide-index">{{ index + 1 }}</span>
                  <span>{{ step }}</span>
                </div>
              </div>

              <template #actions>
                <RouterLink :to="`/applications/${item.id}`">查看部署详情</RouterLink>
                <RouterLink :to="{ name: 'templates', query: { category: item.category } }">查看关联模板</RouterLink>
              </template>
            </a-card>
          </div>

          <div class="pagination-wrap">
            <a-pagination
              v-model:current="currentPage"
              :page-size="pageSize"
              :total="filteredItems.length"
              :show-size-changer="false"
            />
          </div>
        </template>

        <a-empty v-else description="当前筛选条件下没有匹配案例" />
      </a-spin>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

import { api } from '@/api';
import type { ApplicationCaseItem } from '@/types/api';

import { formatDate, splitGuideSteps, uniqueValues } from './utils';

const route = useRoute();
const items = ref<ApplicationCaseItem[]>([]);
const loading = ref(false);
const error = ref('');
const currentPage = ref(1);
const pageSize = 6;
const filters = reactive({
  q: '',
  category: String(route.query.category ?? ''),
  sort: 'latest',
});

const categories = computed(() => uniqueValues(items.value.map((item) => item.category)));
const totalSteps = computed(() => items.value.reduce((sum, item) => sum + splitGuideSteps(item.guide).length, 0));

const filteredItems = computed(() => {
  const keyword = filters.q.trim().toLowerCase();

  const next = items.value.filter((item) => {
    const matchesKeyword =
      !keyword || [item.title, item.summary, item.guide, item.category].join(' ').toLowerCase().includes(keyword);
    const matchesCategory = !filters.category || item.category === filters.category;
    return matchesKeyword && matchesCategory;
  });

  return next.sort((left, right) => {
    if (filters.sort === 'steps') {
      return splitGuideSteps(right.guide).length - splitGuideSteps(left.guide).length;
    }
    if (filters.sort === 'title') {
      return left.title.localeCompare(right.title, 'zh-CN');
    }
    return new Date(right.updated_at).getTime() - new Date(left.updated_at).getTime();
  });
});

const pagedItems = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  return filteredItems.value.slice(start, start + pageSize);
});

watch(
  () => [filters.q, filters.category, filters.sort],
  () => {
    currentPage.value = 1;
  },
);

watch(
  () => route.query.category,
  (value) => {
    filters.category = String(value ?? '');
  },
);

async function load() {
  loading.value = true;
  error.value = '';

  try {
    items.value = await api.listApplicationCases();
  } catch (loadError) {
    error.value = loadError instanceof Error ? loadError.message : '案例列表获取失败';
  } finally {
    loading.value = false;
  }
}

onMounted(load);
</script>

<style scoped>
.hero-panel {
  padding: 28px;
  display: grid;
  gap: 24px;
  grid-template-columns: 1.5fr 1fr;
  margin-bottom: 20px;
}

.hero-title {
  margin-top: 14px;
}

.hero-subtitle {
  max-width: 720px;
}

.hero-stats {
  display: grid;
  gap: 12px;
}

.hero-stat-card {
  border-radius: 20px;
  background: linear-gradient(180deg, #fbfdff, #eef7fb);
}

.block {
  padding: 24px;
}

.toolbar {
  display: grid;
  gap: 12px;
  grid-template-columns: minmax(260px, 1.8fr) repeat(2, minmax(180px, 1fr));
  margin-bottom: 22px;
}

.toolbar-search,
.toolbar-select {
  width: 100%;
}

.case-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.case-card {
  border-radius: 22px;
}

.card-title-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.card-title {
  font-weight: 700;
}

.card-date {
  font-size: 12px;
  color: var(--text-secondary);
}

.card-summary {
  min-height: 48px;
  color: var(--text-secondary);
  margin-bottom: 14px;
}

.guide-panel {
  padding: 14px;
  border-radius: 18px;
  background: var(--surface-soft);
  border: 1px solid var(--line);
}

.guide-title {
  font-weight: 700;
  margin-bottom: 10px;
}

.guide-step {
  display: flex;
  gap: 10px;
  align-items: flex-start;
  margin-bottom: 8px;
  color: var(--text-secondary);
}

.guide-step:last-child {
  margin-bottom: 0;
}

.guide-index {
  width: 22px;
  height: 22px;
  display: inline-grid;
  place-items: center;
  border-radius: 999px;
  background: rgba(19, 163, 127, 0.14);
  color: #0f7b61;
  font-size: 12px;
  font-weight: 700;
  flex: none;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

@media (max-width: 960px) {
  .hero-panel {
    grid-template-columns: 1fr;
  }

  .toolbar {
    grid-template-columns: 1fr;
  }

  .case-grid {
    grid-template-columns: 1fr;
  }
}
</style>

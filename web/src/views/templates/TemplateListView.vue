<template>
  <div class="page-shell">
    <section class="page-card hero-panel">
      <div>
        <a-tag color="blue">Plan / 任务模板</a-tag>
        <h1 class="section-title hero-title">把已有流程模板化，降低任务部署成本</h1>
        <p class="section-subtitle hero-subtitle">
          参考现有模型与数据集页的信息组织方式，把模板库补成可筛选、可对比、可快速跳转的标准入口。
        </p>
      </div>
      <div class="hero-stats">
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="模板总数" :value="items.length" />
        </a-card>
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="覆盖场景" :value="scenes.length" />
        </a-card>
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="累计使用" :value="totalUsage" />
        </a-card>
      </div>
    </section>

    <section class="page-card block">
      <div class="section-head">
        <div>
          <h2 class="section-title">模板目录</h2>
          <p class="section-subtitle">按关键词、分类和场景快速定位合适的执行方案。</p>
        </div>
        <span class="pill-meta">命中 {{ filteredItems.length }} 项</span>
      </div>

      <div class="toolbar">
        <a-input-search
          v-model:value="filters.q"
          class="toolbar-search"
          allow-clear
          placeholder="搜索模板名称、摘要、说明或关联资源"
        />
        <a-select v-model:value="filters.category" class="toolbar-select">
          <a-select-option value="">全部分类</a-select-option>
          <a-select-option v-for="item in categories" :key="item" :value="item">{{ item }}</a-select-option>
        </a-select>
        <a-select v-model:value="filters.scene" class="toolbar-select">
          <a-select-option value="">全部场景</a-select-option>
          <a-select-option v-for="item in scenes" :key="item" :value="item">{{ item }}</a-select-option>
        </a-select>
        <a-select v-model:value="filters.sort" class="toolbar-select">
          <a-select-option value="usage">按热度排序</a-select-option>
          <a-select-option value="latest">最近更新</a-select-option>
          <a-select-option value="name">名称顺序</a-select-option>
        </a-select>
      </div>

      <a-spin :spinning="loading">
        <a-result
          v-if="error"
          status="warning"
          title="模板加载失败"
          :sub-title="error"
        >
          <template #extra>
            <a-button type="primary" @click="load">重新加载</a-button>
          </template>
        </a-result>

        <template v-else-if="filteredItems.length">
          <div class="template-grid">
            <a-card v-for="item in pagedItems" :key="item.id" hoverable class="template-card">
              <template #title>
                <div class="card-title-row">
                  <span class="card-title">{{ item.name }}</span>
                  <span class="card-date">{{ formatDate(item.updated_at) }}</span>
                </div>
              </template>

              <a-space wrap style="margin-bottom: 12px">
                <span class="pill-meta">{{ item.category }}</span>
                <span class="pill-meta">{{ item.scene }}</span>
                <span class="pill-meta">使用 {{ item.usage_count }}</span>
              </a-space>

              <p class="card-summary">{{ item.summary }}</p>

              <div class="guide-panel">
                <div class="guide-title">执行步骤</div>
                <div v-for="(step, index) in previewSteps(item.guide)" :key="`${item.id}-${index}`" class="guide-step">
                  <span class="guide-index">{{ index + 1 }}</span>
                  <span>{{ step }}</span>
                </div>
              </div>

              <div class="resource-panel">
                <div class="guide-title">关联资源</div>
                <a-space wrap>
                  <RouterLink
                    v-for="resource in item.resource_ref.slice(0, 3)"
                    :key="resource"
                    :to="{ name: 'search', query: { q: resource } }"
                  >
                    <span class="pill-meta resource-chip">{{ resource }}</span>
                  </RouterLink>
                  <span v-if="!item.resource_ref.length" class="empty-text">暂无关联资源</span>
                </a-space>
              </div>

              <template #actions>
                <RouterLink :to="`/templates/${item.id}`">查看模板详情</RouterLink>
                <RouterLink :to="{ name: 'search', query: { q: item.name, type: 'task-template' } }">进入全局搜索</RouterLink>
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

        <a-empty
          v-else
          description="当前筛选条件下没有匹配的任务模板"
        />
      </a-spin>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

import { api } from '@/api';
import type { TaskTemplateItem } from '@/types/api';

import { formatDate, splitGuideSteps, uniqueValues } from './utils';

const route = useRoute();
const items = ref<TaskTemplateItem[]>([]);
const loading = ref(false);
const error = ref('');
const currentPage = ref(1);
const pageSize = 6;
const filters = reactive({
  q: String(route.query.q ?? ''),
  category: String(route.query.category ?? ''),
  scene: String(route.query.scene ?? ''),
  sort: 'usage',
});

const categories = computed(() => uniqueValues(items.value.map((item) => item.category)));
const scenes = computed(() => uniqueValues(items.value.map((item) => item.scene)));
const totalUsage = computed(() => items.value.reduce((sum, item) => sum + item.usage_count, 0));

const filteredItems = computed(() => {
  const keyword = filters.q.trim().toLowerCase();

  const next = items.value.filter((item) => {
    const matchesKeyword =
      !keyword ||
      [item.name, item.summary, item.description, item.guide, item.resource_ref.join(' ')]
        .join(' ')
        .toLowerCase()
        .includes(keyword);

    const matchesCategory = !filters.category || item.category === filters.category;
    const matchesScene = !filters.scene || item.scene === filters.scene;

    return matchesKeyword && matchesCategory && matchesScene;
  });

  return next.sort((left, right) => {
    if (filters.sort === 'latest') {
      return new Date(right.updated_at).getTime() - new Date(left.updated_at).getTime();
    }
    if (filters.sort === 'name') {
      return left.name.localeCompare(right.name, 'zh-CN');
    }
    return right.usage_count - left.usage_count;
  });
});

const pagedItems = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  return filteredItems.value.slice(start, start + pageSize);
});

watch(
  () => [filters.q, filters.category, filters.scene, filters.sort],
  () => {
    currentPage.value = 1;
  },
);

watch(
  () => [route.query.q, route.query.category, route.query.scene],
  ([q, category, scene]) => {
    filters.q = String(q ?? '');
    filters.category = String(category ?? '');
    filters.scene = String(scene ?? '');
  },
);

async function load() {
  loading.value = true;
  error.value = '';

  try {
    items.value = await api.listTemplates();
  } catch (loadError) {
    error.value = loadError instanceof Error ? loadError.message : '模板列表获取失败';
  } finally {
    loading.value = false;
  }
}

function previewSteps(guide: string) {
  return splitGuideSteps(guide).slice(0, 3);
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
  background: linear-gradient(180deg, #f9fbff, #eef4fb);
}

.block {
  padding: 24px;
}

.toolbar {
  display: grid;
  gap: 12px;
  grid-template-columns: minmax(260px, 1.8fr) repeat(3, minmax(160px, 1fr));
  margin-bottom: 22px;
}

.toolbar-search,
.toolbar-select {
  width: 100%;
}

.template-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.template-card {
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

.guide-panel,
.resource-panel {
  padding: 14px;
  border-radius: 18px;
  background: var(--surface-soft);
  border: 1px solid var(--line);
}

.resource-panel {
  margin-top: 12px;
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
  background: rgba(22, 119, 255, 0.12);
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 700;
  flex: none;
}

.resource-chip {
  cursor: pointer;
}

.empty-text {
  color: var(--text-secondary);
  font-size: 13px;
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

  .template-grid {
    grid-template-columns: 1fr;
  }
}
</style>

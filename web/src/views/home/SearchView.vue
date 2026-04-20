<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">全局搜索</h1>
          <p class="section-subtitle">覆盖模型、数据集、文档、模板、技能、讨论与用户，并补齐排序和基础筛选能力。</p>
        </div>
        <span class="pill-meta">结果 {{ total }}</span>
      </div>

      <a-space direction="vertical" style="width: 100%" size="middle">
        <a-input-search
          v-model:value="keyword"
          enter-button="搜索"
          size="large"
          placeholder="搜索模型、数据集、文档、任务模板或开发者"
          @search="submit"
        />

        <a-radio-group v-model:value="type">
          <a-radio-button value="">全部</a-radio-button>
          <a-radio-button value="model">模型</a-radio-button>
          <a-radio-button value="dataset">数据集</a-radio-button>
          <a-radio-button value="task-template">模板</a-radio-button>
          <a-radio-button value="skill">技能</a-radio-button>
          <a-radio-button value="discussion">讨论</a-radio-button>
          <a-radio-button value="doc">文档</a-radio-button>
          <a-radio-button value="user">用户</a-radio-button>
        </a-radio-group>

        <div class="filter-grid">
          <a-select v-model:value="tagValues" mode="multiple" allow-clear placeholder="标签筛选">
            <a-select-option v-for="item in filterOptions.tags" :key="item" :value="item">{{ item }}</a-select-option>
          </a-select>
          <a-select v-model:value="robotType" allow-clear placeholder="适用机器人">
            <a-select-option v-for="item in filterOptions.robot_types" :key="item" :value="item">{{ item }}</a-select-option>
          </a-select>
          <a-select v-model:value="sort">
            <a-select-option value="hot">按热度</a-select-option>
            <a-select-option value="latest">按时间</a-select-option>
            <a-select-option value="name">按名称</a-select-option>
          </a-select>
          <a-select v-model:value="updatedWithin">
            <a-select-option :value="0">全部时间</a-select-option>
            <a-select-option :value="7">最近 7 天</a-select-option>
            <a-select-option :value="30">最近 30 天</a-select-option>
            <a-select-option :value="90">最近 90 天</a-select-option>
          </a-select>
          <a-button type="primary" @click="submit">应用筛选</a-button>
        </div>

        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :lg="8">
            <a-card title="热门搜索">
              <a-space wrap>
                <a-tag v-for="item in hotQueries" :key="item.query" color="blue" style="cursor: pointer" @click="applyHot(item.query)">
                  {{ item.query }} · {{ item.count }}
                </a-tag>
              </a-space>
            </a-card>
          </a-col>
          <a-col :xs="24" :lg="8">
            <a-card title="推荐搜索">
              <a-space wrap>
                <a-tag v-for="item in recommendedQueries" :key="item.query" color="gold" style="cursor: pointer" @click="applyHot(item.query)">
                  {{ item.query }}
                </a-tag>
              </a-space>
            </a-card>
          </a-col>
          <a-col :xs="24" :lg="8">
            <a-card title="最近搜索">
              <a-space wrap>
                <a-tag v-for="item in recentQueries" :key="item" style="cursor: pointer" @click="applyHot(item)">
                  {{ item }}
                </a-tag>
              </a-space>
            </a-card>
          </a-col>
        </a-row>

        <a-row v-if="keyword.trim()" :gutter="[16, 16]">
          <a-col :xs="24" :lg="9">
            <a-card title="结果概览">
              <a-empty v-if="!typeCounts.length" description="暂无聚合结果" />
              <a-list v-else :data-source="typeCounts" size="small">
                <template #renderItem="{ item }">
                  <a-list-item class="summary-row" @click="applyType(item.type)">
                    <span>{{ item.label }}</span>
                    <a-tag>{{ item.count }}</a-tag>
                  </a-list-item>
                </template>
              </a-list>
            </a-card>
          </a-col>
          <a-col :xs="24" :lg="15">
            <a-card title="相关关键词">
              <a-space wrap>
                <a-tag
                  v-for="item in suggestedQueries"
                  :key="item.query"
                  color="processing"
                  style="cursor: pointer"
                  @click="applyHot(item.query)"
                >
                  {{ item.query }}
                </a-tag>
              </a-space>
            </a-card>
          </a-col>
        </a-row>

        <a-alert v-if="error" type="warning" show-icon :message="error" />

        <a-spin :spinning="loading">
          <a-empty v-if="!loading && !results.length && keyword.trim()" description="没有匹配结果" />
          <a-empty v-else-if="!loading && !keyword.trim()" description="输入关键词后开始搜索" />

          <a-list v-else :data-source="results" item-layout="vertical">
            <template #renderItem="{ item }">
              <a-list-item class="search-item">
                <RouterLink :to="item.route" class="search-link">
                  <a-space direction="vertical" :size="6">
                    <a-space wrap>
                      <a-tag color="blue">{{ typeLabel(item.type) }}</a-tag>
                      <span class="search-title">{{ item.title }}</span>
                      <span class="search-date">更新于 {{ formatDate(item.updated_at) }}</span>
                    </a-space>
                    <div class="search-summary">{{ item.summary || '暂无摘要' }}</div>
                    <a-space wrap>
                      <span v-for="tag in item.tags ?? []" :key="tag" class="pill-meta">{{ tag }}</span>
                      <span v-if="item.score_hint" class="pill-meta">热度 {{ item.score_hint }}</span>
                    </a-space>
                  </a-space>
                </RouterLink>
              </a-list-item>
            </template>
          </a-list>
        </a-spin>

        <a-row v-if="keyword.trim() && (sameTypeItems.length || relatedItems.length)" :gutter="[16, 16]">
          <a-col :xs="24" :lg="12">
            <a-card :title="sameTypeTitle">
              <a-empty v-if="!sameTypeItems.length" description="暂无同类推荐" />
              <a-list v-else :data-source="sameTypeItems" size="small">
                <template #renderItem="{ item }">
                  <a-list-item>
                    <RouterLink :to="item.route">
                      <div class="list-title">{{ item.title }}</div>
                      <div class="list-desc">{{ item.summary }}</div>
                    </RouterLink>
                  </a-list-item>
                </template>
              </a-list>
            </a-card>
          </a-col>
          <a-col :xs="24" :lg="12">
            <a-card title="关联推荐">
              <a-empty v-if="!relatedItems.length" description="暂无关联推荐" />
              <a-list v-else :data-source="relatedItems" size="small">
                <template #renderItem="{ item }">
                  <a-list-item>
                    <RouterLink :to="item.route">
                      <div class="list-title">{{ item.title }}</div>
                      <div class="list-desc">{{ typeLabel(item.type) }} · {{ item.summary }}</div>
                    </RouterLink>
                  </a-list-item>
                </template>
              </a-list>
            </a-card>
          </a-col>
        </a-row>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import type { FilterOptionsResponse, SearchHotItem, SearchItem, SearchSuggestionItem, SearchTypeCountItem } from '@/types/api';

const route = useRoute();
const router = useRouter();
const keyword = ref(String(route.query.q ?? ''));
const type = ref(String(route.query.type ?? ''));
const tagValues = ref<string[]>(readCSVQuery(route.query.tags));
const robotType = ref(String(route.query.robot_type ?? ''));
const sort = ref(String(route.query.sort ?? 'hot'));
const updatedWithin = ref(Number(route.query.updated_within ?? 0));
const results = ref<SearchItem[]>([]);
const sameTypeItems = ref<SearchItem[]>([]);
const relatedItems = ref<SearchItem[]>([]);
const suggestedQueries = ref<SearchSuggestionItem[]>([]);
const typeCounts = ref<SearchTypeCountItem[]>([]);
const focusType = ref('');
const hotQueries = ref<SearchHotItem[]>([]);
const recommendedQueries = ref<SearchSuggestionItem[]>([]);
const recentQueries = ref<string[]>(readRecentQueries());
const loading = ref(false);
const error = ref('');
const total = ref(0);
const filterOptions = ref<FilterOptionsResponse>({
  tags: [],
  model_tags: [],
  dataset_tags: [],
  robot_types: [],
  dataset_scenes: [],
  template_categories: [],
  template_scenes: [],
  application_case_categories: [],
});

watch(
  () => route.fullPath,
  async () => {
    keyword.value = String(route.query.q ?? '');
    type.value = String(route.query.type ?? '');
    tagValues.value = readCSVQuery(route.query.tags);
    robotType.value = String(route.query.robot_type ?? '');
    sort.value = String(route.query.sort ?? 'hot');
    updatedWithin.value = Number(route.query.updated_within ?? 0);
    const [hotQueryItems, recommendedQueryItems] = await Promise.all([
      api.hotQueries(),
      api.recommendedQueries(),
    ]);
    hotQueries.value = hotQueryItems;
    recommendedQueries.value = recommendedQueryItems;

    if (!keyword.value.trim()) {
      results.value = [];
      total.value = 0;
      sameTypeItems.value = [];
      relatedItems.value = [];
      suggestedQueries.value = [];
      typeCounts.value = [];
      focusType.value = '';
      error.value = '';
      return;
    }

    loading.value = true;
    error.value = '';
    try {
      const response = await api.search({
        q: keyword.value,
        type: type.value || undefined,
        tags: tagValues.value.join(',') || undefined,
        robot_type: robotType.value || undefined,
        sort: sort.value,
        updated_within: updatedWithin.value || undefined,
      });
      results.value = response.items;
      total.value = response.total;
      sameTypeItems.value = response.same_type_items;
      relatedItems.value = response.related_items;
      suggestedQueries.value = response.suggested_queries;
      typeCounts.value = response.type_counts;
      focusType.value = response.focus_type ?? '';
    } catch (loadError) {
      results.value = [];
      total.value = 0;
      sameTypeItems.value = [];
      relatedItems.value = [];
      suggestedQueries.value = [];
      typeCounts.value = [];
      focusType.value = '';
      error.value = loadError instanceof Error ? loadError.message : '搜索失败';
    } finally {
      loading.value = false;
    }
  },
  { immediate: true },
);

function submit() {
  saveRecentQuery(keyword.value);
  recentQueries.value = readRecentQueries();
  router.push({
    name: 'search',
    query: {
      q: keyword.value || undefined,
      type: type.value || undefined,
      tags: tagValues.value.length ? tagValues.value.join(',') : undefined,
      robot_type: robotType.value || undefined,
      sort: sort.value || undefined,
      updated_within: updatedWithin.value || undefined,
    },
  });
}

function typeLabel(value: string) {
  return {
    model: '模型',
    dataset: '数据集',
    'task-template': '模板',
    skill: '技能',
    discussion: '讨论',
    doc: '文档',
    user: '用户',
  }[value] ?? value;
}

function applyHot(value: string) {
  keyword.value = value;
  submit();
}

function applyType(value: string) {
  type.value = value;
  submit();
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString('zh-CN');
}

const sameTypeTitle = computed(() => {
  return focusType.value ? `${typeLabel(focusType.value)}同类推荐` : '同类推荐';
});

function saveRecentQuery(value: string) {
  const trimmed = value.trim();
  if (!trimmed) return;
  const next = [trimmed, ...readRecentQueries().filter((item) => item !== trimmed)].slice(0, 8);
  localStorage.setItem('open-community-recent-searches', JSON.stringify(next));
}

function readRecentQueries() {
  const raw = localStorage.getItem('open-community-recent-searches');
  if (!raw) return [];
  try {
    return JSON.parse(raw) as string[];
  } catch {
    return [];
  }
}

function readCSVQuery(value: unknown) {
  return String(value ?? '')
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean);
}

onMounted(async () => {
  filterOptions.value = await api.getFilterOptions();
});
</script>

<style scoped>
.block {
  padding: 24px;
}

.filter-grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(5, minmax(0, 1fr));
}

.search-item {
  border-radius: 16px;
  border: 1px solid var(--line);
  margin-bottom: 12px;
  background: var(--surface-soft);
}

.search-link {
  display: block;
}

.search-title {
  font-size: 18px;
  font-weight: 700;
}

.search-date {
  color: var(--text-secondary);
  font-size: 12px;
}

.search-summary {
  color: var(--text-secondary);
}

.summary-row {
  cursor: pointer;
}

@media (max-width: 960px) {
  .filter-grid {
    grid-template-columns: 1fr;
  }
}
</style>

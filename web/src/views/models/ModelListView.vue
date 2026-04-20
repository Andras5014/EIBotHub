<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">模型仓库</h1>
          <p class="section-subtitle">面向搬运、巡检、移动作业等机器人场景的模型资源库，补齐视图切换、排序与分页。</p>
        </div>
        <a-space>
          <RouterLink to="/models/compare">
            <a-button>模型对比</a-button>
          </RouterLink>
          <RouterLink to="/models/upload">
            <a-button type="primary">上传模型</a-button>
          </RouterLink>
        </a-space>
      </div>

      <a-space direction="vertical" style="width: 100%" size="large">
        <div class="toolbar">
          <a-form :model="filters" layout="inline" @finish="applyFilters" class="toolbar-form">
            <a-form-item label="关键词" name="q" class="toolbar-item toolbar-item-keyword">
              <a-input v-model:value="filters.q" class="toolbar-control" placeholder="模型名 / 标签 / 描述" />
            </a-form-item>
            <a-form-item label="标签" name="tags" class="toolbar-item toolbar-item-tags">
              <a-select
                v-model:value="tagValues"
                mode="multiple"
                allow-clear
                show-search
                :filter-option="filterSelectOption"
                :options="modelTagOptions"
                :loading="filterLoading"
                :max-tag-count="'responsive'"
                class="toolbar-control"
                placeholder="选择或搜索标签"
              />
            </a-form-item>
            <a-form-item label="机器人类型" name="robot_type" class="toolbar-item toolbar-item-robot">
              <a-select
                v-model:value="filters.robot_type"
                allow-clear
                show-search
                :filter-option="filterSelectOption"
                :options="robotTypeOptions"
                :loading="filterLoading"
                class="toolbar-control"
                placeholder="选择或搜索机器人类型"
              />
            </a-form-item>
            <a-form-item label="排序" name="sort" class="toolbar-item toolbar-item-sort">
              <a-select v-model:value="filters.sort" class="toolbar-control">
                <a-select-option value="downloads">按热度</a-select-option>
                <a-select-option value="latest">最近更新</a-select-option>
                <a-select-option value="name">按名称</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item class="toolbar-item toolbar-submit">
              <a-button html-type="submit" type="primary" class="filter-button">筛选</a-button>
            </a-form-item>
          </a-form>

          <div class="toolbar-side">
            <a-radio-group v-model:value="viewMode" button-style="solid" class="view-mode-switch">
              <a-radio-button value="card">卡片</a-radio-button>
              <a-radio-button value="list">列表</a-radio-button>
            </a-radio-group>
          </div>
        </div>

        <a-alert v-if="error" type="warning" show-icon :message="error" />

        <a-spin :spinning="loading">
          <ResourceCards v-if="viewMode === 'card' && items.length" :items="items" />

          <a-list v-else-if="items.length" :data-source="items">
            <template #renderItem="{ item }">
              <a-list-item class="model-item">
                <a-list-item-meta>
                  <template #title>
                    <RouterLink :to="`/models/${item.id}`">{{ item.name }}</RouterLink>
                  </template>
                  <template #description>
                    <div class="item-summary">{{ item.summary }}</div>
                    <a-space wrap style="margin-top: 8px">
                      <span v-for="tag in item.tags" :key="tag" class="pill-meta">{{ tag }}</span>
                      <span v-if="item.robot_type" class="pill-meta">{{ item.robot_type }}</span>
                      <span class="pill-meta">下载 {{ item.downloads }}</span>
                      <span class="pill-meta">更新 {{ formatDate(item.updated_at) }}</span>
                    </a-space>
                  </template>
                </a-list-item-meta>
              </a-list-item>
            </template>
          </a-list>

          <a-empty v-else description="当前筛选条件下没有匹配模型" />
        </a-spin>

        <div class="pagination-wrap" v-if="total > filters.page_size">
          <a-pagination
            v-model:current="filters.page"
            :page-size="filters.page_size"
            :total="total"
            :show-size-changer="false"
            @change="handlePageChange"
          />
        </div>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import ResourceCards from '@/components/ResourceCards.vue';
import type { FilterOptionsResponse, ResourceCard } from '@/types/api';

const viewMode = ref<'card' | 'list'>('card');
const loading = ref(false);
const error = ref('');
const total = ref(0);
const items = ref<ResourceCard[]>([]);
const tagValues = ref<string[]>([]);
const filterLoading = ref(false);
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
const filters = reactive({
  q: '',
  tags: '',
  robot_type: '',
  sort: 'downloads',
  page: 1,
  page_size: 12,
});
const modelTagOptions = computed(() => filterOptions.value.model_tags.map((item) => ({ label: item, value: item })));
const robotTypeOptions = computed(() => filterOptions.value.robot_types.map((item) => ({ label: item, value: item })));

async function load() {
  loading.value = true;
  error.value = '';
  try {
    const response = await api.listModels({
      ...filters,
      tags: tagValues.value.join(','),
    });
    items.value = response.items;
    total.value = response.total;
  } catch (loadError) {
    items.value = [];
    total.value = 0;
    error.value = loadError instanceof Error ? loadError.message : '模型列表加载失败';
  } finally {
    loading.value = false;
  }
}

async function applyFilters() {
  filters.page = 1;
  await load();
}

async function loadFilterOptions() {
  filterLoading.value = true;
  try {
    filterOptions.value = await api.getFilterOptions();
  } catch (loadError) {
    filterOptions.value = {
      tags: [],
      model_tags: [],
      dataset_tags: [],
      robot_types: [],
      dataset_scenes: [],
      template_categories: [],
      template_scenes: [],
      application_case_categories: [],
    };
    error.value = loadError instanceof Error ? `筛选项加载失败：${loadError.message}` : '筛选项加载失败';
  } finally {
    filterLoading.value = false;
  }
}

async function handlePageChange(page: number) {
  filters.page = page;
  await load();
}

function filterSelectOption(input: string, option: { label?: string; value?: string | number } | undefined) {
  const keyword = input.trim().toLowerCase();
  if (!keyword) {
    return true;
  }
  return String(option?.label ?? option?.value ?? '')
    .toLowerCase()
    .includes(keyword);
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString('zh-CN');
}

onMounted(async () => {
  await loadFilterOptions();
  await load();
});
</script>

<style scoped>
.block {
  padding: 24px;
}

.toolbar {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 16px 20px;
  align-items: end;
}

.toolbar-form {
  display: grid;
  grid-template-columns:
    minmax(220px, 1.35fr)
    minmax(220px, 1.25fr)
    minmax(200px, 1fr)
    minmax(150px, 0.8fr)
    auto;
  gap: 12px 16px;
  min-width: 0;
}

.toolbar-form :deep(.ant-form-item) {
  margin-right: 0;
  margin-bottom: 0;
}

.toolbar-form :deep(.ant-form-item-control) {
  min-width: 0;
}

.toolbar-form :deep(.ant-form-item-control-input-content) {
  min-width: 0;
}

.toolbar-item {
  min-width: 0;
}

.toolbar-control {
  width: 100%;
}

.toolbar-submit {
  align-self: end;
}

.filter-button {
  min-width: 88px;
}

.toolbar-side {
  display: flex;
  justify-content: flex-end;
  align-items: end;
}

.view-mode-switch {
  flex: none;
}

.model-item {
  border: 1px solid var(--line);
  border-radius: 18px;
  margin-bottom: 12px;
  background: var(--surface-soft);
}

.item-summary {
  color: var(--text-secondary);
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 960px) {
  .toolbar {
    grid-template-columns: 1fr;
  }

  .toolbar-form {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .toolbar-submit {
    grid-column: 1 / -1;
  }

  .toolbar-side {
    justify-content: flex-start;
  }
}

@media (max-width: 720px) {
  .toolbar-form {
    grid-template-columns: 1fr;
  }

  .filter-button {
    width: 100%;
  }

  .view-mode-switch {
    width: 100%;
  }

  .view-mode-switch :deep(.ant-radio-button-wrapper) {
    width: 50%;
    text-align: center;
  }
}
</style>

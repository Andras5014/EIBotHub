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
            <a-form-item label="关键词" name="q">
              <a-input v-model:value="filters.q" placeholder="模型名 / 标签 / 描述" />
            </a-form-item>
            <a-form-item label="标签" name="tags">
              <a-input v-model:value="filters.tags" placeholder="多个标签逗号分隔" />
            </a-form-item>
            <a-form-item label="机器人类型" name="robot_type">
              <a-input v-model:value="filters.robot_type" placeholder="如：搬运" />
            </a-form-item>
            <a-form-item label="排序" name="sort">
              <a-select v-model:value="filters.sort" style="width: 140px">
                <a-select-option value="downloads">按热度</a-select-option>
                <a-select-option value="latest">最近更新</a-select-option>
                <a-select-option value="name">按名称</a-select-option>
              </a-select>
            </a-form-item>
            <a-form-item>
              <a-button html-type="submit" type="primary">筛选</a-button>
            </a-form-item>
          </a-form>

          <a-radio-group v-model:value="viewMode" button-style="solid">
            <a-radio-button value="card">卡片</a-radio-button>
            <a-radio-button value="list">列表</a-radio-button>
          </a-radio-group>
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
import { onMounted, reactive, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import ResourceCards from '@/components/ResourceCards.vue';
import type { ResourceCard } from '@/types/api';

const viewMode = ref<'card' | 'list'>('card');
const loading = ref(false);
const error = ref('');
const total = ref(0);
const items = ref<ResourceCard[]>([]);
const filters = reactive({
  q: '',
  tags: '',
  robot_type: '',
  sort: 'downloads',
  page: 1,
  page_size: 12,
});

async function load() {
  loading.value = true;
  error.value = '';
  try {
    const response = await api.listModels(filters);
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

async function handlePageChange(page: number) {
  filters.page = page;
  await load();
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString('zh-CN');
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.toolbar {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
}

.toolbar-form {
  flex: 1;
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
    flex-direction: column;
  }
}
</style>

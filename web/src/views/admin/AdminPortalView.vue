<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">门户运营</h1>
          <p class="section-subtitle">统一维护首页模块开关和推荐位，直接影响首页导读内容。</p>
        </div>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="9">
          <a-card title="首页模块开关" class="inner-card">
            <a-list :data-source="modules">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="item.label" :description="item.module_key" />
                  <a-switch :checked="item.enabled" @change="toggleModule(item.module_key, $event)" />
                </a-list-item>
              </template>
            </a-list>
          </a-card>

          <a-card title="新增推荐位" class="inner-card" style="margin-top: 16px">
            <a-form :model="createForm" layout="vertical" @finish="searchCandidates">
              <a-form-item label="资源类型">
                <a-select v-model:value="createForm.resource_type" :options="typeOptions" />
              </a-form-item>
              <a-form-item label="关键词">
                <a-input v-model:value="createForm.keyword" placeholder="搜索推荐资源" />
              </a-form-item>
              <a-form-item label="排序值">
                <a-input-number v-model:value="createForm.sort_order" :min="0" :max="1000" style="width: 100%" />
              </a-form-item>
              <a-form-item label="启用状态">
                <a-switch v-model:checked="createForm.enabled" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" html-type="submit" :loading="searching">搜索候选项</a-button>
              </a-form-item>
            </a-form>

            <a-list v-if="candidates.length" :data-source="candidates" size="small">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="item.title" :description="item.summary" />
                  <a-button type="link" @click="addFeatured(item)">加入</a-button>
                </a-list-item>
              </template>
            </a-list>
            <a-empty v-else description="先搜索候选项再加入推荐位" />
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="15">
          <a-card title="当前推荐位" class="inner-card">
            <a-empty v-if="!featured.length" description="还没有配置推荐位" />
            <div v-else class="featured-grid">
              <a-card v-for="item in featured" :key="item.id" class="featured-card">
                <div class="section-head featured-head">
                  <div>
                    <div class="list-title">{{ item.title }}</div>
                    <div class="list-desc">{{ item.summary }}</div>
                  </div>
                  <a-tag color="blue">{{ typeLabel(item.resource_type) }}</a-tag>
                </div>
                <a-space wrap style="margin-bottom: 12px">
                  <RouterLink v-if="item.route" :to="item.route">查看资源</RouterLink>
                  <span class="pill-meta">资源 ID {{ item.resource_id }}</span>
                </a-space>
                <a-row :gutter="[12, 12]">
                  <a-col :span="12">
                    <div class="field-label">排序值</div>
                    <a-input-number v-model:value="item.sort_order" :min="0" :max="1000" style="width: 100%" />
                  </a-col>
                  <a-col :span="12">
                    <div class="field-label">启用</div>
                    <a-switch v-model:checked="item.enabled" />
                  </a-col>
                </a-row>
                <a-space style="margin-top: 14px">
                  <a-button type="primary" @click="saveFeatured(item)">保存</a-button>
                  <a-button danger @click="removeFeatured(item.id)">删除</a-button>
                </a-space>
              </a-card>
            </div>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import type { ApplicationCaseItem, FeaturedResourceItem, ModuleSettingItem, SearchItem, TaskTemplateItem } from '@/types/api';

type CandidateItem = {
  id: number;
  title: string;
  summary: string;
  resource_type: string;
};

const modules = ref<ModuleSettingItem[]>([]);
const featured = ref<FeaturedResourceItem[]>([]);
const candidates = ref<CandidateItem[]>([]);
const searching = ref(false);
const createForm = reactive({
  resource_type: 'model',
  keyword: '',
  sort_order: 10,
  enabled: true,
});

const typeOptions = [
  { label: '模型', value: 'model' },
  { label: '数据集', value: 'dataset' },
  { label: '任务模板', value: 'task-template' },
  { label: '具身案例', value: 'application-case' },
];

async function load() {
  const [moduleItems, featuredItems] = await Promise.all([
    api.getAdminPortalModules(),
    api.getAdminFeaturedResources(),
  ]);
  modules.value = moduleItems;
  featured.value = featuredItems;
}

async function toggleModule(key: string, enabled: boolean) {
  await api.updateAdminPortalModule(key, { enabled });
  message.success('首页模块状态已更新');
  await load();
}

async function searchCandidates() {
  searching.value = true;
  try {
    candidates.value = await loadCandidates();
  } finally {
    searching.value = false;
  }
}

async function loadCandidates() {
  const keyword = createForm.keyword.trim();
  if (!keyword) return [];

  if (createForm.resource_type === 'task-template') {
    const items = await api.listTemplates();
    return filterTemplates(items, keyword);
  }
  if (createForm.resource_type === 'application-case') {
    const items = await api.listApplicationCases();
    return filterCases(items, keyword);
  }

  const response = await api.search({
    q: keyword,
    type: createForm.resource_type,
    sort: 'hot',
  });
  return response.items.map((item: SearchItem) => ({
    id: item.id,
    title: item.title,
    summary: item.summary,
    resource_type: mapSearchType(item.type),
  }));
}

function filterTemplates(items: TaskTemplateItem[], keyword: string) {
  return items
    .filter((item) => `${item.name} ${item.summary} ${item.category}`.includes(keyword))
    .map((item) => ({
      id: item.id,
      title: item.name,
      summary: item.summary,
      resource_type: 'task-template',
    }));
}

function filterCases(items: ApplicationCaseItem[], keyword: string) {
  return items
    .filter((item) => `${item.title} ${item.summary} ${item.category}`.includes(keyword))
    .map((item) => ({
      id: item.id,
      title: item.title,
      summary: item.summary,
      resource_type: 'application-case',
    }));
}

function mapSearchType(value: string) {
  if (value === 'task-template') return 'task-template';
  return value;
}

async function addFeatured(item: CandidateItem) {
  await api.createAdminFeaturedResource({
    resource_type: item.resource_type,
    resource_id: item.id,
    sort_order: createForm.sort_order,
    enabled: createForm.enabled,
  });
  message.success('推荐位已添加');
  candidates.value = [];
  createForm.keyword = '';
  await load();
}

async function saveFeatured(item: FeaturedResourceItem) {
  await api.updateAdminFeaturedResource(item.id, {
    resource_type: item.resource_type,
    resource_id: item.resource_id,
    sort_order: item.sort_order,
    enabled: item.enabled,
  });
  message.success('推荐位已保存');
  await load();
}

async function removeFeatured(id: number) {
  await api.deleteAdminFeaturedResource(id);
  message.success('推荐位已删除');
  await load();
}

function typeLabel(type: string) {
  return {
    model: '模型',
    dataset: '数据集',
    'task-template': '模板',
    'application-case': '案例',
  }[type] ?? type;
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.inner-card,
.featured-card {
  border-radius: 18px;
}

.featured-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.featured-head {
  margin-bottom: 12px;
}

.field-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

@media (max-width: 960px) {
  .featured-grid {
    grid-template-columns: 1fr;
  }
}
</style>

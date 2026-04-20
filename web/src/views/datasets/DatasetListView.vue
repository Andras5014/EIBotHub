<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">数据集仓库</h1>
          <p class="section-subtitle">支持样本预览、协议确认和下载记录追踪。</p>
        </div>
        <RouterLink to="/datasets/upload">
          <a-button type="primary">上传数据集</a-button>
        </RouterLink>
      </div>

      <a-form :model="filters" layout="inline" @finish="load">
        <a-form-item label="关键词" name="q">
          <a-input v-model:value="filters.q" placeholder="名称 / 标签 / 场景" />
        </a-form-item>
        <a-form-item label="标签" name="tags">
          <a-select v-model:value="tagValues" mode="multiple" allow-clear style="width: 220px" placeholder="选择标签">
            <a-select-option v-for="item in filterOptions.dataset_tags" :key="item" :value="item">{{ item }}</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="场景" name="scene">
          <a-select v-model:value="sceneValue" allow-clear style="width: 180px" placeholder="选择场景">
            <a-select-option v-for="item in filterOptions.dataset_scenes" :key="item" :value="item">{{ item }}</a-select-option>
          </a-select>
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

      <div style="margin-top: 20px">
        <ResourceCards :items="items" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import ResourceCards from '@/components/ResourceCards.vue';
import type { FilterOptionsResponse, ResourceCard } from '@/types/api';

const filters = reactive({
  q: '',
  scene: '',
  sort: 'downloads',
});
const items = ref<ResourceCard[]>([]);
const tagValues = ref<string[]>([]);
const sceneValue = ref('');
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

async function load() {
  const response = await api.listDatasets({
    ...filters,
    tags: tagValues.value.join(','),
    scene: sceneValue.value || undefined,
  });
  items.value = response.items;
}

onMounted(async () => {
  filterOptions.value = await api.getFilterOptions();
  await load();
});
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

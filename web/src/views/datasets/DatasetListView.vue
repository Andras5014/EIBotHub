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
          <a-input v-model:value="filters.tags" placeholder="如：inspection" />
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
import type { ResourceCard } from '@/types/api';

const filters = reactive({
  q: '',
  tags: '',
});
const items = ref<ResourceCard[]>([]);

async function load() {
  const response = await api.listDatasets(filters);
  items.value = response.items;
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

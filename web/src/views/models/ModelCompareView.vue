<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">模型对比</h1>
          <p class="section-subtitle">选择 2-3 个模型做并排对比，重点看摘要、规格、依赖和版本。</p>
        </div>
      </div>

      <a-space direction="vertical" style="width: 100%">
        <a-select
          v-model:value="selectedIds"
          mode="multiple"
          :max-tag-count="3"
          :options="options"
          style="width: 100%"
          placeholder="选择要对比的模型"
        />

        <a-row :gutter="[16, 16]">
          <a-col v-for="item in compared" :key="item.id" :xs="24" :lg="12" :xl="8">
            <a-card :title="item.name" class="compare-card">
              <p class="section-subtitle">{{ item.summary }}</p>
              <a-divider />
              <p><strong>适用机器人：</strong>{{ item.robot_type || '未填写' }}</p>
              <p><strong>输入规格：</strong>{{ item.input_spec || '未填写' }}</p>
              <p><strong>输出规格：</strong>{{ item.output_spec || '未填写' }}</p>
              <p><strong>许可证：</strong>{{ item.license || '未填写' }}</p>
              <p><strong>依赖：</strong>{{ item.dependencies.join(' / ') || '无' }}</p>
              <p><strong>版本数：</strong>{{ item.versions.length }}</p>
              <template #actions>
                <RouterLink :to="`/models/${item.id}`">查看详情</RouterLink>
              </template>
            </a-card>
          </a-col>
        </a-row>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import type { ModelDetail, ResourceCard } from '@/types/api';

const allModels = ref<ResourceCard[]>([]);
const compared = ref<ModelDetail[]>([]);
const selectedIds = ref<number[]>(readSelectedIds());

const options = computed(() =>
  allModels.value.map((item) => ({
    label: item.name,
    value: item.id,
  })),
);

watch(
  selectedIds,
  async (ids) => {
    localStorage.setItem('open-community-model-compare', JSON.stringify(ids.slice(0, 3)));
    compared.value = await Promise.all(ids.slice(0, 3).map((id) => api.getModel(id)));
  },
  { immediate: true },
);

onMounted(async () => {
  const response = await api.listModels({});
  allModels.value = response.items;
  if (selectedIds.value.length === 0 && response.items.length >= 2) {
    selectedIds.value = response.items.slice(0, 2).map((item) => item.id);
  }
});

function readSelectedIds() {
  const raw = localStorage.getItem('open-community-model-compare');
  if (!raw) return [] as number[];
  try {
    return JSON.parse(raw) as number[];
  } catch {
    return [];
  }
}
</script>

<style scoped>
.block {
  padding: 24px;
}

.compare-card {
  height: 100%;
}
</style>

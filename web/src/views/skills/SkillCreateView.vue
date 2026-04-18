<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">发布技能</h1>
          <p class="section-subtitle">把任务流或可复用逻辑沉淀成技能，供其他开发者评分和 Fork。</p>
        </div>
      </div>

      <a-form :model="form" layout="vertical" @finish="submit">
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="12"><a-form-item label="技能名称"><a-input v-model:value="form.name" /></a-form-item></a-col>
          <a-col :xs="24" :lg="12"><a-form-item label="摘要"><a-input v-model:value="form.summary" /></a-form-item></a-col>
        </a-row>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="12"><a-form-item label="分类"><a-input v-model:value="form.category" /></a-form-item></a-col>
          <a-col :xs="24" :lg="12"><a-form-item label="场景"><a-input v-model:value="form.scene" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="描述"><a-textarea v-model:value="form.description" :rows="4" /></a-form-item>
        <a-form-item label="指南"><a-textarea v-model:value="form.guide" :rows="4" /></a-form-item>
        <a-form-item label="关联资源"><a-input v-model:value="form.resource_ref" placeholder="逗号分隔" /></a-form-item>
        <a-form-item><a-button type="primary" html-type="submit">发布技能</a-button></a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive } from 'vue';
import { message } from 'ant-design-vue';
import { useRouter } from 'vue-router';

import { api } from '@/api';

const router = useRouter();
const form = reactive({
  name: '',
  summary: '',
  description: '',
  category: '',
  scene: '',
  guide: '',
  resource_ref: '',
});

async function submit() {
  try {
    const skill = await api.createSkill(form);
    message.success('技能已发布');
    await router.push(`/skills/${skill.id}`);
  } catch (error) {
    message.error(error instanceof Error ? error.message : '发布失败');
  }
}
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

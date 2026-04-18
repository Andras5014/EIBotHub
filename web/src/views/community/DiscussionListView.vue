<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">互动讨论</h1>
          <p class="section-subtitle">围绕接入、部署、训练和资源组合方式的社区讨论区。</p>
        </div>
      </div>

      <a-card title="发起讨论" style="margin-bottom: 16px">
        <a-form :model="form" layout="vertical" @finish="submit">
          <a-row :gutter="[16, 0]">
            <a-col :xs="24" :lg="12"><a-form-item label="标题"><a-input v-model:value="form.title" /></a-form-item></a-col>
            <a-col :xs="24" :lg="12"><a-form-item label="分类"><a-input v-model:value="form.category" /></a-form-item></a-col>
          </a-row>
          <a-form-item label="摘要"><a-input v-model:value="form.summary" /></a-form-item>
          <a-form-item label="内容"><a-textarea v-model:value="form.content" :rows="4" /></a-form-item>
          <a-form-item><a-button type="primary" html-type="submit">发布讨论</a-button></a-form-item>
        </a-form>
      </a-card>

      <a-list :data-source="items">
        <template #renderItem="{ item }">
          <a-list-item class="discussion-item">
            <RouterLink :to="`/discussions/${item.id}`" class="discussion-link">
              <a-list-item-meta :title="item.title" :description="item.summary" />
              <a-space wrap>
                <span class="pill-meta">{{ item.category }}</span>
                <span class="pill-meta">{{ item.user_name }}</span>
                <span class="pill-meta">评论 {{ item.comment_count }}</span>
              </a-space>
            </RouterLink>
          </a-list-item>
        </template>
      </a-list>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import type { DiscussionItem } from '@/types/api';

const items = ref<DiscussionItem[]>([]);
const form = reactive({
  title: '',
  summary: '',
  content: '',
  category: '接入讨论',
});

async function load() {
  items.value = await api.listDiscussions();
}

async function submit() {
  try {
    await api.createDiscussion(form);
    form.title = '';
    form.summary = '';
    form.content = '';
    await load();
    message.success('讨论已发布');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '发布失败');
  }
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.discussion-item {
  border: 1px solid var(--line);
  border-radius: 18px;
  margin-bottom: 12px;
  background: var(--surface-soft);
}

.discussion-link {
  width: 100%;
}
</style>

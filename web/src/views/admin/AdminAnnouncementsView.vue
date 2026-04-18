<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">公告管理</h1>
          <p class="section-subtitle">维护首页公告和运营导读位。</p>
        </div>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="10">
          <a-card title="发布公告">
            <a-form :model="form" layout="vertical" @finish="submit">
              <a-form-item label="标题" name="title"><a-input v-model:value="form.title" /></a-form-item>
              <a-form-item label="摘要" name="summary"><a-textarea v-model:value="form.summary" :rows="3" /></a-form-item>
              <a-form-item label="链接" name="link"><a-input v-model:value="form.link" /></a-form-item>
              <a-form-item name="pinned"><a-checkbox v-model:checked="form.pinned">置顶</a-checkbox></a-form-item>
              <a-form-item><a-button type="primary" html-type="submit">发布</a-button></a-form-item>
            </a-form>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="14">
          <a-card title="现有公告">
            <a-list :data-source="items">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="item.title" :description="item.summary" />
                  <a-tag v-if="item.pinned" color="blue">置顶</a-tag>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';

import { api } from '@/api';
import type { AnnouncementItem } from '@/types/api';

const form = reactive({
  title: '',
  summary: '',
  link: '',
  pinned: false,
});
const items = ref<AnnouncementItem[]>([]);

async function load() {
  items.value = await api.getAdminAnnouncements();
}

async function submit() {
  await api.createAnnouncement(form);
  form.title = '';
  form.summary = '';
  form.link = '';
  form.pinned = false;
  message.success('公告已发布');
  await load();
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

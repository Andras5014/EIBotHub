<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">社区 Wiki</h1>
          <p class="section-subtitle">由社区成员共同维护的技术知识库，支持修订历史与在线编辑。</p>
        </div>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="7">
          <a-card title="词条列表" class="inner-card">
            <a-space direction="vertical" style="width: 100%">
              <RouterLink v-if="canEdit" to="/wiki">
                <a-button block>新建词条</a-button>
              </RouterLink>
              <RouterLink v-else to="/login" :query="{ redirect: '/wiki' }">
                <a-button block>登录后新建词条</a-button>
              </RouterLink>
              <a-list :data-source="pages">
                <template #renderItem="{ item }">
                  <a-list-item>
                    <RouterLink :to="`/wiki/${item.id}`" @click="selected = item">
                      <div class="list-title">{{ item.title }}</div>
                      <div class="list-desc">{{ item.summary }}</div>
                    </RouterLink>
                  </a-list-item>
                </template>
              </a-list>
            </a-space>
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="10">
          <a-card title="词条内容" class="inner-card">
            <template v-if="selected">
              <h2>{{ selected.title }}</h2>
              <p class="section-subtitle">{{ selected.summary }}</p>
              <a-divider />
              <p style="white-space: pre-wrap">{{ selected.content }}</p>
              <a-space wrap style="margin-top: 12px">
                <span class="pill-meta">编辑者：{{ selected.editor_name }}</span>
                <span class="pill-meta">更新时间：{{ formatDate(selected.updated_at) }}</span>
              </a-space>
            </template>
            <a-empty v-else :description="canEdit ? '请先选择词条或新建词条' : '请先选择词条'" />
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="7">
          <a-card :title="canEdit ? '编辑 / 修订' : '修订历史'" class="inner-card">
            <template v-if="canEdit">
              <a-form :model="form" layout="vertical" @finish="submit">
                <a-form-item label="标题"><a-input v-model:value="form.title" /></a-form-item>
                <a-form-item label="摘要"><a-input v-model:value="form.summary" /></a-form-item>
                <a-form-item label="内容"><a-textarea v-model:value="form.content" :rows="8" /></a-form-item>
                <a-form-item label="修订说明"><a-input v-model:value="form.comment" /></a-form-item>
                <a-form-item>
                  <a-button type="primary" html-type="submit">{{ selected ? '更新词条' : '创建词条' }}</a-button>
                </a-form-item>
              </a-form>
              <a-divider />
            </template>
            <a-alert
              v-else
              type="info"
              show-icon
              style="margin-bottom: 16px"
              message="未登录时只能浏览 Wiki"
              description="登录后才可以新建词条、编辑内容和提交修订。"
            />
            <a-list :data-source="revisions" header="修订历史">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="item.editor_name" :description="`${item.comment || '无说明'} · ${formatDate(item.created_at)}`" />
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
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import { useAuthStore } from '@/stores/auth';
import type { WikiPageItem, WikiRevisionItem } from '@/types/api';

const auth = useAuthStore();
const route = useRoute();
const router = useRouter();
const pages = ref<WikiPageItem[]>([]);
const revisions = ref<WikiRevisionItem[]>([]);
const selected = ref<WikiPageItem>();
const canEdit = computed(() => auth.isAuthenticated);
const form = reactive({
  title: '',
  summary: '',
  content: '',
  comment: '',
});

async function loadPages() {
  pages.value = await api.listWikiPages();
  if (!canEdit.value && !route.params.id && pages.value.length) {
    await router.replace(`/wiki/${pages.value[0].id}`);
  }
}

async function loadSelected(id?: string) {
  if (!id) {
    selected.value = undefined;
    revisions.value = [];
    form.title = '';
    form.summary = '';
    form.content = '';
    form.comment = '';
    return;
  }
  const [page, revisionItems] = await Promise.all([
    api.getWikiPage(id),
    api.getWikiRevisions(id),
  ]);
  selected.value = page;
  revisions.value = revisionItems;
  form.title = page.title;
  form.summary = page.summary;
  form.content = page.content;
  form.comment = '';
}

async function submit() {
  if (!canEdit.value) {
    message.info('请先登录后再编辑 Wiki');
    return;
  }
  try {
    if (selected.value) {
      await api.updateWikiPage(selected.value.id, form);
      message.success('词条已更新');
      await loadSelected(String(selected.value.id));
    } else {
      const created = await api.createWikiPage(form);
      message.success('词条已创建');
      await router.push(`/wiki/${created.id}`);
    }
    await loadPages();
  } catch (error) {
    message.error(error instanceof Error ? error.message : 'Wiki 操作失败');
  }
}

function formatDate(value: string) {
  return new Date(value).toLocaleString();
}

watch(
  () => route.params.id,
  async (id) => {
    await loadSelected(id ? String(id) : undefined);
  },
  { immediate: true },
);

onMounted(async () => {
  await loadPages();
});
</script>

<style scoped>
.block {
  padding: 24px;
}

.inner-card {
  border-radius: 18px;
}
</style>

<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">文档中心</h1>
          <p class="section-subtitle">平台文档、技术文档与 FAQ 统一归集，补上目录搜索和面包屑导航。</p>
        </div>
        <RouterLink to="/wiki">
          <a-button>社区 Wiki</a-button>
        </RouterLink>
      </div>

      <a-tabs v-model:activeKey="docType" @change="load">
        <a-tab-pane key="platform" tab="平台文档" />
        <a-tab-pane key="technical" tab="技术文档" />
      </a-tabs>

      <a-breadcrumb style="margin-bottom: 16px">
        <a-breadcrumb-item>文档中心</a-breadcrumb-item>
        <a-breadcrumb-item>{{ docType === 'platform' ? '平台文档' : '技术文档' }}</a-breadcrumb-item>
        <a-breadcrumb-item v-if="selected">{{ selected.title }}</a-breadcrumb-item>
      </a-breadcrumb>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="7">
          <a-card title="分类" class="inner-card">
            <a-space wrap>
              <span v-for="item in categories" :key="item.id" class="pill-meta">{{ item.name }}</span>
            </a-space>
          </a-card>
          <a-card title="文档目录" class="inner-card" style="margin-top: 16px">
            <a-input v-model:value="keyword" placeholder="搜索文档标题或摘要" style="margin-bottom: 12px" />
            <a-list :data-source="filteredDocs">
              <template #renderItem="{ item }">
                <a-list-item>
                  <RouterLink :to="`/docs/${item.id}`" @click="selected = item">
                    <div class="list-title">{{ item.title }}</div>
                    <div class="list-desc">{{ item.summary }}</div>
                  </RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="11">
          <a-card title="文档内容" class="inner-card">
            <template v-if="selected">
              <h2>{{ selected.title }}</h2>
              <p class="section-subtitle">{{ selected.summary }}</p>
              <a-divider />
              <p style="white-space: pre-wrap">{{ selected.content }}</p>
            </template>
            <a-empty v-else description="请选择一篇文档" />
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="6">
          <a-card title="FAQ" class="inner-card">
            <a-collapse ghost>
              <a-collapse-panel v-for="item in faqs" :key="item.id" :header="item.question">
                <p>{{ item.answer }}</p>
              </a-collapse-panel>
            </a-collapse>
          </a-card>
          <a-card title="最新 Wiki" class="inner-card" style="margin-top: 16px">
            <a-list :data-source="wikiPages">
              <template #renderItem="{ item }">
                <a-list-item>
                  <RouterLink :to="`/wiki/${item.id}`">{{ item.title }}</RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
          <a-card title="视频教程" class="inner-card" style="margin-top: 16px">
            <a-list :data-source="videos">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a :href="item.link" target="_blank" rel="noreferrer">
                    <div class="list-title">{{ item.title }}</div>
                    <div class="list-desc">{{ item.summary }}</div>
                  </a>
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
import { computed, onMounted, ref, watch } from 'vue';
import { RouterLink, useRoute } from 'vue-router';

import { api } from '@/api';
import type { DocumentCategoryItem, DocumentItem, FAQItem, VideoTutorialItem, WikiPageItem } from '@/types/api';

const route = useRoute();
const docType = ref<'platform' | 'technical'>('platform');
const keyword = ref('');
const categories = ref<DocumentCategoryItem[]>([]);
const docs = ref<DocumentItem[]>([]);
const faqs = ref<FAQItem[]>([]);
const videos = ref<VideoTutorialItem[]>([]);
const wikiPages = ref<WikiPageItem[]>([]);
const selected = ref<DocumentItem>();

const filteredDocs = computed(() =>
  docs.value.filter((item) => {
    const query = keyword.value.trim().toLowerCase();
    if (!query) return true;
    return `${item.title} ${item.summary}`.toLowerCase().includes(query);
  }),
);

async function load() {
  const [categoryData, docData, faqData, videoData, wikiData] = await Promise.all([
    api.listDocCategories(docType.value),
    api.listDocs(docType.value),
    api.listFaqs(),
    api.listVideos(),
    api.listWikiPages(),
  ]);
  categories.value = categoryData;
  docs.value = docData;
  faqs.value = faqData;
  videos.value = videoData;
  wikiPages.value = wikiData.slice(0, 5);
  if (route.params.id) {
    selected.value = await api.getDoc(String(route.params.id));
  } else {
    selected.value = docs.value[0];
  }
}

watch(
  () => route.params.id,
  async (id) => {
    if (id) {
      selected.value = await api.getDoc(String(id));
      return;
    }
    selected.value = filteredDocs.value[0];
  },
);

watch(filteredDocs, (items) => {
  if (!items.length) {
    selected.value = undefined;
    return;
  }
  if (!selected.value || !items.some((item) => item.id === selected.value?.id)) {
    selected.value = items[0];
  }
});

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.inner-card {
  border-radius: 18px;
}
</style>

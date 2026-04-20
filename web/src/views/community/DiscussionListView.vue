<template>
  <div class="page-shell">
    <section class="page-card hero-panel">
      <div class="hero-head">
        <div>
          <a-tag color="blue">Community Threads</a-tag>
          <h1 class="section-title hero-title">互动讨论</h1>
          <p class="section-subtitle hero-subtitle">
            发帖只保留标题、标签和正文，用标签聚合讨论，再按热度、最新或评论活跃度切换同一条帖子流的展示顺序。
          </p>
        </div>
        <a-button type="primary" size="large" class="hero-create-button" @click="openCreateDialog">创建帖子</a-button>
      </div>
      <div class="hero-stats">
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="帖子总数" :value="items.length" />
        </a-card>
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="可筛选标签" :value="tagOptions.length" />
        </a-card>
        <a-card :bordered="false" class="hero-stat-card">
          <a-statistic title="当前命中" :value="filteredPosts.length" />
        </a-card>
      </div>
    </section>

    <a-modal
      v-model:open="composerOpen"
      title="创建帖子"
      :footer="null"
      width="860px"
      destroy-on-close
      @cancel="closeCreateDialog"
    >
      <p class="composer-tip modal-tip">只填标题、一个标签和正文。摘要会自动从正文里提炼出来。</p>
      <a-form :model="form" layout="vertical" @finish="submit">
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="16">
            <a-form-item label="标题" name="title">
              <a-input v-model:value="form.title" placeholder="一句话说清你要讨论的问题" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :lg="8">
            <a-form-item label="标签" name="tag">
              <a-input v-model:value="form.tag" placeholder="例如：部署联调" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="正文" name="content">
          <a-textarea
            v-model:value="form.content"
            :rows="6"
            placeholder="把你的问题、背景、已尝试过的方案和想获得的建议写清楚"
          />
        </a-form-item>
        <a-form-item style="margin-bottom: 0">
          <a-space wrap>
            <a-button type="primary" html-type="submit" :loading="submitting">发布帖子</a-button>
            <a-button @click="closeCreateDialog">取消</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </a-modal>

    <section class="page-card block filter-block">
      <div class="section-head">
        <div>
          <h2 class="section-title">帖子筛选</h2>
          <p class="section-subtitle">支持按关键词和标签筛选，再按你关心的排序方式浏览全部帖子。</p>
        </div>
        <span class="pill-meta">命中 {{ filteredPosts.length }} 帖</span>
      </div>

      <div class="filter-toolbar">
        <a-input-search
          v-model:value="keyword"
          allow-clear
          class="filter-search"
          placeholder="搜索标题、正文摘要或标签"
        />
        <a-select v-model:value="selectedTag" class="filter-select">
          <a-select-option value="">全部标签</a-select-option>
          <a-select-option v-for="item in tagOptions" :key="item" :value="item">{{ item }}</a-select-option>
        </a-select>
        <a-select v-model:value="listSort" class="filter-select">
          <a-select-option value="hot">按热度</a-select-option>
          <a-select-option value="latest">按最新</a-select-option>
          <a-select-option value="comments">按评论数</a-select-option>
        </a-select>
      </div>

      <a-alert v-if="error" type="warning" show-icon :message="error" style="margin-top: 16px" />
    </section>

    <section class="page-card block">
      <div class="section-head">
        <div>
          <h2 class="section-title">{{ currentListTitle }}</h2>
          <p class="section-subtitle">按当前筛选条件查看帖子流，排序方式可直接切换。</p>
        </div>
        <span class="pill-meta">排序 {{ sortLabelMap[listSort] }}</span>
      </div>

      <a-empty v-if="!sortedPosts.length" description="当前筛选条件下没有匹配帖子" />

      <a-list v-else :data-source="pagedPosts" class="post-list">
        <template #renderItem="{ item }">
          <a-list-item class="post-item">
            <RouterLink :to="`/discussions/${item.id}`" class="post-link">
              <div class="post-title-row">
                <div class="post-title">{{ item.title }}</div>
                <div class="post-date">{{ formatDate(item.updated_at) }}</div>
              </div>
              <div class="post-summary">{{ item.summary }}</div>
              <a-space wrap class="post-meta">
                <span class="pill-meta">{{ item.tag }}</span>
                <span class="pill-meta">{{ item.user_name }}</span>
                <span class="pill-meta">热度 {{ item.hot_score }}</span>
                <span class="pill-meta">评论 {{ item.comment_count }}</span>
                <span class="pill-meta">发布 {{ formatRelativeDate(item.created_at) }}</span>
              </a-space>
            </RouterLink>
          </a-list-item>
        </template>
      </a-list>

      <div v-if="sortedPosts.length > pageSize" class="pagination-wrap">
        <a-pagination
          v-model:current="currentPage"
          :page-size="pageSize"
          :total="sortedPosts.length"
          :show-size-changer="false"
        />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink, useRouter } from 'vue-router';

import { api } from '@/api';
import { useAuthStore } from '@/stores/auth';
import type { DiscussionItem } from '@/types/api';

type DiscussionSortKey = 'hot' | 'latest' | 'comments';

const auth = useAuthStore();
const router = useRouter();
const items = ref<DiscussionItem[]>([]);
const error = ref('');
const submitting = ref(false);
const composerOpen = ref(false);
const keyword = ref('');
const selectedTag = ref('');
const listSort = ref<DiscussionSortKey>('hot');
const currentPage = ref(1);
const pageSize = 8;
const sortLabelMap: Record<DiscussionSortKey, string> = {
  hot: '按热度',
  latest: '按最新',
  comments: '按评论数',
};
const form = reactive({
  title: '',
  tag: '',
  content: '',
});

const tagOptions = computed(() => {
  const counts = new Map<string, number>();
  for (const item of items.value) {
    const tag = item.tag.trim();
    if (!tag) continue;
    counts.set(tag, (counts.get(tag) ?? 0) + 1);
  }
  return [...counts.entries()]
    .sort((left, right) => right[1] - left[1] || left[0].localeCompare(right[0], 'zh-CN'))
    .map(([value]) => value);
});

const filteredPosts = computed(() => {
  const search = keyword.value.trim().toLowerCase();
  return items.value.filter((item) => {
    const matchesTag = !selectedTag.value || item.tag === selectedTag.value;
    if (!search) {
      return matchesTag;
    }
    const text = `${item.title} ${item.summary} ${item.content} ${item.tag}`.toLowerCase();
    return matchesTag && text.includes(search);
  });
});

const sortedPosts = computed(() => sortPosts(filteredPosts.value, listSort.value));
const currentListTitle = computed(() => {
  if (listSort.value === 'latest') return '最新发布';
  if (listSort.value === 'comments') return '评论最多';
  return '热门讨论';
});
const pagedPosts = computed(() => {
  const start = (currentPage.value - 1) * pageSize;
  return sortedPosts.value.slice(start, start + pageSize);
});

watch(
  () => [keyword.value, selectedTag.value, listSort.value],
  () => {
    currentPage.value = 1;
  },
);

async function load() {
  error.value = '';
  try {
    items.value = await api.listDiscussions();
  } catch (loadError) {
    items.value = [];
    error.value = loadError instanceof Error ? loadError.message : '讨论列表加载失败';
  }
}

async function submit() {
  if (!auth.isAuthenticated) {
    message.info('请先登录后再发布帖子');
    return;
  }
  submitting.value = true;
  try {
    await api.createDiscussion({
      title: form.title.trim(),
      tag: form.tag.trim(),
      content: form.content.trim(),
    });
    form.title = '';
    form.tag = '';
    form.content = '';
    composerOpen.value = false;
    await load();
    message.success('帖子已发布');
  } catch (submitError) {
    message.error(submitError instanceof Error ? submitError.message : '发布失败');
  } finally {
    submitting.value = false;
  }
}

function sortPosts(posts: DiscussionItem[], sortKey: DiscussionSortKey) {
  const next = [...posts];
  next.sort((left, right) => {
    if (sortKey === 'latest') {
      return new Date(right.created_at).getTime() - new Date(left.created_at).getTime();
    }
    if (sortKey === 'comments') {
      return right.comment_count - left.comment_count || right.hot_score - left.hot_score;
    }
    return right.hot_score - left.hot_score || right.comment_count - left.comment_count;
  });
  return next;
}

function formatDate(value: string) {
  return new Date(value).toLocaleDateString('zh-CN');
}

function formatRelativeDate(value: string) {
  const target = new Date(value).getTime();
  const diff = Date.now() - target;
  const hours = Math.floor(diff / (1000 * 60 * 60));
  if (hours < 1) return '刚刚';
  if (hours < 24) return `${hours} 小时前`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days} 天前`;
  return formatDate(value);
}

function openCreateDialog() {
  if (!auth.isAuthenticated) {
    message.info('请先登录后再创建帖子');
    void router.push({ path: '/login', query: { redirect: '/discussions' } });
    return;
  }
  composerOpen.value = true;
}

function closeCreateDialog() {
  composerOpen.value = false;
}

onMounted(load);
</script>

<style scoped>
.hero-panel {
  padding: 28px;
  display: grid;
  gap: 20px;
  grid-template-columns: 1fr;
  margin-bottom: 20px;
}

.hero-title {
  margin-top: 14px;
}

.hero-head {
  display: flex;
  justify-content: space-between;
  gap: 20px;
  align-items: flex-start;
}

.hero-subtitle {
  max-width: 760px;
  line-height: 1.75;
}

.hero-create-button {
  flex: none;
}

.hero-stats {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.hero-stat-card {
  border-radius: 20px;
  background: linear-gradient(180deg, #fbfdff, #eef7fb);
}

.block {
  padding: 24px;
}

.filter-block {
  margin-bottom: 20px;
}

.composer-tip {
  color: var(--text-secondary);
  font-size: 13px;
}

.modal-tip {
  margin: 0 0 16px;
}

.filter-toolbar {
  display: grid;
  grid-template-columns: minmax(260px, 1.6fr) repeat(2, minmax(180px, 0.8fr));
  gap: 12px;
}

.filter-search,
.filter-select {
  width: 100%;
}

.post-title {
  color: var(--text-main);
  font-size: 16px;
  font-weight: 700;
}

.post-summary {
  margin-top: 8px;
  color: var(--text-secondary);
  line-height: 1.7;
}

.post-meta {
  margin-top: 12px;
}

.post-list :deep(.ant-list-item) {
  padding: 0;
}

.post-item {
  border: 1px solid rgba(220, 230, 242, 0.82);
  border-radius: 18px;
  margin-bottom: 14px;
  background: linear-gradient(180deg, #fff 0%, #f9fbff 100%);
}

.post-link {
  display: block;
  width: 100%;
  padding: 18px;
}

.post-title-row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
}

.post-date {
  flex: none;
  color: var(--text-secondary);
  font-size: 12px;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

@media (max-width: 960px) {
  .hero-head {
    flex-direction: column;
  }

  .filter-toolbar {
    grid-template-columns: 1fr;
  }

  .hero-stats {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .post-title-row {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>

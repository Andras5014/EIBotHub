<template>
  <div class="page-shell">
    <div v-if="isEnabled('hero')" class="hero page-card">
      <div class="hero-copy">
        <a-tag color="blue">OpenLoong 风格信息门户</a-tag>
        <h1>围绕模型、数据集、模板与文档的开放社区</h1>
        <p>{{ home?.platform_intro }}</p>
        <a-space wrap>
          <RouterLink to="/models/upload">
            <a-button type="primary" size="large">上传模型</a-button>
          </RouterLink>
          <RouterLink to="/datasets/upload">
            <a-button size="large">上传数据集</a-button>
          </RouterLink>
          <RouterLink to="/search">
            <a-button type="link" size="large">进入全局搜索</a-button>
          </RouterLink>
        </a-space>
      </div>
      <div class="hero-side">
        <a-card title="平台亮点" :bordered="false" class="soft-card">
          <a-list :data-source="home?.highlights ?? []" size="small">
            <template #renderItem="{ item }">
              <a-list-item>{{ item }}</a-list-item>
            </template>
          </a-list>
        </a-card>
      </div>
    </div>

    <a-row v-if="isEnabled('models') || isEnabled('announcements')" :gutter="[16, 16]" style="margin-top: 20px">
      <a-col v-if="isEnabled('models')" :xs="24" :lg="isEnabled('announcements') ? 16 : 24">
        <div class="page-card section-block">
          <div class="section-head">
            <div>
              <h2 class="section-title">热门模型</h2>
              <p class="section-subtitle">优先展示下载热度高、描述完整、适配场景明确的模型资源。</p>
            </div>
            <RouterLink to="/models">查看更多</RouterLink>
          </div>
          <ResourceCards :items="home?.hot_models ?? []" />
        </div>
      </a-col>
      <a-col v-if="isEnabled('announcements')" :xs="24" :lg="isEnabled('models') ? 8 : 24">
        <div class="page-card section-block">
          <div class="section-head">
            <div>
              <h2 class="section-title">公告动态</h2>
              <p class="section-subtitle">首页采用论坛门户式分区布局，公告作为导读入口固定展示。</p>
            </div>
          </div>
          <a-timeline>
            <a-timeline-item v-for="item in home?.announcements ?? []" :key="item.id" :color="item.pinned ? 'blue' : 'gray'">
              <div class="announce-title">{{ item.title }}</div>
              <div class="announce-summary">{{ item.summary }}</div>
            </a-timeline-item>
          </a-timeline>
        </div>
      </a-col>
    </a-row>

    <div v-if="isEnabled('resources')" class="page-card section-block" style="margin-top: 20px">
      <div class="section-head">
        <div>
          <h2 class="section-title">数据集与模板</h2>
          <p class="section-subtitle">组合式展示数据、任务模板和具身应用案例，保持较高信息密度和清晰分区。</p>
        </div>
      </div>
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="10">
          <h3>热门数据集</h3>
          <ResourceCards :items="home?.hot_datasets ?? []" />
        </a-col>
        <a-col :xs="24" :lg="7">
          <h3>任务模板</h3>
          <a-list :data-source="home?.task_templates ?? []" bordered class="list-panel">
            <template #renderItem="{ item }">
              <a-list-item>
                <RouterLink :to="`/templates/${item.id}`">
                  <div class="list-title">{{ item.name }}</div>
                  <div class="list-desc">{{ item.summary }}</div>
                </RouterLink>
              </a-list-item>
            </template>
          </a-list>
        </a-col>
        <a-col :xs="24" :lg="7">
          <h3>具身案例</h3>
          <a-list :data-source="home?.application_cases ?? []" bordered class="list-panel">
            <template #renderItem="{ item }">
              <a-list-item>
                <RouterLink :to="`/applications/${item.id}`">
                  <div class="list-title">{{ item.title }}</div>
                  <div class="list-desc">{{ item.summary }}</div>
                </RouterLink>
              </a-list-item>
            </template>
          </a-list>
        </a-col>
      </a-row>
    </div>

    <div v-if="isEnabled('community')" class="page-card section-block" style="margin-top: 20px">
      <div class="section-head">
        <div>
          <h2 class="section-title">社区共创</h2>
          <p class="section-subtitle">把 Wiki 词条、协作讨论和贡献榜拉进首页，形成更强的社区导读入口。</p>
        </div>
      </div>
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="12">
          <h3>最新 Wiki</h3>
          <a-list :data-source="wikiPages" bordered class="list-panel">
            <template #renderItem="{ item }">
              <a-list-item>
                <RouterLink :to="`/wiki/${item.id}`">
                  <div class="list-title">{{ item.title }}</div>
                  <div class="list-desc">{{ item.summary }}</div>
                </RouterLink>
              </a-list-item>
            </template>
          </a-list>
        </a-col>
        <a-col :xs="24" :lg="12">
          <h3>贡献排行榜</h3>
          <a-list :data-source="rankings" bordered class="list-panel">
            <template #renderItem="{ item, index }">
              <a-list-item>
                <RouterLink :to="`/community/users/${item.user_id}`">
                  <div class="list-title">{{ index + 1 }}. {{ item.user_name }}</div>
                  <div class="list-desc">积分 {{ item.points }}</div>
                </RouterLink>
              </a-list-item>
            </template>
          </a-list>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import ResourceCards from '@/components/ResourceCards.vue';
import type { ContributorRankingItem, HomePayload, WikiPageItem } from '@/types/api';

const home = ref<HomePayload>();
const wikiPages = ref<WikiPageItem[]>([]);
const rankings = ref<ContributorRankingItem[]>([]);

onMounted(async () => {
  const [homeData, wikiData, rankingData] = await Promise.all([
    api.home(),
    api.listWikiPages(),
    api.contributorRankings(),
  ]);
  home.value = homeData;
  wikiPages.value = wikiData.slice(0, 5);
  rankings.value = rankingData.slice(0, 5);
});

function isEnabled(key: string) {
  const item = home.value?.module_settings?.find((setting) => setting.module_key === key);
  return item ? item.enabled : true;
}
</script>

<style scoped>
.hero {
  padding: 28px;
  display: grid;
  gap: 24px;
  grid-template-columns: 1.7fr 1fr;
}

.hero-copy h1 {
  margin: 12px 0 12px;
  font-size: 40px;
  line-height: 1.1;
}

.hero-copy p {
  font-size: 16px;
  color: var(--text-secondary);
  max-width: 720px;
}

.soft-card {
  border-radius: 22px;
  background: linear-gradient(180deg, #f8fbff, #eef4fb);
}

.section-block {
  padding: 24px;
}

.section-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  margin-bottom: 18px;
}

.announce-title,
.list-title {
  font-weight: 700;
}

.announce-summary,
.list-desc {
  margin-top: 4px;
  color: var(--text-secondary);
}

.list-panel {
  border-radius: 18px;
  overflow: hidden;
}

@media (max-width: 960px) {
  .hero {
    grid-template-columns: 1fr;
  }

  .hero-copy h1 {
    font-size: 32px;
  }
}
</style>

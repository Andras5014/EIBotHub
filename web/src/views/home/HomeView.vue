<template>
  <div class="page-shell">
    <template v-for="section in orderedSections" :key="section.key">
      <div v-if="section.key === 'hero'" class="hero page-card">
      <div class="hero-copy">
        <a-tag color="blue">{{ heroConfig.tagline }}</a-tag>
        <h1>{{ heroConfig.title }}</h1>
        <p>{{ heroConfig.description }}</p>
        <a-space wrap>
          <RouterLink to="/models/upload">
            <a-button type="primary" size="large">{{ heroConfig.primary_button }}</a-button>
          </RouterLink>
          <RouterLink to="/datasets/upload">
            <a-button size="large">{{ heroConfig.secondary_button }}</a-button>
          </RouterLink>
          <RouterLink to="/search">
            <a-button type="link" size="large">{{ heroConfig.search_button }}</a-button>
          </RouterLink>
        </a-space>
      </div>
      <div class="hero-side">
        <a-card title="平台亮点" :bordered="false" class="soft-card">
          <a-list :data-source="highlightItems" size="small">
            <template #renderItem="{ item }">
              <a-list-item>{{ item }}</a-list-item>
            </template>
          </a-list>
        </a-card>
      </div>
      </div>

      <a-row v-else-if="section.key === 'discovery'" :gutter="[16, 16]" style="margin-top: 20px">
      <a-col v-if="isEnabled('models')" :xs="24" :lg="isEnabled('announcements') ? 16 : 24">
        <div class="page-card section-block">
          <div class="section-head">
            <div>
              <h2 class="section-title">热门模型</h2>
              <p class="section-subtitle">优先展示下载热度高、描述完整、适配场景明确的模型资源。</p>
            </div>
            <RouterLink to="/models">查看更多</RouterLink>
          </div>
          <ResourceCards :items="hotModels" />
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
            <a-timeline-item v-for="item in announcements" :key="item.id" :color="item.pinned ? 'blue' : 'gray'">
              <div class="announce-title">{{ item.title }}</div>
              <div class="announce-summary">{{ item.summary }}</div>
            </a-timeline-item>
          </a-timeline>
        </div>
      </a-col>
      </a-row>

      <div v-else-if="section.key === 'scenes'" class="page-card section-block" style="margin-top: 20px">
      <div class="section-head">
        <div>
          <h2 class="section-title">应用场景</h2>
          <p class="section-subtitle">把场景当成资源导流入口，先看目标场景，再反查模型、数据集、模板和案例。</p>
        </div>
        <RouterLink to="/scenes">查看全部</RouterLink>
      </div>
      <div class="scene-grid">
        <RouterLink v-for="item in scenePages" :key="item.id" :to="`/scenes/${item.slug}`" class="scene-card">
          <div class="scene-kicker">{{ item.tagline || '场景专题' }}</div>
          <div class="scene-title">{{ item.name }}</div>
          <div class="scene-summary">{{ item.summary }}</div>
        </RouterLink>
      </div>
      </div>

      <div v-else-if="section.key === 'resources'" class="page-card section-block" style="margin-top: 20px">
      <div class="section-head">
        <div>
          <h2 class="section-title">数据集与模板</h2>
          <p class="section-subtitle">组合式展示数据、任务模板和具身应用案例，保持较高信息密度和清晰分区。</p>
        </div>
      </div>
      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="10">
          <h3>热门数据集</h3>
          <ResourceCards :items="hotDatasets" />
        </a-col>
        <a-col :xs="24" :lg="7">
          <h3>任务模板</h3>
          <a-list :data-source="taskTemplates" bordered class="list-panel">
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
          <a-list :data-source="applicationCases" bordered class="list-panel">
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

      <div v-else-if="section.key === 'community'" class="page-card section-block" style="margin-top: 20px">
      <div class="section-head">
        <div>
          <h2 class="section-title">社区共创</h2>
          <p class="section-subtitle">把 Wiki 词条、热门讨论和贡献榜拉进首页，形成更强的社区导读入口。</p>
        </div>
      </div>
      <a-row :gutter="[16, 16]" class="community-grid">
        <a-col :xs="24" :lg="isRankingEnabled ? 8 : 12">
          <div class="community-panel">
            <div class="community-panel-head">
              <div>
                <h3 class="community-panel-title">最新 Wiki</h3>
                <p class="community-panel-subtitle">持续收敛接入经验、模板说明和联调手册，保持社区知识沉淀入口稳定可读。</p>
              </div>
            </div>
            <a-list :data-source="visibleWikiPages" bordered class="list-panel community-list-panel">
              <template #renderItem="{ item }">
                <a-list-item>
                  <RouterLink :to="`/wiki/${item.id}`">
                    <div class="list-title">{{ item.title }}</div>
                    <div class="list-desc">{{ item.summary }}</div>
                  </RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </div>
        </a-col>
        <a-col :xs="24" :lg="isRankingEnabled ? 8 : 12">
          <div class="community-panel">
            <div class="community-panel-head">
              <div>
                <h3 class="community-panel-title">热门互动讨论</h3>
                <p class="community-panel-subtitle">优先展示当前社区里回复更活跃、热度更高的讨论帖子。</p>
              </div>
              <RouterLink to="/discussions">查看全部</RouterLink>
            </div>
            <a-list :data-source="visibleDiscussions" bordered class="list-panel community-list-panel">
              <template #renderItem="{ item }">
                <a-list-item>
                  <RouterLink :to="`/discussions/${item.id}`">
                    <div class="list-title">{{ item.title }}</div>
                    <div class="list-desc">{{ item.summary }}</div>
                    <div class="list-meta-line">
                      <span class="pill-meta">{{ item.tag }}</span>
                      <span class="pill-meta">评论 {{ item.comment_count }}</span>
                    </div>
                  </RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </div>
        </a-col>
        <a-col v-if="isRankingEnabled" :xs="24" :lg="8">
          <div class="community-panel">
            <div class="community-panel-head">
              <div>
                <h3 class="community-panel-title">{{ rankingTitle }}</h3>
                <p class="community-panel-subtitle">{{ rankingSubtitle }}</p>
              </div>
            </div>
            <a-list :data-source="visibleRankings" bordered class="list-panel community-list-panel">
              <template #renderItem="{ item, index }">
                <a-list-item>
                  <RouterLink :to="`/community/users/${item.user_id}`">
                    <div class="list-title">{{ index + 1 }}. {{ item.user_name }}</div>
                    <div class="list-desc">积分 {{ item.points }}</div>
                  </RouterLink>
                </a-list-item>
              </template>
            </a-list>
          </div>
        </a-col>
      </a-row>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import ResourceCards from '@/components/ResourceCards.vue';
import type { ApplicationCaseItem, AnnouncementItem, ContributorRankingItem, DiscussionItem, HomePayload, ResourceCard, ScenePageItem, TaskTemplateItem, WikiPageItem } from '@/types/api';

const home = ref<HomePayload>();
const wikiPages = ref<WikiPageItem[]>([]);
const rankings = ref<ContributorRankingItem[]>([]);
const discussions = ref<DiscussionItem[]>([]);
const HOME_SECTION_LIMITS = {
  highlights: 4,
  hotModels: 4,
  announcements: 4,
  scenePages: 6,
  hotDatasets: 2,
  taskTemplates: 4,
  applicationCases: 4,
  wikiPages: 4,
  discussions: 4,
  rankings: 4,
} as const;

const heroConfig = computed(() => home.value?.hero_config ?? {
  tagline: 'OpenLoong 风格信息门户',
  title: '围绕模型、数据集、模板与文档的开放社区',
  description: '开放社区聚合模型、数据集、任务模板、文档与具身应用案例，帮助开发者围绕机器人场景快速完成接入、训练、部署和复用。',
  primary_button: '上传模型',
  secondary_button: '上传数据集',
  search_button: '进入全局搜索',
});
const rankingTitle = computed(() => home.value?.rankings_config?.title || '贡献排行榜');
const rankingSubtitle = computed(() => home.value?.rankings_config?.subtitle || '基于积分展示近期社区贡献活跃度。');
const rankingLimit = computed(() => home.value?.rankings_config?.limit || 5);
const isRankingEnabled = computed(() => home.value?.rankings_config?.enabled ?? true);
const highlightItems = computed(() => limitItems(home.value?.highlights, HOME_SECTION_LIMITS.highlights));
const hotModels = computed(() => limitItems<ResourceCard>(home.value?.hot_models, HOME_SECTION_LIMITS.hotModels));
const announcements = computed(() => limitItems<AnnouncementItem>(home.value?.announcements, HOME_SECTION_LIMITS.announcements));
const scenePages = computed(() => limitItems<ScenePageItem>(home.value?.scene_pages, HOME_SECTION_LIMITS.scenePages));
const hotDatasets = computed(() => limitItems<ResourceCard>(home.value?.hot_datasets, HOME_SECTION_LIMITS.hotDatasets));
const taskTemplates = computed(() => limitItems<TaskTemplateItem>(home.value?.task_templates, HOME_SECTION_LIMITS.taskTemplates));
const applicationCases = computed(() => limitItems<ApplicationCaseItem>(home.value?.application_cases, HOME_SECTION_LIMITS.applicationCases));
const visibleWikiPages = computed(() => limitItems<WikiPageItem>(wikiPages.value, HOME_SECTION_LIMITS.wikiPages));
const visibleDiscussions = computed(() => limitItems<DiscussionItem>(discussions.value, HOME_SECTION_LIMITS.discussions));
const visibleRankings = computed(() =>
  limitItems<ContributorRankingItem>(rankings.value, Math.min(rankingLimit.value, HOME_SECTION_LIMITS.rankings)),
);
const orderedSections = computed(() => {
  if (!home.value) {
      return [
      { key: 'hero', sortOrder: 10 },
      { key: 'discovery', sortOrder: 20 },
      { key: 'scenes', sortOrder: 40 },
      { key: 'resources', sortOrder: 50 },
      { key: 'community', sortOrder: 60 },
    ];
  }
  const sortOrder = (key: string, fallback: number) =>
    home.value?.module_settings?.find((setting) => setting.module_key === key)?.sort_order ?? fallback;

  return [
    { key: 'hero', sortOrder: sortOrder('hero', 10), enabled: isEnabled('hero') },
    { key: 'discovery', sortOrder: Math.min(sortOrder('models', 20), sortOrder('announcements', 30)), enabled: isEnabled('models') || isEnabled('announcements') },
    { key: 'scenes', sortOrder: sortOrder('scenes', 40), enabled: isEnabled('scenes') && Boolean(home.value?.scene_pages?.length) },
    { key: 'resources', sortOrder: sortOrder('resources', 50), enabled: isEnabled('resources') },
    { key: 'community', sortOrder: sortOrder('community', 60), enabled: isEnabled('community') },
  ]
    .filter((item) => item.enabled)
    .sort((left, right) => left.sortOrder - right.sortOrder);
});

onMounted(async () => {
  const [homeData, wikiData, rankingData, discussionData] = await Promise.all([
    api.home(),
    api.listWikiPages(),
    api.contributorRankings(),
    api.listDiscussions({ sort: 'hot', limit: HOME_SECTION_LIMITS.discussions }),
  ]);
  home.value = homeData;
  wikiPages.value = wikiData;
  rankings.value = rankingData;
  discussions.value = discussionData;
});

function isEnabled(key: string) {
  const item = home.value?.module_settings?.find((setting) => setting.module_key === key);
  return item ? item.enabled : true;
}

function limitItems<T>(items: T[] | undefined, limit: number) {
  return (items ?? []).slice(0, limit);
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

.scene-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.scene-card {
  display: block;
  padding: 18px;
  border-radius: 18px;
  border: 1px solid rgba(220, 230, 242, 0.82);
  background: linear-gradient(180deg, #fff 0%, #f9fbff 100%);
}

.scene-kicker {
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.scene-title {
  margin-top: 10px;
  color: var(--text-main);
  font-size: 20px;
  font-weight: 700;
}

.scene-summary {
  margin-top: 10px;
  color: var(--text-secondary);
  line-height: 1.7;
}

.announce-summary,
.list-desc {
  margin-top: 4px;
  color: var(--text-secondary);
}

.list-meta-line {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 10px;
}

.community-grid {
  align-items: stretch;
}

.community-panel {
  display: flex;
  height: 100%;
  flex-direction: column;
}

.community-panel-head {
  display: flex;
  min-height: 86px;
  align-items: flex-start;
  margin-bottom: 12px;
}

.community-panel-title {
  margin: 0;
  color: var(--text-main);
  font-size: 20px;
  font-weight: 700;
}

.community-panel-subtitle {
  margin: 8px 0 0;
  color: var(--text-secondary);
  line-height: 1.7;
}

.list-panel {
  border-radius: 18px;
  overflow: hidden;
}

.community-list-panel {
  flex: 1;
}

@media (max-width: 960px) {
  .hero {
    grid-template-columns: 1fr;
  }

  .hero-copy h1 {
    font-size: 32px;
  }

  .scene-grid {
    grid-template-columns: 1fr;
  }

  .community-panel-head {
    min-height: 0;
  }
}
</style>

<template>
  <div class="page-shell">
    <a-spin :spinning="loading">
      <a-result
        v-if="error"
        status="warning"
        title="模板详情加载失败"
        :sub-title="error"
      >
        <template #extra>
          <a-space>
            <RouterLink to="/templates">
              <a-button>返回模板库</a-button>
            </RouterLink>
            <a-button type="primary" @click="load">重试</a-button>
          </a-space>
        </template>
      </a-result>

      <template v-else-if="detail">
        <section class="page-card hero-panel">
          <div>
            <a-space wrap>
              <span class="pill-meta">{{ detail.category }}</span>
              <span class="pill-meta">{{ detail.scene }}</span>
              <span class="pill-meta">使用 {{ detail.usage_count }}</span>
              <span class="pill-meta">更新于 {{ formatDate(detail.updated_at) }}</span>
            </a-space>
            <h1 class="section-title hero-title">{{ detail.name }}</h1>
            <p class="section-subtitle hero-subtitle">{{ detail.summary }}</p>
          </div>
          <a-space wrap>
            <RouterLink to="/templates">
              <a-button>返回列表</a-button>
            </RouterLink>
            <RouterLink :to="{ name: 'search', query: { q: detail.name, type: 'task-template' } }">
              <a-button>全局搜索</a-button>
            </RouterLink>
            <RouterLink :to="{ name: 'applications', query: { category: detail.category } }">
              <a-button type="primary">查看同类案例</a-button>
            </RouterLink>
          </a-space>
        </section>

        <a-alert
          v-if="!auth.isAuthenticated"
          type="info"
          show-icon
          class="login-alert"
          message="登录后可参与评分与评论，当前仍可浏览模板说明与关联案例。"
        />

        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :lg="15">
            <a-card title="模板说明" class="content-card">
              <p class="content-copy">{{ detail.description }}</p>
            </a-card>

            <a-card title="执行步骤" class="content-card" style="margin-top: 16px">
              <div v-if="guideSteps.length" class="step-list">
                <div v-for="(step, index) in guideSteps" :key="`${detail.id}-${index}`" class="step-item">
                  <span class="step-index">{{ index + 1 }}</span>
                  <div>
                    <div class="step-title">步骤 {{ index + 1 }}</div>
                    <div class="step-copy">{{ step }}</div>
                  </div>
                </div>
              </div>
              <a-empty v-else description="当前模板还没有填写详细步骤" />
            </a-card>

            <a-card title="模板评论" class="content-card" style="margin-top: 16px">
              <a-form :model="commentForm" layout="vertical" @finish="submitComment">
                <a-form-item label="评论内容">
                  <a-textarea
                    v-model:value="commentForm.content"
                    :rows="3"
                    :disabled="!auth.isAuthenticated"
                    placeholder="补充执行经验、常见问题或踩坑记录"
                  />
                </a-form-item>
                <a-form-item>
                  <a-button type="primary" html-type="submit" :disabled="!auth.isAuthenticated">发表评论</a-button>
                </a-form-item>
              </a-form>

              <a-list v-if="comments.length" :data-source="comments">
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta :description="item.content">
                      <template #title>
                        <div class="comment-title">
                          <span>{{ item.user_name }}</span>
                          <span class="comment-date">{{ formatDate(item.created_at) }}</span>
                        </div>
                      </template>
                    </a-list-item-meta>
                  </a-list-item>
                </template>
              </a-list>
              <a-empty v-else description="还没有评论，欢迎补充使用经验" />
            </a-card>
          </a-col>

          <a-col :xs="24" :lg="9">
            <a-card title="模板概览" class="content-card">
              <div class="overview-grid">
                <div class="overview-item">
                  <div class="overview-label">模板分类</div>
                  <div class="overview-value">{{ detail.category }}</div>
                </div>
                <div class="overview-item">
                  <div class="overview-label">应用场景</div>
                  <div class="overview-value">{{ detail.scene }}</div>
                </div>
                <div class="overview-item">
                  <div class="overview-label">累计使用</div>
                  <div class="overview-value">{{ detail.usage_count }}</div>
                </div>
                <div class="overview-item">
                  <div class="overview-label">最近更新</div>
                  <div class="overview-value">{{ formatDate(detail.updated_at) }}</div>
                </div>
              </div>
            </a-card>

            <a-card title="关联资源" class="content-card" style="margin-top: 16px">
              <a-space v-if="detail.resource_ref.length" wrap>
                <RouterLink
                  v-for="resource in detail.resource_ref"
                  :key="resource"
                  :to="{ name: 'search', query: { q: resource } }"
                >
                  <span class="pill-meta resource-chip">{{ resource }}</span>
                </RouterLink>
              </a-space>
              <a-empty v-else description="暂无关联资源" />
            </a-card>

            <a-card title="模板评分" class="content-card" style="margin-top: 16px">
              <a-statistic :value="ratings?.average ?? 0" :precision="1" title="平均分" />
              <p class="section-subtitle">共 {{ ratings?.count ?? 0 }} 条评分</p>
              <a-divider />

              <a-form :model="ratingForm" layout="vertical" @finish="submitRating">
                <a-form-item label="评分">
                  <a-rate v-model:value="ratingForm.score" :disabled="!auth.isAuthenticated" />
                </a-form-item>
                <a-form-item label="反馈">
                  <a-textarea
                    v-model:value="ratingForm.feedback"
                    :rows="3"
                    :disabled="!auth.isAuthenticated"
                    placeholder="说明适用设备、场地限制或执行建议"
                  />
                </a-form-item>
                <a-form-item>
                  <a-button type="primary" html-type="submit" :disabled="!auth.isAuthenticated">提交评分</a-button>
                </a-form-item>
              </a-form>

              <a-list v-if="ratings?.items.length" :data-source="ratings.items.slice(0, 3)" size="small">
                <template #renderItem="{ item }">
                  <a-list-item>
                    <a-list-item-meta :description="item.feedback || '该用户未填写反馈'">
                      <template #title>
                        <div class="comment-title">
                          <span>{{ item.user_name }} · {{ item.score }} 分</span>
                          <span class="comment-date">{{ formatDate(item.created_at) }}</span>
                        </div>
                      </template>
                    </a-list-item-meta>
                  </a-list-item>
                </template>
              </a-list>
            </a-card>
          </a-col>
        </a-row>

        <a-card
          v-if="relatedCases.length"
          title="相关应用案例"
          class="content-card"
          style="margin-top: 16px"
        >
          <div class="related-grid">
            <a-card v-for="item in relatedCases" :key="item.id" hoverable class="related-card">
              <template #title>{{ item.title }}</template>
              <p class="card-summary">{{ item.summary }}</p>
              <span class="pill-meta">{{ item.category }}</span>
              <template #actions>
                <RouterLink :to="`/applications/${item.id}`">查看案例</RouterLink>
              </template>
            </a-card>
          </div>
        </a-card>
      </template>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { message } from 'ant-design-vue';
import { RouterLink, useRoute } from 'vue-router';

import { api } from '@/api';
import { useAuthStore } from '@/stores/auth';
import type { ApplicationCaseItem, CommentItem, RatingSummary, TaskTemplateItem } from '@/types/api';

import { formatDate, splitGuideSteps } from './utils';

const route = useRoute();
const auth = useAuthStore();
const detail = ref<TaskTemplateItem>();
const ratings = ref<RatingSummary>();
const comments = ref<CommentItem[]>([]);
const relatedCases = ref<ApplicationCaseItem[]>([]);
const loading = ref(false);
const error = ref('');
const ratingForm = reactive({ score: 5, feedback: '' });
const commentForm = reactive({ content: '' });

const guideSteps = computed(() => splitGuideSteps(detail.value?.guide));

watch(
  () => route.params.id,
  () => {
    void load();
  },
  { immediate: true },
);

async function load() {
  loading.value = true;
  error.value = '';

  try {
    const id = String(route.params.id);
    const [template, ratingSummary, commentItems, cases] = await Promise.all([
      api.getTemplate(id),
      api.getTemplateRatings(id),
      api.getTemplateComments(id),
      api.listApplicationCases(),
    ]);

    detail.value = template;
    ratings.value = ratingSummary;
    comments.value = commentItems;
    relatedCases.value = cases.filter((item) => item.category === template.category).slice(0, 3);
  } catch (loadError) {
    error.value = loadError instanceof Error ? loadError.message : '模板详情获取失败';
  } finally {
    loading.value = false;
  }
}

async function submitRating() {
  if (!detail.value) return;
  if (!auth.isAuthenticated) {
    message.info('请先登录后再评分');
    return;
  }
  try {
    ratings.value = await api.rateTemplate(detail.value.id, ratingForm);
    ratingForm.feedback = '';
    message.success('评分已提交');
  } catch (submitError) {
    message.error(submitError instanceof Error ? submitError.message : '评分失败');
  }
}

async function submitComment() {
  if (!detail.value) return;
  if (!auth.isAuthenticated) {
    message.info('请先登录后再评论');
    return;
  }
  try {
    await api.commentTemplate(detail.value.id, commentForm);
    comments.value = await api.getTemplateComments(detail.value.id);
    commentForm.content = '';
    message.success('评论已提交');
  } catch (submitError) {
    message.error(submitError instanceof Error ? submitError.message : '评论失败');
  }
}
</script>

<style scoped>
.hero-panel {
  padding: 28px;
  display: flex;
  justify-content: space-between;
  gap: 20px;
  align-items: flex-start;
  margin-bottom: 16px;
}

.hero-title {
  margin-top: 14px;
}

.hero-subtitle {
  max-width: 760px;
}

.login-alert {
  margin-bottom: 16px;
}

.content-card {
  border-radius: 22px;
}

.content-copy {
  margin: 0;
  color: var(--text-secondary);
  white-space: pre-wrap;
}

.step-list {
  display: grid;
  gap: 12px;
}

.step-item {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  padding: 14px;
  border-radius: 18px;
  background: var(--surface-soft);
  border: 1px solid var(--line);
}

.step-index {
  width: 28px;
  height: 28px;
  display: inline-grid;
  place-items: center;
  border-radius: 999px;
  background: rgba(22, 119, 255, 0.12);
  color: var(--brand-strong);
  font-weight: 700;
  flex: none;
}

.step-title {
  font-weight: 700;
  margin-bottom: 4px;
}

.step-copy {
  color: var(--text-secondary);
}

.overview-grid {
  display: grid;
  gap: 12px;
}

.overview-item {
  padding: 14px;
  border-radius: 18px;
  background: var(--surface-soft);
  border: 1px solid var(--line);
}

.overview-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

.overview-value {
  font-size: 16px;
  font-weight: 700;
}

.comment-title {
  display: flex;
  justify-content: space-between;
  gap: 12px;
}

.comment-date {
  color: var(--text-secondary);
  font-size: 12px;
}

.resource-chip {
  cursor: pointer;
}

.related-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.related-card {
  border-radius: 18px;
}

.card-summary {
  color: var(--text-secondary);
  min-height: 44px;
}

@media (max-width: 960px) {
  .hero-panel {
    flex-direction: column;
  }

  .related-grid {
    grid-template-columns: 1fr;
  }
}
</style>

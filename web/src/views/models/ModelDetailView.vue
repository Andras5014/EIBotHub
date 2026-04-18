<template>
  <div class="page-shell" v-if="detail">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <a-space wrap>
            <a-tag color="blue">模型</a-tag>
            <span v-for="tag in detail.tags" :key="tag" class="pill-meta">{{ tag }}</span>
          </a-space>
          <h1 class="section-title" style="margin-top: 10px">{{ detail.name }}</h1>
          <p class="section-subtitle">{{ detail.summary }}</p>
        </div>
        <a-space>
          <a-button @click="toggleFavorite">{{ detail.favorited ? '取消收藏' : '收藏模型' }}</a-button>
          <a-button @click="addToCompare">加入对比</a-button>
          <a-button type="primary" @click="download">记录下载并打开文件</a-button>
        </a-space>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="16">
          <a-card title="说明" class="inner-card">
            <p>{{ detail.description }}</p>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="8">
          <a-card title="元信息" class="inner-card">
            <p><strong>适用机器人：</strong>{{ detail.robot_type || '未填写' }}</p>
            <p><strong>输入规格：</strong>{{ detail.input_spec || '未填写' }}</p>
            <p><strong>输出规格：</strong>{{ detail.output_spec || '未填写' }}</p>
            <p><strong>许可证：</strong>{{ detail.license || '未填写' }}</p>
            <p><strong>下载量：</strong>{{ detail.downloads }}</p>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="版本与依赖" class="inner-card" style="margin-top: 16px">
        <a-space wrap style="margin-bottom: 12px">
          <span v-for="dependency in detail.dependencies" :key="dependency" class="pill-meta">{{ dependency }}</span>
        </a-space>
        <a-table :data-source="detail.versions" :pagination="false" row-key="id" size="small">
          <a-table-column title="版本" data-index="version" />
          <a-table-column title="文件" data-index="file_name" />
          <a-table-column title="变更" data-index="changelog" />
          <a-table-column title="创建时间" data-index="created_at">
            <template #default="{ record }">{{ formatDate(record.created_at) }}</template>
          </a-table-column>
        </a-table>
      </a-card>

      <a-row :gutter="[16, 16]" style="margin-top: 16px">
        <a-col :xs="24" :lg="10">
          <a-card title="用户评分" class="inner-card">
            <a-statistic :value="ratings?.average ?? 0" :precision="1" title="平均分" />
            <p class="section-subtitle">共 {{ ratings?.count ?? 0 }} 条评分</p>
            <a-divider />
            <a-form :model="ratingForm" layout="vertical" @finish="submitRating">
              <a-form-item label="评分" name="score">
                <a-rate v-model:value="ratingForm.score" />
              </a-form-item>
              <a-form-item label="反馈" name="feedback">
                <a-textarea v-model:value="ratingForm.feedback" :rows="3" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" html-type="submit">提交评分</a-button>
              </a-form-item>
            </a-form>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="14">
          <a-card title="模型评测" class="inner-card">
            <a-form :model="evaluationForm" layout="vertical" @finish="submitEvaluation">
              <a-row :gutter="[12, 0]">
                <a-col :span="12"><a-form-item label="基准名称" name="benchmark"><a-input v-model:value="evaluationForm.benchmark" /></a-form-item></a-col>
                <a-col :span="12"><a-form-item label="评分" name="score"><a-input-number v-model:value="evaluationForm.score" :min="0" :max="100" style="width: 100%" /></a-form-item></a-col>
              </a-row>
              <a-form-item label="评测摘要" name="summary"><a-input v-model:value="evaluationForm.summary" /></a-form-item>
              <a-form-item label="评测说明" name="notes"><a-textarea v-model:value="evaluationForm.notes" :rows="3" /></a-form-item>
              <a-form-item><a-button type="primary" html-type="submit">提交评测</a-button></a-form-item>
            </a-form>

            <a-list :data-source="evaluations" style="margin-top: 12px">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="`${item.benchmark} · ${item.score}`" :description="`${item.summary} · ${item.user_name}`" />
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="模型评论" class="inner-card" style="margin-top: 16px">
        <a-form :model="commentForm" layout="vertical" @finish="submitComment">
          <a-form-item label="评论内容">
            <a-textarea v-model:value="commentForm.content" :rows="3" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" html-type="submit">发表评论</a-button>
          </a-form-item>
        </a-form>
        <a-list :data-source="comments">
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta :title="item.user_name" :description="item.content" />
            </a-list-item>
          </template>
        </a-list>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import type { CommentItem, ModelDetail, ModelEvaluationItem, RatingSummary } from '@/types/api';

const route = useRoute();
const router = useRouter();
const detail = ref<ModelDetail>();
const evaluations = ref<ModelEvaluationItem[]>([]);
const ratings = ref<RatingSummary>();
const comments = ref<CommentItem[]>([]);
const ratingForm = reactive({
  score: 5,
  feedback: '',
});
const commentForm = reactive({
  content: '',
});
const evaluationForm = reactive({
  benchmark: '自定义基准',
  summary: '',
  score: 85,
  notes: '',
});

async function load() {
  const id = String(route.params.id);
  const [model, evaluationItems, ratingSummary, commentItems] = await Promise.all([
    api.getModel(id),
    api.getModelEvaluations(id),
    api.getModelRatings(id),
    api.getModelComments(id),
  ]);
  detail.value = model;
  evaluations.value = evaluationItems;
  ratings.value = ratingSummary;
  comments.value = commentItems;
}

async function toggleFavorite() {
  if (!detail.value) return;
  await api.toggleFavorite({
    resource_type: 'model',
    resource_id: detail.value.id,
    title: detail.value.name,
  });
  await load();
  message.success('收藏状态已更新');
}

async function download() {
  if (!detail.value) return;
  await api.downloadModel(detail.value.id);
  const fileUrl = detail.value.versions[0]?.file_url;
  if (fileUrl) {
    window.open(fileUrl, '_blank');
  }
  message.success('已记录下载');
}

function addToCompare() {
  if (!detail.value) return;
  const key = 'open-community-model-compare';
  const current = readCompareIds();
  const next = [detail.value.id, ...current.filter((item) => item !== detail.value?.id)].slice(0, 3);
  localStorage.setItem(key, JSON.stringify(next));
  message.success('已加入模型对比');
  router.push('/models/compare');
}

async function submitRating() {
  if (!detail.value) return;
  try {
    ratings.value = await api.rateModel(detail.value.id, ratingForm);
    ratingForm.feedback = '';
    message.success('评分已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '评分失败');
  }
}

async function submitEvaluation() {
  if (!detail.value) return;
  try {
    await api.createModelEvaluation(detail.value.id, evaluationForm);
    evaluationForm.summary = '';
    evaluationForm.notes = '';
    await load();
    message.success('评测已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '评测失败');
  }
}

async function submitComment() {
  if (!detail.value) return;
  try {
    await api.commentModel(detail.value.id, commentForm);
    commentForm.content = '';
    comments.value = await api.getModelComments(detail.value.id);
    message.success('评论已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '评论失败');
  }
}

function formatDate(value: string) {
  return new Date(value).toLocaleString();
}

function readCompareIds() {
  const raw = localStorage.getItem('open-community-model-compare');
  if (!raw) return [] as number[];
  try {
    return JSON.parse(raw) as number[];
  } catch {
    return [];
  }
}

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

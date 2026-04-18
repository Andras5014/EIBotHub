<template>
  <div class="page-shell" v-if="detail">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <a-space wrap>
            <a-tag color="cyan">数据集</a-tag>
            <span v-for="tag in detail.tags" :key="tag" class="pill-meta">{{ tag }}</span>
          </a-space>
          <h1 class="section-title" style="margin-top: 10px">{{ detail.name }}</h1>
          <p class="section-subtitle">{{ detail.summary }}</p>
        </div>
        <a-space>
          <a-button @click="toggleFavorite">{{ detail.favorited ? '取消收藏' : '收藏数据集' }}</a-button>
          <a-button @click="confirmAgreement">确认协议</a-button>
          <a-button v-if="needsAccessRequest" @click="submitAccessRequest">申请下载权限</a-button>
          <a-button type="primary" @click="download">记录下载并打开文件</a-button>
        </a-space>
      </div>

      <a-alert
        v-if="needsAccessRequest"
        type="warning"
        show-icon
        style="margin-bottom: 16px"
        :message="`当前数据集为 ${detail.privacy} 级别，下载前需要通过访问审批。`"
        :description="accessRequestDescription"
      />

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="15">
          <a-card title="数据集说明" class="inner-card">
            <p>{{ detail.description }}</p>
            <a-divider />
            <strong>下载协议</strong>
            <p style="white-space: pre-wrap">{{ detail.agreement_text }}</p>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="9">
          <a-card title="元信息" class="inner-card">
            <p><strong>场景：</strong>{{ detail.scene }}</p>
            <p><strong>采集设备：</strong>{{ detail.device }}</p>
            <p><strong>样本量：</strong>{{ detail.sample_count }}</p>
            <p><strong>权限等级：</strong>{{ detail.privacy }}</p>
            <p><strong>下载量：</strong>{{ detail.downloads }}</p>
          </a-card>
        </a-col>
      </a-row>

      <a-card title="样本预览" class="inner-card" style="margin-top: 16px">
        <div class="sample-grid">
          <a-card v-for="item in samples" :key="item.id" class="sample-card" size="small">
            <template #title>
              <a-space wrap>
                <span>{{ item.title }}</span>
                <span class="pill-meta">{{ sampleTypeLabel(item.sample_type) }}</span>
              </a-space>
            </template>

            <template v-if="item.sample_type === 'text'">
              <p class="sample-copy">{{ item.preview_text }}</p>
            </template>

            <template v-else-if="item.sample_type === 'image' && item.preview_url">
              <img :src="item.preview_url" :alt="item.title" class="sample-image" />
              <p class="sample-copy">{{ item.preview_text }}</p>
            </template>

            <template v-else-if="item.sample_type === 'video'">
              <video v-if="item.preview_url" class="sample-video" :src="item.preview_url" controls preload="metadata"></video>
              <div class="sample-placeholder">
                <strong>视频样本预览</strong>
                <p>{{ item.preview_text }}</p>
                <a v-if="item.preview_url" :href="item.preview_url" target="_blank" rel="noreferrer">打开视频文件</a>
              </div>
            </template>

            <template v-else-if="item.sample_type === 'pointcloud'">
              <PointCloudPreview v-if="item.preview_url" :src="item.preview_url" />
              <div class="sample-placeholder">
                <strong>点云样本投影预览</strong>
                <p>{{ item.preview_text }}</p>
                <a v-if="item.preview_url" :href="item.preview_url" target="_blank" rel="noreferrer">下载点云文件</a>
              </div>
            </template>

            <template v-else>
              <div class="sample-placeholder">
                <strong>文件样本</strong>
                <p>{{ item.file_name || item.preview_text }}</p>
                <a v-if="item.preview_url" :href="item.preview_url" target="_blank" rel="noreferrer">打开文件</a>
              </div>
            </template>
          </a-card>
        </div>
      </a-card>

      <a-card title="版本信息" class="inner-card" style="margin-top: 16px">
        <a-table :data-source="detail.versions" :pagination="false" row-key="id" size="small">
          <a-table-column title="版本" data-index="version" />
          <a-table-column title="文件" data-index="file_name" />
          <a-table-column title="变更" data-index="changelog" />
          <a-table-column title="创建时间" data-index="created_at">
            <template #default="{ record }">{{ formatDate(record.created_at) }}</template>
          </a-table-column>
        </a-table>
      </a-card>

      <a-card v-if="auth.isAuthenticated" title="分批下载任务" class="inner-card" style="margin-top: 16px">
        <a-space style="margin-bottom: 16px">
          <a-input-number v-model:value="packageParts" :min="1" :max="10" />
          <a-button type="primary" :loading="creatingPackage" @click="createDownloadPackage">生成分批下载任务</a-button>
        </a-space>

        <a-list v-if="downloadPackages.length" :data-source="downloadPackages">
          <template #renderItem="{ item }">
            <a-list-item>
              <a-list-item-meta :description="`共 ${item.total_parts} 个分片 · ${formatDate(item.created_at)}`">
                <template #title>
                  <a-space wrap>
                    <span>任务 #{{ item.id }}</span>
                    <a-tag color="green">{{ item.status }}</a-tag>
                    <a v-if="item.bundle_url" :href="item.bundle_url" target="_blank" rel="noreferrer">下载总包说明</a>
                  </a-space>
                </template>
              </a-list-item-meta>
              <template #actions>
                <a v-for="(link, index) in item.part_links" :key="link" :href="link" target="_blank" rel="noreferrer">
                  分片 {{ index + 1 }}
                </a>
              </template>
            </a-list-item>
          </template>
        </a-list>
        <a-empty v-else description="还没有生成分批下载任务" />
      </a-card>

      <a-card v-if="stats" title="数据统计" class="inner-card" style="margin-top: 16px">
        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :lg="8"><a-statistic title="总下载量" :value="stats.download_count" /></a-col>
          <a-col :xs="24" :lg="8"><a-statistic title="样本量" :value="stats.sample_count" /></a-col>
          <a-col :xs="24" :lg="8"><a-statistic title="样本类型数" :value="stats.sample_type_mix.length" /></a-col>
        </a-row>
        <a-divider />
        <a-row :gutter="[16, 16]">
          <a-col :xs="24" :lg="12">
            <h3>样本类型分布</h3>
            <a-list :data-source="stats.sample_type_mix">
              <template #renderItem="{ item }">
                <a-list-item>{{ item.label }} · {{ item.value }}</a-list-item>
              </template>
            </a-list>
          </a-col>
          <a-col :xs="24" :lg="12">
            <h3>最近下载趋势</h3>
            <a-list :data-source="stats.download_trend">
              <template #renderItem="{ item }">
                <a-list-item>{{ item.label }} · {{ item.value }}</a-list-item>
              </template>
            </a-list>
          </a-col>
        </a-row>
      </a-card>

      <a-row :gutter="[16, 16]" style="margin-top: 16px">
        <a-col :xs="24" :lg="10">
          <a-card title="数据集评分" class="inner-card">
            <a-statistic :value="ratings?.average ?? 0" :precision="1" title="平均分" />
            <p class="section-subtitle">共 {{ ratings?.count ?? 0 }} 条评分</p>
            <a-divider />
            <a-form :model="ratingForm" layout="vertical" @finish="submitRating">
              <a-form-item label="评分"><a-rate v-model:value="ratingForm.score" /></a-form-item>
              <a-form-item label="反馈"><a-textarea v-model:value="ratingForm.feedback" :rows="3" /></a-form-item>
              <a-form-item><a-button type="primary" html-type="submit">提交评分</a-button></a-form-item>
            </a-form>
          </a-card>
        </a-col>
        <a-col :xs="24" :lg="14">
          <a-card title="数据集评论" class="inner-card">
            <a-form :model="commentForm" layout="vertical" @finish="submitComment">
              <a-form-item label="评论内容"><a-textarea v-model:value="commentForm.content" :rows="3" /></a-form-item>
              <a-form-item><a-button type="primary" html-type="submit">发表评论</a-button></a-form-item>
            </a-form>
            <a-list :data-source="comments">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :title="item.user_name" :description="item.content" />
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
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { useRoute } from 'vue-router';

import { api } from '@/api';
import PointCloudPreview from '@/components/PointCloudPreview.vue';
import { useAuthStore } from '@/stores/auth';
import type {
  CommentItem,
  DatasetAccessRequestItem,
  DatasetDetail,
  DatasetSample,
  DatasetStatsResponse,
  DownloadPackageTaskItem,
  RatingSummary,
} from '@/types/api';

const route = useRoute();
const auth = useAuthStore();
const detail = ref<DatasetDetail>();
const stats = ref<DatasetStatsResponse>();
const ratings = ref<RatingSummary>();
const comments = ref<CommentItem[]>([]);
const samples = ref<DatasetSample[]>([]);
const downloadPackages = ref<DownloadPackageTaskItem[]>([]);
const accessRequest = ref<DatasetAccessRequestItem | null>(null);
const creatingPackage = ref(false);
const packageParts = ref(3);
const ratingForm = reactive({ score: 5, feedback: '' });
const commentForm = reactive({ content: '' });

async function load() {
  const id = String(route.params.id);
  const [dataset, sampleItems, statsData, ratingSummary, commentItems] = await Promise.all([
    api.getDataset(id),
    api.getDatasetSamples(id),
    api.getDatasetStats(id),
    api.getDatasetRatings(id),
    api.getDatasetComments(id),
  ]);
  detail.value = dataset;
  samples.value = sampleItems;
  stats.value = statsData;
  ratings.value = ratingSummary;
  comments.value = commentItems;
  if (auth.isAuthenticated) {
    await loadAccessRequest();
    await loadDownloadPackages();
  }
}

async function loadAccessRequest() {
  if (!detail.value || !auth.isAuthenticated) return;
  try {
    accessRequest.value = await api.getMyDatasetAccessRequest(detail.value.id);
  } catch {
    accessRequest.value = null;
  }
}

async function loadDownloadPackages() {
  if (!detail.value) return;
  try {
    downloadPackages.value = await api.getDatasetDownloadPackages(detail.value.id);
  } catch {
    downloadPackages.value = [];
  }
}

async function toggleFavorite() {
  if (!detail.value) return;
  await api.toggleFavorite({
    resource_type: 'dataset',
    resource_id: detail.value.id,
    title: detail.value.name,
  });
  await load();
  message.success('收藏状态已更新');
}

async function confirmAgreement() {
  if (!detail.value) return;
  await api.confirmDatasetAgreement(detail.value.id);
  message.success('已确认协议');
}

async function download() {
  if (!detail.value) return;
  await api.downloadDataset(detail.value.id);
  const fileUrl = detail.value.versions[0]?.file_url;
  if (fileUrl) {
    window.open(fileUrl, '_blank');
  }
  message.success('已记录下载');
}

async function createDownloadPackage() {
  if (!detail.value) return;
  creatingPackage.value = true;
  try {
    await api.createDatasetDownloadPackage(detail.value.id, { parts: packageParts.value ?? 3 });
    message.success('分批下载任务已生成');
    await loadDownloadPackages();
  } catch (error) {
    message.error(error instanceof Error ? error.message : '生成下载任务失败');
  } finally {
    creatingPackage.value = false;
  }
}

async function submitAccessRequest() {
  if (!detail.value) return;
  try {
    accessRequest.value = await api.createDatasetAccessRequest(detail.value.id, {
      reason: `申请访问数据集《${detail.value.name}》，用于场景验证与模型联调。`,
    });
    message.success('访问申请已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '访问申请提交失败');
  }
}

async function submitRating() {
  if (!detail.value) return;
  try {
    ratings.value = await api.rateDataset(detail.value.id, ratingForm);
    ratingForm.feedback = '';
    message.success('评分已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '评分失败');
  }
}

async function submitComment() {
  if (!detail.value) return;
  try {
    await api.commentDataset(detail.value.id, commentForm);
    commentForm.content = '';
    comments.value = await api.getDatasetComments(detail.value.id);
    message.success('评论已提交');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '评论失败');
  }
}

function sampleTypeLabel(type: string) {
  return {
    text: '文本',
    image: '图片',
    video: '视频',
    pointcloud: '点云',
    file: '文件',
  }[type] ?? type;
}

const needsAccessRequest = computed(() => ['internal', 'restricted'].includes(detail.value?.privacy ?? ''));

const accessRequestDescription = computed(() => {
  if (!accessRequest.value) {
    return '你还没有提交访问申请。';
  }
  const statusText = {
    pending: '申请已提交，等待管理员审核。',
    approved: '访问申请已通过，现在可以下载数据集。',
    rejected: '访问申请已被驳回，可修改说明后重新申请。',
  }[accessRequest.value.status] ?? accessRequest.value.status;
  return accessRequest.value.review_comment ? `${statusText} 审核意见：${accessRequest.value.review_comment}` : statusText;
});

function formatDate(value: string) {
  return new Date(value).toLocaleString('zh-CN');
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

.sample-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.sample-card {
  border-radius: 16px;
}

.sample-copy {
  color: var(--text-secondary);
  white-space: pre-wrap;
}

.sample-image {
  width: 100%;
  border-radius: 12px;
  border: 1px solid var(--line);
  margin-bottom: 12px;
}

.sample-video {
  width: 100%;
  border-radius: 12px;
  border: 1px solid var(--line);
  margin-bottom: 12px;
  background: #000;
}

.sample-placeholder {
  min-height: 140px;
  display: grid;
  align-content: start;
  gap: 8px;
  padding: 12px;
  border-radius: 12px;
  background: var(--surface-soft);
  border: 1px dashed var(--line);
}
</style>

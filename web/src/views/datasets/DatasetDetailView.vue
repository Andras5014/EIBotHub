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
          <a-button v-if="canSubmitAccessRequest" @click="submitAccessRequest">申请下载权限</a-button>
          <a-button type="primary" @click="download">记录下载并打开文件</a-button>
        </a-space>
      </div>

      <a-alert
        v-if="needsAccessRequest"
        type="warning"
        show-icon
        style="margin-bottom: 16px"
        :message="accessAlertMessage"
        :description="accessRequestDescription"
      />

      <a-alert
        v-if="detail?.privacy === 'internal' && hasApprovedVerification"
        type="success"
        show-icon
        style="margin-bottom: 16px"
        message="当前数据集为 internal 级别，你已通过开发者认证，确认协议后可直接下载。"
      />

      <a-card v-if="needsAccessRequest && auth.isAuthenticated" title="审批轨迹" class="inner-card" style="margin-bottom: 16px">
        <a-empty v-if="!accessHistory.length" description="还没有访问申请记录" />
        <a-timeline v-else class="access-timeline">
          <a-timeline-item v-for="item in accessHistory" :key="item.id" :color="statusColor(item.status)">
            <div class="timeline-title">{{ statusLabel(item.status) }} · {{ formatDate(item.created_at) }}</div>
            <div class="timeline-desc">申请说明：{{ item.reason }}</div>
            <div v-if="item.required_approvals > 1" class="timeline-desc">审批进度：{{ item.approval_stage }}/{{ item.required_approvals }}</div>
            <div v-if="item.sla_hours > 0" class="timeline-desc">处理时限：{{ slaSummary(item) }}</div>
            <div v-if="item.reviewed_at" class="timeline-desc">处理时间：{{ formatDate(item.reviewed_at) }}</div>
            <div v-if="item.status === 'approved'" class="timeline-desc">授权信息：{{ accessGrantSummary(item) }}</div>
            <div v-if="item.review_comment" class="timeline-desc">审核意见：{{ item.review_comment }}</div>
          </a-timeline-item>
        </a-timeline>
      </a-card>

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
            <p v-if="accessRequest?.status === 'approved'"><strong>授权状态：</strong>{{ accessRequest.authorization_active ? '有效' : '已失效' }}</p>
            <p v-if="accessRequest?.status === 'approved' && accessRequest.approval_expires_at"><strong>授权到期：</strong>{{ formatDate(accessRequest.approval_expires_at) }}</p>
            <p v-if="accessRequest?.status === 'approved'"><strong>剩余下载：</strong>{{ remainingDownloadText(accessRequest) }}</p>
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
                <strong>点云样本三维预览</strong>
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
        <a-alert
          v-if="detail.review_comment && detail.status === 'rejected'"
          type="warning"
          show-icon
          style="margin-bottom: 16px"
          :message="`驳回原因：${detail.review_comment}`"
        />
        <a-table :data-source="detail.versions" :pagination="false" row-key="id" size="small">
          <a-table-column title="版本" data-index="version" />
          <a-table-column title="文件" data-index="file_name" />
          <a-table-column title="变更" data-index="changelog" />
          <a-table-column title="创建时间" data-index="created_at">
            <template #default="{ record }">{{ formatDate(record.created_at) }}</template>
          </a-table-column>
        </a-table>

        <div v-if="detail.versions.length > 1" class="compare-section">
          <div class="compare-head">
            <div>
              <h3>历史版本对比</h3>
              <p class="section-subtitle">选择两个版本，快速比较数据包文件与版本说明。</p>
            </div>
          </div>
          <a-row :gutter="[12, 12]" class="compare-toolbar">
            <a-col :xs="24" :md="12">
              <a-select v-model:value="leftVersionId" style="width: 100%" placeholder="选择左侧版本">
                <a-select-option v-for="item in detail.versions" :key="item.id" :value="item.id">{{ item.version }}</a-select-option>
              </a-select>
            </a-col>
            <a-col :xs="24" :md="12">
              <a-select v-model:value="rightVersionId" style="width: 100%" placeholder="选择右侧版本">
                <a-select-option v-for="item in detail.versions" :key="item.id" :value="item.id">{{ item.version }}</a-select-option>
              </a-select>
            </a-col>
          </a-row>

          <a-row v-if="leftVersion && rightVersion" :gutter="[16, 16]" class="compare-grid">
            <a-col :xs="24" :lg="12">
              <a-card class="compare-card" :title="leftVersion.version">
                <p><strong>文件：</strong>{{ leftVersion.file_name }}</p>
                <p><strong>时间：</strong>{{ formatDate(leftVersion.created_at) }}</p>
                <p><strong>说明：</strong>{{ leftVersion.changelog || '无' }}</p>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="12">
              <a-card class="compare-card" :title="rightVersion.version">
                <p><strong>文件：</strong>{{ rightVersion.file_name }}</p>
                <p><strong>时间：</strong>{{ formatDate(rightVersion.created_at) }}</p>
                <p><strong>说明：</strong>{{ rightVersion.changelog || '无' }}</p>
              </a-card>
            </a-col>
            <a-col :span="24">
              <a-card class="compare-summary-card" title="差异摘要">
                <a-space wrap>
                  <a-tag :color="leftVersion.file_name === rightVersion.file_name ? 'default' : 'blue'">
                    {{ leftVersion.file_name === rightVersion.file_name ? '文件未变化' : '文件已变化' }}
                  </a-tag>
                  <a-tag :color="leftVersion.changelog === rightVersion.changelog ? 'default' : 'purple'">
                    {{ leftVersion.changelog === rightVersion.changelog ? '版本说明一致' : '版本说明已变化' }}
                  </a-tag>
                  <a-tag color="gold">
                    时间跨度 {{ compareDuration }}
                  </a-tag>
                </a-space>
              </a-card>
            </a-col>
          </a-row>
        </div>
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
  VerificationStatusItem,
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
const accessHistory = ref<DatasetAccessRequestItem[]>([]);
const verification = ref<VerificationStatusItem | null>(null);
const leftVersionId = ref<number>();
const rightVersionId = ref<number>();
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
  leftVersionId.value = dataset.versions[1]?.id ?? dataset.versions[0]?.id;
  rightVersionId.value = dataset.versions[0]?.id;
  if (auth.isAuthenticated) {
    verification.value = await api.myVerification();
    await loadAccessRequest();
    await loadAccessHistory();
    await loadDownloadPackages();
  } else {
    accessHistory.value = [];
    verification.value = null;
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

async function loadAccessHistory() {
  if (!detail.value || !auth.isAuthenticated) return;
  try {
    accessHistory.value = await api.getMyDatasetAccessHistory(detail.value.id);
  } catch {
    accessHistory.value = [];
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
    await loadAccessHistory();
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

const hasApprovedVerification = computed(() => verification.value?.status === 'approved');
const needsAccessRequest = computed(() => {
  if (detail.value?.privacy === 'restricted') {
    return true;
  }
  if (detail.value?.privacy === 'internal') {
    return !hasApprovedVerification.value;
  }
  return false;
});
const canSubmitAccessRequest = computed(() => {
  if (!needsAccessRequest.value || !auth.isAuthenticated) {
    return false;
  }
  return !accessRequest.value || accessRequest.value.status === 'rejected';
});
const accessAlertMessage = computed(() => {
  if (detail.value?.privacy === 'internal') {
    return '当前数据集为 internal 级别，未通过开发者认证的用户需要通过访问审批。';
  }
  return `当前数据集为 ${detail.value?.privacy} 级别，下载前需要通过访问审批。`;
});

const accessRequestDescription = computed(() => {
  if (!accessRequest.value) {
    return needsAccessRequest.value ? '你还没有提交访问申请。' : '当前权限条件下无需额外审批。';
  }
  const statusText = {
    pending:
      accessRequest.value.required_approvals > 1 && accessRequest.value.approval_stage > 0
        ? `当前申请已完成第 ${accessRequest.value.approval_stage}/${accessRequest.value.required_approvals} 级审批，仍在等待后续审核。${accessRequest.value.sla_hours > 0 ? slaSummary(accessRequest.value) : ''}`
        : accessRequest.value.sla_hours > 0
          ? `申请已提交，等待管理员审核。${slaSummary(accessRequest.value)}`
          : '申请已提交，等待管理员审核。',
    approved: accessRequest.value.authorization_active ? '访问申请已通过，现在可以下载数据集。' : '访问申请已通过，但当前授权已失效。',
    rejected: '访问申请已被驳回，可修改说明后重新申请。',
  }[accessRequest.value.status] ?? accessRequest.value.status;
  const grantText = accessRequest.value.status === 'approved' ? ` 授权信息：${accessGrantSummary(accessRequest.value)}` : '';
  return accessRequest.value.review_comment ? `${statusText}${grantText} 审核意见：${accessRequest.value.review_comment}` : `${statusText}${grantText}`;
});
const leftVersion = computed(() => detail.value?.versions.find((item) => item.id === leftVersionId.value));
const rightVersion = computed(() => detail.value?.versions.find((item) => item.id === rightVersionId.value));
const compareDuration = computed(() => {
  if (!leftVersion.value || !rightVersion.value) {
    return '0 天';
  }
  const diff = Math.abs(new Date(rightVersion.value.created_at).getTime() - new Date(leftVersion.value.created_at).getTime());
  const days = Math.max(1, Math.round(diff / (24 * 60 * 60 * 1000)));
  return `${days} 天`;
});

function statusLabel(status?: string) {
  return {
    pending: '待审核',
    approved: '已通过',
    rejected: '已驳回',
  }[status ?? ''] ?? status;
}

function statusColor(status?: string) {
  return {
    pending: 'gold',
    approved: 'green',
    rejected: 'red',
  }[status ?? ''] ?? 'gray';
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('zh-CN');
}

function remainingDownloadText(item: DatasetAccessRequestItem) {
  return item.download_limit > 0 ? `${item.download_count}/${item.download_limit}，剩余 ${item.remaining_downloads ?? 0} 次` : '不限次数';
}

function accessGrantSummary(item: DatasetAccessRequestItem) {
  const limitText = item.download_limit > 0 ? `剩余 ${item.remaining_downloads ?? 0} 次` : '不限次数';
  const expireText = item.approval_expires_at ? formatDate(item.approval_expires_at) : '长期有效';
  const statusText = item.authorization_active ? '授权有效' : item.is_expired ? '授权已过期' : '次数已用尽';
  return `${statusText} · ${limitText} · 有效期至 ${expireText}`;
}

function slaSummary(item: DatasetAccessRequestItem) {
  if (!item.sla_hours || !item.sla_deadline_at) {
    return '无 SLA 要求';
  }
  return item.sla_overdue
    ? `SLA ${item.sla_hours} 小时，已超时 ${formatMinutes(-item.sla_remaining_minutes)}`
    : `SLA ${item.sla_hours} 小时，剩余 ${formatMinutes(item.sla_remaining_minutes)}`;
}

function formatMinutes(value: number) {
  const minutes = Math.max(value, 0);
  const hours = Math.floor(minutes / 60);
  const rest = minutes % 60;
  if (hours > 0) {
    return `${hours} 小时 ${rest} 分`;
  }
  return `${rest} 分`;
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

.access-timeline {
  margin-top: 4px;
}

.timeline-title {
  font-weight: 600;
  color: var(--text-main);
}

.timeline-desc {
  margin-top: 6px;
  color: var(--text-secondary);
}

.compare-section {
  margin-top: 20px;
}

.compare-toolbar {
  margin-top: 12px;
}

.compare-grid {
  margin-top: 16px;
}

.compare-card,
.compare-summary-card {
  border-radius: 16px;
}
</style>

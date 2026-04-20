<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">数据集访问审批</h1>
          <p class="section-subtitle">审核 `internal / restricted` 数据集的下载申请，通过后用户才能下载或生成分批下载任务。</p>
        </div>
      </div>

      <div class="filter-bar">
        <a-select v-model:value="filters.status" :options="statusOptions" style="width: 180px" />
        <a-select v-model:value="filters.privacy" :options="privacyOptions" style="width: 180px" />
        <a-select v-model:value="filters.sla_status" :options="slaStatusOptions" style="width: 180px" />
        <div class="owner-filter">
          <UserSearchSelect v-model="filters.owner_id" placeholder="搜索负责人" />
        </div>
        <a-input-search
          v-model:value="filters.q"
          placeholder="搜索数据集、申请人、负责人或申请理由"
          allow-clear
          style="max-width: 360px"
          @search="load"
        />
        <a-button type="primary" @click="load">筛选</a-button>
        <a-button @click="resetFilters">重置</a-button>
      </div>

      <div class="batch-bar">
        <a-space wrap>
          <a-tag>已选 {{ selectedRequestIds.length }} 项</a-tag>
          <a-tag color="orange">超时待审 {{ overduePendingCount }}</a-tag>
          <a-button :disabled="!pendingItems.length" @click="selectAllPending">全选当前待审</a-button>
          <a-button :disabled="!selectedRequestIds.length" @click="clearSelection">清空选择</a-button>
          <a-button type="primary" :disabled="!selectedRequestIds.length" @click="openReviewModal('batch', 'approved')">批量通过</a-button>
          <a-button danger :disabled="!selectedRequestIds.length" @click="openReviewModal('batch', 'rejected')">批量驳回</a-button>
        </a-space>
      </div>

      <a-list :data-source="items">
        <template #renderItem="{ item }">
          <a-list-item class="request-item">
            <a-list-item-meta>
              <template #title>
                <a-space wrap>
                  <span>{{ item.user_name || `用户 ${item.user_id}` }}</span>
                  <a-tag color="blue">{{ item.dataset_name || `数据集 ${item.dataset_id}` }}</a-tag>
                  <a-tag v-if="item.dataset_privacy" color="purple">{{ item.dataset_privacy }}</a-tag>
                  <a-tag v-if="item.dataset_owner_name" color="cyan">负责人 {{ item.dataset_owner_name }}</a-tag>
                  <a-tag :color="statusColor(item.status)">{{ item.status }}</a-tag>
                  <a-tag v-if="item.required_approvals > 1" color="geekblue">审批进度 {{ item.approval_stage }}/{{ item.required_approvals }}</a-tag>
                  <a-tag v-if="item.status === 'pending'" :color="item.sla_overdue ? 'red' : 'blue'">{{ slaStatusText(item) }}</a-tag>
                  <a-checkbox
                    v-if="item.status === 'pending'"
                    :checked="selectedRequestIds.includes(item.id)"
                    @change="toggleSelection(item.id, $event.target.checked)"
                  >
                    批量选择
                  </a-checkbox>
                </a-space>
              </template>
              <template #description>
                <div>申请理由：{{ item.reason }}</div>
                <div v-if="item.review_comment">审核意见：{{ item.review_comment }}</div>
                <div v-if="item.required_approvals > 1" class="section-subtitle">多级审批：当前已完成 {{ item.approval_stage }}/{{ item.required_approvals }} 级</div>
                <div v-if="item.status === 'approved'">
                  授权信息：{{ accessGrantSummary(item) }}
                </div>
                <div v-if="item.sla_hours > 0" class="section-subtitle">SLA：{{ slaSummary(item) }}</div>
                <div>数据集 ID：{{ item.dataset_id }}</div>
                <div>申请时间：{{ formatDate(item.created_at) }}</div>
              </template>
            </a-list-item-meta>
            <template v-if="item.status === 'pending'" #actions>
              <a-button type="link" @click="openReviewModal('single', 'approved', item.id)">通过</a-button>
              <a-button danger type="link" @click="openReviewModal('single', 'rejected', item.id)">驳回</a-button>
            </template>
          </a-list-item>
        </template>
      </a-list>

      <a-modal
        v-model:open="reviewModal.open"
        :title="reviewModal.mode === 'batch' ? '批量审核访问申请' : '审核访问申请'"
        ok-text="提交审核"
        cancel-text="取消"
        @ok="submitReview"
      >
        <a-form layout="vertical">
          <a-form-item label="审核结果">
            <a-tag :color="reviewModal.decision === 'approved' ? 'green' : 'red'">
              {{ reviewModal.decision === 'approved' ? '通过' : '驳回' }}
            </a-tag>
          </a-form-item>
          <a-form-item label="备注模板">
            <a-select :value="reviewModal.template" :options="reviewTemplateOptions" @change="applyReviewTemplate" />
          </a-form-item>
          <a-form-item label="审核说明">
            <a-textarea v-model:value="reviewModal.comment" :rows="4" placeholder="请输入审核说明" />
          </a-form-item>
          <template v-if="reviewModal.decision === 'approved'">
            <a-form-item label="有效期天数">
              <a-input-number v-model:value="reviewModal.valid_days" :min="0" :max="3650" style="width: 100%" />
              <div class="field-hint">`0` 表示长期有效</div>
            </a-form-item>
            <a-form-item label="最大下载次数">
              <a-input-number v-model:value="reviewModal.download_limit" :min="0" :max="10000" style="width: 100%" />
              <div class="field-hint">`0` 表示不限制次数；生成分批下载任务同样会消耗一次次数</div>
            </a-form-item>
          </template>
        </a-form>
      </a-modal>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';

import { api } from '@/api';
import UserSearchSelect from '@/components/UserSearchSelect.vue';
import type { DatasetAccessRequestItem } from '@/types/api';

const items = ref<DatasetAccessRequestItem[]>([]);
const selectedRequestIds = ref<number[]>([]);
const filters = reactive({
  status: '',
  privacy: '',
  sla_status: '',
  owner_id: undefined as number | undefined,
  q: '',
});
const statusOptions = [
  { label: '全部状态', value: '' },
  { label: '待审核', value: 'pending' },
  { label: '已通过', value: 'approved' },
  { label: '已驳回', value: 'rejected' },
];
const privacyOptions = [
  { label: '全部权限', value: '' },
  { label: '内部', value: 'internal' },
  { label: '受限', value: 'restricted' },
];
const slaStatusOptions = [
  { label: '全部 SLA', value: '' },
  { label: 'SLA 内', value: 'ontrack' },
  { label: '已超时', value: 'overdue' },
];
const pendingItems = computed(() => items.value.filter((item) => item.status === 'pending'));
const overduePendingCount = computed(() => pendingItems.value.filter((item) => item.sla_overdue).length);
const approvedReviewTemplates = [
  { label: '允许下载', value: '访问申请通过，可继续下载与生成分批下载任务。' },
  { label: '场景已确认', value: '场景说明完整，访问申请通过。' },
  { label: '联调准入', value: '资料已满足联调要求，允许访问该数据集。' },
];
const rejectedReviewTemplates = [
  { label: '缺少场景说明', value: '缺少明确的使用场景说明，请补充后重新申请。' },
  { label: '材料不完整', value: '申请材料不完整，请补充数据用途和联调范围后再提交。' },
  { label: '权限不匹配', value: '当前申请人与数据权限要求不匹配，请联系负责人补充说明。' },
];
const reviewTemplateOptions = computed(() => (reviewModal.decision === 'approved' ? approvedReviewTemplates : rejectedReviewTemplates));
const reviewModal = reactive({
  open: false,
  mode: 'single' as 'single' | 'batch',
  decision: 'approved' as 'approved' | 'rejected',
  targetId: undefined as number | undefined,
  template: '',
  comment: '',
  valid_days: 0,
  download_limit: 0,
});

async function load() {
  items.value = await api.getAdminDatasetAccessRequests({
    status: filters.status || undefined,
    privacy: filters.privacy || undefined,
    sla_status: filters.sla_status || undefined,
    owner_id: filters.owner_id || undefined,
    q: filters.q.trim() || undefined,
  });
  selectedRequestIds.value = selectedRequestIds.value.filter((id) => items.value.some((item) => item.id === id && item.status === 'pending'));
}

function openReviewModal(mode: 'single' | 'batch', decision: 'approved' | 'rejected', targetId?: number) {
  reviewModal.open = true;
  reviewModal.mode = mode;
  reviewModal.decision = decision;
  reviewModal.targetId = targetId;
  reviewModal.template = reviewTemplateOptions.value[0]?.value ?? '';
  reviewModal.comment = reviewModal.template;
  reviewModal.valid_days = 0;
  reviewModal.download_limit = 0;
}

function applyReviewTemplate(value: string) {
  reviewModal.template = value;
  reviewModal.comment = value;
}

async function submitReview() {
  const comment = reviewModal.comment.trim();
  if (!comment) {
    message.error('请填写审核说明');
    return;
  }

  if (reviewModal.mode === 'single' && reviewModal.targetId) {
    await api.reviewAdminDatasetAccessRequest(reviewModal.targetId, {
      decision: reviewModal.decision,
      comment,
      valid_days: reviewModal.valid_days,
      download_limit: reviewModal.download_limit,
    });
    message.success('访问申请状态已更新');
  } else {
    await api.batchReviewAdminDatasetAccessRequests({
      ids: selectedRequestIds.value,
      decision: reviewModal.decision,
      comment,
      valid_days: reviewModal.valid_days,
      download_limit: reviewModal.download_limit,
    });
    message.success('批量审核已完成');
  }

  reviewModal.open = false;
  reviewModal.targetId = undefined;
  await load();
}

function toggleSelection(id: number, checked: boolean) {
  if (checked) {
    if (!selectedRequestIds.value.includes(id)) {
      selectedRequestIds.value = [...selectedRequestIds.value, id];
    }
    return;
  }
  selectedRequestIds.value = selectedRequestIds.value.filter((item) => item !== id);
}

function selectAllPending() {
  selectedRequestIds.value = pendingItems.value.map((item) => item.id);
}

function clearSelection() {
  selectedRequestIds.value = [];
}

function resetFilters() {
  filters.status = '';
  filters.privacy = '';
  filters.sla_status = '';
  filters.owner_id = undefined;
  filters.q = '';
  clearSelection();
  load();
}

function statusColor(status: string) {
  return {
    pending: 'gold',
    approved: 'green',
    rejected: 'red',
  }[status] ?? 'default';
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('zh-CN');
}

function accessGrantSummary(item: DatasetAccessRequestItem) {
  const limitText = item.download_limit > 0 ? `${item.download_count}/${item.download_limit}` : '不限次数';
  const expireText = item.approval_expires_at ? formatDate(item.approval_expires_at) : '长期有效';
  return `${limitText} · 有效期至 ${expireText}`;
}

function slaStatusText(item: DatasetAccessRequestItem) {
  if (item.sla_overdue) {
    return `SLA 超时 ${formatMinutes(-item.sla_remaining_minutes)}`;
  }
  return `SLA 剩余 ${formatMinutes(item.sla_remaining_minutes)}`;
}

function slaSummary(item: DatasetAccessRequestItem) {
  if (!item.sla_hours || !item.sla_deadline_at) {
    return '无 SLA 要求';
  }
  return `${item.sla_hours} 小时 · 截止 ${formatDate(item.sla_deadline_at)} · ${item.sla_overdue ? `已超时 ${formatMinutes(-item.sla_remaining_minutes)}` : `剩余 ${formatMinutes(item.sla_remaining_minutes)}`}`;
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

.filter-bar {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  align-items: center;
  margin-bottom: 16px;
}

.owner-filter {
  width: 220px;
}

.batch-bar {
  margin-bottom: 16px;
}

.field-hint {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 12px;
}

.request-item {
  border-radius: 16px;
  border: 1px solid var(--line);
  margin-bottom: 12px;
  background: var(--surface-soft);
}
</style>

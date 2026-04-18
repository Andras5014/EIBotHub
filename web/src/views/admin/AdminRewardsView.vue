<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">积分运营</h1>
          <p class="section-subtitle">维护兑换权益、人工修正积分，并查看最近的积分运营记录。</p>
        </div>
      </div>

      <div class="grid-cards" v-if="overview">
        <a-card><a-statistic title="权益总数" :value="overview.benefits" /></a-card>
        <a-card><a-statistic title="启用权益" :value="overview.active_benefits" /></a-card>
        <a-card><a-statistic title="兑换记录" :value="overview.redemptions" /></a-card>
        <a-card><a-statistic title="账本记录" :value="overview.ledger_entries" /></a-card>
        <a-card><a-statistic title="净积分" :value="overview.net_points" /></a-card>
      </div>

      <a-row :gutter="[16, 16]" style="margin-top: 18px">
        <a-col :xs="24" :lg="13">
          <a-card title="权益配置" class="inner-card">
            <a-list :data-source="benefits">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta>
                    <template #title>
                      <a-space wrap>
                        <span>{{ item.name }}</span>
                        <a-tag :color="item.active ? 'green' : 'default'">{{ item.active ? '启用' : '停用' }}</a-tag>
                        <a-tag color="gold">{{ item.cost_points }} 积分</a-tag>
                      </a-space>
                    </template>
                    <template #description>{{ item.summary }}</template>
                  </a-list-item-meta>
                  <template #actions>
                    <a-button type="link" @click="openEdit(item)">编辑</a-button>
                  </template>
                </a-list-item>
              </template>
            </a-list>
          </a-card>

          <a-card title="贡献排行榜" class="inner-card" style="margin-top: 16px">
            <a-list :data-source="rankings">
              <template #renderItem="{ item, index }">
                <a-list-item>{{ index + 1 }}. {{ item.user_name }} · {{ item.points }}</a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="11">
          <a-card title="人工积分修正" class="inner-card">
            <a-form :model="adjustForm" layout="vertical" @finish="submitAdjustment">
              <a-form-item label="目标用户" name="user_id">
                <UserSearchSelect v-model="adjustForm.user_id" placeholder="搜索用户名" />
              </a-form-item>
              <a-form-item label="积分变更值" name="points">
                <a-input-number v-model:value="adjustForm.points" style="width: 100%" :min="-100000" :max="100000" />
              </a-form-item>
              <a-form-item label="变更说明" name="remark">
                <a-textarea v-model:value="adjustForm.remark" :rows="3" />
              </a-form-item>
              <a-form-item>
                <a-button type="primary" html-type="submit" :loading="adjusting">提交修正</a-button>
              </a-form-item>
            </a-form>
          </a-card>

          <a-card title="最近修正记录" class="inner-card" style="margin-top: 16px">
            <a-list :data-source="adjustments">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta
                    :title="`${item.user_name} · ${item.points > 0 ? '+' : ''}${item.points}`"
                    :description="`${item.remark} · ${formatDate(item.created_at)}`"
                  />
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>
      </a-row>
    </div>

    <a-modal v-model:open="editOpen" title="编辑权益" ok-text="保存" cancel-text="取消" @ok="saveBenefit">
      <a-form :model="editForm" layout="vertical">
        <a-form-item label="权益名称"><a-input v-model:value="editForm.name" /></a-form-item>
        <a-form-item label="权益摘要"><a-textarea v-model:value="editForm.summary" :rows="3" /></a-form-item>
        <a-form-item label="所需积分">
          <a-input-number v-model:value="editForm.cost_points" style="width: 100%" :min="1" :max="100000" />
        </a-form-item>
        <a-form-item label="启用状态">
          <a-switch v-model:checked="editForm.active" checked-children="启用" un-checked-children="停用" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';

import { api } from '@/api';
import UserSearchSelect from '@/components/UserSearchSelect.vue';
import type {
  AdminRewardAdjustmentItem,
  AdminRewardOverview,
  ContributorRankingItem,
  RewardBenefitItem,
} from '@/types/api';

const overview = ref<AdminRewardOverview>();
const benefits = ref<RewardBenefitItem[]>([]);
const rankings = ref<ContributorRankingItem[]>([]);
const adjustments = ref<AdminRewardAdjustmentItem[]>([]);
const adjusting = ref(false);
const editOpen = ref(false);
const editingBenefitId = ref<number>();

const adjustForm = reactive({
  user_id: undefined as number | undefined,
  points: 0,
  remark: '',
});

const editForm = reactive({
  name: '',
  summary: '',
  cost_points: 1,
  active: true,
});

async function load() {
  const [overviewData, benefitData, rankingData, adjustmentData] = await Promise.all([
    api.getAdminRewardOverview(),
    api.getAdminRewardBenefits(),
    api.contributorRankings(),
    api.getAdminRewardAdjustments(),
  ]);
  overview.value = overviewData;
  benefits.value = benefitData;
  rankings.value = rankingData;
  adjustments.value = adjustmentData;
}

function openEdit(item: RewardBenefitItem) {
  editingBenefitId.value = item.id;
  editForm.name = item.name;
  editForm.summary = item.summary;
  editForm.cost_points = item.cost_points;
  editForm.active = item.active;
  editOpen.value = true;
}

async function saveBenefit() {
  if (!editingBenefitId.value) {
    return;
  }
  await api.updateAdminRewardBenefit(editingBenefitId.value, editForm);
  message.success('权益配置已更新');
  editOpen.value = false;
  await load();
}

async function submitAdjustment() {
  adjusting.value = true;
  try {
    await api.createAdminRewardAdjustment({
      user_id: adjustForm.user_id ?? 0,
      points: adjustForm.points,
      remark: adjustForm.remark,
    });
    adjustForm.user_id = undefined;
    adjustForm.points = 0;
    adjustForm.remark = '';
    message.success('积分修正已写入');
    await load();
  } catch (error) {
    message.error(error instanceof Error ? error.message : '积分修正失败');
  } finally {
    adjusting.value = false;
  }
}

function formatDate(value: string) {
  return new Date(value).toLocaleString();
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

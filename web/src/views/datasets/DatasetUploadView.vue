<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">{{ isEditMode ? '编辑数据集' : '上传数据集' }}</h1>
          <p class="section-subtitle">
            {{ isEditMode ? '保存后状态会重置为草稿，可继续提交审核。' : '上传完成后自动提交审核，协议文本会在详情页展示。' }}
          </p>
        </div>
      </div>

      <a-alert
        v-if="isEditMode"
        type="info"
        show-icon
        message="当前页面用于编辑数据集元数据与文本样本说明；版本文件和样本文件维护仍通过版本与上传流程处理。"
        style="margin-bottom: 16px"
      />

      <a-form ref="formRef" :model="form" layout="vertical">
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="12"><a-form-item label="数据集名称" name="name" :rules="requiredRule('请输入数据集名称')"><a-input v-model:value="form.name" /></a-form-item></a-col>
          <a-col :xs="24" :lg="12"><a-form-item label="摘要" name="summary" :rules="requiredRule('请输入摘要')"><a-input v-model:value="form.summary" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="描述" name="description" :rules="requiredRule('请输入描述')"><a-textarea v-model:value="form.description" :rows="4" /></a-form-item>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="8">
            <a-form-item label="标签" name="tags">
              <a-select
                v-model:value="form.tags"
                mode="tags"
                allow-clear
                show-search
                :options="datasetTagOptions"
                placeholder="选择或输入标签"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="样本量" name="sampleCount" :rules="requiredRule('请输入样本量')"><a-input v-model:value="form.sampleCount" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="采集设备" name="device" :rules="requiredRule('请输入采集设备')"><a-input v-model:value="form.device" /></a-form-item></a-col>
        </a-row>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="8">
            <a-form-item label="场景" name="scene" :rules="requiredRule('请选择场景')">
              <a-select
                v-model:value="form.scene"
                show-search
                allow-clear
                :options="sceneOptions"
                placeholder="选择数据集场景"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :lg="8">
            <a-form-item label="权限等级" name="privacy" :rules="requiredRule('请选择权限等级')">
              <a-select v-model:value="form.privacy" :options="privacyOptions.map((item) => ({ label: `${item.name} · ${item.description}`, value: item.code }))" />
            </a-form-item>
          </a-col>
          <a-col v-if="!isEditMode" :xs="24" :lg="8"><a-form-item label="版本" name="version" :rules="requiredRule('请输入版本号')"><a-input v-model:value="form.version" /></a-form-item></a-col>
        </a-row>
        <a-form-item v-if="!isEditMode" label="版本说明" name="changelog" :rules="requiredRule('请输入版本说明')"><a-input v-model:value="form.changelog" /></a-form-item>
        <a-form-item label="协议模板">
          <a-space>
            <a-select v-model:value="selectedAgreementTemplateId" style="min-width: 260px" :options="agreementTemplates.map((item) => ({ label: item.name, value: item.id }))" />
            <a-button @click="applyAgreementTemplate(selectedAgreementTemplateId)">应用模板</a-button>
          </a-space>
        </a-form-item>
        <a-form-item label="协议文本" name="agreementText" :rules="requiredRule('请输入协议文本')"><a-textarea v-model:value="form.agreementText" :rows="3" /></a-form-item>
        <a-form-item label="样本预览文本（每行一条）" name="samplePreview"><a-textarea v-model:value="form.samplePreview" :rows="3" /></a-form-item>
        <template v-if="!isEditMode">
          <a-form-item label="样本预览文件"><input type="file" multiple accept="image/*,video/*,.ply,.pcd,.las,.laz,.obj" @change="onSampleFilesChange" /></a-form-item>
          <a-form-item label="数据集文件"><input type="file" @change="onFileChange" /></a-form-item>
        </template>
        <a-form-item>
          <a-space>
            <a-button type="primary" :loading="submitting" @click="submitPrimary">
              {{ isEditMode ? '保存修改' : '上传并提交审核' }}
            </a-button>
            <a-button v-if="isEditMode" :loading="submittingAndSubmitting" @click="submitAndReview">保存并提交审核</a-button>
          </a-space>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import type { FormInstance } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import type { AgreementTemplateItem, DatasetPrivacyOptionItem, FilterOptionsResponse } from '@/types/api';

const route = useRoute();
const router = useRouter();
const formRef = ref<FormInstance>();
const submitting = ref(false);
const submittingAndSubmitting = ref(false);
const file = ref<File>();
const sampleFiles = ref<File[]>([]);
const agreementTemplates = ref<AgreementTemplateItem[]>([]);
const privacyOptions = ref<DatasetPrivacyOptionItem[]>([]);
const filterOptions = ref<FilterOptionsResponse>({
  tags: [],
  model_tags: [],
  dataset_tags: [],
  robot_types: [],
  dataset_scenes: [],
  template_categories: [],
  template_scenes: [],
  application_case_categories: [],
});
const selectedAgreementTemplateId = ref<number>();
const editingId = computed(() => Number(route.params.id ?? 0));
const isEditMode = computed(() => editingId.value > 0);
const form = reactive({
  name: '',
  summary: '',
  description: '',
  tags: [] as string[],
  sampleCount: '10',
  device: 'RGB Camera',
  scene: '巡逻巡检',
  privacy: 'public',
  version: 'v1.0.0',
  changelog: 'initial release',
  agreementText: '下载并使用本数据集前需同意开放社区使用协议。',
  samplePreview: '示例文本 1\n示例文本 2',
});
const datasetTagOptions = computed(() => filterOptions.value.dataset_tags.map((item) => ({ label: item, value: item })));
const sceneOptions = computed(() => {
  const values = [...filterOptions.value.dataset_scenes];
  if (form.scene && !values.includes(form.scene)) {
    values.unshift(form.scene);
  }
  return values.map((item) => ({ label: item, value: item }));
});

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement;
  file.value = target.files?.[0];
}

function onSampleFilesChange(event: Event) {
  const target = event.target as HTMLInputElement;
  sampleFiles.value = Array.from(target.files ?? []);
}

function applyAgreementTemplate(templateId?: number) {
  if (!templateId) return;
  const target = agreementTemplates.value.find((item) => item.id === templateId);
  if (target) {
    form.agreementText = target.content;
  }
}

async function submitPrimary() {
  await handleSubmit(false);
}

async function submitAndReview() {
  await handleSubmit(true);
}

async function handleSubmit(submitAfterSave: boolean) {
  try {
    await formRef.value?.validate();
  } catch {
    return;
  }
  if (!isEditMode.value && !file.value) {
    message.error('请选择数据集文件');
    return;
  }

  if (submitAfterSave) {
    submittingAndSubmitting.value = true;
  } else {
    submitting.value = true;
  }

  try {
    if (isEditMode.value) {
      const updated = await api.updateDataset(editingId.value, {
        name: form.name,
        summary: form.summary,
        description: form.description,
        tags: form.tags.join(','),
        sample_count: Number(form.sampleCount),
        device: form.device,
        scene: form.scene,
        privacy: form.privacy,
        agreement_text: form.agreementText,
        sample_preview: form.samplePreview,
      });
      if (submitAfterSave) {
        await api.submitDataset(updated.id);
        message.success('数据集已保存并重新提交审核');
      } else {
        message.success('数据集信息已保存');
      }
      await router.push(`/datasets/${updated.id}`);
      return;
    }

    const payload = new FormData();
    payload.append('name', form.name);
    payload.append('summary', form.summary);
    payload.append('description', form.description);
    payload.append('tags', form.tags.join(','));
    payload.append('sample_count', form.sampleCount);
    payload.append('device', form.device);
    payload.append('scene', form.scene);
    payload.append('privacy', form.privacy);
    payload.append('version', form.version);
    payload.append('changelog', form.changelog);
    payload.append('agreement_text', form.agreementText);
    payload.append('sample_preview', form.samplePreview);
    sampleFiles.value.forEach((sampleFile) => payload.append('sample_files', sampleFile));
    if (file.value) payload.append('file', file.value);
    const created = await api.createDataset(payload);
    await api.submitDataset(created.id);
    message.success('数据集已上传并提交审核');
    await router.push(`/datasets/${created.id}`);
  } catch (error) {
    message.error(error instanceof Error ? error.message : isEditMode.value ? '数据集保存失败' : '数据集上传失败');
  } finally {
    submitting.value = false;
    submittingAndSubmitting.value = false;
  }
}

function requiredRule(messageText: string) {
  return [{ required: true, message: messageText }];
}

onMounted(async () => {
  const [options, filterOptionData] = await Promise.all([api.getDatasetOptions(), api.getFilterOptions()]);
  agreementTemplates.value = options.agreement_templates;
  privacyOptions.value = options.privacy_options;
  filterOptions.value = filterOptionData;
  if (agreementTemplates.value.length) {
    selectedAgreementTemplateId.value = agreementTemplates.value[0].id;
  }

  if (!isEditMode.value) {
    return;
  }

  const detail = await api.getDataset(editingId.value);
  form.name = detail.name;
  form.summary = detail.summary;
  form.description = detail.description ?? '';
  form.tags = detail.tags;
  form.sampleCount = String(detail.sample_count);
  form.device = detail.device;
  form.scene = detail.scene;
  form.privacy = detail.privacy;
  form.agreementText = detail.agreement_text;
  form.samplePreview = detail.samples
    .filter((item) => item.sample_type === 'text')
    .map((item) => item.preview_text)
    .join('\n');
});
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

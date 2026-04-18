<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">上传数据集</h1>
          <p class="section-subtitle">上传完成后自动提交审核，协议文本会在详情页展示。</p>
        </div>
      </div>

      <a-form :model="form" layout="vertical" @finish="submit">
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="12"><a-form-item label="数据集名称" name="name" :rules="requiredRule('请输入数据集名称')"><a-input v-model:value="form.name" /></a-form-item></a-col>
          <a-col :xs="24" :lg="12"><a-form-item label="摘要" name="summary" :rules="requiredRule('请输入摘要')"><a-input v-model:value="form.summary" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="描述" name="description" :rules="requiredRule('请输入描述')"><a-textarea v-model:value="form.description" :rows="4" /></a-form-item>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="8"><a-form-item label="标签" name="tags"><a-input v-model:value="form.tags" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="样本量" name="sampleCount" :rules="requiredRule('请输入样本量')"><a-input v-model:value="form.sampleCount" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="采集设备" name="device" :rules="requiredRule('请输入采集设备')"><a-input v-model:value="form.device" /></a-form-item></a-col>
        </a-row>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="8"><a-form-item label="场景" name="scene" :rules="requiredRule('请输入场景')"><a-input v-model:value="form.scene" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8">
            <a-form-item label="权限等级" name="privacy" :rules="requiredRule('请选择权限等级')">
              <a-select v-model:value="form.privacy" :options="privacyOptions.map((item) => ({ label: `${item.name} · ${item.description}`, value: item.code }))" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="版本" name="version" :rules="requiredRule('请输入版本号')"><a-input v-model:value="form.version" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="版本说明" name="changelog" :rules="requiredRule('请输入版本说明')"><a-input v-model:value="form.changelog" /></a-form-item>
        <a-form-item label="协议模板">
          <a-space>
            <a-select v-model:value="selectedAgreementTemplateId" style="min-width: 260px" :options="agreementTemplates.map((item) => ({ label: item.name, value: item.id }))" />
            <a-button @click="applyAgreementTemplate(selectedAgreementTemplateId)">应用模板</a-button>
          </a-space>
        </a-form-item>
        <a-form-item label="协议文本" name="agreementText" :rules="requiredRule('请输入协议文本')"><a-textarea v-model:value="form.agreementText" :rows="3" /></a-form-item>
        <a-form-item label="样本预览文本（每行一条）" name="samplePreview"><a-textarea v-model:value="form.samplePreview" :rows="3" /></a-form-item>
        <a-form-item label="样本预览文件"><input type="file" multiple accept="image/*,video/*,.ply,.pcd,.las,.laz,.obj" @change="onSampleFilesChange" /></a-form-item>
        <a-form-item label="数据集文件"><input type="file" @change="onFileChange" /></a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="submitting">上传并提交审核</a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { useRouter } from 'vue-router';

import { api } from '@/api';
import type { AgreementTemplateItem, DatasetPrivacyOptionItem } from '@/types/api';

const router = useRouter();
const submitting = ref(false);
const file = ref<File>();
const sampleFiles = ref<File[]>([]);
const agreementTemplates = ref<AgreementTemplateItem[]>([]);
const privacyOptions = ref<DatasetPrivacyOptionItem[]>([]);
const selectedAgreementTemplateId = ref<number>();
const form = reactive({
  name: '',
  summary: '',
  description: '',
  tags: '',
  sampleCount: '10',
  device: 'RGB Camera',
  scene: '巡逻巡检',
  privacy: 'public',
  version: 'v1.0.0',
  changelog: 'initial release',
  agreementText: '下载并使用本数据集前需同意开放社区使用协议。',
  samplePreview: '示例文本 1\n示例文本 2',
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

async function submit() {
  if (!file.value) {
    message.error('请选择数据集文件');
    return;
  }
  submitting.value = true;
  try {
    const payload = new FormData();
    payload.append('name', form.name);
    payload.append('summary', form.summary);
    payload.append('description', form.description);
    payload.append('tags', form.tags);
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
    message.error(error instanceof Error ? error.message : '数据集上传失败');
  } finally {
    submitting.value = false;
  }
}

function requiredRule(messageText: string) {
  return [{ required: true, message: messageText }];
}

onMounted(async () => {
  const options = await api.getDatasetOptions();
  agreementTemplates.value = options.agreement_templates;
  privacyOptions.value = options.privacy_options;
  if (agreementTemplates.value.length) {
    selectedAgreementTemplateId.value = agreementTemplates.value[0].id;
  }
});
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">{{ isEditMode ? '编辑模型' : '上传模型' }}</h1>
          <p class="section-subtitle">
            {{ isEditMode ? '保存后状态会重置为草稿，可继续提交审核。' : '上传完成后会自动提交到后台审核。' }}
          </p>
        </div>
      </div>

      <a-alert
        v-if="isEditMode"
        type="info"
        show-icon
        message="当前页面用于编辑模型元数据；模型文件与版本变更请使用版本管理能力。"
        style="margin-bottom: 16px"
      />

      <a-form ref="formRef" :model="form" layout="vertical">
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="12"><a-form-item label="模型名称" name="name" :rules="requiredRule('请输入模型名称')"><a-input v-model:value="form.name" /></a-form-item></a-col>
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
                :options="tagOptions"
                placeholder="选择或输入标签"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :lg="8">
            <a-form-item label="机器人类型" name="robotType" :rules="requiredRule('请选择机器人类型')">
              <a-select
                v-model:value="form.robotType"
                show-search
                allow-clear
                :options="robotTypeOptions"
                placeholder="选择机器人类型"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :lg="8">
            <a-form-item label="许可证" name="license" :rules="requiredRule('请选择许可证')">
              <a-select
                v-model:value="form.license"
                show-search
                :options="licenseOptions"
                placeholder="选择许可证"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="8"><a-form-item label="输入规格" name="inputSpec" :rules="requiredRule('请输入输入规格')"><a-input v-model:value="form.inputSpec" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="输出规格" name="outputSpec" :rules="requiredRule('请输入输出规格')"><a-input v-model:value="form.outputSpec" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8">
            <a-form-item label="依赖" name="dependencies">
              <a-select
                v-model:value="form.dependencies"
                mode="tags"
                allow-clear
                show-search
                :options="dependencyOptions"
                placeholder="选择或输入依赖"
              />
            </a-form-item>
          </a-col>
        </a-row>
        <template v-if="!isEditMode">
          <a-row :gutter="[16, 0]">
            <a-col :xs="24" :lg="8"><a-form-item label="版本" name="version" :rules="requiredRule('请输入版本号')"><a-input v-model:value="form.version" /></a-form-item></a-col>
            <a-col :xs="24" :lg="16"><a-form-item label="版本说明" name="changelog" :rules="requiredRule('请输入版本说明')"><a-input v-model:value="form.changelog" /></a-form-item></a-col>
          </a-row>
          <a-form-item label="模型文件">
            <input type="file" @change="onFileChange" />
          </a-form-item>
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
import type { FilterOptionsResponse } from '@/types/api';

const route = useRoute();
const router = useRouter();
const formRef = ref<FormInstance>();
const submitting = ref(false);
const submittingAndSubmitting = ref(false);
const file = ref<File>();
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
const editingId = computed(() => Number(route.params.id ?? 0));
const isEditMode = computed(() => editingId.value > 0);
const form = reactive({
  name: '',
  summary: '',
  description: '',
  tags: [] as string[],
  robotType: '',
  inputSpec: '',
  outputSpec: '',
  license: 'Apache-2.0',
  dependencies: [] as string[],
  version: 'v1.0.0',
  changelog: 'initial release',
});
const tagOptions = computed(() => filterOptions.value.model_tags.map((item) => ({ label: item, value: item })));
const robotTypeOptions = computed(() => {
  const values = [...filterOptions.value.robot_types];
  if (form.robotType && !values.includes(form.robotType)) {
    values.unshift(form.robotType);
  }
  return values.map((item) => ({ label: item, value: item }));
});
const licenseOptions = computed(() =>
  ['Apache-2.0', 'MIT', 'BSD-3-Clause', 'GPL-3.0', 'MPL-2.0', 'Proprietary'].map((item) => ({
    label: item,
    value: item,
  })),
);
const dependencyOptions = computed(() => {
  const discovered = ['opencv', 'onnxruntime', 'pytorch', 'tensorrt', 'ffmpeg', 'numpy'];
  const values = [...discovered];
  for (const item of form.dependencies) {
    if (!values.includes(item)) {
      values.unshift(item);
    }
  }
  return values.map((item) => ({ label: item, value: item }));
});

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement;
  file.value = target.files?.[0];
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
    message.error('请选择模型文件');
    return;
  }

  if (submitAfterSave) {
    submittingAndSubmitting.value = true;
  } else {
    submitting.value = true;
  }

  try {
    if (isEditMode.value) {
      const updated = await api.updateModel(editingId.value, {
        name: form.name,
        summary: form.summary,
        description: form.description,
        tags: form.tags.join(','),
        robot_type: form.robotType,
        input_spec: form.inputSpec,
        output_spec: form.outputSpec,
        license: form.license,
        dependencies: form.dependencies.join(','),
      });
      if (submitAfterSave) {
        await api.submitModel(updated.id);
        message.success('模型已保存并重新提交审核');
      } else {
        message.success('模型信息已保存');
      }
      await router.push(`/models/${updated.id}`);
      return;
    }

    const payload = new FormData();
    payload.append('name', form.name);
    payload.append('summary', form.summary);
    payload.append('description', form.description);
    payload.append('tags', form.tags.join(','));
    payload.append('robot_type', form.robotType);
    payload.append('input_spec', form.inputSpec);
    payload.append('output_spec', form.outputSpec);
    payload.append('license', form.license);
    payload.append('dependencies', form.dependencies.join(','));
    payload.append('version', form.version);
    payload.append('changelog', form.changelog);
    if (file.value) payload.append('file', file.value);
    const created = await api.createModel(payload);
    await api.submitModel(created.id);
    message.success('模型已上传并提交审核');
    await router.push(`/models/${created.id}`);
  } catch (error) {
    message.error(error instanceof Error ? error.message : isEditMode.value ? '模型保存失败' : '模型上传失败');
  } finally {
    submitting.value = false;
    submittingAndSubmitting.value = false;
  }
}

function requiredRule(messageText: string) {
  return [{ required: true, message: messageText }];
}

onMounted(async () => {
  filterOptions.value = await api.getFilterOptions();
  if (!isEditMode.value) {
    return;
  }
  const detail = await api.getModel(editingId.value);
  form.name = detail.name;
  form.summary = detail.summary;
  form.description = detail.description ?? '';
  form.tags = detail.tags;
  form.robotType = detail.robot_type ?? '';
  form.inputSpec = detail.input_spec;
  form.outputSpec = detail.output_spec;
  form.license = detail.license;
  form.dependencies = detail.dependencies;
});
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

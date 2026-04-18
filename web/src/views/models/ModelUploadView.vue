<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">上传模型</h1>
          <p class="section-subtitle">上传完成后会自动提交到后台审核。</p>
        </div>
      </div>

      <a-form :model="form" layout="vertical" @finish="submit">
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="12"><a-form-item label="模型名称" name="name" :rules="requiredRule('请输入模型名称')"><a-input v-model:value="form.name" /></a-form-item></a-col>
          <a-col :xs="24" :lg="12"><a-form-item label="摘要" name="summary" :rules="requiredRule('请输入摘要')"><a-input v-model:value="form.summary" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="描述" name="description" :rules="requiredRule('请输入描述')"><a-textarea v-model:value="form.description" :rows="4" /></a-form-item>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="8"><a-form-item label="标签" name="tags"><a-input v-model:value="form.tags" placeholder="逗号分隔" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="机器人类型" name="robotType" :rules="requiredRule('请输入机器人类型')"><a-input v-model:value="form.robotType" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="许可证" name="license" :rules="requiredRule('请输入许可证')"><a-input v-model:value="form.license" /></a-form-item></a-col>
        </a-row>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="8"><a-form-item label="输入规格" name="inputSpec" :rules="requiredRule('请输入输入规格')"><a-input v-model:value="form.inputSpec" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="输出规格" name="outputSpec" :rules="requiredRule('请输入输出规格')"><a-input v-model:value="form.outputSpec" /></a-form-item></a-col>
          <a-col :xs="24" :lg="8"><a-form-item label="依赖" name="dependencies"><a-input v-model:value="form.dependencies" placeholder="逗号分隔" /></a-form-item></a-col>
        </a-row>
        <a-row :gutter="[16, 0]">
          <a-col :xs="24" :lg="8"><a-form-item label="版本" name="version" :rules="requiredRule('请输入版本号')"><a-input v-model:value="form.version" /></a-form-item></a-col>
          <a-col :xs="24" :lg="16"><a-form-item label="版本说明" name="changelog" :rules="requiredRule('请输入版本说明')"><a-input v-model:value="form.changelog" /></a-form-item></a-col>
        </a-row>
        <a-form-item label="模型文件">
          <input type="file" @change="onFileChange" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="submitting">上传并提交审核</a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { useRouter } from 'vue-router';

import { api } from '@/api';

const router = useRouter();
const submitting = ref(false);
const file = ref<File>();
const form = reactive({
  name: '',
  summary: '',
  description: '',
  tags: '',
  robotType: '',
  inputSpec: '',
  outputSpec: '',
  license: 'Apache-2.0',
  dependencies: '',
  version: 'v1.0.0',
  changelog: 'initial release',
});

function onFileChange(event: Event) {
  const target = event.target as HTMLInputElement;
  file.value = target.files?.[0];
}

async function submit() {
  if (!file.value) {
    message.error('请选择模型文件');
    return;
  }
  submitting.value = true;
  try {
    const payload = new FormData();
    payload.append('name', form.name);
    payload.append('summary', form.summary);
    payload.append('description', form.description);
    payload.append('tags', form.tags);
    payload.append('robot_type', form.robotType);
    payload.append('input_spec', form.inputSpec);
    payload.append('output_spec', form.outputSpec);
    payload.append('license', form.license);
    payload.append('dependencies', form.dependencies);
    payload.append('version', form.version);
    payload.append('changelog', form.changelog);
    if (file.value) payload.append('file', file.value);
    const created = await api.createModel(payload);
    await api.submitModel(created.id);
    message.success('模型已上传并提交审核');
    await router.push(`/models/${created.id}`);
  } catch (error) {
    message.error(error instanceof Error ? error.message : '模型上传失败');
  } finally {
    submitting.value = false;
  }
}

function requiredRule(messageText: string) {
  return [{ required: true, message: messageText }];
}
</script>

<style scoped>
.block {
  padding: 24px;
}
</style>

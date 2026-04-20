<template>
  <div class="page-shell">
    <div class="auth-wrap">
      <div class="page-card auth-card">
        <div class="login-brand">
          <BrandLogo />
        </div>
        <h1 class="section-title">登录 EIBotHub</h1>
        <p class="section-subtitle">默认账号：`demo@example.com / Demo123!`，后台账号：`admin@opencommunity.local / Admin123!`</p>

        <a-form :model="form" layout="vertical" @finish="submit">
          <a-form-item label="邮箱" name="email" :rules="[{ required: true, message: '请输入邮箱' }]">
            <a-input v-model:value="form.email" placeholder="请输入邮箱" />
          </a-form-item>
          <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
            <a-input-password v-model:value="form.password" placeholder="请输入密码" />
          </a-form-item>
          <a-form-item>
            <a-space>
              <a-button type="primary" html-type="submit" :loading="submitting">登录</a-button>
              <a-button @click="registering = !registering">
                {{ registering ? '切回登录' : '我要注册' }}
              </a-button>
            </a-space>
          </a-form-item>
        </a-form>

        <a-form v-if="registering" :model="registerForm" layout="vertical" @finish="register">
          <a-divider>注册账号</a-divider>
          <a-form-item label="用户名" name="username" :rules="[{ required: true, message: '请输入用户名' }]">
            <a-input v-model:value="registerForm.username" />
          </a-form-item>
          <a-form-item label="邮箱" name="email" :rules="[{ required: true, message: '请输入邮箱' }]">
            <a-input v-model:value="registerForm.email" />
          </a-form-item>
          <a-form-item label="密码" name="password" :rules="[{ required: true, message: '请输入密码' }]">
            <a-input-password v-model:value="registerForm.password" />
          </a-form-item>
          <a-form-item>
            <a-button type="primary" html-type="submit" :loading="submitting">注册并登录</a-button>
          </a-form-item>
        </a-form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import { message } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import BrandLogo from '@/components/BrandLogo.vue';
import { useAuthStore } from '@/stores/auth';

const auth = useAuthStore();
const route = useRoute();
const router = useRouter();
const submitting = ref(false);
const registering = ref(false);

const form = reactive({
  email: 'demo@example.com',
  password: 'Demo123!',
});

const registerForm = reactive({
  username: 'new-user',
  email: 'new-user@example.com',
  password: 'Change123!',
});

async function submit() {
  submitting.value = true;
  try {
    await auth.login(form);
    message.success('登录成功');
    await router.push(String(route.query.redirect ?? '/'));
  } catch (error) {
    message.error(error instanceof Error ? error.message : '登录失败');
  } finally {
    submitting.value = false;
  }
}

async function register() {
  submitting.value = true;
  try {
    await auth.register(registerForm);
    message.success('注册成功');
    await router.push('/');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '注册失败');
  } finally {
    submitting.value = false;
  }
}
</script>

<style scoped>
.auth-wrap {
  min-height: calc(100vh - 180px);
  display: grid;
  place-items: center;
}

.auth-card {
  width: min(560px, 100%);
  padding: 28px;
}

.login-brand {
  margin-bottom: 20px;
}
</style>

<template>
  <a-layout style="min-height: 100vh">
    <a-layout-header class="shell-header">
      <div class="shell-inner">
        <RouterLink class="brand" to="/">
          <span class="brand-mark">OC</span>
          <div>
            <div class="brand-title">开放社区</div>
            <div class="brand-subtitle">Embodied AI Portal</div>
          </div>
        </RouterLink>

        <div class="nav-links">
          <RouterLink v-for="item in navItems" :key="item.to" class="nav-link" :to="item.to">
            {{ item.label }}
          </RouterLink>
        </div>

        <a-space :size="12" class="nav-actions">
          <RouterLink to="/search">
            <a-button type="default">全局搜索</a-button>
          </RouterLink>
          <RouterLink v-if="auth.isAdmin" to="/admin">
            <a-button type="default">后台</a-button>
          </RouterLink>
          <RouterLink v-if="auth.isAuthenticated" to="/me">
            <a-badge :count="unreadCount" :offset="[4, 2]">
              <a-button type="primary">{{ auth.user?.username ?? '个人中心' }}</a-button>
            </a-badge>
          </RouterLink>
          <RouterLink v-else to="/login">
            <a-button type="primary">登录 / 注册</a-button>
          </RouterLink>
        </a-space>
      </div>
    </a-layout-header>

    <a-layout-content>
      <RouterView />
    </a-layout-content>
  </a-layout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { RouterLink, RouterView, useRoute } from 'vue-router';

import { api } from '@/api';
import { useAuthStore } from '@/stores/auth';
import type { NotificationRecord } from '@/types/api';

const auth = useAuthStore();
const route = useRoute();
const notifications = ref<NotificationRecord[]>([]);

const unreadCount = computed(() => notifications.value.filter((item) => !item.Read).length);

const navItems = [
  { label: '首页', to: '/' },
  { label: '模型仓库', to: '/models' },
  { label: '数据集', to: '/datasets' },
  { label: '任务模板', to: '/templates' },
  { label: '技能分享', to: '/skills' },
  { label: '互动讨论', to: '/discussions' },
  { label: '私信', to: '/messages' },
  { label: '协作空间', to: '/workspaces' },
  { label: '具身案例', to: '/applications' },
  { label: '文档中心', to: '/docs' },
  { label: '开发者接入', to: '/developers' },
];

async function loadNotifications() {
  if (!auth.isAuthenticated) {
    notifications.value = [];
    return;
  }
  try {
    notifications.value = await api.getNotifications();
  } catch {
    notifications.value = [];
  }
}

watch(
  () => [auth.isAuthenticated, route.fullPath],
  async () => {
    await loadNotifications();
  },
  { immediate: true },
);

onMounted(async () => {
  await loadNotifications();
});
</script>

<style scoped>
.shell-header {
  height: auto;
  padding: 0;
  background: rgba(255, 255, 255, 0.88);
  backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--line);
  position: sticky;
  top: 0;
  z-index: 10;
}

.shell-inner {
  width: min(1200px, calc(100vw - 32px));
  margin: 0 auto;
  min-height: 76px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  min-width: 220px;
}

.brand-mark {
  width: 42px;
  height: 42px;
  display: grid;
  place-items: center;
  border-radius: 14px;
  background: linear-gradient(135deg, #1677ff, #0f4fae);
  color: #fff;
  font-weight: 800;
  letter-spacing: 1px;
}

.brand-title {
  color: var(--text-main);
  font-size: 18px;
  font-weight: 700;
}

.brand-subtitle {
  color: var(--text-secondary);
  font-size: 12px;
}

.nav-links {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.nav-link {
  color: var(--text-secondary);
  font-weight: 600;
}

.nav-link.router-link-active {
  color: var(--brand-strong);
}

@media (max-width: 980px) {
  .shell-inner {
    width: min(100vw - 20px, 1200px);
    padding: 12px 0;
    align-items: flex-start;
    flex-direction: column;
  }

  .nav-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>

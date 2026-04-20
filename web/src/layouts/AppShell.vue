<template>
  <a-layout style="min-height: 100vh">
    <a-layout-header class="shell-header">
      <div class="shell-inner">
        <RouterLink class="brand" to="/">
          <BrandLogo compact />
        </RouterLink>

        <div class="nav-links">
          <RouterLink v-for="item in primaryNavItems" :key="item.to" class="nav-link" :to="item.to">
            {{ item.label }}
          </RouterLink>
          <a-dropdown v-if="secondaryNavItems.length">
            <a class="nav-link more-link" @click.prevent>
              更多
            </a>
            <template #overlay>
              <a-menu @click="handleMoreMenuClick">
                <a-menu-item v-for="item in secondaryNavItems" :key="item.to">
                  {{ item.label }}
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>

        <a-space :size="12" class="nav-actions">
          <RouterLink to="/search">
            <a-button type="default">全局搜索</a-button>
          </RouterLink>
          <RouterLink v-if="auth.hasBackofficeAccess" :to="adminEntryRoute">
            <a-button type="default">后台</a-button>
          </RouterLink>
          <a-popover v-if="auth.isAuthenticated" trigger="click" placement="bottomRight" overlay-class-name="notification-popover">
            <template #content>
              <div class="notification-preview">
                <div class="notification-preview-head">
                  <div>
                    <div class="notification-preview-title">提醒列表</div>
                    <div class="notification-preview-subtitle">未读 {{ unreadCount }} 条</div>
                  </div>
                  <a-button size="small" :disabled="!unreadCount" @click="markNotificationsRead">全部已读</a-button>
                </div>

                <a-empty v-if="!notificationPreviewItems.length" description="当前没有提醒" />

                <div v-else class="notification-preview-list">
                  <button
                    v-for="item in notificationPreviewItems"
                    :key="item.ID"
                    type="button"
                    :class="['notification-preview-item', { unread: !item.Read, clickable: Boolean(item.Link) }]"
                    @click="openNotificationItem(item)"
                  >
                    <div class="notification-preview-item-head">
                      <div class="notification-preview-item-title">
                        <span v-if="!item.Read" class="notification-preview-dot" />
                        <span>{{ item.Title }}</span>
                      </div>
                      <span class="notification-preview-date">{{ formatDate(item.CreatedAt) }}</span>
                    </div>
                    <div class="notification-preview-content">{{ item.Content }}</div>
                  </button>
                </div>

                <div class="notification-preview-actions">
                  <a-button type="link" @click="openNotificationCenter">查看全部提醒</a-button>
                </div>
              </div>
            </template>
            <a-badge :count="unreadCount" :offset="[6, 2]">
              <a-button>提醒</a-button>
            </a-badge>
          </a-popover>
          <a-dropdown v-if="auth.isAuthenticated">
            <a-button type="primary">
              {{ auth.user?.username ?? '个人中心' }}
            </a-button>
            <template #overlay>
              <a-menu @click="handleUserMenuClick">
                <a-menu-item key="/me">个人中心</a-menu-item>
                <a-menu-item key="settings-profile">个人设置</a-menu-item>
                <a-menu-item key="settings-verification">提交认证</a-menu-item>
                <a-menu-item key="notifications">通知中心</a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout">退出登录</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
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
import { message } from 'ant-design-vue';
import { RouterLink, RouterView, useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import BrandLogo from '@/components/BrandLogo.vue';
import { PERMISSIONS } from '@/constants/permissions';
import { useAuthStore } from '@/stores/auth';
import type { NotificationRecord } from '@/types/api';

const auth = useAuthStore();
const route = useRoute();
const router = useRouter();
const notifications = ref<NotificationRecord[]>([]);

const unreadCount = computed(() => notifications.value.filter((item) => !item.Read).length);
const notificationPreviewItems = computed(() => {
  const unreadItems = notifications.value.filter((item) => !item.Read);
  return (unreadItems.length ? unreadItems : notifications.value).slice(0, 5);
});

const navItems = [
  { label: '首页', to: '/' },
  { label: '模型仓库', to: '/models' },
  { label: '数据集', to: '/datasets' },
  { label: '场景专区', to: '/scenes' },
  { label: '任务模板', to: '/templates' },
  { label: '技能分享', to: '/skills' },
  { label: '互动讨论', to: '/discussions' },
  { label: '私信', to: '/messages' },
  { label: '协作空间', to: '/workspaces' },
  { label: '具身案例', to: '/applications' },
  { label: '文档中心', to: '/docs' },
  { label: '开发者接入', to: '/developers' },
];
const primaryNavItems = computed(() => navItems.slice(0, 8));
const secondaryNavItems = computed(() => navItems.slice(8));
const adminEntryRoute = computed(() => {
  if (auth.hasPermission(PERMISSIONS.portalAccess)) return '/admin/portal';
  if (auth.hasPermission(PERMISSIONS.reviewAccess) || auth.hasPermission(PERMISSIONS.datasetAccessAccess)) return '/admin/reviews';
  return '/admin';
});

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

function handleMoreMenuClick({ key }: { key: string | number }) {
  void router.push(String(key));
}

async function handleUserMenuClick({ key }: { key: string | number }) {
  const action = String(key);
  if (action === 'logout') {
    await auth.logout();
    notifications.value = [];
    await router.push('/');
    return;
  }
  if (action === 'settings-profile') {
    await router.push({ path: '/me', query: { settings: '1', settingsTab: 'profile' } });
    return;
  }
  if (action === 'settings-verification') {
    await router.push({ path: '/me', query: { settings: '1', settingsTab: 'verification' } });
    return;
  }
  if (action === 'notifications') {
    openNotificationCenter();
    return;
  }
  await router.push(action);
}

async function markNotificationsRead() {
  if (!unreadCount.value) return;
  await api.readNotifications();
  await loadNotifications();
  message.success('通知已标记已读');
}

function openNotificationCenter() {
  void router.push({ path: '/me', query: { tab: 'community', subtab: 'notifications' } });
}

function openNotificationItem(item: NotificationRecord) {
  if (item.Link) {
    void router.push(item.Link);
    return;
  }
  if (item.Type === 'message') {
    void router.push('/messages');
    return;
  }
  openNotificationCenter();
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('zh-CN');
}
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
  min-width: 250px;
}

.nav-links {
  display: flex;
  gap: 16px;
  align-items: center;
  flex: 1;
  min-width: 0;
  overflow: hidden;
}

.nav-link {
  color: var(--text-secondary);
  font-weight: 600;
  white-space: nowrap;
  flex: none;
}

.nav-link.router-link-active {
  color: var(--brand-strong);
}

.more-link {
  display: inline-flex;
  align-items: center;
}

.nav-actions {
  flex: none;
}

.notification-preview {
  width: min(360px, calc(100vw - 48px));
}

.notification-preview-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
  margin-bottom: 12px;
}

.notification-preview-title {
  color: var(--text-main);
  font-weight: 700;
}

.notification-preview-subtitle,
.notification-preview-date {
  color: var(--text-secondary);
  font-size: 12px;
}

.notification-preview-list {
  display: grid;
  gap: 10px;
}

.notification-preview-item {
  width: 100%;
  text-align: left;
  border: 1px solid rgba(220, 230, 242, 0.82);
  border-radius: 14px;
  padding: 12px;
  background: #fff;
  cursor: default;
}

.notification-preview-item.unread {
  border-color: rgba(255, 77, 79, 0.28);
  background: linear-gradient(180deg, #fff8f8 0%, #ffffff 100%);
}

.notification-preview-item.clickable {
  cursor: pointer;
}

.notification-preview-item.clickable:hover {
  border-color: rgba(22, 119, 255, 0.28);
  background: linear-gradient(180deg, #f8fbff 0%, #ffffff 100%);
}

.notification-preview-item-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
}

.notification-preview-item-title {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: var(--text-main);
  font-weight: 700;
}

.notification-preview-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #ff4d4f;
  flex: none;
}

.notification-preview-content {
  margin-top: 8px;
  color: var(--text-secondary);
  line-height: 1.7;
}

.notification-preview-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

@media (max-width: 980px) {
  .shell-inner {
    width: min(100vw - 20px, 1200px);
    padding: 12px 0;
    align-items: flex-start;
    flex-direction: column;
  }

  .nav-links {
    flex-wrap: wrap;
    overflow: visible;
  }

  .nav-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>

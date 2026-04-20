<template>
  <div class="profile-card" :class="{ sticky }">
    <div class="profile-cover"></div>
    <div class="profile-body">
      <a-avatar :src="avatar || undefined" :size="112" class="profile-avatar">
        {{ fallbackInitial }}
      </a-avatar>

      <div class="profile-head">
        <h1 class="profile-name">{{ name }}</h1>
        <div class="profile-handle">@{{ handle }}</div>
      </div>

      <a-space v-if="badges.length" wrap class="profile-badges">
        <a-tag v-for="badge in badges" :key="`${badge.color}-${badge.text}`" :color="badge.color">
          {{ badge.text }}
        </a-tag>
      </a-space>

      <p class="profile-bio">{{ bio || '这个开发者暂未填写简介。' }}</p>

      <div v-if="$slots.actions" class="profile-actions">
        <slot name="actions" />
      </div>

      <div v-if="stats.length" class="profile-stats">
        <div v-for="item in stats" :key="item.label" class="profile-stat">
          <div class="profile-stat-value">{{ item.value }}</div>
          <div class="profile-stat-label">{{ item.label }}</div>
        </div>
      </div>

      <div v-if="$slots.default" class="profile-meta">
        <slot />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

type BadgeItem = {
  text: string;
  color: string;
};

type StatItem = {
  label: string;
  value: string | number;
};

const props = withDefaults(
  defineProps<{
    name: string;
    handle: string;
    bio?: string;
    avatar?: string;
    badges?: BadgeItem[];
    stats?: StatItem[];
    sticky?: boolean;
  }>(),
  {
    bio: '',
    avatar: '',
    badges: () => [],
    stats: () => [],
    sticky: false,
  },
);

const fallbackInitial = computed(() => props.name.trim().charAt(0).toUpperCase() || 'U');
</script>

<style scoped>
.profile-card {
  position: relative;
  isolation: isolate;
  overflow: hidden;
  border: 1px solid var(--line);
  border-radius: 24px;
  background: linear-gradient(180deg, #f4f8ff 0, #fff 144px);
  box-shadow: var(--shadow);
}

.profile-card.sticky {
  position: sticky;
  top: 24px;
  z-index: 6;
}

.profile-cover {
  position: relative;
  z-index: 0;
  height: 126px;
  background:
    radial-gradient(circle at top left, rgba(19, 163, 127, 0.2), transparent 34%),
    radial-gradient(circle at top right, rgba(22, 119, 255, 0.28), transparent 36%),
    linear-gradient(135deg, #dbeafe 0%, #ecfdf5 48%, #f8fafc 100%);
  border-bottom: 1px solid rgba(220, 230, 242, 0.85);
}

.profile-body {
  position: relative;
  z-index: 1;
  padding: 0 24px 24px;
}

.profile-avatar {
  margin-top: -56px;
  border: 4px solid #fff;
  box-shadow: 0 12px 28px rgba(22, 50, 79, 0.18);
  background: linear-gradient(135deg, #1677ff, #13a37f);
  color: #fff;
  font-size: 40px;
  font-weight: 700;
}

.profile-head {
  margin-top: 16px;
}

.profile-name {
  margin: 0;
  color: var(--text-main);
  font-size: 28px;
  line-height: 1.1;
}

.profile-handle {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 15px;
}

.profile-badges {
  margin-top: 16px;
}

.profile-bio {
  margin: 18px 0 0;
  color: var(--text-secondary);
  line-height: 1.75;
  white-space: pre-wrap;
}

.profile-actions {
  margin-top: 20px;
}

.profile-actions :deep(.ant-space) {
  width: 100%;
}

.profile-actions :deep(.ant-btn) {
  width: 100%;
}

.profile-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
  margin-top: 20px;
}

.profile-stat {
  padding: 14px 16px;
  border-radius: 16px;
  background: var(--surface-soft);
  border: 1px solid rgba(220, 230, 242, 0.82);
}

.profile-stat-value {
  color: var(--text-main);
  font-size: 20px;
  font-weight: 700;
}

.profile-stat-label {
  margin-top: 4px;
  color: var(--text-secondary);
  font-size: 12px;
}

.profile-meta {
  margin-top: 18px;
  padding-top: 18px;
  border-top: 1px solid rgba(220, 230, 242, 0.85);
  color: var(--text-secondary);
  font-size: 13px;
  line-height: 1.7;
}

@media (max-width: 960px) {
  .profile-card.sticky {
    position: relative;
    top: 0;
  }
}
</style>

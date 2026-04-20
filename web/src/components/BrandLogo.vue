<template>
  <div class="brand-logo" :class="{ compact }">
    <svg class="brand-icon" viewBox="0 0 72 72" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
      <defs>
        <linearGradient id="brand-logo-gradient" x1="8" y1="8" x2="64" y2="64" gradientUnits="userSpaceOnUse">
          <stop offset="0" stop-color="#1677ff" />
          <stop offset="0.55" stop-color="#0f4fae" />
          <stop offset="1" stop-color="#13a37f" />
        </linearGradient>
      </defs>
      <rect x="6" y="6" width="60" height="60" rx="20" fill="url(#brand-logo-gradient)" />
      <path d="M22 22H40" stroke="#fff" stroke-width="5" stroke-linecap="round" />
      <path d="M22 36H36" stroke="#fff" stroke-width="5" stroke-linecap="round" />
      <path d="M22 50H40" stroke="#fff" stroke-width="5" stroke-linecap="round" />
      <path d="M45 24V48" stroke="#d9fff5" stroke-width="3.5" stroke-linecap="round" />
      <path d="M45 24H52" stroke="#d9fff5" stroke-width="3.5" stroke-linecap="round" />
      <path d="M45 36H56" stroke="#d9fff5" stroke-width="3.5" stroke-linecap="round" />
      <path d="M45 48H52" stroke="#d9fff5" stroke-width="3.5" stroke-linecap="round" />
      <circle cx="56" cy="24" r="4" fill="#d9fff5" />
      <circle cx="60" cy="36" r="4" fill="#9ef3dd" />
      <circle cx="56" cy="48" r="4" fill="#d9fff5" />
    </svg>

    <div v-if="!iconOnly" class="brand-copy">
      <div class="brand-wordmark">EIBotHub</div>
      <div v-if="showSubtitle" class="brand-caption">Embodied Intelligence Bot Hub</div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    compact?: boolean;
    iconOnly?: boolean;
    showSubtitle?: boolean;
  }>(),
  {
    compact: false,
    iconOnly: false,
    showSubtitle: true,
  },
);
</script>

<style scoped>
.brand-logo {
  display: inline-flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.brand-icon {
  width: 48px;
  height: 48px;
  flex: none;
  filter: drop-shadow(0 10px 20px rgba(22, 50, 79, 0.18));
}

.brand-copy {
  min-width: 0;
}

.brand-wordmark {
  color: var(--text-main);
  font-size: 20px;
  font-weight: 800;
  letter-spacing: 0.01em;
  line-height: 1;
}

.brand-caption {
  margin-top: 4px;
  color: var(--text-secondary);
  font-size: 11px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  white-space: nowrap;
}

.brand-logo.compact {
  gap: 10px;
}

.brand-logo.compact .brand-icon {
  width: 42px;
  height: 42px;
}

.brand-logo.compact .brand-copy {
  transform: translateY(12px);
}

.brand-logo.compact .brand-wordmark {
  font-size: 18px;
  transform: translateY(18px);
}

.brand-logo.compact .brand-caption {
  margin-top: 4px;
}

@media (max-width: 768px) {
  .brand-caption {
    display: none;
  }
}
</style>

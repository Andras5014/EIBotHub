<template>
  <AppShell />
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth';

import AppShell from '@/layouts/AppShell.vue';

const auth = useAuthStore();

if (auth.token) {
  void auth.refreshMe().catch(async () => {
    await auth.logout();
  });
}
</script>

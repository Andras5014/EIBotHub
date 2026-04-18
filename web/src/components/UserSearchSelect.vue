<template>
  <a-select
    :value="modelValue"
    show-search
    allow-clear
    :filter-option="false"
    :options="options"
    :placeholder="placeholder"
    style="width: 100%"
    @search="handleSearch"
    @change="handleChange"
  />
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

import { api } from '@/api';

const props = defineProps<{
  modelValue?: number;
  placeholder?: string;
}>();

const emit = defineEmits<{
  'update:modelValue': [value?: number];
}>();

const options = ref<Array<{ label: string; value: number }>>([]);

watch(
  () => props.modelValue,
  async (value) => {
    if (!value) {
      return;
    }
    const exists = options.value.some((item) => item.value === value);
    if (exists) {
      return;
    }
    try {
      const user = await api.getPublicUser(value);
      options.value = [{ label: user.username, value: user.id }, ...options.value];
    } catch {
      // Ignore missing users here; the calling form will still validate on submit.
    }
  },
  { immediate: true },
);

async function handleSearch(keyword: string) {
  if (!keyword.trim()) {
    options.value = [];
    return;
  }
  const result = await api.search({ q: keyword, type: 'user' });
  options.value = result.items.map((item) => ({
    label: `${item.title}${item.summary ? ` · ${item.summary}` : ''}`,
    value: item.id,
  }));
}

function handleChange(value?: number) {
  emit('update:modelValue', value);
}
</script>

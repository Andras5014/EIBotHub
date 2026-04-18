<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">协作空间</h1>
          <p class="section-subtitle">围绕模板、模型和数据集联调的协作空间，支持成员协作和消息同步。</p>
        </div>
      </div>

      <a-row :gutter="[16, 16]">
        <a-col :xs="24" :lg="8">
          <a-card title="创建协作空间" class="inner-card">
            <a-form :model="createForm" layout="vertical" @finish="createWorkspace">
              <a-form-item label="名称"><a-input v-model:value="createForm.name" /></a-form-item>
              <a-form-item label="摘要"><a-textarea v-model:value="createForm.summary" :rows="3" /></a-form-item>
              <a-form-item><a-button type="primary" html-type="submit">创建</a-button></a-form-item>
            </a-form>
          </a-card>
          <a-card title="空间列表" class="inner-card" style="margin-top: 16px">
            <a-list :data-source="spaces">
              <template #renderItem="{ item }">
                <a-list-item :class="{ active: selected?.id === item.id }" @click="selectWorkspace(item.id)" style="cursor: pointer">
                  <a-list-item-meta :title="item.name" :description="`${item.summary} · 成员 ${item.member_count}`" />
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-col>

        <a-col :xs="24" :lg="16">
          <a-card title="空间详情" class="inner-card">
            <template v-if="selected">
              <h2>{{ selected.name }}</h2>
              <p class="section-subtitle">{{ selected.summary }}</p>
              <a-space wrap style="margin-bottom: 12px">
                <span class="pill-meta">成员 {{ selected.member_count }}</span>
                <span class="pill-meta">空间 ID {{ selected.id }}</span>
              </a-space>
              <a-row :gutter="[16, 16]" style="margin-top: 12px">
                <a-col :xs="24" :lg="10">
                  <a-card title="成员" size="small">
                    <a-list :data-source="selected.members">
                      <template #renderItem="{ item }">
                        <a-list-item>{{ item.user_name }}</a-list-item>
                      </template>
                    </a-list>
                    <a-divider />
                    <a-form :model="memberForm" layout="vertical" @finish="addMember">
                      <a-form-item label="新增成员">
                        <UserSearchSelect v-model="memberForm.user_id" placeholder="搜索用户名" />
                      </a-form-item>
                      <a-form-item><a-button html-type="submit">添加成员</a-button></a-form-item>
                    </a-form>
                  </a-card>
                </a-col>
                <a-col :xs="24" :lg="14">
                  <a-card title="空间消息" size="small">
                    <a-list :data-source="selected.messages">
                      <template #renderItem="{ item }">
                        <a-list-item>
                          <a-list-item-meta :title="item.sender_name" :description="`${item.content} · ${formatDate(item.created_at)}`" />
                        </a-list-item>
                      </template>
                    </a-list>
                    <a-divider />
                    <a-form :model="messageForm" layout="vertical" @finish="sendWorkspaceMessage">
                      <a-form-item label="消息内容"><a-textarea v-model:value="messageForm.content" :rows="3" /></a-form-item>
                      <a-form-item><a-button type="primary" html-type="submit">发送消息</a-button></a-form-item>
                    </a-form>
                  </a-card>
                </a-col>
              </a-row>
            </template>
            <a-empty v-else description="请选择一个协作空间" />
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue';
import { message } from 'ant-design-vue';
import { useRoute, useRouter } from 'vue-router';

import { api } from '@/api';
import UserSearchSelect from '@/components/UserSearchSelect.vue';
import type { WorkspaceDetail, WorkspaceItem } from '@/types/api';

const route = useRoute();
const router = useRouter();
const spaces = ref<WorkspaceItem[]>([]);
const selected = ref<WorkspaceDetail>();
const createForm = reactive({
  name: '',
  summary: '',
});
const memberForm = reactive({
  user_id: undefined as number | undefined,
});
const messageForm = reactive({
  content: '',
});

async function loadSpaces() {
  spaces.value = await api.listWorkspaces();
}

async function selectWorkspace(id: number) {
  selected.value = await api.getWorkspace(id);
  router.replace({ query: { ...route.query, id: String(id) } });
}

async function createWorkspace() {
  try {
    const workspace = await api.createWorkspace(createForm);
    createForm.name = '';
    createForm.summary = '';
    await loadSpaces();
    selected.value = workspace;
    message.success('协作空间已创建');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '创建失败');
  }
}

async function addMember() {
  if (!selected.value) return;
  if (!memberForm.user_id) {
    message.error('请选择要添加的成员');
    return;
  }
  try {
    await api.addWorkspaceMember(selected.value.id, { user_id: memberForm.user_id });
    selected.value = await api.getWorkspace(selected.value.id);
    memberForm.user_id = undefined;
    message.success('成员已添加');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '添加成员失败');
  }
}

async function sendWorkspaceMessage() {
  if (!selected.value) return;
  try {
    await api.sendWorkspaceMessage(selected.value.id, messageForm);
    messageForm.content = '';
    selected.value = await api.getWorkspace(selected.value.id);
    await loadSpaces();
    message.success('消息已发送');
  } catch (error) {
    message.error(error instanceof Error ? error.message : '发送失败');
  }
}

watch(
  () => [route.query.name, route.query.summary],
  ([name, summary]) => {
    if (typeof name === 'string') {
      createForm.name = name;
    }
    if (typeof summary === 'string') {
      createForm.summary = summary;
    }
  },
  { immediate: true },
);

watch(
  () => route.query.id,
  async (id) => {
    if (id) {
      await selectWorkspace(Number(id));
    }
  },
);

onMounted(async () => {
  await loadSpaces();
  const workspaceID = route.query.id;
  if (workspaceID) {
    await selectWorkspace(Number(workspaceID));
    return;
  }
  if (spaces.value.length > 0) {
    await selectWorkspace(spaces.value[0].id);
  }
});

function formatDate(value: string) {
  return new Date(value).toLocaleString();
}
</script>

<style scoped>
.block {
  padding: 24px;
}

.inner-card {
  border-radius: 18px;
}

.active {
  background: var(--surface-soft);
  border-radius: 14px;
}
</style>

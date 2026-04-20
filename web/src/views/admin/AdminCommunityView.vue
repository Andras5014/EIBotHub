<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">社区治理</h1>
          <p class="section-subtitle">统一治理技能、讨论、评论、私信会话与协作空间，补齐社区风险处理闭环。</p>
        </div>
      </div>

      <div class="grid-cards" v-if="overview">
        <a-card><a-statistic title="技能数" :value="overview.skills" /></a-card>
        <a-card><a-statistic title="讨论数" :value="overview.discussions" /></a-card>
        <a-card><a-statistic title="评论数" :value="overview.comments" /></a-card>
        <a-card><a-statistic title="私信会话" :value="conversations.length" /></a-card>
        <a-card><a-statistic title="协作空间" :value="workspaces.length" /></a-card>
      </div>

      <a-row :gutter="[16, 16]" style="margin-top: 18px">
        <a-col :xs="24" :lg="10">
          <a-input v-model:value="keyword" placeholder="按标题、摘要、用户名、消息内容或空间信息筛选" />
        </a-col>
        <a-col :xs="24" :lg="6">
          <a-select v-model:value="skillStatusFilter" style="width: 100%" :options="skillStatusOptions" />
        </a-col>
        <a-col :xs="24" :lg="8">
          <a-select v-model:value="commentTypeFilter" style="width: 100%" :options="commentTypeOptions" />
        </a-col>
      </a-row>

      <a-tabs style="margin-top: 18px">
        <a-tab-pane v-if="canModerateContent" key="skills" tab="技能">
          <a-list :data-source="filteredSkills">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-space wrap>
                      <span>{{ item.name }}</span>
                      <a-tag :color="item.status === 'published' ? 'green' : 'red'">{{ item.status }}</a-tag>
                    </a-space>
                  </template>
                  <template #description>
                    <div>{{ item.summary }}</div>
                    <a-space wrap style="margin-top: 6px">
                      <RouterLink :to="`/community/users/${item.owner_id}`">查看作者</RouterLink>
                      <span>作者：{{ item.owner_name }}</span>
                      <span>更新时间：{{ formatDate(item.updated_at) }}</span>
                    </a-space>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button danger type="link" @click="hideSkill(item.id)">下架</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
        <a-tab-pane v-if="canModerateContent" key="discussions" tab="讨论">
          <a-list :data-source="filteredDiscussions">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <RouterLink :to="`/discussions/${item.id}`">{{ item.title }}</RouterLink>
                  </template>
                  <template #description>
                    <a-space wrap>
                      <span class="pill-meta">{{ item.tag }}</span>
                      <RouterLink :to="`/community/users/${item.user_id}`">查看作者</RouterLink>
                      <span>作者：{{ item.user_name }}</span>
                      <span>评论：{{ item.comment_count }}</span>
                      <span>更新时间：{{ formatDate(item.updated_at) }}</span>
                    </a-space>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button danger type="link" @click="removeDiscussion(item.id)">删除</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
        <a-tab-pane v-if="canModerateContent" key="comments" tab="评论">
          <a-list :data-source="filteredComments">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-space wrap>
                      <RouterLink :to="resourceRoute(item.resource_type, item.resource_id)">
                        {{ resourceLabel(item.resource_type) }} #{{ item.resource_id }}
                      </RouterLink>
                      <RouterLink :to="`/community/users/${item.user_id}`">查看作者</RouterLink>
                      <span>作者：{{ item.user_name }}</span>
                    </a-space>
                  </template>
                  <template #description>
                    <div>{{ item.content }}</div>
                    <a-space wrap style="margin-top: 6px">
                      <span>创建时间：{{ formatDate(item.created_at) }}</span>
                    </a-space>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-button danger type="link" @click="removeComment(item.id)">删除</a-button>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
        <a-tab-pane v-if="canModerateConversations" key="conversations" tab="私信">
          <a-list :data-source="filteredConversations">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-space wrap>
                      <span>{{ item.title || '未命名会话' }}</span>
                      <a-tag :color="item.active ? 'green' : 'red'">{{ item.active ? '正常' : '已封禁' }}</a-tag>
                      <span class="pill-meta">{{ conversationKindLabel(item.kind) }}</span>
                    </a-space>
                  </template>
                  <template #description>
                    <div>{{ item.latest_message || '暂无消息' }}</div>
                    <div v-if="item.blocked_reason" class="muted-danger">封禁原因：{{ item.blocked_reason }}</div>
                    <a-space wrap style="margin-top: 6px">
                      <a-tag v-for="name in item.participant_names" :key="name">{{ name }}</a-tag>
                      <span>更新时间：{{ formatDate(item.updated_at) }}</span>
                    </a-space>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-space class="govern-action" wrap>
                    <a-input
                      v-if="item.active"
                      v-model:value="conversationReasonForms[item.id]"
                      placeholder="输入封禁原因"
                      style="width: 220px"
                    />
                    <a-button v-if="item.active" danger @click="blockConversation(item)">封禁</a-button>
                    <a-button v-else @click="unblockConversation(item.id)">解封</a-button>
                  </a-space>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
        <a-tab-pane v-if="canModerateWorkspaces" key="workspaces" tab="协作空间">
          <a-list :data-source="filteredWorkspaces">
            <template #renderItem="{ item }">
              <a-list-item>
                <a-list-item-meta>
                  <template #title>
                    <a-space wrap>
                      <span>{{ item.name }}</span>
                      <a-tag :color="item.active ? 'green' : 'red'">{{ item.active ? '正常' : '已封禁' }}</a-tag>
                      <span class="pill-meta">所有者：{{ item.owner_name }}</span>
                    </a-space>
                  </template>
                  <template #description>
                    <div>{{ item.summary }}</div>
                    <div v-if="item.blocked_reason" class="muted-danger">封禁原因：{{ item.blocked_reason }}</div>
                    <a-space wrap style="margin-top: 6px">
                      <span>成员数：{{ item.member_count }}</span>
                      <span>更新时间：{{ formatDate(item.updated_at) }}</span>
                    </a-space>
                    <div class="member-wrap" v-if="item.members.length">
                      <a-space wrap>
                        <span v-for="member in item.members" :key="member.user_id" class="member-chip">
                          <a-tag :color="member.user_id === item.owner_id ? 'geekblue' : 'default'">
                            {{ member.user_name }} #{{ member.user_id }}
                          </a-tag>
                          <a-button
                            v-if="member.user_id !== item.owner_id"
                            type="link"
                            danger
                            size="small"
                            @click="removeWorkspaceMember(item.id, member.user_id)"
                          >
                            移除
                          </a-button>
                        </span>
                      </a-space>
                    </div>
                  </template>
                </a-list-item-meta>
                <template #actions>
                  <a-space class="govern-action" wrap>
                    <a-input
                      v-if="item.active"
                      v-model:value="workspaceReasonForms[item.id]"
                      placeholder="输入封禁原因"
                      style="width: 220px"
                    />
                    <a-button v-if="item.active" danger @click="blockWorkspace(item)">封禁</a-button>
                    <a-button v-else @click="unblockWorkspace(item.id)">解封</a-button>
                  </a-space>
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import { PERMISSIONS } from '@/constants/permissions';
import { useAuthStore } from '@/stores/auth';
import type {
  AdminCommentModerationItem,
  AdminCommunityOverview,
  AdminConversationModerationItem,
  AdminDiscussionModerationItem,
  AdminSkillModerationItem,
  AdminWorkspaceModerationItem,
} from '@/types/api';

const auth = useAuthStore();
const overview = ref<AdminCommunityOverview>();
const skills = ref<AdminSkillModerationItem[]>([]);
const discussions = ref<AdminDiscussionModerationItem[]>([]);
const comments = ref<AdminCommentModerationItem[]>([]);
const conversations = ref<AdminConversationModerationItem[]>([]);
const workspaces = ref<AdminWorkspaceModerationItem[]>([]);
const keyword = ref('');
const skillStatusFilter = ref<'all' | string>('all');
const commentTypeFilter = ref<'all' | string>('all');
const conversationReasonForms = reactive<Record<number, string>>({});
const workspaceReasonForms = reactive<Record<number, string>>({});
const canModerateContent = computed(() => auth.hasPermission(PERMISSIONS.communityContentModerate));
const canModerateConversations = computed(() => auth.hasPermission(PERMISSIONS.conversationModerate));
const canModerateWorkspaces = computed(() => auth.hasPermission(PERMISSIONS.workspaceModerate));

const skillStatusOptions = [
  { label: '全部技能状态', value: 'all' },
  { label: '已发布', value: 'published' },
  { label: '已驳回', value: 'rejected' },
];

const commentTypeOptions = [
  { label: '全部评论类型', value: 'all' },
  { label: '模型', value: 'model' },
  { label: '数据集', value: 'dataset' },
  { label: '模板', value: 'task-template' },
  { label: '技能', value: 'skill' },
  { label: '讨论', value: 'discussion' },
];

const filteredSkills = computed(() =>
  skills.value.filter((item) => {
    const matchKeyword =
      !keyword.value.trim() ||
      `${item.name} ${item.summary} ${item.owner_name}`.includes(keyword.value.trim());
    const matchStatus = skillStatusFilter.value === 'all' || item.status === skillStatusFilter.value;
    return matchKeyword && matchStatus;
  }),
);

const filteredDiscussions = computed(() =>
  discussions.value.filter((item) =>
    !keyword.value.trim() || `${item.title} ${item.tag} ${item.user_name}`.includes(keyword.value.trim()),
  ),
);

const filteredComments = computed(() =>
  comments.value.filter((item) => {
    const matchKeyword =
      !keyword.value.trim() ||
      `${item.user_name} ${item.content} ${item.resource_type}`.includes(keyword.value.trim());
    const matchType = commentTypeFilter.value === 'all' || item.resource_type === commentTypeFilter.value;
    return matchKeyword && matchType;
  }),
);

const filteredConversations = computed(() =>
  conversations.value.filter((item) =>
    !keyword.value.trim() ||
    `${item.title} ${item.latest_message} ${item.participant_names.join(' ')} ${item.blocked_reason}`.includes(
      keyword.value.trim(),
    ),
  ),
);

const filteredWorkspaces = computed(() =>
  workspaces.value.filter((item) =>
    !keyword.value.trim() ||
    `${item.name} ${item.summary} ${item.owner_name} ${item.blocked_reason} ${item.members
      .map((member) => member.user_name)
      .join(' ')}`.includes(keyword.value.trim()),
  ),
);

async function load() {
  const [communityOverview, skillItems, discussionItems, commentItems, conversationItems, workspaceItems] =
    await Promise.all([
    api.getAdminCommunityOverview(),
    canModerateContent.value ? api.getAdminCommunitySkills() : Promise.resolve([]),
    canModerateContent.value ? api.getAdminCommunityDiscussions() : Promise.resolve([]),
    canModerateContent.value ? api.getAdminCommunityComments() : Promise.resolve([]),
    canModerateConversations.value ? api.getAdminCommunityConversations() : Promise.resolve([]),
    canModerateWorkspaces.value ? api.getAdminCommunityWorkspaces() : Promise.resolve([]),
  ]);
  overview.value = communityOverview;
  skills.value = skillItems;
  discussions.value = discussionItems;
  comments.value = commentItems;
  conversations.value = conversationItems;
  workspaces.value = workspaceItems;
}

async function hideSkill(id: number) {
  if (!(await confirmDanger('确认下架技能？', '下架后该技能会退出公开展示。'))) return;
  await api.hideAdminSkill(id);
  message.success('技能已下架');
  await load();
}

async function removeDiscussion(id: number) {
  if (!(await confirmDanger('确认删除讨论？', '删除后该讨论将无法恢复。'))) return;
  await api.removeAdminDiscussion(id);
  message.success('讨论已删除');
  await load();
}

async function removeComment(id: number) {
  if (!(await confirmDanger('确认删除评论？', '删除后该评论将无法恢复。'))) return;
  await api.removeAdminComment(id);
  message.success('评论已删除');
  await load();
}

async function blockConversation(item: AdminConversationModerationItem) {
  const reason = conversationReasonForms[item.id]?.trim();
  if (!reason) {
    message.error('请输入封禁原因');
    return;
  }
  if (!(await confirmDanger('确认封禁私信会话？', '封禁后当前会话成员将无法继续查看和发送消息。'))) return;
  await api.blockAdminConversation(item.id, { reason });
  conversationReasonForms[item.id] = '';
  message.success('私信会话已封禁');
  await load();
}

async function unblockConversation(id: number) {
  if (!(await confirmDanger('确认解封私信会话？', '解封后会话成员将恢复访问和发送能力。'))) return;
  await api.unblockAdminConversation(id);
  message.success('私信会话已解封');
  await load();
}

async function blockWorkspace(item: AdminWorkspaceModerationItem) {
  const reason = workspaceReasonForms[item.id]?.trim();
  if (!reason) {
    message.error('请输入封禁原因');
    return;
  }
  if (!(await confirmDanger('确认封禁协作空间？', '封禁后成员将无法继续访问空间和发送消息。'))) return;
  await api.blockAdminWorkspace(item.id, { reason });
  workspaceReasonForms[item.id] = '';
  message.success('协作空间已封禁');
  await load();
}

async function unblockWorkspace(id: number) {
  if (!(await confirmDanger('确认解封协作空间？', '解封后成员将恢复访问和协作。'))) return;
  await api.unblockAdminWorkspace(id);
  message.success('协作空间已解封');
  await load();
}

async function removeWorkspaceMember(workspaceID: number, userID: number) {
  if (!(await confirmDanger('确认移除成员？', '移除后该成员将失去当前协作空间访问权限。'))) return;
  await api.removeAdminWorkspaceMember(workspaceID, { user_id: userID });
  message.success('成员已移除');
  await load();
}

function confirmDanger(title: string, content: string) {
  return new Promise<boolean>((resolve) => {
    Modal.confirm({
      title,
      content,
      okType: 'danger',
      okText: '确认',
      cancelText: '取消',
      onOk: () => resolve(true),
      onCancel: () => resolve(false),
    });
  });
}

function resourceRoute(resourceType: string, resourceID: number) {
  return {
    model: `/models/${resourceID}`,
    dataset: `/datasets/${resourceID}`,
    'task-template': `/templates/${resourceID}`,
    skill: `/skills/${resourceID}`,
    discussion: `/discussions/${resourceID}`,
  }[resourceType] ?? '/admin/community';
}

function resourceLabel(resourceType: string) {
  return {
    model: '模型',
    dataset: '数据集',
    'task-template': '模板',
    skill: '技能',
    discussion: '讨论',
  }[resourceType] ?? resourceType;
}

function conversationKindLabel(kind: string) {
  return {
    direct: '私信',
    workspace: '空间会话',
  }[kind] ?? kind;
}

function formatDate(value: string) {
  return new Date(value).toLocaleString();
}

onMounted(load);
</script>

<style scoped>
.block {
  padding: 24px;
}

.govern-action {
  justify-content: flex-end;
}

.muted-danger {
  margin-top: 6px;
  color: #cf1322;
}

.member-wrap {
  margin-top: 10px;
}

.member-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
</style>

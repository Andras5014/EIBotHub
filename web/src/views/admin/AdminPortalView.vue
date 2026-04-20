<template>
  <div class="page-shell">
    <section class="page-card hero-panel">
      <div class="hero-copy">
        <div class="hero-kicker">Portal Operation</div>
        <h1 class="section-title hero-title">门户运营</h1>
        <p class="section-subtitle hero-subtitle">
          这个页面不再只是一堆配置表单，而是把“控制什么、显示在哪里、应该怎么改”先说明清楚，再按首页结构分区维护。
        </p>
      </div>
      <div class="guide-grid">
        <div v-for="item in portalGuideCards" :key="item.key" class="guide-card">
          <div class="guide-label">{{ item.label }}</div>
          <div class="guide-title">{{ item.title }}</div>
          <div class="guide-description">{{ item.description }}</div>
          <div class="guide-location">前台位置：{{ item.location }}</div>
        </div>
      </div>
    </section>

    <section class="page-card block">
      <a-tabs v-model:activeKey="activeTab" class="portal-tabs">
        <a-tab-pane key="structure" tab="首页结构">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="11">
              <a-card v-if="canManagePortal" class="inner-card portal-card">
                <div class="card-head">
                  <div class="card-kicker">首页结构</div>
                  <h3 class="card-title">首页模块开关</h3>
                  <p class="card-description">控制首页各模块是否展示，以及从上到下的出现顺序。</p>
                </div>
                <div class="module-stack">
                  <div v-for="item in modules" :key="item.module_key" class="module-entry">
                    <div class="module-entry-copy">
                      <div class="module-entry-title">{{ item.label }}</div>
                      <div class="module-entry-location">前台位置：{{ moduleHelp(item.module_key).location }}</div>
                      <div class="module-entry-description">{{ moduleHelp(item.module_key).description }}</div>
                    </div>
                    <div class="module-entry-controls">
                      <div class="module-control">
                        <div class="field-label">顺序</div>
                        <a-input-number v-model:value="item.sort_order" :min="0" :max="1000" style="width: 100%" />
                      </div>
                      <div class="module-control module-control-switch">
                        <div class="field-label">启用</div>
                        <a-switch v-model:checked="item.enabled" />
                      </div>
                      <div class="module-control module-control-action">
                        <div class="field-label">保存</div>
                        <a-button type="primary" block @click="saveModule(item)">保存</a-button>
                      </div>
                    </div>
                  </div>
                </div>
              </a-card>
            </a-col>

            <a-col :xs="24" :lg="13">
              <a-card v-if="canManagePortal" class="inner-card portal-card">
                <div class="card-head">
                  <div class="card-kicker">首页首屏</div>
                  <h3 class="card-title">首页亮点文案</h3>
                  <p class="card-description">显示在首页首屏右侧“平台亮点”列表，适合写平台价值、能力范围和使用收益。</p>
                </div>

                <a-form :model="highlightForm" layout="vertical" @finish="saveHighlight">
                  <a-form-item label="亮点文案">
                    <a-input v-model:value="highlightForm.text" placeholder="例如：统一管理模型、数据集和任务模板" />
                  </a-form-item>
                  <a-row :gutter="[12, 0]">
                    <a-col :xs="24" :md="12">
                      <a-form-item label="排序值">
                        <a-input-number v-model:value="highlightForm.sort_order" :min="0" :max="1000" style="width: 100%" />
                      </a-form-item>
                    </a-col>
                    <a-col :xs="24" :md="12">
                      <a-form-item label="启用">
                        <a-switch v-model:checked="highlightForm.enabled" />
                      </a-form-item>
                    </a-col>
                  </a-row>
                  <a-space>
                    <a-button type="primary" html-type="submit">{{ highlightEditingId ? '保存' : '新增' }}</a-button>
                    <a-button @click="resetHighlightForm">重置</a-button>
                  </a-space>
                </a-form>

                <a-divider />
                <a-list :data-source="highlights" size="small">
                  <template #renderItem="{ item }">
                    <a-list-item>
                      <a-list-item-meta :title="item.text" :description="`排序 ${item.sort_order}`" />
                      <a-space>
                        <a-tag :color="item.enabled ? 'green' : 'default'">{{ item.enabled ? '启用' : '停用' }}</a-tag>
                        <a-button type="link" @click="startEditHighlight(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteHighlight(item.id)">删除</a-button>
                      </a-space>
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="hero" tab="导读区">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="15">
              <a-card v-if="canManagePortal" class="inner-card portal-card">
                <div class="card-head">
                  <div class="card-kicker">首页最上方</div>
                  <h3 class="card-title">首页导读区</h3>
                  <p class="card-description">控制首页最顶部蓝色标签、主标题、说明文案，以及三个主操作按钮的文字。</p>
                </div>
                <a-form :model="heroConfig" layout="vertical" @finish="saveHeroConfig">
                  <a-form-item label="导读标签">
                    <a-input v-model:value="heroConfig.tagline" />
                  </a-form-item>
                  <a-form-item label="主标题">
                    <a-input v-model:value="heroConfig.title" />
                  </a-form-item>
                  <a-form-item label="导读说明">
                    <a-textarea v-model:value="heroConfig.description" :rows="4" />
                  </a-form-item>
                  <a-row :gutter="[12, 0]">
                    <a-col :xs="24" :md="8">
                      <a-form-item label="主按钮">
                        <a-input v-model:value="heroConfig.primary_button" />
                      </a-form-item>
                    </a-col>
                    <a-col :xs="24" :md="8">
                      <a-form-item label="次按钮">
                        <a-input v-model:value="heroConfig.secondary_button" />
                      </a-form-item>
                    </a-col>
                    <a-col :xs="24" :md="8">
                      <a-form-item label="搜索按钮">
                        <a-input v-model:value="heroConfig.search_button" />
                      </a-form-item>
                    </a-col>
                  </a-row>
                  <a-button type="primary" html-type="submit">保存导读区</a-button>
                </a-form>
              </a-card>
            </a-col>

            <a-col :xs="24" :lg="9">
              <a-card class="inner-card preview-card">
                <div class="card-head">
                  <div class="card-kicker">前台预览</div>
                  <h3 class="card-title">首页导读区会这样显示</h3>
                  <p class="card-description">这里不是最终像素级预览，而是帮助你确认当前文案分别会出现在什么位置。</p>
                </div>
                <div class="hero-preview">
                  <a-tag color="blue">{{ heroConfig.tagline }}</a-tag>
                  <div class="hero-preview-title">{{ heroConfig.title }}</div>
                  <div class="hero-preview-description">{{ heroConfig.description }}</div>
                  <a-space wrap>
                    <a-button type="primary">{{ heroConfig.primary_button }}</a-button>
                    <a-button>{{ heroConfig.secondary_button }}</a-button>
                    <a-button type="link">{{ heroConfig.search_button }}</a-button>
                  </a-space>
                </div>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="scenes" tab="场景与榜单">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="16">
              <a-card v-if="canManagePortal" class="inner-card portal-card">
                <div class="card-head">
                  <div class="card-kicker">首页应用场景</div>
                  <h3 class="card-title">场景页配置</h3>
                  <p class="card-description">控制首页“应用场景”卡片区，以及 `/scenes` 场景列表页的展示内容。</p>
                </div>
                <a-form :model="sceneForm" layout="vertical" @finish="saveScenePage">
                  <a-row :gutter="[12, 0]">
                    <a-col :xs="24" :md="12">
                      <a-form-item label="场景标识">
                        <a-input v-model:value="sceneForm.slug" placeholder="例如：warehouse-ops" />
                      </a-form-item>
                    </a-col>
                    <a-col :xs="24" :md="12">
                      <a-form-item label="场景名称">
                        <a-input v-model:value="sceneForm.name" placeholder="例如：仓储搬运" />
                      </a-form-item>
                    </a-col>
                  </a-row>
                  <a-form-item label="场景标签">
                    <a-input v-model:value="sceneForm.tagline" placeholder="例如：Warehouse Operation" />
                  </a-form-item>
                  <a-form-item label="场景摘要">
                    <a-input v-model:value="sceneForm.summary" placeholder="一句话描述场景价值" />
                  </a-form-item>
                  <a-form-item label="场景说明">
                    <a-textarea v-model:value="sceneForm.description" :rows="4" />
                  </a-form-item>
                  <a-row :gutter="[12, 0]">
                    <a-col :xs="24" :md="12">
                      <a-form-item label="排序值">
                        <a-input-number v-model:value="sceneForm.sort_order" :min="0" :max="1000" style="width: 100%" />
                      </a-form-item>
                    </a-col>
                    <a-col :xs="24" :md="12">
                      <a-form-item label="启用">
                        <a-switch v-model:checked="sceneForm.enabled" />
                      </a-form-item>
                    </a-col>
                  </a-row>
                  <a-space>
                    <a-button type="primary" html-type="submit">{{ sceneEditingId ? '保存' : '新增' }}</a-button>
                    <a-button @click="resetSceneForm">重置</a-button>
                  </a-space>
                </a-form>

                <a-divider />
                <a-list :data-source="scenePages" size="small">
                  <template #renderItem="{ item }">
                    <a-list-item>
                      <a-list-item-meta :title="item.name" :description="`${item.slug} · 排序 ${item.sort_order}`" />
                      <a-space>
                        <a-tag :color="item.enabled ? 'green' : 'default'">{{ item.enabled ? '启用' : '停用' }}</a-tag>
                        <RouterLink :to="`/scenes/${item.slug}`">预览</RouterLink>
                        <a-button type="link" @click="startEditScenePage(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteScenePage(item.id)">删除</a-button>
                      </a-space>
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>

            <a-col :xs="24" :lg="8">
              <a-card v-if="canManagePortal" class="inner-card portal-card">
                <div class="card-head">
                  <div class="card-kicker">首页社区区</div>
                  <h3 class="card-title">榜单配置</h3>
                  <p class="card-description">控制首页“贡献排行榜”的标题、副标题、展示数量和启用状态。</p>
                </div>
                <a-form :model="rankingConfig" layout="vertical" @finish="saveRankingConfig">
                  <a-form-item label="榜单标题">
                    <a-input v-model:value="rankingConfig.title" />
                  </a-form-item>
                  <a-form-item label="榜单副标题">
                    <a-input v-model:value="rankingConfig.subtitle" />
                  </a-form-item>
                  <a-row :gutter="[12, 0]">
                    <a-col :xs="24" :md="12">
                      <a-form-item label="展示数量">
                        <a-input-number v-model:value="rankingConfig.limit" :min="1" :max="20" style="width: 100%" />
                      </a-form-item>
                    </a-col>
                    <a-col :xs="24" :md="12">
                      <a-form-item label="启用">
                        <a-switch v-model:checked="rankingConfig.enabled" />
                      </a-form-item>
                    </a-col>
                  </a-row>
                  <a-button type="primary" html-type="submit">保存榜单配置</a-button>
                </a-form>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="featured" tab="推荐位">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="8">
              <a-card v-if="canManageFeaturedResources" class="inner-card portal-card">
                <div class="card-head">
                  <div class="card-kicker">首页推荐资源</div>
                  <h3 class="card-title">新增推荐位</h3>
                  <p class="card-description">控制首页的热门模型、热门数据集、任务模板和具身案例推荐内容。</p>
                </div>
                <a-form :model="createForm" layout="vertical" @finish="searchCandidates">
                  <a-form-item label="资源类型">
                    <a-select v-model:value="createForm.resource_type" :options="typeOptions" />
                  </a-form-item>
                  <a-form-item label="关键词">
                    <a-input v-model:value="createForm.keyword" placeholder="搜索推荐资源" />
                  </a-form-item>
                  <a-form-item label="排序值">
                    <a-input-number v-model:value="createForm.sort_order" :min="0" :max="1000" style="width: 100%" />
                  </a-form-item>
                  <a-form-item label="推荐标签">
                    <a-input v-model:value="createForm.badge_label" placeholder="例如：编辑推荐 / 场景精选" />
                  </a-form-item>
                  <a-form-item label="启用状态">
                    <a-switch v-model:checked="createForm.enabled" />
                  </a-form-item>
                  <a-form-item>
                    <a-button type="primary" html-type="submit" :loading="searching">搜索候选项</a-button>
                  </a-form-item>
                </a-form>

                <a-list v-if="candidates.length" :data-source="candidates" size="small">
                  <template #renderItem="{ item }">
                    <a-list-item>
                      <a-list-item-meta :title="item.title" :description="item.summary" />
                      <a-button type="link" @click="addFeatured(item)">加入</a-button>
                    </a-list-item>
                  </template>
                </a-list>
                <a-empty v-else description="先搜索候选项再加入推荐位" />
              </a-card>
            </a-col>

            <a-col :xs="24" :lg="16">
              <a-card v-if="canManageFeaturedResources" class="inner-card portal-card">
                <div class="card-head">
                  <div class="card-kicker">当前生效配置</div>
                  <h3 class="card-title">当前推荐位</h3>
                  <p class="card-description">每条推荐位都会显示在首页对应资源区。这里直接维护标签、排序和启用状态。</p>
                </div>
                <a-empty v-if="!featured.length" description="还没有配置推荐位" />
                <a-list v-else :data-source="featured">
                  <template #renderItem="{ item }">
                    <a-list-item class="featured-item">
                      <div class="featured-item-body">
                        <div class="featured-item-head">
                          <div>
                            <div class="list-title">{{ item.title }}</div>
                            <div class="list-desc">{{ item.summary }}</div>
                            <div class="list-hint">首页位置：{{ featuredDisplayArea(item.resource_type) }}</div>
                          </div>
                          <a-tag color="blue">{{ typeLabel(item.resource_type) }}</a-tag>
                        </div>
                        <a-space wrap style="margin-top: 10px">
                          <RouterLink v-if="item.route" :to="item.route">查看资源</RouterLink>
                          <span class="pill-meta">资源 ID {{ item.resource_id }}</span>
                        </a-space>
                        <div class="featured-editor-grid">
                          <div>
                            <div class="field-label">推荐标签</div>
                            <a-input v-model:value="item.badge_label" placeholder="例如：编辑推荐 / 标准模板" />
                          </div>
                          <div>
                            <div class="field-label">排序值</div>
                            <a-input-number v-model:value="item.sort_order" :min="0" :max="1000" style="width: 100%" />
                          </div>
                          <div>
                            <div class="field-label">启用</div>
                            <a-switch v-model:checked="item.enabled" />
                          </div>
                        </div>
                        <a-space style="margin-top: 14px">
                          <a-button type="primary" @click="saveFeatured(item)">保存</a-button>
                          <a-button danger @click="removeFeatured(item.id)">删除</a-button>
                        </a-space>
                      </div>
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="search" tab="搜索运营">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="8">
              <a-card v-if="canManageSearchKeywords" class="inner-card portal-card">
                <div class="card-head">
                  <div class="card-kicker">搜索入口</div>
                  <h3 class="card-title">搜索运营词</h3>
                  <p class="card-description">控制搜索页的热门词和推荐词，让首页搜索和全局搜索更容易被用户点击。</p>
                </div>
                <a-form :model="keywordForm" layout="vertical" @finish="saveKeyword">
                  <a-form-item label="词类型">
                    <a-select v-model:value="keywordForm.keyword_type" :options="keywordTypeOptions" />
                  </a-form-item>
                  <a-form-item label="关键词">
                    <a-input v-model:value="keywordForm.query" placeholder="例如：巡检 / 夜间安防" />
                  </a-form-item>
                  <a-row :gutter="[12, 0]">
                    <a-col :xs="24" :md="12">
                      <a-form-item label="排序值">
                        <a-input-number v-model:value="keywordForm.sort_order" :min="0" :max="1000" style="width: 100%" />
                      </a-form-item>
                    </a-col>
                    <a-col :xs="24" :md="12">
                      <a-form-item label="启用">
                        <a-switch v-model:checked="keywordForm.enabled" />
                      </a-form-item>
                    </a-col>
                  </a-row>
                  <a-space>
                    <a-button type="primary" html-type="submit">{{ keywordEditingId ? '保存' : '新增' }}</a-button>
                    <a-button @click="resetKeywordForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>

            <a-col :xs="24" :lg="16">
              <a-row :gutter="[16, 16]">
                <a-col :xs="24" :md="12">
                  <a-card v-if="canManageSearchKeywords" class="inner-card portal-card">
                    <div class="card-head">
                      <div class="card-kicker">前台位置</div>
                      <h3 class="card-title">热门词</h3>
                      <p class="card-description">显示在搜索页热门搜索区，适合运营当前最想让用户点击的关键词。</p>
                    </div>
                    <a-list :data-source="hotKeywords" size="small">
                      <template #renderItem="{ item }">
                        <a-list-item>
                          <a-list-item-meta :title="item.query" :description="`排序 ${item.sort_order}`" />
                          <a-space>
                            <a-tag :color="item.enabled ? 'green' : 'default'">{{ item.enabled ? '启用' : '停用' }}</a-tag>
                            <a-button type="link" @click="startEditKeyword(item)">编辑</a-button>
                            <a-button danger type="link" @click="deleteKeyword(item.id)">删除</a-button>
                          </a-space>
                        </a-list-item>
                      </template>
                    </a-list>
                    <a-empty v-if="!hotKeywords.length" description="还没有热门词" />
                  </a-card>
                </a-col>

                <a-col :xs="24" :md="12">
                  <a-card v-if="canManageSearchKeywords" class="inner-card portal-card">
                    <div class="card-head">
                      <div class="card-kicker">前台位置</div>
                      <h3 class="card-title">推荐词</h3>
                      <p class="card-description">显示在搜索页推荐词区，适合放场景词、专题词和运营引导词。</p>
                    </div>
                    <a-list :data-source="recommendedKeywords" size="small">
                      <template #renderItem="{ item }">
                        <a-list-item>
                          <a-list-item-meta :title="item.query" :description="`排序 ${item.sort_order}`" />
                          <a-space>
                            <a-tag :color="item.enabled ? 'green' : 'default'">{{ item.enabled ? '启用' : '停用' }}</a-tag>
                            <a-button type="link" @click="startEditKeyword(item)">编辑</a-button>
                            <a-button danger type="link" @click="deleteKeyword(item.id)">删除</a-button>
                          </a-space>
                        </a-list-item>
                      </template>
                    </a-list>
                    <a-empty v-if="!recommendedKeywords.length" description="还没有推荐词" />
                  </a-card>
                </a-col>
              </a-row>
            </a-col>
          </a-row>
        </a-tab-pane>
      </a-tabs>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue';
import { message, Modal } from 'ant-design-vue';
import { RouterLink } from 'vue-router';

import { api } from '@/api';
import { PERMISSIONS } from '@/constants/permissions';
import { useAuthStore } from '@/stores/auth';
import type { ApplicationCaseItem, FeaturedResourceItem, HomeHeroConfigItem, HomeHighlightItem, ModuleSettingItem, RankingConfigItem, ScenePageItem, SearchItem, SearchKeywordConfigItem, TaskTemplateItem } from '@/types/api';

type CandidateItem = {
  id: number;
  title: string;
  summary: string;
  resource_type: string;
};

const auth = useAuthStore();
const activeTab = ref('structure');
const modules = ref<ModuleSettingItem[]>([]);
const heroConfig = reactive<HomeHeroConfigItem>({
  tagline: 'OpenLoong 风格信息门户',
  title: '围绕模型、数据集、模板与文档的开放社区',
  description: '开放社区聚合模型、数据集、任务模板、文档与具身应用案例，帮助开发者围绕机器人场景快速完成接入、训练、部署和复用。',
  primary_button: '上传模型',
  secondary_button: '上传数据集',
  search_button: '进入全局搜索',
});
const highlights = ref<HomeHighlightItem[]>([]);
const rankingConfig = reactive<RankingConfigItem>({
  title: '贡献排行榜',
  subtitle: '基于积分展示近期社区贡献活跃度。',
  limit: 5,
  enabled: true,
});
const scenePages = ref<ScenePageItem[]>([]);
const featured = ref<FeaturedResourceItem[]>([]);
const searchKeywords = ref<SearchKeywordConfigItem[]>([]);
const candidates = ref<CandidateItem[]>([]);
const searching = ref(false);
const highlightEditingId = ref<number>();
const keywordEditingId = ref<number>();
const highlightForm = reactive({
  text: '',
  sort_order: 10,
  enabled: true,
});
const createForm = reactive({
  resource_type: 'model',
  keyword: '',
  badge_label: '',
  sort_order: 10,
  enabled: true,
});
const sceneEditingId = ref<number>();
const sceneForm = reactive({
  slug: '',
  name: '',
  tagline: '',
  summary: '',
  description: '',
  sort_order: 10,
  enabled: true,
});
const keywordForm = reactive({
  query: '',
  keyword_type: 'hot',
  sort_order: 10,
  enabled: true,
});

const typeOptions = [
  { label: '模型', value: 'model' },
  { label: '数据集', value: 'dataset' },
  { label: '任务模板', value: 'task-template' },
  { label: '具身案例', value: 'application-case' },
];
const keywordTypeOptions = [
  { label: '热门词', value: 'hot' },
  { label: '推荐词', value: 'recommended' },
];
const portalGuideCards = [
  {
    key: 'structure',
    label: '首页结构',
    title: '模块开关与亮点文案',
    description: '控制首页每个分区显示顺序，以及首屏右侧的亮点列表。',
    location: '首页导读区、模型推荐、公告动态、应用场景、数据集与模板、社区共创',
  },
  {
    key: 'hero',
    label: '导读区',
    title: '首页最上方主视觉',
    description: '控制首页最顶部蓝色标签、主标题、说明和三个主按钮文案。',
    location: '首页首屏最上方',
  },
  {
    key: 'scenes',
    label: '场景与榜单',
    title: '场景卡片与贡献榜',
    description: '控制首页应用场景卡片，以及社区共创里的贡献榜标题与显示数量。',
    location: '首页应用场景区、首页社区共创区',
  },
  {
    key: 'featured',
    label: '推荐位',
    title: '首页热门资源',
    description: '控制首页热门模型、热门数据集、任务模板和具身案例的推荐内容。',
    location: '首页资源推荐区',
  },
  {
    key: 'search',
    label: '搜索运营',
    title: '热门词与推荐词',
    description: '控制搜索页的热门词、推荐词，影响用户从首页和搜索页进入内容的路径。',
    location: '搜索页、首页搜索入口联动',
  },
] as const;
const canManagePortal = computed(() => auth.hasPermission(PERMISSIONS.portalManage));
const canManageSearchKeywords = computed(() => auth.hasPermission(PERMISSIONS.searchKeywordManage));
const canManageFeaturedResources = computed(() => auth.hasPermission(PERMISSIONS.featuredResourceManage));
const hotKeywords = computed(() => searchKeywords.value.filter((item) => item.keyword_type === 'hot'));
const recommendedKeywords = computed(() => searchKeywords.value.filter((item) => item.keyword_type === 'recommended'));

async function load() {
  const [moduleItems, heroConfigData, highlightItems, sceneItems, rankingConfigData, featuredItems, keywordItems] = await Promise.all([
    api.getAdminPortalModules(),
    api.getAdminHomeHeroConfig(),
    api.getAdminHomeHighlights(),
    api.getAdminScenePages(),
    api.getAdminRankingConfig(),
    api.getAdminFeaturedResources(),
    api.getAdminSearchKeywords(),
  ]);
  modules.value = moduleItems;
  heroConfig.tagline = heroConfigData.tagline;
  heroConfig.title = heroConfigData.title;
  heroConfig.description = heroConfigData.description;
  heroConfig.primary_button = heroConfigData.primary_button;
  heroConfig.secondary_button = heroConfigData.secondary_button;
  heroConfig.search_button = heroConfigData.search_button;
  highlights.value = highlightItems;
  scenePages.value = sceneItems;
  rankingConfig.title = rankingConfigData.title;
  rankingConfig.subtitle = rankingConfigData.subtitle;
  rankingConfig.limit = rankingConfigData.limit;
  rankingConfig.enabled = rankingConfigData.enabled;
  featured.value = featuredItems;
  searchKeywords.value = keywordItems;
}

async function saveModule(item: ModuleSettingItem) {
  await api.updateAdminPortalModule(item.module_key, {
    enabled: item.enabled,
    sort_order: item.sort_order,
  });
  message.success('首页模块配置已保存');
  await load();
}

async function saveHighlight() {
  if (highlightEditingId.value) {
    await api.updateAdminHomeHighlight(highlightEditingId.value, highlightForm);
    message.success('首页亮点已更新');
  } else {
    await api.createAdminHomeHighlight(highlightForm);
    message.success('首页亮点已添加');
  }
  resetHighlightForm();
  await load();
}

async function saveHeroConfig() {
  await api.updateAdminHomeHeroConfig({
    tagline: heroConfig.tagline,
    title: heroConfig.title,
    description: heroConfig.description,
    primary_button: heroConfig.primary_button,
    secondary_button: heroConfig.secondary_button,
    search_button: heroConfig.search_button,
  });
  message.success('首页导读区已更新');
  await load();
}

async function deleteHighlight(id: number) {
  if (!(await confirmDanger('确认删除首页亮点？', '删除后该首页亮点将无法恢复。'))) return;
  await api.deleteAdminHomeHighlight(id);
  message.success('首页亮点已删除');
  await load();
}

function startEditHighlight(item: HomeHighlightItem) {
  highlightEditingId.value = item.id;
  highlightForm.text = item.text;
  highlightForm.sort_order = item.sort_order;
  highlightForm.enabled = item.enabled;
}

function resetHighlightForm() {
  highlightEditingId.value = undefined;
  highlightForm.text = '';
  highlightForm.sort_order = 10;
  highlightForm.enabled = true;
}

async function saveRankingConfig() {
  await api.updateAdminRankingConfig({
    title: rankingConfig.title,
    subtitle: rankingConfig.subtitle,
    limit: rankingConfig.limit,
    enabled: rankingConfig.enabled,
  });
  message.success('榜单配置已保存');
  await load();
}

async function saveScenePage() {
  if (sceneEditingId.value) {
    await api.updateAdminScenePage(sceneEditingId.value, sceneForm);
    message.success('场景页已更新');
  } else {
    await api.createAdminScenePage(sceneForm);
    message.success('场景页已创建');
  }
  resetSceneForm();
  await load();
}

function startEditScenePage(item: ScenePageItem) {
  sceneEditingId.value = item.id;
  sceneForm.slug = item.slug;
  sceneForm.name = item.name;
  sceneForm.tagline = item.tagline;
  sceneForm.summary = item.summary;
  sceneForm.description = item.description;
  sceneForm.sort_order = item.sort_order;
  sceneForm.enabled = item.enabled;
}

async function deleteScenePage(id: number) {
  if (!(await confirmDanger('确认删除场景页？', '删除后该场景页配置将无法恢复。'))) return;
  await api.deleteAdminScenePage(id);
  message.success('场景页已删除');
  if (sceneEditingId.value === id) {
    resetSceneForm();
  }
  await load();
}

function resetSceneForm() {
  sceneEditingId.value = undefined;
  sceneForm.slug = '';
  sceneForm.name = '';
  sceneForm.tagline = '';
  sceneForm.summary = '';
  sceneForm.description = '';
  sceneForm.sort_order = 10;
  sceneForm.enabled = true;
}

async function searchCandidates() {
  searching.value = true;
  try {
    candidates.value = await loadCandidates();
  } finally {
    searching.value = false;
  }
}

async function loadCandidates() {
  const keyword = createForm.keyword.trim();
  if (!keyword) return [];

  if (createForm.resource_type === 'task-template') {
    const items = await api.listTemplates();
    return filterTemplates(items, keyword);
  }
  if (createForm.resource_type === 'application-case') {
    const items = await api.listApplicationCases();
    return filterCases(items, keyword);
  }

  const response = await api.search({
    q: keyword,
    type: createForm.resource_type,
    sort: 'hot',
  });
  return response.items.map((item: SearchItem) => ({
    id: item.id,
    title: item.title,
    summary: item.summary,
    resource_type: mapSearchType(item.type),
  }));
}

function filterTemplates(items: TaskTemplateItem[], keyword: string) {
  return items
    .filter((item) => `${item.name} ${item.summary} ${item.category}`.includes(keyword))
    .map((item) => ({
      id: item.id,
      title: item.name,
      summary: item.summary,
      resource_type: 'task-template',
    }));
}

function filterCases(items: ApplicationCaseItem[], keyword: string) {
  return items
    .filter((item) => `${item.title} ${item.summary} ${item.category}`.includes(keyword))
    .map((item) => ({
      id: item.id,
      title: item.title,
      summary: item.summary,
      resource_type: 'application-case',
    }));
}

function mapSearchType(value: string) {
  if (value === 'task-template') return 'task-template';
  return value;
}

async function addFeatured(item: CandidateItem) {
  await api.createAdminFeaturedResource({
    resource_type: item.resource_type,
    resource_id: item.id,
    badge_label: createForm.badge_label,
    sort_order: createForm.sort_order,
    enabled: createForm.enabled,
  });
  message.success('推荐位已添加');
  candidates.value = [];
  createForm.keyword = '';
  createForm.badge_label = '';
  await load();
}

async function saveFeatured(item: FeaturedResourceItem) {
  await api.updateAdminFeaturedResource(item.id, {
    resource_type: item.resource_type,
    resource_id: item.resource_id,
    badge_label: item.badge_label,
    sort_order: item.sort_order,
    enabled: item.enabled,
  });
  message.success('推荐位已保存');
  await load();
}

async function removeFeatured(id: number) {
  if (!(await confirmDanger('确认删除推荐位？', '删除后该推荐位将从首页与运营位移除。'))) return;
  await api.deleteAdminFeaturedResource(id);
  message.success('推荐位已删除');
  await load();
}

async function saveKeyword() {
  if (keywordEditingId.value) {
    await api.updateAdminSearchKeyword(keywordEditingId.value, keywordForm);
    message.success('搜索运营词已更新');
  } else {
    await api.createAdminSearchKeyword(keywordForm);
    message.success('搜索运营词已添加');
  }
  resetKeywordForm();
  await load();
}

async function deleteKeyword(id: number) {
  if (!(await confirmDanger('确认删除搜索运营词？', '删除后该词将不再参与热门词或推荐词展示。'))) return;
  await api.deleteAdminSearchKeyword(id);
  message.success('搜索运营词已删除');
  await load();
}

function startEditKeyword(item: SearchKeywordConfigItem) {
  keywordEditingId.value = item.id;
  keywordForm.query = item.query;
  keywordForm.keyword_type = item.keyword_type;
  keywordForm.sort_order = item.sort_order;
  keywordForm.enabled = item.enabled;
}

function resetKeywordForm() {
  keywordEditingId.value = undefined;
  keywordForm.query = '';
  keywordForm.keyword_type = 'hot';
  keywordForm.sort_order = 10;
  keywordForm.enabled = true;
}

function typeLabel(type: string) {
  return {
    model: '模型',
    dataset: '数据集',
    'task-template': '模板',
    'application-case': '案例',
  }[type] ?? type;
}

function keywordTypeLabel(type: string) {
  return {
    hot: '热门词',
    recommended: '推荐词',
  }[type] ?? type;
}

function moduleHelp(key: string) {
  return {
    hero: {
      description: '首页首屏导读区，显示标签、主标题、说明和按钮。',
      location: '首页顶部第一屏',
    },
    models: {
      description: '首页热门模型区，展示运营挑选的模型资源。',
      location: '首页模型推荐区',
    },
    announcements: {
      description: '首页公告动态区，展示置顶或最新公告。',
      location: '首页公告动态区',
    },
    scenes: {
      description: '首页应用场景区，展示场景入口卡片。',
      location: '首页应用场景区',
    },
    resources: {
      description: '首页数据集、模板、案例组合展示区。',
      location: '首页数据集与模板区',
    },
    community: {
      description: '首页社区共创区，展示 Wiki、讨论和排行榜。',
      location: '首页社区共创区',
    },
  }[key] ?? {
    description: '控制首页某个分区的显示状态和顺序。',
    location: '首页内容区',
  };
}

function featuredDisplayArea(type: string) {
  return {
    model: '首页「热门模型」',
    dataset: '首页「热门数据集」',
    'task-template': '首页「任务模板」',
    'application-case': '首页「具身案例」',
  }[type] ?? '首页推荐区';
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

onMounted(load);
</script>

<style scoped>
.hero-panel {
  padding: 28px;
  display: grid;
  gap: 20px;
  margin-bottom: 20px;
  background:
    radial-gradient(circle at top right, rgba(22, 119, 255, 0.16), transparent 30%),
    linear-gradient(145deg, #f8fbff 0%, #ffffff 58%, #f2fbf7 100%);
}

.hero-kicker,
.card-kicker,
.guide-label {
  color: var(--brand-strong);
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
}

.hero-title {
  margin-top: 12px;
}

.hero-subtitle {
  max-width: 920px;
  line-height: 1.8;
}

.guide-grid {
  display: grid;
  gap: 12px;
  grid-template-columns: repeat(5, minmax(0, 1fr));
}

.guide-card {
  padding: 16px 18px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.92);
  border: 1px solid rgba(220, 230, 242, 0.82);
}

.guide-title {
  margin-top: 10px;
  color: var(--text-main);
  font-size: 18px;
  font-weight: 700;
}

.guide-description,
.guide-location,
.card-description,
.list-hint {
  margin-top: 8px;
  color: var(--text-secondary);
  line-height: 1.7;
}

.block {
  padding: 24px;
}

.portal-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 20px;
}

.portal-card,
.preview-card,
.inner-card {
  border-radius: 18px;
}

.portal-card :deep(.ant-card-body),
.preview-card :deep(.ant-card-body) {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.card-head {
  margin-bottom: 2px;
}

.card-title {
  margin: 10px 0 0;
  color: var(--text-main);
  font-size: 22px;
  font-weight: 700;
}

.module-item,
.featured-item {
  align-items: flex-start;
}

.featured-item-body {
  width: 100%;
  min-width: 0;
}

.module-stack {
  display: grid;
  gap: 14px;
}

.module-entry {
  display: grid;
  gap: 14px;
  padding: 18px;
  border: 1px solid rgba(220, 230, 242, 0.82);
  border-radius: 18px;
  background: linear-gradient(180deg, #ffffff 0%, #f9fbff 100%);
}

.module-entry-copy {
  min-width: 0;
}

.module-entry-title {
  color: var(--text-main);
  font-size: 18px;
  font-weight: 700;
}

.module-entry-location,
.module-entry-description {
  margin-top: 8px;
  color: var(--text-secondary);
  line-height: 1.7;
}

.module-entry-controls {
  display: grid;
  gap: 12px;
  grid-template-columns: minmax(140px, 0.8fr) 100px minmax(120px, 0.7fr);
}

.module-control {
  min-width: 0;
}

.module-control-switch {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.module-control-action {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.list-desc {
  overflow-wrap: anywhere;
}

.hero-preview {
  display: grid;
  gap: 12px;
  padding: 18px;
  border-radius: 18px;
  background: linear-gradient(180deg, #fbfdff 0%, #eef4fb 100%);
  border: 1px solid rgba(220, 230, 242, 0.82);
}

.hero-preview-title {
  color: var(--text-main);
  font-size: 28px;
  font-weight: 700;
  line-height: 1.2;
}

.hero-preview-description {
  color: var(--text-secondary);
  line-height: 1.8;
}

.featured-item-head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
}

.featured-editor-grid {
  display: grid;
  gap: 12px;
  grid-template-columns: minmax(0, 1.4fr) minmax(140px, 0.8fr) minmax(120px, 0.6fr);
  margin-top: 14px;
}

.field-label {
  font-size: 12px;
  color: var(--text-secondary);
  margin-bottom: 6px;
}

@media (max-width: 1280px) {
  .guide-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .block :deep(.ant-row > .ant-col-8),
  .block :deep(.ant-row > .ant-col-12) {
    flex: 0 0 100%;
    max-width: 100%;
  }

  .guide-grid,
  .module-entry-controls,
  .featured-editor-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 960px) {
  .featured-item-head {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>

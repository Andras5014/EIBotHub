<template>
  <div class="page-shell">
    <div class="page-card block">
      <div class="section-head">
        <div>
          <h1 class="section-title">内容管理</h1>
          <p class="section-subtitle">统一维护模板、案例、文档分类、文档、FAQ 和视频教程，并支持批量状态操作与日志查看。</p>
        </div>
      </div>

      <a-tabs>
        <a-tab-pane key="templates" tab="任务模板">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="10">
              <a-card title="模板列表" class="inner-card">
                <div class="batch-bar">
                  <a-select v-model:value="templateBatchStatus" :options="statusOptions" style="width: 160px" />
                  <a-button :disabled="!selectedTemplateIds.length" @click="batchUpdateTemplates">批量更新状态</a-button>
                  <a-button danger :disabled="!selectedTemplateIds.length" @click="batchDeleteTemplates">批量删除</a-button>
                </div>
                <a-list :data-source="templates">
                  <template #renderItem="{ item }">
                    <a-list-item class="clickable-item">
                      <template #actions>
                        <a-button type="link" @click="startEditTemplate(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteTemplate(item.id)">删除</a-button>
                      </template>
                      <a-list-item-meta :title="item.name" :description="item.summary" />
                      <template #extra>
                        <a-space>
                          <a-tag>{{ item.status || 'published' }}</a-tag>
                          <a-checkbox :checked="selectedTemplateIds.includes(item.id)" @change="toggleTemplate(item.id, $event.target.checked)" />
                        </a-space>
                      </template>
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="14">
              <a-card :title="templateEditingId ? '编辑模板' : '新建模板'" class="inner-card">
                <a-form :model="templateForm" layout="vertical" @finish="saveTemplate">
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12"><a-form-item label="模板名称"><a-input v-model:value="templateForm.name" /></a-form-item></a-col>
                    <a-col :span="12"><a-form-item label="状态"><a-select v-model:value="templateForm.status" :options="statusOptions" /></a-form-item></a-col>
                  </a-row>
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12"><a-form-item label="分类"><a-input v-model:value="templateForm.category" /></a-form-item></a-col>
                    <a-col :span="12"><a-form-item label="场景"><a-input v-model:value="templateForm.scene" /></a-form-item></a-col>
                  </a-row>
                  <a-form-item label="摘要"><a-input v-model:value="templateForm.summary" /></a-form-item>
                  <a-form-item label="说明"><a-textarea v-model:value="templateForm.description" :rows="3" /></a-form-item>
                  <a-form-item label="操作步骤"><a-textarea v-model:value="templateForm.guide" :rows="4" /></a-form-item>
                  <a-form-item label="关联资源"><a-input v-model:value="templateForm.resource_ref" placeholder="多个资源逗号分隔" /></a-form-item>
                  <a-space>
                    <a-button type="primary" html-type="submit">保存</a-button>
                    <a-button @click="resetTemplateForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="cases" tab="具身案例">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="10">
              <a-card title="案例列表" class="inner-card">
                <div class="batch-bar">
                  <a-select v-model:value="caseBatchStatus" :options="statusOptions" style="width: 160px" />
                  <a-button :disabled="!selectedCaseIds.length" @click="batchUpdateCases">批量更新状态</a-button>
                  <a-button danger :disabled="!selectedCaseIds.length" @click="batchDeleteCases">批量删除</a-button>
                </div>
                <a-list :data-source="applicationCases">
                  <template #renderItem="{ item }">
                    <a-list-item class="clickable-item">
                      <template #actions>
                        <a-button type="link" @click="startEditCase(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteCase(item.id)">删除</a-button>
                      </template>
                      <a-list-item-meta :title="item.title" :description="item.summary" />
                      <template #extra>
                        <a-space>
                          <a-tag>{{ item.status || 'published' }}</a-tag>
                          <a-checkbox :checked="selectedCaseIds.includes(item.id)" @change="toggleCase(item.id, $event.target.checked)" />
                        </a-space>
                      </template>
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="14">
              <a-card :title="applicationCaseEditingId ? '编辑案例' : '新建案例'" class="inner-card">
                <a-form :model="applicationCaseForm" layout="vertical" @finish="saveApplicationCase">
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12"><a-form-item label="案例标题"><a-input v-model:value="applicationCaseForm.title" /></a-form-item></a-col>
                    <a-col :span="12"><a-form-item label="状态"><a-select v-model:value="applicationCaseForm.status" :options="statusOptions" /></a-form-item></a-col>
                  </a-row>
                  <a-form-item label="分类"><a-input v-model:value="applicationCaseForm.category" /></a-form-item>
                  <a-form-item label="摘要"><a-input v-model:value="applicationCaseForm.summary" /></a-form-item>
                  <a-form-item label="部署指南"><a-textarea v-model:value="applicationCaseForm.guide" :rows="5" /></a-form-item>
                  <a-form-item label="封面图"><a-input v-model:value="applicationCaseForm.cover_image" /></a-form-item>
                  <a-space>
                    <a-button type="primary" html-type="submit">保存</a-button>
                    <a-button @click="resetApplicationCaseForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="categories" tab="文档分类">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="10">
              <a-card title="分类列表" class="inner-card">
                <a-list :data-source="docCategories">
                  <template #renderItem="{ item }">
                    <a-list-item class="clickable-item">
                      <template #actions>
                        <a-button type="link" @click="startEditCategory(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteCategory(item.id)">删除</a-button>
                      </template>
                      <a-list-item-meta :title="item.name" :description="item.doc_type" />
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="14">
              <a-card :title="docCategoryEditingId ? '编辑分类' : '新建分类'" class="inner-card">
                <a-form :model="docCategoryForm" layout="vertical" @finish="saveCategory">
                  <a-form-item label="分类名称"><a-input v-model:value="docCategoryForm.name" /></a-form-item>
                  <a-form-item label="文档类型">
                    <a-select v-model:value="docCategoryForm.doc_type" :options="docTypeOptions" />
                  </a-form-item>
                  <a-space>
                    <a-button type="primary" html-type="submit">保存</a-button>
                    <a-button @click="resetCategoryForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="docs" tab="文档">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="10">
              <a-card title="文档列表" class="inner-card">
                <div class="batch-bar">
                  <a-select v-model:value="documentBatchStatus" :options="statusOptions" style="width: 160px" />
                  <a-button :disabled="!selectedDocumentIds.length" @click="batchUpdateDocuments">批量更新状态</a-button>
                  <a-button danger :disabled="!selectedDocumentIds.length" @click="batchDeleteDocuments">批量删除</a-button>
                </div>
                <a-list :data-source="documents">
                  <template #renderItem="{ item }">
                    <a-list-item class="clickable-item">
                      <template #actions>
                        <a-button type="link" @click="startEditDocument(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteDocument(item.id)">删除</a-button>
                      </template>
                      <a-list-item-meta :title="item.title" :description="item.summary" />
                      <template #extra>
                        <a-space>
                          <a-tag>{{ item.status || 'published' }}</a-tag>
                          <a-checkbox :checked="selectedDocumentIds.includes(item.id)" @change="toggleDocument(item.id, $event.target.checked)" />
                        </a-space>
                      </template>
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="14">
              <a-card :title="documentEditingId ? '编辑文档' : '新建文档'" class="inner-card">
                <a-form :model="documentForm" layout="vertical" @finish="saveDocument">
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12">
                      <a-form-item label="文档分类">
                        <a-select v-model:value="documentForm.category_id" :options="docCategoryOptions" />
                      </a-form-item>
                    </a-col>
                    <a-col :span="12">
                      <a-form-item label="状态">
                        <a-select v-model:value="documentForm.status" :options="statusOptions" />
                      </a-form-item>
                    </a-col>
                  </a-row>
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12">
                      <a-form-item label="文档类型">
                        <a-select v-model:value="documentForm.doc_type" :options="docTypeOptions" />
                      </a-form-item>
                    </a-col>
                    <a-col :span="12"><a-form-item label="标题"><a-input v-model:value="documentForm.title" /></a-form-item></a-col>
                  </a-row>
                  <a-form-item label="摘要"><a-input v-model:value="documentForm.summary" /></a-form-item>
                  <a-form-item label="正文"><a-textarea v-model:value="documentForm.content" :rows="8" /></a-form-item>
                  <a-space>
                    <a-button type="primary" html-type="submit">保存</a-button>
                    <a-button @click="resetDocumentForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="faqs" tab="FAQ">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="10">
              <a-card title="FAQ 列表" class="inner-card">
                <a-list :data-source="faqs">
                  <template #renderItem="{ item }">
                    <a-list-item class="clickable-item">
                      <template #actions>
                        <a-button type="link" @click="startEditFaq(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteFaq(item.id)">删除</a-button>
                      </template>
                      <a-list-item-meta :title="item.question" :description="item.answer" />
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="14">
              <a-card :title="faqEditingId ? '编辑 FAQ' : '新建 FAQ'" class="inner-card">
                <a-form :model="faqForm" layout="vertical" @finish="saveFaq">
                  <a-form-item label="问题"><a-input v-model:value="faqForm.question" /></a-form-item>
                  <a-form-item label="答案"><a-textarea v-model:value="faqForm.answer" :rows="6" /></a-form-item>
                  <a-space>
                    <a-button type="primary" html-type="submit">保存</a-button>
                    <a-button @click="resetFaqForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="videos" tab="视频教程">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="10">
              <a-card title="视频列表" class="inner-card">
                <div class="batch-bar">
                  <a-radio-group v-model:value="videoBatchActive">
                    <a-radio-button :value="true">批量启用</a-radio-button>
                    <a-radio-button :value="false">批量停用</a-radio-button>
                  </a-radio-group>
                  <a-button :disabled="!selectedVideoIds.length" @click="batchUpdateVideos">批量更新</a-button>
                  <a-button danger :disabled="!selectedVideoIds.length" @click="batchDeleteVideos">批量删除</a-button>
                </div>
                <a-list :data-source="videos">
                  <template #renderItem="{ item }">
                    <a-list-item class="clickable-item">
                      <template #actions>
                        <a-button type="link" @click="startEditVideo(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteVideo(item.id)">删除</a-button>
                      </template>
                      <a-list-item-meta :title="item.title" :description="item.summary" />
                      <template #extra>
                        <a-space>
                          <a-tag :color="item.active ? 'green' : 'default'">{{ item.active ? '启用' : '停用' }}</a-tag>
                          <a-checkbox :checked="selectedVideoIds.includes(item.id)" @change="toggleVideo(item.id, $event.target.checked)" />
                        </a-space>
                      </template>
                    </a-list-item>
                  </template>
                </a-list>
              </a-card>
            </a-col>
            <a-col :xs="24" :lg="14">
              <a-card :title="videoEditingId ? '编辑视频' : '新建视频'" class="inner-card">
                <a-form :model="videoForm" layout="vertical" @finish="saveVideo">
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12"><a-form-item label="标题"><a-input v-model:value="videoForm.title" /></a-form-item></a-col>
                    <a-col :span="12"><a-form-item label="分类"><a-input v-model:value="videoForm.category" /></a-form-item></a-col>
                  </a-row>
                  <a-form-item label="摘要"><a-input v-model:value="videoForm.summary" /></a-form-item>
                  <a-form-item label="视频链接"><a-input v-model:value="videoForm.link" /></a-form-item>
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12"><a-form-item label="排序"><a-input-number v-model:value="videoForm.sort_order" :min="0" :max="1000" style="width: 100%" /></a-form-item></a-col>
                    <a-col :span="12"><a-form-item label="启用"><a-switch v-model:checked="videoForm.active" /></a-form-item></a-col>
                  </a-row>
                  <a-space>
                    <a-button type="primary" html-type="submit">保存</a-button>
                    <a-button @click="resetVideoForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="dataset-configs" tab="数据集配置">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :lg="12">
              <a-card title="协议模板" class="inner-card">
                <a-list :data-source="agreementTemplates">
                  <template #renderItem="{ item }">
                    <a-list-item class="clickable-item">
                      <template #actions>
                        <a-button type="link" @click="startEditAgreementTemplate(item)">编辑</a-button>
                        <a-button danger type="link" @click="deleteAgreementTemplate(item.id)">删除</a-button>
                      </template>
                      <a-list-item-meta :title="item.name" :description="item.content" />
                      <template #extra>
                        <a-tag :color="item.active ? 'green' : 'default'">{{ item.active ? '启用' : '停用' }}</a-tag>
                      </template>
                    </a-list-item>
                  </template>
                </a-list>
                <a-divider />
                <a-form :model="agreementTemplateForm" layout="vertical" @finish="saveAgreementTemplate">
                  <a-form-item label="模板名称"><a-input v-model:value="agreementTemplateForm.name" /></a-form-item>
                  <a-form-item label="模板内容"><a-textarea v-model:value="agreementTemplateForm.content" :rows="5" /></a-form-item>
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12"><a-form-item label="排序"><a-input-number v-model:value="agreementTemplateForm.sort_order" :min="0" :max="1000" style="width: 100%" /></a-form-item></a-col>
                    <a-col :span="12"><a-form-item label="启用"><a-switch v-model:checked="agreementTemplateForm.active" /></a-form-item></a-col>
                  </a-row>
                  <a-space>
                    <a-button type="primary" html-type="submit">保存</a-button>
                    <a-button @click="resetAgreementTemplateForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>

            <a-col :xs="24" :lg="12">
              <a-card title="权限级别" class="inner-card">
                <a-list :data-source="privacyOptions">
                  <template #renderItem="{ item }">
                    <a-list-item class="clickable-item">
                      <template #actions>
                        <a-button type="link" @click="startEditPrivacyOption(item)">编辑</a-button>
                        <a-button danger type="link" @click="deletePrivacyOption(item.id)">删除</a-button>
                      </template>
                      <a-list-item-meta :title="`${item.name} · ${item.code}`" :description="item.description" />
                      <template #extra>
                        <a-tag :color="item.active ? 'green' : 'default'">{{ item.active ? '启用' : '停用' }}</a-tag>
                      </template>
                    </a-list-item>
                  </template>
                </a-list>
                <a-divider />
                <a-form :model="privacyOptionForm" layout="vertical" @finish="savePrivacyOption">
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12"><a-form-item label="代码"><a-input v-model:value="privacyOptionForm.code" /></a-form-item></a-col>
                    <a-col :span="12"><a-form-item label="名称"><a-input v-model:value="privacyOptionForm.name" /></a-form-item></a-col>
                  </a-row>
                  <a-form-item label="说明"><a-textarea v-model:value="privacyOptionForm.description" :rows="4" /></a-form-item>
                  <a-row :gutter="[16, 0]">
                    <a-col :span="12"><a-form-item label="排序"><a-input-number v-model:value="privacyOptionForm.sort_order" :min="0" :max="1000" style="width: 100%" /></a-form-item></a-col>
                    <a-col :span="12"><a-form-item label="启用"><a-switch v-model:checked="privacyOptionForm.active" /></a-form-item></a-col>
                  </a-row>
                  <a-space>
                    <a-button type="primary" html-type="submit">保存</a-button>
                    <a-button @click="resetPrivacyOptionForm">重置</a-button>
                  </a-space>
                </a-form>
              </a-card>
            </a-col>
          </a-row>
        </a-tab-pane>

        <a-tab-pane key="logs" tab="操作日志">
          <a-card title="最近操作" class="inner-card">
            <a-list :data-source="operationLogs">
              <template #renderItem="{ item }">
                <a-list-item>
                  <a-list-item-meta :description="`${item.resource_type} #${item.resource_id || '-'} · ${item.detail || '无附加信息'}`">
                    <template #title>
                      <a-space wrap>
                        <span>{{ item.summary }}</span>
                        <a-tag color="blue">{{ item.action }}</a-tag>
                        <span class="log-meta">{{ item.admin_name || `管理员 ${item.admin_user_id}` }}</span>
                        <span class="log-meta">{{ formatDate(item.created_at) }}</span>
                      </a-space>
                    </template>
                  </a-list-item-meta>
                </a-list-item>
              </template>
            </a-list>
          </a-card>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue';
import { message } from 'ant-design-vue';

import { api } from '@/api';
import type {
  AgreementTemplateItem,
  AdminOperationLogItem,
  ApplicationCaseItem,
  DatasetPrivacyOptionItem,
  DocumentCategoryItem,
  DocumentItem,
  FAQItem,
  TaskTemplateItem,
  VideoTutorialItem,
} from '@/types/api';

const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '待审核', value: 'pending' },
  { label: '已发布', value: 'published' },
  { label: '已驳回', value: 'rejected' },
];

const docTypeOptions = [
  { label: '平台文档', value: 'platform' },
  { label: '技术文档', value: 'technical' },
];

const templates = ref<TaskTemplateItem[]>([]);
const applicationCases = ref<ApplicationCaseItem[]>([]);
const docCategories = ref<DocumentCategoryItem[]>([]);
const documents = ref<DocumentItem[]>([]);
const faqs = ref<FAQItem[]>([]);
const videos = ref<VideoTutorialItem[]>([]);
const agreementTemplates = ref<AgreementTemplateItem[]>([]);
const privacyOptions = ref<DatasetPrivacyOptionItem[]>([]);
const operationLogs = ref<AdminOperationLogItem[]>([]);

const selectedTemplateIds = ref<number[]>([]);
const selectedCaseIds = ref<number[]>([]);
const selectedDocumentIds = ref<number[]>([]);
const selectedVideoIds = ref<number[]>([]);
const templateBatchStatus = ref('published');
const caseBatchStatus = ref('published');
const documentBatchStatus = ref('published');
const videoBatchActive = ref(true);

const templateEditingId = ref<number>();
const applicationCaseEditingId = ref<number>();
const docCategoryEditingId = ref<number>();
const documentEditingId = ref<number>();
const faqEditingId = ref<number>();
const videoEditingId = ref<number>();
const agreementTemplateEditingId = ref<number>();
const privacyOptionEditingId = ref<number>();

const templateForm = reactive({
  name: '',
  summary: '',
  description: '',
  category: '',
  scene: '',
  guide: '',
  resource_ref: '',
  status: 'published',
});

const applicationCaseForm = reactive({
  title: '',
  summary: '',
  category: '',
  guide: '',
  cover_image: '',
  status: 'published',
});

const docCategoryForm = reactive({
  name: '',
  doc_type: 'platform',
});

const documentForm = reactive({
  category_id: undefined as number | undefined,
  title: '',
  summary: '',
  content: '',
  doc_type: 'platform',
  status: 'published',
});

const faqForm = reactive({
  question: '',
  answer: '',
});

const videoForm = reactive({
  title: '',
  summary: '',
  link: '',
  category: '',
  sort_order: 10,
  active: true,
});

const agreementTemplateForm = reactive({
  name: '',
  content: '',
  sort_order: 10,
  active: true,
});

const privacyOptionForm = reactive({
  code: '',
  name: '',
  description: '',
  sort_order: 10,
  active: true,
});

const docCategoryOptions = ref<Array<{ label: string; value: number }>>([]);

async function load() {
  const [templateItems, caseItems, categoryItems, documentItems, faqItems, videoItems, agreementTemplateItems, privacyOptionItems, logItems] = await Promise.all([
    api.getAdminTemplates(),
    api.getAdminApplicationCases(),
    api.getAdminDocCategories(),
    api.getAdminDocuments(),
    api.getAdminFaqs(),
    api.getAdminVideos(),
    api.getAdminAgreementTemplates(),
    api.getAdminPrivacyOptions(),
    api.getAdminOperationLogs(),
  ]);
  templates.value = templateItems;
  applicationCases.value = caseItems;
  docCategories.value = categoryItems;
  documents.value = documentItems;
  faqs.value = faqItems;
  videos.value = videoItems;
  agreementTemplates.value = agreementTemplateItems;
  privacyOptions.value = privacyOptionItems;
  operationLogs.value = logItems;
  docCategoryOptions.value = categoryItems.map((item) => ({
    label: `${item.name} · ${item.doc_type}`,
    value: item.id,
  }));
}

function toggleSelection(target: number[], id: number, checked: boolean) {
  if (checked && !target.includes(id)) target.push(id);
  if (!checked) {
    const index = target.indexOf(id);
    if (index >= 0) target.splice(index, 1);
  }
}

function toggleTemplate(id: number, checked: boolean) {
  toggleSelection(selectedTemplateIds.value, id, checked);
}

function toggleCase(id: number, checked: boolean) {
  toggleSelection(selectedCaseIds.value, id, checked);
}

function toggleDocument(id: number, checked: boolean) {
  toggleSelection(selectedDocumentIds.value, id, checked);
}

function toggleVideo(id: number, checked: boolean) {
  toggleSelection(selectedVideoIds.value, id, checked);
}

async function saveTemplate() {
  if (templateEditingId.value) {
    await api.updateAdminTemplate(templateEditingId.value, templateForm);
    message.success('模板已更新');
  } else {
    await api.createAdminTemplate(templateForm);
    message.success('模板已创建');
  }
  resetTemplateForm();
  await load();
}

async function deleteTemplate(id: number) {
  await api.deleteAdminTemplate(id);
  message.success('模板已删除');
  await load();
}

async function batchUpdateTemplates() {
  await api.batchAdminTemplateStatus({ ids: selectedTemplateIds.value, status: templateBatchStatus.value });
  selectedTemplateIds.value = [];
  message.success('模板状态已批量更新');
  await load();
}

async function batchDeleteTemplates() {
  await api.batchDeleteAdminTemplates({ ids: selectedTemplateIds.value });
  selectedTemplateIds.value = [];
  message.success('模板已批量删除');
  await load();
}

function startEditTemplate(item: TaskTemplateItem) {
  templateEditingId.value = item.id;
  templateForm.name = item.name;
  templateForm.summary = item.summary;
  templateForm.description = item.description;
  templateForm.category = item.category;
  templateForm.scene = item.scene;
  templateForm.guide = item.guide;
  templateForm.resource_ref = item.resource_ref.join(',');
  templateForm.status = item.status || 'published';
}

function resetTemplateForm() {
  templateEditingId.value = undefined;
  templateForm.name = '';
  templateForm.summary = '';
  templateForm.description = '';
  templateForm.category = '';
  templateForm.scene = '';
  templateForm.guide = '';
  templateForm.resource_ref = '';
  templateForm.status = 'published';
}

async function saveApplicationCase() {
  if (applicationCaseEditingId.value) {
    await api.updateAdminApplicationCase(applicationCaseEditingId.value, applicationCaseForm);
    message.success('案例已更新');
  } else {
    await api.createAdminApplicationCase(applicationCaseForm);
    message.success('案例已创建');
  }
  resetApplicationCaseForm();
  await load();
}

async function deleteCase(id: number) {
  await api.deleteAdminApplicationCase(id);
  message.success('案例已删除');
  await load();
}

async function batchUpdateCases() {
  await api.batchAdminApplicationCaseStatus({ ids: selectedCaseIds.value, status: caseBatchStatus.value });
  selectedCaseIds.value = [];
  message.success('案例状态已批量更新');
  await load();
}

async function batchDeleteCases() {
  await api.batchDeleteAdminApplicationCases({ ids: selectedCaseIds.value });
  selectedCaseIds.value = [];
  message.success('案例已批量删除');
  await load();
}

function startEditCase(item: ApplicationCaseItem) {
  applicationCaseEditingId.value = item.id;
  applicationCaseForm.title = item.title;
  applicationCaseForm.summary = item.summary;
  applicationCaseForm.category = item.category;
  applicationCaseForm.guide = item.guide;
  applicationCaseForm.cover_image = item.cover_image ?? '';
  applicationCaseForm.status = item.status || 'published';
}

function resetApplicationCaseForm() {
  applicationCaseEditingId.value = undefined;
  applicationCaseForm.title = '';
  applicationCaseForm.summary = '';
  applicationCaseForm.category = '';
  applicationCaseForm.guide = '';
  applicationCaseForm.cover_image = '';
  applicationCaseForm.status = 'published';
}

async function saveCategory() {
  if (docCategoryEditingId.value) {
    await api.updateAdminDocCategory(docCategoryEditingId.value, docCategoryForm);
    message.success('分类已更新');
  } else {
    await api.createAdminDocCategory(docCategoryForm);
    message.success('分类已创建');
  }
  resetCategoryForm();
  await load();
}

async function deleteCategory(id: number) {
  try {
    await api.deleteAdminDocCategory(id);
    message.success('分类已删除');
    await load();
  } catch (error) {
    message.error(error instanceof Error ? error.message : '分类删除失败');
  }
}

function startEditCategory(item: DocumentCategoryItem) {
  docCategoryEditingId.value = item.id;
  docCategoryForm.name = item.name;
  docCategoryForm.doc_type = item.doc_type;
}

function resetCategoryForm() {
  docCategoryEditingId.value = undefined;
  docCategoryForm.name = '';
  docCategoryForm.doc_type = 'platform';
}

async function saveDocument() {
  if (!documentForm.category_id) {
    message.error('请选择文档分类');
    return;
  }
  const payload = {
    category_id: documentForm.category_id,
    title: documentForm.title,
    summary: documentForm.summary,
    content: documentForm.content,
    doc_type: documentForm.doc_type,
    status: documentForm.status,
  };
  if (documentEditingId.value) {
    await api.updateAdminDocument(documentEditingId.value, payload);
    message.success('文档已更新');
  } else {
    await api.createAdminDocument(payload);
    message.success('文档已创建');
  }
  resetDocumentForm();
  await load();
}

async function deleteDocument(id: number) {
  await api.deleteAdminDocument(id);
  message.success('文档已删除');
  await load();
}

async function batchUpdateDocuments() {
  await api.batchAdminDocumentStatus({ ids: selectedDocumentIds.value, status: documentBatchStatus.value });
  selectedDocumentIds.value = [];
  message.success('文档状态已批量更新');
  await load();
}

async function batchDeleteDocuments() {
  await api.batchDeleteAdminDocuments({ ids: selectedDocumentIds.value });
  selectedDocumentIds.value = [];
  message.success('文档已批量删除');
  await load();
}

function startEditDocument(item: DocumentItem) {
  documentEditingId.value = item.id;
  documentForm.category_id = item.category_id;
  documentForm.title = item.title;
  documentForm.summary = item.summary;
  documentForm.content = item.content;
  documentForm.doc_type = item.doc_type;
  documentForm.status = item.status || 'published';
}

function resetDocumentForm() {
  documentEditingId.value = undefined;
  documentForm.category_id = docCategories.value[0]?.id;
  documentForm.title = '';
  documentForm.summary = '';
  documentForm.content = '';
  documentForm.doc_type = 'platform';
  documentForm.status = 'published';
}

async function saveFaq() {
  if (faqEditingId.value) {
    await api.updateAdminFaq(faqEditingId.value, faqForm);
    message.success('FAQ 已更新');
  } else {
    await api.createAdminFaq(faqForm);
    message.success('FAQ 已创建');
  }
  resetFaqForm();
  await load();
}

async function deleteFaq(id: number) {
  await api.deleteAdminFaq(id);
  message.success('FAQ 已删除');
  await load();
}

function startEditFaq(item: FAQItem) {
  faqEditingId.value = item.id;
  faqForm.question = item.question;
  faqForm.answer = item.answer;
}

function resetFaqForm() {
  faqEditingId.value = undefined;
  faqForm.question = '';
  faqForm.answer = '';
}

async function saveVideo() {
  if (videoEditingId.value) {
    await api.updateAdminVideo(videoEditingId.value, videoForm);
    message.success('视频已更新');
  } else {
    await api.createAdminVideo(videoForm);
    message.success('视频已创建');
  }
  resetVideoForm();
  await load();
}

async function deleteVideo(id: number) {
  await api.deleteAdminVideo(id);
  message.success('视频已删除');
  await load();
}

async function batchUpdateVideos() {
  await api.batchAdminVideoStatus({ ids: selectedVideoIds.value, active: videoBatchActive.value });
  selectedVideoIds.value = [];
  message.success('视频状态已批量更新');
  await load();
}

async function batchDeleteVideos() {
  await api.batchDeleteAdminVideos({ ids: selectedVideoIds.value });
  selectedVideoIds.value = [];
  message.success('视频已批量删除');
  await load();
}

function startEditVideo(item: VideoTutorialItem) {
  videoEditingId.value = item.id;
  videoForm.title = item.title;
  videoForm.summary = item.summary;
  videoForm.link = item.link;
  videoForm.category = item.category;
  videoForm.sort_order = item.sort_order;
  videoForm.active = item.active;
}

function resetVideoForm() {
  videoEditingId.value = undefined;
  videoForm.title = '';
  videoForm.summary = '';
  videoForm.link = '';
  videoForm.category = '';
  videoForm.sort_order = 10;
  videoForm.active = true;
}

async function saveAgreementTemplate() {
  if (agreementTemplateEditingId.value) {
    await api.updateAdminAgreementTemplate(agreementTemplateEditingId.value, agreementTemplateForm);
    message.success('协议模板已更新');
  } else {
    await api.createAdminAgreementTemplate(agreementTemplateForm);
    message.success('协议模板已创建');
  }
  resetAgreementTemplateForm();
  await load();
}

async function deleteAgreementTemplate(id: number) {
  await api.deleteAdminAgreementTemplate(id);
  message.success('协议模板已删除');
  await load();
}

function startEditAgreementTemplate(item: AgreementTemplateItem) {
  agreementTemplateEditingId.value = item.id;
  agreementTemplateForm.name = item.name;
  agreementTemplateForm.content = item.content;
  agreementTemplateForm.sort_order = item.sort_order;
  agreementTemplateForm.active = item.active;
}

function resetAgreementTemplateForm() {
  agreementTemplateEditingId.value = undefined;
  agreementTemplateForm.name = '';
  agreementTemplateForm.content = '';
  agreementTemplateForm.sort_order = 10;
  agreementTemplateForm.active = true;
}

async function savePrivacyOption() {
  if (privacyOptionEditingId.value) {
    await api.updateAdminPrivacyOption(privacyOptionEditingId.value, privacyOptionForm);
    message.success('权限级别已更新');
  } else {
    await api.createAdminPrivacyOption(privacyOptionForm);
    message.success('权限级别已创建');
  }
  resetPrivacyOptionForm();
  await load();
}

async function deletePrivacyOption(id: number) {
  await api.deleteAdminPrivacyOption(id);
  message.success('权限级别已删除');
  await load();
}

function startEditPrivacyOption(item: DatasetPrivacyOptionItem) {
  privacyOptionEditingId.value = item.id;
  privacyOptionForm.code = item.code;
  privacyOptionForm.name = item.name;
  privacyOptionForm.description = item.description;
  privacyOptionForm.sort_order = item.sort_order;
  privacyOptionForm.active = item.active;
}

function resetPrivacyOptionForm() {
  privacyOptionEditingId.value = undefined;
  privacyOptionForm.code = '';
  privacyOptionForm.name = '';
  privacyOptionForm.description = '';
  privacyOptionForm.sort_order = 10;
  privacyOptionForm.active = true;
}

function formatDate(value: string) {
  return new Date(value).toLocaleString('zh-CN');
}

onMounted(async () => {
  await load();
  resetDocumentForm();
});
</script>

<style scoped>
.block {
  padding: 24px;
}

.inner-card {
  border-radius: 18px;
}

.clickable-item {
  cursor: pointer;
}

.batch-bar {
  display: flex;
  gap: 12px;
  align-items: center;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.log-meta {
  color: var(--text-secondary);
  font-size: 12px;
}
</style>

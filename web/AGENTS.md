# Web Agent Guide

## 适用范围
- 本文件适用于 `web/` 下的 Vue 3、TypeScript、Ant Design Vue、Pinia、Vue Router 相关代码与页面实现。
- 默认目录约定为 `src/views`、`src/components`、`src/api`、`src/stores`、`src/router`、`src/composables`。

## 代码与页面规范
- 统一使用 `script setup lang="ts"`，组件名使用 PascalCase，组合式函数使用 `useXxx`，Pinia store 使用 `useXxxStore`。
- 页面组件保持薄层，接口请求放到 `src/api`，跨组件状态与业务流程优先放到 store 或 composable。
- 优先复用 Ant Design Vue 组件，不先造自定义基础控件。
- 列表页默认采用“筛选区 + 工具栏 + 表格 + 分页”的结构；简单编辑优先使用 `Drawer` 或 `Modal`；详情页优先拆成 `Card`、`Descriptions` 等信息块。
- 首次实现就要覆盖 loading、empty、error、disabled、permission denied 等状态。
- 避免把路由守卫、复杂状态流转、接口串联和视图细节全部塞进一个单文件组件。

## 视觉与交互规范
- 样式界面参考 `https://forum.openloong.org.cn/` 的信息组织方式与视觉气质，而不是逐像素照搬。
- 优先参考其浅色背景、顶部导航、导读或热点模块、内容分区、卡片化版块、列表与统计并存的社区型布局。
- 当页面需要做首页、门户页、社区页、资讯页或工作台时，优先采用“头部导航 + 核心导读区 + 热门或最新内容区 + 分区卡片区”的信息架构。
- 保持较高信息密度，但模块边界必须清晰；标题、摘要、标签、数字指标和主要操作要一眼可扫。
- 避免默认紫色系、过度阴影和厚重暗色风格；优先使用克制、清爽、技术社区感的配色与排版。
- 响应式需要同时考虑桌面和移动端，桌面优先保证信息组织，移动端优先保证可读性和操作顺序。

## 交付要求
- 新页面需说明使用了哪些 API、状态来源和交互假设。
- 交付前至少执行 `pnpm lint`、`pnpm type-check` 和 `pnpm build`。
- 如果只是借鉴 OpenLoong 论坛的布局语言，应在说明中指出借鉴的是导航、分区、卡片和内容密度，不要直接复制其论坛实现细节。

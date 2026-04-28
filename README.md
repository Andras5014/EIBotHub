# EIBotHub

面向具身智能与机器人开发场景的开放社区原型项目，提供模型仓库、数据集仓库、任务模板、技能分享、文档中心、Wiki、讨论区、私信协作与后台运营能力。

## 项目简介

EIBotHub 目标是把机器人场景下分散的模型、数据集、任务模板、技能、文档和社区互动能力聚合到同一平台中，形成一个可浏览、可上传、可审核、可下载、可协作的开放社区。

当前仓库采用前后端分离架构：

- `web/`：Vue 3 + TypeScript + Ant Design Vue 前端
- `server/`：Go + Gin + Gorm 后端服务
- `docs/`：需求规划、迭代记录、优化建议

## 功能概览

### 门户与内容

- 门户首页聚合公告、推荐资源、数据集、模板、案例、Wiki 和贡献榜
- 全局搜索支持跨资源检索、筛选、排序、热门搜索和最近搜索
- 文档中心覆盖平台文档、技术文档、FAQ、视频教程
- Wiki 支持词条浏览、编辑、修订历史、锁定与回滚

### 资源仓库

- 模型仓库：列表、详情、上传、版本管理、下载记录、评分评论、评测
- 数据集仓库：列表、详情、上传、版本管理、多媒体样本预览、协议确认、分批下载、下载审批
- 任务模板：模板浏览、说明、评分评论
- 具身案例：案例展示、部署指南
- 技能分享：技能发布、详情、Fork、评分评论

### 用户与社区

- 用户注册、登录、注销
- 个人中心：资料、上传资源、收藏、下载记录、通知、关注关系、认证、积分
- 讨论区：发帖、评论、回复
- 私信会话与协作空间
- 开发者认证、企业认证、积分与贡献榜

### 后台运营

- 仪表盘
- 模型 / 数据集审核
- 首页模块开关与推荐位配置
- 公告管理
- 模板、案例、文档、FAQ、视频教程管理
- 协议模板与数据集权限级别管理
- 社区内容治理
- 数据集访问审批
- Wiki 治理
- 积分权益与操作日志

## 技术栈

### 前端

- Vue 3
- TypeScript
- Vite
- Vue Router
- Pinia
- Ant Design Vue
- Axios

### 后端

- Go
- Gin
- Gorm
- SQLite

## 目录结构

```text
EIBotHub/
├─ web/                     # Vue 3 前端
│  ├─ src/
│  │  ├─ api/               # 接口封装
│  │  ├─ components/        # 通用组件
│  │  ├─ layouts/           # 页面布局
│  │  ├─ router/            # 路由
│  │  ├─ stores/            # 状态管理
│  │  ├─ styles/            # 全局样式
│  │  ├─ types/             # TS 类型定义
│  │  └─ views/             # 页面
│  └─ vite.config.ts
├─ server/                  # Go 后端
│  ├─ cmd/api/              # 启动入口
│  └─ internal/
│     ├─ app/               # 应用装配与种子数据
│     ├─ config/            # 配置
│     ├─ dto/               # DTO
│     ├─ handler/           # HTTP 处理层
│     ├─ middleware/        # 中间件
│     ├─ model/             # 数据模型
│     ├─ repository/        # 数据访问层
│     ├─ service/           # 业务服务层
│     └─ support/           # 通用支持能力
├─ docs/                    # 规划与迭代文档
├─ AGENTS.md
└─ README.md
```

## 快速开始

### 环境要求

- Node.js 18+
- pnpm 8+
- Go 1.25+

### 1. 启动前端

```powershell
cd web
pnpm install
pnpm dev
```

默认地址：`http://127.0.0.1:5173`

### 2. 启动后端

```powershell
cd server
go run ./cmd/api
```

默认地址：`http://127.0.0.1:8080`

### 3. 本地校验

前端：

```powershell
cd web
pnpm type-check
pnpm build
```

后端：

```powershell
cd server
go test ./...
go vet ./...
```

## 默认配置

后端默认配置位于 `server/internal/config/config.go`。部署包会优先读取
`deploy/config.json`，也可以通过 `CONFIG_FILE` 指定其他 JSON 配置文件。

- `port` / `APP_PORT`：服务端口，默认 `8080`
- `db_path` / `DB_PATH`：数据库路径，默认 `server/data/opencommunity.db`
- `storage_dir` / `STORAGE_DIR`：本地文件目录，默认 `server/storage`
- `app_secret` / `APP_SECRET`：本地 JWT/Token 密钥
- `seed_demo` / `SEED_DEMO`：是否补充演示数据，默认 `true`
- `gin_mode` / `GIN_MODE`：Gin 运行模式，部署时建议 `release`

前端开发服务器默认：

- 端口：`5173`
- `/api` 代理到 `http://localhost:8080`
- `/storage` 代理到 `http://localhost:8080`

## 演示账号

项目默认会写入演示数据与演示账号：

- 普通用户：`demo@example.com / Demo123!`
- 管理员：`admin@opencommunity.local / Admin123!`

仅用于本地开发、联调和演示，请勿用于生产环境。

## 文档索引

- [开放社区研发拆分规划](./docs/%E5%BC%80%E6%94%BE%E7%A4%BE%E5%8C%BA%E7%A0%94%E5%8F%91%E6%8B%86%E5%88%86%E8%A7%84%E5%88%92.md)
- [迭代文档 2026-04-17](./docs/%E8%BF%AD%E4%BB%A3%E6%96%87%E6%A1%A3-2026-04-17.md)
- [功能完善优化建议 2026-04-18](./docs/%E5%8A%9F%E8%83%BD%E5%AE%8C%E5%96%84%E4%BC%98%E5%8C%96%E5%BB%BA%E8%AE%AE-2026-04-18.md)

## 当前状态

当前版本已经具备可演示、可联调的完整主流程，适合继续向以下方向推进：

- 权限体系与审核流细化
- 文件存储与下载授权升级
- 模型 / 数据集编辑闭环
- 搜索运营配置增强
- 数据集审批历史与权限矩阵
- 社区治理和协作治理能力增强

## 开发约定

- 本仓库默认采用前后端分离开发模式。
- 页面交互和 API 类型尽量显式定义，不直接泄露后端模型。
- 后端采用 `handler -> service -> repository` 分层。
- 本地运行数据、缓存、构建产物和日志已通过 `.gitignore` 排除，不会提交到远程仓库。

## License

本项目当前包含 [LICENSE](./LICENSE) 文件，具体使用方式请结合仓库约定与后续业务要求确认。

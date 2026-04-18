# EIBotHub

开放社区原型项目，面向具身智能 / 机器人开发场景，提供模型、数据集、任务模板、技能、文档与社区协作能力。

## 项目结构

- `web/`：Vue 3 + TypeScript + Ant Design Vue 前端
- `server/`：Go + Gin + Gorm 后端
- `docs/`：需求、迭代记录、规划文档

## 当前能力

- 门户首页、全局搜索、登录注册、个人中心
- 模型仓库、数据集仓库、任务模板、具身案例
- 文档中心、FAQ、视频教程、Wiki
- 技能分享、讨论区、私信、协作空间
- 开发者认证、企业认证、积分与排行榜
- 后台审核、推荐位、公告、内容管理、数据集访问审批

## 技术栈

### 前端

- Vue 3
- TypeScript
- Vue Router
- Pinia
- Ant Design Vue
- Vite

### 后端

- Go
- Gin
- Gorm
- SQLite

## 本地开发

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

### 3. 常用校验命令

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

后端默认配置见 `server/internal/config/config.go`：

- `APP_PORT`：默认 `8080`
- `DB_PATH`：默认 `server/data/opencommunity.db`
- `STORAGE_DIR`：默认 `server/storage`
- `APP_SECRET`：本地默认密钥
- `SEED_DEMO`：默认 `true`，启动时补充演示数据

前端开发服务器默认端口为 `5173`，并代理：

- `/api` -> `http://localhost:8080`
- `/storage` -> `http://localhost:8080`

## 演示账号

项目默认会补充演示账号：

- 普通用户：`demo@example.com / Demo123!`
- 管理员：`admin@opencommunity.local / Admin123!`

仅用于本地开发与演示，不应用于生产环境。

## 文档索引

- [研发拆分规划](./docs/%E5%BC%80%E6%94%BE%E7%A4%BE%E5%8C%BA%E7%A0%94%E5%8F%91%E6%8B%86%E5%88%86%E8%A7%84%E5%88%92.md)
- [迭代文档 2026-04-17](./docs/%E8%BF%AD%E4%BB%A3%E6%96%87%E6%A1%A3-2026-04-17.md)
- [功能完善优化建议 2026-04-18](./docs/%E5%8A%9F%E8%83%BD%E5%AE%8C%E5%96%84%E4%BC%98%E5%8C%96%E5%BB%BA%E8%AE%AE-2026-04-18.md)

## 提交说明

已在 `.gitignore` 中忽略本地运行时数据与构建产物，包括：

- IDE 配置目录
- 本地数据库与缓存
- `server/data/`、`server/storage/`
- `web/dist/`、日志文件、`node_modules/`

提交远程仓库前，建议至少执行一次：

```powershell
cd web
pnpm build

cd ..\server
go test ./...
```

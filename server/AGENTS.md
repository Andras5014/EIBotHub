# Server Agent Guide

## 适用范围
- 本文件适用于 `server/` 下的 Go、Gin、Gorm 相关代码与文档生成。
- 新增后端功能时，默认沿用 `cmd/api`、`internal/handler`、`internal/service`、`internal/repository`、`internal/model`、`internal/dto`、`internal/middleware` 的分层结构。

## 分层职责
- `handler` 只负责路由注册、参数绑定、基础校验、认证上下文读取、调用 service、返回统一响应。
- `service` 负责业务编排、权限判断、事务边界、跨 repository 聚合和依赖存储的业务校验。
- `repository` 负责全部 Gorm 访问逻辑，包括查询条件、分页、排序、`Preload`、事务内持久化和必要的 join 或 scope。
- `model` 只描述持久化结构，不直接作为接口返回体。
- `dto` 明确定义请求、响应、分页和筛选结构，禁止把数据库模型直接暴露给前端。

## 开发规范
- 包名保持小写、简短，代码必须通过 `gofmt`。
- `context.Context` 需要从 handler 传到 service 和 repository，不要在中间层丢失。
- 事务默认放在 service 边界控制，一个用例只允许一个清晰的事务入口。
- 不直接向客户端暴露原始 SQL 或 Gorm 错误，先映射成稳定的业务错误或统一错误码。
- handler 中不要堆积业务逻辑；repository 中不要写权限判断、响应拼装或控制层分支。
- 新增接口时优先定义 DTO、错误语义和响应结构，再补 handler、service、repository。

## 测试与交付
- 优先为 service 和工具函数编写表驱动测试。
- 对容易回归的查询逻辑补 repository 测试，对状态码和响应结构补 handler 测试。
- 修复缺陷时，能补回归测试就补回归测试。
- 交付前至少执行 `go test ./...` 和 `go vet ./...`，并在说明中标出迁移、配置项、事务边界或非直观查询行为。

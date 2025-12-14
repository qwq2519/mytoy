# Repository Guidelines

## 项目结构与模块
- `main.go`：程序入口，组装 DI 容器与 Gin 路由。
- `bootstrap/`：依赖注入配置（samber/do），数据库初始化。
- `controller/`：HTTP 入口与路由注册（`router.go`）。
- `model/`、`service/`、`util/`：领域模型、业务逻辑、通用工具的占位，按模块扩展（chat、character、knowledge/RAG、config、assets）。
- `storage/`：所有持久化资源；`db/`（SQLite）、`files/`（images/documents/exports）、`tmp/`（临时）。
- `config/`：配置文件（如 `config.toml`），支持环境变量覆盖。

## 构建、测试与开发命令
- `go mod tidy`：拉取并整理依赖。
- `go run .`：启动 HTTP 服务（默认 `:8080`，可用 `PORT` 覆盖）。
- `APP_DB_PATH=storage/db/app.sqlite go run .`：指定 SQLite 路径运行。
- `go test ./...`：运行全部 Go 测试。

## 代码风格与命名
- 格式：提交前执行 `gofmt`（或 `go fmt ./...`）。
- 包名：全小写短名，文件按功能命名，如 `chat_service.go`、`chat_controller.go`。
- 导出/非导出：类型与函数用大驼峰导出，小驼峰内部使用。
- 错误：返回包装后的 error，Handler 避免 panic。

## 测试规范
- 框架：标准库 `testing`。
- 命名：`_test.go` 中的 `TestXxx`，与被测代码同目录。
- 范围：优先服务层单测；Handler 可用 `httptest`。
- 运行：`go test ./...`，表驱动用例更易维护。

## 架构概览
- 模式：MVC + Service + DI.Controller 通过 `samber/do` 获取 Service；Service 依赖 GORM 与配置。
- 存储：所有数据集中在 `storage/`，便于备份与迁移。

## Codex 响应规范
- 语言：所有回复使用简体中文，保持简洁明了。
- 流程：每次回答先简述理解与计划（方案确认），再给正式答复/输出。
- 选项：能提供方案或思路时至少给出 3 种可行选项；若无法满足需说明原因。
- 结构：优先用条目/编号方便决策；引用代码或路径时使用反引号标注。
- 透明：遇到不确定或需要更多信息时先提出澄清，不盲目臆测。
- 具体实现一个功能时，提供一个简洁的实现方式，并列出多个可扩展点供用户选择是否需要，避免过度设计。

## 额外工作流程要求
- 执行任务前优先使用代码搜索工具（如 `rg`）确认是否已有相关,类似的实现或文件。
- 生成并维护 `codeReview.md`：作为 Go 后端项目的代码审查参考规范，涵盖功能正确性（逻辑、边界、异常）、Go 语言规范（命名、编码习惯）、代码简洁与复用、HTTP 接口参数校验、性能与安全、可测试性、日志处理等。
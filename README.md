# MyToy AI Assistant Backend

一个用于个人练手的轻量级后端项目，目标是实现类似 SillyTavern 的本地 AI 助手能力：

- 集成多个llm供应商
- 角色聊天
- 个人资料收藏
- RAG检索。

采用轻量化技术栈，并且后续会扩展更多小功能。

---

## 核心目标

- 提供一个 **REST API 后端**，用于：
  - 多会话 AI 聊天（基础 AI 助手）
  - 角色 / 世界观 / Author's Note 管理
  - 个人资料库 + 文档切分 + 向量检索（RAG）
  - 系统配置与资源管理（数据目录、导出、清理等）
- 保持技术栈 **轻量化**，适合作为个人项目练手和二次扩展。

---

## 技术栈

- 语言：Go 
- Web 框架：Gin
- ORM：GORM
- 数据库：SQLite（单文件，本地轻量化）
- 依赖注入：`samber/do`
- 配置：TOML

---

## 目录结构概览

当前项目目录结构（后续会随着模块增加而扩展）：

```text
.
├── bootstrap/           # 应用启动和依赖注入容器
│   └── container.go
├── config/              # 配置文件与配置加载（计划）
│   └── config.toml      # 默认配置，可以按需填充
├── controller/          # 控制器层（C），负责 HTTP 路由处理
│   └── router.go        # Gin 路由入口，挂载各模块路由
├── model/               # 模型层（M），GORM 实体定义（计划）
├── service/             # 业务服务层，封装核心逻辑（计划）
├── storage/             # 所有持久化数据与资源的根目录
│   ├── db/              # SQLite 数据库文件
│   ├── files/           # 文件资源（图片、文档、导出文件等）
│   │   └── images/
│   └── tmp/             # 临时文件
├── util/                # 通用工具函数（计划）
├── go.mod
└── main.go              # 应用入口，启动 HTTP 服务
```

> 说明：  
> - 当前只是骨架，`model/`、`service/`、`util/` 将随着业务模块逐步完善。  
> - `storage/` 是唯一的数据根目录，方便备份和迁移。

---

## 运行前准备

1. **安装 Go**

   推荐使用 Go 1.21+：

   - https://go.dev/dl/

2. **克隆项目**

   ```bash
   git clone <your-repo-url> mytoy
   cd mytoy
   ```

3. **拉取依赖**

   ```bash
   go mod tidy
   ```

   这一步会自动下载 Gin、GORM、SQLite 驱动、samber/do 等依赖。

---

## 架构设计概览

### MVC + Router

- **Model（M）**：在 `model/` 中定义领域模型和 GORM 实体，例如：
  - `ChatSession`、`ChatMessage`
  - `Character`、`WorldInfoEntry`
  - `Document`、`DocumentChunk`、`Embedding`
- **Service**：在 `service/` 中实现业务逻辑，组合多个 model/repository，屏蔽具体实现细节。
- **Controller（C）**：在 `controller/` 中实现 HTTP 处理函数（Gin Handler），负责：
  - 参数解析与校验
  - 调用 service 层
  - 返回统一格式的 JSON 响应
- **Router**：`controller/router.go` 中提供 `NewRouter` 函数：
  - 创建 Gin 引擎
  - 注册中间件
  - 注册各业务模块路由（如 chat、character、knowledge 等）

### 依赖注入（DI）

- `bootstrap/container.go` 使用 `samber/do` 初始化 DI 容器：
  - 当前注册的依赖：
    - `*gorm.DB`：通过 `gorm.Open(sqlite.Open(...))` 创建 SQLite 连接
  - 后续可以在此注册：
    - 配置对象（AppConfig）
    - 各业务 service（ChatService、CharacterService 等）
- 在 controller / service 中通过 `do.Invoke` / `do.MustInvoke` 获取依赖，实现松耦合。

### 存储设计

- SQLite 数据库文件：
  - 默认路径：`storage/db/app.sqlite`
  - 可通过环境变量 `APP_DB_PATH` 覆盖
- 文件资源：
  - 图片：`storage/files/images`
  - 文档（计划）：`storage/files/documents`
  - 导出数据（计划）：`storage/files/exports`
- 临时文件：`storage/tmp`

---

## 功能模块规划（后续迭代）

> 以下模块是项目的目标设计，当前代码仅实现了基础骨架和健康检查接口。

1. **Chat 模块（聊天 / 会话）**
   - 会话管理：创建 / 删除 / 重命名会话
   - 消息管理：发送消息、查看历史
   - 与 LLM Provider 集成：封装调用 OpenAI 等接口
   - RAG 集成：根据会话上下文检索知识库片段

2. **Character 模块（角色 / 世界观）**
   - 角色卡管理（人物设定、场景、示例对话等）
   - World Info / Lorebook：世界设定条目，按关键词或向量匹配

3. **Knowledge / RAG 模块（资料库与检索）**
   - 资料来源：文件上传、URL 抓取、自定义笔记
   - 文档切分：按段落或长度拆分
   - 向量化：调用外部 Embedding API 生成向量，存入 SQLite 或外部向量库
   - 检索与融合：根据当前对话上下文检索相关片段，并拼接到 Prompt 中

4. **Config 模块（配置管理）**
   - 系统配置：LLM API Key、默认模型、RAG 参数等
   - 配置存储：数据库 / 配置文件
   - 提供 API 给前端修改部分配置

5. **Asset 模块（资源管理）**
   - 图片、导出文件、附件的上传与管理
   - 聊天记录导出（JSON/Markdown）

6. **System 模块（系统工具）**
   - 健康检查、版本信息
   - 数据清理（删除孤立记录、无用文件）
   - 简单统计（会话数量、文档数量等）

---

## 开发约定（建议）

这些约定不是强制的，但有助于项目结构清晰、易扩展：

- 按模块拆分文件：
  - `model/chat.go`、`service/chat_service.go`、`controller/chat_controller.go` 等
- Controller 层只做 HTTP 相关逻辑，不写业务逻辑
- Service 层只依赖 model/repository 与其他 service，不依赖 Gin
- 所有对外 API 返回统一的 JSON 结构（例如 `{code, message, data}`）
- 新增模块时：
  1. 先定义 model
  2. 实现 service
  3. 在 controller 中暴露 HTTP 路由
  4. 在 `controller/router.go` 注册路由
  5. 在 `bootstrap/container.go` 注册依赖（如需要）

---

## 许可

个人练手项目，未正式指定开源协议。  
如需公开开源，可根据需要添加 `LICENSE` 文件。 


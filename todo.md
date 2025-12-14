# TODO 列表（添加前需与用户确认）

## 规则
- 仅在多步或复杂任务时使用本文件；每个 todo 必须先与用户确认方案后再写入。
- 每个 todo 记录状态，完成后需用户确认/CR 才能标记完成。
- 全部 todo 完成并经用户确认后才可清空列表。

## 当前 TODO
- [ ] 配置管理器：分段 struct + Extra map，atomic.Value 快照，原子写入文件
- [ ] 默认配置文件：生成 `config/server.toml`、`config/database.toml`、`config/logging.toml`
- [ ] 预留配置路由/服务占位：后续扩展 REST 接口读写配置
- [ ] 启动整合：容器注册配置管理器，入口使用配置初始化端口/DB 等
- [ ] 拆分 config 包为模块文件，并将 toml 读写工具迁移到 util 包

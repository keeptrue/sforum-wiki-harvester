# SForum Wiki Harvester

Harvest high-quality forum discussions into a site knowledge base via AI tools.

SForum 插件：把论坛里被时间验证的高质量讨论，沉淀为站点自建知识库。当前版本通过
Manifest V3 `aiTools[]` 向宿主 AI 补全循环贡献只读工具；词条提炼与入库面在后续
版本通过 `sforum.ai.complete@1` 与 Job Registry 实现。

- ID: `wiki.harvester`
- Name: **SForum Wiki Harvester**
- Type: `plugin`（非内置，operator 安装）
- Manifest contract: `sforum.manifest@3` / Protocol V2（`hashicorp-go-plugin`）
- Website: https://github.com/zhuchunshu/sforum

## AI 工具

| 工具 ID | 模型可见名 | 说明 |
| --- | --- | --- |
| `wiki.harvester.status.read` | `wiki-harvest-status` | 查询知识库采集状态：已沉淀词条数、追踪来源主题数，可按 `topicId` 过滤。 |

## 安装

1. 构建后端二进制：

   ```bash
   cd backend
   go build -o plugin .
   ```

2. 刷新 digest 并校验（需要 SForum 开发环境）：

   ```bash
   cd apps/api
   go run ./cmd/sforum extension digest --write /path/to/sforum-wiki-harvester
   go run ./cmd/sforum extension test --allow-scaffold /path/to/sforum-wiki-harvester
   ```

3. 通过 admin 上传安装，或在宿主配置 `EXTERNAL_EXTENSION_ROOTS` 指向插件父目录。
   安装需通过 exact-artifact trust review；启用后重启 API 生效。

4. 在站点 AI 设置 `reply.tools.allowed` 中显式列出 `wiki-harvest-status`，
   模型才能看到该工具。

## 开发

```bash
cd backend
go vet ./...
go test ./...
go build -o plugin .
```

## License

MIT

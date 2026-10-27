# SForum Wiki Harvester

Harvest high-quality forum discussions into a site knowledge base via AI tools.

SForum 插件：把论坛里被时间验证的高质量讨论，沉淀为站点自建知识库。当前版本通过
Manifest V3 `aiTools[]` 向宿主 AI 补全循环贡献只读工具；词条提炼与入库面在后续
版本通过 `sforum.ai.complete@1` 与 Job Registry 实现。

- ID: `wiki.harvester`
- Name: **SForum Wiki Harvester**
- Type: `plugin`（非内置，operator 安装）
- Manifest contract: `sforum.manifest@3` / Protocol V2（`hashicorp-go-plugin`）
- 兼容性: `sforumVersion: >=3.0.0`
- Website: https://github.com/keeptrue/sforum

## AI 工具

| 工具 ID | 模型可见名 | 说明 |
| --- | --- | --- |
| `wiki.harvester.status.read` | `wiki-harvest-status` | 查询知识库采集状态：已沉淀词条数、追踪来源主题数，可按 `topicId` 过滤。 |

## 安装

从 [Releases](https://github.com/keeptrue/sforum-wiki-harvester/releases) 下载与服务器
平台匹配的压缩包，在后台「插件管理 → 安装插件」上传，或在宿主配置
`EXTERNAL_EXTENSION_ROOTS` 指向解压后的父目录。安装需通过 exact-artifact trust
review；启用后重启 API 生效。最后在站点 AI 设置 `reply.tools.allowed` 中显式列出
`wiki-harvest-status`，模型才能看到该工具。

## 兼容性

`sforumVersion` 声明为 `>=3.0.0`。宿主自 **v3.0.17** 起真正执行该约束：声明范围不
覆盖运行中的 `coreVersion` 时，安装返回 422、启用/升级/重启返回 409，错误码
`extension.core_incompatible`，并在响应里带上 `required` / `coreVersion` /
`buildVersion`。

注意 `extension validate` 与 `extension test` **都不检查** `sforumVersion`，两者全绿
不代表装得上；CI 里的 `Host core compatibility gate` 直接调用宿主的
`CheckSForumVersion`，与安装闸门同源。

## 开发

```bash
cd backend
go vet ./...
go test ./...
go build -o plugin .
```

`backend/go.mod` 通过 `replace` 指向 `upstream/sforum` 子模块，改动宿主面后需要同步
子模块指针。

刷新摘要与校验（需要一个 SForum 开发环境）：

```bash
cd /path/to/SForum/apps/api
go run ./cmd/sforum extension digest --write /path/to/sforum-wiki-harvester
go run ./cmd/sforum extension validate  /path/to/sforum-wiki-harvester
go run ./cmd/sforum extension test      /path/to/sforum-wiki-harvester
```

## CI / 发布

- `.github/workflows/ci.yml`：推 `master` / PR 时检出 SForum 源码，覆盖 `replace` 后跑
  `go vet`、`go test`、构建二进制、刷新摘要、`extension validate` / `test`，最后跑核心
  兼容闸门。
- `.github/workflows/release.yml`：推 `v*` 标签即发布。`verify` 先做同一套校验并核对
  标签与 manifest `version` 一致，`build` 按 linux/darwin × amd64/arm64 四平台各自构建
  并**按平台**刷新摘要、打包，`publish` 汇总 `SHA256SUMS` 挂到 GitHub Release。

发布一个新版本：改 `sforum.extension.json` 的 `version`，提交后打同名标签。

```bash
git tag v0.1.1 && git push origin v0.1.1
```

## License

MIT

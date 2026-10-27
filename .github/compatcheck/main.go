// Command compatcheck 用宿主自己的兼容求值函数判定本扩展能否装上运行中的 SForum。
//
// `extension validate` 与 `extension test` 都不检查 sforumVersion，真正拦住安装/启用
// 的是生命周期里的 ensureCoreCompatible。CI 把本文件复制进宿主源码树的
// apps/api/cmd/compatcheck 再 `go run`，因此判定逻辑与宿主逐字一致，不存在复刻漂移。
//
// 用法：compatcheck <sforum.extension.json 路径>
package main

import (
	"encoding/json"
	"fmt"
	"os"

	extensionmanifest "github.com/keeptrue/sforum/apps/api/app/Support/ExtensionManifest"
	platformversion "github.com/keeptrue/sforum/apps/api/version"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: compatcheck <manifest.json>")
		os.Exit(2)
	}
	body, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "read manifest: %v\n", err)
		os.Exit(1)
	}
	var manifest extensionmanifest.Manifest
	if err := json.Unmarshal(body, &manifest); err != nil {
		fmt.Fprintf(os.Stderr, "parse manifest: %v\n", err)
		os.Exit(1)
	}

	identity := platformversion.Identity()
	if !identity.Valid {
		// 宿主在 coreVersion 不是严格 semver 时也会跳过强制，这里保持一致。
		fmt.Printf("core version %q is not strict semver; host skips enforcement, so does this gate\n", identity.CoreVersion)
		return
	}
	if err := extensionmanifest.CheckSForumVersion(manifest, identity.CoreVersion); err != nil {
		fmt.Fprintf(os.Stderr, "core compatibility failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("core compatibility ok: sforumVersion=%q covers core=%s\n", manifest.SForumVersion, identity.CoreVersion)
}

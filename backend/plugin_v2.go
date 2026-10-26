package main

import (
	pluginv2 "github.com/keeptrue/sforum/apps/api/sdk/plugin/v2"

	"github.com/keeptrue/sforum-wiki-harvester/backend/wiki"
)

// wikiHarvesterPluginV2 装配 Wiki Harvester 的 AI 工具注册表。
// 工具声明与 manifest `aiTools[]` 一一对应，ID 必须落在 wiki.harvester 命名空间下。
type wikiHarvesterPluginV2 struct {
	*pluginv2.Server
}

func newWikiHarvesterPluginV2() *wikiHarvesterPluginV2 {
	registry, err := pluginv2.NewAIToolRegistry(wiki.ToolDefinitions()...)
	if err != nil {
		// 声明是编译期常量，失败只可能是作者代码错误，直接 panic 暴露。
		panic(err)
	}
	server := pluginv2.NewServer().WithAIToolRegistry(registry)
	return &wikiHarvesterPluginV2{Server: server}
}

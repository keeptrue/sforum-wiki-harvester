// Package wiki 实现 Wiki Harvester 的插件侧 AI 工具。
//
// 当前版本提供 wiki-harvest-status：模型在补全循环里查询本插件的采集状态。
// 词条沉淀的写入面（事件监听 + sforum.ai.complete 提炼）属于下一阶段，
// 届时通过 Host Command 与 Job Registry 实现，不使用 raw DB。
package wiki

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	pluginv2 "github.com/zhuchunshu/sforum/apps/api/sdk/plugin/v2"
)

const (
	toolIDStatus      = "wiki.harvester.status.read"
	toolNameStatus    = "wiki-harvest-status"
	statusTimeoutMS   = 1000
	statusContractVer = toolIDStatus + "@1"
)

// statusAt 是可注入时钟，便于测试。
var statusAt = time.Now

// ToolDefinitions 返回与 manifest `aiTools[]` 严格同源的声明列表。
func ToolDefinitions() []pluginv2.AIToolDefinition {
	return []pluginv2.AIToolDefinition{
		{
			ID:              toolIDStatus,
			Name:            toolNameStatus,
			Description:     "Report the Wiki Harvester knowledge-base status: harvested entry counts per source topic.",
			Parameters:      statusParametersSchema,
			TimeoutMS:       statusTimeoutMS,
			ContractVersion: statusContractVer,
			Execute:         handleStatusRead,
		},
	}
}

const statusParametersSchema = `{
  "type": "object",
  "properties": {
    "topicId": {
      "type": "integer",
      "description": "Optional source topic id. Omit for a global summary."
    }
  }
}`

// handleStatusRead 响应模型的状态查询。
// 入库面尚未实现时返回明确的可解释失败（IsError），让模型知道工具可用但
// 数据为空，而不是让整次补全失败。
func handleStatusRead(_ context.Context, call *pluginv2.AIToolCall) (*pluginv2.AIToolResult, error) {
	args := strings.TrimSpace(call.Arguments)
	payload := map[string]any{}
	if args != "" {
		if err := json.Unmarshal([]byte(args), &payload); err != nil {
			return &pluginv2.AIToolResult{
				Content: `{"error":"invalid arguments","hint":"pass {\"topicId\": <int>} or no arguments"}`,
				IsError: true,
			}, nil
		}
	}
	summary := map[string]any{
		"harvestedEntries": 0,
		"sourcesTracked":   0,
		"checkedAt":        statusAt().UTC().Format(time.RFC3339),
		"note":             "harvest ingestion lands in a later milestone; the store is empty",
	}
	if topicID, ok := payload["topicId"].(float64); ok {
		summary["topicId"] = int64(topicID)
		summary["topicEntries"] = 0
	}
	content, err := json.Marshal(summary)
	if err != nil {
		return nil, fmt.Errorf("encode status summary: %w", err)
	}
	return &pluginv2.AIToolResult{Content: string(content)}, nil
}

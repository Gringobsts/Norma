package harness

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/Gringobsts/Norma/llm"
	"github.com/Gringobsts/Norma/permission"
	"github.com/Gringobsts/Norma/tool"
)

// execOne must turn a panic into an error result and refuse, at execution time,
// direct calls to deferred tools and calls to DisabledTools while settling.
func TestExecOneRecoversAndGates(t *testing.T) {
	ran := false
	mk := func(name string, panics bool) tool.CoreTool {
		return tool.Build(tool.Spec{
			Name:   name,
			Schema: map[string]any{"type": "object"},
			Permissions: func(context.Context, json.RawMessage, permission.Context) permission.Decision {
				return permission.Decision{Behavior: permission.Allow}
			},
			Run: func(context.Context, json.RawMessage, *tool.ToolContext) (tool.Result, error) {
				if panics {
					panic("boom")
				}
				ran = true
				return tool.Text("ran"), nil
			},
		})
	}
	reg := tool.NewRegistry(mk("Boom", true), mk("Hidden", false), mk("Bash2", false))
	l := &loop{ctx: context.Background(), toolCtx: context.Background(), in: QueryInput{
		Tools:          reg,
		PermissionMode: permission.ModeBypass,
		DeferredTools:  []string{"Hidden"},
		Settlement:     &Settlement{Prompt: "x", DisabledTools: []string{"Bash2"}},
	}}
	call := func(name string) llm.ContentBlock {
		res, _ := l.execOne(llm.ContentBlock{Type: llm.BlockToolUse, ID: "t", Name: name, Input: json.RawMessage(`{}`)}, nil)
		return res
	}
	if !call("Boom").IsError {
		t.Fatal("panicking tool must yield an error result")
	}
	if !call("Hidden").IsError || ran {
		t.Fatal("deferred tool must not run when called directly")
	}
	if call("Bash2").IsError || !ran {
		t.Fatal("tool should run outside settlement")
	}
	ran = false
	l.settling = true
	if !call("Bash2").IsError || ran {
		t.Fatal("disabled tool must not run while settling")
	}
}

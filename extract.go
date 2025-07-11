// wafcheck/extract.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type tfPlan struct {
	ResourceChanges []struct {
		Type   string `json:"type"`
		Name   string `json:"name"`
		Change struct {
			Actions []string `json:"actions"`
			After   struct {
				Name  string     `json:"name"`
				Rules []MockRule `json:"rules"`
			} `json:"after"`
		} `json:"change"`
	} `json:"resource_changes"`
}

func RunExtract(planPath, zoneFilter string, onlyChanged bool, outPath string) error {
	b, err := os.ReadFile(planPath)
	if err != nil {
		return fmt.Errorf("reading plan file: %w", err)
	}

	var plan tfPlan
	if err := json.Unmarshal(b, &plan); err != nil {
		return fmt.Errorf("unmarshaling plan: %w", err)
	}

	rules := []MockRule{}
	for _, rc := range plan.ResourceChanges {
		if rc.Type != "cloudflare_ruleset" {
			continue
		}
		if zoneFilter != "" && !strings.Contains(rc.Change.After.Name, zoneFilter) {
			continue
		}
		if onlyChanged && !hasChange(rc.Change.Actions) {
			continue
		}
		for _, rule := range rc.Change.After.Rules {
			rules = append(rules, rule)
		}
	}

	ruleset := struct {
		Rules []MockRule `json:"rules"`
	}{Rules: rules}

	data, err := json.MarshalIndent(ruleset, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling output: %w", err)
	}
	if err := os.WriteFile(outPath, data, 0644); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	return nil
}

func hasChange(actions []string) bool {
	for _, a := range actions {
		if a == "create" || a == "update" || a == "delete" {
			return true
		}
	}
	return false
}

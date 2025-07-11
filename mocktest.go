// wafcheck/mocktest.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type MockRequest struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	IP     string `json:"ip"`
}

type MockRule struct {
	Expression  string   `json:"expression"`
	Action      string   `json:"action"`
	Description string   `json:"description"`
	IPSet       []string `json:"ip_set,omitempty"`
}

func matchExpression(expr string, req MockRequest, ipSet []string) bool {
	expr = strings.TrimSpace(expr)

	if strings.Contains(expr, "http.request.uri.path contains") {
		re := regexp.MustCompile(`(?i)contains \"(.*?)\"`)
		if m := re.FindStringSubmatch(expr); len(m) > 1 {
			return strings.Contains(req.Path, m[1])
		}
	}

	if strings.Contains(expr, "http.request.uri.path eq") {
		re := regexp.MustCompile(`(?i)eq \"(.*?)\"`)
		if m := re.FindStringSubmatch(expr); len(m) > 1 {
			return req.Path == m[1]
		}
	}

	if strings.Contains(expr, "http.request.method eq") {
		re := regexp.MustCompile(`(?i)method eq \"(.*?)\"`)
		if m := re.FindStringSubmatch(expr); len(m) > 1 {
			return strings.EqualFold(req.Method, m[1])
		}
	}

	if strings.Contains(expr, "ip.src in") {
		for _, ip := range ipSet {
			if req.IP == ip {
				return true
			}
		}
		return false
	}

	if strings.Contains(expr, "not ip.src in") {
		for _, ip := range ipSet {
			if req.IP == ip {
				return false
			}
		}
		return true
	}

	return false
}

func RunMockTest(rulesFile, reqsFile string) error {
	rData, err := os.ReadFile(rulesFile)
	if err != nil {
		return fmt.Errorf("reading rules file: %w", err)
	}
	qData, err := os.ReadFile(reqsFile)
	if err != nil {
		return fmt.Errorf("reading requests file: %w", err)
	}

	var ruleset struct {
		Rules []MockRule `json:"rules"`
	}
	var requests []MockRequest

	json.Unmarshal(rData, &ruleset)
	json.Unmarshal(qData, &requests)

	for _, req := range requests {
		matched := false
		for _, rule := range ruleset.Rules {
			if matchExpression(rule.Expression, req, rule.IPSet) {
				fmt.Printf("🚨 Matched: %s %s by rule: %s → action: %s\n", req.Method, req.Path, rule.Description, rule.Action)
				matched = true
				break
			}
		}
		if !matched {
			fmt.Printf("✅ No match: %s %s\n", req.Method, req.Path)
		}
	}
	return nil
}

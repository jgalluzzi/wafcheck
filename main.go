// cmd/wafcheck/main.go
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: wafcheck <command> [args]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "mocktest":
		mockCmd := flag.NewFlagSet("mocktest", flag.ExitOnError)
		rulesFile := mockCmd.String("rules", "", "Path to ruleset file")
		reqsFile := mockCmd.String("requests", "", "Path to mock requests JSON")
		mockCmd.Parse(os.Args[2:])

		if err := RunMockTest(*rulesFile, *reqsFile); err != nil {
			log.Fatal(err)
		}

	case "extract":
		extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)
		planFile := extractCmd.String("plan", "", "Path to tfplan.json file")
		zoneFilter := extractCmd.String("zone", "", "Zone name to filter")
		onlyChanged := extractCmd.Bool("only-changed", false, "Only extract changed rulesets")
		output := extractCmd.String("out", "rules.json", "Output file for ruleset")
		extractCmd.Parse(os.Args[2:])

		if err := RunExtract(*planFile, *zoneFilter, *onlyChanged, *output); err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("Unknown command. Available: mocktest, extract")
		os.Exit(1)
	}
}

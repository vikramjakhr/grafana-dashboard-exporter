package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs"
	"github.com/vikramjakhr/grafana-dashboard-exporter/config"
	_ "github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs/all"
	_ "github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs/all"
	"log"
)

var fQuiet = flag.Bool("quiet", false,
	"run in quiet mode")
var fTest = flag.Bool("test", false, "gather metrics, print them out, and exit")
var fConfig = flag.String("config", "", "configuration file to load")
var fVersion = flag.Bool("version", false, "display the version")
var fUsage = flag.String("usage", "",
	"print usage for a plugin, ie, 'grafana-dashboard-exporter --usage s3'")
var fSampleConfig = flag.Bool("sample-config", false,
	"print out full sample configuration")
var fPidfile = flag.String("pidfile", "", "file to write our pid to")
var fInputFilters = flag.String("input-filter", "",
	"filter the inputs to enable, separator is :")
var fInputList = flag.Bool("input-list", false,
	"print available input plugins.")
var fOutputFilters = flag.String("output-filter", "",
	"filter the outputs to enable, separator is :")
var fOutputList = flag.Bool("output-list", false,
	"print available output plugins.")

var (
	nextVersion = "1.0.0"
	version     string
	commit      string
	branch      string
)

func init() {
	// If commit or branch are not set, make that clear.
	if commit == "" {
		commit = "unknown"
	}
	if branch == "" {
		branch = "unknown"
	}
}

func displayVersion() string {
	if version == "" {
		return fmt.Sprintf("v%s~%s", nextVersion, commit)
	}
	return "v" + version
}

const usage = `Grafana-dashboard-exporter, The plugin-driven agent for exporting grafana dashboards json.

Usage:

  grafana-dashboard-exporter [commands|flags]

The commands & flags are:

  config              print out full sample configuration to stdout
  version             print the version to stdout

  --config <file>     configuration file to load
  --test              gather metrics once, print them to stdout, and exit
  --usage             print usage for a plugin, ie, 'grafana-dashboard-exporter --usage s3'
  --quiet             run in quiet mode

Examples:

  # generate a grafana-dashboard-exporter config file:
  grafana-dashboard-exporter config > gde.conf

  # run a single grafana-dashboard-exporter collection, outputing json to stdout
  grafana-dashboard-exporter --config gde.conf --test

  # run grafana-dashboard-exporter with all plugins defined in config file
  grafana-dashboard-exporter --config gde.conf
`

func usageExit(rc int) {
	fmt.Println(usage)
	os.Exit(rc)
}

func main() {
	flag.Usage = func() { usageExit(0) }
	flag.Parse()
	args := flag.Args()

	inputFilters, outputFilters := []string{}, []string{}
	if *fInputFilters != "" {
		inputFilters = strings.Split(":"+strings.TrimSpace(*fInputFilters)+":", ":")
	}
	if *fOutputFilters != "" {
		outputFilters = strings.Split(":"+strings.TrimSpace(*fOutputFilters)+":", ":")
	}

	if len(args) > 0 {
		switch args[0] {
		case "version":
			fmt.Printf("Grafana-dashboard-exporter %s (git: %s %s)\n", displayVersion(), branch, commit)
			return
		case "config":
			config.PrintSampleConfig(
				inputFilters,
				outputFilters,
			)
			return
		}
	}

	// switch for flags which just do something and exit immediately
	switch {
	case *fOutputList:
		fmt.Println("Available Output Plugins:")
		for k, _ := range outputs.Outputs {
			fmt.Printf("  %s\n", k)
		}
		return
	case *fInputList:
		fmt.Println("Available Input Plugins:")
		for k, _ := range inputs.Inputs {
			fmt.Printf("  %s\n", k)
		}
		return
	case *fVersion:
		fmt.Printf("grafana-dashboard-exporter %s (git: %s %s)\n", displayVersion(), branch, commit)
		return
	case *fSampleConfig:
		config.PrintSampleConfig(
			inputFilters,
			outputFilters,
		)
		return
	case *fUsage != "":
		err := config.PrintInputConfig(*fUsage)
		err2 := config.PrintOutputConfig(*fUsage)
		if err != nil && err2 != nil {
			log.Fatalf("E! %s and %s", err, err2)
		}
		return
	}
}

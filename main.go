package main

import (
	"flag"
	"fmt"
	"os"
)

var fQuiet = flag.Bool("quiet", false,
	"run in quiet mode")
var fTest = flag.Bool("test", false, "gather metrics, print them out, and exit")
var fConfig = flag.String("config", "", "configuration file to load")
var fVersion = flag.Bool("version", false, "display the version")
var fUsage = flag.String("usage", "",
	"print usage for a plugin, ie, 'grafana-dashboard-exporter --usage s3'")

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

	if len(args) > 0 {
		switch args[0] {
		case "version":
			fmt.Printf("Grafana-dashboard-exporter %s (git: %s %s)\n", displayVersion(), branch, commit)
			return
		case "config":
			return
		}
	}

	// switch for flags which just do something and exit immediately
	switch {
	case *fVersion:
		fmt.Printf("grafana-dashboard-exporter %s (git: %s %s)\n", displayVersion(), branch, commit)
		return
	case *fUsage != "":
		usageExit(0)
		return
	}
}

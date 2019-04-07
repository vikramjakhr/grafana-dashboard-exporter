package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vikramjakhr/grafana-dashboard-exporter/agent"
	"github.com/vikramjakhr/grafana-dashboard-exporter/config"
	"github.com/vikramjakhr/grafana-dashboard-exporter/logger"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs"
	_ "github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs/all"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	_ "github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs/all"
	"log"
	"os/signal"
	"syscall"
)

var fDebug = flag.Bool("debug", false,
	"turn on debug logging")
var fQuiet = flag.Bool("quiet", false,
	"run in quiet mode")
var fTest = flag.Bool("test", false, "gather metrics, print them out, and exit")
var fConfig = flag.String("config", "", "configuration file to load")
var fVersion = flag.Bool("version", false, "display the version")
var fUsage = flag.String("usage", "",
	"print usage for a plugin, ie, 'gde --usage s3'")
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
var fConfigDirectory = flag.String("config-directory", "",
	"directory containing additional *.conf files")

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

const usage = `GDE, The plugin-driven agent for exporting grafana dashboards json.

Usage:

  gde [commands|flags]

The commands & flags are:

  config              print out full sample configuration to stdout
  version             print the version to stdout

  --config <file>     configuration file to load
  --test              gather metrics once, print them to stdout, and exit
  --usage             print usage for a plugin, ie, 'gde --usage s3'
  --quiet             run in quiet mode

Examples:

  # generate a gde config file:
  gde config > gde.conf

  # run a single gde collection, outputing json to stdout
  gde --config gde.conf --test

  # run gde with all plugins defined in config file
  gde --config gde.conf
`

func usageExit(rc int) {
	fmt.Println(usage)
	os.Exit(rc)
}

var stop chan struct{}

func reloadLoop(
	stop chan struct{},
	inputFilters []string,
	outputFilters []string,
) {
	reload := make(chan bool, 1)
	reload <- true
	for <-reload {
		reload <- false

		// If no other options are specified, load the config file and run.
		c := config.NewConfig()
		c.OutputFilters = outputFilters
		c.InputFilters = inputFilters
		err := c.LoadConfig(*fConfig)
		if err != nil {
			log.Fatal("E! " + err.Error())
		}

		if *fConfigDirectory != "" {
			err = c.LoadDirectory(*fConfigDirectory)
			if err != nil {
				log.Fatal("E! " + err.Error())
			}
		}
		if !*fTest && len(c.Outputs) == 0 {
			log.Fatalf("E! Error: no outputs found, did you provide a valid config file?")
		}
		if len(c.Inputs) == 0 {
			log.Fatalf("E! Error: no inputs found, did you provide a valid config file?")
		}

		if int64(c.Agent.Interval.Duration) <= 0 {
			log.Fatalf("E! Agent interval must be positive, found %s",
				c.Agent.Interval)
		}

		ag, err := agent.NewAgent(c)
		if err != nil {
			log.Fatal("E! " + err.Error())
		}

		// Setup logging
		logger.SetupLogging(
			ag.Config.Agent.Debug || *fDebug,
			ag.Config.Agent.Quiet || *fQuiet,
			ag.Config.Agent.Logfile,
		)

		if *fTest {
			err = ag.Test()
			if err != nil {
				log.Fatal("E! " + err.Error())
			}
			os.Exit(0)
		}

		err = ag.Connect()
		if err != nil {
			log.Fatal("E! " + err.Error())
		}

		shutdown := make(chan struct{})
		signals := make(chan os.Signal)
		signal.Notify(signals, os.Interrupt, syscall.SIGHUP)
		go func() {
			select {
			case sig := <-signals:
				if sig == os.Interrupt {
					close(shutdown)
				}
				if sig == syscall.SIGHUP {
					log.Printf("I! Reloading GDE config\n")
					<-reload
					reload <- true
					close(shutdown)
				}
			case <-stop:
				close(shutdown)
			}
		}()

		log.Printf("I! Starting GDE %s\n", displayVersion())
		log.Printf("I! Loaded outputs: %s", strings.Join(c.OutputNames(), " "))
		log.Printf("I! Loaded inputs: %s", strings.Join(c.InputNames(), " "))

		if *fPidfile != "" {
			f, err := os.OpenFile(*fPidfile, os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Printf("E! Unable to create pidfile: %s", err)
			} else {
				fmt.Fprintf(f, "%d\n", os.Getpid())

				f.Close()

				defer func() {
					err := os.Remove(*fPidfile)
					if err != nil {
						log.Printf("E! Unable to remove pidfile: %s", err)
					}
				}()
			}
		}

		ag.Run(shutdown)
	}
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
			fmt.Printf("GDE %s (git: %s %s)\n", displayVersion(), branch, commit)
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
		fmt.Printf("gde %s (git: %s %s)\n", displayVersion(), branch, commit)
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

	stop = make(chan struct{})
	reloadLoop(
		stop,
		inputFilters,
		outputFilters,
	)
}

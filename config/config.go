package config

import (
	"time"
	"sort"
	"fmt"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/outputs"
	"strings"
	"github.com/vikramjakhr/grafana-dashboard-exporter/plugins/inputs"
	"regexp"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
	"errors"
)

var (
	// Default input plugins
	inputDefaults = []string{"grafana"}

	// Default output plugins
	outputDefaults = []string{"file"}

	// envVarRe is a regex to find environment variables in the config file
	envVarRe = regexp.MustCompile(`\$\w+`)

	envVarEscaper = strings.NewReplacer(
		`"`, `\"`,
		`\`, `\\`,
	)
)

var header = `# Grafana-dashboard-exporter Configuration
#
# Grafana-dashboard-exporter is entirely plugin driven. All json are gathered from the
# declared inputs, and sent to the declared outputs.
#
# Use 'grafana-dashboard-exporter -config gde.conf -test' to see what json a config
# file would generate.
#
# Environment variables can be used anywhere in this config file, simply prepend
# them with $. For strings the variable must be within quotes (ie, "$STR_VAR"),
# for numbers and booleans they should be plain (ie, $INT_VAR, $BOOL_VAR)

# Configuration for grafana-dashboard-exporter agent
[agent]
  ## Default data collection interval for all inputs
  interval = "10s"
  ## Rounds collection interval to 'interval'
  ## ie, if interval="10s" then always collect on :00, :10, :20, etc.
  round_interval = true

  ## Logging configuration:
  ## Run grafana-dashboard-exporter with debug log messages.
  debug = false
  ## Run grafana-dashboard-exporter in quiet mode (error log messages only).
  quiet = false
  ## Specify the log file name. The empty string means to log to stderr.
  logfile = ""


###############################################################################
#                            OUTPUT PLUGINS                                   #
###############################################################################
`

var inputHeader = `

###############################################################################
#                            INPUT PLUGINS                                    #
###############################################################################
`

type Config struct {
	Agent   *AgentConfig
	Inputs  []*RunningInput
	Outputs []*RunningOutput
}

func NewConfig() *Config {
	c := &Config{
		// Agent defaults:
		Agent: &AgentConfig{
			Interval:      10 * time.Second,
			RoundInterval: true,
		},

		Inputs:  make([]*RunningInput, 0),
		Outputs: make([]*RunningOutput, 0),
	}
	return c
}

type RunningInput struct {
	Input gde.Input
}

// RunningOutput contains the output configuration
type RunningOutput struct {
	Name   string
	Output gde.Output
}

type AgentConfig struct {
	// Interval at which to gather information
	Interval time.Duration

	// RoundInterval rounds collection interval to 'interval'.
	//     ie, if Interval=10s then always collect on :00, :10, :20, etc.
	RoundInterval bool

	// Debug is the option for running in debug mode
	Debug bool

	// Logfile specifies the file to send logs to
	Logfile string

	// Quiet is the option for running in quiet mode
	Quiet bool
}

func PrintSampleConfig(
	inputFilters []string,
	outputFilters []string,
) {
	fmt.Printf(header)

	// print output plugins
	if len(outputFilters) != 0 {
		printFilteredOutputs(outputFilters, false)
	} else {
		printFilteredOutputs(outputDefaults, false)
		// Print non-default outputs, commented
		var pnames []string
		for pname := range outputs.Outputs {
			if !sliceContains(pname, outputDefaults) {
				pnames = append(pnames, pname)
			}
		}
		sort.Strings(pnames)
		printFilteredOutputs(pnames, true)
	}

	// print input plugins
	fmt.Printf(inputHeader)
	if len(inputFilters) != 0 {
		printFilteredInputs(inputFilters, false)
	} else {
		printFilteredInputs(inputDefaults, false)
		// Print non-default inputs, commented
		var pnames []string
		for pname := range inputs.Inputs {
			if !sliceContains(pname, inputDefaults) {
				pnames = append(pnames, pname)
			}
		}
		sort.Strings(pnames)
		printFilteredInputs(pnames, true)
	}
}

func printFilteredInputs(inputFilters []string, commented bool) {
	// Filter inputs
	var pnames []string
	for pname := range inputs.Inputs {
		if sliceContains(pname, inputFilters) {
			pnames = append(pnames, pname)
		}
	}
	sort.Strings(pnames)

	// Print Inputs
	for _, pname := range pnames {
		creator := inputs.Inputs[pname]
		input := creator()
		printConfig(pname, input, "inputs", commented)
	}
}

func printFilteredOutputs(outputFilters []string, commented bool) {
	// Filter outputs
	var onames []string
	for oname := range outputs.Outputs {
		if sliceContains(oname, outputFilters) {
			onames = append(onames, oname)
		}
	}
	sort.Strings(onames)

	// Print Outputs
	for _, oname := range onames {
		creator := outputs.Outputs[oname]
		output := creator()
		printConfig(oname, output, "outputs", commented)
	}
}

type printer interface {
	Description() string
	SampleConfig() string
}

func printConfig(name string, p printer, op string, commented bool) {
	comment := ""
	if commented {
		comment = "# "
	}
	fmt.Printf("\n%s# %s\n%s[[%s.%s]]", comment, p.Description(), comment,
		op, name)

	config := p.SampleConfig()
	if config == "" {
		fmt.Printf("\n%s  # no configuration\n\n", comment)
	} else {
		lines := strings.Split(config, "\n")
		for i, line := range lines {
			if i == 0 || i == len(lines)-1 {
				fmt.Print("\n")
				continue
			}
			fmt.Print(strings.TrimRight(comment+line, " ") + "\n")
		}
	}
}

func sliceContains(name string, list []string) bool {
	for _, b := range list {
		if b == name {
			return true
		}
	}
	return false
}

// PrintInputConfig prints the config usage of a single input.
func PrintInputConfig(name string) error {
	if creator, ok := inputs.Inputs[name]; ok {
		printConfig(name, creator(), "inputs", false)
	} else {
		return errors.New(fmt.Sprintf("Input %s not found", name))
	}
	return nil
}

// PrintOutputConfig prints the config usage of a single output.
func PrintOutputConfig(name string) error {
	if creator, ok := outputs.Outputs[name]; ok {
		printConfig(name, creator(), "outputs", false)
	} else {
		return errors.New(fmt.Sprintf("Output %s not found", name))
	}
	return nil
}

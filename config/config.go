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
	"os"
	"log"
	"io/ioutil"
	"bytes"
	"github.com/influxdata/toml/ast"
	"github.com/influxdata/toml"
	"path/filepath"
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
	Agent         *AgentConfig
	InputFilters  []string
	OutputFilters []string

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

		Inputs:        make([]*RunningInput, 0),
		Outputs:       make([]*RunningOutput, 0),
		InputFilters:  make([]string, 0),
		OutputFilters: make([]string, 0),
	}
	return c
}

// InputConfig containing a name
type InputConfig struct {
	Name     string
	Interval time.Duration
}

func (r *RunningInput) Name() string {
	return "inputs." + r.Config.Name
}

type RunningInput struct {
	Input  gde.Input
	Config *InputConfig
}

// OutputConfig containing name
type OutputConfig struct {
	Name string
}

// RunningOutput contains the output configuration
type RunningOutput struct {
	Name   string
	Output gde.Output
	Config *OutputConfig
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

func getDefaultConfigPath() (string, error) {
	envfile := os.Getenv("GDE_CONFIG_PATH")
	homefile := os.ExpandEnv("${HOME}/.gde/gde.conf")
	etcfile := "/etc/gde/gde.conf"
	for _, path := range []string{envfile, homefile, etcfile} {
		if _, err := os.Stat(path); err == nil {
			log.Printf("I! Using config file: %s", path)
			return path, nil
		}
	}

	// if we got here, we didn't find a file in a default location
	return "", fmt.Errorf("No config file specified, and could not find one"+
		" in $GDE_CONFIG_PATH, %s, or %s", homefile, etcfile)
}

// LoadConfig loads the given config file and applies it to c
func (c *Config) LoadConfig(path string) error {
	var err error
	if path == "" {
		if path, err = getDefaultConfigPath(); err != nil {
			return err
		}
	}
	tbl, err := parseFile(path)
	if err != nil {
		return fmt.Errorf("Error parsing %s, %s", path, err)
	}

	// Parse agent table:
	if val, ok := tbl.Fields["agent"]; ok {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("%s: invalid configuration", path)
		}
		if err = toml.UnmarshalTable(subTable, c.Agent); err != nil {
			log.Printf("E! Could not parse [agent] config\n")
			return fmt.Errorf("Error parsing %s, %s", path, err)
		}
	}

	// Parse all the rest of the plugins:
	for name, val := range tbl.Fields {
		subTable, ok := val.(*ast.Table)
		if !ok {
			return fmt.Errorf("%s: invalid configuration", path)
		}

		switch name {
		case "agent":
		case "outputs":
			for pluginName, pluginVal := range subTable.Fields {
				switch pluginSubTable := pluginVal.(type) {
				case []*ast.Table:
					for _, t := range pluginSubTable {
						if err = c.addOutput(pluginName, t); err != nil {
							return fmt.Errorf("Error parsing %s, %s", path, err)
						}
					}
				default:
					return fmt.Errorf("Unsupported config format: %s, file %s",
						pluginName, path)
				}
			}
		case "inputs":
			for pluginName, pluginVal := range subTable.Fields {
				switch pluginSubTable := pluginVal.(type) {
				case []*ast.Table:
					for _, t := range pluginSubTable {
						if err = c.addInput(pluginName, t); err != nil {
							return fmt.Errorf("Error parsing %s, %s", path, err)
						}
					}
				default:
					return fmt.Errorf("Unsupported config format: %s, file %s",
						pluginName, path)
				}
			}
			// Assume it's an input input for legacy config file support if no other
			// identifiers are present
		default:
			if err = c.addInput(name, subTable); err != nil {
				return fmt.Errorf("Error parsing %s, %s", path, err)
			}
		}
	}
	return nil
}

// parseFile loads a TOML configuration from a provided path and
// returns the AST produced from the TOML parser. When loading the file, it
// will find environment variables and replace them.
func parseFile(fpath string) (*ast.Table, error) {
	contents, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	env_vars := envVarRe.FindAll(contents, -1)
	for _, env_var := range env_vars {
		env_val, ok := os.LookupEnv(strings.TrimPrefix(string(env_var), "$"))
		if ok {
			env_val = escapeEnv(env_val)
			contents = bytes.Replace(contents, env_var, []byte(env_val), 1)
		}
	}

	return toml.Parse(contents)
}

// escapeEnv escapes a value for inserting into a TOML string.
func escapeEnv(value string) string {
	return envVarEscaper.Replace(value)
}

func (c *Config) addOutput(name string, table *ast.Table) error {
	if len(c.OutputFilters) > 0 && !sliceContains(name, c.OutputFilters) {
		return nil
	}
	creator, ok := outputs.Outputs[name]
	if !ok {
		return fmt.Errorf("Undefined but requested output: %s", name)
	}
	output := creator()

	if err := toml.UnmarshalTable(table, output); err != nil {
		return err
	}
	ro := &RunningOutput{
		Name:   name,
		Output: output,
	}

	c.Outputs = append(c.Outputs, ro)
	return nil
}

func (c *Config) addInput(name string, table *ast.Table) error {
	if len(c.InputFilters) > 0 && !sliceContains(name, c.InputFilters) {
		return nil
	}

	creator, ok := inputs.Inputs[name]
	if !ok {
		return fmt.Errorf("Undefined but requested input: %s", name)
	}
	input := creator()

	if err := toml.UnmarshalTable(table, input); err != nil {
		return err
	}

	rp := &RunningInput{
		Input: input,
	}
	c.Inputs = append(c.Inputs, rp)
	return nil
}

func (c *Config) LoadDirectory(path string) error {
	walkfn := func(thispath string, info os.FileInfo, _ error) error {
		if info == nil {
			log.Printf("W! Telegraf is not permitted to read %s", thispath)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if len(name) < 6 || name[len(name)-5:] != ".conf" {
			return nil
		}
		err := c.LoadConfig(thispath)
		if err != nil {
			return err
		}
		return nil
	}
	return filepath.Walk(path, walkfn)
}

// Inputs returns a list of strings of the configured inputs.
func (c *Config) InputNames() []string {
	var name []string
	for _, input := range c.Inputs {
		name = append(name, input.Name())
	}
	return name
}

// Outputs returns a list of strings of the configured outputs.
func (c *Config) OutputNames() []string {
	var name []string
	for _, output := range c.Outputs {
		name = append(name, output.Name)
	}
	return name
}

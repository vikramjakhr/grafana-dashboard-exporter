package agent

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/vikramjakhr/grafana-dashboard-exporter/config"
	"github.com/vikramjakhr/grafana-dashboard-exporter"
)

// Agent runs GDE and collects data based on the given config
type Agent struct {
	Config *config.Config
}

// NewAgent returns an Agent struct based off the given Config
func NewAgent(config *config.Config) (*Agent, error) {
	a := &Agent{
		Config: config,
	}
	return a, nil
}

// Connect connects to all configured outputs
func (a *Agent) Connect() error {
	for _, o := range a.Config.Outputs {

		log.Printf("D! Attempting connection to output: %s\n", o.Name)
		err := o.Output.Connect()
		if err != nil {
			log.Printf("E! Failed to connect to output %s,"+
				"error was '%s' \n", o.Name, err)
			if err != nil {
				return err
			}
		}
		log.Printf("D! Successfully connected to output: %s\n", o.Name)
	}
	return nil
}

func panicRecover(input *config.RunningInput) {
	if err := recover(); err != nil {
		trace := make([]byte, 2048)
		runtime.Stack(trace, true)
		log.Printf("E! FATAL: Input [%s] panicked: %s, Stack:\n%s\n",
			input.Name(), err, trace)
		log.Println("E! PLEASE REPORT THIS PANIC ON GITHUB with " +
			"stack trace, configuration, and OS information: " +
			"https://github.com/vikramjakhr/grafana-dashboard-exporter/issues/new")
	}
}

// gatherer runs the inputs that have been configured with their own
// reporting interval.
func (a *Agent) gatherer(
	shutdown chan struct{},
	input *config.RunningInput,
	interval time.Duration,
	metricC chan gde.Metric,
) {
	defer panicRecover(input)

	acc := NewAccumulator(input, metricC)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		gatherWithTimeout(shutdown, input, acc, interval)

		select {
		case <-shutdown:
			return
		case <-ticker.C:
			continue
		}
	}
}

// gatherWithTimeout gathers from the given input, with the given timeout.
//   when the given timeout is reached, gatherWithTimeout logs an error message
//   but continues waiting for it to return. This is to avoid leaving behind
//   hung processes, and to prevent re-calling the same hung process over and
//   over.
func gatherWithTimeout(
	shutdown chan struct{},
	input *config.RunningInput,
	acc *accumulator,
	timeout time.Duration,
) {
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	done := make(chan error)
	go func() {
		done <- input.Input.Process(acc)
	}()

	for {
		select {
		case err := <-done:
			if err != nil {
				acc.AddError(err)
			}
			return
		case <-ticker.C:
			err := fmt.Errorf("took longer to collect than collection interval (%s)",
				timeout)
			acc.AddError(err)
			continue
		case <-shutdown:
			return
		}
	}
}

// Test verifies that we can 'Gather' from all inputs with their configured
// Config struct
func (a *Agent) Test() error {
	shutdown := make(chan struct{})
	defer close(shutdown)
	metricC := make(chan gde.Metric)

	// dummy receiver for the point channel
	go func() {
		for {
			select {
			case <-metricC:
				// do nothing
			case <-shutdown:
				return
			}
		}
	}()

	for _, input := range a.Config.Inputs {

		acc := NewAccumulator(input, metricC)

		fmt.Printf("* Plugin: %s, Collection 1\n", input.Name())
		if input.Config.Interval != 0 {
			fmt.Printf("* Internal: %s\n", input.Config.Interval)
		}

		if err := input.Input.Process(acc); err != nil {
			return err
		}
	}
	return nil
}

// flush writes a list of metrics to all configured outputs
func (a *Agent) flush(metric gde.Metric) {
	var wg sync.WaitGroup

	wg.Add(len(a.Config.Outputs))
	for _, o := range a.Config.Outputs {
		go func(output *config.RunningOutput) {
			defer wg.Done()
			err := output.Output.Write(metric)
			if err != nil {
				log.Printf("E! Error writing to output [%s]: %s\n",
					output.Name, err.Error())
			}
		}(o)
	}

	wg.Wait()
}

func (a *Agent) flusher(shutdown chan struct{}, metricC chan gde.Metric) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-shutdown:
				return
			case m := <-metricC:
				for _, o := range a.Config.Outputs {
					o.Output.Write(m)
				}
			}
		}
	}()
	return nil
}

// Run runs the agent daemon, gathering every Interval
func (a *Agent) Run(shutdown chan struct{}) error {
	var wg sync.WaitGroup

	log.Printf("I! Agent Config: Interval:%s, Quiet:%#v \n",
		a.Config.Agent.Interval, a.Config.Agent.Quiet)

	// channel shared between all input threads for accumulating metrics
	metricC := make(chan gde.Metric, 100)

	// Round collection to nearest interval by sleeping
	if a.Config.Agent.RoundInterval {
		i := int64(a.Config.Agent.Interval.Duration)
		time.Sleep(time.Duration(i - (time.Now().UnixNano() % i)))
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.flusher(shutdown, metricC); err != nil {
			log.Printf("E! Flusher routine failed, exiting: %s\n", err.Error())
			close(shutdown)
		}
	}()

	wg.Add(len(a.Config.Inputs))
	for _, input := range a.Config.Inputs {
		interval := a.Config.Agent.Interval.Duration
		// overwrite global interval if this plugin has it's own.
		if input.Config.Interval != 0 {
			interval = input.Config.Interval
		}
		go func(in *config.RunningInput, interv time.Duration) {
			defer wg.Done()
			a.gatherer(shutdown, in, interv, metricC)
		}(input, interval)
	}

	wg.Wait()
	return nil
}

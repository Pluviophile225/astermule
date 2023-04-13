package app

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kasterism/astermule/cmd/app/options"
	"github.com/kasterism/astermule/pkg/dag"
	"github.com/kasterism/astermule/pkg/handlers"
	"github.com/kasterism/astermule/pkg/parser"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ActionMap is used globally to store the list of function entry parameters
var ActionMap map[string]map[string]string

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
var onlyOneSignalHandler = make(chan struct{})
var shutdownHandler chan os.Signal

func NewRootCommand() *cobra.Command {
	opts := options.NewOptions()
	// rootCmd represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:   "astermule",
		Short: "astermule project executor",
		Long:  `Collecting microservice's data according to the directed acyclic graph given by the command line parameter`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCommand(cmd, opts)
		},
	}

	opts.Parse(rootCmd)

	return rootCmd
}

func runCommand(cmd *cobra.Command, opts *options.Options) error {
	return Run(setSignal(), setLogger(), opts)
}

func Run(ctx context.Context, logger *logrus.Logger, opts *options.Options) error {
	graph := dag.NewDAG()
	err := json.Unmarshal([]byte(opts.DagStr), graph)
	if err != nil {
		logger.Fatalln("The dag is not canonical and cannot be resolved")
		return err
	}

	err = graph.Preflight()
	if err != nil {
		logger.Fatalln("Preflight errors:", err)
		return err
	}
	ActionMap = make(map[string]map[string]string)
	p := parser.NewSimpleParser(ActionMap)
	controlPlane, ExitName := p.Parse(graph)
	err = handlers.StartServer(&controlPlane, opts.Address, opts.Port, opts.Target, opts.EntryParam, ActionMap, ExitName)
	if err != nil {
		return err
	}

	fmt.Println(graph)
	return nil
}

func setLogger() *logrus.Logger {
	const logKey = "package"

	logger := logrus.New()

	handlers.SetLogger(logger.WithField(logKey, "handler"))
	dag.SetLogger(logger.WithField(logKey, "dag"))
	parser.SetLogger(logger.WithField(logKey, "parser"))

	return logger
}

func setSignal() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		stopCh := setupSignalHandler()
		<-stopCh
		cancel()
	}()
	return ctx
}

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
// Only one of SetupSignalContext and SetupSignalHandler should be called, and only can
// be called once.
func setupSignalHandler() <-chan struct{} {
	return setupSignalContext().Done()
}

// SetupSignalContext is same as SetupSignalHandler, but a context.Context is returned.
// Only one of SetupSignalContext and SetupSignalHandler should be called, and only can
// be called once.
func setupSignalContext() context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	shutdownHandler = make(chan os.Signal, 2)

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(shutdownHandler, shutdownSignals...)
	go func() {
		<-shutdownHandler
		cancel()
		<-shutdownHandler
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

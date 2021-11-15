package cmd

import (
	"fmt"
	"os"

	"github.com/fairyhunter13/envcompact/internal/app"
	"github.com/fairyhunter13/envcompact/internal/customlog"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	name    = `envcompact`
	version = `v0.2.0`
)

var (
	option  = new(app.Option)
	rootCmd = &cobra.Command{
		Use:   "envcompact",
		Short: "Envcompact is a tool to compact multi-line environment files to single-line environment files.",
		Long: `Envcompact is a tool to compact multi-line environment files to single-line.
Envcompact makes parsing multi-line environments safer since the standard environment file only supports single-line.
Multi-line environment uses double quote (") and single quote (') as the start of multi-line value
and closes the multi-line value with the respective quote.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := zapcore.ErrorLevel
			if option.Verbose {
				level = zapcore.DebugLevel
			}

			core := zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				os.Stderr,
				zap.NewAtomicLevelAt(level),
			)
			options := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(1)}
			logEngine := zap.New(core, options...)
			if option.Silent {
				logEngine = zap.NewNop()
			}

			customlog.Set(logEngine)
		},
		Run: func(cmd *cobra.Command, args []string) {
			if option.PrintVersion {
				fmt.Printf("%s: %s\r\n", name, version)
				return
			}
			application := app.New(
				app.WithInputPath(option.Input),
				app.WithVerbosity(option.Verbose, option.Silent),
			)
			if err := application.Init(); err != nil {
				customlog.Get().Fatal(
					"Error in initializing the application.",
					zap.Error(err),
				)
				return
			}
			defer application.Close()

			if err := application.Run(); err != nil {
				customlog.Get().Fatal(
					"Error running the application.",
					zap.Error(err),
				)
			}
		},
	}
)

func init() {
	// Root flags only
	rootCmd.Flags().StringVarP(&option.Input, "input", "i", "", "input file to read from (default os.Stdin)")
	rootCmd.Flags().BoolVarP(&option.PrintVersion, "version", "V", false, "print version then exit")

	// All child flags
	rootCmd.PersistentFlags().BoolVarP(&option.Verbose, "verbose", "v", false, "verbosity of logging")
	rootCmd.PersistentFlags().BoolVarP(&option.Silent, "silent", "s", false, "disable logging based on this value")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		customlog.Get().Fatal(
			"Error in executing the root command.",
			zap.Error(err),
		)
	}
}

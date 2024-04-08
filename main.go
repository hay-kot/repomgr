package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/hay-kot/repomgr/app/commands"
)

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

func main() {
	ctrl := &commands.Controller{
		Flags: &commands.Flags{},
	}
	app := &cli.App{
		Name:    "Repo Manager",
		Usage:   "Repository Management TUI/CLI for working with Github Projects",
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level (debug, info, warn, error, fatal, panic)",
				Value:       "debug",
				Destination: &ctrl.Flags.LogLevel,
				EnvVars:     []string{"REPOMGR_LOG_LEVEL"},
			},
			&cli.PathFlag{
				Name:    "log-file",
				Usage:   "log file",
				Value:   "repomgr.log",
				EnvVars: []string{"REPOMGR_LOG_FILE"},
			},
		},
		Before: func(ctx *cli.Context) error {
			var writer io.Writer

			logFile := ctx.Path("log-file")
			if logFile == "" {
				writer = io.Discard
			} else {
				f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
				writer = f
			}

			log.Logger = log.Output(zerolog.ConsoleWriter{Out: writer})

			level, err := zerolog.ParseLevel(ctrl.Flags.LogLevel)
			if err != nil {
				return err
			}

			zerolog.SetGlobalLevel(level)
			return nil
		},
		Action: func(ctx *cli.Context) error {
			count := 0
			for {
				log.Info().
					Int("count", count).
					Msg("Hello World")
				time.Sleep(200 * time.Millisecond)
				count++
			}
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run Repo Manager")
	}
}

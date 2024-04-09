package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/hay-kot/repomgr/app/commands"
	"github.com/hay-kot/repomgr/app/core/config"
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
	appctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg := &config.Config{}

	app := &cli.App{
		Name:    "Repo Manager",
		Usage:   "Repository Management TUI/CLI for working with Github Projects",
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level (debug, info, warn, error, fatal, panic)",
				Value:   "debug",
				EnvVars: []string{"REPOMGR_LOG_LEVEL"},
			},
			&cli.PathFlag{
				Name:    "log-file",
				Usage:   "log file",
				Value:   "",
				EnvVars: []string{"REPOMGR_LOG_FILE"},
			},
			&cli.PathFlag{
				Name:    "config",
				Usage:   "config file",
				Value:   "",
				EnvVars: []string{"REPOMGR_CONFIG"},
				Action: func(ctx *cli.Context, p cli.Path) error {
					f, err := os.Open(p)
					if err != nil {
						return fmt.Errorf("failed to open config file: %w", err)
					}

					defer f.Close()

					cfg, err = config.New(f)
					if err != nil {
						return err
					}

					return nil
				},
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

			// TODO: remove color from logs in prod, but keep it in dev
			// for nice tail output
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: writer})

			level, err := zerolog.ParseLevel(ctx.String("log-level"))
			if err != nil {
				return err
			}

			zerolog.SetGlobalLevel(level)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "cache",
				Usage: "cache controls for the database",
				Action: func(ctx *cli.Context) error {
					ctrl := commands.NewController(cfg)
					return ctrl.Cache(appctx)
				},
			},
			{
				Name:   "dump-config",
				Hidden: true,
				Action: func(ctx *cli.Context) error {
					cfgstr, err := cfg.Dump()
					if err != nil {
						return err
					}

					fmt.Println(cfgstr)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run Repo Manager")
	}
}

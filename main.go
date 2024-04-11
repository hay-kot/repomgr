package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/hay-kot/repomgr/app/commands"
	"github.com/hay-kot/repomgr/app/commands/ui"
	"github.com/hay-kot/repomgr/app/console"
	"github.com/hay-kot/repomgr/app/core/config"
	"github.com/hay-kot/repomgr/app/core/services"
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

	cons := console.NewConsole(os.Stdout, true)

	app := &cli.App{
		Name:    "Repo Manager",
		Usage:   "Repository Management TUI/CLI for working with Github Projects",
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "config file",
				Value:   config.ExpandPath("", "~/.config/repomgr/config.toml"),
				EnvVars: []string{"REPOMGR_CONFIG"},
			},
		},
		Before: func(ctx *cli.Context) error {
			p := ctx.String("config")
			f, err := os.Open(p)
			if err != nil {
				return fmt.Errorf("failed to open config file: %w", err)
			}

			defer f.Close()

			absolutePath, err := filepath.Abs(p)
			if err != nil {
				return err
			}

			cfg, err = config.New(absolutePath, f)
			if err != nil {
				return err
			}

			err = cfg.PrepareDirectories()
			if err != nil {
				return err
			}

			var writer io.Writer

			logFile := cfg.Logs.File
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
			if cfg.Logs.Format == "text" {
				log.Logger = log.Output(zerolog.ConsoleWriter{
					Out:     writer,
					NoColor: !cfg.Logs.Color,
				})
			} else if cfg.Logs.Format == "json" {
				log.Logger = log.Output(writer)
			}

			zerolog.SetGlobalLevel(cfg.Logs.Level)
			log.Debug().Str("config", absolutePath).Msg("loaded config")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "cache",
				Usage: "cache controls for the database",
				Action: func(ctx *cli.Context) error {
					sqldb, err := sql.Open("sqlite", cfg.Database.DNS())
					if err != nil {
						return err
					}

					defer sqldb.Close()
					service, err := services.NewRepositoryService(sqldb)
					if err != nil {
						return err
					}

					ctrl := commands.NewController(cfg, service)
					return ctrl.Cache(appctx)
				},
			},
			{
				Name:  "search",
				Usage: "search for repositories",
				Action: func(ctx *cli.Context) error {
					sqldb, err := sql.Open("sqlite", cfg.Database.DNS())
					if err != nil {
						return err
					}

					defer sqldb.Close()
					service, err := services.NewRepositoryService(sqldb)
					if err != nil {
						return err
					}

					ctrl := commands.NewController(cfg, service)
					r, err := ctrl.Search(appctx)
					if err != nil {
						return err
					}

					fmt.Println(r.DisplayName())
					return nil
				},
			},
			{
				Name:   "dev",
				Hidden: true,
				Subcommands: []*cli.Command{
					{
						Name:  "config",
						Usage: "dumps the config to console with default supsutitions",
						Action: func(ctx *cli.Context) error {
							cfgstr, err := cfg.Dump()
							if err != nil {
								return err
							}

							fmt.Println(cfgstr)
							return nil
						},
					},
					{
						Name:  "spinner",
						Usage: "test spinner",
						Action: func(ctx *cli.Context) error {
							return ui.NewSpinnerFunc("loading...", func(ch chan<- string) error {
								for i := 0; i < 10; i++ {
									ch <- fmt.Sprintf("loading... %d", i)
									time.Sleep(300 * time.Millisecond)
								}

								return nil
							})
						},
					},
					{
						Name:  "error",
						Usage: "test/dump console outputs",
						Action: func(ctx *cli.Context) error {
							return errors.New("failed to run config file")
						},
					},
					{
						Name:  "console",
						Usage: "test/dump console outputs",
						Action: func(ctx *cli.Context) error {
							cons.UnknownError("An unexpected error occurred", fmt.Errorf("this is an error"))
							cons.LineBreak()
							cons.List("List of Items Title", []console.ListItem{
								{StatusOk: true, Status: "Item 1 (Success) "},
								{StatusOk: false, Status: "Item 2 (Error) "},
							})
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		cons.UnknownError("An unexpected error occurred", err)
		os.Exit(1)
	}
}

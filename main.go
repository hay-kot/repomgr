package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/hay-kot/repomgr/app/commands"
	"github.com/hay-kot/repomgr/app/debugger"
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
	streamer := debugger.NewStreamer()

	app := &cli.App{
		Name:    "Repo Manager",
		Usage:   "Repository Management TUI/CLI for working with Github Projects",
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cwd",
				Usage:       "current working directory",
				Value:       ".",
				Destination: &ctrl.Flags.WorkingDirectory,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level (debug, info, warn, error, fatal, panic)",
				Value:       "debug",
				Destination: &ctrl.Flags.LogLevel,
			},
		},
		Before: func(ctx *cli.Context) error {
			writer := io.MultiWriter(streamer)
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: writer})

			level, err := zerolog.ParseLevel(ctrl.Flags.LogLevel)
			if err != nil {
				return err
			}

			zerolog.SetGlobalLevel(level)

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "hello",
				Usage: "Says hello world",
				Action: func(ctx *cli.Context) error {
					go func() {
						err := streamer.Start("8080")
						if err != nil {
              panic(err)  
						}
					}()

					count := 0
					for {
						log.Info().
							Int("count", count).
							Msg("Hello World")
						time.Sleep(200 * time.Millisecond)
						count++
					}
				},
			},
			{
				Name:  "attach",
				Usage: "Attach to log endpoint for debugging",
				Action: func(ctx *cli.Context) error {
					u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/logs"}

					c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
					if err != nil {
						return err
					}

					defer c.Close()

					for {
						_, msg, err := c.ReadMessage()
						if err != nil {
							return err
						}
						fmt.Print(string(msg))
					}
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run Repo Manager")
	}
}

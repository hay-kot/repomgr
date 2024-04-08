// Package commands contains the CLI commands for the application
package commands

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type Controller struct {
	Flags *Flags
}

func (c *Controller) HelloWorld(ctx *cli.Context) error {
	fmt.Println("Hello World!")
	return nil
}

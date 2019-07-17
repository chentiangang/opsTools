package com

import (
	"fmt"

	"github.com/urfave/cli"
)

func String(s string, c *cli.Context) (value string, err error) {
	if c.String(s) == "" {
		err := fmt.Sprintf("ERROR: %s cannot is empty, use \"--%s value\"", s, s)
		return "", cli.NewExitError(err, 1)
	}
	return c.String(s), nil
}

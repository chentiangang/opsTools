package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/urfave/cli"
)

const (
	MAC_BUILD     = "go build"
	LINUX_BUILD   = "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build"
	WINDOWS_BUILD = "CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build"
)

func main() {
	app := cli.NewApp()

	app.Name = "Mac对linux,windows的编译工具"
	app.Commands = []cli.Command{
		{
			Name:   "mac",
			Usage:  "Building Mac applications",
			Action: build,
		},
		{
			Name:   "linux",
			Usage:  "Building Linux applications",
			Action: build,
		},
		{
			Name:   "win",
			Usage:  "Building Windows applications",
			Action: build,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func build(ctx *cli.Context) error {
	var cmd string
	switch ctx.Command.Name {
	case "", "mac":
		cmd = MAC_BUILD
	case "linux":
		cmd = LINUX_BUILD
	case "windows", "win":
		cmd = WINDOWS_BUILD
	default:
		return cli.NewExitError("ERROR: unknown command", 1)
	}

	if ctx.Args().First() != "" {
		cmd = cmd + " " + ctx.Args().First()
	}

	err := execmd(cmd, ctx)
	if err != nil {
		return err
	}
	return nil
}

func execmd(cmd string, ctx *cli.Context) error {
	cm := exec.Command("/bin/bash", "-c", cmd)
	cm.Stderr = os.Stderr
	cm.Stdout = os.Stdout

	err := cm.Start()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	err = cm.Wait()
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return nil
}


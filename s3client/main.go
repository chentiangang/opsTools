package main

import (
	"log"
	"os"
	"project/http-deploy/butterfly/command/s3client/get"
	"project/http-deploy/butterfly/command/s3client/list"
	"project/http-deploy/butterfly/command/s3client/put"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "AWS S3 Client Tools; get, put, list."
	app.Commands = []cli.Command{

		{
			Name:  "put",
			Usage: "Upload file to AWS S3.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "src,s",
					Usage: "指定源文件路径",
				},

				cli.StringFlag{
					Name:  "dest,d",
					Usage: "指定目标文件位置",
				},
				cli.StringFlag{
					Name:  "bucket,b",
					Usage: "指定bucket",
				},
				cli.StringFlag{
					Name:  "region",
					Usage: "指定 AWS S3所在区域",
					Value: "ap-northeast-1",
				},
				cli.BoolFlag{
					Name:  "verbose,v",
					Usage: "显示详细日志输出",
				},
			},
			Action: put.Put,
		},
		{
			Name:  "get",
			Usage: "Download file from AWS S3.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file,f",
					Usage: "指定要下载的文件路径",
				},

				cli.StringFlag{
					Name:  "bucket,b",
					Usage: "指定bucket",
					Value: "55-code-store",
				},
				cli.StringFlag{
					Name:  "region",
					Usage: "指定AWS S3所在区域",
					Value: "ap-northeast-1",
				},
			},
			Action: get.Get,
		},
		{
			Name:   "list",
			Usage:  "List Object from AWS S3.",
			Action: list.List,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "bucket,b",
					Usage: "指定bucket",
					Value: "55-code-store",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

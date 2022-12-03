package cmd

import (
	"fmt"
	"github.com/Xacnio/tweetcapture-go/pkg/tweetcapture"
	"github.com/urfave/cli"
	"os"
)

var app *cli.App

func RunCmd() {
	app = cli.NewApp()
	app.Name = "tweetcapture-go"
	app.Usage = "simple tweet screenshot tool"
	app.UsageText = app.Name + " [options] [Tweet URL]"
	app.Version = "v1.0.2"

	registerArgs(app)
	registerAction(app)
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func registerArgs(app *cli.App) {
	app.Flags = []cli.Flag{
		cli.UintFlag{Name: "nightmode, n", Value: 0, Usage: "0-Light / 1-Dim / 2-Lights out"},
		cli.BoolFlag{Name: "seleniumdebug, sd", Hidden: true},
		cli.StringFlag{Name: "output, o", Value: "", Usage: "File path to save screenshot"},
	}
}

func registerAction(app *cli.App) {
	app.Action = func(c *cli.Context) error {
		url := c.Args().Get(0)
		if len(url) < 1 {
			return cli.ShowAppHelp(c)
		}
		tss := tweetcapture.NewTweetScreenshot(url)
		tss.Opts.NightMode = uint8(c.Uint("nightmode"))
		tss.Opts.SeleniumDebug = c.Bool("seleniumdebug")
		tss.SavePath = c.String("output")
		err := tss.Screenshot()
		if err != nil {
			return cli.NewExitError(err, 26)
		} else {
			fmt.Printf("Tweet: %s\nScreenshot Saved: %s\n", tss.URL, tss.SavePath)
		}
		return nil
	}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/skip2/go-qrcode"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetLevel(log.InfoLevel)
	log.SetReportTimestamp(false)
	log.SetTimeFormat("2006.01.02 15:04:05")
	log.DebugLevelStyle = log.DebugLevelStyle.MaxWidth(5)
	log.FatalLevelStyle = log.FatalLevelStyle.MaxWidth(5)
}

func main() {
	app := &cli.App{
		Name:      "mqr",
		Usage:     "generate QR codes in terminal",
		ArgsUsage: "[data-to-encode]",
		Version:   versionInfo(),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "enable debug mode",
				Action: func(_ *cli.Context, debug bool) error {
					if debug {
						log.SetLevel(log.DebugLevel)
						log.SetReportTimestamp(true)
					}
					return nil
				},
			},
			&cli.IntFlag{
				Name:    "recovery-level",
				Usage:   fmt.Sprintf("error recovery level [%d-%d]", qrcode.Low, qrcode.Highest),
				Value:   1,
				Aliases: []string{"r"},
				Action: func(_ *cli.Context, level int) error {
					if level < int(qrcode.Low) || level > int(qrcode.Highest) {
						return fmt.Errorf("invalid recovery level, should be in range [%d-%d]",
							qrcode.Low, qrcode.Highest)
					}
					return nil
				},
			},
			&cli.BoolFlag{
				Name:    "big",
				Usage:   "print big",
				Value:   false,
				Aliases: []string{"b"},
			},
			&cli.BoolFlag{
				Name:    "inverted",
				Usage:   "print with inverted colors",
				Value:   false,
				Aliases: []string{"i"},
			},
			&cli.BoolFlag{
				Name:    "tui",
				Usage:   "force TUI",
				Value:   false,
				Aliases: []string{"t"},
			},
		},
		EnableBashCompletion: true,
		HideHelpCommand:      true,
		BashComplete:         cli.DefaultAppComplete,
		Action:               run,
		Authors: []*cli.Author{
			{
				Name:  "Artem Yadelskyi",
				Email: "mymmrac@gmail.com",
			},
		},
		UseShortOptionHandling: true,
		Suggest:                true,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		cancel()
	}()

	if err := app.RunContext(ctx, os.Args); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		log.Fatalf("Error: %s", err)
	}
}

func versionInfo() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	var (
		vcsRevision string
		vcsTime     time.Time
		vcsModified bool
	)

	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcsRevision = setting.Value
		case "vcs.time":
			vcsTime, _ = time.Parse(time.RFC3339, setting.Value)
		case "vcs.modified":
			vcsModified, _ = strconv.ParseBool(setting.Value)
		}
	}

	version := fmt.Sprintf("%s, build with %s", info.Main.Version, info.GoVersion)

	if vcsRevision != "" {
		version += ", revision " + vcsRevision
	}
	if !vcsTime.IsZero() {
		version += ", at " + vcsTime.Local().String()
	}
	if vcsModified {
		version += ", modified"
	}

	return version
}

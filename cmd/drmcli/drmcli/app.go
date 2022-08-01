package drmcli

import (
	"strconv"

	"git.sr.ht/~emersion/go-drm"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

var App = cli.App{
	Name: "drmcli",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "device",
			Usage:    "path to DRI",
			Value:    "/dev/dri/card0",
			Aliases:  []string{"D"},
			Required: true,
		},
		&cli.BoolFlag{
			Name:    "raw",
			Usage:   "don't convert enums to strings if true",
			Value:   false,
			Aliases: []string{"r"},
		},
	},
	Commands: []*cli.Command{
		{
			Name:        "resources",
			Action:      cmdResources,
			Description: "print *drm.ModeCard containing top-level information (mostly IDs)",
		},
		{
			Name: "connector",
			Subcommands: []*cli.Command{
				{
					Name:        "info",
					Action:      cmdConnectorInfo,
					ArgsUsage:   "<connector ID>",
					Description: "print *drm.ModeConnector containing information for a connector",
				},
				{
					Name:        "properties",
					Action:      cmdConnectorProperties,
					ArgsUsage:   "<connector ID>",
					Description: "get all properties of a connector",
				},
			},
		},
		{
			Name:      "property",
			Action:    cmdProperty,
			Usage:     "query information about a property",
			ArgsUsage: "<property ID>",
		},
	},
}

func firstArgID(ctx *cli.Context) (drm.ObjectID, error) {
	return nArgID(ctx, 0)
}

func nArgID(ctx *cli.Context, n int) (drm.ObjectID, error) {
	str := ctx.Args().Get(n)
	if str == "" {
		cli.ShowSubcommandHelp(ctx)
		return 0, cli.Exit("invalid usage", 1)
	}

	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "invalid ID")
	}

	return drm.ObjectID(id), nil

}

package main

import (
	"os"
	"sidecar/pkg"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func main() {
	pkg.InitLogger()

	app := &cli.App{
		Name:  "sidecar",
		Usage: "Easily inject environment variables",
		Commands: []*cli.Command{
			{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "Create new profile",
				Action: func(cCtx *cli.Context) error {
					if err := pkg.CreateProfile(cCtx.Args().First()); err != nil {
						return err
					}

					log.Info().Msgf("> Profile %s created", cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "show",
				Aliases: []string{"s"},
				Usage:   "Show profile",
				Action: func(cCtx *cli.Context) error {
					profile, err := pkg.ShowProfile(cCtx.Args().First())
					if err != nil {
						return err
					}

					pkg.DisplayProfile(cCtx.Args().First(), profile)
					return nil
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List existing profiles",
				Action: func(cCtx *cli.Context) error {
					profiles, err := pkg.ListProfiles()
					if err != nil {
						return err
					}

					pkg.DisplayProfiles(profiles)
					return nil
				},
			},
			{
				Name:      "delete",
				Aliases:   []string{"rm"},
				Usage:     "Delete profile",
				ArgsUsage: "name of the profile",
				Action: func(cCtx *cli.Context) error {
					if err := pkg.DeleteProfile(cCtx.Args().First()); err != nil {
						return err
					}

					log.Info().Msgf("> Profile %s deleted", cCtx.Args().First())
					return nil
				},
			},

			{
				Name:    "add",
				Aliases: []string{},
				Usage:   "Add values to profile",
				Action: func(cCtx *cli.Context) error {
					if err := pkg.AddToProfile(cCtx.Args().First(), cCtx.Args().Tail()...); err != nil {
						return err
					}

					log.Info().Msgf("> Profile %s updated", cCtx.Args().First())
					return nil
				},
			},
			{
				Name:    "remove",
				Aliases: []string{},
				Usage:   "Remove values to profile",
				Action: func(cCtx *cli.Context) error {
					if err := pkg.RemoveFromProfile(cCtx.Args().First(), cCtx.Args().Tail()...); err != nil {
						return err
					}

					log.Info().Msgf("> Profile %s updated", cCtx.Args().First())
					return nil
				},
			},

			{
				Name:    "execute",
				Aliases: []string{"exe"},
				Usage:   "Execute a command with sidecar",
				Action: func(cCtx *cli.Context) error {
					return pkg.Execute(cCtx.Args().First(), cCtx.Args().Get(1))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error().Msg(err.Error())
	}
}

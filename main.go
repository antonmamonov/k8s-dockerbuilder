// Copyright 2022 Anton Mamonov <hi@antonmamonov.com> GNU GENERAL PUBLIC LICENSE
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/antonmamonov/k8s-dockerbuilder/build"
	"github.com/cristalhq/acmd"
)

func main() {

	cmds := []acmd.Command{
		{
			Name:        "now",
			Description: "prints current time",
			ExecFunc: func(ctx context.Context, args []string) error {
				fmt.Printf("now: %s\n", time.Now())
				return nil
			},
		},
		{
			Name:        "build",
			Description: "Build the current git repository and push it to the registry",
			ExecFunc: func(ctx context.Context, args []string) error {

				// first arg should be the docker image destination
				if len(args) == 0 {
					return fmt.Errorf("Please provide the docker image destination as an argument")
				}

				cfg := build.BuildConfig{
					DockerImageDestination: args[0],
				}

				buildError := build.Build(&cfg)

				if buildError != nil {
					return buildError
				}

				return nil
			},
		},
	}

	// all the acmd.Config fields are optional
	r := acmd.RunnerOf(cmds, acmd.Config{
		AppName:        "kubebuild",
		AppDescription: "Build & Push Docker images inside your Kubernetes cluster <hi@antonmamonov.com>",
		Version:        "v0.0.1",
	})

	if err := r.Run(); err != nil {
		r.Exit(err)
	}
}

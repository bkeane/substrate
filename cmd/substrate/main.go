package main

import (
	"context"
	"fmt"
	"os"

	"github.com/bkeane/substrate/cmd/substrate/list"
	"github.com/bkeane/substrate/cmd/substrate/render"

	"github.com/alexflint/go-arg"
	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
	if value, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if level, err := zerolog.ParseLevel(value); err == nil {
			zerolog.SetGlobalLevel(level)
		}
	}
}

type Root struct {
	Render *render.Root `arg:"subcommand:render" help:"render substrate configuration given name and features"`
	List   *list.Root   `arg:"subcommand:list" help:"list available substrates and their features"`
}

func (r *Root) Route(ctx context.Context, awsconfig aws.Config) (*string, error) {
	switch {
	case r.Render != nil:
		return r.Render.Route(ctx, awsconfig)
	case r.List != nil:
		return r.List.Route(ctx, awsconfig)
	}
	return nil, nil
}

func main() {
	ctx := context.Background()

	var root Root
	var output *string
	var err error

	arg.MustParse(&root)

	awsconfig, err := awscfg.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	output, err = root.Route(ctx, awsconfig)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	if output != nil {
		fmt.Println(*output)
	}
}

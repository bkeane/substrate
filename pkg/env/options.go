package env

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	envlib "github.com/caarlos0/env/v11"
)

type Options struct {
	// include org name in paths such as routes.
	OrgPrefixedPaths bool `env:"SUBSTRATE_PREFIX_PATHS_WITH_ORG"`
	// include org name in resource names such as IAM policies.
	OrgPrefixedNames bool `env:"SUBSTRATE_PREFIX_NAMES_WITH_ORG"`
}

func (o *Options) Parse(ctx context.Context, awsconfig aws.Config, envMap map[string]string) error {
	if err := o.parse(envMap); err != nil {
		return err
	}

	if err := o.deduce(); err != nil {
		return err
	}

	if err := o.validate(); err != nil {
		return err
	}

	return nil
}

func (o *Options) parse(envMap map[string]string) error {
	if envMap == nil {
		if err := envlib.Parse(o); err != nil {
			return err
		}
		return nil
	}

	return envlib.ParseWithOptions(o, envlib.Options{
		Environment: envMap,
	})
}

func (o *Options) deduce() error {
	return nil
}

func (o *Options) validate() error {
	return nil
}

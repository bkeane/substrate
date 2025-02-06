package env

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	envlib "github.com/caarlos0/env/v11"
)

type Lambda struct {
	// Lambda region, defaults to caller region.
	Region *string `env:"SUBSTRATE_LAMBDA_REGION"`
}

func (l *Lambda) Parse(ctx context.Context, awsconfig aws.Config, envMap map[string]string) error {
	if err := l.parse(envMap); err != nil {
		return err
	}

	if err := l.deduce(awsconfig); err != nil {
		return err
	}

	if err := l.validate(); err != nil {
		return err
	}

	return nil
}

func (l *Lambda) parse(envMap map[string]string) error {
	if envMap == nil {
		if err := envlib.Parse(l); err != nil {
			return err
		}
		return nil
	}

	return envlib.ParseWithOptions(l, envlib.Options{
		Environment: envMap,
	})
}

func (l *Lambda) deduce(awsconfig aws.Config) error {
	if l.Region == nil {
		l.Region = aws.String(awsconfig.Region)
	}

	return nil
}

func (l *Lambda) validate() error {
	return nil
}

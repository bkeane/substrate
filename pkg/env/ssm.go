package env

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	envlib "github.com/caarlos0/env/v11"
)

type SSM struct {
	// SSM region, defaults to caller region.
	Region *string `env:"SUBSTRATE_SSM_REGION"`
}

func (s *SSM) Parse(ctx context.Context, awsconfig aws.Config, envMap map[string]string) error {
	if err := s.parse(envMap); err != nil {
		return err
	}

	if err := s.deduce(awsconfig); err != nil {
		return err
	}

	if err := s.validate(); err != nil {
		return err
	}

	return nil
}

func (s *SSM) parse(envMap map[string]string) error {
	if envMap == nil {
		if err := envlib.Parse(s); err != nil {
			return err
		}
		return nil
	}

	return envlib.ParseWithOptions(s, envlib.Options{
		Environment: envMap,
	})
}

func (s *SSM) deduce(awsconfig aws.Config) error {
	if s.Region == nil {
		s.Region = aws.String(awsconfig.Region)
	}

	return nil
}

func (s *SSM) validate() error {
	return nil
}

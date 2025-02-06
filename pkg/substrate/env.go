package substrate

import (
	"context"

	"github.com/bkeane/substrate/pkg/env"

	"github.com/aws/aws-sdk-go-v2/aws"
	envlib "github.com/caarlos0/env/v11"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type Substrate struct {
	// The name of the substrate
	Name string `env:"SUBSTRATE_NAME"`
	// This account's name
	Source string `env:"SUBSTRATE_SOURCE"`
	// This account's possible destinations
	Destinations []string `env:"SUBSTRATE_DESTINATIONS"`
	// Available feature flags on the substrate
	Features []string `env:"SUBSTRATE_FEATURES"`

	Account     env.Account
	Options     env.Options
	EventBridge env.EventBridge
	ApiGateway  env.ApiGateway
	Lambda      env.Lambda
	ECR         env.ECR
	SSM         env.SSM
	VPC         env.VPC
	awsconfig   aws.Config `json:"-"`
}

func (s *Substrate) Parse(ctx context.Context, envMap map[string]string) error {
	if err := s.parse(envMap); err != nil {
		return err
	}

	if err := s.deduce(); err != nil {
		return err
	}

	if err := s.validate(); err != nil {
		return err
	}

	return nil
}

func (s *Substrate) parse(envMap map[string]string) error {
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

func (s *Substrate) deduce() error {
	return nil
}

func (s *Substrate) validate() error {
	return v.ValidateStruct(s,
		v.Field(&s.Name, v.Required),
		v.Field(&s.Source, v.Required),
		v.Field(&s.Destinations, v.Required),
		v.Field(&s.Features, v.Required),
	)
}

package env

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	envlib "github.com/caarlos0/env/v11"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type Account struct {
	// This account's ID
	Id *string
	// This account's region
	Region *string
	// This account's name
	Name string `env:"SUBSTRATE_SOURCE"`
}

func (a *Account) Parse(ctx context.Context, awsconfig aws.Config, envMap map[string]string) error {
	if err := a.parse(envMap); err != nil {
		return err
	}

	if err := a.deduce(ctx, awsconfig); err != nil {
		return err
	}

	if err := a.validate(); err != nil {
		return err
	}

	return nil
}

func (a *Account) parse(envMap map[string]string) error {
	if envMap == nil {
		if err := envlib.Parse(a); err != nil {
			return err
		}
		return nil
	}

	return envlib.ParseWithOptions(a, envlib.Options{
		Environment: envMap,
	})
}

func (a *Account) deduce(ctx context.Context, awsconfig aws.Config) error {
	stsc := sts.NewFromConfig(awsconfig)
	caller, err := stsc.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return err
	}

	a.Id = caller.Account
	a.Region = aws.String(awsconfig.Region)

	return nil
}

func (a *Account) validate() error {
	return v.ValidateStruct(a,
		v.Field(&a.Id, v.Required),
		v.Field(&a.Region, v.Required),
		v.Field(&a.Name, v.Required),
	)
}

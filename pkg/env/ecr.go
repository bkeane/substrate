package env

import (
	"context"
	"fmt"

	"github.com/bkeane/substrate/pkg/registry"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	envlib "github.com/caarlos0/env/v11"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type ECR struct {
	awsconfig aws.Config `json:"-"`
	// ECR registry id, defaults to caller account ECR registry.
	Id *string `env:"SUBSTRATE_ECR_REGISTRY_ID"`
	// ECR registry region, defaults to caller region.
	Region *string `env:"SUBSTRATE_ECR_REGISTRY_REGION"`
}

func (e *ECR) Parse(ctx context.Context, awsconfig aws.Config, envMap map[string]string) error {
	if err := e.parse(envMap); err != nil {
		return err
	}

	if err := e.deduce(ctx, awsconfig); err != nil {
		return err
	}

	if err := e.validate(); err != nil {
		return err
	}

	return nil
}

func (e *ECR) parse(envMap map[string]string) error {
	if envMap == nil {
		if err := envlib.Parse(e); err != nil {
			return err
		}
		return nil
	}

	return envlib.ParseWithOptions(e, envlib.Options{
		Environment: envMap,
	})
}

func (e *ECR) deduce(ctx context.Context, awsconfig aws.Config) error {
	if e.Id == nil {
		ecrc := ecr.NewFromConfig(awsconfig)
		ecr, err := ecrc.DescribeRegistry(ctx, &ecr.DescribeRegistryInput{})
		if err != nil {
			return err
		}
		e.Id = ecr.RegistryId
	}

	if e.Region == nil {
		e.Region = aws.String(awsconfig.Region)
	}

	e.awsconfig = awsconfig

	if err := e.validate(); err != nil {
		return err
	}

	return nil
}

func (e *ECR) validate() error {
	return v.ValidateStruct(e,
		v.Field(&e.awsconfig, v.Required), // awsconfig is required for registry client
		v.Field(&e.Id, v.Required),
		v.Field(&e.Region, v.Required),
	)
}

func (e *ECR) RegistryUrl() string {
	return fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", *e.Id, *e.Region)
}

func (e *ECR) RegistryClient(ctx context.Context) (*registry.Registry, error) {
	client, err := registry.Init(ctx, e.awsconfig, e.RegistryUrl())
	if err != nil {
		return nil, fmt.Errorf("failed while initializing registry client: %w", err)
	}

	return client, nil
}

func (e *ECR) FetchByUri(ctx context.Context, imageUri string) (registry.ImagePointer, error) {
	image, err := registry.GetImageFromUri(ctx, e.awsconfig, imageUri)
	if err != nil {
		return registry.ImagePointer{}, fmt.Errorf("failed while fetching %s: %w", imageUri, err)
	}

	return image, nil
}

func (e *ECR) FetchByName(ctx context.Context, repository, reference string) (registry.ImagePointer, error) {
	reg, err := e.RegistryClient(ctx)
	if err != nil {
		return registry.ImagePointer{}, fmt.Errorf("failed while fetching %s:%s: %w", repository, reference, err)
	}

	return reg.GetImageFromName(ctx, repository, reference)
}

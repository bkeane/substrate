package env

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	envlib "github.com/caarlos0/env/v11"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type VPC struct {
	// security group ids for deploying to private vpc.
	SecurityGroupIds []string `env:"SUBSTRATE_VPC_SECURITY_GROUP_IDS"`
	// subnet ids for deploying to private vpc.
	SubnetIds []string `env:"SUBSTRATE_VPC_SUBNET_IDS"`
}

func (c *VPC) Parse(ctx context.Context, awsconfig aws.Config, envMap map[string]string) error {
	if err := c.parse(envMap); err != nil {
		return err
	}

	if err := c.deduce(); err != nil {
		return err
	}

	if err := c.validate(); err != nil {
		return err
	}

	return nil
}

func (c *VPC) parse(envMap map[string]string) error {
	if envMap == nil {
		if err := envlib.Parse(c); err != nil {
			return err
		}
		return nil
	}

	return envlib.ParseWithOptions(c, envlib.Options{
		Environment: envMap,
	})
}

func (c *VPC) deduce() error {
	return nil
}

func (c *VPC) validate() error {
	return v.ValidateStruct(c,
		v.Field(&c.SecurityGroupIds, v.When(c.SubnetIds != nil, v.Required)),
		v.Field(&c.SubnetIds, v.When(c.SecurityGroupIds != nil, v.Required)),
	)
}

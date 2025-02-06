package env

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	envlib "github.com/caarlos0/env/v11"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type ApiGateway struct {
	// API Gateway Id for https based invokes.
	Id *string `env:"SUBSTRATE_APIGATEWAY_ID"`
	// API Gateway region, defaults to caller region.
	Region *string `env:"SUBSTRATE_APIGATEWAY_REGION"`
	// "true": enable, "false": disable, unset: leave in current state.
	Enable *bool `env:"SUBSTRATE_APIGATEWAY_ENABLE"`
	// One of "NONE", "AWS_IAM", "CUSTOM", "JWT".
	AuthType *string `env:"SUBSTRATE_APIGATEWAY_AUTH_TYPE"`
	// Required if auth type is "CUSTOM" or "JWT".
	AuthorizerId *string `env:"SUBSTRATE_APIGATEWAY_AUTHORIZER_ID"`
}

func (a *ApiGateway) Parse(ctx context.Context, awsconfig aws.Config, envMap map[string]string) error {
	if err := a.parse(envMap); err != nil {
		return err
	}

	if err := a.deduce(awsconfig); err != nil {
		return err
	}

	if err := a.validate(); err != nil {
		return err
	}

	return nil
}

func (a *ApiGateway) parse(envMap map[string]string) error {
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

func (a *ApiGateway) deduce(awsconfig aws.Config) error {
	if a.Region == nil {
		a.Region = aws.String(awsconfig.Region)
	}
	return nil
}

func (a *ApiGateway) validate() error {
	return v.ValidateStruct(a,
		v.Field(&a.Id, v.When(a.Enable != nil && *a.Enable, v.Required)),
		v.Field(&a.AuthType, v.When(a.AuthType != nil, v.In("NONE", "AWS_IAM", "CUSTOM", "JWT")), v.When(a.AuthorizerId != nil, v.In("JWT", "CUSTOM"))),
		v.Field(&a.AuthorizerId, v.When(a.AuthType != nil && (*a.AuthType == "CUSTOM" || *a.AuthType == "JWT"), v.Required)),
	)
}

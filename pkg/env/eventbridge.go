package env

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	envlib "github.com/caarlos0/env/v11"
	v "github.com/go-ozzo/ozzo-validation/v4"
)

type EventBridge struct {
	// EventBridge bus name for eventbridge based invokes, defaults to "default".
	BusName string `env:"SUBSTRATE_EVENTBRIDGE_BUS_NAME"`
	// EventBridge region, defaults to caller region.
	Region *string `env:"SUBSTRATE_EVENTBRIDGE_REGION"`
	// "true": enable, "false": disable, unset: leave in current state.
	Enable *bool `env:"SUBSTRATE_EVENTBRIDGE_ENABLE"`
}

func (e *EventBridge) Parse(ctx context.Context, awsconfig aws.Config, envMap map[string]string) error {
	if err := e.parse(envMap); err != nil {
		return err
	}

	if err := e.deduce(awsconfig); err != nil {
		return err
	}

	if err := e.validate(); err != nil {
		return err
	}

	return nil
}

func (e *EventBridge) parse(envMap map[string]string) error {
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

func (e *EventBridge) deduce(awsconfig aws.Config) error {
	if e.BusName == "" {
		e.BusName = "default"
	}

	if e.Region == nil {
		e.Region = aws.String(awsconfig.Region)
	}

	return nil
}

func (e *EventBridge) validate() error {
	return v.ValidateStruct(e,
		v.Field(&e.BusName, v.When(e.Enable != nil && *e.Enable, v.Required)),
	)
}

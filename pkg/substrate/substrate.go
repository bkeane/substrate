package substrate

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/bkeane/substrate/pkg/event"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	substratePath        = "/substrate/"
	substratePathPattern = `^/substrate/(?P<name>[^/]+)(?:/(?P<feature>[^/]+))?$`
)

type SubstrateDetail struct {
	Name         string
	Features     []string
	Destinations []string
}

func Index(ctx context.Context, awsconfig aws.Config) ([]SubstrateDetail, error) {
	client := ssm.NewFromConfig(awsconfig)

	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(substratePath),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}

	var substrates []string
	paginator := ssm.NewGetParametersByPathPaginator(client, input)
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, param := range output.Parameters {
			name, _, err := extract(*param.Name)
			if err != nil {
				return nil, err
			}

			if !slices.Contains(substrates, name) {
				substrates = append(substrates, name)
			}
		}
	}

	var details []SubstrateDetail
	for _, name := range substrates {
		substrate, err := Parse(ctx, awsconfig, name)
		if err != nil {
			return nil, err
		}

		details = append(details, SubstrateDetail{
			Name:         substrate.Name,
			Features:     substrate.Features,
			Destinations: substrate.Destinations,
		})
	}

	return details, nil
}

func Parse(ctx context.Context, awsconfig aws.Config, substrate string, features ...string) (*Substrate, error) {
	s := &Substrate{
		awsconfig: awsconfig,
	}

	env, err := merged(ctx, awsconfig, substrate, features...)
	if err != nil {
		return nil, err
	}

	if err := s.Parse(ctx, env); err != nil {
		return nil, err
	}

	if err := s.Account.Parse(ctx, awsconfig, env); err != nil {
		return nil, err
	}

	if err := s.Options.Parse(ctx, awsconfig, env); err != nil {
		return nil, err
	}

	if err := s.ECR.Parse(ctx, awsconfig, env); err != nil {
		return nil, err
	}

	if err := s.EventBridge.Parse(ctx, awsconfig, env); err != nil {
		return nil, err
	}

	if err := s.ApiGateway.Parse(ctx, awsconfig, env); err != nil {
		return nil, err
	}

	if err := s.Lambda.Parse(ctx, awsconfig, env); err != nil {
		return nil, err
	}

	if err := s.SSM.Parse(ctx, awsconfig, env); err != nil {
		return nil, err
	}

	if err := s.VPC.Parse(ctx, awsconfig, env); err != nil {
		return nil, err
	}

	return s, nil
}

func Rx(ctx context.Context, awsconfig aws.Config, data json.RawMessage) (proto.Message, *Substrate, error) {
	protojopt := protojson.UnmarshalOptions{}
	protoopt := proto.UnmarshalOptions{}

	var outer event.EventBridgeEvent
	if err := protojopt.Unmarshal(data, &outer); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal json to eventbridge event: %w", err)
	}

	header := outer.GetDetail().GetHeader()

	msg, err := anypb.UnmarshalNew(outer.GetDetail().GetBody(), protoopt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal eventbridge event to proto type: %w", err)
	}

	substrate, err := Parse(ctx, awsconfig, header.Substrate, header.Features...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to resolve substrate: %w", err)
	}

	return msg, substrate, nil
}

func (s *Substrate) Tx(ctx context.Context, msg proto.Message, destinations []string, features []string) (*string, error) {
	for _, destination := range destinations {
		if !slices.Contains(s.Destinations, destination) {
			return nil, fmt.Errorf("invalid destination: %s", destination)
		}
	}

	for _, feature := range features {
		if !slices.Contains(s.Features, feature) {
			return nil, fmt.Errorf("invalid feature: %s", feature)
		}
	}

	header := &event.Header{
		Substrate:   s.Name,
		Features:    features,
		Source:      s.Source,
		Destination: destinations,
	}

	event, err := event.ToEventBridgeEvent(msg, header)
	if err != nil {
		return nil, err
	}

	out, err := event.ToString()
	if err != nil {
		return nil, err
	}

	o, _ := os.Stdout.Stat()
	if (o.Mode() & os.ModeCharDevice) != os.ModeCharDevice {
		log.Debug().Msg("detected piped output, not publishing event")
		return out, nil
	}

	return out, event.Publish(ctx, s.awsconfig, s.EventBridge.BusName)
}

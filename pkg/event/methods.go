package event

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

var (
	pjsonm *protojson.MarshalOptions
)

func init() {
	pjsonm = &protojson.MarshalOptions{}
}

func ToEventBridgeEvent[T protoreflect.ProtoMessage](body T, header *Header) (*EventBridgeEvent, error) {
	binaryPath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	payload, err := anypb.New(body)
	if err != nil {
		return nil, err
	}

	return &EventBridgeEvent{
		Source:     filepath.Base(binaryPath),
		DetailType: string(body.ProtoReflect().Descriptor().FullName()),
		Detail: &Transport{
			Header: header,
			Body:   payload,
		},
	}, nil
}

func (e *EventBridgeEvent) ToBytes() ([]byte, error) {
	eventBytes, err := pjsonm.Marshal(e)
	if err != nil {
		return nil, err
	}

	return eventBytes, nil
}

func (e *EventBridgeEvent) ToString() (*string, error) {
	eventBytes, err := e.ToBytes()
	if err != nil {
		return nil, err
	}

	eventString := string(eventBytes)
	return &eventString, nil
}

func (e *EventBridgeEvent) Publish(ctx context.Context, awsconfig aws.Config, busName string) error {
	detailBytes, err := pjsonm.Marshal(e.Detail)
	if err != nil {
		return err
	}

	entry := []types.PutEventsRequestEntry{
		{
			Source:       aws.String(e.Source),
			EventBusName: aws.String(busName),
			DetailType:   aws.String(e.DetailType),
			Detail:       aws.String(string(detailBytes)),
		},
	}

	input := &eventbridge.PutEventsInput{
		Entries: entry,
	}

	ebc := eventbridge.NewFromConfig(awsconfig)
	_, err = ebc.PutEvents(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func FromEventBridgeEvent[T protoreflect.ProtoMessage](event *EventBridgeEvent, into T) error {
	return event.GetDetail().GetBody().UnmarshalTo(into)
}

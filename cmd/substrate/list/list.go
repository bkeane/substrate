package list

import (
	"context"
	"strings"

	"github.com/bkeane/substrate/pkg/substrate"
	"github.com/charmbracelet/lipgloss/table"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type Root struct{}

func (r *Root) Route(ctx context.Context, awsconfig aws.Config) (*string, error) {
	substrates, err := substrate.Index(ctx, awsconfig)
	if err != nil {
		return nil, err
	}

	tbl := table.New()
	tbl.Headers("substrate", "features", "destinations")

	for _, detail := range substrates {
		tbl.Row(detail.Name, strings.Join(detail.Features, ", "), strings.Join(detail.Destinations, ", "))
	}

	result := tbl.Render()
	return &result, nil
}

package substrate

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	dotenvlib "github.com/joho/godotenv"
)

func merged(ctx context.Context, awsconfig aws.Config, substrateName string, features ...string) (map[string]string, error) {
	env, err := fetch(ctx, awsconfig, substratePath+substrateName)
	if err != nil {
		return nil, err
	}

	for _, feature := range features {
		featureEnv, err := fetch(ctx, awsconfig, substratePath+substrateName+"/"+feature)
		if err != nil {
			return nil, err
		}
		env = merge(env, featureEnv)
	}

	return env, nil
}

func fetch(ctx context.Context, awsconfig aws.Config, path string) (map[string]string, error) {
	client := ssm.NewFromConfig(awsconfig)
	param, err := client.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(path),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		// Check specifically for ParameterNotFound error
		var pnf *types.ParameterNotFound
		if errors.As(err, &pnf) {
			return nil, fmt.Errorf("ssm parameter '%s' not found in this AWS account", path)
		}
		return nil, err
	}
	return dotenvlib.Parse(strings.NewReader(*param.Parameter.Value))
}

func merge(base map[string]string, variant map[string]string) map[string]string {
	for k, v := range variant {
		base[k] = v
	}
	return base
}

func extract(path string) (name string, feature string, err error) {
	re := regexp.MustCompile(substratePathPattern)
	matches := re.FindStringSubmatch(path)

	if matches == nil {
		return "", "", fmt.Errorf("invalid substrate path: %s", path)
	}

	// Get the named capture groups
	nameIndex := re.SubexpIndex("name")
	featureIndex := re.SubexpIndex("feature")

	return matches[nameIndex], matches[featureIndex], nil
}

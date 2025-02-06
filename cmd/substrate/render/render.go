package render

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/bkeane/substrate/pkg/substrate"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/charmbracelet/lipgloss/table"
)

type Root struct {
	Name     string   `arg:"positional" help:"substrate name" default:"agar"`
	Features []string `arg:"-f,--features" help:"feature flags"`
}

func (r *Root) Route(ctx context.Context, awsconfig aws.Config) (*string, error) {
	rendered, err := substrate.Parse(ctx, awsconfig, r.Name, r.Features...)
	if err != nil {
		return nil, fmt.Errorf("failed to render substrate: %w", err)
	}

	flattened := flatten(rendered)

	tbl := table.New()
	tbl.Headers("env", "field", "value")

	keys := make([]string, 0, len(flattened))
	for k := range flattened {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		values := flattened[k]
		tbl.Row(values[0], "substrate."+k, values[1])
	}

	result := tbl.Render()
	return &result, nil
}

func flatten(s any) map[string][]string {
	result := make(map[string][]string)
	val := reflect.ValueOf(s)

	// Handle pointer input by getting underlying value
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return result
	}

	var flatten func(v reflect.Value, prefix string)
	flatten = func(v reflect.Value, prefix string) {
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			structField := t.Field(i)
			name := structField.Name

			// Skip unexported fields
			if !structField.IsExported() {
				continue
			}

			if prefix != "" {
				name = prefix + "." + name
			}

			// Get env tag
			envTag := structField.Tag.Get("env")
			if envTag == "" {
				envTag = "-"
			}

			switch field.Kind() {
			case reflect.Struct:
				flatten(field, name)
			case reflect.Ptr:
				if !field.IsNil() {
					elem := field.Elem()
					if elem.Kind() == reflect.Struct {
						flatten(elem, name)
					} else {
						result[name] = []string{envTag, fmt.Sprintf("%v", elem.Interface())}
					}
				} else {
					result[name] = []string{envTag, "<nil>"}
				}
			case reflect.Slice, reflect.Array:
				result[name] = []string{envTag, fmt.Sprintf("%v", field.Interface())}
			default:
				result[name] = []string{envTag, fmt.Sprintf("%v", field.Interface())}
			}
		}
	}

	flatten(val, "")
	return result
}

package utils

import (
	"bytes"
	"context"

	"github.com/a-h/templ"
)

func RenderTemplToString(ctx context.Context, component templ.Component) (string, error) {
	var buf bytes.Buffer
	err := component.Render(ctx, &buf)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

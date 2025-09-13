package context

import "context"

// @WireSet("Infrastructure")
func NewContext() context.Context {
	ctx := context.Background()

	return ctx
}

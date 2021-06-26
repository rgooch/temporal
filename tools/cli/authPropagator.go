package cli

import (
	"context"
	"fmt"
	"os"

	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

type (
	// contextKey is an unexported type used as key for items stored in the
	// Context object
	contextKey struct{}

	// propagator implements the custom context propagator
	propagator struct{}

	// Values is a struct holding values
	AuthValue struct {
		Username string
		Groups []string
	}
)

// PropagateKey is the key used to store the value in the Context object
var PropagateKey = contextKey{}

// propagationKey is the key used by the propagator to pass values through the
// Temporal server headers
const propagationKey = "authentication-header"

// NewContextPropagator returns a context propagator that propagates a set of
// string key-value pairs across a workflow
func NewContextPropagator() workflow.ContextPropagator {
	fmt.Fprintln(os.Stderr, "HACK: NewContextPropagator()")
	return &propagator{}
}

// Inject injects values from context into headers for propagation
func (s *propagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	fmt.Fprintln(os.Stderr, "HACK: Inject()")
	value := ctx.Value(PropagateKey)
	payload, err := converter.GetDefaultDataConverter().ToPayload(value)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

// InjectFromWorkflow injects values from context into headers for propagation
func (s *propagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	fmt.Fprintln(os.Stderr, "HACK: InjectFromWorkflow()")
	value := ctx.Value(PropagateKey)
	payload, err := converter.GetDefaultDataConverter().ToPayload(value)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

// Extract extracts values from headers and puts them into context
func (s *propagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	fmt.Fprintln(os.Stderr, "HACK: Extract()")
	if value, ok := reader.Get(propagationKey); ok {
		var authValue AuthValue
		if err := converter.GetDefaultDataConverter().FromPayload(value,
			&authValue); err != nil {
			return ctx, nil
		}
		ctx = context.WithValue(ctx, PropagateKey, authValue)
	}

	return ctx, nil
}

// ExtractToWorkflow extracts values from headers and puts them into context
func (s *propagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	fmt.Fprintln(os.Stderr, "HACK: ExtractFromWorkflow()")
	if value, ok := reader.Get(propagationKey); ok {
		var authValue AuthValue
		if err := converter.GetDefaultDataConverter().FromPayload(value, &authValue); err != nil {
			return ctx, nil
		}
		ctx = workflow.WithValue(ctx, PropagateKey, authValue)
	}

	return ctx, nil
}

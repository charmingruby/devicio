package client

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/charmingruby/devicio/service/processor/pkg/instrumentation"
)

const (
	lowLatency int = iota
	mediumLatency
	highLatency
	unreliableLatency

	errProbability = 0.1
)

var (
	ErrUnstable = errors.New("the api is unstable")
	ErrUnknown  = errors.New("unknown error")

	latency = map[int]int{
		lowLatency:        100,
		mediumLatency:     200,
		highLatency:       300,
		unreliableLatency: 400,
	}
)

type UnstableAPI struct {
}

func NewUnstableAPI() UnstableAPI {
	return UnstableAPI{}
}

func (a *UnstableAPI) VolatileCall(ctx context.Context) (context.Context, error) {
	ctx, complete := instrumentation.Tracer.Span(ctx, "external.UnstableAPI.VolatileCall")
	defer complete()

	ctx, err := a.simulateLatency(ctx)
	if err != nil {
		return ctx, err
	}

	ctx, err = a.simulateErr(ctx)
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (a *UnstableAPI) simulateLatency(ctx context.Context) (context.Context, error) {
	traceID := instrumentation.Tracer.GetTraceIDFromContext(ctx)

	ctx, complete := instrumentation.Tracer.Span(ctx, "external.UnstableAPI.simulateLatency")
	defer complete()

	latency := latency[rand.Intn(len(latency))]

	instrumentation.Logger.Debug("Simulating latency", "latency", latency, "traceId", traceID)

	time.Sleep(time.Duration(latency) * time.Millisecond)

	if latency == unreliableLatency {
		return ctx, ErrUnstable
	}

	return ctx, nil
}

func (a *UnstableAPI) simulateErr(ctx context.Context) (context.Context, error) {
	ctx, complete := instrumentation.Tracer.Span(ctx, "external.UnstableAPI.simulateErr")
	defer complete()

	traceID := instrumentation.Tracer.GetTraceIDFromContext(ctx)

	shouldErr := rand.Float64() < errProbability

	instrumentation.Logger.Debug("Simulating error", "shouldErr", shouldErr, "traceId", traceID)

	if shouldErr {
		return ctx, ErrUnknown
	}

	return ctx, nil
}

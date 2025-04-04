package client

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/service/processor/pkg/logger"
)

const (
	lowLatency int = iota
	mediumLatency
	highLatency
	unreliableLatency

	errProbability = 0.2
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
	tracer observability.Tracer
}

func NewUnstableAPI(tracer observability.Tracer) UnstableAPI {
	return UnstableAPI{
		tracer: tracer,
	}
}

func (a *UnstableAPI) VolatileCall(ctx context.Context) (context.Context, error) {
	ctx, complete := a.tracer.Span(ctx, "external.UnstableAPI.VolatileCall")
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
	traceID := a.tracer.GetTraceIDFromContext(ctx)

	ctx, complete := a.tracer.Span(ctx, "external.UnstableAPI.simulateLatency")
	defer complete()

	latency := latency[rand.Intn(len(latency))]

	logger.Log.Info(fmt.Sprintf("got latency latency=%dms,traceId=%s,", latency, traceID))

	time.Sleep(time.Duration(latency) * time.Millisecond)

	if latency == unreliableLatency {
		return ctx, ErrUnstable
	}

	return ctx, nil
}

func (a *UnstableAPI) simulateErr(ctx context.Context) (context.Context, error) {
	ctx, complete := a.tracer.Span(ctx, "external.UnstableAPI.simulateErr")
	defer complete()

	traceID := a.tracer.GetTraceIDFromContext(ctx)

	shouldErr := rand.Float64() < errProbability

	logger.Log.Info(fmt.Sprintf("should err=%t,traceID=%s", shouldErr, traceID))

	if shouldErr {
		return ctx, ErrUnknown
	}

	return ctx, nil
}

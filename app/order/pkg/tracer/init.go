package tracer

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"

	"TikTokMall/app/order/conf"
)

func InitJaeger() (opentracing.Tracer, io.Closer, error) {
	config := conf.GetConf()
	cfg := &jaegercfg.Configuration{
		ServiceName: config.Service.Name,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  config.Jaeger.SamplerType,
			Param: config.Jaeger.SamplerParam,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           config.Jaeger.LogSpans,
			LocalAgentHostPort: fmt.Sprintf("%s:%d", config.Jaeger.Host, config.Jaeger.Port),
		},
	}

	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jaeger.StdLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("无法创建Jaeger tracer: %w", err)
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

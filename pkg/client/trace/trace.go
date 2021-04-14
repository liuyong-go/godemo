package trace

import (
	"fmt"
	"io"

	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var tracer opentracing.Tracer
var closer io.Closer

func NewTrace() (opentracing.Tracer, io.Closer) {
	if tracer != nil {
		return tracer, closer
	}
	conf := conf.NewYTrace().GetConfig()
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           conf.LogSpans,
			LocalAgentHostPort: conf.LocalAgentAddress,
		},
	}
	tracer, closer, err := cfg.New(conf.ServiceName, config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

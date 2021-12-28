package core

import (
	"github.com/openzipkin/zipkin-go"
	"github.com/rs/zerolog/log"
	"time"
)
import httpreporter "github.com/openzipkin/zipkin-go/reporter/http"

const endpointURL = "http://127.0.0.1:9411"

func GetTracer(serviceName string, ip string) (*zipkin.Tracer, error) {
	//创建一个供 tracer 使用的 reporter
	reporter := httpreporter.NewReporter(endpointURL)
	defer reporter.Close()
	//创建一个endpoint
	endpoint, err := zipkin.NewEndpoint(serviceName, ip)
	if err != nil {
		log.Err(err).Msgf("unable to create local endpoint: %+v\n", err)
		return nil, err
	}
	//设置取样策略
	//sampler := zipkin.NewModuloSampler(1)
	sampler, err := zipkin.NewBoundarySampler(0.01, time.Now().UnixNano())
	if err != nil {
		log.Err(err).Msgf("unable to create sampler: %+v\n", sampler)
		return nil, err
	}
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint), zipkin.WithSampler(sampler), zipkin.WithTraceID128Bit(true))
	//opentracing.SetGlobalTracer(tracer)
	if err != nil {
		log.Err(err).Msgf("unable to create tracer: %+v\n", tracer)
		return nil, err
	}
	return tracer, nil
}

func GetCollector() {

}

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/liuyong-go/godemo/pkg/client/trace"
	"github.com/liuyong-go/godemo/pkg/conf"
	"github.com/opentracing/opentracing-go"
)

func foo3(req string, ctx context.Context) (reply string) {
	//1.创建子span
	span, _ := opentracing.StartSpanFromContext(ctx, "span_foo3")
	defer func() {
		//4.接口调用完，在tag中设置request和reply
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	println(req)
	//2.模拟处理耗时
	time.Sleep(time.Second / 2)
	//3.返回reply
	reply = "foo3Reply"
	return
}

//跟foo3一样逻辑
func foo4(req string, ctx context.Context) (reply string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "span_foo4")
	defer func() {
		span.SetTag("request", req)
		span.SetTag("reply", reply)
		span.Finish()
	}()

	println(req)
	time.Sleep(time.Second / 2)
	reply = "foo4Reply"
	return
}
func main() {
	conf.ConfPath = "/Users/liuyong/go/src/godemo/toml/toml_dev/"
	tracer, closer := trace.NewTrace()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer) //StartspanFromContext创建新span时会用到

	span := tracer.StartSpan("span_root")
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	r1 := foo3("Hello foo3", ctx)
	r2 := foo4("Hello foo4", ctx)
	fmt.Println(r1, r2)
	span.Finish()
}

package quickproxy

import (
	"github.com/elazarl/goproxy"
	"log"
	. "net/http"
	"time"
)

const (
	DEFAULT_PORT = "8888"
	LISTEN       = "127.0.0.1"
	VERBOSE      = false
)

type (
	requestCallback  func(request *Request)
	responseCallback func(response *Response)
	doneCallback     func(doneReuqestData *DoneRequestData)
)

type (
	requestCallbacks  []requestCallback
	responseCallbacks []responseCallback
	doneCallbacks     []doneCallback
)

type RoundTrip struct {
	RequestTime  time.Time
	ResponseTime time.Time
}

type DoneRequestData struct {
	Request       *Request
	Response      *Response
	RoundTripTime *RoundTrip
}

var (
	port                string = DEFAULT_PORT
	listen              string = LISTEN
	verbose             bool   = VERBOSE
	sessions                   = map[int64]*RoundTrip{}
	onRequestCallbacks         = requestCallbacks{}
	onResponseCallbacks        = responseCallbacks{}
	onDoneCallbacks            = doneCallbacks{}
)

func NewRoundTrip(requestTime time.Time) *RoundTrip {
	roundTrip := new(RoundTrip)
	roundTrip.RequestTime = requestTime
	return roundTrip
}

func (r *RoundTrip) ElapsedTime() time.Duration {
	return r.ResponseTime.Sub(r.RequestTime)
}

func Prepare(options map[string]string) {
	if options["port"] != "" {
		port = options["port"]
	}

	if options["listen"] != "" {
		listen = options["listen"]
	}

	if options["verbose"] != "" {
		if options["verbose"] == "true" {
			verbose = true
		} else {
			verbose = false
		}
	}
}

func OnRequest(f requestCallback) {
	onRequestCallbacks = append(onRequestCallbacks, f)
}

func OnResponse(f responseCallback) {
	onResponseCallbacks = append(onResponseCallbacks, f)
}

func OnDone(f doneCallback) {
	onDoneCallbacks = append(onDoneCallbacks, f)
}

func Run() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(r *Request, ctx *goproxy.ProxyCtx) (*Request, *Response) {
			sessions[ctx.Session] = NewRoundTrip(time.Now())

			for _, f := range onRequestCallbacks {
				f(ctx.Req)
			}

			return r, nil
		})

	proxy.OnResponse().DoFunc(
		func(r *Response, ctx *goproxy.ProxyCtx) *Response {
			for _, f := range onResponseCallbacks {
				f(ctx.Resp)
			}

			if ctx.Resp != nil {
				roundTrip := sessions[ctx.Session]
				roundTrip.ResponseTime = time.Now()

				for _, f := range onDoneCallbacks {
					f(&DoneRequestData{ctx.Req, ctx.Resp, roundTrip})
				}
			}

			return r
		})

	log.Fatal(ListenAndServe(listen+":"+port, proxy))
}

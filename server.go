package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"github.com/vladivolo/go_daemon_skeleton/stats"
	log "github.com/vladivolo/lumber"
	"io"
	"time"
)

type ResponseVersion struct {
	Binary         string
	Version        string
	Maintainer     string
	AutoBuildTag   string
	BuildDate      string
	BuildHost      string
	BuildGoVersion string
	BuildCommand   string
}

type RequestHandlerInfo struct {
	get_handler    func(*fasthttp.RequestCtx)
	post_handler   func(*fasthttp.RequestCtx)
	delete_handler func(*fasthttp.RequestCtx)
}

var (
	workersCloser []io.Closer
	VersionInfo   ResponseVersion
)

var Handlers = map[string]RequestHandlerInfo{
	"/ping":          RequestHandlerInfo{GetPing, PostPing, nil},
	"/service/stats": RequestHandlerInfo{GetStats, nil, nil},
	"/version":       RequestHandlerInfo{GetVersion, nil, nil},
}

func StartHttpServer(addr string, workers_count int) {
	for i := 0; i < workers_count; i++ {
		lc := listen(addr)
		workersCloser = append(workersCloser, lc)
	}
}

func CloseHttpServer() {
	for _, cl := range workersCloser {
		cl.Close()
	}
	time.Sleep(2 * time.Second) // waiting for it to be inflated
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	stats.HttpRequestsInc()

	ctx.Response.Header.Set("daemon-api-version", "1.0")

	r, ok := Handlers[string(ctx.Path())]
	if ok == true {
		switch string(ctx.Method()) {
		case "GET":
			if r.get_handler != nil {
				r.get_handler(ctx)
				return
			}
		case "POST":
			if r.post_handler != nil {
				r.post_handler(ctx)
				return
			}
		case "DELETE":
			if r.delete_handler != nil {
				r.delete_handler(ctx)
				return
			}
		}
	}

	stats.HttpUnknownRequestsInc()
	ctx.Error("Not found", fasthttp.StatusNotFound)
}

func listen(addr string) io.Closer {
	ln, err := reuseport.Listen("tcp4", addr)
	if err != nil {
		log.Fatal("error in reuseport listener: %s", err)
	}

	go func() {
		s := fasthttp.Server{
			Handler: requestHandler,
			//DisableKeepalive: true,
		}
		if err = s.Serve(ln); err != nil {
			log.Fatal("error in fasthttp Server: %s", err)
		}
		log.Info("fasthttp.Serve(): EXIT")
	}()

	return ln
}

func PostPing(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "PONG")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetPing(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "PONG")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetStats(ctx *fasthttp.RequestCtx) {
	r, err := stats.Stats()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	fmt.Fprintf(ctx, r)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetVersion(ctx *fasthttp.RequestCtx) {
	bin, err := json.Marshal(VersionInfo)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	fmt.Fprintf(ctx, string(bin))
	ctx.SetStatusCode(fasthttp.StatusOK)
}

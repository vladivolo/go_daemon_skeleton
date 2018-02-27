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

var (
	workersCloser []io.Closer
	VersionInfo   ResponseVersion
)

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

	switch string(ctx.Path()) {
	case "/ping":
		PingHandler(ctx)
	case "/service/stats":
		StatsHandler(ctx)
	case "/version", "/version/": // it's better to configure the redirect in nginx
		VersionHandler(ctx)
	default:
		stats.HttpUnknownRequestsInc()
		ctx.Error("Not found", fasthttp.StatusNotFound)
		return
	}
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
		log.Info("fasthttp.Serve: EXIT")
	}()

	return ln
}

func PingHandler(ctx *fasthttp.RequestCtx) {
	log.Debug("PingHandler() User-Agent: %s RemoteAddr %s", ctx.UserAgent(), ctx.RemoteAddr())

	switch string(ctx.Method()) {
	case "POST":
		PostPing(ctx)
	case "GET":
		GetPing(ctx)
	default:
	}

	return
}

func PostPing(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "PONG")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func GetPing(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "PONG")
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func StatsHandler(ctx *fasthttp.RequestCtx) {
	r, err := stats.Stats()
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}
	fmt.Fprintf(ctx, r)
	ctx.SetStatusCode(fasthttp.StatusOK)
}

func VersionHandler(ctx *fasthttp.RequestCtx) {
	if bin, err := json.Marshal(VersionInfo); err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	} else {
		fmt.Fprintf(ctx, string(bin))
		ctx.SetStatusCode(fasthttp.StatusOK)
		return
	}
}

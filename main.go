package main

import (
	"flag"
	"github.com/shudiwsh2009/reservation_thxl_go/service"
	"github.com/shudiwsh2009/reservation_thxl_go/web"
	"log"
	"os"
	"os/signal"
	"runtime"
)

func main() {
	var webAddress, debugAssetsPort, confPath string
	var isDebug, isSmock bool
	flag.StringVar(&webAddress, "web", ":9000", "Web address server listening on")
	flag.StringVar(&debugAssetsPort, "devWeb", "", "Web address server listening on (like :9010)")
	flag.StringVar(&confPath, "conf", "deploy/thxl.conf", "Configuration file path for service")
	flag.BoolVar(&isDebug, "debug", false, "Debug mode")
	flag.BoolVar(&isSmock, "smock", true, "Smock server")
	flag.Parse()

	service.InitService(confPath, isSmock)
	server := web.NewServer(isDebug)
	if isDebug && debugAssetsPort != "" {
		server.SetAssetDomain("//localhost" + debugAssetsPort)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		s := <-c
		log.Printf("Got signal: %s", s)
		server.Cleanup()
	}()

	log.Fatal(server.ListenAndServe(webAddress))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

package main

import (
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/service"
	"bitbucket.org/shudiwsh2009/reservation_thxl_go/web"
	"flag"
	"log"
	"runtime"
)

func main() {
	var webAddr, confPath string
	var isDebug, isSmock bool
	flag.StringVar(&webAddr, "web", ":9000", "Web address server listening on")
	flag.StringVar(&confPath, "conf", "./deploy/thxl.conf", "Configuration file path for backends")
	flag.BoolVar(&isDebug, "debug", false, "Debug mode")
	flag.BoolVar(&isSmock, "smock", true, "Smock server")
	flag.Parse()

	service.InitService(confPath, isSmock)
	server := web.NewServer(confPath, isDebug)
	log.Fatal(server.ListenAndServe(webAddr))
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

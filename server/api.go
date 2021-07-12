package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bigscreen/manga-scrapper/config"
	"github.com/bigscreen/manga-scrapper/dep"
	"github.com/bigscreen/manga-scrapper/logger"
	"github.com/codegangsta/negroni"
)

func StartAPIServer() {
	logger.InfoF("Starting Mangajack api server in port %d", config.Port())

	deps := dep.InitAppDependencies()
	muxRouter := Router(deps)
	handlerFunc := muxRouter.ServeHTTP

	n := negroni.New(negroni.NewRecovery())
	n.UseHandlerFunc(handlerFunc)
	portInfo := ":" + strconv.Itoa(config.Port())
	server := &http.Server{Addr: portInfo, Handler: n}
	go listenServer(server)
	waitForShutdown(server)
}

func listenServer(apiServer *http.Server) {
	err := apiServer.ListenAndServe()
	if err != http.ErrServerClosed {
		logger.Fatal("Mangajack fails to serve, err:", err)
	}
}

func waitForShutdown(apiServer *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Info("Mangajack is shutting down")
	_ = apiServer.Shutdown(context.Background())
	logger.Info("Mangajack shutdown complete")
}

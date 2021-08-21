package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sbpann/go-docker-multi-stage-build-graceful-shutdown-example/helpers"
)

func main() {
	restarted := 0
	maxRestart := 4
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "foobar",
		})
	})

	server := helpers.HttpServerFromGinEngine(r)

	startServer := func(server *http.Server, callback func()) {
		callback()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("start server error: %s\n", err)
		}
	}
	// run server in another goroutine to avoid blocking
	go startServer(server, func() { log.Println("Starting server") })

	serverRestarter := func(osSignal chan os.Signal, server *http.Server) *http.Server {
		newServer := helpers.HttpServerFromGinEngine(r)
		helpers.SignalNotify(&helpers.SignalNotifyArgs{
			OSSignal: osSignal,
			CatchingSignals: []os.Signal{
				syscall.SIGINT,
				syscall.SIGTERM,
			},
			Callback: func() {
				ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
				defer cancel()
				log.Println("Shutdown server...")
				if err := server.Shutdown(ctx); err != nil {
					log.Fatalf("shutdown server error: %s", err)
				}
				log.Println("Shutdown successfully")
				log.Printf("Server has restarted %d times\n", restarted)
				restarted++
				if restarted <= maxRestart {
					go startServer(newServer, func() { log.Println("Restarting server...") })
				}
			},
		})
		return newServer
	}
	osSignal := make(chan os.Signal)
	for restarted <= maxRestart {
		// returned newServer to be shutdown next time
		server = serverRestarter(osSignal, server)
	}
	close(osSignal)
	log.Println("Maximum restart time is reached")
	log.Println("Program exited")
}

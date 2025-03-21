package main

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	messageCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "message_count",
		Help: "Total number of messages sent.",
	})
	_ = time.Now()
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	prometheus.MustRegister(messageCounter)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			log.Fatal().Err(err).Msg("Failed to start metrics server")
		}
	}()

	var wg sync.WaitGroup
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stopChan:
				log.Info().Msg("Program Closed")
				return
			default:
				log.Info().Msg("Hi All Here")
				messageCounter.Inc()
				time.Sleep(1 * time.Second)
			}
		}
	}()

	wg.Wait()
}

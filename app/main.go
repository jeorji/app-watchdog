package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	defaultPort        = "8080"
	defaultFailRateStr = "0"
	defaultDelayStr    = "0"
)

type Config struct {
	Addr       string
	FailRate   float64
	MaxDelayMs int
}

func main() {
	log.Printf("build: version=%s commit=%s date=%s",
		ReleaseVersion, ReleaseCommit, ReleaseDate,
	)

	cfg := loadConfig()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/healthz", healthHandler(cfg))

	log.Printf("starting http server on %s", cfg.Addr)

    if err := http.ListenAndServe(cfg.Addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func loadConfig() Config {
    port := os.Getenv("APP_PORT")
    if port == "" {
        port = defaultPort
    }

    failRateStr := os.Getenv("HEALTH_FAIL_RATE")
	if failRateStr == "" {
		failRateStr = defaultFailRateStr
	}

	failRate, err := strconv.ParseFloat(failRateStr, 64)
	if err != nil {
		log.Printf("invalid HEALTH_FAIL_RATE=%q, using default %s", failRateStr, defaultFailRateStr)
		failRate, _ = strconv.ParseFloat(defaultFailRateStr, 64)
	}
	if failRate < 0 {
		failRate = 0
	}
	if failRate > 1 {
		failRate = 1
	}

	delayStr := os.Getenv("HEALTH_MAX_DELAY_MS")
	if delayStr == "" {
		delayStr = defaultDelayStr
	}

	maxDelayMs, err := strconv.Atoi(delayStr)
	if err != nil || maxDelayMs < 0 {
		log.Printf("invalid HEALTH_MAX_DELAY_MS=%q, using default %s", delayStr, defaultDelayStr)
	}

	cfg := Config{
        Addr:       "0.0.0.0:" + port,
		FailRate:   failRate,
		MaxDelayMs: maxDelayMs,
	}

	log.Printf("config: addr=%s fail_rate=%.3f max_delay_ms=%d", cfg.Addr, cfg.FailRate, cfg.MaxDelayMs)
	return cfg
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Hello World!\n"))
	if err != nil {
		log.Printf("error writing response /: %v", err)
	}
}

func healthHandler(cfg Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cfg.MaxDelayMs > 0 {
			delay := rand.Intn(cfg.MaxDelayMs + 1) // 0..MaxDelayMs
			time.Sleep(time.Duration(delay) * time.Millisecond)
			log.Printf("healthz: delay=%dms", delay)
		}

		x := rand.Float64()
		if x < cfg.FailRate {
			log.Printf("healthz: forced failure (x=%.3f < %.3f)", x, cfg.FailRate)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("error\n"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	}
}

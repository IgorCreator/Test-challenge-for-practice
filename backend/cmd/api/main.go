package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"nimble-challenge/backend/internal/auth"
	"nimble-challenge/backend/internal/config"
	"nimble-challenge/backend/internal/crypto"
	"nimble-challenge/backend/internal/db"
	"nimble-challenge/backend/internal/graphql"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	cipher, err := crypto.NewCipherFromBase64(cfg.EncryptionKeyB64)
	if err != nil {
		log.Fatalf("crypto: %v", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB,
	)
	store, err := db.NewStore(context.Background(), dsn, cipher)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer store.Close()

	if err := store.EnsureDemoData(context.Background(), cfg.StoreSlug, cfg.StoreName, cfg.MerchantUser, cfg.MerchantPass, cfg.CustomerUser, cfg.CustomerPass); err != nil {
		log.Fatalf("seed: %v", err)
	}

	handler := graphql.NewHandler(store)
	handler = auth.Middleware(store)(handler)
	handler = withCORS(handler)
	handler = withRateLimit(handler, 120, time.Minute)

	mux := http.NewServeMux()
	mux.Handle("/graphql", handler)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.AppPort),
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
		TLSConfig:         loadTLSConfig(cfg),
	}

	go func() {
		log.Printf("API listening on https://localhost:%d/graphql", cfg.AppPort)
		if err := serveTLS(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("serve: %v", err)
		}
	}()

	waitForShutdown(srv)
}

func serveTLS(srv *http.Server) error {
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	tlsListener := tls.NewListener(ln, srv.TLSConfig)
	return srv.Serve(tlsListener)
}

func loadTLSConfig(cfg config.Config) *tls.Config {
	if cfg.TLSCertPath != "" && cfg.TLSKeyPath != "" {
		if _, err := os.Stat(cfg.TLSCertPath); err == nil {
			if _, err := os.Stat(cfg.TLSKeyPath); err == nil {
				cert, err := tls.LoadX509KeyPair(cfg.TLSCertPath, cfg.TLSKeyPath)
				if err == nil {
					return &tls.Config{Certificates: []tls.Certificate{cert}}
				}
			}
		}
	}
	cert, err := crypto.SelfSignedTLSCert()
	if err != nil {
		log.Fatalf("self-signed cert: %v", err)
	}
	return &tls.Config{Certificates: []tls.Certificate{cert}}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func withRateLimit(next http.Handler, maxTokens int, refill time.Duration) http.Handler {
	type bucket struct {
		tokens int
		last   time.Time
	}
	var (
		mu        sync.Mutex
		bucketMap = map[string]*bucket{}
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		mu.Lock()
		b, ok := bucketMap[ip]
		if !ok {
			b = &bucket{tokens: maxTokens, last: time.Now()}
			bucketMap[ip] = b
		}
		now := time.Now()
		elapsed := now.Sub(b.last)
		if elapsed >= refill {
			b.tokens = maxTokens
			b.last = now
		}
		if b.tokens <= 0 {
			mu.Unlock()
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		b.tokens--
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func waitForShutdown(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

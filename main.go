package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/imsteev/whoxy/integrations/clerk"
	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port          string `env:"PORT,default=9000"`
	RedisHost     string `env:"REDIS_HOST,default=localhost"`
	RedisPort     string `env:"REDIS_PORT,default=6379"`
	RedisPassword string `env:"REDIS_PASSWORD,default="`
	RedisDB       int    `env:"REDIS_DB,default=0"`
}

func main() {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + cfg.RedisPort,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	clerkIntegration := &clerk.ClerkIntegration{
		Redis: redisClient,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("whoxy"))
	})

	http.HandleFunc("/clerk", func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		dest, err := clerkIntegration.GetDestination(body)
		if err != nil {
			log.Fatal(err)
		}

		req, err := ForwardPostRequest(r, dest, body)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
	})

	http.HandleFunc("POST /mappings", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		var requests []struct {
			ClientIP    string `json:"client_ip"`
			Destination string `json:"destination"`
			ServiceName string `json:"service_name"`
		}

		if err := json.Unmarshal(body, &requests); err != nil {
			log.Fatal(err)
		}

		for _, request := range requests {
			redisClient.Set(context.Background(), request.ServiceName+":"+request.ClientIP, request.Destination, 0)
		}
	})

	// DELETE /mappings/clerk/12.34.56.78
	http.HandleFunc("DELETE /mappings/:service_name/:client_ip", func(w http.ResponseWriter, r *http.Request) {
		serviceName := r.URL.Query().Get("service_name")
		clientIP := r.URL.Query().Get("client_ip")
		redisClient.Del(context.Background(), serviceName+":"+clientIP)
		w.WriteHeader(http.StatusNoContent)
	})

	addr := ":" + cfg.Port
	log.Printf("listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

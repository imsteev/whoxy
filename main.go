package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/imsteev/whoxy/integrations/clerk"
	"github.com/redis/go-redis/v9"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port     string `env:"PORT"`
	RedisUrl string `env:"REDIS_URL"`
}

type MappingRequest struct {
	InternalSlug string `json:"internal_slug"`
	ExternalSlug string `json:"external_slug"`
}

func main() {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("redis url: %s\n", cfg.RedisUrl)

	redisOpts, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		log.Fatal(err)
	}

	redisOpts.OnConnect = func(ctx context.Context, cn *redis.Conn) error {
		log.Printf("connected to redis\n")
		return nil
	}

	redisClient := redis.NewClient(redisOpts)

	clerkIntegration := &clerk.ClerkIntegration{
		Redis: redisClient,
	}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "whoxy")
	})

	r.POST("/:service", func(c *gin.Context) {
		serviceName := c.Param("service")
		log.Printf("Processing request for service: %s\n", serviceName)

		body, err := c.GetRawData()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		dest, err := clerkIntegration.GetDestination(body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		req, err := ForwardPostRequest(c.Request, dest, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		c.Status(resp.StatusCode)
	})

	r.POST("/:service/mappings", func(c *gin.Context) {
		service := c.Param("service")

		var requests []MappingRequest
		if err := c.ShouldBindJSON(&requests); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, request := range requests {
			redisClient.Set(context.Background(), service+":"+request.InternalSlug, request.ExternalSlug, 0)
		}

		c.JSON(http.StatusOK, gin.H{"message": "mappings created successfully"})
	})

	r.DELETE("/:service/mappings/:internal_slug", func(c *gin.Context) {
		service := c.Param("service")
		internalSlug := c.Param("internal_slug")

		err := redisClient.Del(context.Background(), service+":"+internalSlug).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	})

	port := cfg.Port
	if port == "" {
		port = "9000"
	}

	if _, err := strconv.Atoi(port); err != nil {
		log.Fatal("Invalid port number")
	}

	addr := ":" + port
	log.Printf("listening on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

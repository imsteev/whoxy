package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"slices"
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

type RoutingRequest struct {
	EventKey       string `json:"event_key"`
	DestinationUrl string `json:"destination_url"`
}

var supportedServices = []string{"clerk"}

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

	clerkIntegration := &clerk.ClerkIntegration{}

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "whoxy")
	})

	// webhooks are forwarded to this endpoint
	r.POST("/webhooks/:service", func(c *gin.Context) {
		serviceName := c.Param("service")
		log.Printf("Processing request for service: %s\n", serviceName)

		clonedReq := c.Request.Clone(context.Background())
		eventKey, err := clerkIntegration.GetEventKey(clonedReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		destinationUrl, err := redisClient.Get(context.Background(), serviceName+":"+eventKey).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		forwardResp, err := ForwardPostRequest(c.Request, destinationUrl, body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(forwardResp.StatusCode)
	})

	r.POST("/webhooks/:service/routes", func(c *gin.Context) {
		service := c.Param("service")

		if !slices.Contains(supportedServices, service) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "service not supported"})
			return
		}

		var requests []RoutingRequest
		if err := c.ShouldBindJSON(&requests); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, request := range requests {
			redisClient.Set(context.Background(), service+":"+request.EventKey, request.DestinationUrl, 0)
		}

		c.JSON(http.StatusOK, gin.H{"message": "mappings created successfully"})
	})

	r.DELETE("/webhooks/:service/routes/:event_key", func(c *gin.Context) {
		service := c.Param("service")
		eventKey := c.Param("event_key")

		if !slices.Contains(supportedServices, service) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "service not supported"})
			return
		}

		err := redisClient.Del(context.Background(), service+":"+eventKey).Err()
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

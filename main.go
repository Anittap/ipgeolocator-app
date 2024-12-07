package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
)

type GeoData struct {
	Cached       string `json:"cached"`
	ApiServer    string `json:"apiServer"`
	Version      string `json:"version"`
	IP           string `json:"ip"`
	Continent    string `json:"continent_name"`
	Country      string `json:"country_name"`
	City         string `json:"city"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	ISP          string `json:"isp"`
	Organization string `json:"organization"`
}

var (
	memcachedClient *memcache.Client
	apiKey          string
	apiServer       string
	version         = "version1"
)

func fetchAPIKeyFromSecretsManager(secretName string) (string, error) {
	region := os.Getenv("REGION_NAME")
	if region == "" {
		return "", fmt.Errorf("missing REGION_NAME environment variable")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return "", fmt.Errorf("error loading AWS configuration: %w", err)
	}

	svc := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("error retrieving secret: %w", err)
	}

	// Parse the secret value (assuming JSON format)
	var secretData map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secretData); err != nil {
		return "", fmt.Errorf("error parsing secret JSON: %w", err)
	}

	apiKey, exists := secretData["API_KEY"]
	if !exists {
		return "", fmt.Errorf("API_KEY not found in the secret")
	}

	return apiKey, nil
}

func getFromCache(ip string) (*GeoData, error) {
	item, err := memcachedClient.Get(ip)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		return nil, fmt.Errorf("error retrieving from cache: %w", err)
	}

	var geoData GeoData
	err = json.Unmarshal(item.Value, &geoData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling cached data: %w", err)
	}

	geoData.Cached = "True"
	geoData.ApiServer = apiServer
	geoData.Version = version

	return &geoData, nil
}

func fetchFromAPI(ip string) (*GeoData, error) {
	url := fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=%s&ip=%s", apiKey, ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	var geoData GeoData
	err = json.NewDecoder(resp.Body).Decode(&geoData)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling API response: %w", err)
	}

	geoData.Cached = "False"
	geoData.ApiServer = apiServer
	geoData.Version = version

	return &geoData, nil
}

func setToCache(ip string, geoData *GeoData) error {
	data, err := json.Marshal(geoData)
	if err != nil {
		return fmt.Errorf("error marshaling data for cache: %w", err)
	}

	err = memcachedClient.Set(&memcache.Item{
		Key:        ip,
		Value:      data,
		Expiration: int32(time.Hour.Seconds()),
	})
	if err != nil {
		return fmt.Errorf("error setting data to cache: %w", err)
	}

	return nil
}

func main() {
	apiServer = os.Getenv("HOSTNAME")
	memcachedHost := os.Getenv("MEMCACHED_HOST")
	memcachedPort := os.Getenv("MEMCACHED_PORT")
	appPort := os.Getenv("APP_PORT")
	secretName := os.Getenv("SECRET_NAME")
	useSecretsManager := os.Getenv("API_KEY_FROM_SECRETSMANAGER") == "True"

	if apiServer == "" {
		apiServer = "none"
	}
	if appPort == "" {
		appPort = "8080"
	}

	if memcachedHost == "" || memcachedPort == "" {
		log.Fatal("Environment variables MEMCACHED_HOST and MEMCACHED_PORT must be set")
	}

	if useSecretsManager {
		if secretName == "" {
			log.Fatal("SECRET_NAME must be set when using Secrets Manager")
		}

		var err error
		apiKey, err = fetchAPIKeyFromSecretsManager(secretName)
		if err != nil {
			log.Fatalf("Error fetching API key from Secrets Manager: %v", err)
		}
	} else {
		apiKey = os.Getenv("API_KEY")
		if apiKey == "" {
			log.Fatal("Environment variable API_KEY must be set if not using Secrets Manager")
		}
	}

	memcachedClient = memcache.New(fmt.Sprintf("%s:%s", memcachedHost, memcachedPort))

	router := gin.Default()

	router.GET("/ip/:ip", func(c *gin.Context) {
		ip := c.Param("ip")

		// Try fetching from cache
		geoData, err := getFromCache(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if geoData != nil {
			c.JSON(http.StatusOK, geoData)
			return
		}

		// Fetch from API
		geoData, err = fetchFromAPI(ip)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Cache the result
		err = setToCache(ip, geoData)
		if err != nil {
			log.Printf("Error caching data: %v", err)
		}

		c.JSON(http.StatusOK, geoData)
	})

	router.GET("/status", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.Run(fmt.Sprintf(":%s", appPort))
}

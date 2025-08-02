package videos

import (
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"strings"
	"time"
)

func ListAll(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"success": true,
	})
}

// MakeAnonymousRequest makes a GET request while masking server information
func MakeAnonymousRequest(url string) ([]byte, error) {
	// Create a custom HTTP client with modifications to hide server info
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			DisableKeepAlives: true, // Prevent connection reuse
		},
	}

	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers to mask the request origin and mimic a regular browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	// Remove or mask headers that might reveal server information
	req.Header.Del("X-Forwarded-For")
	req.Header.Del("X-Real-IP")
	req.Header.Del("X-Forwarded-Proto")
	req.Header.Del("X-Forwarded-Host")

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle compressed responses
	var reader io.Reader = resp.Body

	// Check if response is gzip compressed
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	// Read the response body
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// ListNewest - Your Fiber handler that uses the anonymous request
func ListNewest(c *fiber.Ctx) error {
	// Get the URL from query parameter or use a default
	targetURL := c.Query("url", "https://nguon.vsphim.com/api/danh-sach/phim-moi-cap-nhat?page=1") // Default URL for testing

	// Make the anonymous request
	responseBody, err := MakeAnonymousRequest(targetURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch data",
			"message": err.Error(),
		})
	}

	// Try to parse as JSON first
	var jsonData interface{}
	if err := json.Unmarshal(responseBody, &jsonData); err == nil {
		// If it's valid JSON, return it as JSON
		return c.JSON(fiber.Map{
			"success": true,
			"data":    jsonData,
			"url":     targetURL,
		})
	}

	// If not JSON, return as string
	return c.JSON(fiber.Map{
		"success": true,
		"data":    string(responseBody),
		"url":     targetURL,
	})
}

func Details(c *fiber.Ctx) error {
	var request struct {
		Slug string `json:"slug"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Invalid request body",
			"message": err.Error(),
		})
	}

	if request.Slug == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Missing 'slug' in request body",
		})
	}

	// Construct the target URL for details
	url := fmt.Sprintf("https://nguon.vsphim.com/api/phim/%v", request.Slug)
	targetURL := c.Query("url", url)

	// Make the anonymous request
	responseBody, err := MakeAnonymousRequest(targetURL)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Failed to fetch data",
			"message": err.Error(),
		})
	}

	// Try to parse the response as JSON
	var jsonData interface{}
	if err := json.Unmarshal(responseBody, &jsonData); err == nil {
		return c.JSON(fiber.Map{
			"success": true,
			"data":    jsonData,
			"url":     targetURL,
		})
	}

	// Fallback to returning raw response as string
	return c.JSON(fiber.Map{
		"success": true,
		"data":    string(responseBody),
		"url":     targetURL,
	})
}

// Alternative version using a proxy service (more anonymous)
func MakeProxyRequest(url string) ([]byte, error) {
	// You can use free proxy services or set up your own
	// This example shows how to structure it
	proxyURL := "https://api.allorigins.win/get?url=" + url

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", proxyURL, nil)
	if err != nil {
		return nil, err
	}

	// Standard browser headers
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Enhanced version with more anonymization features
func MakeStealthRequest(url string, options ...RequestOption) ([]byte, error) {
	config := &RequestConfig{
		Timeout:      30 * time.Second,
		UserAgent:    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		DisableCache: true,
		RandomizeUA:  false,
	}

	// Apply options
	for _, opt := range options {
		opt(config)
	}

	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			DisableKeepAlives: true,
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set headers based on config
	req.Header.Set("User-Agent", config.UserAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	if config.DisableCache {
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")
	}

	// Remove identifying headers
	req.Header.Del("X-Forwarded-For")
	req.Header.Del("X-Real-IP")
	req.Header.Del("CF-Connecting-IP")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Configuration types for the enhanced version
type RequestConfig struct {
	Timeout      time.Duration
	UserAgent    string
	DisableCache bool
	RandomizeUA  bool
}

type RequestOption func(*RequestConfig)

func WithTimeout(timeout time.Duration) RequestOption {
	return func(c *RequestConfig) {
		c.Timeout = timeout
	}
}

func WithUserAgent(ua string) RequestOption {
	return func(c *RequestConfig) {
		c.UserAgent = ua
	}
}

func WithRandomUA() RequestOption {
	return func(c *RequestConfig) {
		c.RandomizeUA = true
		userAgents := []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36",
		}
		c.UserAgent = userAgents[time.Now().Unix()%int64(len(userAgents))]
	}
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	blue = "\x1b[34m(+)\x1b[0m"
)

type counter struct {
	count int
	mu    sync.Mutex
}

type promoGenerator struct {
}

func (p *promoGenerator) generatePromo() {
	apiurl := "https://api.discord.gx.games/v1/direct-fulfillment"
	partnerUserID := generateUUID()

	payload, err := json.Marshal(map[string]string{
		"partnerUserId": partnerUserID,
	})
	if err != nil {
		log.Fatalf("Error marshalling JSON: %s\n", err)
	}

	proxyUrl, err := url.Parse("your-proxy")

	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	req, err := http.NewRequest("POST", apiurl, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatalf("Error creating HTTP request: %s\n", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Sec-Ch-Ua", "\"Opera GX\";v=\"105\", \"Chromium\";v=\"119\", \"Not?A_Brand\";v=\"24\"")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 OPR/105.0.0.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making HTTP request: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s\n", err)
	}

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			log.Fatalf("Error unmarshalling JSON: %s\n", err)
		}
		token, ok := result["token"].(string)
		if ok {
			c.mu.Lock()
			c.count++
			c.mu.Unlock()

			link := fmt.Sprintf("https://discord.com/billing/partner-promotions/1180231712274387115/%s", token)
			file, err := os.OpenFile("promos.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Fatalf("Error opening file: %s\n", err)
			}
			defer file.Close()

			if _, err := fmt.Fprintf(file, "%s\n", link); err != nil {
				log.Fatalf("Error writing to file: %s\n", err)
			}

			fmt.Printf("%s Generated Promo Link : %s\n", getTimestamp(), link)
		}
	} else if resp.StatusCode == http.StatusTooManyRequests {
		fmt.Printf("%s You are being rate-limited!\n", getTimestamp())
	} else {
		fmt.Printf("%s Request failed : %d\n", getTimestamp(), resp.StatusCode)
	}
}

func generateUUID() string {
	return fmt.Sprintf("%s-%s-%s-%s-%s", randomString(8), randomString(4), randomString(4), randomString(4), randomString(12))
}

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var result strings.Builder
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charSet))
		result.WriteByte(charSet[randomIndex])
	}
	return result.String()
}

func getTimestamp() string {
	timeIDK := time.Now().Format("15:04:05")
	return fmt.Sprintf("[\x1b[90m%s\x1b[0m]", timeIDK)
}

var c = counter{}

func main() {
	var numThreads int
	fmt.Printf("%s %s Enter Number Of Threads : ", getTimestamp(), blue)
	fmt.Scan(&numThreads)

	var wg sync.WaitGroup
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			generator := promoGenerator{}
			for {
				generator.generatePromo()
			}
		}()
	}
	wg.Wait()
}

package main

import (
	"fmt"
	"io"
	"encoding/json"
	"net/http"
	"time"
    "github.com/joho/godotenv"
	"log"
	"os"

)

func main() {

	ip, err := getIP()
	if err != nil {
		fmt.Println(err)
	}
	location, err := getLocation(ip)
    if err != nil {
        fmt.Printf("Failed to get location. Error: %v", err)
        return
    }
    fmt.Printf("Location: %s, %s", location.City, location.CountryName)
	getPrayerTimes(location.City)
}



func getIP() (string, error) {
    resp, err := http.Get("https://api.ipify.org?format=json")
    if err != nil {
        return "", fmt.Errorf("failed to get public IP. Error: %v", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response body. Error: %v", err)
    }
    var data map[string]string
    err = json.Unmarshal(body, &data)
    if err != nil {
        return "", fmt.Errorf("failed to unmarshal response body. Error: %v", err)
    }
    ip := data["ip"]
    return ip, nil
}

type Location struct {
    CountryName string `json:"country_name"`
    City        string `json:"city"`
}

func getLocation(ip string) (*Location, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessKey := os.Getenv("ACCESS_KEY")

	url := fmt.Sprintf("http://api.ipstack.com/%s?access_key=%s", ip, accessKey)
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to get location. Error: %v", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body. Error: %v", err)
    }
    var location Location
    err = json.Unmarshal(body, &location)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal response body. Error: %v", err)
    }
	fmt.Println(location)
    return &location, nil

}

type PrayerTimes struct {
    Date       string `json:"date"`
    Fajr       string `json:"Fajr"`
    Sunrise    string `json:"Sunrise"`
    Dhuhr      string `json:"Dhuhr"`
    Asr        string `json:"Asr"`
    Sunset     string `json:"Sunset"`
    Maghrib    string `json:"Maghrib"`
    Isha       string `json:"Isha"`
    Imsak      string `json:"Imsak"`
    Midnight   string `json:"Midnight"`
}

func getPrayerTimes(city string) {
	prayerTimes, err := callAPI(city)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Prayer times for %s on %s:\n", city, prayerTimes.Date)
	fmt.Printf("Fajr: %s\n", formatTime(prayerTimes.Fajr))
	fmt.Printf("Sunrise: %s\n", formatTime(prayerTimes.Sunrise))
	fmt.Printf("Dhuhr: %s\n", formatTime(prayerTimes.Dhuhr))
	fmt.Printf("Asr: %s\n", formatTime(prayerTimes.Asr))
	fmt.Printf("Sunset: %s\n", formatTime(prayerTimes.Sunset))
	fmt.Printf("Maghrib: %s\n", formatTime(prayerTimes.Maghrib))
	fmt.Printf("Isha: %s\n", formatTime(prayerTimes.Isha))
	fmt.Printf("Imsak: %s\n", formatTime(prayerTimes.Imsak))
	fmt.Printf("Midnight: %s\n", formatTime(prayerTimes.Midnight))
}

func formatTime(timeStr string) string {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return timeStr
	}
	return t.Format("03:04 PM")
}

func callAPI(city string) (*PrayerTimes, error) {
	url := fmt.Sprintf("http://api.aladhan.com/v1/timingsByCity?city=%s&country=egypt&method=2", city)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get prayer times. Error: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body. Error: %v", err)
	}
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body. Error: %v", err)
	}
	prayerTimes := data["data"].(map[string]interface{})["timings"].(map[string]interface{})
	date := time.Now().Format("2006-01-02")
	return &PrayerTimes{
		Date:       date,
		Fajr:       prayerTimes["Fajr"].(string),
		Sunrise:    prayerTimes["Sunrise"].(string),
		Dhuhr:      prayerTimes["Dhuhr"].(string),
		Asr:        prayerTimes["Asr"].(string),
		Sunset:     prayerTimes["Sunset"].(string),
		Maghrib:    prayerTimes["Maghrib"].(string),
		Isha:       prayerTimes["Isha"].(string),
		Imsak:      prayerTimes["Imsak"].(string),
		Midnight:   prayerTimes["Midnight"].(string),
	}, nil
}
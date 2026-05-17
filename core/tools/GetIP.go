package tools

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type IPGeoData struct {
	Country   string  `json:"country"`
	City      string  `json:"city"`
	ISP       string  `json:"isp"`
	Query     string  `json:"query"`
	Region    string  `json:"regionName"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Org       string  `json:"org"`
	AS        string  `json:"as"`  
	Browser   string  `json:"browser"`
	OS        string  `json:"os"`
	Device    string  `json:"device"`
	UserAgent string  `json:"user_agent"`
}
func GetIPDetails(userAgent string) (*IPGeoData, error) {
	// URL untuk lookup IP
	url := "http://ip-api.com/json/"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	// Decode hasil response ke struct
	var data IPGeoData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Parse User-Agent untuk informasi lebih lanjut
	data.UserAgent = userAgent
	data.Browser, data.OS, data.Device = ParseUserAgent(userAgent)

	return &data, nil
}

func ParseUserAgent(userAgent string) (browser string, os string, device string) {
	// Deteksi browser
	if strings.Contains(userAgent, "Chrome") {
		browser = "Chrome"
	} else if strings.Contains(userAgent, "Firefox") {
		browser = "Firefox"
	} else if strings.Contains(userAgent, "Safari") {
		browser = "Safari"
	} else {
		browser = "Unknown"
	}

	// Deteksi OS
	if strings.Contains(userAgent, "Windows") {
		os = "Windows"
	} else if strings.Contains(userAgent, "Macintosh") {
		os = "MacOS"
	} else if strings.Contains(userAgent, "Linux") {
		os = "Linux"
	} else if strings.Contains(userAgent, "Android") {
		os = "Android"
	} else if strings.Contains(userAgent, "iPhone") {
		os = "iOS"
	} else {
		os = "Unknown"
	}
	if strings.Contains(userAgent, "Mobi") {
		device = "Mobile"
	} else {
		device = "Desktop"
	}

	return browser, os, device
}

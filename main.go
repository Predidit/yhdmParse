package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func decodeText(text string) (string, error) {
	if !stringContains(text, "{") {
		var decodedText string
		key := 1561
		length := len(text)
		for i := 0; i < length; i += 2 {
			charCode, err := strconv.ParseInt(text[i:i+2], 16, 64)
			if err != nil {
				return "", err
			}
			charCode = (charCode + 1048576 - int64(key) - int64(length/2-1-i/2)) % 256
			decodedText = string(rune(charCode)) + decodedText
		}
		text = decodedText
	}
	decodedValue, err := url.QueryUnescape(text)
	if err != nil {
		return "", err
	}
	return decodedValue, nil
}

func stringContains(str, substr string) bool {
	return strings.Index(str, substr) != -1
}

func parseVideoUrl(responseText string) (string, error) {
	decodedText, err := decodeText(responseText)
	if err != nil {
		return "", err
	}
	var videoInfo map[string]interface{}
	err = json.Unmarshal([]byte(decodedText), &videoInfo)
	if err != nil {
		return "", err
	}
	return videoInfo["vurl"].(string), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <responseText>")
		os.Exit(1)
	}
	responseText := os.Args[1]
	vurl, err := parseVideoUrl(responseText)
	if err != nil {
		fmt.Printf("Error parsing video URL: %s\n", err)
		os.Exit(1)
	}
	fmt.Println(vurl)
}

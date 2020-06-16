package lib

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Realface(savename string) {
	var err error
	var dst *os.File
	var req *http.Request
	var resp *http.Response
	var durat time.Duration
	req, err = http.NewRequest("GET", "https://thispersondoesnotexist.com/image", nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	client := &http.Client{Timeout: durat}
	resp, err = client.Do(req)
	if err != nil {
		log.Println(err)
	}
	dst, err = os.Create("./img/" + savename)
	if err != nil {
		log.Println(err)
	}
	_, err = io.Copy(dst, resp.Body)
	if err != nil {
		log.Println(err)
	}
}

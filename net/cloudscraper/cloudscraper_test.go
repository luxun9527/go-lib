package cloudscraper

import (
	"bytes"
	"encoding/json"
	hawk "github.com/juiced-aio/hawk-go"
	http "github.com/useflyent/fhttp"
	"github.com/useflyent/fhttp/cookiejar"
	"io"
	"log"
	"testing"
)

func TestCloudscraper(t *testing.T) {

	// Client has to be from fhttp and up to CloudFlare's standards, this can include ja3 fingerprint/http2 settings.
	client := http.Client{}
	// Client also will need a cookie jar.
	cookieJar, _ := cookiejar.New(nil)
	client.Jar = cookieJar
	scraper := hawk.CFInit(client, "YOUR_KEY_HERE", true)
	b := map[string]interface{}{"alias": "", "is_public": true, "group_id": 1, "url": "https://www.api.com"}
	d, _ := json.Marshal(b)
	// You will have to create your own function if you want to solve captchas.
	scraper.CaptchaFunction = func(originalURL string, siteKey string) (string, error) {
		// CaptchaFunction should return the token as a string.
		return "", nil
	}
	req, _ := http.NewRequest("POST", "https://goo.su/api/links/create", bytes.NewBuffer(d))

	req.Header = http.Header{
		"sec-ch-ua":                 {`"Chromium";v="92", " Not A;Brand";v="99", "Google Chrome";v="92"`},
		"sec-ch-ua-mobile":          {`?0`},
		"upgrade-insecure-requests": {`1`},
		"user-agent":                {`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36`},
		"accept":                    {`text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9`},
		"sec-fetch-site":            {`none`},
		"sec-fetch-mode":            {`navigate`},
		"sec-fetch-user":            {`?1`},
		"sec-fetch-dest":            {`document`},
		"accept-encoding":           {`gzip, deflate`},
		"accept-language":           {`en-US,en;q=0.9`},
		http.HeaderOrderKey:         {"sec-ch-ua", "sec-ch-ua-mobile", "upgrade-insecure-requests", "user-agent", "accept", "sec-fetch-site", "sec-fetch-mode", "sec-fetch-user", "sec-fetch-dest", "accept-encoding", "accept-language"},
		http.PHeaderOrderKey:        {":method", ":authority", ":scheme", ":path"},
		"x-goo-api-token":           {"W4exIOZZym2nxYwr3puHz2orStoMjdM1QmQVXvZidskftEgOi4VLqfufRKv1"},
	}

	resp, err := scraper.Do(req)
	if err != nil {
		log.Printf("err =%v\n", err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("err =%v", err)
	}
	log.Println(string(body))
}

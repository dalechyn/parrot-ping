package main

import "net/http"

func ping(url string, c chan<- PingRequest) {
	resp, err := http.Get(url)
	p := PingRequest{url, 1}
	if err != nil {
	}
	if resp.StatusCode == 200 {
		p.result = 0
	}
	c <- p
}


package main

import "net/http"

func ping(id int, url string, c chan<- URLRequest) {
	/*
	 * ping - reads the head of http response, returns 1 if website
	 * is reachable
	 */
	resp, err := http.Get(url)
	p := URLRequest{id, url, -1}
	if err != nil {
		p.result = 404
	} else {
		p.result = resp.StatusCode
	}
	c <- p
}


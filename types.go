package main

type URLRequest struct {
	id int
	url string
	result int
}

type URLValidation struct {
	url string
	err error
}

type WatchWorker struct {
	id int
	url string
	delay float64
	err error
}


package main

import (
	"net/url"
)

func ValidateURL(link string) error {
	/*
	 * Checks the integrity of the URL
	 */
	 _, err := url.ParseRequestURI(link)
	 return err
}


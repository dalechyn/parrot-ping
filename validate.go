package main

import (
	"regexp"
	"errors"
)

func ValidateURL(url string) error {
	/*
	 * Checks the integrity of the URL
	 */
	 if !regexp.MustCompile(`^https?:\/\/.*\..*`).MatchString(url) {
		return errors.New(`URL doesn't contain http/https protocol
							or is bad`)
	 }
	 return nil
}


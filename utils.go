package main

func ValidateAndFixURL(url string) URLValidation  {
	/*
	 * Validates and tries to fix url if it is wrong, returns error
	 * if it can't fix it
	 */
	if ValidateURL(url) != nil {
		url := "http://" + url
		if e := ValidateURL(url); e != nil {
			return URLValidation{"", e}
		} else {
			return URLValidation{url, nil}
		}
	} else {
		return URLValidation{url, nil}
	}
}

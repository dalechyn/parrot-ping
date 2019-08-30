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

func getExtension(fileURL string) string {
	var ext string
	for i := len(fileURL) - 1; i >= 0; i-- {
		if fileURL[i] != '.' {
			// TODO: Find how to append to the beggining and not use this
			ext = string(fileURL[i]) + ext
		} else {
			break
		}
	}
	return ext
}

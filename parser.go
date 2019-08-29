package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"errors"
)

func parseURLs(args []string, output chan<- URLValidation) {
	for _, badURL := range args {
		go func () {
			output <- ValidateAndFixURL(badURL)
		}()
	}
}

func parseURLsFromFile(fileURL string, output chan<- []URLValidation) {
	/*
	 * reads the file and processes file content to the URL
	 * format if the links in files are valid
	 */
	content, err := ioutil.ReadFile(fileURL)
	result := []URLValidation{}

	if err != nil {
		output <- result
		return
	}

	// trim and validate all URLs concurrently
	URLs := strings.Split(strings.Join(strings.Fields(string(content)), " "),	" ")

	urlChannel := make(chan URLValidation, len(URLs))

	parseURLs(URLs, urlChannel)

	for i := 0; i < len(URLs); i++ {
		result = append(result, <-urlChannel)
	}
	close(urlChannel)
	output <- result
}

func parseWatchWorkers(watchArgs []string,
						workersChannel chan<- WatchWorker) {
	/*
	 * parses arguments and puts the program in the infinite
	 * watch mode
	 */
	for i, arg := range watchArgs {
		/*
		 * creating WatchWorker and fixing badURL
		 */
		go func() {
			res := WatchWorker{i, "", -1, nil}
			// parsing urls and delays
			splitted := strings.Split(arg, ":")
			if len(splitted) != 2 {
				res.err = errors.New("Wrong argument: " + arg + "in --watch")
				workersChannel <- res
				return
			}
			res.url = splitted[0]
			// fixing url
			urlVal := ValidateAndFixURL(splitted[0])
			if urlVal.err != nil {
				res.err = urlVal.err
				workersChannel <- res
				return
			}
			res.url = urlVal.url
			// parsing delay
			delay, err := strconv.ParseFloat(splitted[1], 64)
			if err != nil {
				res.err = err
				workersChannel <- res
				return
			}
			res.delay = delay
		}()
	}
}


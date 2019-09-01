package main

import (
	"io/ioutil"
	"strings"
	"strconv"
	"errors"
)

func parseFile(fileURL string, func parseFunc (string) []interface{}) []interface {} {
	content, err := ioutil.ReadFile(fileURL)

	if err != nil {
		return nil
	}

	formatted := strings.Split(strings.Join(strings.Fields(string(content)), " "),	" ")
	return parseFunc(formatted)
}

func parseURLs(urls []string) []URLValidation {
	/*
	 * parses all providen urls concurrently
	 */
	res := make([]URLValidation, len(urls))
	output := make(chan URLValidation, len(urls))
	for _, badURL := range urls {
		go func(badURL string) {
			output <- ValidateAndFixURL(badURL)
		}(badURL)
	}
	for i := 0; i < len(urls); i++ {
		res[i] = <-output
	}
	close(output)
	return res
}

func parseURLsFromFile(fileURL string) []URLValidation {
	/*
	 * reads the file and processes file content to the URL
	 * format if the links in files are valid
	 */
	return parseFile(fileURL, parseURLs)
}

func parseWatchers(watchArgs []string) []WatchWorker {
	/*
	 * parses arguments and puts the program in the infinite
	 * watch mode
	 */
	res := make([]WatchWorker, len(watchArgs))
	output := make(chan WatchWorker, len(watchArgs))
	for i, arg := range watchArgs {
		/*
		 * creating WatchWorker and fixing badURL
		 */
		go func(i int, arg string) {
			res := WatchWorker{i, "", -1, nil}
			// parsing urls and delays
			splitted := strings.Split(arg, ":")
			if len(splitted) != 2 {
				res.err = errors.New("Wrong argument: " + arg + "in --watch")
				output <- res
				return
			}
			res.url = splitted[0]
			// fixing url
			urlVal := ValidateAndFixURL(splitted[0])
			if urlVal.err != nil {
				res.err = urlVal.err
				output <- res
				return
			}
			res.url = urlVal.url
			// parsing delay
			delay, err := strconv.Atoi(splitted[1])
			if err != nil {
				res.err = err
				output <- res
				return
			}
			res.delay = delay
			output <- res
		}(i, arg)
	}
	for i := 0; i < len(watchArgs); i++ {
		res[i] = <-output
	}
	close(output)
	return res
}

func parseWatchersFromFile(fileURL string) []WatchWorker {
	/*
	 * parses watchers from file, function name speaks for itself
	 */
	return parseFile(fileURL, parseWatchers)
}

func parseBridges(bridgeArgs []string) []Bridge {
	res := make([]Bridge, len(bridgeArgs))
	output := make(chan Bridge, len(bridgeArgs))
	for _, arg := range bridgeArgs {
		go func(i int, arg string) {
			res := Bridge{"", "", nil, nil}
			splitted := strings.Split(arg, ":")
			if len(splitted) != 2 {
				res.err = errors.New("Wrong argument: " + arg + "in --bridge")
				output <- res
			}
			if _, exists := supportedBridges[splitted[0]]; !exists {
				res.err = errors.New("Not supported bridge: " + splitted[0])
				output <- res
			}
			res.bridgeType = splitted[0]
			res.token = splitted[1]
			output <- res
		}
	}
	for i := 0; i < len(bridgeArgs); i++ {
		res[i] = <-output
	}
	return res
}

func parseBridgesFromFile(fileURL string) []Bridge {
	content, err := ioutil.ReadFile(fileURL)

	if err != nil {
		return nil
	}

	Bridges := strings.Split(strings.Join(strings.Fields(string(content)), " "),	" ")

	return parseBridges(Bridges)
}

func parseBridgesFromFile(fileURL string) []Bridge {
	return parseFile(fileURL, parseBridges)
}

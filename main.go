package main

import
(
	"fmt"
	"os"
	"strings"
	"regexp"
	"time"
	"sync"
)

/*
 * URLsources array by it's nature has always contain valid urls
 */

var outerURLs []string
var outerWatchers []WatchWorker

func main() {
	/*
	 * main() function reads console arguments and routes them to the
	 * aliasRouter and to the argRouter
	 */
	aliasRouter := map[string] string {
		"b": "bridge",
		"i": "input",
		"p": "ping",
		"w": "watch",
	}

	argRouter := map[string] func([]string) {
		"help": func (helpArgs []string)  {
			if len(helpArgs) != 0 {
				Error(`--help doesn't take any arguments`)
				return
			}
			fmt.Println("parrot-watcher", version, ":")
			fmt.Println("--help: show help")
			fmt.Println("-b or --bridge: pass messages to bridge")
			fmt.Println("-i or --input: add a file as source, must have .{b,p,w}list extension")
			fmt.Println("-p or --ping: ping website from list")
			fmt.Println("-w or --watch: watch on sites in url:delay format")
			fmt.Println("\t delay in sec")
		},
		"input": func (inputArgs []string) {
			/*
			 * provide path to list of files for next steps
			 */
			extensionRouter := map[string] func(string) {
				"plist": func (fileURL string) {
					for _, url := range parseURLsFromFile(fileURL) {
						if url.err != nil {
							Warning(
								"--input incorrect url from file %s, Error: %s",
								fileURL,
								url.err)
						} else {
							outerURLs = append(outerURLs, url.url)
						}
					}
				},
				"wlist": func (fileURL string) {
					for _, watcher := range parseWatchersFromFile(fileURL) {
						if watcher.err != nil {
							Warning("--input incorrect watcher from file %s, Error: %s", fileURL, watcher.err)
						} else {
							outerWatchers = append(outerWatchers, watcher)
						}
					}
				},
			}

			/* 
			 * as we start fileParsing in goroutines, we have to wait
			 * until they all finish
			 */
			var wg sync.WaitGroup
			for _, arg := range inputArgs {
				ext := getExtension(arg)
				if _, exists := extensionRouter[ext]; !exists {
					Warning("--input wrong file extension, must be .{p,w}list")
				} else {
					wg.Add(1)
					go func(fileName string) {
						defer wg.Done()
						extensionRouter[ext](fileName)
					}(arg)
				}
			}
			wg.Wait()
		},
		"ping": func (pingArgs []string) {
			/*
			 * pings list of provided urls and list of sites provided
			 * by --input
			 */

			// fixing bar urls
			outputURLs := make(chan URLValidation, len(pingArgs))
			URLs := []string{}
			for _, res := range parseURLs(pingArgs) {
				if res.err != nil {
					Warning("--ping: %s: %s", res.url, res.err)
				} else {
					URLs = append(URLs, res.url)
				}
			}
			close(outputURLs)

			// making requests
			outputReq := make(chan URLRequest, len(URLs) + len(outerURLs))
			for i, url := range append(URLs, outerURLs...) {
				go func(i int, url string) {
					outputReq <- ping(i, url)
				}(i, url)
			}
			fmt.Printf("PING OUTPUT:\n\n")
			for i := 0; i < len(URLs) + len(outerURLs); i++ {
				res := <-outputReq
				fmt.Printf("ID: %d | URL: %s | RES: %d\n", res.id, res.url, res.result)
			}
			print("\n\n")
			close(outputReq)
		},
		"watch": func (watchArgs []string) {
			/*
			 * parses arguments and puts the program in the infinite
			 * watch mode
			 */
			watchersValidated := []WatchWorker{}
			for _, worker := range parseWatchers(watchArgs) {
				if worker.err != nil {
					Warning("--watch %s", worker.err)
				} else {
					watchersValidated = append(watchersValidated, worker)
				}
			}
			output := make(chan URLRequest, len(watchersValidated))
			for _, watcher := range append(watchersValidated, outerWatchers...) {
				go func(watcher WatchWorker) {
					for {
						time.Sleep(time.Duration(watcher.delay) * time.Millisecond)
						output <- ping(watcher.id, watcher.url)
					}
				}(watcher)
			}
			fmt.Printf("WATCH OUTPUT:\n\n")
			for pingRes := range output {
				fmt.Printf("Watcher %d | URL: %s | Result %d\n", pingRes.id,
							pingRes.url, pingRes.result)
			}
			close(output)
		},
	}

	if len(os.Args) == 1 {
		argRouter["help"]([]string{})
		return
	}

	clArgs := regexp.MustCompile(` *-{1,2}`).Split(strings.Join(os.Args[1:], " "), -1)

	clArgsSplitted := make([][]string, 0)

	for i := range clArgs[1:] {
		clArgsSplitted = append(clArgsSplitted, strings.Split(clArgs[1:][i], " "))
	}
	if len(clArgsSplitted) == 0 {
		Error("No arguments provided")
		return
	}
	for i := range clArgsSplitted {
		parseFunc, exists := argRouter[aliasRouter[clArgsSplitted[i][0]]]
		if !exists {
			parseFunc, exists = argRouter[clArgsSplitted[i][0]]
			if !exists {
				Error("Wrong argument, write --help for help")
			return
			}
		}
		parseFunc(clArgsSplitted[i][1:])
	}
}


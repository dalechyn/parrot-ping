package main

import
(
	"fmt"
	"os"
	"strings"
	"regexp"
)

/*
 * URLsources array by it's nature has always contain valid urls
 */

var URLOthers []string

func main() {
	/*
	 * main() function reads console arguments and routes them to the
	 * aliasRouter and to the argRouter
	 */
	aliasRouter := map[string] string {
		"p": "ping",
		"i": "input",
	}

	argRouter := map[string] func ([]string) {
		"help": func (helpArgs []string)  {
			if len(helpArgs) != 0 {
				Error(`--help doesn't take any arguments`)
				return
			}
			fmt.Println("parrot-watcher", version, ":")
			fmt.Println("--list: show watching websites")
			fmt.Println("--ping: ping website from list")
			fmt.Println("--input: add a file as url source")
			fmt.Println("--watch: watch on sites in url:delay format")
			fmt.Println("\t delay in sec")
		},
		"ping": func (pingArgs []string) {
			/*
			 * pings list of provided urls and list of sites provided
			 * by --input
			 */

			// fixing bar urls
			outputURLs := make(chan URLValidation, len(pingArgs))
			parseURLs(pingArgs, outputURLs)
			URLs := []string{}
			for i := 0; i < len(pingArgs); i++ {
				if res := <-outputURLs; res.err != nil {
					Warning("--ping: %s", res.err)
				} else {
					URLs = append(URLs, res.url)
				}
			}
			close(outputURLs)

			// making requests
			outputReq := make(chan URLRequest, len(URLs) + len(URLOthers))
			for i, url := range append(URLs, URLOthers...) {
				go ping(i, url, outputReq)
			}
			for i := 0; i < len(URLs) + len(URLOthers); i++ {
				res := <-outputReq
				fmt.Printf("ID: %d | URL: %s | RES: %d\n", res.id, res.url, res.result)
			}
		},
		"input": func (inputArgs []string) {
			/*
			 * provide path to list of files for next steps
			 */
			output := make(chan []URLValidation, len(inputArgs))

			for _, arg := range inputArgs {
				go parseURLsFromFile(arg, output)
			}

			for i := 0; i < len(inputArgs); i++ {
				for _, url := range <-output {
					if url.err != nil {
						Warning("--input %s", url.err)
					} else {
						URLOthers = append(URLOthers, url.url)
					}
				}
			}
		},
		"watch": func (watchArgs []string) {
			/*
			 * parses arguments and puts the program in the infinite
			 * watch mode
			 */
			//watchList := []WatchWorker{}
			workersChannel := make(chan WatchWorker, len(watchArgs))
			parseWatchWorkers(watchArgs, workersChannel)
			for i := 1; i < len(watchArgs); i++ {
				if worker := <-workersChannel; worker.err != nil {
				Warning("--watch %s", worker.err)
				} else {
					// I will code that tommorow
				}
			}
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


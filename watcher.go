package main

import
(
	"fmt"
	"os"
	"strings"
	"regexp"
)

type PingRequest struct {
	url string
	result byte
}

type WatchUnit struct {
	id int
	url string
}

func pingHandler (pingArgs []string) {
	results := make(chan PingRequest, len(pingArgs))
	for i := 0; i < len(pingArgs); i++ {
		go ping(pingArgs[i], results)
	}
	for i := 0; i < len(pingArgs); i++ {
		res := <-results
		fmt.Println(res.url, res.result)
	}
}

func main(){
	aliasList := map[string] string {
		"p": "ping",
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
			fmt.Println("--input: prompt a input file for ")
		},
		"ping": pingHandler,
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

	for i := range clArgsSplitted {
		parseFunc, exists := argRouter[clArgsSplitted[i][0]]
		if !exists {
			parseFunc, exists = argRouter[aliasList[clArgsSplitted[i][0]]]
			if !exists {
				fmt.Println("Wrong argument, write --help for help")
			return
			}
		}
		parseFunc(clArgsSplitted[i][1:])
	}
}


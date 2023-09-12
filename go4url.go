package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var (
	noDisplayErrors bool
)

func init() {
	flag.Usage = func() {
		fmt.Println(`
		
██████   ██████   ██   ██ ██    ██ ██████  ██      
██       ██    ██ ██   ██ ██    ██ ██   ██ ██      
██   ███ ██    ██ ███████ ██    ██ ██████  ██      
██    ██ ██    ██      ██ ██    ██ ██   ██ ██      
 ██████   ██████       ██  ██████  ██   ██ ███████ 
												   
																																																																						   
			
			`)

		fmt.Println("go4url v.0.1")
		fmt.Println("Author: ninposec")
		fmt.Println("")
		fmt.Println("Usage: cat urls.txt | go4url")
		fmt.Println("Extract URLs from Input e.g. JS Files")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
}

func main() {
	// Parse command-line flags
	urlsFile := flag.String("urls", "", "File containing URLs")
	concurrency := flag.Int("c", 1, "Concurrency level")
	flag.BoolVar(&noDisplayErrors, "nd", false, "Ignore and suppress error messages")
	flag.Parse()

	var urls []string

	// Read URLs from the -urls flag if provided
	if *urlsFile != "" {
		fileURLs, err := readURLsFromFile(*urlsFile)
		if err != nil {
			printError("Failed to read URLs from file:", err)
			os.Exit(1)
		}
		urls = append(urls, fileURLs...)
	}

	// Read URLs from STDIN if no URLs are provided via the -urls flag
	if len(urls) == 0 {
		stdinURLs, err := readURLsFromStdin()
		if err != nil {
			printError("Failed to read URLs from stdin:", err)
			os.Exit(1)
		}
		urls = append(urls, stdinURLs...)
	}

	if len(urls) == 0 {
		fmt.Println("No URLs provided.")
		os.Exit(1)
	}

	uniqueURLs := uniqueStrings(urls)
	fullURLs := make(chan string)
	var wg sync.WaitGroup

	// Spawn worker goroutines for processing full URLs
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for url := range fullURLs {
				endpoints, err := extractFullURLs(url)
				if err != nil {
					if !noDisplayErrors {
						printError(fmt.Sprintf("Failed to extract URLs for %s:", url), err)
					}
					continue
				}
				printEndpoints(endpoints)
			}
		}()
	}

	// Feed URLs to the worker goroutines
	for _, u := range uniqueURLs {
		fullURLs <- u
	}

	close(fullURLs)

	wg.Wait()
}

func readURLsFromFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		urls = append(urls, strings.TrimSpace(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	return urls, nil
}

func readURLsFromStdin() ([]string, error) {
	var urls []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		urls = append(urls, strings.TrimSpace(line))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading input from stdin: %s", err)
	}

	return urls, nil
}

func extractFullURLs(u string) ([]string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	response, err := client.Get(u)
	if err != nil {
		if !noDisplayErrors && strings.Contains(err.Error(), "no such host") {
			// Ignore "no such host" error
			return nil, nil
		}
		return nil, fmt.Errorf("failed to make HTTP request: %s", err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %s", err)
	}

	// Extract full URLs
	fullRegex := regexp.MustCompile(`https?://[^\s"'>]+`)
	fullMatches := fullRegex.FindAllString(string(body), -1)

	var urls []string
	for _, match := range fullMatches {
		urls = append(urls, match)
	}

	// Extract relative URL paths starting and ending with "/"
	relativeRegex := regexp.MustCompile(`(?m)^/[\w/.-]+/$`)
	relativeMatches := relativeRegex.FindAllString(string(body), -1)

	for _, match := range relativeMatches {
		urls = append(urls, match)
	}

	return urls, nil
}

func printEndpoints(endpoints []string) {
	sort.Strings(endpoints)
	uniqueEndpoints := uniqueStrings(endpoints)

	for _, endpoint := range uniqueEndpoints {
		fmt.Println(endpoint)
	}
}

func uniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func getDomainFromURL(u string) string {
	parsedURL, _ := url.Parse(u)
	return parsedURL.Scheme + "://" + parsedURL.Hostname()
}

func printError(msg string, err error) {
	if !noDisplayErrors {
		fmt.Printf("%s %v\n", msg, err)
	}
}

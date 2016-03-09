package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func httpGET(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	src, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(src), nil
}

func extractThreads(src string) []string {
	var result []string
	reg := regexp.MustCompile(`[<a href="res/[0-9]{8,12}" class="replylink">`)
	matches := reg.FindAllStringSubmatch(src, -1)

	for _, s := range matches {
		tmp := s[0]
		tmp = tmp[:strings.Index(tmp, "\"")]
		result = append(result, tmp)
	}

	return result
}

func checkThread(srcThread string, pic string) bool {
	if strings.Contains(srcThread, pic) {
		return true
	}
	return false
}

var (
	flagFile  = flag.String("file", "", "The imagename you want to find the thread for")
	flagQuery = flag.String("query", "", "Filter all threads with this word")
	flagChan  = flag.String("chan", "4chan", "The imageboard you want to search on (-list for supported boards")
	flagList  = flag.Bool("list", false, "List the supported imageboards.")
	flagBoard = flag.String("board", "b", "The specific board to search")
)

func main() {
	flag.Parse()

	var sBaseURL = "http://boards.4chan.org"
	var sBoard = "s"
	var sURL = sBaseURL + "/" + sBoard + "/"

	i := 0
	for i < 15 {
		url := sURL + strconv.Itoa(i)
		fmt.Println("Download Page: " + url)
		srcPage, err := httpGET(url)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Println("Extract Threads from Page " + strconv.Itoa(i))
		threads := extractThreads(srcPage)

		for _, thread := range threads {
			srcThread, err := httpGET(sURL + thread)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Println("Check Thread: " + thread)
			result := checkThread(srcThread, *flagFile)
			if result {
				fmt.Println("Pic found in Thread: " + sURL + thread)
				return
			}
		}

		i++
	}

}

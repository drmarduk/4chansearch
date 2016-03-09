package main

import (
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

func main() {
	var iStart int = 0
	var iEnd int = 11
	var sPicture = "1384107007405.jpg"
	var sBaseURL = "http://boards.4chan.org"
	var sBoard = "s"
	var sURL = sBaseURL + "/" + sBoard + "/"

	// get user input

	// get Threads
	// Base URL: http://$chan.de
	fmt.Print("Please insert Board URL (no / at the end pl0x): ")
	fmt.Scanf("%s", &sBaseURL)

	// Board
	fmt.Print("Which Chan (no /b/, only b): ")
	fmt.Scanf("%s", &sBoard)

	// search foo
	fmt.Print("Which Picture do you need? Pls only filename (should be uniq enough, if not, your imageboard sucks)")
	fmt.Scanf("%s", &sPicture)

	fmt.Print("How many pages are on " + sBaseURL + "/" + sBoard + "/ ?: ")
	fmt.Scanf("%d", &iEnd)

	fmt.Println("kk, check.")
	fmt.Println("BaseURL: " + sBaseURL)
	fmt.Println("Channel: " + sBoard)
	fmt.Println("maxPage: " + strconv.Itoa(iEnd))

	for i := iStart; i < iEnd; i++ {

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
			result := checkThread(srcThread, sPicture)
			if result {
				fmt.Println("Pic found in Thread: " + sURL + thread)
				return
			}
		}
	}

}

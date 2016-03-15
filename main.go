package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Chan interface {
	getMaxPage() int
	getUrl(string, int) string
	getCatalogUrl(string) string
	getThreads(string, string) []string
	hasCatalog() bool
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

var (
	flagFile  = flag.String("file", "", "The imagename you want to find the thread for")
	flagQuery = flag.String("query", "", "Filter all threads with this word")
	flagChan  = flag.String("chan", "4chan", "The imageboard you want to search on (-list for supported boards")
	flagList  = flag.Bool("list", false, "List the supported imageboards.")
	flagBoard = flag.String("board", "b", "The specific board to search")
)

func main() {
	flag.Parse()

	var obj Chan
	switch *flagChan {
	case "4chan":
		obj = New4Chan()
	}

	var threads []string
	if !obj.hasCatalog() {
		for i := 0; i < obj.getMaxPage(); i++ {
			url := obj.getUrl(*flagBoard, i)

			srcPage, err := httpGET(url)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			threads = obj.getThreads(*flagBoard, srcPage)

			// make
			go dlThreads(obj, threads)
		}
	} else {
		// download catalog
		src, err := httpGET(obj.getCatalogUrl(*flagBoard))
		if err != nil {
			fmt.Println("Error while downloading catalog.")
			return
		}
		threads = obj.getThreads(*flagBoard, src)
		fmt.Printf("got %d threads\n", len(threads))
		dlThreads(obj, threads)
	}
}

func dlThreads(obj Chan, threads []string) {
	for _, t := range threads {
		src, err := httpGET(t)
		if err != nil {
			fmt.Println("Error while downloading " + t)
			continue
		}
		if checkThread(src, *flagFile) {
			fmt.Println("Found: " + t)
			os.Exit(0) // well, good idea?
		}
	}
}

func checkThread(srcThread string, pic string) bool {
	if strings.Contains(srcThread, pic) {
		return true
	}
	return false
}

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

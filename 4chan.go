package main

import "strconv"

type FourChan struct {
}

func New4Chan() FourChan {
	return FourChan{}
}

func (c FourChan) hasCatalog() bool {
	return true
}

func (c FourChan) getMaxPage() int {
	return 15
}

func (c FourChan) getUrl(b string, i int) string {
	return "https://boards.4chan.org/" + b + "/" + strconv.Itoa(i)
}

func (c FourChan) getThreads(src string) []string {
	var result []string
	return result
}

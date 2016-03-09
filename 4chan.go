package main

import (
	"regexp"
	"strconv"
	"strings"
)

type FourChan struct {
}

func New4Chan() FourChan {
	return FourChan{}
}

func (c FourChan) hasCatalog() bool {
	return true
}

func (c FourChan) getCatalogUrl(b string) string {
	return "https://boards.4chan.org/" + b + "/catalog"
}

func (c FourChan) getMaxPage() int {
	return 15
}

func (c FourChan) getUrl(b string, i int) string {
	return "https://boards.4chan.org/" + b + "/" + strconv.Itoa(i)
}

func (c FourChan) getThreads(src string) []string {
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

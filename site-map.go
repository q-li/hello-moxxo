package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
)

// A SiteMap records all the mapping between an URL
// to all the internal (as within the same domain) URLs that the page contains
type SiteMap map[string]map[string]bool

// PutInMap puts the mapping parent -> parsed in site.(*siteMap),
// and also checks whether parsed is already visited.
// If we know parsed is an external link, return true
// so that don't follow the link and fetch an external page
func (site *SiteMap) PutInMap(parent, parsed *url.URL, internal bool) bool {
	// ignore fragments
	parent.Fragment = ""
	parsed.Fragment = ""

	// e.g. https://monzo.com/about/ -> https://monzo.com/about
	// eliminates unnecessary duplicates
	parentStr := strings.TrimRight(parent.String(), "/")
	parsedStr := strings.TrimRight(parsed.String(), "/")

	siteM := *site
	if _, ok := siteM[parentStr]; !ok {
		siteM[parentStr] = map[string]bool{}
	}
	//fmt.Printf("[parser] put in %s -> %s\n", parentStr, parsedStr)
	siteM[parentStr][parsedStr] = true

	if !internal {
		return true
	}

	_, exists := siteM[parsedStr]
	if !exists {
		siteM[parsedStr] = map[string]bool{}
	}
	return exists
}

func (site *SiteMap) String() string {
	str := ""
	for k, vs := range *site {
		str += fmt.Sprintln(k, len(vs))
		for subk := range vs {
			str += fmt.Sprintf("\t%s\n", subk)
		}
	}
	return str
}

// DumpToFile dumps site.(SiteMap) to a file specified by filname
func (site *SiteMap) DumpToFile(filename string) error {
	f, err := os.Create("./" + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	line := ""
	inx := 0
	for k, vs := range *site {
		line = fmt.Sprintf("-> %d %s (outdegree: %d)\n", inx, k, len(vs))
		inx++

		_, err := w.WriteString(line)
		if err != nil {
			return err
		}

		for subk := range vs {
			line = fmt.Sprintf("\t%s\n", subk)
			_, err := w.WriteString(line)
			if err != nil {
				return err
			}
		}
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

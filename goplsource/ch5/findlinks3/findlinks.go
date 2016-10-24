// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 139.

// Findlinks3 crawls the web, starting with the URLs on the command line.
package main

import (
	"fmt"
	"log"
	"os"

	"gopl.io/ch5/links"
)

//!+breadthFirst
// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		fmt.Println("len of worklist: " , len(worklist))
		items := worklist
		worklist = nil
		for _, item := range items {

			//fmt.Println("len of seen: " , len(seen))
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

//!-breadthFirst

//!+crawl
func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

func crawl_ext(url string) []string{
	fmt.Println(url)
	s := "http://baidu.com.cn"
	s2 := "http://www.163.com.cn"
	var lks []string
	lks = append(lks, s, s2)
	return lks
}


//!-crawl

//!+main
func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl_ext, os.Args[1:])
}

//!-main

package main

import (
	"fmt"
	"net/http"
	"log"
	// "io"
	"golang.org/x/net/html"
	"strings"
	"os"
)

func Extract(url string, linkBuff chan string, errorChannel chan error)  { // ([]string, error)
	fmt.Printf("\nStarting : %s \n", url)
	resp, err := http.Get(url)
	if err != nil {
		errorChannel <- err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		errorChannel <- fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		errorChannel <- fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	// var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				// links = append(links, link.String())
				linkUrl := link.String()
				if i := strings.Index(linkUrl, "burntsushi.net"); i > -1 {
					// errorChannel <- fmt.Errorf("The %s", url, resp.Status)
					fmt.Printf("Good -> %s", linkUrl)
					linkBuff <- linkUrl
				}

				// fmt.Printf("The url %s is not supported", url)
				continue
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	// return links, nil
}

//!-Extract

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

// func treatLink(url string, ) {

// }


func main() {
	linkBuff := make(chan string)
	errorChannel := make(chan error)

	go Extract("https://blog.burntsushi.net/about/", linkBuff, errorChannel) // links, err := 

	for {
		select{
		case link := <- linkBuff:
			go Extract(link, linkBuff, errorChannel)
		case err := <- errorChannel:
			fmt.Printf("LOL!")
			log.Fatal(err)
			os.Exit(1)
		}
	}

	// if err == nil {
	// 	fmt.Printf("Printing:\n")
	// 	for _, link := range links {
	// 		fmt.Printf("%s\n", link)

	// 	}
	// } else {
	// 	log.Fatal(err)
	// }

//	if err != nil {
//		fmt.Printf("PANIC!")
//	}

//	defer resp.Body.Close()

//	for link := range getLinks(resp.Body) {
//		fmt.Printf(link)
//	}
//	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("DONE! \n")
}

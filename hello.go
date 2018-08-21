package main

import (
	"fmt"
	"net/http"
	// "log"
	// "io"
	"golang.org/x/net/html"
	"strings"
	// "os"
	"time"
	"sync"
)

func findInArray(wasAlready []string, search string) bool {
	found := false
	for i := range wasAlready {
		if wasAlready[i] == search {
			found = true
			break
		}
	}
	return found
}

func Extract(url string, linkBuff chan string, errorChannel chan error, syncChannel chan bool, wasAlready *[]string, wg *sync.WaitGroup)  { // ([]string, error)
	defer func() {<- syncChannel}()
	defer wg.Done()
	time.Sleep(time.Millisecond * 150)

	fmt.Printf("\nStarting : %s , length of channel %v, length of links:%v\n", url, len(syncChannel), len(linkBuff))
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("----->%s\n", err)
		// panic(err)
		errorChannel <- err
		return

	}

	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()

		errorChannel <- fmt.Errorf("getting %s: %s", url, resp.Status)
		return

	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		errorChannel <- fmt.Errorf("parsing %s as HTML: %v", url, err)
		return
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
				
				if indx:=strings.Index(linkUrl, "#"); indx > -1 {
					linkUrl = linkUrl[:indx]
				}

				if i := strings.Index(linkUrl, "burntsushi.net"); i > -1 {
					if !findInArray(*wasAlready, linkUrl){
						fmt.Printf("Good -> %s\n", linkUrl)
						
						*wasAlready = append(*wasAlready, url)

						linkBuff <- linkUrl
					}
				}

			}
		}
	}
	forEachNode(doc, visitNode, nil)
	// return links, nil
	return
}

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

func main() {
	linkBuff := make(chan string, 100)
	errorChannel := make(chan error)
	syncChannel := make(chan bool, 10)
	var wasAlready []string
	timeout := time.NewTimer(5 * time.Second)

	var wg sync.WaitGroup
    
	for i:=0;i<10;i++{
		fmt.Printf("\nInitial filling channel: %v", i)
		syncChannel<-true
	}

	wg.Add(1)
	go Extract("https://blog.burntsushi.net/about/", linkBuff, errorChannel, syncChannel, &wasAlready, &wg) // links, err := 

	L:
	for {
		select{
		case link := <- linkBuff:
			wg.Add(1)
			go Extract(link, linkBuff, errorChannel, syncChannel, &wasAlready, &wg)
			timeout.Reset(5 * time.Second)
		case err := <- errorChannel:			
			fmt.Printf("LOL! %s", err)
		case <- timeout.C:
			fmt.Printf("\n=================== TimeOut!\n")
			break L
		}

	}
	// wg.Wait()

	fmt.Printf("\nDONE! \n\n")

	for _, link := range wasAlready {
		fmt.Printf(" -> , %v", link)
	}

	fmt.Printf("\nUNDONE! \n\n")

}

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
			// fmt.Printf("Found: %s\n", search)
			found = true
			break
		}
	}
	return found
}

func Extract(url string, linkBuff chan string, errorChannel chan error, syncChannel chan bool, wasAlready []string, wg *sync.WaitGroup)  { // ([]string, error)
	defer func() {<- syncChannel}()
	defer func() {wg.Done()}()
	time.Sleep(time.Millisecond * 150)
	

	fmt.Printf("\nStarting : %s , length of channel %v, length of links:%v\n", url, len(syncChannel), len(linkBuff))
	resp, err := http.Get(url)

	// defer resp.Body.Close()

	// defer func() {resp.Body.Close()}()


	if err != nil {
		fmt.Printf("----->%s\n", err)
		// return

		// panic(err)
		errorChannel <- err
		return

	}

	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		wg.Done()
		// return

		errorChannel <- fmt.Errorf("getting %s: %s", url, resp.Status)
		return

	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		errorChannel <- fmt.Errorf("parsing %s as HTML: %v", url, err)
		return
	}

	// resp.Body.Close()
	// fmt.Printf("LOL %v", 6)
	// var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			// fmt.Printf("LOL %v", 7)
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
					// errorChannel <- fmt.Errorf("The %s", url, resp.Status)
					
					// if len(linkBuff) == cap(linkBuff) {
					// 	fullChannel <- true
					// }
					// fmt.Printf("...success!")
					if !findInArray(wasAlready, linkUrl){
						fmt.Printf("Good -> %s\n", linkUrl)
						wasAlready = append(wasAlready, url)
						// for _, link := range wasAlready {
						// 	fmt.Printf(" => %v", link)
						// }

						linkBuff <- linkUrl
					}
				}

				// fmt.Printf("The url %s is not supported", url)
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	// return links, nil
	return
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

// func refiller(syncChannel chan bool) {
// 	fmt.Printf("Starting sync worker")
// 	for {
		
// 		if len(syncChannel) < cap(syncChannel){
// 			fmt.Printf("\nlen(syncChannel):%v, is less than cap(syncChannel):%v,",len(syncChannel) ,cap(syncChannel))
// 			fmt.Printf("Refilling\n")
// 			syncChannel <- true
// 		} else{
// 			// fmt.Printf("Not,refilling\n")
// 		}
// 	}
// }


func main() {
	linkBuff := make(chan string, 100)
	errorChannel := make(chan error)
	syncChannel := make(chan bool, 10)
	var wasAlready []string

	var wg sync.WaitGroup
    

	for i:=0;i<10;i++{
		fmt.Printf("\nInitial filling channel: %v", i)
		syncChannel<-true
	}

	// go refiller(syncChannel)
	// fmt.Printf("\ncontinuint")


	wg.Add(1)
	go Extract("https://blog.burntsushi.net/about/", linkBuff, errorChannel, syncChannel, wasAlready, &wg) // links, err := 

	L:
	for {
		select{
		case link := <- linkBuff:
			wg.Add(1)
			go Extract(link, linkBuff, errorChannel, syncChannel, wasAlready, &wg)
		case err := <- errorChannel:			
			fmt.Printf("LOL! %s", err)
			// log.Fatal(err)
			break L

			// os.Exit(1)
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



	fmt.Printf("\nDONE! \n\n")

	for _, link := range wasAlready {
		fmt.Printf(" -> , %v", link)
	}

	fmt.Printf("\nUNDONE! \n\n")

}

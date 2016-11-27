package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	x, y, n, intPrime uint64
	z                 string
	algorithm         string
	writer            io.Writer
	primeSlice        = []uint64{}
	usageBool         bool
	ok                bool
	usetoml           bool
	//mapToPrimes       = map[uint64]Primers{}
	url = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
)

// Items is a list of Primes
type Items struct {
	Items []Item `json:"itmes"`
}

// Item is a list of Primes
type Item struct {
	Title       string  `json:"title"`
	UnitPrice   float32 `json:"unit_price"`
	Size        float32 `json:"size"`
	Description string  `json:"description"`
}

func init() {

	flag.BoolVar(&usageBool, "u", false, "Show the usage parameters.") //#3
}

func main() {

	flag.Parse()

	if usageBool {
		usage()
		os.Exit(0)
	}

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	// use CSS selector found with the browser inspector
	// for each, use index and item

	//fmt.Print(doc.Find("ul"))
	doc.Find("ul").Each(func(index int, ul *goquery.Selection) {
		//fmt.Printf("-ppopoppo-Post #%d: kkjssjs-item_as_text %s - \n", index, ul.Text())
		ul.Find("li").Each(func(i int, li *goquery.Selection) {
			//	linkTag := item.Find("a")
			h3Selector := li.Find("h3")
			aaSelector := h3Selector.Find("a")
			aatag, _ := aaSelector.Attr("href")

			if aatag != "" {
				fmt.Printf("-xxxxxxxll-Post #%d:\n hfjdhfdjhfd-aaSelector %s\n hgjdhdddjhdf-aatag %s\n", i, strings.TrimSpace(aaSelector.Text()), aatag)
				doc2, err := goquery.NewDocument(aatag)
				if err != nil {
					fmt.Println(err)
				}
				size, _ := getSize(aatag, "b")
				fmt.Printf("size in bytes:  -->%v<--, \n ", size)
				fmt.Printf("size in kilobytes:  -->%v<--, \n ", float32(size/1024))

				//	doc2.Find("div#addItem_149117").Each(func(ii int, div *goquery.Selection) {
				doc2.Find("[id^=addItem]").Each(func(ii int, div *goquery.Selection) {
					//fmt.Printf("detail page div  :  %v\n", div.Text())
					div.Find("p").Each(func(i int, pp *goquery.Selection) {
						pptag, _ := pp.Attr("class")
						ppstring := pp.Text()
						fmt.Printf("detail page b1:  %v, value -->%v<-- \n ", pptag, strings.TrimSpace(ppstring))
					})
				})
				doc2.Find("[class^=productTitleDescriptionContainer]").Each(func(ii int, div *goquery.Selection) {
					fmt.Printf("detail page div lll  :  %v\n", strings.TrimSpace(div.Text()))
					div.Find("p").Each(func(i int, pp *goquery.Selection) {
						pptag, _ := pp.Attr("class")
						ppstring := pp.Text()
						fmt.Printf("detail page x1:  %v, value -->%v<-- \n ", pptag, strings.TrimSpace(ppstring))
					})
				})

				doc2.Find("productcontent").Each(func(ii int, prod *goquery.Selection) {
					htmlselector := prod.Find("htmlcontent")
					htmlselector.Find("[class^=productDataItemHeader]").Each(func(ii int, item *goquery.Selection) {

						//fmt.Printf("itemselector xx:  -->%v<-- \n ", item.Text())

						if item.Text() == "Description" {
							itemSelector := item.NextUntil("h3")
							itemSelector.Each(func(ii int, desc *goquery.Selection) {
								//fmt.Printf("Description hurrah hurrah hurrah")
								fmt.Printf("desc :  -->%v<-- \n ", strings.TrimSpace(desc.Text()))
							})
						}

					})
					//fmt.Printf("itemselector :  -->%v<-- \n ", item.Text())

					//htmlselector.NextUntil(selector)
					//htmlselector := htmlselector.Find("[class^=productTitleDescriptionContainer]")

				})

			}

		})
		//		title := item.Text()
		//	linkTag := item.Find("a")
		//		link, _ := linkTag.Attr("href")

	})

	// scrape web page.
	//
	// if response, err := http.Get(url); err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	defer response.Body.Close()
	// 	_, err := io.Copy(os.Stdout, response.Body)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	fmt.Println("-->>Server Starting.....<<--")
	//r := mux.NewRouter()
	//r.HandleFunc("/primes/{algorithm}/{prime}", PrimeHandler)
	//r.HandleFunc("/primes/xml/{algorithm}/{prime}", PrimeXMLHandler)
	// Preload array with up to 5 million in background
	//go loadCache(sizeToCache)
	//	http.ListenAndServe(":8081", r)

}

func getSize(url string, unit string) (float32, error) {

	res, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		// handle error
	}
	l := len(body)
	return float32(l), nil

}

func usage() {

	fmt.Printf("\n\tUsage\n\t=====\n\ta Web Service to consume a webpage, process some data and present it.\n\n\t ")
	fmt.Printf("\n\tYou can choose the method of calculating the Prime numbers ; either the \"Sieve of Aitkin\" or the \"Sieve of Eratosthenes (Segmented)\t")
	fmt.Printf("\n\tTo Choose Aitkin the url format is http://your.host.com/primes/aitkin/15 ")
	fmt.Printf("\n\tTo Choose Eratosthenes the url format is http://your.host.com/primes/segmented/15 ")
	fmt.Printf("\n\tThe output Can also be represented as XML; ")
	fmt.Printf("\n\tThe URL for XML will be http://your.host.com/primes/xml/aitkin/15\n\n")

}

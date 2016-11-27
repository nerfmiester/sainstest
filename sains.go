package main

import (
	"encoding/json"
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
	UnitPrice   string  `json:"unit_price"`
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
	jItem := Item{}
	jItems := []Item{}
	// use CSS selector found with the browser inspector
	// for each, use index and item

	//fmt.Print(doc.Find("ul"))
	doc.Find("ul").Each(func(index int, ul *goquery.Selection) {
		ul.Find("li").Each(func(i int, li *goquery.Selection) {
			// Get the title
			h3Selector := li.Find("h3")
			aaSelector := h3Selector.Find("a")
			aatag, _ := aaSelector.Attr("href")
			if aatag != "" {
				jItem.Title = strings.TrimSpace(aaSelector.Text())

				// Get the size
				size, _ := getSize(aatag, "b")
				jItem.Size = float32(size / 1024)

				// Get the page with the details
				doc2, err := goquery.NewDocument(aatag)
				if err != nil {
					fmt.Println(err)
				}

				// Get unit price
				doc2.Find("[id^=addItem]").Each(func(ii int, div *goquery.Selection) {
					div.Find("p").Each(func(i int, pp *goquery.Selection) {
						pptag, _ := pp.Attr("class")
						if pptag == "pricePerUnit" {
							jItem.UnitPrice = strings.TrimSpace(pp.Text())
						}
					})
				})

				// Get the Description
				doc2.Find("productcontent").Each(func(ii int, prod *goquery.Selection) {
					htmlselector := prod.Find("htmlcontent")
					htmlselector.Find("[class^=productDataItemHeader]").Each(func(ii int, item *goquery.Selection) {
						if item.Text() == "Description" {
							itemSelector := item.NextUntil("h3")
							itemSelector.Each(func(ii int, desc *goquery.Selection) {
								jItem.Description = strings.TrimSpace(desc.Text())
							})
						}
					})
				})

				jItems = append(jItems, jItem)
				jItem = Item{}

			}
		})
	})

	res2B, _ := json.Marshal(jItems)
	fmt.Println(string(res2B))
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

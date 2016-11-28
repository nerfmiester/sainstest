package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	x, y, n, intPrime uint64
	itemStruct        = Items{}
	totUnitPrice      = float64(0)
	doc, doc2         *goquery.Document
	err               error
	usageBool         bool
	jsonItems         []byte
	size              string
	//mapToPrimes       = map[uint64]Primers{}
	url = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
)

// Items is a list of Primes
type Items struct {
	Items          []Item  `json:"results"`
	TotalUnitPrice float64 `json:"total"`
}

// Item is a list of Primes
type Item struct {
	Title       string  `json:"title"`
	Size        string  `json:"size"`
	UnitPrice   float64 `json:"unit_price"`
	Description string  `json:"description"`
}
type findInDocT func(index int, sel *goquery.Selection) *goquery.Selection

func init() {

	flag.BoolVar(&usageBool, "u", false, "Show the usage parameters.") //#3
}

func main() {

	flag.Parse()

	if usageBool {
		usage()
		os.Exit(0)
	}

	doc, err = goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	jItem := Item{}
	jItems := []Item{}

	// use CSS selector found with the browser inspector
	// for each, use index and item

	doc.Find("[id^=productLister]").Each(func(index int, ul *goquery.Selection) {
		ul.Find("li").Each(func(i int, li *goquery.Selection) {
			// Get the title
			h3Selector := li.Find("h3")
			aaSelector := h3Selector.Find("a")
			aatag, _ := aaSelector.Attr("href")
			if aatag != "" {
				jItem.Title = strings.TrimSpace(aaSelector.Text())

				// Get the size
				size, err = getSize(aatag, "b")
				if err != nil {
					fmt.Println(err)
				}
				jItem.Size = size

				// Get the page with the details
				doc2, err = goquery.NewDocument(aatag)
				if err != nil {
					fmt.Println(err)
				}

				// Get unit price
				doc2.Find("[id^=addItem]").Each(func(ii int, itemPrice *goquery.Selection) {
					itemPrice.Find("[class$=pricePerUnit]").Each(func(i int, pp *goquery.Selection) {

						pptag, _ := pp.Attr("class")

						if pptag == "pricePerUnit" {
							jItem.UnitPrice, err = getPrice(strings.TrimSpace(pp.Text()))
							totUnitPrice += jItem.UnitPrice
							if err != nil {
								fmt.Println(err)
							}
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
	itemStruct.Items = jItems
	itemStruct.TotalUnitPrice = toFixed(totUnitPrice, 2)

	jsonItems, err = jsonEncoder(itemStruct)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Returned JSON doc is -->%v<--", string(jsonItems))
	}

}

func getSize(url string, unit string) (string, error) {

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
	return strconv.Itoa(round(float64(len(body))/1024)) + unit, nil
}

func jsonEncoder(jItems Items) ([]byte, error) {
	jsonItems, err := json.Marshal(jItems)
	if err != nil {
		return nil, err
	}
	return jsonItems, nil

}
func getPrice(s string) (float64, error) {
	sFloat, err := strconv.ParseFloat(s[(strings.Index(s, "£")+2):strings.LastIndex(s, "/")], 64)
	if err != nil {
		return 0.0, err
	}

	return sFloat, nil

}
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func usage() {

	fmt.Printf("\n\tUsage\n\t=====\n\ta Web Service to consume a webpage, process some data and present it.\n\n\t ")
	fmt.Printf("\n\tYou can choose the method of calculating the Prime numbers ; either the \"Sieve of Aitkin\" or the \"Sieve of Eratosthenes (Segmented)\t")
	fmt.Printf("\n\tTo Choose Aitkin the url format is http://your.host.com/primes/aitkin/15 ")
	fmt.Printf("\n\tTo Choose Eratosthenes the url format is http://your.host.com/primes/segmented/15 ")
	fmt.Printf("\n\tThe output Can also be represented as XML; ")
	fmt.Printf("\n\tThe URL for XML will be http://your.host.com/primes/xml/aitkin/15\n\n")

}

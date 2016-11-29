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
	resp         *http.Response
	itemStruct   = Items{}
	totUnitPrice = float64(0)
	doc, doc2    *goquery.Document
	sFloat       float64
	err          error
	usageBool    bool
	jsonItems    []byte
	size         string
	url          = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
)

// Items is a list of Items and the total Unit Price of the items
type Items struct {
	Items          []Item  `json:"results"`
	TotalUnitPrice float64 `json:"total"`
}

// Item is a description of an Item
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
	// process external flags
	flag.Parse()
	// if only requiring to see the usage.
	if usageBool {
		usage()
		os.Exit(0)
	}
	// Get the document to process based on a url, here hard coded but could be an external list.
	doc, err = goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	// Get the returned Jason and display it
	jsonVal, errP := process(doc)
	if errP != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Returned JSON doc is -->%v<--", string(jsonVal))
	}

}
func process(doc *goquery.Document) ([]byte, error) {

	// define the empty structures
	jItem := Item{}
	jItems := []Item{}
	// using the qoquery package search for elements in the document.
	// See the goquery documentation for a clearer explination - but suffice to say like jQuery
	doc.Find("[id^=productLister]").Each(func(index int, ul *goquery.Selection) {
		ul.Find("li").Each(func(i int, li *goquery.Selection) {
			// Get the title
			h3Selector := li.Find("h3")
			aaSelector := h3Selector.Find("a")
			aatag, _ := aaSelector.Attr("href")
			if aatag != "" {
				jItem.Title = strings.TrimSpace(aaSelector.Text())

				// Get the size
				resp, err = http.Get(aatag)
				body, _ := ioutil.ReadAll(resp.Body)
				size, err = getSize(body, "b")
				if err != nil {
					fmt.Println(err)
				}
				defer resp.Body.Close()
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
	// Encode into JSON
	jsonItems, err = jsonEncoder(itemStruct)

	return jsonItems, nil

}

// Function to return the size of a body of a page
func getSize(body []byte, unit string) (string, error) {
	return strconv.Itoa(round(float64(len(body))/1024)) + unit, nil
}

// Encode to JSON use empty interface for reuse?!
func jsonEncoder(jItems interface{}) ([]byte, error) {
	jsonItems, err = json.Marshal(jItems)
	if err != nil {
		return nil, err
	}
	return jsonItems, nil

}

// Extract the price and convert to a float64
func getPrice(s string) (float64, error) {
	sFloat, err = strconv.ParseFloat(s[(strings.Index(s, "£")+2):strings.LastIndex(s, "/")], 64)
	if err != nil {
		return 0.0, err
	}

	return sFloat, nil

}

// Change value to a fixed precision
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// Use math package to round correctly.
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// Show usage.
func usage() {

	fmt.Printf("\n\tUsage\n\t=====\n\ta Web Service to consume a webpage, process some data and present it as JSON output.\n\n\t ")

}

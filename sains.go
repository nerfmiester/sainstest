package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/caleblloyd/primesieve"
	"github.com/gorilla/mux"
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
	mapToPrimes       = map[uint64]Primers{}
	url               = "http://hiring-tests.s3-website-eu-west-1.amazonaws.com/2015_Developer_Scrape/5_products.html"
)

// Primers is a list of Primes
type Primers struct {
	Initial string   `json:"initial"`
	Primes  []uint64 `json:"primes"`
}

// Xprimes is a list of Primes
type Xprimes struct {
	XMLName xml.Name `xml:"primeNumbers"`
	Initial string   `xml:"initial"`
	Primes  []uint64 `xml:"primes"`
}

func init() {

	flag.BoolVar(&usageBool, "u", false, "Show the usage parameters.") //#3
}

// Filter Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan uint64, out chan<- uint64, prime uint64) {
	for {
		i := <-in // Receive value from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
		}
	}
}

const sizeToCache = 50000000

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
				//	doc2.Find("div#addItem_149117").Each(func(ii int, div *goquery.Selection) {
				doc2.Find("[id^=addItem]").Each(func(ii int, div *goquery.Selection) {
					//fmt.Printf("detail page div  :  %v\n", div.Text())
					div.Find("p").Each(func(i int, pp *goquery.Selection) {
						pptag, _ := pp.Attr("class")
						ppstring := pp.Text()
						fmt.Printf("detail page :  %v, value -->%v<-- \n ", pptag, strings.TrimSpace(ppstring))
					})
				})
				doc2.Find("[class^=productTitleDescriptionContainer]").Each(func(ii int, div *goquery.Selection) {
					fmt.Printf("detail page div lll  :  %v\n", strings.TrimSpace(div.Text()))
					div.Find("p").Each(func(i int, pp *goquery.Selection) {
						pptag, _ := pp.Attr("class")
						ppstring := pp.Text()
						fmt.Printf("detail page :  %v, value -->%v<-- \n ", pptag, strings.TrimSpace(ppstring))
					})

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
	r := mux.NewRouter()
	r.HandleFunc("/primes/{algorithm}/{prime}", PrimeHandler)
	r.HandleFunc("/primes/xml/{algorithm}/{prime}", PrimeXMLHandler)
	// Preload array with up to 5 million in background
	go loadCache(sizeToCache)
	http.ListenAndServe(":8081", r)

}

// PrimeHandler = the function will Orchistrate prime number creation and return JSON.
func PrimeHandler(w http.ResponseWriter, r *http.Request) {
	val := Primers{}
	err := errors.New("")
	prime := mux.Vars(r)["prime"]
	algorithm := mux.Vars(r)["algorithm"]
	if intPrime, err = strconv.ParseUint(prime, 10, 64); err != nil {
		fmt.Println(err)
	}
	if intPrime <= sizeToCache {
		if val, ok = mapToPrimes[intPrime]; ok {
		} else {
			if algorithm == "segmented" {
				val = workerSegmented(intPrime)
			} else {
				val = workerAitkin(intPrime)
			}
		}
	} else {
		val = workerAitkin(intPrime)
	}
	j, err := json.Marshal(val)

	if err != nil {
		fmt.Println(err)
	}
	w.Write([]byte(j))

}

// PrimeXMLHandler = the function will Orchistrate prime number creation and return XML notation.
func PrimeXMLHandler(w http.ResponseWriter, r *http.Request) {
	val := Primers{}
	top := Xprimes{}
	err := errors.New("")
	prime := mux.Vars(r)["prime"]
	if intPrime, err = strconv.ParseUint(prime, 10, 64); err != nil {
		fmt.Println(err)
	}
	if intPrime <= sizeToCache {
		if val, ok = mapToPrimes[intPrime]; ok {
		} else {
			val = workerAitkin(intPrime)
		}
	} else {
		val = workerAitkin(intPrime)
	}
	top.Initial = prime
	top.Primes = val.Primes

	output, err := xml.Marshal(&top)
	if err != nil {
		fmt.Println("Error marshling to xml", err)
		return
	}

	w.Write([]byte(xml.Header + string(output)))

}

func loadCache(size uint64) {
	for i := uint64(1); i <= size; i++ {
		mapToPrimes[i] = workerAitkin(i)
	}
}

func workerSegmented(toPrime uint64) Primers {
	primers := Primers{}
	primers.Initial = strconv.FormatUint(toPrime, 10)
	primers.Primes = primesieve.ListMax(uint64(toPrime))
	return primers
}

func workerAitkin(toPrime uint64) Primers {

	// Many thanks to gofool for this implementation
	// https://raw.githubusercontent.com/agis-/gofool/master/atkin.go

	var x, y, n uint64
	nsqrt := math.Sqrt(float64(toPrime))

	isPrime := make([]bool, (sizeToCache + 1))
	for x = 1; float64(x) <= nsqrt; x++ {
		for y = 1; float64(y) <= nsqrt; y++ {
			n = 4*(x*x) + y*y
			if n <= toPrime && (n%12 == 1 || n%12 == 5) {
				isPrime[n] = !isPrime[n]
			}
			n = 3*(x*x) + y*y
			if n <= toPrime && n%12 == 7 {
				isPrime[n] = !isPrime[n]
			}
			n = 3*(x*x) - y*y
			if x > y && n <= toPrime && n%12 == 11 {
				isPrime[n] = !isPrime[n]
			}
		}
	}

	for n = 5; float64(n) <= nsqrt; n++ {
		if isPrime[n] {
			for y = n * n; y < toPrime; y += n * n {
				isPrime[y] = false
			}
		}
	}
	//fmt.Println("len of isPrime = ", len(isPrime))
	isPrime[2] = true
	isPrime[3] = true

	primes := make([]uint64, 0, 1270606)
	for x = 0; x < uint64(len(isPrime))-1; x++ {
		if isPrime[x] {
			primes = append(primes, x)
		}
	}

	// primes is now a slice that contains all the
	// primes numbers up to isPrime

	primers := Primers{}
	primers.Initial = strconv.FormatUint(toPrime, 10)
	primers.Primes = primes

	return primers

}

func usage() {

	fmt.Printf("\n\tUsage\n\t=====\n\ta Web Service to consume a webpage, process some data and present it.\n\n\t ")
	fmt.Printf("\n\tYou can choose the method of calculating the Prime numbers ; either the \"Sieve of Aitkin\" or the \"Sieve of Eratosthenes (Segmented)\t")
	fmt.Printf("\n\tTo Choose Aitkin the url format is http://your.host.com/primes/aitkin/15 ")
	fmt.Printf("\n\tTo Choose Eratosthenes the url format is http://your.host.com/primes/segmented/15 ")
	fmt.Printf("\n\tThe output Can also be represented as XML; ")
	fmt.Printf("\n\tThe URL for XML will be http://your.host.com/primes/xml/aitkin/15\n\n")

}

//

package main

import (
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func LookupPrice() {
	fileInfos, err := ioutil.ReadDir(dataDirFlag)
	if err != nil {
		log.Fatal(err)
	}
	for _, fileInfo := range fileInfos {
		if !strings.HasPrefix(fileInfo.Name(), "doomtrooper") {
			continue
		}
		file, err := os.Open(filepath.Join(dataDirFlag, fileInfo.Name()))
		if err != nil {
			log.Fatal(err)
		}

		doc, err := goquery.NewDocumentFromReader(file)
		if err != nil {
			log.Fatal(err)
		}

		// Find the review items
		doc.Find("#listingTable").Find("tr").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			titleLink, ok := s.Find(".td-title").Find("a").Attr("href")
			if !ok {
				return
			}
			state := strings.TrimSpace(s.Find(".td-state").Text())
			if strings.Contains(titleLink, cardFlag) && state == "Prodáno" {
				title := strings.TrimSpace(s.Find(".td-title").Text())
				dateRE := regexp.MustCompile(`(?m)^.*Datum\ vložení\ do\ archivu:(.*)$`)
				dateMatch := dateRE.FindStringSubmatch(title)
				sellerRE := regexp.MustCompile(`(?m)^.*Prodejce:(.*)$`)
				sellerMatch := sellerRE.FindStringSubmatch(title)
				price := strings.TrimSpace(s.Find(".td-price").Text())
				fmt.Printf("%s - %s - %s\n", price, strings.TrimSpace(dateMatch[1]), strings.TrimSpace(sellerMatch[1]))
			}
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

var (
	cardFlag string
	dataDirFlag string
)


func init() {
	flag.StringVar(&cardFlag, "card", "", "the card the lookup prices for; lower-case, dash-separate (e.g. komunikacni-sum)")
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	flag.StringVar(&dataDirFlag, "data-dir", wd, "the directory to process data from")
}

func main() {
	flag.Parse()
	if cardFlag == "" {

	}
	LookupPrice()
}

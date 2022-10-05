package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func check(error error) {
	if error != nil {
		fmt.Println(error)
	}
}

func getHTML(url string) *http.Response {
	res, error := http.Get(url)
	check(error)

	if res.StatusCode > 400 {
		fmt.Println("status code:", res.StatusCode)
	}

	return res
}

func writeCSV(scrapedData []string) {
	filename := "ebayData.csv"

	file, error := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	check(error)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	error = writer.Write(scrapedData)
	check(error)
}

func scrapePageData(doc *goquery.Document) {
	doc.Find("ul.srp-results>li.s-item").Each(func(index int, item *goquery.Selection) {
		a := item.Find("a.s-item__link")

		title := strings.TrimSpace(a.Text())
		url, _ := a.Attr("href")

		priceSpan := strings.TrimSpace(item.Find("span.s-item__price").Text())
		price := strings.Trim("$ ", priceSpan)

		scrapedData := []string{title, price, url}

		writeCSV(scrapedData)
	})
}

func main() {
	url := "https://www.ebay.com/sch/i.html?_from=R40&_trksid=p2334524.m570.l2632&_nkw=funko+pop+marvel+lot&_sacat=246&LH_TitleDesc=0&_odkw=funko+pop+marvel&_osacat=0"

	var previousUrl string

	for {
		res := getHTML(url)
		defer res.Body.Close()

		doc, error := goquery.NewDocumentFromReader(res.Body)
		check(error)

		scrapePageData(doc)

		href, _ := doc.Find("nav.pagination > a.padination__next").Attr("href")

		if href == previousUrl {
			break
		} else {
			url = href
			previousUrl = href
		}
	}

}

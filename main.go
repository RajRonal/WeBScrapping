package main

import (
	"encoding/csv"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

func checkError(e error) {
	if e != nil {
		logrus.Error(e.Error())
	}
}
func writeCsv(scrapedData []string) {
	filename := "output.csv"
	file, fileError := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	checkError(fileError)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writeError := writer.Write(scrapedData)
	checkError(writeError)

}

// FUNCTION TO SCRAPE DATA FROM WEBSITE
func scrapePageData(doc *goquery.Document) {
	doc.Find("ul.srp-results > li.s-item").Each(func(index int, item *goquery.Selection) {
		a := item.Find("a.s-item__link")
		title := strings.TrimSpace(a.Text())
		//fmt.Println(title)
		url, _ := a.Attr("href")
		priceSpan := item.Find("span.s-item__price").Text()
		price := strings.TrimSpace(priceSpan)
		scrapedData := []string{title, price, url}
		writeCsv(scrapedData)
	})

	return
}

//REQUEST THE HTML PAGE
func getHtml(url string) *http.Response {
	response, resError := http.Get(url)
	checkError(resError)
	if response.StatusCode != 200 {
		logrus.Error("status code error: %d %s", response.StatusCode, response.Status)

	}
	return response
}

func main() {
	url := "https://www.ebay.com/sch/i.html?_from=R40&_nkw=beatles+puzzle&_sacat=0&LH_TitleDesc=0&rt=nc&_ipg=120/"
	//response := getHtml(url)
	response, resError := http.Get(url)
	checkError(resError)
	if response.StatusCode != 200 {
		logrus.Error("status code error: %d %s", response.StatusCode, response.Status)

	}
	//LOAD THE HTML PAGE
	document, err := goquery.NewDocumentFromReader(response.Body)
	checkError(err)
	scrapePageData(document)
	defer response.Body.Close()
}

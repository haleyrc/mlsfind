package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: mlsfind MLSID")
		os.Exit(2)
	}

	url := fmt.Sprintf("http://www.slcmls.com/Search/ListingDetail.aspx?org_id=nystlc&mls_property_id=%s", os.Args[1])
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Error from slcmls.com: %d - %s\n", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		log.Fatalln(err)
	}

	// doc.Find("#photo-carousel").Each(func(i int, s *goquery.Selection) {
	// 	if s == nil {
	// 		log.Fatalln("S was nil for:", string(b))
	// 		os.Exit(1)
	// 	}
	// 	s.Children().Each(func(j int, s *goquery.Selection) {
	// 		val, ok := s.Children().First().Attr("src")
	// 		if !ok {
	// 			return
	// 		}
	// 		fmt.Println(val)
	// 	})
	// })

	doc.Find(".details-text-data").Each(func(i int, s *goquery.Selection) {
		label := s.Children().First()
		if label.Text() == "Tax Amount:" {
			fmt.Println(label.Text(), label.SiblingsFiltered("span").First().Text())
		}
	})
}

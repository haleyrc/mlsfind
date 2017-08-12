package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mlsfind MLSID...")
		os.Exit(2)
	}

	for i := 1; i < len(os.Args); i++ {
		print(os.Args[i])
		fmt.Println()
	}
}

func print(mls string) {
	url := fmt.Sprintf("http://www.slcmls.com/Search/ListingDetail.aspx?org_id=nystlc&mls_property_id=%s", mls)
	fmt.Println("MLS#:\t\t\t", mls)
	fmt.Println("Link:\t\t\t", url)

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

	fmt.Println("Price:\t\t\t", doc.Find(".price").First().Text())
	fmt.Println("Square Feet:\t\t", strings.TrimSuffix(doc.Find(".sqft").First().Text(), " sqft"))
	fmt.Println("Bedrooms/Bathrooms:\t", doc.Find(".bed-baths").First().Children().Remove().End().Text())

	doc.Find(".details-text-data").Each(func(i int, s *goquery.Selection) {
		label := s.Children().First()
		switch label.Text() {
		case "Tax Amount:", "Lot Size:":
			fmt.Println(label.Text()+"\t\t", label.SiblingsFiltered("span").First().Text())
		case "Tax Assessment:":
			fmt.Printf(label.Text()+"\t\t $%s\n", strings.TrimPrefix(label.Parent().Text(), label.Text()))
		case "Year Built:":
			fmt.Println(label.Text()+"\t\t", strings.TrimPrefix(label.Parent().Text(), label.Text()))
		}
	})
}

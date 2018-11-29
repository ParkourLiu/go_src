package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"fmt"
)

func ExampleScrape() {
	doc, err := goquery.NewDocument("http://www.qiushibaike.com") //http://www.qiushibaike.com
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".article").Each(func(i int, s *goquery.Selection) {
		if s.Find(".thumb").Nodes == nil && s.Find(".video_holder").Nodes == nil {
			content := s.Find(".content").Text()
			list := []rune(content)
			for i, v := range list {
				if i%50 == 0 {
					fmt.Println()
				}
				if string(v) != "\n" {
					fmt.Print(string(v))
				} else if i < len(list)-1 {
					if string(v) == "\n" && string(list[i+1]) == "\n" {
						fmt.Println()
					}
				}

			}
		}
	})
}

func main() {
	ExampleScrape()
}

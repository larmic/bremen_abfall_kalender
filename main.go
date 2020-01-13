package main

import (
	"fmt"
	"log"

	"github.com/schollz/progressbar/v2"

	"bremen_trash/client"
	"bremen_trash/html/bremen/stadtreinigung"
	"bremen_trash/html/bremen/stadtreinigung/parser"
)

var (
	bremerStadtreinigungRootUrl  = "http://213.168.213.236/bremereb/bify/"
	bremerStadtreinigungIndexUrl = bremerStadtreinigungRootUrl + "index.jsp"
	c                            = client.NewClient()
)

func main() {
	// Load first letters
	fmt.Println("Loading street first letters")
	bar := progressbar.NewOptions(1, progressbar.OptionSetRenderBlankState(true))
	firstLetters := loadFirstLetterOfStreets()
	_ = bar.Finish()
	fmt.Println()
	fmt.Printf("%d street first letters loaded", len(firstLetters))
	fmt.Println()
	fmt.Println()

	// Load streets
	fmt.Println("Loading streets")
	bar = progressbar.NewOptions(3700, progressbar.OptionSetRenderBlankState(true))
	streets := loadStreets(firstLetters, bar)
	_ = bar.Finish()
	fmt.Println()
	fmt.Printf("%d streets loaded", len(streets))
	fmt.Println()
	fmt.Println()

	type GarbageCollectionUrl struct {
		Street               string
		HouseNumber          string
		GarbageCollectionUrl string
	}

	fmt.Println("Loading house numbers")
	garbageCollectionUrls := make([]GarbageCollectionUrl, 0)
	bar = progressbar.NewOptions(len(streets), progressbar.OptionSetRenderBlankState(true))
	for _, street := range streets {
		content, err := c.GetContent(street.Url)

		if err != nil {
			log.Fatal(err)
		}

		houseNumbers := stadtreinigung.ParseHouseNumber(content, bremerStadtreinigungRootUrl)

		for _, houseNumber := range houseNumbers {
			garbageCollectionUrls = append(garbageCollectionUrls, GarbageCollectionUrl{street.Value, houseNumber.Value, houseNumber.Url})
		}

		_ = bar.Add(1)
		//loadDates(houseNumbers, content)
	}

	_ = bar.Finish()
	fmt.Println()
	fmt.Printf("%d house numbers loaded", len(garbageCollectionUrls))
	fmt.Println()
	fmt.Println()

	fmt.Println("Loading garbage collection dates")
	bar = progressbar.NewOptions(len(garbageCollectionUrls), progressbar.OptionSetRenderBlankState(true))
	dates := 0
	for _, garbageCollectionUrl := range garbageCollectionUrls {
		gcd := loadDates(garbageCollectionUrl.GarbageCollectionUrl)
		_ = bar.Add(1)
		dates = dates + len(gcd)
	}
	_ = bar.Finish()
	fmt.Println()
	fmt.Printf("%d garbage collection dates loaded", dates)
	fmt.Println()
}

func loadDates(garbageCollectionUrl string) []stadtreinigung.GarageCollection {
	garbageContent, err := c.GetContent(garbageCollectionUrl)

	if err != nil {
		log.Fatal(err)
	}

	if garbageContent == "" {
		log.Fatal(`Dates content is empty. Url: `, garbageCollectionUrl)
	}

	return stadtreinigung.ParseGarbageCollectionDates(garbageContent)
}

func loadStreets(firstLetters []parser.Dto, bar *progressbar.ProgressBar) []parser.Dto {
	streets := make([]parser.Dto, 0)

	for _, firstLetter := range firstLetters {
		//fmt.Println(`Found first letter of street`, firstLetter.FirstLetter, firstLetter.Url)

		content, err := c.GetContent(firstLetter.Url)

		if err != nil {
			log.Fatal(err)
		}

		firstLetterStreets := stadtreinigung.ParseStreetPage(content, firstLetter, bremerStadtreinigungRootUrl)

		//if err != nil {
		//	fmt.Printf(`Error while parsing streets of %s. Error is '%s'. Url will be ignored.`, firstLetter.Url, err)
		//	fmt.Println()
		//}

		for _, element := range firstLetterStreets {
			streets = append(streets, element)
		}

		_ = bar.Add(len(firstLetterStreets))
	}
	return streets
}

func loadFirstLetterOfStreets() []parser.Dto {
	content, err := c.GetContent(bremerStadtreinigungIndexUrl)

	if err != nil {
		log.Fatal(err)
	}

	firstLetters := stadtreinigung.ParseIndexPage(content, bremerStadtreinigungRootUrl)
	return firstLetters
}

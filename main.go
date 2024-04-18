package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/gocolly/colly"
)

func scrap(imgs *[]string, source string) {
	f, err := os.Create("collector.txt")

	if err != nil {
		panic(err)
	}

	defer f.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	scrapper := colly.NewCollector()

	scrapper.OnHTML("img", func(e *colly.HTMLElement) {
		if len(e.Attr("src")) >= 5 && e.Attr("src")[0:5] == "https" {
			//fmt.Println("Found the image:", e.Attr("src"))
			f.WriteString(e.Attr("src") + "\n")

			*imgs = append(*imgs, e.Attr("src"))
		}else {
			
			url_parsed, err := url.Parse(source)
			if err != nil {
				panic(err)
			}

			f.WriteString("https://" + url_parsed.Host + e.Attr("src") + "\n")

			*imgs = append(*imgs, "https://" + url_parsed.Host + e.Attr("src"))
		}
	})

	scrapper.OnScraped(func(r *colly.Response) {
		defer wg.Done()
	})


	scrapper.Visit(source)

	wg.Wait()
}

func main() {
	imgs := []string{}
	var source string
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tplt := template.Must(template.ParseFiles("index.html"))

		tplt.Execute(w, nil)
	})

	http.HandleFunc("/href", func(w http.ResponseWriter, r *http.Request) {
		source = r.PostFormValue("link")
		
		tpltAnother := template.Must(template.ParseFiles("output.html"))
		fmt.Println(source)

		scrap(&imgs, source)
		
		fmt.Println(imgs)

		tpltAnother.Execute(w, nil)
	})

	http.HandleFunc("/get-button", func(w http.ResponseWriter, r *http.Request) {
		
	})
	http.ListenAndServe(":8080", nil)
}
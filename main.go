package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"sync"
	"io"
	"log"

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
		
		//fmt.Println(imgs)

		tpltAnother.Execute(w, nil)
	})

	http.HandleFunc("/get-button", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Println("HIIIII")

		// err1 := os.RemoveAll("/images")
		// if err1 != nil {
		// 	log.Fatal("Cannot delete this folder :(")
		// }

		err := os.Mkdir("images", 0755)

		if err != nil {
			log.Fatal("Cannot find this path for folder :(")
		}

		for ind, img_path := range imgs {
			var name string
			if string(img_path[len(img_path) - 4]) == "." {
				name = fmt.Sprintf("imgoo_%d%s", ind, img_path[len(img_path) - 4:])
			} else {
				name = fmt.Sprintf("imgoo_%d.jpg", ind)
			}
			file, err := os.Create(fmt.Sprintf("images/%s", name))
			if err != nil {
				log.Fatal("Cannot create this file :(")
			}

			client := http.Client{
				CheckRedirect: func(r *http.Request, via []*http.Request) error {
					r.URL.Opaque = r.URL.Path
					return nil
				},
			}
			
			resp, err := client.Get(img_path)
			if err != nil {
				log.Fatal("Cannot get this client :(")
			}
			defer resp.Body.Close()

			size, err := io.Copy(file, resp.Body)
			if err != nil {
				log.Fatal("Cannot copy size :(")
			}
 
    		defer file.Close()
 
    		fmt.Printf("Downloaded a file %s with size %d\n", name, size)
			fmt.Fprintf(w, "<strong>Images downloaded</strong><br>")

		}
	})
	http.ListenAndServe(":8080", nil)
}
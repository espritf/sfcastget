package main

import (
	"github.com/gocolly/colly"
	"log"
	"os"
	"strings"
)

const download_dir = "download"

func main() {

    email := os.Args[1]
    password := os.Args[2]
    course_path := os.Args[3]

    l := colly.NewCollector()
    l.OnHTML("input[name=_csrf_token]", func(e *colly.HTMLElement) {
        c := l.Clone()
        err := c.Post("https://symfonycasts.com/login", map[string]string{
            "email":       email,
            "password":    password,
            "_csrf_token": e.Attr("value"),
        })

        if err != nil {
            log.Fatal(err)
        }

        c.OnHTML("ul.chapter-list li a[href]", func(e *colly.HTMLElement) {
            url := e.Request.AbsoluteURL(e.Attr("href") + "/download/video")
            log.Println(url)

            path := download_dir + "/" + course_path
            err := os.MkdirAll(path, os.ModePerm)
            if err != nil {
                log.Fatal(err)
            }

            d := c.Clone()

            d.OnResponse(func(r *colly.Response) {
                if strings.Contains(r.Headers.Get("Content-Type"), "video/mp4") {
                    log.Println("Download", r.Request.URL)

                    filename := r.FileName()
                    download_filename := path + "/" + filename

                    log.Println("Download path", download_filename)
                    r.Save(download_filename)
                }
            })

            d.Visit(url)
        })

        c.Visit("https://symfonycasts.com/screencast/" + course_path)

    })
    l.Visit("https://symfonycasts.com/login")
}

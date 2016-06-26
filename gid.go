//script to download imgur image(s), given a url
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type ImageFilename struct {
	Url      string
	FileName string
}

func getLinksFromAlbum(imgURL *url.URL, dir string, ordered bool) ([]ImageFilename, error) {
	images := make([]ImageFilename, 0)
	doc, err := goquery.NewDocument(imgURL.String())
	if err != nil {
		return images, err
	}

	doc.Find(".post-images").Find(".post-image-container").Each(func(i int, s *goquery.Selection) {
		id, exists := s.Attr("id")
		if !exists {
			return
		}

		var fileName = fmt.Sprintf("%s.jpg", id)
		if ordered {
			fileName = fmt.Sprintf("%s-%d-%s.jpg", imgURL.Path[3:], i, id)
		}

		images = append(images, ImageFilename{fmt.Sprintf("http://i.imgur.com/%s.jpg", id), fileName})
	})

	return images, nil
}

func downloadImage(url, dir, name string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Downloading...", url, name)
	return ioutil.WriteFile(fmt.Sprintf("%s/%s", dir, name), contents, 0644)
}

func main() {
	//parse args for url and directory
	urlPtr := flag.String("url", "", "imgur url with images")
	dirPtr := flag.String("d", "images", "directory to download images to")
	orderPtr := flag.Bool("o", true, "keep ordering of images via naming")
	flag.Parse()

	//parse imgur URL
	imgURL, err := url.Parse(*urlPtr)
	if err != nil {
		fmt.Println("bad URL, please try again.")
		panic(err)
	}

	//check if its imgur
	//TODO: make this more robust
	if !strings.Contains(imgURL.String(), "imgur.com") {
		fmt.Print("not an imgur link, please try again.")
		os.Exit(0)
	}

	//downgrade to http in case https link is given for performance
	imgURL.Scheme = "https"

	//determine whether it's an album or just a single image
	switch {
	case imgURL.Path[:3] == "/a/":
		log.Println("Downloading album into", *dirPtr)
		images, err := getLinksFromAlbum(imgURL, *dirPtr, *orderPtr)
		if err != nil {
			log.Fatalf("Error downloading image, error: %s", err)
		}

		var wg sync.WaitGroup
		for _, u := range images {
			wg.Add(1)
			go func(url, fileName string) {
				defer wg.Done()
				downloadImage(url, *dirPtr, fileName)
			}(u.Url, u.FileName)
		}
		wg.Wait()
		log.Println("Download success!")
	default:
		log.Println("Downloading image into", *dirPtr)
		filename := path.Base(imgURL.String())

		if filename == "" {
			log.Fatal("could not resolve filename from url")
		}

		err := downloadImage(imgURL.String(), *dirPtr, filename)

		if os.IsNotExist(err) {
			//make directory if not exist
			os.Mkdir(*dirPtr, 0644)
			err = downloadImage(imgURL.String(), *dirPtr, filename)
		}

		if err != nil {
			log.Fatalf("Error downloading image, error: %s", err)
		}

		log.Println("Download success!")
	}
}

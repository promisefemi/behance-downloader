package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
)

//"strconv"
//"strings"

type Author struct {
	Name, ProfileLink string
}
type BehanceImage struct {
	FileName, URL string
}
type Project struct {
	Author       Author
	ProjectTitle string
	Images       []BehanceImage
}

func ProcessLink(url string) (Project, error) {
	project := Project{}

	resp, err := http.Get(url)
	if err != nil {
		return project, err
	}

	if resp.StatusCode != 200 {
		return project, fmt.Errorf("Cannot reach page, please check URL")
	}

	defer resp.Body.Close()

	// Initialize HTML Parser

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	project.ProjectTitle = doc.Find("title").Text()

	projectModules := doc.Find("#project-modules img")
	authorSelection := doc.Find(".Project-ownerName-A8O")
	project.Author = getAuthor(authorSelection)
	projectModules.Each(func(i int, s *goquery.Selection) {
		parsedImage, err := getImage(s)
		if err == nil {
			project.Images = append(project.Images, parsedImage)
		}
	})
	return project, nil
}

func ProcessImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(resp.Status)
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return responseBody, nil
}

func getAuthor(authorSelection *goquery.Selection) Author {
	author := Author{}

	profileLink, ffd := authorSelection.Attr("href")
	if ffd {
		author.ProfileLink = profileLink
		name := authorSelection.Find("span").Text()
		if name != "" {
			author.Name = name
		}
	}
	return author
}

func getImage(s *goquery.Selection) (BehanceImage, error) {
	img := BehanceImage{}

	imgURL, isThere := s.Attr("data-src")
	if !isThere {
		imgURL, isThere = s.Attr("src")
	}

	if isThere {
		fileName := getFileName(imgURL)
		if fileName == "blank.png" {
			return img, fmt.Errorf("Image not found")
		} else {
			img.FileName = fileName
			img.URL = imgURL
		}
	} else {
		return img, fmt.Errorf("Image not found")
	}

	return img, nil
}

func getFileName(imgURL string) string {
	urlPath, _ := url.Parse(imgURL)
	fileName := path.Base(urlPath.Path)
	return fileName
}

func main() {

	projectChan := make(chan *Project)
	behanceImage := &Project{}
	go func() {
		project, err := ProcessLink("https://www.behance.net/gallery/154852101/Hamburg-Noir-II?tracking_source=search_projects")
		if err != nil {
			//fmt.Println(err)
		}
		projectChan <- &project
	}()

breakLabel:
	for {
		for _, r := range `-\|/` {
			fmt.Printf("%s", fmt.Sprintf("\r%c", r))
			time.Sleep(time.Duration(1) * time.Second)
		}
		select {
		case behanceImage = <-projectChan:
			fmt.Printf("\r\n")
			break breakLabel
		default:
		}
	}

	item, err := json.MarshalIndent(behanceImage, "   ", "")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%s", item)
}

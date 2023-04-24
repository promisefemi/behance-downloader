package core

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

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
	Author            Author
	ProjectTitle, URL string
	Images            []BehanceImage
}

func ProcessLink(url string) (Project, error) {
	//Process Link and return a Project object and error.
	//Could return error or project

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
	if err != nil {
		return project, err
	}
	project.ProjectTitle = strings.ReplaceAll(doc.Find("title").Text(), "on BehanceAdobe", "")

	projectModules := doc.Find("#project-modules img")
	authorSelection := doc.Find(".Project-ownerName-A8O")
	project.Author = getAuthor(authorSelection)
	projectModules.Each(func(i int, s *goquery.Selection) {
		parsedImage, err := getImage(s)
		if err == nil {
			project.Images = append(project.Images, parsedImage)
		}
	})
	project.URL = url
	return project, nil
}

// func ProcessImage(url string) ([]byte, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf(resp.Status)
// 	}
// 	responseBody, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return responseBody, nil
// }

func ProcessDownload(url string) ([]byte, string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, "", "", fmt.Errorf("Cannot reach page, please check URL")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", "", err
	}
	return body, resp.Header.Get("Content-Type"), resp.Header.Get("Content-Length"), nil
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
	// urlPath, _ := url.Parse(imgURL)
	fileName := path.Base(imgURL)
	return fileName
}

package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func ProcessLink(url string) (Project, error) {
	//Process Link and return a Project object and error.
	//Could return error or project

	project := Project{}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, url, nil)
	check(err)
	request.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	request.Header.Set("Sec-Fetch-Site", "same-origin")
	request.Header.Set("Sec-Fetch-Mode", "navigate")
	request.Header.Set("Sec-Fetch-Dest", "document")

	resp, err := client.Do(request)
	if err != nil {
		return project, err
	}

	if resp.StatusCode != 200 {
		return project, fmt.Errorf("Cannot reach page, please check URL")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	check(err)
	file, err := os.OpenFile("file.html", os.O_RDONLY|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	check(err)

	_, err = file.Write(body)
	check(err)

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

func printToTerminal(payload interface{}) {
	if _, yes := payload.(string); yes {
		fmt.Println(payload)
	}

	if _, yes := payload.([]byte); yes {
		fmt.Printf("%s\n", payload)
	}

	if _, yes := payload.(int); yes {
		fmt.Println(payload)
	}

	data, err := json.MarshalIndent(payload, " ", "  ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

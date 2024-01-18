package core

import (
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/promisefemi/behance-downloader/core"
	"github.com/spf13/cobra"
)

var list bool
var author bool
var source string

var rootCommand = &cobra.Command{
	Use:   "behance",
	Short: "Scrapper and Downloader for Behance projects",
	Long:  `A lightweight scrapper and downloader for Behance projects, github at https://github.com/promisefemi/behance-downloader	`,
	Args: func(cmd *cobra.Command, arg []string) error {

		if len(arg) <= 0 {
			return fmt.Errorf("minimum of one argument required")
		}

		for _, link := range arg {
			parsedUrl, err := url.Parse(link)
			if err != nil {
				return fmt.Errorf("some of the urls are invalid ")
			}

			if parsedUrl.Host != "www.behance.net" {
				return fmt.Errorf("invalid URL, URL is not from Behance")
			}
		}

		return nil

	},
	Run: func(cmd *cobra.Command, arg []string) {
		// totalwait := 4
		done := make(chan bool)
		possibleErrors := make([]error, 0)
		wg := sync.WaitGroup{}

		for _, link := range arg {
			wg.Add(1)
			go func(link string) {
				behanceImage, err := ProcessLink(link)
				if err != nil {
					fmt.Printf("\r\n")
					exit(fmt.Sprintf("%s", err))
				}

				if list || author {
					if list {
						ListDetails(behanceImage)
					} else if author {
						PrintAuthor(behanceImage.Author, true)
					}
					wg.Done()
				} else {
					fmt.Printf("\rDownloading %s...\n", behanceImage.ProjectTitle)
					destination := ""
					if source != "" {
						destination = fmt.Sprintf("%s/%s", source, behanceImage.ProjectTitle)
					} else {
						destination = behanceImage.ProjectTitle
					}
					err := CreateFolder(destination)
					if err != nil {
						possibleErrors = append(possibleErrors, fmt.Errorf("error creating folder - '%s'", err))
						wg.Done()
						return
					}
					downloadWait := sync.WaitGroup{}
					downloadWait.Add(len(behanceImage.Images))
					for _, image := range behanceImage.Images {
						go func(image core.BehanceImage) {
							err := DownloadImage(destination, image, &downloadWait)
							if err != nil {
								possibleErrors = append(possibleErrors, fmt.Errorf("error downloading file  - '%s'", image.FileName))
								return
							}
							fmt.Printf("\rDownloaded %s\n", image.FileName)
						}(image)
					}
					downloadWait.Wait()
					wg.Done()
				}
			}(link)
		}

		go func() {
			wg.Wait()
			done <- true
			return
		}()

	breakLabel:
		for {
			select {
			case <-done:
				break breakLabel
			default:
			}
			for _, r := range `-\|/` {
				fmt.Printf("%s", fmt.Sprintf("\r%c", r))
				// time.Sleep(time.Duration(1) * time.Second)
			}
		}
		fmt.Printf("\r")
		if list || author {
			fmt.Printf("%s \n", "Detail listed successfully...")

		} else {
			if len(possibleErrors) > 0 {
				fmt.Println(possibleErrors)
			}
			// fmt.Printf("\r %s \n", "All Projects Downloaded Successfully...")
		}
		// ProcessLink()
		// fmt.Println(behanceImage)

	},
}

func Execute() {
	rootCommand.Flags().BoolVarP(&list, "list", "l", false, "List all images in the project")
	rootCommand.Flags().BoolVarP(&author, "author", "a", false, "Get the name and link of the project author")
	rootCommand.Flags().StringVarP(&source, "src", "s", "", "File destination")
	if err := rootCommand.Execute(); err != nil {
		exit(fmt.Sprintf("%s", err))
	}
}

func exit(error string) {
	fmt.Println(error)
	os.Exit(1)
}

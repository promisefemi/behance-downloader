package cli

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var list bool
var author bool

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
				return fmt.Errorf("invalid Url, Url is not from Behance")
			}
		}

		return nil

	},
	Run: func(cmd *cobra.Command, arg []string) {
		// totalwait := 4
		done := make(chan bool)
		wg := sync.WaitGroup{}

		for _, link := range arg {
			wg.Add(1)
			go func(link string) {
				behance := ProcessLink(link, &wg)
				fmt.Println(behance)
			}(link)
		}

		go func() {
			wg.Wait()
			done <- true
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
				time.Sleep(time.Duration(1) * time.Second)
			}
		}
		fmt.Printf("\r %s \n", "All Projects Downloaded Successfully...")
		// ProcessLink()
		// fmt.Println(behanceImage)

	},
}

func Execute() {
	rootCommand.Flags().BoolVarP(&list, "list", "l", false, "List all images in the project")
	rootCommand.Flags().BoolVarP(&author, "author", "a", false, "Get the name and link of the project author")
	if err := rootCommand.Execute(); err != nil {
		exit(fmt.Sprintf("%s", err))
	}
}

func exit(error string) {
	fmt.Println(error)
	os.Exit(1)
}

package cli

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/promisefemi/behance-downloader/core"
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "behance",
	Short: "Scrapper and Downloader for Behance projects",
	Long:  `A lightweight scrapper and downloader for Behance projects, github at https://github.com/promisefemi/behance-downloader	`,
	Run: func(cmd *cobra.Command, arg []string) {
		if len(arg) <= 0 {
			exit("No links parsed")
		}
		behanceURL := arg[0]

		parsedUrl, err := url.Parse(behanceURL)
		if err != nil {
			exit("Invalid URL")
		}

		if parsedUrl.Host != "www.behance.net" {
			exit("Invalid Url, Url is not from Behance")
		}

		projectChan := make(chan *core.Project)
		errorChan := make(chan error)
		behanceImage := &core.Project{}
		go func() {
			project, err := core.ProcessLink(behanceURL)
			if err != nil {
				errorChan <- err
			} else {
				projectChan <- &project
			}
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
			case err = <-errorChan:
				fmt.Printf("\r\n")
				exit(fmt.Sprintf("%s", err))
				break breakLabel
			default:
			}
		}

		fmt.Println(behanceImage)

	},
}

func Execute() {
	if err := rootCommand.Execute(); err != nil {
		exit(fmt.Sprintf("%s", err))
	}
}

func exit(error string) {
	fmt.Println(error)
	os.Exit(1)
}

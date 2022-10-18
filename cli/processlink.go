package cli

import (
	"fmt"
	"sync"

	"github.com/promisefemi/behance-downloader/core"
)

func ProcessLink(link string, wg *sync.WaitGroup) *core.Project {
	projectChan := make(chan *core.Project)
	errorChan := make(chan error)
	// behanceImage := &core.Project{}

	go func() {
		project, err := core.ProcessLink(link)
		if err != nil {
			errorChan <- err
		} else {
			projectChan <- &project
		}
	}()

	for {
		select {
		case behanceImage := <-projectChan:
			fmt.Printf("\r\n")
			wg.Done()
			return behanceImage
			// break
		case err := <-errorChan:
			fmt.Printf("\r\n")
			exit(fmt.Sprintf("%s", err))
		}
	}

	// fmt.Println("asdfasdh")
	// behanceImage := <-projectChan
	// fmt.Println("asdfasdh")

	// err := <-errorChan

	// if err != nil {
	// 	fmt.Printf("\r\n")
	// 	exit(fmt.Sprintf("%s", err))
	// }

	// if behanceImage.Author.Name != "" {
	// 	fmt.Printf("\r\n")
	// 	wg.Done()
	// 	return behanceImage
	// }

	// return nil

}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"log"
	"os"
	"strings"
)

var client *aspace.ASClient

func main() {
	outFile, _ := os.Create("handle-updates.txt")
	defer outFile.Close()
	log.SetOutput(outFile)

	var err error
	client, err = aspace.NewClient("go-aspace.yml", "dev", 20)
	if err != nil {
		panic(err)
	}

	for _, i := range []int{2, 3, 6} {
		getDOs(i)
	}
}

func getDOs(repoID int) {
	iresults, err := client.Search(repoID, "digital_object", "*", 1)
	if err != nil {
		panic(err)
	}

	for page := 1; page <= iresults.LastPage; page++ {
		fmt.Printf("Repository %d Page %d of %d\n", repoID, page, iresults.LastPage)

		pageResults, err := client.Search(repoID, "digital_object", "*", page)
		if err != nil {
			panic(err)
		}
		for _, result := range pageResults.Results {
			do := aspace.DigitalObject{}
			err = json.Unmarshal([]byte(fmt.Sprint(result["json"])), &do)
			if err != nil {
				panic(err)
			}
			if containsHandle(do.FileVersions) == true {
				//repoID, doID, err := aspace.URISplit(do.URI)

				newFV := []aspace.FileVersion{}
				for _, fv := range do.FileVersions {
					if strings.Contains(fv.FileURI, "http://hdl.handle.net/2333.1/") == true {
						fv.FileURI = strings.ReplaceAll(fv.FileURI, "http", "https")
						newFV = append(newFV, fv)
					} else {
						newFV = append(newFV, fv)
					}
				}
				do.FileVersions = newFV
				fmt.Printf("[INFO] %s %v\n", do.URI, do.FileVersions)
			}
		}
	}
}

func containsHandle(fileVersions []aspace.FileVersion) bool {
	for _, fv := range fileVersions {
		if strings.Contains(fv.FileURI, "http://hdl.handle.net/2333.1/") == true {
			return true
		}
	}
	return false
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/nyudlts/go-aspace"
	"log"
	"os"
	"strings"
)

func main() {
	client, err := aspace.NewClient("go-aspace.yml", "dev", 20)
	if err != nil {
		panic(err)
	}

	iresults, err := client.Search(2, "digital_object", "*", 1)
	if err != nil {
		panic(err)
	}

	outFile, _ := os.Create("handles.txt")
	defer outFile.Close()
	log.SetOutput(outFile)

	for i := 1; i <= iresults.LastPage; i++ {
		fmt.Printf("Page %d of %d\n", i, iresults.LastPage)
		numUpdates := 0
		page, err := client.Search(3, "digital_object", "*", i)
		if err != nil {
			panic(err)
		}

		for _, result := range page.Results {
			do := aspace.DigitalObject{}
			err = json.Unmarshal([]byte(fmt.Sprint(result["json"])), &do)
			if err != nil {
				panic(err)
			}

			repoID, doID, err := aspace.URISplit(do.URI)
			if err != nil {
				panic(err)
			}

			fileVersions := do.FileVersions
			if containsHandle(fileVersions) == true {
				numUpdates = numUpdates + 1
				do.FileVersions = updateFileVersions(fileVersions)
				msg, err := client.UpdateDigitalObject(repoID, doID, do)
				if err != nil {
					log.Printf("[ERROR] %s\n", strings.ReplaceAll(err.Error(), "\n", ""))
					continue
				}
				log.Printf("[INFO] %s\n", strings.ReplaceAll(msg, "\n", ""))
			}
		}

		fmt.Println(numUpdates, " DOs updated")
	}
}

func containsHandle(fileVersions []aspace.FileVersion) bool {
	for _, fv := range fileVersions {
		if strings.Contains(fv.FileURI, "hdl.handle") == true {
			return true
		}
	}
	return false
}

func updateFileVersions(fileVersions []aspace.FileVersion) []aspace.FileVersion {
	newFileVersions := []aspace.FileVersion{}
	for _, fv := range fileVersions {
		if strings.Contains(fv.FileURI, "hdl.handle") == true {
			fv.FileURI = strings.ReplaceAll(fv.FileURI, "http", "https")
			newFileVersions = append(fileVersions, fv)
		} else {
			newFileVersions = append(fileVersions, fv)
		}
	}
	return newFileVersions
}

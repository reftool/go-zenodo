package main

import (
	"fmt"
	"time"

	"github.com/reftool/gozenodo"
)

func main() {
	gozenodo.SetAccessToken("h4kmNNh9QL2ZD1YqiVW8MSqZapbAmvriroH2Bh57sjRk3SWkVjgzKO7Gn2c1")
	gozenodo.SetSandboxMode(true)

	// Create blank deposition
	deposition, err := gozenodo.CreateDeposition()
	if err != nil {
		panic(err)
	}
	fmt.Println("Created Deposition with ID: ", deposition.ID)

	// Add a file to the deposition bucket
	newFileUpload, err := gozenodo.UploadFile(deposition.Links.Bucket, "test.txt", "examples/test.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("Uploaded file: Key=", newFileUpload.Key)

	// Update the title of the deposition
	deposition.Title = "test title 123"
	deposition, err = gozenodo.UpdateDeposition(deposition)
	if err != nil {
		panic(err)
	}
	fmt.Println("New Deposition Title: ", deposition.Title)

	// Get deposition by id
	test, err := gozenodo.GetDeposition(deposition.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Found Deposition by ID: ", test.ID)

	// Wait for a while because zenodo takes its time to update
	fmt.Println("waiting a few seconds...")
	time.Sleep(5 * time.Second)

	// List all of your depositions
	all, err := gozenodo.ListDepositions()
	if err != nil {
		panic(err)
	}

	for _, d := range all {
		fmt.Println("Found Deposition: ", d.ID)
	}

	// Delete deposition by id
	err = gozenodo.DeleteDeposition(deposition.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("Deleted deposition by id: ", deposition.ID)
}

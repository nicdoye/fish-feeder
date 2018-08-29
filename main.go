package main

import (
	"fmt"
	"os"

	"github.com/nicdoye/fish-feeder/cmd"
)

func main() {

	if err := cmd.Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	/*
		tempFile, err := ioutil.TempFile("", TempFilePrefix)
		if err != nil {
			log.Fatal(err)
		}

		// errcheck moans about this
		defer os.Remove(tempFile.Name())

		// TODO sprintf + commandline
		fileURL := "https://github.com/gohugoio/hugo/releases/download/v0.40.3/hugo_0.40.3_checksums.txt"

		err = DownloadFile(tempFile, fileURL)
		if err != nil {
			panic(err)
		}

		PackageMap := make(map[string]gofish.Package)
		err = createPackageMap(tempFile, &PackageMap)

		if err != nil {
			panic(err)
		}
		fmt.Printf("%v", PackageMap)

		var fish = Fish{
			Name:        Name,
			Description: Description,
			License:     License,
			Homepage:    Homepage,
			Caveats:     Caveats,
			Version:     Version,
			PackageMap:  PackageMap,
		}

		writeFish(fish)
	*/
}

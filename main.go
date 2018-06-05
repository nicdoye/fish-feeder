package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"

	"github.com/fishworks/gofish"
)

// OSType ...
type OSType string

// OSFoodName ...
type OSFoodName string

// OSInfo ...
type OSInfo struct {
	Type     OSType
	FoodName OSFoodName
}

// Fish Fishy MvFishFace
type Fish struct {
	Name        string
	Description string
	License     string
	Homepage    string
	Caveats     string
	Version     string

	PackageMap map[string]gofish.Package
}

// Enum
const (
	unix    OSType = "unix"
	windows OSType = "windows"
)

// Enum - FN suffix is a temporary hack
const (
	dragonflybsdFN OSFoodName = "dragonflybsd"
	freebsdFN      OSFoodName = "freebsd"
	linuxFN        OSFoodName = "linux"
	netbsdFN       OSFoodName = "netbsd"
	openbsdFN      OSFoodName = "openbsd"
	windowsFN      OSFoodName = "windows"
	// This one is darwin, not macos
	macosFN OSFoodName = "darwin"
)

var osMap = map[string]OSInfo{
	"DragonFlyBSD": {Type: unix, FoodName: dragonflybsdFN},
	"FreeBSD":      {Type: unix, FoodName: freebsdFN},
	"Linux":        {Type: unix, FoodName: linuxFN},
	"NetBSD":       {Type: unix, FoodName: netbsdFN},
	"OpenBSD":      {Type: unix, FoodName: openbsdFN},
	"Windows":      {Type: windows, FoodName: windowsFN},
	"macOS":        {Type: unix, FoodName: macosFN},
}

var archMap = map[string]string{
	"32bit": "386",
	"64bit": "amd64",
	"ARM":   "arm",
	"ARM64": "arm64",
}

// Name ...
const Name string = "hugo"

// Description ...
const Description = "The world’s fastest framework for building websites."

// License ...License
const License = "Apache-2.0"

// Homepage ...Homepage
const Homepage = "https://gohugo.io/"

// Caveats ...
const Caveats = ""

// Version ...
const Version = "0.40.3"

var resources = map[OSType]gofish.Resource{
	unix:    {Path: "name", InstallPath: "\"bin\" .. name", Executable: true},
	windows: {Path: "name .. \".exe\"", InstallPath: "\"bin\\\\\" .. name .. \".exe\"", Executable: true},
}

// DownloadFile ...
func DownloadFile(file *os.File, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	// errcheck moans about this
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)

	return err
}

func main() {
	tempFile, err := ioutil.TempFile("", "hugo-updater-")
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
}

func writeFish(fish Fish) {
	/*
				macOS-64bit.tag.gz -> os = macos


				UNIX
				path = name,
				installpath = "bin/" .. name,
				executable = true

				Windows
				path = name .. ".exe",
		    installpath = "bin\\" .. name .. ".exe"
	*/
	t := template.Must(template.ParseFiles("hugo.tpl"))
	err := t.Execute(os.Stdout, fish)
	if err != nil {
		panic(err)
	}
}

func makePackage(fileName string, sha256 string, packagePtr *gofish.Package) error {
	//	fileName is software_version_osString-archString.[tar.gz|deb|zip]

	osArchSuffixStr := regexp.MustCompile(`(_+)`).Split(fileName, 3)[2]
	osArchStr := regexp.MustCompile(`(.+)`).Split(osArchSuffixStr, 2)[0]
	osArchPair := regexp.MustCompile(`(-+)`).Split(osArchStr, 2)

	osInfo := osMap[osArchPair[0]]
	resource := resources[osInfo.Type]
	var resourceArrayPtr []*gofish.Resource

	resourceArrayPtr = append(resourceArrayPtr, &resource)

	*packagePtr = gofish.Package{
		OS:        (string)(osInfo.Type),
		Arch:      archMap[osArchPair[1]],
		Resources: resourceArrayPtr,
		URL:       "", // "https://github.com/gohugoio/" .. name .. "/releases/download/v" .. version .. "/" .. "hugo_" .. version .. "_macOS-64bit.tar.gz",
		Mirrors:   nil,
		SHA256:    sha256,
	}

	return nil
}

func createPackageMap(tempFile io.ReadSeeker, PackageMapPtr *map[string]gofish.Package) error {
	// Rewind
	_, err := tempFile.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(tempFile)
	for scanner.Scan() {
		reversedPair := regexp.MustCompile(`( +)`).Split(scanner.Text(), 2)
		var sha256 = reversedPair[0]
		var packageFileName = reversedPair[1]
		var packagePackage = gofish.Package{}
		err = makePackage(packageFileName, sha256, &packagePackage)

		if err == nil {
			(*PackageMapPtr)[packageFileName] = packagePackage
		}
	}

	return scanner.Err()
}

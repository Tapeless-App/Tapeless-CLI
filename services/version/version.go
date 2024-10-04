package version

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"tapeless.app/tapeless-cli/env"
	"tapeless.app/tapeless-cli/util"
)

var ()

type GithubReleaseResponse struct {
	TagName string `json:"tag_name"`
}

func CheckLatestVersion() {

	installedVersion := env.Version

	if installedVersion == "Development" {
		return
	}

	var versionResponse GithubReleaseResponse

	util.MakeUnAuthRequestAndParseResponse("GET", "https://api.github.com/repos/Tapeless-App/Tapeless-CLI/releases/latest", nil, &versionResponse)

	latestVersion := versionResponse.TagName

	if shouldPerformUpdate(installedVersion, latestVersion) {
		fmt.Println("You may want to update the CLI to the latest version.")
		switch os := runtime.GOOS; os {
		case "darwin", "linux":
			fmt.Println("If you installed the CLI using Homebrew, run 'brew update && brew upgrade tapeless'")
			fmt.Println("Or download the latest version from 'https://github.com/Tapeless-App/Tapeless-CLI/releases/latest'")
		default:
			fmt.Println("Please download the latest version from 'https://github.com/Tapeless-App/Tapeless-CLI/releases/latest'")
		}
	}

}

func shouldPerformUpdate(installedVersion string, latestVersion string) bool {
	// Remove "v" prefix if present
	installedVersionTrimmed := strings.TrimPrefix(installedVersion, "v")
	latestVersionTrimmed := strings.TrimPrefix(latestVersion, "v")

	// Split versions into components
	installedVersionParts := strings.Split(installedVersionTrimmed, ".")
	latestVersionParts := strings.Split(latestVersionTrimmed, ".")

	// Find the longest length between the two versions
	if len(installedVersionParts) != len(latestVersionParts) {
		fmt.Println("Cannot compare versions with different format: installed version: ", installedVersion, " latest version: ", latestVersion)
		return true
	}

	if len(installedVersionParts) < 3 || len(latestVersionParts) < 3 {
		fmt.Println("Version should be major, minor and patch - cannot apply format to installed version: ", installedVersion, " or latest version: ", latestVersion)
		return true
	}

	installedMajor, err := strconv.Atoi(installedVersionParts[0])
	if err != nil {
		fmt.Println("Cannot convert installed major version to integer: ", installedVersionParts[0])
		return true
	}

	lastestMajor, err := strconv.Atoi(latestVersionParts[0])
	if err != nil {
		fmt.Println("Cannot convert latest major version to integer: ", latestVersionParts[0])
		return true
	}

	if lastestMajor > installedMajor {
		fmt.Println("New major version. Installed version: ", installedVersion, " latest version: ", latestVersion)
		return true
	}

	installedMinor, err := strconv.Atoi(installedVersionParts[1])

	if err != nil {
		fmt.Println("Cannot convert installed minor version to integer: ", installedVersionParts[1])
		return true
	}

	latestMinor, err := strconv.Atoi(latestVersionParts[1])

	if err != nil {
		fmt.Println("Cannot convert latest minor version to integer: ", latestVersionParts[1])
		return true
	}

	if lastestMajor == installedMajor && latestMinor > installedMinor {
		fmt.Println("New minor version. Installed version: ", installedVersion, " latest version: ", latestVersion)
		return true
	}

	installedPatch, err := strconv.Atoi(installedVersionParts[2])

	if err != nil {
		fmt.Println("Cannot convert installed patch version to integer: ", installedVersionParts[2])
		return true
	}

	lastestPatch, err := strconv.Atoi(latestVersionParts[2])

	if err != nil {
		fmt.Println("Cannot convert latest patch version to integer: ", latestVersionParts[2])
		return true
	}

	if lastestMajor == installedMajor && latestMinor == installedMinor && lastestPatch > installedPatch {
		fmt.Println("New patch version. Installed version: ", installedVersion, " latest version: ", latestVersion)
		return true
	}

	return false

}

package utils

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	WowAddonsDirectory string   `json:"wow_addon_directory"`
	ElvuiFolders       []string `json:"elvui_folders"`
	API                string   `json:"api"`
}

type WorldOfWarcraft struct {
	Directory string
	Addons    []Addon
}

type Addon struct {
	Folder string
}

type API struct {
	API   string
	Elvui Elvui
}

type Elvui struct {
	Slug        string   `json:"slug" bson:"slug"`
	Url         string   `json:"url" bson:"url"`
	Version     string   `json:"version" bson:"version"`
	Directories []string `json:"directories" bson:"directories"`
}

type TocFile struct {
	Version string
}

func (c Config) CheckConfig() error {
	if c.API == "" {
		return errors.New("API url not present in config file")
	}

	if len(c.ElvuiFolders) == 0 {
		return errors.New("Elvui folder not present in config file")
	}

	if c.WowAddonsDirectory == "" {
		return errors.New("Wow addon directory not present in config file")
	}

	return nil
}

func IsElvuiPresent(cf Config, ad []Addon) (bool, []Addon) {
	var b bool
	var addons []Addon

	for _, a := range ad {
		for _, f := range cf.ElvuiFolders {
			if a.Folder == f {
				b = true
				addons = append(addons, Addon{Folder: a.Folder})
			}
		}
	}

	return b, addons
}

func (w *WorldOfWarcraft) PopulateWowAddons() error {
	f, err := os.Open(w.Directory)
	if err != nil {
		return err
	}

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	for _, v := range files {
		if v.IsDir() {
			w.Addons = append(w.Addons, Addon{
				Folder: v.Name(),
			})
		}
	}

	return nil
}

func (a Addon) GetVersion(cf Config) (string, error) {
	tocFile := filepath.Join(cf.WowAddonsDirectory, a.Folder, a.Folder+"_Mainline.toc")
	f, err := os.Open(tocFile)
	if err != nil {
		return "", fmt.Errorf("cannot open TOC file %s: %w", tocFile, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	const versionPrefix = "## Version:"
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, versionPrefix) {
			version := strings.TrimSpace(strings.TrimPrefix(line, versionPrefix))
			version = strings.TrimPrefix(version, "v")
			if version == "" {
				return "", errors.New("version found but empty")
			}
			return version, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return "", errors.New("version not found in TOC file")
}

func (api *API) GetElvuiInformation() error {
	res, err := http.Get(api.API)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return fmt.Errorf("invalid HTTP status: %d", res.StatusCode)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var elvuis []Elvui
	if err := json.Unmarshal(resBody, &elvuis); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(elvuis) == 0 {
		return errors.New("no data received from API")
	}

	for _, elvui := range elvuis {
		if elvui.Slug == "elvui" {
			api.Elvui = elvui
			return nil
		}
	}

	return errors.New("ElvUI not found in API data")
}

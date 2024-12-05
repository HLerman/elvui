package utils

import (
	"bufio"
	"encoding/json"
	"errors"
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
	f, err := os.Open(filepath.Join(cf.WowAddonsDirectory, a.Folder, a.Folder+"_Mainline.toc"))
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "##") {
			line = strings.TrimPrefix(line, "##")
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				if key == "Version" {
					if value[0:1] == "v" {
						return value[1:], nil
					}
					return value, nil
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", errors.New("no version")
}

func (api *API) GetElvuiInformation() error {
	res, err := http.Get(api.API)
	if err != nil {
		return err
	}

	if res.StatusCode > 299 {
		return err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var elvuis []Elvui
	err = json.Unmarshal(resBody, &elvuis)
	if err != nil {
		return err
	}

	if len(elvuis) == 0 {
		return errors.New("no data from API")
	}

	var id *int
	for i, e := range elvuis {
		if e.Slug == "elvui" {
			id = &i
		}
	}

	if id == nil {
		return errors.New("no elvui from api")
	}

	api.Elvui = elvuis[*id]
	return nil
}

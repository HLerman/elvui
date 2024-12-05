package main

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/HLerman/elvui/utils"
	"github.com/hashicorp/go-version"
)

func main() {
	log, err := os.Create("log_elvui.log")
	defer log.Close()
	if err != nil {
		fmt.Fprintf(log, "error when creatin log file: %s\n", err)
		os.Exit(1)
	}

	o, err := os.ReadFile("settings.json")
	if err != nil {
		fmt.Fprintf(log, "error when opening config file: %s\n", err)
		os.Exit(1)
	}

	var cf utils.Config
	err = json.Unmarshal(o, &cf)
	if err != nil {
		fmt.Fprintf(log, "error when Unmarshal config file: %s\n", err)
		os.Exit(1)
	}

	if err = cf.CheckConfig(); err != nil {
		fmt.Fprintf(log, err.Error())
		os.Exit(1)
	}

	fmt.Fprintf(log, "Wow addons directory: %s\nAPI: %s\n", cf.WowAddonsDirectory, cf.API)

	wow := utils.WorldOfWarcraft{
		Directory: cf.WowAddonsDirectory,
	}

	api := utils.API{
		API: cf.API,
	}

	err = api.GetElvuiInformation()
	if err != nil {
		fmt.Fprintf(log, "error making http request: %s\n", err)
		os.Exit(1)
	}

	err = wow.PopulateWowAddons()
	if err != nil {
		fmt.Fprintf(log, "error during populate: %s\n", err)
		os.Exit(1)
	}

	var isOudated bool
	if present, addons := utils.IsElvuiPresent(cf, wow.Addons); present {
		for _, a := range addons {
			v, err := a.GetVersion(cf)
			if err != nil {
				fmt.Fprintf(log, "error getVersion: %s\n", err)
				os.Exit(1)
			}

			inGame, err := version.NewVersion(v)
			if err != nil {
				fmt.Fprintf(log, "error set version: %s\n", err)
				os.Exit(1)
			}

			inSite, err := version.NewVersion(api.Elvui.Version)
			if err != nil {
				fmt.Fprintf(log, "error set version: %s\n", err)
				os.Exit(1)
			}

			fmt.Fprintln(log, "addon:", a.Folder, "in_game_version:", inGame.String())
			fmt.Fprintln(log, "addon:", a.Folder, "site_version:", inSite.String())
			if inGame.LessThan(inSite) {
				isOudated = true
				fmt.Fprintln(log, "Outdated !")
				break
			}
		}

		if isOudated {
			// remove elvui folders
			for _, f := range api.Elvui.Directories {
				err := os.RemoveAll(filepath.Join(cf.WowAddonsDirectory, f))
				if err != nil {
					fmt.Fprintf(log, "error when delete folder: %s\n", err)
					os.Exit(1)
				}

			}

			// download newest version of elvui
			err := DownloadElvui(cf, api.Elvui)
			if err != nil {
				fmt.Fprintln(log, err)
				os.Exit(1)
			}
		}
	}

	if !isOudated {
		fmt.Fprintln(log, "nothing to do, addon is up-to-date")
	}
}

func DownloadElvui(cf utils.Config, elvui utils.Elvui) error {
	elvuiFile := filepath.Join(os.Getenv("TEMP"), "elvui.zip")
	out, err := os.Create(elvuiFile)
	if err != nil {
		return errors.Join(errors.New("error when creatin elvui zip file:"), err)
	}
	defer out.Close()

	resp, err := http.Get(elvui.Url)
	if err != nil {
		return errors.Join(errors.New("error when downloading elvui zip file:"), err)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("error when downloading elvui zip file: http status:%s", resp.Status))
	}

	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Join(errors.New("error when copying elvui zip file:"), err)
	}

	zip, err := zip.OpenReader(elvuiFile)
	if err != nil {
		return errors.Join(errors.New("error when opening elvui zip file:"), err)
	}

	defer zip.Close()

	for _, f := range zip.File {
		rc, err := f.Open()
		if err != nil {
			return errors.Join(errors.New("error when opening file in zip file:"), err)
		}

		defer rc.Close()

		if f.FileInfo().IsDir() {
			err := os.Mkdir(filepath.Join(cf.WowAddonsDirectory, f.Name), 0755)
			if err != nil {
				return errors.Join(errors.New("error when creatin folder:"), err)
			}

			continue
		}

		out, err := os.Create(filepath.Join(cf.WowAddonsDirectory, f.Name))
		if err != nil {
			return errors.Join(errors.New("error when creatin file:"), err)
		}
		defer out.Close()

		_, err = io.Copy(out, rc)
		if err != nil {
			return errors.Join(errors.New("error when copying elvui file to the destination folder:"), err)
		}
	}

	return nil
}

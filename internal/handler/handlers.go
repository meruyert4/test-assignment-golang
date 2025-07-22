package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"test-assignment/internal/archiver"
	"test-assignment/internal/config"
	"test-assignment/internal/resolver"
	"test-assignment/internal/sshclient"
)

func HandleCreate(configPath string) error {
	if !isLikelyJSON(configPath) {
		return fmt.Errorf("file is not a json")
	}
	cfg, err := config.ParsePacketConfig(configPath)
	if err != nil {
		return err
	}

	var allFiles []string
	for _, target := range cfg.Targets {
		if target.IsString {
			files, err := resolver.FindFilesWithExclude(target.Path, nil)
			if err != nil {
				return err
			}
			allFiles = append(allFiles, files...)
		} else {
			files, err := resolver.FindFilesWithExclude(target.Path, target.Exclude)
			if err != nil {
				return err
			}
			allFiles = append(allFiles, files...)
		}
	}

	archiveName := cfg.Name + "-" + cfg.Ver + ".zip"
	if err := archiver.CreateZipArchive(allFiles, archiveName); err != nil {
		return err
	}

	sshCfg := config.LoadSSHConfig()
	remotePath := filepath.Join("/tmp", archiveName)
	if err := sshclient.UploadFile(archiveName, remotePath, sshCfg); err != nil {
		return err
	}

	fmt.Println("Archive uploaded to", remotePath)
	return nil
}

func HandleUpdate(configPath string) error {
	if !isLikelyJSON(configPath) {
		return fmt.Errorf("file is not a json")
	}
	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer f.Close()
	var probe interface{}
	dec := json.NewDecoder(f)
	if err := dec.Decode(&probe); err != nil {
		return fmt.Errorf("config file is not a valid json: %v", err)
	}

	cfg, err := config.ParsePackagesConfig(configPath)
	if err != nil {
		return err
	}

	sshCfg := config.LoadSSHConfig()
	for _, pkg := range cfg.Packages {
		client, err := sshclient.ConnectRaw(sshCfg)
		if err != nil {
			return err
		}
		files, err := client.ReadDir("/tmp")
		client.Close()
		if err != nil {
			return err
		}
		var candidates []string
		verRe := regexp.MustCompile(`^` + regexp.QuoteMeta(pkg.Name) + `-(.+)\.zip$`)
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			m := verRe.FindStringSubmatch(f.Name())
			if len(m) == 2 {
				if versionMatches(m[1], pkg.Ver) {
					candidates = append(candidates, m[1])
				}
			}
		}
		if len(candidates) == 0 {
			return fmt.Errorf("no archive found for %s with version condition %s", pkg.Name, pkg.Ver)
		}
		// выбрать максимальную подходящую версию
		sort.Slice(candidates, func(i, j int) bool { return versionLess(candidates[j], candidates[i]) })
		chosenVer := candidates[0]
		archiveName := pkg.Name + "-" + chosenVer + ".zip"
		remotePath := filepath.Join("/tmp", archiveName)
		localPath := archiveName
		if err := sshclient.DownloadFile(remotePath, localPath, sshCfg); err != nil {
			return err
		}
		fmt.Println("Downloaded:", localPath)
		if err := archiver.Unzip(localPath, "."); err != nil {
			return err
		}
	}
	return nil
}

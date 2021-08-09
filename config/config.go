package config

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-git/go-git/v5"
)

var (
	Global   config
	Manifest manifestConfig
	Version  string = "develop"
)

type manifestConfig struct {
	Name           string   `toml:"name"`
	Replications   int      `toml:"replications"`
	ClusterKey     string   `toml:"cluster_key"`
	StatsNode      string   `toml:"stats_node"`
	BootstrapPeers []string `toml:"bootstrap_peers"`
	Mirrors        []string `toml:"mirrors"`
}

type config struct {
	General  general
	Database database
	Ipfs     ipfs
	Mail     mail
	Manifest manifest
	Web      web
}

type general struct {
	Version string
}

type database struct {
	Path string
}

type mail struct {
	Domain     string
	PrivateKey string
	Sender     string
}

type manifest struct {
	Url  string
	Path string
}

type web struct {
	Addr string
}

type ipfs struct {
	Path       string
	PrivateKey string
	PeerID     string
	Addr       string
}

func Init() (err error) {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	Global = parseConfigEnv(
		&config{
			General: general{
				Version: Version,
			},
			Database: database{
				Path: "arkstat.db",
			},
			Manifest: manifest{
				Url:  "https://github.com/arken/core-manifest.git",
				Path: filepath.Join(user.HomeDir, ".arkstat", "manifest"),
			},
			Web: web{
				Addr: ":8080",
			},
			Ipfs: ipfs{
				Path:       filepath.Join(user.HomeDir, ".arkstat", "ipfs"),
				PeerID:     "",
				PrivateKey: "",
				Addr:       "",
			},
		},
	)
	Manifest, err = parseConfigManifest(Global.Manifest.Path, Global.Manifest.Url)
	return err
}

func parseConfigEnv(input *config) (result config) {
	numSubStructs := reflect.ValueOf(input).Elem().NumField()
	for i := 0; i < numSubStructs; i++ {
		iter := reflect.ValueOf(input).Elem().Field(i)
		subStruct := strings.ToUpper(iter.Type().Name())

		structType := iter.Type()
		for j := 0; j < iter.NumField(); j++ {
			fieldVal := iter.Field(j).String()
			if fieldVal != "Version" {
				fieldName := structType.Field(j).Name
				for _, prefix := range []string{"ARKSTAT"} {
					evName := prefix + "_" + subStruct + "_" + strings.ToUpper(fieldName)
					evVal, evExists := os.LookupEnv(evName)
					if evExists && evVal != fieldVal {
						iter.FieldByName(fieldName).SetString(evVal)
					}
				}
			}
		}
	}
	return *input
}

func parseConfigManifest(path, url string) (result manifestConfig, err error) {
	r, err := git.PlainOpen(path)
	if err != nil && err.Error() == "repository does not exist" {
		r, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
	}
	if err != nil {
		return result, err
	}
	w, err := r.Worktree()
	if err != nil {
		return result, err
	}
	err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	if err != nil && err.Error() != "already up-to-date" {
		return result, err
	}
	_, err = toml.DecodeFile(filepath.Join(Global.Manifest.Path, "config.toml"), &result)
	return result, err
}

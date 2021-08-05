package config

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
)

var (
	Global   config
	Manifest manifestConfig
	Version  string = "develop"
)

type manifestConfig struct {
	Name           string   `yaml:"name"`
	Replications   int      `yaml:"replications"`
	ClusterKey     string   `yaml:"cluster_key"`
	StatsNode      string   `yaml:"stats_node"`
	BootstrapPeers []string `yaml:"bootstrap_peers"`
	Mirrors        []string `yaml:"mirrors"`
}

type config struct {
	General  general
	Database database
	Ipfs     ipfs
	Keyset   keyset
	Web      web
}

type general struct {
	Version string
}

type database struct {
	Path string
}

type keyset struct {
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

func init() {
	var err error
	home, _ := os.UserHomeDir()
	Global = parseConfigEnv(
		&config{
			General: general{
				Version: Version,
			},
			Database: database{
				Path: "arkstat.db",
			},
			Keyset: keyset{
				Url:  "https://github.com/arken/core-keyset.git",
				Path: filepath.Join(home, ".config", "arkstat", "keyset"),
			},
			Web: web{
				Addr: ":8080",
			},
			Ipfs: ipfs{
				Path:       filepath.Join(home, ".config", "arkstat", "ipfs"),
				PeerID:     "",
				PrivateKey: "",
				Addr:       "",
			},
		},
	)
	Manifest, err = parseConfigManifest(Global.Keyset.Path, Global.Keyset.Url)
	if err != nil {
		log.Fatal(err)
	}
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
	bytes, err := os.ReadFile(filepath.Join(Global.Keyset.Path, "config.yml"))
	if os.IsNotExist(err) {
		return result, err
	}
	err = yaml.Unmarshal(bytes, &result)
	return result, err
}

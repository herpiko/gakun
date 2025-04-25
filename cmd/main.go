package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli/v3"
)

type Profiles map[string]map[string]string

type Config struct {
	Profiles  Profiles
	UpdatedAt int64
}

func main() {
	config, err := parseConfig(os.ExpandEnv("$HOME/.config/gakun/config.json"))
	if err != nil {
		panic(err)
	}

	var host string
	var path string

	cmd := &cli.Command{
		Name:  "gakun",
		Usage: "SSH key manager",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Usage:       "host",
				Destination: &host,
			},
			&cli.StringFlag{
				Name:        "path",
				Usage:       "path",
				Destination: &path,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "add",
				Aliases: []string{"s"},
				Usage:   "add host and key",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if host == "" {
						err = errors.New("host is required")
						return err
					}
					// TODO check for host string validation

					if path == "" {
						err = errors.New("path to key is required")
						return err
					}
					_, err := os.ReadFile(path)
					if err != nil {
						if os.IsNotExist(err) {
							return errors.New("SSH key path is not valid")
						} else {
							return err
						}
					}

					profile := cmd.Args().First()
					if config.Profiles == nil {
						config.Profiles = Profiles{}
					}
					if config.Profiles[profile] == nil {
						config.Profiles[profile] = map[string]string{}
					}
					config.Profiles[profile][host] = path

					// TODO saveConfig

					return nil
				},
			},
			{
				Name:    "use",
				Aliases: []string{"u"},
				Usage:   "use key",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("completed task: ", cmd.Args().First())

					// TODO read from config the apply to ~/.ssh/config

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}

func parseConfig(filePath string) (*Config, error) {
	var config Config

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return nil, err
			}
			defer f.Close()

			config.UpdatedAt = time.Now().Unix()
			data, err = json.Marshal(config)
			if err != nil {
				return nil, err
			}
			_, err = f.WriteString(string(data))
			if err != nil {
				return nil, err
			}
		}

	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func saveConfig(config *Config) error {
	return nil
}

func applyProfile(name string, host string, keyPath string) *error {
	return nil
}

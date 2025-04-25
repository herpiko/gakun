package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli/v3"
)

type Profiles map[string]map[string]string

type Config struct {
	Profiles  Profiles `json:"profiles"`
	UpdatedAt int64    `json:"updated_at"`
}

var configPath string
var sshConfigPath string

func main() {
	configPath = os.ExpandEnv("$HOME/.config/gakun/config.json")
	sshConfigPath = os.ExpandEnv("$HOME/.ssh/config")

	var err error

	gakun := Gakun{}
	err = gakun.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}

	var host string
	var key string

	cmd := &cli.Command{
		Name:  "gakun",
		Usage: "SSH key manager",
		Commands: []*cli.Command{
			{
				Name:  "add",
				Usage: "Add host and key to a profile. Example: 'gakun add work gitlab.com ~/.ssh/id_rsa_work'",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "host",
						Aliases:     []string{"h"},
						Usage:       "host",
						Destination: &host,
					},
					&cli.StringFlag{
						Name:        "key",
						Aliases:     []string{"k"},
						Usage:       "key",
						Destination: &key,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if host == "" {
						err = errors.New("host is required")
						return err
					}
					// TODO check for host string validation

					if key == "" {
						err = errors.New("path to key is required")
						return err
					}
					_, err := os.ReadFile(key)
					if err != nil {
						if os.IsNotExist(err) {
							return errors.New("SSH key path is not valid")
						} else {
							return err
						}
					}

					profile := cmd.Args().First()

					err = gakun.Add(profile, host, key)
					if err != nil {
						return nil
					}

					return nil
				},
			},
			{
				Name:  "use",
				Usage: "Use SSH key for certain host. Example: 'gakun use work -h gitlab.com'",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "host",
						Aliases:     []string{"h"},
						Usage:       "host",
						Destination: &host,
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var err error
					profile := cmd.Args().First()
					err = gakun.Use(profile, host)
					if err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "ls",
				Usage: "List profiles",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var err error
					err = gakun.List()
					if err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

type Gakun struct {
	Config Config
}

func (g *Gakun) Add(profile string, host string, key string) error {
	var err error
	if g.Config.Profiles == nil {
		g.Config.Profiles = Profiles{}
	}
	if g.Config.Profiles[profile] == nil {
		g.Config.Profiles[profile] = map[string]string{}
	}
	g.Config.Profiles[profile][host] = key

	err = g.saveConfig(configPath, &g.Config)
	if err != nil {
		return err
	}

	return nil
}

func (g *Gakun) Use(profile string, host string) error {
	var err error
	key := g.Config.Profiles[profile][host]

	if key == "" {
		err = errors.New("There is no such profile and host combination. Please type gakun ls to show your profiles and hosts.")
		return err
	}

	data, err := g.readFileWithSkipSection(sshConfigPath)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(sshConfigPath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	newConfig := "###### gakun begin"
	newConfig += "\nHost " + host
	newConfig += "\n  Hostname " + host
	newConfig += "\n  IdentityFile " + key
	newConfig += "\n###### gakun end"

	_, err = f.WriteString(data + newConfig)
	if err != nil {
		return err
	}

	fmt.Println("Key " + key + " is now active for " + host + " ✓")

	return nil
}

func (g *Gakun) List() error {
	var err error

	for profile := range g.Config.Profiles {
		fmt.Println("\n" + profile + ":")
		for host := range g.Config.Profiles[profile] {
			fmt.Println("   " + host + " → " + g.Config.Profiles[profile][host])
		}
	}

	return err
}

func (g *Gakun) LoadConfig(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
			if err != nil {
				return err
			}
			defer f.Close()

			g.Config.UpdatedAt = time.Now().Unix()
			data, err = json.Marshal(g.Config)
			if err != nil {
				return err
			}
			_, err = f.WriteString(string(data))
			if err != nil {
				return err
			}
		}

	}

	err = json.Unmarshal(data, &g.Config)
	if err != nil {
		return err
	}
	return nil
}

func (g *Gakun) saveConfig(filePath string, config *Config) error {
	config.UpdatedAt = time.Now().Unix()
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(string(data))
	if err != nil {
		return err
	}
	return nil
}

func (g *Gakun) applyProfile(name string, host string, keyPath string) *error {
	return nil
}

func (g *Gakun) readFileWithSkipSection(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var result strings.Builder

	// Flag to track if we're in the section to skip
	skipSection := false

	// Read line by line
	for scanner.Scan() {
		line := scanner.Text()

		// Check if we're entering the skip section
		if strings.Contains(line, "gakun begin") {
			skipSection = true
			continue
		}

		// Check if we're exiting the skip section
		if strings.Contains(line, "gakun end") {
			skipSection = false
			continue
		}

		// Only append if we're not in the skip section
		if !skipSection {
			result.WriteString(line + "\n")
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file: %w", err)
	}

	return result.String(), nil
}

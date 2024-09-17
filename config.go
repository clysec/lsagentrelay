package main

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Lansweeper struct {
		Url       string `yaml:"url"`
		Port      int    `yaml:"port"`
		IgnoreSSL bool   `yaml:"ignore_ssl"`
	}
	Listen struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
		Tls  struct {
			Enabled bool   `yaml:"enabled"`
			Cert    string `yaml:"cert"`
			Key     string `yaml:"key"`
		} `yaml:"tls"`
	}
	Rewrite struct {
		LansweeperRegex string `yaml:"LansweeperHostnameRegex"`
		ProxyHostname   string `yaml:"ProxyHostname"`
	} `yaml:"rewrite"`

	CompiledRegexp *regexp.Regexp
}

func (c *Config) Read(path string) error {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	c.CompiledRegexp = regexp.MustCompile(c.Rewrite.LansweeperRegex)

	return nil
}

func (c *Config) GetListener() string {
	return fmt.Sprintf("%s:%d", c.Listen.Host, c.Listen.Port)
}

func (c *Config) GetLansweeperUrl() string {
	return fmt.Sprintf("https://%s:%d/lsagent", c.Lansweeper.Url, c.Lansweeper.Port)
}

package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func Parse(
	param *string,
) (Config, error) {
	var cfg Config
	var setFile bool
	if param != nil {
		for _, parameter := range strings.Split(*param, ",") {
			par := strings.SplitN(parameter, "=", 2)
			if len(par) != 1 {
				continue
			}
			switch {
			case strings.HasSuffix(par[0], ".yaml") || strings.HasSuffix(par[0], ".yml"):
				if cfgData, err := os.ReadFile(par[0]); err == nil {
					if err = yaml.Unmarshal(cfgData, &cfg); err != nil {
						return cfg, err
					}
				}
				setFile = true
			case strings.HasSuffix(par[0], ".toml"):
				if cfgData, err := os.ReadFile(par[0]); err == nil {
					if err = toml.Unmarshal(cfgData, &cfg); err != nil {
						return cfg, err
					}
				}
				setFile = true
			}
		}
	}
	if !setFile {
		if cfgData, cfgErr := os.ReadFile(`litepb.yml`); cfgErr == nil {
			if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
				return cfg, err
			}
		}
		if cfgData, cfgErr := os.ReadFile(`litepb.yaml`); cfgErr == nil {
			if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
				return cfg, err
			}
		}
		if cfgData, cfgErr := os.ReadFile(`litepb.toml`); cfgErr == nil {
			if err := toml.Unmarshal(cfgData, &cfg); err != nil {
				return cfg, err
			}
		}
	}
	if param != nil {
		for _, parameter := range strings.Split(*param, ",") {
			par := strings.SplitN(parameter, "=", 2)
			if len(par) > 1 {
				if err := yaml.Unmarshal([]byte(fmt.Sprintf("%s: %s", par[0], par[1])), &cfg); err != nil {
					return cfg, err
				}
			}
		}
	}
	return cfg, nil
}

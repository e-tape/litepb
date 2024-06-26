package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func Parse(
	cfg *Config,
	param *string,
) error {
	if cfgData, cfgErr := os.ReadFile(`litepb.yaml`); cfgErr == nil {
		if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
			return err
		}
	}
	if param == nil {
		return nil
	}
	parameters := strings.Split(*param, ",")
	for _, parameter := range parameters {
		par := strings.SplitN(parameter, "=", 2)
		if len(par) == 1 {
			if cfgData, err := os.ReadFile(par[0]); err == nil {
				if err = yaml.Unmarshal(cfgData, &cfg); err != nil {
					return err
				}
			}
			continue
		}
		if err := yaml.Unmarshal([]byte(fmt.Sprintf("%s: %s", par[0], par[1])), &cfg); err != nil {
			return err
		}
	}
	return nil
}

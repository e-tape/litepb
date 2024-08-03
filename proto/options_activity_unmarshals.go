package litepb

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func (a *Activity) UnmarshalYAML(value *yaml.Node) error {
	if v, ok := activityMap[strings.ToLower(value.Value)]; ok {
		*a = v
		return nil
	}
	keys := make([]string, 0, len(activityMap))
	for k := range activityMap {
		keys = append(keys, k)
	}
	return fmt.Errorf("unknown Activity value [%s] with type [%s] available: %s", value.Value, value.Tag, strings.Join(keys, ", "))
}

func (a *Activity) UnmarshalTOML(value any) error {
	vs := fmt.Sprintf("%v", value)
	if v, ok := activityMap[strings.ToLower(vs)]; ok {
		*a = v
		return nil
	}
	keys := make([]string, 0, len(activityMap))
	for k := range activityMap {
		keys = append(keys, k)
	}
	return fmt.Errorf("unknown Activity value [%v] available: %s", value, strings.Join(keys, ", "))
}

var (
	activityMap = map[string]Activity{
		"1":        Activity_Active,
		"0":        Activity_Inactive,
		"active":   Activity_Active,
		"inactive": Activity_Inactive,
		"on":       Activity_Active,
		"off":      Activity_Inactive,
		"true":     Activity_Active,
		"false":    Activity_Inactive,
	}
)

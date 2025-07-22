package domain

import (
	"encoding/json"
	"errors"
)

type Target struct {
	Path    string   `json:"path"`
	Exclude []string `json:"exclude,omitempty"`
}

type TargetOrString struct {
	Target
	IsString bool
}

func (t *TargetOrString) UnmarshalJSON(data []byte) error {
	// Если это строка
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		t.Path = s
		t.IsString = true
		return nil
	}
	// Если это объект
	var obj Target
	if err := json.Unmarshal(data, &obj); err == nil {
		t.Target = obj
		t.IsString = false
		return nil
	}
	return errors.New("invalid target format")
}

type PacketConfig struct {
	Name    string           `json:"name"`
	Ver     string           `json:"ver"`
	Targets []TargetOrString `json:"targets"`
	Packets []struct {
		Name string `json:"name"`
		Ver  string `json:"ver"`
	} `json:"packets"`
}

type PackagesConfig struct {
	Packages []struct {
		Name string `json:"name"`
		Ver  string `json:"ver,omitempty"`
	} `json:"packages"`
}

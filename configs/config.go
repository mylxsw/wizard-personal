package configs

import (
	"encoding/json"

	"github.com/mylxsw/container"
)

type Config struct {
	Listen            string        `json:"listen"`
	UseLocalDashboard bool          `json:"use_local_dashboard"`
	WorkDir           string        `json:"work_dir"`
	Storage           StorageConfig `json:"storage"`
}

type StorageConfig struct {
	Type     string `json:"type"`
	BasePath string `json:"base_path"`
}

func (conf *Config) Serialize() string {
	rs, _ := json.Marshal(conf)
	return string(rs)
}

// Get return config object from container
func Get(cc container.Container) *Config {
	return cc.MustGet(&Config{}).(*Config)
}

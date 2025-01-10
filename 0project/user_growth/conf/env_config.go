package conf

import (
	"encoding/json"
	"log"
	"os"
)

var GlobalConfig *ProjectConfig

const envConfigName = "USER_GROWTH_CONFIG"

type ProjectConfig struct {
	DB struct {
		Engine          string `json:"engine"`
		UserName        string `json:"user_name"`
		Password        string `json:"password"`
		Host            string `json:"host"`
		Port            string `json:"port"`
		Database        string `json:"database"`
		Charset         string `json:"charset"`
		ShowSql         bool   `json:"show_sql"`
		MaxIdleConns    int    `json:"max_idle_conns"`
		MaxOpenConns    int    `json:"max_open_conns"`
		ConnMaxLifetime int    `json:"conn_max_lifetime"`
	}
}

func LoadConfig() {
	LoadEnvConfig()

}
func LoadEnvConfig() {
	pc := &ProjectConfig{}

	strConfigs := os.Getenv(envConfigName)

	if len(strConfigs) > 0 {
		err := json.Unmarshal([]byte(strConfigs), pc)
		if err != nil {
			log.Fatalf("conf.LoadEnvConfig(%s) error=%v", envConfigName, err)
			return
		}
	}
	GlobalConfig = pc
}

package amconfig

import (
	"crypto/tls"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

const (
	SessionStarterHTTPPort = "service.sessionstarter.http.port"

	AppBaseURL   = "app.baseurl"
	AppDBUser    = "app.db.user"
	AppDBPass    = "app.db.pass"
	AppDBName    = "app.db.name"
	AppDBHost    = "app.db.host"
	SMQHTTPPort  = "smq.http.port"
	SMQURL       = "smq.url"
	MockHTTPPort = "mock.http.port"

	SMQJobsTopicName    = "smq.jobs.topic.name"
	SMQResultsTopicName = "smq.results.topic.name"
)

type IConfig interface {
	GetString(string) string
	GetInt(string) int
	GetInt64(string) int64
	GetBool(string) bool
}

type ITLSConfig interface {
	GetCustomTLSConfig(string) (*tls.Config, error)
}

type ViperConfig struct {
	vconfig IConfig
}

func (cfg ViperConfig) GetString(s string) string {
	return cfg.vconfig.GetString(s)
}
func (cfg ViperConfig) GetInt(s string) int {
	return cfg.vconfig.GetInt(s)
}
func (cfg ViperConfig) GetInt64(s string) int64 {
	return cfg.vconfig.GetInt64(s)
}
func (cfg ViperConfig) GetBool(s string) bool {
	return cfg.vconfig.GetBool(s)
}

func NewViperConfig() (IConfig, error) {
	cfg, err := createViperConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func createViperConfig() (IConfig, error) {

	path, err := getConfigFileDir()
	if err != nil {
		log.Println("No configuration folder found")
	}

	viper.SetConfigName("app")
	viper.AddConfigPath(*path)
	viper.AddConfigPath(".")

	// ----- Env bindings -----
	_ = viper.BindEnv(SessionStarterHTTPPort, "SESSIONSTARTER_HTTP_PORT")

	_ = viper.BindEnv(AppBaseURL, "APP_BASE_URL")

	_ = viper.BindEnv(AppDBUser, "APP_DB_USER")
	_ = viper.BindEnv(AppDBPass, "APP_DB_PASS")
	_ = viper.BindEnv(AppDBName, "APP_DB_NAME")
	_ = viper.BindEnv(AppDBHost, "APP_DB_HOST")

	_ = viper.BindEnv(SMQHTTPPort, "SMQ_HTTP_PORT")
	_ = viper.BindEnv(SMQURL, "SMQ_URL")
	_ = viper.BindEnv(SMQJobsTopicName, "SMQ_JOBS_TOPIC_NAME")
	_ = viper.BindEnv(SMQResultsTopicName, "SMQ_RESULTS_TOPIC_NAME")

	_ = viper.BindEnv(MockHTTPPort, "MOCK_HTTP_PORT")

	viper.AutomaticEnv()
	_ = viper.ReadInConfig()

	configFileUsed := viper.ConfigFileUsed()
	if len(configFileUsed) == 0 {
		log.Println("no configuration file found")
	} else {
		log.Println("configuration file used")
	}
	return viper.GetViper(), nil
}

func getConfigFileDir() (*string, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dir, err := filepath.Abs(filepath.Dir(ex))
	if err != nil {
		return nil, err
	}
	return &dir, err
}

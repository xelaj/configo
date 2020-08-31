package prometheus

import (
	"net/url"
	"path/filepath"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"github.com/xelaj/go-dry"
	"github.com/xelaj/configo"
)

type AppConfig struct {
	GlobalConfig   GlobalConfig    `param:"global" validate:"required"`
	AlertingConfig AlertingConfig  `param:"alerting"`
	RuleFiles      []string        `param:"rule_files"`
	ScrapeConfigs  []*ScrapeConfig `param:"scrape_configs"`
}

type ScrapeConfig struct {
	JobName         string        `param:"job_name" validate:"required"`
	HonorLabels     bool          `param:"honor_labels"`
	HonorTimestamps bool          `param:"honor_timestamps"`
	Params          url.Values    `param:"params"`
	ScrapeInterval  time.Duration `param:"scrape_interval"`
	ScrapeTimeout   time.Duration `param:"scrape_timeout"`
	MetricsPath     string        `param:"metrics_path" validate:"file"`
	Scheme          string        `param:"scheme"`
	SampleLimit     uint          `param:"sample_limit"`
	TargetLimit     uint          `param:"target_limit"`
}

type GlobalConfig struct {
	ScrapeInterval     time.Duration     `param:"scrape_interval"`
	ScrapeTimeout      time.Duration     `param:"scrape_timeout"`
	EvaluationInterval time.Duration     `param:"evaluation_interval"`
	QueryLogFile       string            `param:"query_log_file" validate:"file"`
	ExternalLabels     map[string]string `param:"external_labels"`
}

type AlertingConfig struct {
	AlertRelabelConfigs map[string]string `param:"alert_relabel_configs"`
	AlertmanagerConfigs []string          `param:"alertmanagers"`
}

var (
	Config *AppConfig = &AppConfig{}
)

func TryIt() {
	_, filename, _, _ := runtime.Caller(0)
	system := filepath.Join(filepath.Dir(filename), "system")
	user := filepath.Join(filepath.Dir(filename), "user")

	godotenv.Load(filepath.Join(filepath.Dir(filename), "session.env"))
	err := configo.InitConfigWithExplicitConfigPaths("simplest", system, user, Config)
	dry.PanicIfErr(err)
	pp.Println(Config)
}

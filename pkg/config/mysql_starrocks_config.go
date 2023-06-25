package config

import (
	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
	"github.com/liuxinwang/go-mysql-starrocks/pkg/rule"
	"github.com/siddontang/go-log/log"
	"os"
	"path/filepath"
)

type MysqlSrConfig struct {
	Name       string
	Input      *InputConfig
	Mysql      *MysqlConfig
	Starrocks  *StarrocksConfig
	Filter     []*FilterConfig
	Rules      []*rule.MysqlToSrRule `toml:"rule"`
	Logger     *log.Logger
	ConfigFile string
	SyncParam  *SyncParamConfig `toml:"sync-param"`
}

func (config *MysqlSrConfig) ReadMysqlSrConf(filename string) (*MysqlSrConfig, error) {
	var err error
	if _, err = toml.DecodeFile(filename, &config); err != nil {
		return nil, errors.Trace(err)
	}
	return config, err
}

func NewMysqlSrConfig(configFile *string) *MysqlSrConfig {
	c := &MysqlSrConfig{}
	fileName, err := filepath.Abs(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	c, err = c.ReadMysqlSrConf(fileName)
	if err != nil {
		log.Fatal(err)
	}
	if c.Name == "" {
		log.Errorf("The configuration file \"name\" variable cannot be empty")
		os.Exit(0)
	}
	if c.SyncParam == nil {
		log.Errorf("The configuration file \"[sync-param]\" variable cannot be empty")
		os.Exit(0)
	}
	if c.SyncParam.ChannelSize < 100 {
		log.Warnf("The [sync-param] configuration parameter \"channel-size\" should not be less than 100, and reset configured channel-size = 100")
		c.SyncParam.ChannelSize = 100
	}
	if c.SyncParam.FlushDelaySecond < 1 {
		log.Warnf("The [sync-param] configuration parameter \"flush-delay-second\" should not be less than 1, and reset configured flush-delay-second = 1")
		c.SyncParam.FlushDelaySecond = 1
	}
	c.ConfigFile = fileName
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range c.Rules {
		r.RuleType = "init" // default 'init'
	}
	return c
}

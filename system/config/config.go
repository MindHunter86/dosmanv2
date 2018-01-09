package config

import "errors"
import "io/ioutil"

import "gopkg.in/yaml.v2"

// 2delete:
const (
	CTX_APP_CONFIG = uint8(iota)
	CTX_APP_LOGGER
)


type SysConfig struct {
	Base struct {
		Debug bool
		Log_level string
		Broker struct {
			Buffer int
		}
		Plugins struct {
			Basedir string
			Loadlist []string
		}
	}
	Vklogger struct {
		Login, Password string
		Cookie_storage struct {
			Path string
		}
		Http_api struct {
			Schema, Listen, Host string
		}
	}
	Mysql struct {
		Host, Username, Password, Database, Migrations_dir string
	}
}

func (self *SysConfig) Parse() (*SysConfig, error) {
	cnfile, e := ioutil.ReadFile("./conf/settings.yaml")
	if e != nil { return nil,errors.New("Could not read configuration file! | " + e.Error()) }

	if e = yaml.UnmarshalStrict(cnfile, self); e != nil {
		return nil,errors.New("Could not parse YAML from confifuration file! | " + e.Error())
	}
	return self,nil
}

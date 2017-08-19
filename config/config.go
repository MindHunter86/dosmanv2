package config

import "errors"
import "io/ioutil"

import "gopkg.in/yaml.v2"


type AppConfig struct {
	Base struct {
		Debug boolg
		Log_level string
		Http struct {
			Listen string
		}
	}
	Telegram struct {
		Token string
	}
	Mysql struct {
		Host, Username, Password, Database string
	}
}

func (self *AppConfig) Parse() (*AppConfig, error) {
	cnfile, e := ioutil.ReadFile("./conf/settings.yaml")
	if e != nil { return nil,errors.New("Could not read configuration file! | " + e.Error()) }

	if e = yaml.UnmarshalStrict(cnfile, self); e != nil {
		return nil,errors.New("Could not parse YAML from confifuration file! | " + e.Error())
	}
	return self,nil
}

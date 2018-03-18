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
			Buffer int }
		Plugins struct {
			Basedir string
			Loadlist []string }
	}
	Vklogger struct {
		Login, Password string
		Cookie_storage struct {
			Path string }
		Http_api struct {
			Schema, Listen, Host string }
	}
	Randoshka struct {
		Http struct {
			Schema, Listen, Host string
			Session_Keypairs struct {
				Current struct {
					Authentication, Encryption string }
				Previous struct {
					Authentication, Encryption string }
			}}
	}
	Mysql struct {
		Host, Username, Password, Database, Migrations_dir string }
	Telegram struct {
		Token string
		Timeout int }
	Sysparser struct {
		Http_Client struct {
			Ua, Origin string
			Timeout int }
		Storage struct {
			Bolt_Path string }
		Mysql struct {
			Hostname,Username,Password,Database string }
		Sysmru struct {
			Authentication struct {
				Creds struct {
					Username, Password, Next string}
				Login_Url, Login_Post string
				Test_String string }
			Session_Robber struct {
				Vulnerable_Url string }
			Calendar_Url string
			Parse_Until int }
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

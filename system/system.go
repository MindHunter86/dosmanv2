package system

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"plugin"

	"dosmanv2/modules"
	"dosmanv2/system/config"

	"github.com/rs/zerolog"
)


// System structs:
type System struct {
	wg *sync.WaitGroup
	log *zerolog.Logger
	cfg *config.SysConfig

	mods *modules.Modules
}


// System API:
func (m *System) SetLogger(logger *zerolog.Logger) *System {
	m.log = logger
	return m
}

func (m *System) Configure() (*System, error) {
	var e error

	// parse configuration file:
	if m.cfg, e = new(config.SysConfig).Parse(); e != nil { return nil,e }

	// define new modulelist:
	m.mods = new(modules.Modules)
	if m.mods.Config, e = new(config.SysConfig).Parse(); e != nil { return nil,e }

	m.mods.Hub = make(map[string]*modules.BaseModule)
	m.mods.Debug = m.mods.Config.Base.Debug
	m.mods.Logger = m.log
	m.mods.DonePipe = make(chan struct{})
	m.log.Debug().Msg("Config subsystem has been successfully configured!")

	// plugins loader:
	for _,pluginName := range m.cfg.Base.Plugins.Loadlist {
		m.log.Debug().Str("plugin", pluginName).Msg("Trying to stat plugin file...")

		pluginFile, e := os.Stat(m.cfg.Base.Plugins.Basedir + "/" + pluginName + ".so"); if os.IsNotExist(e) {
			m.log.Warn().Str("plugin", pluginName).Msg("Could not find plugin file!"); return nil,e
		} else if e != nil { m.log.Error().Str("plugin", pluginName).Err(e).Msg("Undefined error!") ; return nil,e }

		if e = m.preloadPlugin(pluginFile.Name()); e != nil { return nil,e }
	}

	return m,nil
}

func (m *System) Bootstrap() error {
	// define global error var for modError pipe:
	var e error
	var modErrorPipe chan *modules.ModuleError = make(chan *modules.ModuleError)

	// define kernel signal listener:
	var kernelSignal chan os.Signal = make(chan os.Signal)
	signal.Notify(kernelSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)

	// bootstrap configured modules:
	for modName, modPointer := range m.mods.Hub {
		if m.mods.Hub[modName].GetModuleStatus() == modules.StatusReady {
			go m.moduleBootstrap(modName, modPointer, modErrorPipe)
		}
	}

	m.log.Info().Msg("All modules were bootstrapped! Starting main event handler...")
	// start sys event loop:
LOOP:
	for {
		select {
		case <-kernelSignal:
			m.log.Warn().Msg("Syscall.SIG* has been detected! Closing application...")
			break LOOP
		case modError := <-modErrorPipe:
			e = modError.Error()
			m.log.Error().Str("MODULE", modError.ModuleName()).Err(e).Msg("CRITICAL ERROR!")
			break LOOP
		}
	}

	// TODO: Add buf for modErrorPipe. Check m.mods.WaitGroup
	// TODO: Check module order!!! 

	// stop and unload all modules:
	m.shutdown()

	// return nil or errors from modules (over mods.ErrorPipe):
	return e
}


// System internal methods:
func (m *System) preloadPlugin(plgName string) error {
	m.log.Debug().Str("plugin", plgName).Msg("preloadPlugin started...")

	plg, e := plugin.Open(m.cfg.Base.Plugins.Basedir + "/" + plgName); if e != nil {
		m.log.Warn().Str("plugin", plgName).Err(e).Msg("Could not load plugin!")
		return e
	}

	plgPointer, e := plg.Lookup("Plugin"); if e != nil {
		m.log.Warn().Str("plugin", plgName).Err(e).Msg("Could not find the Configure method!")
		return e
	}

	modPointer, e := plgPointer.(modules.Module).Construct(m.mods, nil); if e != nil {
		m.log.Warn().Str("plugin", plgName).Err(e).Msg("Could not execute the Configure method!")
		return e
	}

	m.mods.Hub[plgName] = &modules.BaseModule{ Module: modPointer }
	m.mods.Hub[plgName].SetModuleStatus(modules.StatusReady)

	m.log.Info().Str("plugin", plgName).Msg("Module has been successfully loaded and configured! Module status: StatusReady")
	return nil
}

func (m *System) moduleBootstrap(modName string, modPointer *modules.BaseModule, modError chan *modules.ModuleError) {
	m.log.Info().Msg("Bootstrapping module \""+modName+"\"...")

	m.mods.Hub[modName].SetModuleStatus(modules.StatusRunning)
	m.mods.WaitGroup.Add(1)

	if e := modPointer.Bootstrap(); e != nil && m.mods.Hub[modName].GetModuleStatus() != modules.StatusStopping {
		modError<- new(modules.ModuleError).SetModuleName(modName).SetError(e)
		m.mods.Hub[modName].SetModuleStatus(modules.StatusFailed)
		m.mods.WaitGroup.Done()
		return
	}

	m.mods.WaitGroup.Done()
	m.mods.Hub[modName].SetModuleStatus(modules.StatusReady)
	m.log.Info().Msg("Module \""+modName+"\" has been successfully stopped and unloaded!")
}

func (m *System) shutdown() {
	// Set status STOP for running modules:
	for modName, modPointer := range m.mods.Hub {
		if m.mods.Hub[modName].GetModuleStatus() == modules.StatusReady {
			m.mods.Hub[modName].SetModuleStatus(modules.StatusStopping)
			if e := modPointer.Destruct(); e != nil {
				m.log.Warn().Str("plugin", modName).Err(e).Msg("Could not correctly shutdown the plugin!")
			}
		}
	}

	// Close done pipe and wait modules unload:
	close(m.mods.DonePipe)
	m.mods.WaitGroup.Wait()
}

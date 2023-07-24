package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/layou233/ZBProxy/common"
	"github.com/layou233/ZBProxy/config"
	"github.com/layou233/ZBProxy/console"
	"github.com/layou233/ZBProxy/service"
	"github.com/layou233/ZBProxy/version"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

func main() {
	log.SetOutput(color.Output)
	//rand.Seed(time.Now().UnixNano())
	console.SetTitle(fmt.Sprintf("ZBProxy %v | Running...", version.Version))
	console.Println(color.HiRedString(` ______  _____   _____   _____    _____  __    __ __    __
|___  / |  _  \ |  _  \ |  _  \  /  _  \ \ \  / / \ \  / /
   / /  | |_| | | |_| | | |_| |  | | | |  \ \/ /   \ \/ /`), color.HiWhiteString(`
  / /   |  _  { |  ___/ |  _  /  | | | |   }  {     \  /
 / /__  | |_| | | |     | | \ \  | |_| |  / /\ \    / /
/_____| |_____/ |_|     |_|  \_\ \_____/ /_/  \_\  /_/`))
	color.HiGreen("Welcome to ZBProxy %s (%s)!\n", version.Version, version.CommitHash)
	color.HiBlack("Build Information: %s, %s/%s, CGO %s\n",
		runtime.Version(), runtime.GOOS, runtime.GOARCH, common.CGOHint)
	// go version.CheckUpdate()

	config.LoadConfig()
	service.Listeners = make([]net.Listener, 0, len(config.Config.Services))

	// hot reload
	// use inotify on Linux
	// use Win32 ReadDirectoryChangesW on Windows
	{
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Panic(err)
		}
		defer watcher.Close()
		err = monitorConfig(watcher)
		if err != nil {
			log.Panic("Config Reload Error : ", err)
		}
	}

	{
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
		<-osSignals
		// stop the program
		service.CleanupServices()
	}
}

func monitorConfig(watcher *fsnotify.Watcher) error {
	ctx, cancel := context.WithCancel(context.Background())
	service.ExecuteServices(ctx)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op.Has(fsnotify.Write) { // config reload
					// wait for the file to finish writing
					for {
						select {
						case <-watcher.Events:
						case <-time.After(time.Millisecond * 100):
							goto NextStep
						}
					}
				NextStep:
					log.Println(color.HiMagentaString("Config Reload : File change detected. Reloading..."))
					if config.LoadLists(true) { // reload success
						log.Println(color.HiMagentaString("Config Reload : Successfully reloaded Lists."))
						cancel()
						service.CleanupServices()
						service.Listeners = make([]net.Listener, 0, len(config.Config.Services))
						ctx, cancel = context.WithCancel(context.Background())
						service.ExecuteServices(ctx)
					} else {
						log.Println(color.HiMagentaString("Config Reload : Failed to reload Lists."))
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println(color.HiRedString("Config Reload Error : ", err))
			}
		}
	}()
	return watcher.Add("ZBProxy.json")
}

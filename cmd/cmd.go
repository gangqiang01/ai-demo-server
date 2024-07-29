package cmd

import (
	"os"

	_ "github.com/edgehook/ithings/common/dbm"
	"github.com/edgehook/ithings/webserver"
	"github.com/kardianos/service"

	"github.com/jwzl/beehive/pkg/core"
	"github.com/spf13/cobra"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
)

var logger service.Logger

var serviceConfig = &service.Config{
	Name:             "AI-DemoService",
	DisplayName:      "AI-Demo Service",
	Description:      "Advantech AI-Demo",
	UserName:         "",
	Arguments:        []string{},
	Executable:       "C:\\Users\\qa\\Desktop\\AI-Demo\\aiDemo.exe",
	Dependencies:     []string{},
	WorkingDirectory: "",
	ChRoot:           "",
	Option:           map[string]interface{}{},
}

func NewAppService() {
	options := make(service.KeyValue)
	options["DelayedAutoStart"] = true
	options["StartType"] = "automatic"
	options["OnFailure"] = "restart"
	options["OnFailureDelayDuration"] = "1s"
	options["OnFailureResetPeriod"] = 10

	serviceConfig.Option = options

	prog := &Program{}
	s, err := service.New(prog, serviceConfig)
	if err != nil {
		klog.Errorf("create windows service with error: %s", err.Error())
	}

	errs := make(chan error, 5)
	// logger, err = s.Logger(errs)
	logger, err = s.SystemLogger(errs)
	if err != nil {
		klog.Errorf("windows service logger with err: %v", err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				klog.Errorf("windows service with err: %v", err)
			}
		}
	}()

	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			klog.Infof("os.Args[2]: %v", os.Args[2])
			if os.Args[2] == "" {
				klog.Errorf("Executable Path is NULL, PLease check")
				return
			}
			serviceConfig.Executable = os.Args[2]
			s.Install()
			klog.Infof("Install Service Success")
			return
		}

		if os.Args[1] == "remove" {
			s.Uninstall()
			klog.Infof("Remove Service Success")
			return
		}
	}

	err = s.Run()
	if err != nil {
		klog.Errorf("windows service run with error: %s", err.Error())
	}

	// interrupt := make(chan os.Signal, 1)
	// signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// <-interrupt

	// klog.Errorf("Exit AppHub-Agent.exe")
	// os.Exit(0)
}

type Program struct{}

func (p *Program) Start(s service.Service) error {
	klog.Infof("==========  Start AppHub-Agent Service ==========")
	if service.Interactive() {
		logger.Infof("AppHub-Agent running %v. Running in terminal.", service.Platform())
		klog.Infof("Running in terminal.")
	} else {
		logger.Infof("AppHub-Agent running %v. Running under service manager.", service.Platform())
		klog.Infof("Running under service manager.")
	}
	go p.run()
	return nil
}

func (p *Program) Stop(s service.Service) error {
	klog.Errorf("==========  Stop AppHub-Agent Service ==========")
	logger.Info("AppHub-Agent Stopping!")
	logs.FlushLogs()
	return nil
}

func (p *Program) Shutdown(s service.Service) error {
	klog.Errorf("==========  Windows shutdown ==========")
	// klog.Errorf("Exit AppHub-Agent.exe")
	logger.Info("OS shutdown, AppHub-Agent Stopping!")
	logs.FlushLogs()
	return nil
}

func (p *Program) run() {
	defer func() {
		if err := recover(); err != nil {
			klog.Errorf("windows service lwm2m failed with %s", err)
		}
	}()

	klog.Infof("###########  Start the lwm2m client...! ###########")

	registerModules()
	// start all modules
	core.Run()
}

var Cmd = &cobra.Command{
	Use:     "ithings",
	Long:    `iot things manager.. `,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: To help debugging, immediately log version
		klog.Infof("###########  Start the ithings...! ###########")
		registerModules()
		// // start all modules
		core.Run()
	},
}

func Execute() error {
	if err := Cmd.Execute(); err != nil {
		return err
	}

	return nil
}

// register all module into beehive.
func registerModules() {
	webserver.Register()
}

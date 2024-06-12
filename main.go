package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bclswl0827/openstation-monitor/monitor"
	"github.com/bclswl0827/openstation-monitor/serial"
)

func parseCommandLine(args *arguments) {
	flag.StringVar(&args.config, "config", "./config.json", "Path to config file")
	flag.Parse()
}

func main() {
	var (
		args arguments
		conf config
	)
	parseCommandLine(&args)
	err := conf.Read(args.config)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("config file has been loaded")

	// Check if the interface exists
	ifName, err := getInterfaceByPattern(conf.IpNet.Pattern, conf.IpNet.Fuzzy)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("specified interface found: %s\n", ifName)

	// Open serial port
	monitorPort, err := serial.Open(conf.Monitor.Device, conf.Monitor.Baudrate)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("serial port %s opened\n", conf.Monitor.Device)
	defer serial.Close(monitorPort)

	// Create monitor dependency
	monitorDependency := &monitor.MonitorDependency{
		Port:  monitorPort,
		State: &monitor.MonitorState{Busy: true},
	}
	monitorDriver := monitor.MonitorDriver(&monitor.MonitorDriverImpl{})

	// Reset and initialize monitor device
	err = monitorDriver.Reset(monitorDependency)
	if err != nil {
		log.Fatalln(err)
	}
	err = monitorDriver.Init(monitorDependency)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("monitor device has been initialized")

	// Attach system signal handler
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Create timer to display IP addresses
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-sig:
			log.Println("received system signal, shutting down...")

			monitorDependency.State.Busy = false
			monitorDependency.State.Error = false
			monitorDriver.Clear(monitorDependency)
			monitorDriver.Display(
				monitorDependency,
				fmt.Sprintf("%s\nDevice Shutdown.", time.Now().UTC().Format("01-02 15:04:05")),
				0, 0,
			)

			return
		case <-ticker.C:
			ipMap, err := getIPv4Addrs()
			if err != nil {
				log.Println(err)
				monitorDependency.State.Error = true
				monitorDriver.Display(monitorDependency, "Failed to get\nIPv4 addresses.", 0, 0)
				continue
			}

			ip, ok := ipMap[ifName]
			if !ok {
				log.Printf("interface %s has no IPv4 address found", ifName)
				monitorDependency.State.Error = true
				monitorDriver.Display(monitorDependency, "Interface has no\nIPv4 address.", 0, 0)
				continue
			}

			log.Printf("got IP address for %s: %s\n", ifName, ip)
			monitorDependency.State.Error = false
			monitorDriver.Display(
				monitorDependency,
				fmt.Sprintf("%s\n%s", time.Now().UTC().Format("01-02 15:04:05"), ip),
				0,
				0,
			)
		}
	}
}

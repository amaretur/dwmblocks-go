package main

import (
	"os"
	"os/signal"
	"os/exec"
	"strings"
	"log"
	"time"
	"syscall"
	"flag"
	"encoding/json"
)

const (
	MIN_SIG uint8 = 34

	// line feed, LF — ASCII control character (0x0A, 
	// 10 in the decimal numeral system, '\n')
	LF byte	= byte(10)
)

// structure describing the info block
type block struct {
	Command		string	`json:"command"`
	Interval	uint16	`json:"interval"`
	Signal		uint8	`json:"signal"`
}

// block execution
func (b *block) Cmd() string {
   
	out, err := exec.Command("bash", "-c", b.Command).Output()
	if err != nil {
		log.Fatal(err)
	}

	// cuts the '/n' (LF) character from the end of the line if it is there 
	if out[len(out)-1] == LF {
		out = out[0:len(out)-1]
	}

	return string(out)
}

// list of all blocks and separator
type config struct {
	Blocks		[]block		`json:"blocks"`
	Separator	string		`json:"separator"` 
}

// stores the result of embedded blokes
var cache []string

// reading and initializing configurations
func getConfig(path string) *config {

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("cannot open config file with path '%s'", path)
	}
	defer f.Close()

	conf := new(config)

	err = json.NewDecoder(f).Decode(&conf)
	if err != nil {
		log.Fatal("incorrect configurations")
	}

	return conf
}

// bar status update
func updStatus(data []string, separator string) {
	exec.Command(
		"xsetroot", 
		"-name", 
		strings.Join(data, separator),
	).Run()
}

// main loop
func loop(conf *config) {

	var t uint16 = 0

	for {

		for i, item := range conf.Blocks {

			if item.Interval == 0 {
				continue
			}

			if t % item.Interval == 0 {
				cache[i] = item.Cmd()
			}
		}

		go updStatus(cache, conf.Separator)
		time.Sleep(1 * time.Second) 

		t++
	}
}

// cache initialization
func initCache(conf *config) {

	cache = make([]string, len(conf.Blocks))

	for i, item := range conf.Blocks {
		cache[i] = item.Cmd()
	}
}

func handleSignal(conf *config, blockId uint8) {

	cache[blockId] = conf.Blocks[blockId].Cmd()

	updStatus(cache, conf.Separator)
}

func setSignals(conf *config) {

	c := make(chan os.Signal)
   
	signalToBlock := make(map[uint8]uint8, len(conf.Blocks))

	// subscription to notifications about received signals
	for i, item := range conf.Blocks {

		// if the signal is zero - do not intercept it
		if item.Signal == 0 {
			continue
		}

		signal.Notify(c, syscall.Signal(item.Signal + MIN_SIG))
		signalToBlock[item.Signal] = uint8(i)
	}

	// waiting for signals to be received
	go func() {
		for {
			s := <-c  

			sig := uint8(s.(syscall.Signal))
			handleSignal(conf, signalToBlock[sig - MIN_SIG])
		}
	}()
}

func main() {

	homeDir, ok := os.LookupEnv("HOME")
	if ! ok {   
		log.Fatal("could not get the value of the environment variable $HOME")
	} 

	var path string

	flag.StringVar(
		&path, 
		"c", 
		homeDir + "/.config/dwmblocks/config.json", 
		"Path to the configuration file")

	flag.Parse()
   
	conf := getConfig(path)

	initCache(conf)
	setSignals(conf)
	loop(conf) 
}

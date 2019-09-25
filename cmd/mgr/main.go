package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nissy/envexpand/yaml"
	"github.com/nissy/mgr"
)

var (
	cfgFile   = flag.String("c", "mgr.toml", "")
	isHelp    = flag.Bool("h", false, "")
	isVersion = flag.Bool("v", false, "")
	version   = "dev"
)

type Mgr struct {
	ToRedis []*mgr.ToRedis `yaml:"to_redis"`
}

func main() {
	if err := run(); err != nil {
		if _, perr := fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error()); perr != nil {
			panic(err)
		}
	}
}

func run() error {
	flag.Parse()

	if *isHelp {
		_, err := fmt.Fprint(os.Stderr, help)
		return err
	}
	if *isVersion {
		fmt.Printf("Version is %s\n", version)
		return nil
	}

	m := &Mgr{}
	if err := yaml.Open(*cfgFile, m); err != nil {
		return err
	}
	for _, v := range m.ToRedis {
		if err := v.Do(); err != nil {
			return err
		}
	}

	fmt.Println("Migration finished.")
	return nil
}

var help = `Usage:
    mgr [options]
Options:
    -c string
        Set configuration file. (default "mgr.toml")
    -h bool
        This help.
    -v bool
        Display the version of mg.
`

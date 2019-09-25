package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nissy/envexpand/yaml"
	"github.com/nissy/mgr"
)

var cfgFile = flag.String("c", "mgr.toml", "")

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

	m := &Mgr{}
	if err := yaml.Open(*cfgFile, m); err != nil {
		return err
	}

	for _, v := range m.ToRedis {
		if err := v.Do(); err != nil {
			return err
		}
	}

	return nil
}

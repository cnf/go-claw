package config

import "flag"
import "path/filepath"

func (c *Config) parseCommandlineArgs() {
    flag.StringVar(&c.cfgfile, "conf", c.cfgfile, "path to our config file.")
    flag.IntVar(&c.httpport, "port", c.httpport, "port the web interface listens on")
    flag.BoolVar(&c.verbose, "v", c.verbose, "turn on verbose logging")
    flag.Parse()
    c.cfgfile, _ = filepath.Abs(c.cfgfile)
}

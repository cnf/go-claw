package config

// SetVerbose sets verbosity level
func (c *Config) SetVerbose(v bool) {
    c.verbose = v
}

// Verbose returns if we are set to verbose
func (c *Config) Verbose() bool {
    return c.verbose
}

// SetCfgFile sets config file path
func (c *Config) SetCfgFile(cfg string) {
    c.cfgfile = cfg
}

// CfgFile returns if we are set to verbose
// func (c *Config) CfgFile() bool {
//     return c.verbose
// }

// HTTPPort returns the HTTP Port number
func (c *Config) HTTPPort() int {
    return c.httpport
}

package config

// Setup prepares the config object with defaults etc
func (c *Config) Setup() {
    c.verbose = false
    c.cfgfile = getConfigPath()
    c.httpport = getDefaultHTTPPort()
    c.parseCommandlineArgs()
}

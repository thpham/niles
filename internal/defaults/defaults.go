package defaults

const (
	TickMin = 3 // minimal time in seconds between that can be set in config file. If not set or less then, Set to this value.
	HistTimeout = 30 // in seconds; must be >= 1
	HistStart = "now-7days" // must be >= 1

	AppName = "niles"

	EnvConfVarName = "NILES_CONF"
	ConfFileName   = "niles.conf"
	SiteConfDir    = "/etc/" + AppName + "/"
	SiteConfFile   = SiteConfDir + ConfFileName

	TemplatesDir        = SiteConfDir + "templates"
	TemplatesSuffix     = "-wft.yaml"
	TemplatesDescSuffix = "-wft.desc"
)

var (
	// default paths
	BinPaths = map[string]string{
		"kubectl":    "/usr/local/bin/kubectl",
	}
)

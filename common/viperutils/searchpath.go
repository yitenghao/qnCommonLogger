package viperutils

import (
	"log"
	"os"
)

var (
	searchPath []string
)

const (
	EnvGOPATH          = "GOPATH"
	EnvNameCfgPath     = "ENV_CFG_PATH"
	OfficialConfigPath = "/etc/officialconfig"
	RelativeConfigPath = "./config"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func init() {
	var altPath = os.Getenv(EnvNameCfgPath)
	var goPath = os.Getenv(EnvGOPATH)
	if altPath != "" {
		// If the user has overridden the path with an envvar, its the only path
		// we will consider

		if !fileExists(altPath) {
			log.Printf("%s %s does not exist\n", EnvNameCfgPath, altPath)
		}

		AddSearchPath(altPath)
	} else {
		// If we get here, we should use the default paths in priority order:
		//
		// *) CWD
		// *) /etc/nbcex

		// CWD
		if fileExists(goPath) {
			AddSearchPath(goPath)
		}

		AddSearchPath(RelativeConfigPath)

		// And finally, the official path
		if fileExists(OfficialConfigPath) {
			AddSearchPath(OfficialConfigPath)
		}
	}
}

func AddSearchPath(path string) {
	searchPath = append(searchPath, path)
}

func GetSearchPath() []string {
	return searchPath
}

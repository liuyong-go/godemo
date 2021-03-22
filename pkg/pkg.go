package pkg

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/liuyong-go/godemo/pkg/constant"
)

const jupiterVersion = "0.2.0"

var (
	startTime string
	goVersion string
)

// build info
/*

 */
var (
	appName         string
	appID           string
	hostName        string
	buildAppVersion string
	buildUser       string
	buildHost       string
	buildStatus     string
	buildTime       string
)

func init() {
	if appName == "" {
		appName = os.Getenv(constant.EnvAppName)
		if appName == "" {
			appName = filepath.Base(os.Args[0])
		}
	}

	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = time.Now().String()
	SetBuildTime(buildTime)
	goVersion = runtime.Version()
	InitEnv()
}

// Name gets application name.
func Name() string {
	return appName
}

//SetName set app anme
func SetName(s string) {
	appName = s
}

//AppID get appID
func AppID() string {
	return appID
}

//SetAppID set appID
func SetAppID(s string) {
	appID = s
}

//AppVersion get buildAppVersion
func AppVersion() string {
	return buildAppVersion
}

//appVersion not defined
// func SetAppVersion(s string) {
// 	appVersion = s
// }

//JupiterVersion get jupiterVersion
func JupiterVersion() string {
	return jupiterVersion
}

// todo: jupiterVersion is const not be set
// func SetJupiterVersion(s string) {
// 	jupiterVersion = s
// }

//BuildTime get buildTime
func BuildTime() string {
	return buildTime
}

//BuildUser get buildUser
func BuildUser() string {
	return buildUser
}

//BuildHost get buildHost
func BuildHost() string {
	return buildHost
}

//SetBuildTime set buildTime
func SetBuildTime(param string) {
	buildTime = strings.Replace(param, "--", " ", 1)
}

// HostName get host name
func HostName() string {
	return hostName
}

//StartTime get start time
func StartTime() string {
	return startTime
}

//GoVersion get go version
func GoVersion() string {
	return goVersion
}

// PrintVersion print formated version info
func PrintVersion() {

}

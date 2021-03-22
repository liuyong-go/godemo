package server

import (
	"context"

	"github.com/liuyong-go/godemo/pkg"
)

type Option func(c *ServiceInfo)

// ServiceInfo represents service info
type ServiceInfo struct {
	Name     string            `json:"name"`
	AppID    string            `json:"appId"`
	Scheme   string            `json:"scheme"`
	Address  string            `json:"address"`
	Weight   float64           `json:"weight"`
	Enable   bool              `json:"enable"`
	Healthy  bool              `json:"healthy"`
	Metadata map[string]string `json:"metadata"`
	Region   string            `json:"region"`
	Zone     string            `json:"zone"`
	// Deployment 部署组: 不同组的流量隔离
	// 比如某些服务给内部调用和第三方调用，可以配置不同的deployment,进行流量隔离
	Deployment string `json:"deployment"`
	// Group 流量组: 流量在Group之间进行负载均衡
	Group    string              `json:"group"`
	Services map[string]*Service `json:"services" toml:"services"`
}

// Service ...
type Service struct {
	Namespace string            `json:"namespace" toml:"namespace"`
	Name      string            `json:"name" toml:"name"`
	Labels    map[string]string `json:"labels" toml:"labels"`
	Methods   []string          `json:"methods" toml:"methods"`
}
type Server interface {
	Serve() error
	Stop() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
}

func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}
func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{
		Name:       pkg.Name(),
		AppID:      pkg.AppID(),
		Weight:     100,
		Enable:     true,
		Healthy:    true,
		Metadata:   make(map[string]string),
		Region:     pkg.AppRegion(),
		Zone:       pkg.AppZone(),
		Deployment: "",
		Group:      "",
	}
	si.Metadata["appMode"] = pkg.AppMode()
	si.Metadata["appHost"] = pkg.AppHost()
	si.Metadata["startTime"] = pkg.StartTime()
	si.Metadata["buildTime"] = pkg.BuildTime()
	si.Metadata["appVersion"] = pkg.AppVersion()
	si.Metadata["jupiterVersion"] = pkg.JupiterVersion()
	return si
}

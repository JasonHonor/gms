package device

import (
	"github.com/gogf/gf/container/garray"
)

type HostSecret struct {
	Host      string `json:"host"`
	Passwd    string `json:"pwd"`
	Community string `json:"community"`
}

type HostPort struct {
	Name     string `json:"name"`
	UpStream bool   `json:"upstream"`
}

type HostConfigItem struct {
	Host          string     `json:"host"`
	User          string     `json:"user"`
	Type          string     `json:"type"`
	Ports         []HostPort `json:"ports"`
	KexAlgorithms string     `json:"kex-algorithms"`
	//has walked ?
	Walked bool
}

type Device interface {
	//connect to device
	Connect()
	//probe system infomation
	Probe()
	//save current configuration
	Save()
	//disconnect from device
	Close()

	GetArpTable() *garray.Array
}

//find secret-info by host
func FindHostSec(hostsec []HostSecret, host HostConfigItem) *HostSecret {
	for _, hostSec := range hostsec {
		if hostSec.Host == host.Host {
			return &hostSec
		}
	}
	return nil
}

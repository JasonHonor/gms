package device

type HostSecret struct {
	Host      string `json:"host"`
	Passwd    string `json:"pwd"`
	Community string `json:"community"`
}

type HostPort struct {
	Name     string `json:"name"`
	UpStream bool   `json:"upstream"`
	Host     string `json:"host"`
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

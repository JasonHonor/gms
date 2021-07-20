package cluster

import (
	"encoding/json"
	"fmt"
	. "gms/device"
	"gms/device/cisco"
	"gms/device/huawei"
	"gms/utils"
	"strings"

	"github.com/gogf/gf/os/gfile"
)

type Cluster struct {
	MacMap      *utils.ArrayMap
	IpMap       *utils.ArrayMap
	PortMap     *utils.ArrayMap
	hosts       []HostConfigItem
	hostSecrets []HostSecret
}

//walk host
func (cluster *Cluster) WalkHost(host HostConfigItem, hostSecrets []HostSecret, useSnmp bool, filters []utils.FilterItem) (Device, error) {

	hostSec := FindHostSec(hostSecrets, host)
	if hostSec == nil {
		return nil, fmt.Errorf("Can't find secret info for host %s", host.Host)
	}

	if useSnmp { //Use SNMP
		sOid := "1.3.6.1.2.1.4.22.1.2"
		if hostSec.Community != "" {
			utils.WalkSnmpOid(host.Host, hostSec.Community,
				sOid, true, filters)
			fmt.Printf("Host %v Count=%d\n", host.Host, utils.SnmpCount)
		}
	} else { //use SSH
		//get device type
		fmt.Printf("================%v================\n", host.Host)
		device, err := cluster.CreateDevice(host, hostSec)
		if err == nil {
			device.Connect()
			device.Probe()
			//device.Save()
			device.Close()
			return device, nil
		} else {
			fmt.Printf("CreateDeviceError:%v\n", err)
		}
	}
	return nil, fmt.Errorf("Execution skipped.")
}

func (cluster *Cluster) CreateDevice(hostConf HostConfigItem, hostSec *HostSecret) (Device, error) {

	var dev Device
	var err error

	if strings.HasPrefix(hostConf.Type, "cisco") {
		dev = cisco.NewC2960(hostConf, hostSec)
	} else if strings.HasPrefix(hostConf.Type, "huawei") {
		dev = huawei.NewS5700(hostConf, hostSec)
	} else if strings.HasPrefix(hostConf.Type, "linksys") {
		dev = cisco.NewC2960(hostConf, hostSec)
	} else {
		err = fmt.Errorf("Unsupported device-type: %s\n", hostConf.Type)
	}

	return dev, err
}

func (cluster *Cluster) Discovery() {
	file := "cache" + utils.PathSep() + "-mac-map.json"
	if !gfile.Exists(file) {
		for _, host := range cluster.hosts {
			//if idx > 3 {
			//	continue
			//}
			if !host.Walked {
				if !strings.HasPrefix(host.Host, "#") {
					dev, err := cluster.WalkHost(host, cluster.hostSecrets, false, nil)
					if err == nil {
						gArr := dev.GetArpTable()
						if gArr != nil {
							gArr.Iterator(func(k int, v interface{}) bool {
								arpItem, ok := v.(utils.ArpItem)
								if ok {
									cluster.MacMap.Set(arpItem.Mac, arpItem)
									cluster.IpMap.Set(arpItem.IP, arpItem)
								}
								return true
							})
						}
					} else {
						fmt.Printf("WalkHostError %v\n", err)
					}
				}
			}
		}
		cluster.MacMap.Save(file)

		fileIp := "cache" + utils.PathSep() + "-ip-map.json"
		cluster.IpMap.Save(fileIp)
	} else {
		file := "cache" + utils.PathSep() + "-mac-map.json"
		cluster.MacMap.Load(file)

		fileIp := "cache" + utils.PathSep() + "-ip-map.json"
		cluster.IpMap.Load(fileIp)

	}
}

func (cluster *Cluster) Init() {

	intBytes := gfile.GetBytes("config.json")
	err := json.Unmarshal(intBytes, &cluster.hosts)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
	}

	secBytes := gfile.GetBytes("secret.json")
	errSec := json.Unmarshal(secBytes, &cluster.hostSecrets)
	if errSec != nil {
		fmt.Printf("Error:%v\n", errSec)
		return
	}

}

func (cluster *Cluster) FindByMac(mac string) string {
	dat, _ := cluster.MacMap.Get(mac)
	if dat == nil {
		return ""
	}
	var ret string
	dat.Iterator(func(k int, v interface{}) bool {
		arpItem := v.(utils.ArpItem)
		ret += arpItem.DumpByMac()
		return true
	})
	return ret
}

func (cluster *Cluster) FindByIP(ip string) string {

	dat, _ := cluster.IpMap.Get(ip)

	if dat == nil {
		fmt.Println("NotFound!")
		return "NotFound!"
	}

	var ret string
	dat.Iterator(func(k int, v interface{}) bool {
		arpItem := v.(utils.ArpItem)
		//fmt.Println("%v", arpItem)
		ret += arpItem.DumpByIP()
		return true
	})
	return ret
}

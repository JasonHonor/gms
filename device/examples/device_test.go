package examples

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/gogf/gf/os/gfile"

	. "gms/device"
	"gms/utils"

	"gms/device/cisco"
	"gms/device/huawei"

	"github.com/gogf/gf/frame/g"
)

func TestSnmpWalkConfigFile(t *testing.T) {
	g.Cfg().SetFileName("config.json")

	hosts := []HostConfigItem{}
	intBytes := gfile.GetBytes("config.json")
	err := json.Unmarshal(intBytes, &hosts)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
	}
	//fmt.Printf("Hosts:%v\n", hosts)

	hostSecrets := []HostSecret{}
	secBytes := gfile.GetBytes("secret.json")
	errSec := json.Unmarshal(secBytes, &hostSecrets)
	if errSec != nil {
		fmt.Printf("Error:%v\n", errSec)
		return
	}

	for _, host := range hosts {
		if !host.Walked {
			if !strings.HasPrefix(host.Host, "#") {
				WalkHost(host, hostSecrets, false, nil)
			}
		}
	}
}

//create device
func CreateDevice(hostConf HostConfigItem, hostSec *HostSecret) (Device, error) {

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

//walk host
func WalkHost(host HostConfigItem, hostSecrets []HostSecret, useSnmp bool, filters []utils.FilterItem) {

	hostSec := FindHostSec(hostSecrets, host)
	if hostSec == nil {
		return
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
		device, err := CreateDevice(host, hostSec)
		if err == nil {
			device.Connect()
			device.Probe()
			//device.Save()
			device.Close()
		} else {
			fmt.Printf("CreateDeviceError:%v\n", err)
		}
	}
}

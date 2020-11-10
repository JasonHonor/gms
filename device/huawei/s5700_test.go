package huawei

import (
	"encoding/json"
	"fmt"
	"gms/service"
	"gms/utils"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"

	. "gms/utils"
)

func findHostSec(hostsec []HostSecret, host string) string {
	for _, hostSec := range hostsec {
		if hostSec.Host == host {
			return hostSec.Passwd
		}
	}
	return ""
}

func TestGet5700Info(t *testing.T) {
	return

	g.Cfg().SetFileName("config.json")

	//return
	hostSecrets := []HostSecret{}
	secBytes := gfile.GetBytes("secret.json")
	errSec := json.Unmarshal(secBytes, &hostSecrets)
	if errSec != nil {
		fmt.Printf("Error:%v\n", errSec)
		return
	}

	fmt.Printf("HostSecrets:%v\n", hostSecrets)

	hosts := []HostConfigItem{}
	intBytes := gfile.GetBytes("config.json")

	err := json.Unmarshal(intBytes, &hosts)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
	}

	host := hosts[0].Host
	pwd := findHostSec(hostSecrets, host)

	if pwd == "" {
		fmt.Printf("NoPassword return \n")
		return
	}

	dev := NewS5700(hosts[0].Host, hosts[0].User, pwd)

	dev.Probe()
	dev.Save()
	//dev.Dump()
	dev.DumpArpTables()

	//FindInterfaceByIP
	ifName := dev.FindInterfaceByIP("172.20.81.11")
	fmt.Printf("IfName %s\n", ifName)

	//arpList := dev.FindArpListByMac("c81f")
	//fmt.Printf("ArpList Cnt=%d %v\n", len(arpList), arpList)
	//register to zabbix
	dev.ArpTable.Iterator(func(k int, v interface{}) bool {
		arpItem, ok1 := v.(ArpItem)
		if ok1 {

			group := service.GetHostGroupIdByName(t, arpItem.Vlan)

			if len(group) > 0 {
				service.CreateHostByGroupId(group, arpItem.IP, arpItem.Mac, t)
			}
		}
		return true
	})
}

func filterByMacPrefixFile(mac, file string) bool {

	var sMac string = ""
	sMac = strings.ReplaceAll(mac, ":", "")
	sMac = strings.ToUpper(sMac)

	//filter by file lines.
	var sFileContent string = ""
	f, err := ioutil.ReadFile(file)
	if err == nil {
		sFileContent = string(f)
	} else {
		return false
	}

	sFileLines := strings.Split(sFileContent, "\n")
	var bFound bool = false
	for _, sFileLine := range sFileLines {
		//fmt.Printf("Filter 1. %s 2. %s\n", sFileLine, sMac)
		if strings.HasPrefix(sMac, sFileLine) {
			bFound = true
			break
		}
	}

	return bFound
}

func TestGet5700InfoByCache(t *testing.T) {
	return

	g.Cfg().SetFileName("config.json")

	dev := NewS5700(g.Cfg().Get("host").(string), g.Cfg().Get("user").(string), g.Cfg().Get("pwd").(string))

	dev.Load()
	//dev.Dump()
	dev.DumpArpTables()

	//register to zabbix
	dev.ArpTable.Iterator(func(k int, v interface{}) bool {
		arpItem, ok1 := v.(ArpItem)
		if ok1 {

			group := service.GetHostGroupIdByName(t, arpItem.Vlan)

			if len(group) > 0 {
				service.CreateHostByGroupId(group, arpItem.IP, arpItem.Mac, t)
			}
		}

		return true
	})

	//FindInterfaceByIP
	ifName := dev.FindInterfaceByIP("172.20.81.11")
	fmt.Printf("IfName %s\n", ifName)

	arpList := dev.FindArpListByMac("c81f")
	fmt.Printf("ArpList Cnt=%d %v\n", len(arpList), arpList)

	dev.PrintDhcpLeaseFile("c81f", "65")
	//fmt.Printf("ArpList Cnt=%d %v\n", len(arpList2), arpList2)
}

func TestGetSnmpInfo(t *testing.T) {
	return

	hOids := NewHuaweiSnmpOids()
	utils.GetSNMPInfo("172.20.65.254", "",
		[]string{hOids.SysName, hOids.SysDescription, hOids.SysUptime, hOids.SysObjectID})

	//cOids := NewCiscoSnmpOids()
	utils.GetSNMPInfo("172.20.65.249", "",
		[]string{hOids.SysName, hOids.SysDescription, hOids.SysUptime, hOids.SysObjectID})
}

func WalkHost(host HostConfigItem, filters []utils.FilterItem) {
	//return
	fmt.Printf("Host %v\n", host.Host)
	/*
		sOid := "1.3.6.1.2.1.4.22.1.2"
		if host.Community != "" {
			utils.WalkSnmpOid(host.Host, host.Community,
				sOid, true, filters)
			fmt.Printf("Host %v Count=%d\n", host.Host, utils.SnmpCount)
		}
	*/
}

func GetFilterData() []FilterItem {
	return nil

	var filters []FilterItem

	filters = append(filters, utils.FilterItem{
		Name:  "GE26",
		Type:  "mac",
		Value: " 2c:44:fd:7d:5d:e0 ",
	})

	filters = append(filters, utils.FilterItem{
		Name:  "GE35",
		Type:  "mac",
		Value: "  e4:11:5b:0c:89:6a  ",
	})

	filters = append(filters, utils.FilterItem{
		Name:  "GE39",
		Type:  "mac",
		Value: "  ac:16:2d:75:2a:76  ",
	})

	filters = append(filters, utils.FilterItem{
		Name:  "GE40",
		Type:  "mac",
		Value: "  24:5e:be:0e:18:f0  ",
	})

	filters = append(filters, utils.FilterItem{
		Name:  "GE41",
		Type:  "mac",
		Value: "   24:5e:be:0e:18:f1   ",
	})

	filters = append(filters, utils.FilterItem{
		Name:  "GE42",
		Type:  "mac",
		Value: "   24:5e:be:0e:18:f3   ",
	})

	filters = append(filters, utils.FilterItem{
		Name:  "GE43",
		Type:  "mac",
		Value: "   24:5e:be:0e:18:f2   ",
	})

	return filters
}

func TestSnmpWalkConfigFile(t *testing.T) {

	hosts := []HostConfigItem{}

	intBytes := gfile.GetBytes("config.json")

	err := json.Unmarshal(intBytes, &hosts)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
	}

	//fmt.Printf("Hosts:%v\n", hosts)

	for _, host := range hosts {
		if !host.Walked {
			if !strings.HasPrefix(host.Host, "#") {
				WalkHost(host, GetFilterData())
			}
		}
	}
}

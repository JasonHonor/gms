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
)

type HostPort struct {
	Name     string `json:"name"`
	UpStream bool   `json:"upstream"`
	Host     string `json:"host"`
}

type HostConfigItem struct {
	Host      string     `json:"host"`
	User      string     `json:"user"`
	Ports     []HostPort `json:"ports"`
	Community string     `json:"community"`
	//has walked ?
	Walked bool
}

func TestGet5700Info(t *testing.T) {
	return

	g.Cfg().SetFileName("config.json")

	dev := NewS5700(g.Cfg().Get("host").(string), g.Cfg().Get("user").(string), g.Cfg().Get("pwd").(string))

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

func TestWalkSnmpOid(t *testing.T) {
	return

	hOids := NewHuaweiSnmpOids()
	/*utils.WalkSnmpOid("172.20.65.254", "",
		hOids.PhysicalNameTable)

	//cOids := NewCiscoSnmpOids()
	utils.WalkSnmpOid("172.20.65.249", "",
		hOids.PhysicalNameTable)

	utils.WalkSnmpOid("172.20.65.254", "",
		hOids.IfDescrTable)

	utils.WalkSnmpOid("172.20.65.249", "",
		hOids.IfDescrTable)*/
	fmt.Println("---------EAST-CORE------------------")
	utils.WalkSnmpOid("172.20.65.254", "",
		hOids.ArpTable, true)

	fmt.Println("---------WEST-CORE------------------")
	utils.WalkSnmpOid("172.20.65.250", "",
		hOids.ArpTable, true)

	fmt.Println("---------WEST-C5------------------")
	utils.WalkSnmpOid("172.20.65.249", "",
		hOids.ArpTable, true)

	fmt.Println("---------WEST-C4------------------")
	utils.WalkSnmpOid("172.20.65.251", "",
		hOids.ArpTable, true)

	fmt.Println("---------WEST-C3------------------")
	utils.WalkSnmpOid("172.20.65.252", "",
		hOids.ArpTable, true)

	fmt.Println("---------EAST-B3------------------")
	utils.WalkSnmpOid("172.20.65.247", "",
		hOids.ArpTable, true)

	fmt.Println("---------WEST-B2------------------")
	utils.WalkSnmpOid("172.20.65.248", "",
		hOids.ArpTable, true)
}

func WalkHost(host HostConfigItem) {
	fmt.Printf("Host %v\n", host.Host)

	sOid := "1.3.6.1.2.1.4.22.1.2"

	if host.Community != "" {
		utils.WalkSnmpOid(host.Host, host.Community,
			sOid, true)

		fmt.Printf("Host %v Count=%d\n", host.Host, utils.SnmpCount)
	}
}

func TestWalkConfigFile(t *testing.T) {

	hosts := []HostConfigItem{}

	intBytes := gfile.GetBytes("config.json")

	err := json.Unmarshal(intBytes, &hosts)
	if err != nil {
		fmt.Printf("Error:%v\n", err)
	}

	fmt.Printf("Hosts:%v\n", hosts)

	for _, host := range hosts {
		if !host.Walked {
			WalkHost(host)
		}
	}
}

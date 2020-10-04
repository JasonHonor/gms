package huawei

import (
	"fmt"
	"gms/utils"
	"testing"
)

/*
func TestGet5700Info(t *testing.T) {

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
}

func TestGet5700InfoByCache(t *testing.T) {

	g.Cfg().SetFileName("config.json")

	dev := NewS5700(g.Cfg().Get("host").(string), g.Cfg().Get("user").(string), g.Cfg().Get("pwd").(string))

	dev.Load()
	//dev.Dump()
	dev.DumpArpTables()

	//FindInterfaceByIP
	ifName := dev.FindInterfaceByIP("172.20.81.11")
	fmt.Printf("IfName %s\n", ifName)

	//arpList := dev.FindArpListByMac("c81f")
	//fmt.Printf("ArpList Cnt=%d %v\n", len(arpList), arpList)
}

func TestGetSnmpInfo(t *testing.T) {
	hOids := NewHuaweiSnmpOids()
	utils.GetSNMPInfo("172.20.65.254", "",
		[]string{hOids.SysName, hOids.SysDescription, hOids.SysUptime, hOids.SysObjectID})

	//cOids := NewCiscoSnmpOids()
	utils.GetSNMPInfo("172.20.65.249", "",
		[]string{hOids.SysName, hOids.SysDescription, hOids.SysUptime, hOids.SysObjectID})
}
*/
func TestWalkSnmpOid(t *testing.T) {
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

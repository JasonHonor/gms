package huawei

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/frame/g"
)

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

	arpList := dev.FindArpListByMac("c81f")
	fmt.Printf("ArpList Cnt=%d %v\n", len(arpList), arpList)
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

	arpList := dev.FindArpListByMac("c81f")
	fmt.Printf("ArpList Cnt=%d %v\n", len(arpList), arpList)
}

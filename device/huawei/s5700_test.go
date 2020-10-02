package huawei

import (
	"fmt"
	. "gms/utils"
	"testing"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/frame/g"
)

func TestGet5700Info(t *testing.T) {

	g.Cfg().SetFileName("config.json")

	dev := S5700{
		SSHClient: SSHClient{
			Host:            g.Cfg().Get("host").(string),
			Port:            22,
			Username:        g.Cfg().Get("user").(string),
			Password:        g.Cfg().Get("pwd").(string),
			MoreTag:         "---- More ----",
			IsMoreLine:      true,
			MoreWant:        " ",
			ColorTag:        "1b5b343244H", //\u001b[42D
			ReadOnlyPrompt:  ">",
			SysEnablePrompt: "]",
			LineBreak:       "\r\n",
			ExitCmd:         "quit",
		},
		InterfaceList:   gset.NewSet(),
		InterfaceIpList: gset.NewSet(),
		ArpTable:        garray.NewArray(),
		IfArpCounts:     gmap.NewStrIntMap(),
		UpStreamIf:      "GigabitEthernet0/0/24",
	}

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

	dev := S5700{
		SSHClient: SSHClient{
			Host:            g.Cfg().Get("host").(string),
			Port:            22,
			Username:        g.Cfg().Get("user").(string),
			Password:        g.Cfg().Get("pwd").(string),
			MoreTag:         "More",
			MoreWant:        " ",
			ColorTag:        "\u001b[42D",
			ReadOnlyPrompt:  ">",
			SysEnablePrompt: "]",
			LineBreak:       "\r\n",
			ExitCmd:         "quit",
		},
		InterfaceList:   gset.NewSet(),
		InterfaceIpList: gset.NewSet(),
		ArpTable:        garray.NewArray(),
		IfArpCounts:     gmap.NewStrIntMap(),
		UpStreamIf:      "GigabitEthernet0/0/24",
	}

	dev.Load()
	//dev.Dump()
	dev.DumpArpTables()

	//FindInterfaceByIP
	ifName := dev.FindInterfaceByIP("172.20.81.11")
	fmt.Printf("IfName %s\n", ifName)

	arpList := dev.FindArpListByMac("c81f")
	fmt.Printf("ArpList Cnt=%d %v\n", len(arpList), arpList)
}

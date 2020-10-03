package huawei

import (
	"fmt"
	. "gms/utils"
	"strconv"
	"strings"

	"encoding/json"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gset"
)

// IPInterface switch interface with ip info.
type IPInterface struct {
	//Name Port name
	Name string

	//IpAndMask IP with Mask
	IPAndMask string

	//PhysicalState
	PhysicalState string

	//ProtoState
	ProtocalState string
}

type Interface struct {
	//Name Port name
	Name string

	//PhysicalState
	PhysicalState string

	//ProtoState
	ProtocalState string

	//InBound utility
	InUti string

	//OutBound utility
	OutUti string

	//InBound errors
	InErrors string

	//OutBound errors
	OutErrors string
}

type ArpItem struct {
	IP        string
	Mac       string
	Vlan      string
	Interface string
}

// S5700 Switch demo
type S5700 struct {
	SSHClient

	//Switch interface list
	InterfaceList *gset.Set

	//Switch ip list
	InterfaceIpList *gset.Set

	UpStreamIf string

	ArpTable *garray.Array

	PhoneTable *garray.Array

	IfArpCounts *gmap.StrIntMap
}

func (dev *S5700) Probe() {

	if dev.InterfaceIpList == nil {
		dev.InterfaceIpList = gset.NewSet()
	}

	if dev.InterfaceList == nil {
		dev.InterfaceList = gset.NewSet()
	}

	results := dev.Execute([]string{
		"disp int bri",
		"disp ip int bri",
		"disp arp",
	})

	for _, ret := range results {
		//fetch first line
		lines := strings.Split(ret, dev.LineBreak)

		if len(lines) > 0 && lines[0] == "disp ip int bri" {
			dev.ParseIpInterface(lines)
		}

		if len(lines) > 0 && lines[0] == "disp int bri" {
			dev.ParseInterface(lines)
		}

		if len(lines) > 0 && lines[0] == "disp arp" {
			dev.ParseArp(lines)
			dev.CalcArpTables()
		}
	}
}

func (dev *S5700) SumUtil(exc string) (string, string, int16, int16) {
	var dataIn, dataOut float64 = 0.0, 0.0
	var downCnt, upCnt int16 = 0, 0
	dev.InterfaceList.Iterator(func(v interface{}) bool {
		intf := v.(Interface)
		if intf.Name != exc {
			util1, err1 := strconv.ParseFloat(strings.Replace(intf.InUti, "%", "", -1), 64)
			if err1 == nil {
				dataIn += util1
			}

			util2, err2 := strconv.ParseFloat(strings.Replace(intf.OutUti, "%", "", -1), 64)
			if err2 == nil {
				dataOut += util2
			}
		}

		if strings.Contains(intf.PhysicalState, "down") {
			downCnt++
		}
		if strings.Contains(intf.PhysicalState, "up") {
			upCnt++
		}

		return true
	})
	return fmt.Sprintf("%v", dataIn), fmt.Sprintf("%v", dataOut), downCnt, upCnt
}

func (dev *S5700) CalcArpTables() {

	dev.ArpTable.Iterator(func(k int, v interface{}) bool {
		arpItem, ok1 := v.(ArpItem)
		if ok1 {
			//fmt.Printf("%s\n", arpItem.Interface)
			nCount, ok2 := dev.IfArpCounts.Search(arpItem.Interface)
			if ok2 {
				dev.IfArpCounts.Set(arpItem.Interface, nCount+1)
			} else {
				dev.IfArpCounts.Set(arpItem.Interface, 1)
			}
		}
		return true
	})
}

func (dev *S5700) Dump() {

	s, err := json.Marshal(dev.InterfaceList)
	if err == nil {
		fmt.Printf("%s\n", s)
	} else {
		fmt.Printf("Err=%v\n", err)
	}

	sIn, sOut, downCnt, upCnt := dev.SumUtil(dev.UpStreamIf)
	fmt.Printf("SumIn=%s SumOut=%s downCnt=%d upCnt=%d\n", sIn, sOut, downCnt, upCnt)

	s1, err1 := json.Marshal(dev.InterfaceIpList)
	if err1 == nil {
		fmt.Printf("%s\n", s1)
	} else {
		fmt.Printf("Err=%v\n", err1)
	}

	s2, err2 := json.Marshal(dev.ArpTable)
	if err2 == nil {
		fmt.Printf("%s\n", s2)
	} else {
		fmt.Printf("Err=%v\n", err2)
	}

	s3, err3 := json.Marshal(dev.IfArpCounts)
	if err3 == nil {
		fmt.Printf("%s\n", s3)
	} else {
		fmt.Printf("Err=%v\n", err3)
	}
}

func (dev *S5700) DumpArpTables() {
	s3, err3 := json.Marshal(dev.IfArpCounts)
	if err3 == nil {
		fmt.Printf("Arp Count=%d %s\n", dev.ArpTable.Len(), s3)
	} else {
		fmt.Printf("Err=%v\n", err3)
	}

	s4, err4 := json.Marshal(dev.PhoneTable)
	if err4 == nil {
		fmt.Printf("PhoneTable count=%d %s\n", dev.PhoneTable.Len(), s4)
	} else {
		fmt.Printf("Err=%v\n", err4)
	}

	fmt.Printf("ArpSum Total:%d Phone:%d PC:%d\n", dev.ArpTable.Len(), dev.PhoneTable.Len(), dev.ArpTable.Len()-dev.PhoneTable.Len())
}

func NewS5700(host, user, pwd string) *S5700 {
	return &S5700{
		SSHClient: SSHClient{
			Host:            host,
			Port:            22,
			Username:        user,
			Password:        pwd,
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
		PhoneTable:      garray.NewArray(),
		IfArpCounts:     gmap.NewStrIntMap(),
		UpStreamIf:      "GigabitEthernet0/0/24",
	}
}

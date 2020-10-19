package cisco

import (
	"encoding/json"
	"fmt"
	"gms/utils"
	"strings"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gset"
)

// C2960 Switch demo
type C2960 struct {
	utils.SSHClient

	//Switch interface list
	InterfaceList *gset.Set

	//Switch ip list
	InterfaceIpList *gset.Set

	UpStreamIf string

	ArpTable *garray.Array

	MacTable *garray.Array

	PhoneTable *garray.Array
}

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

func (dev *C2960) Probe() {

	if dev.InterfaceIpList == nil {
		dev.InterfaceIpList = gset.NewSet()
	}

	if dev.InterfaceList == nil {
		dev.InterfaceList = gset.NewSet()
	}

	results := dev.Execute([]string{
		"show arp",
		"show mac addr",
	})

	fmt.Printf("results %v\n", results)

	for _, ret := range results {
		//fetch first line
		lines := strings.Split(ret, dev.LineBreak)

		if len(lines) > 0 {
			fmt.Printf("line = %s\n", lines[0])
		}

		if len(lines) > 0 && lines[0] == "disp ip int bri" {
			dev.ParseIpInterface(lines)
		}

		if len(lines) > 0 && lines[0] == "disp int bri" {
			dev.ParseInterface(lines)
		}

		if len(lines) > 0 && lines[0] == "show arp" {
			dev.ParseArp(lines)
		}

		if len(lines) > 0 && lines[0] == "show mac addr" {
			dev.ParseMacAddr(lines)
		}
	}
}

func (dev *C2960) Dump() {

	s, err := json.Marshal(dev.InterfaceList)
	if err == nil {
		fmt.Printf("%s\n", s)
	} else {
		fmt.Printf("Err=%v\n", err)
	}

	//sIn, sOut, downCnt, upCnt := dev.SumUtil(dev.UpStreamIf)
	//fmt.Printf("SumIn=%s SumOut=%s downCnt=%d upCnt=%d\n", sIn, sOut, downCnt, upCnt)

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

	/*s3, err3 := json.Marshal(dev.IfArpCounts)
	if err3 == nil {
		fmt.Printf("%s\n", s3)
	} else {
		fmt.Printf("Err=%v\n", err3)
	}*/

	/*s4, err4 := json.Marshal(dev.PhoneTable)
	if err4 == nil {
		fmt.Printf("PhoneTable count=%d %s\n", dev.PhoneTable.Len(), s4)
	} else {
		fmt.Printf("Err=%v\n", err4)
	}*/
	dev.DumpMacTables()
}

func (dev *C2960) DumpMacTables() {
	s3, err3 := json.Marshal(dev.MacTable)
	if err3 == nil {
		fmt.Printf("Arp Count=%d %s\n", dev.MacTable.Len(), s3)
	} else {
		fmt.Printf("Err=%v\n", err3)
	}
}

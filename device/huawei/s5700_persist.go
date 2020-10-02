package huawei

import (
	"encoding/json"
	"fmt"
	"gms/utils"

	"github.com/gogf/gf/os/gfile"
)

func (dev *S5700) Save() {

	s, err := json.Marshal(dev.InterfaceList)
	if err == nil {
		gfile.PutBytes("cache"+utils.PathSep()+dev.Host+"-interface.json", s)
	} else {
		fmt.Printf("Err=%v\n", err)
	}

	s1, err1 := json.Marshal(dev.InterfaceIpList)
	if err1 == nil {
		gfile.PutBytes("cache"+utils.PathSep()+dev.Host+"-ip.json", s1)

	} else {
		fmt.Printf("Err=%v\n", err1)
	}

	s2, err2 := json.Marshal(dev.ArpTable)
	if err2 == nil {
		gfile.PutBytes("cache"+utils.PathSep()+dev.Host+"-arp.json", s2)
	} else {
		fmt.Printf("Err=%v\n", err2)
	}
}

//Load load data from cache file.
func (dev *S5700) Load() {

	ifList := []Interface{}

	intBytes := gfile.GetBytes("cache" + utils.PathSep() + dev.Host + "-interface.json")
	json.Unmarshal(intBytes, &ifList)
	for _, ift := range ifList {
		dev.InterfaceList.Add(ift)
	}

	ipList := []IPInterface{}
	ipBytes := gfile.GetBytes("cache" + utils.PathSep() + dev.Host + "-ip.json")
	json.Unmarshal(ipBytes, &ipList)
	for _, ipt := range ifList {
		dev.InterfaceIpList.Add(ipt)
	}

	arpList := []ArpItem{}
	arpBytes := gfile.GetBytes("cache" + utils.PathSep() + dev.Host + "-arp.json")
	json.Unmarshal(arpBytes, &arpList)
	for _, arp := range arpList {
		dev.ArpTable.Append(arp)
	}

	dev.CalcArpTables()
}

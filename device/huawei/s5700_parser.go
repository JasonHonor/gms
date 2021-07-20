package huawei

import (
	"fmt"
	"gms/utils"
	"strings"

	"github.com/gogf/gf/frame/g"
)

func (dev *S5700) ParseIpInterface(lines []string) {
	var isData bool = false
	for _, li := range lines {
		//skip title line
		if strings.HasPrefix(li, "Interface ") {
			isData = true
			continue
		}

		if isData {
			cols := utils.SplitColumns(li, " ")
			if len(cols) > 3 {

				ipIf := IPInterface{}
				ipIf.Name = cols[0]
				ipIf.IPAndMask = cols[1]
				ipIf.PhysicalState = cols[2]
				ipIf.ProtocalState = cols[3]

				//save interface info
				dev.InterfaceIpList.Add(ipIf)
			}
		}
	}
}

func (dev *S5700) ParseInterface(lines []string) {
	var isData bool = false
	for _, li := range lines {
		//skip title line
		if strings.HasPrefix(li, "Interface ") {
			isData = true
			continue
		}

		if isData {
			cols := utils.SplitColumns(li, " ")
			//fmt.Printf("Cols:%v\n", cols)

			if len(cols) > 6 {

				intf := Interface{}
				intf.Name = cols[0]
				intf.PhysicalState = cols[1]
				intf.ProtocalState = cols[2]
				intf.InUti = cols[3]
				intf.OutUti = cols[4]
				intf.InErrors = cols[5]
				intf.OutErrors = cols[6]

				//save interface info
				dev.InterfaceList.Add(intf)
			}
		}
	}
}

func (dev *S5700) TranslateInterface(name string) string {
	if dev.UpStreamIf == name {
		return "UpStream"
	} else {
		return name
	}
}

func (dev *S5700) ParseArp(lines []string) {
	fmt.Printf("Parse arp table,line-count=%v\n", len(lines))

	var sPhonePrefix string

	oPhone := g.Cfg().Get("phone")

	if oPhone != nil {
		sPhonePrefix = oPhone.(string)
	}

	var isData bool = false
	var lastId int
	for _, li := range lines {
		//skip title line
		if strings.HasPrefix(li, "IP ADDRESS ") {
			isData = true
			continue
		}

		if isData {
			cols := utils.SplitColumns(li, " ")
			//fmt.Printf("Cols:%v\n", cols)

			if len(cols) > 4 {

				arp := utils.ArpItem{}
				arp.Sys = dev.Host
				arp.IP = cols[0]
				arp.Mac = utils.MacFormat(cols[1])
				arp.Vlan = cols[3]
				arp.Interface = dev.TranslateInterface(cols[4])

				//save interface info
				dev.ArpTable.Append(arp)
				lastId = dev.ArpTable.Len() - 1

				if len(sPhonePrefix) > 0 && strings.HasPrefix(arp.Mac, sPhonePrefix) {
					dev.PhoneTable.Append(arp)
				}
			} else {
				obj, found := dev.ArpTable.Get(lastId)
				if found {
					if len(cols) > 0 {
						arp := obj.(utils.ArpItem)
						arp.Vlan = cols[0]
						arp.Sys = dev.Host
						dev.ArpTable.Set(lastId, arp)
					}
				}
			}
		}
	}
}

func (dev *S5700) ParseMacAddr(lines []string) {
	fmt.Printf("Parse mac table,line-count=%v\n", len(lines))

	var isData bool = false
	for idx, li := range lines {
		//skip title line
		if idx == 5 {
			isData = true
			continue
		}

		if isData {
			cols := utils.SplitColumns(li, " ")
			//fmt.Printf("Cols:%v\n", cols)

			if len(cols) > 5 {

				arp := utils.ArpItem{}
				arp.Sys = dev.Host
				arp.IP = ""
				arp.Mac = utils.MacFormat(strings.Trim(cols[0], " "))
				arp.Vlan = strings.Trim(cols[1], " ")
				arp.Interface = dev.TranslateInterface(strings.Trim(cols[4], " "))

				dev.UpdateArpTable(arp.Mac, arp.IP, arp.Interface)
			}
		}
	}
}

func (dev *S5700) UpdateArpTable(mac, vlan, intf string) {
	var bFound bool = false
	dev.ArpTable.Iterator(func(k int, v interface{}) bool {
		arpItem, ok1 := v.(utils.ArpItem)
		if ok1 {
			if strings.Contains(arpItem.Mac, mac) {
				arpItem.Interface = intf
				arpItem.Vlan = vlan
				arpItem.Sys = dev.Host
				dev.ArpTable.Set(k, arpItem)
				bFound = true
				return false
			}
		}
		return true
	})

	if !bFound {
		item := utils.ArpItem{}
		item.Interface = intf
		item.Mac = mac
		item.Vlan = vlan
		item.Sys = dev.Host

		dev.ArpTable.Append(item)
	}
}

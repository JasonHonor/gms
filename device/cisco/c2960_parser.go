package cisco

import (
	//	"fmt"
	"gms/utils"
	"strings"
)

func (dev *C2960) ParseIpInterface(lines []string) {
	var isData bool = false
	for _, li := range lines {
		//skip title line
		if strings.HasPrefix(li, "Protocol  Address ") {
			isData = true
			continue
		}

		if isData {
			cols := utils.SplitColumns(li, " ")

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

func (dev *C2960) ParseInterface(lines []string) {
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

func (dev *C2960) ParseArp(lines []string) {
	//fmt.Printf("ParseMacAddr\n")

	var isData bool = false
	for _, li := range lines {
		//skip title line
		if strings.HasPrefix(li, "Protocol  Address          Age (min)  Hardware Addr   Type   Interface") {
			isData = true
			continue
		}

		if isData {
			cols := utils.SplitColumns(li, " ")
			//fmt.Printf("Cols:%v\n", cols)

			if len(cols) > 4 {

				arp := ArpItem{}
				arp.IP = strings.Trim(cols[1], " ")
				arp.Mac = strings.Trim(cols[3], " ")
				arp.Vlan = strings.Trim(cols[5], " ")

				//save interface info
				dev.ArpTable.Append(arp)

				/*
					if len(sPhonePrefix) > 0 && strings.Contains(arp.Mac, sPhonePrefix) {
						dev.PhoneTable.Append(arp)
					}
				*/
			}
		}
	}
}

func (dev *C2960) UpdateArpTable(mac, vlan, intf string) {
	var bFound bool = false
	dev.ArpTable.Iterator(func(k int, v interface{}) bool {
		arpItem, ok1 := v.(ArpItem)
		if ok1 {
			if strings.Contains(strings.ToUpper(arpItem.Mac), strings.ToUpper(mac)) {
				arpItem.Interface = intf
				arpItem.Vlan = vlan
				dev.ArpTable.Set(k, arpItem)
				bFound = true
				return false
			}
		}
		return true
	})

	if !bFound {
		item := ArpItem{}
		item.Interface = intf
		item.Mac = mac
		item.Vlan = vlan

		dev.ArpTable.Append(item)
	}
}

func (dev *C2960) ParseMacAddr(lines []string) {

	//fmt.Printf("ParseMacAddr\n")

	//sPhonePrefix := g.Cfg().Get("phone").(string)

	var isData bool = false
	for _, li := range lines {
		//skip title line
		if strings.Contains(li, "----    -----------       --------    -----") {
			isData = true
			continue
		}

		if isData {
			cols := utils.SplitColumns(li, "    ")
			//fmt.Printf("Cols:%v\n", cols)

			if len(cols) > 3 {

				arp := ArpItem{}
				arp.IP = ""
				arp.Mac = strings.Trim(cols[1], " ")
				arp.Vlan = strings.Trim(cols[0], " ")
				arp.Interface = strings.Trim(cols[3], " ")

				dev.UpdateArpTable(strings.Trim(cols[1], " "), strings.Trim(cols[0], " "), strings.Trim(cols[3], " "))

				//save interface info
				//dev.ArpTable.Append(arp)

				//if len(sPhonePrefix) > 0 && strings.Contains(arp.Mac, sPhonePrefix) {
				//dev.PhoneTable.Append(arp)
				//}
			}
		}
	}
}

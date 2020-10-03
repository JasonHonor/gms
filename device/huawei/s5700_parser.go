package huawei

import (
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

func (dev *S5700) ParseArp(lines []string) {

	sPhonePrefix := g.Cfg().Get("phone").(string)

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

				arp := ArpItem{}
				arp.IP = cols[0]
				arp.Mac = cols[1]
				arp.Vlan = cols[3]
				arp.Interface = cols[4]

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
						arp := obj.(ArpItem)
						arp.Vlan = cols[0]
						dev.ArpTable.Set(lastId, arp)
					}
				}
			}
		}
	}
}

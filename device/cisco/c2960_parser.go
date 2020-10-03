package cisco

import (
	"gms/utils"
	"strings"

	"github.com/gogf/gf/frame/g"
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

	sPhonePrefix := g.Cfg().Get("phone").(string)

	var isData bool = false
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
				arp.IP = cols[1]
				arp.Mac = cols[3]
				arp.Vlan = cols[5]

				//save interface info
				dev.ArpTable.Append(arp)

				if len(sPhonePrefix) > 0 && strings.Contains(arp.Mac, sPhonePrefix) {
					dev.PhoneTable.Append(arp)
				}
			}
		}
	}
}

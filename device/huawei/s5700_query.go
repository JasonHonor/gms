package huawei

import (
	"fmt"
	"io/ioutil"
	"strings"
)

//FindInterfaceByIP find interface name by ip
func (dev *S5700) FindInterfaceByIP(ip string) string {
	var interfaceName string
	dev.ArpTable.Iterator(func(k int, v interface{}) bool {
		arpItem, ok1 := v.(ArpItem)
		if ok1 {
			if arpItem.IP == ip {
				interfaceName = arpItem.Interface
				return false
			}
		}
		return true
	})
	return interfaceName
}

//FindArpListByMac find arp items by mac substring.
func (dev *S5700) FindArpListByMac(mac string) []ArpItem {

	arpItemList := []ArpItem{}

	dev.ArpTable.Iterator(func(k int, v interface{}) bool {
		arpItem, ok1 := v.(ArpItem)
		if ok1 {
			if strings.Contains(strings.ToUpper(arpItem.Mac), strings.ToUpper(mac)) {
				arpItemList = append(arpItemList, arpItem)
			}
		}
		return true
	})

	return arpItemList
}

func (dev *S5700) FindArpListByMacVlan(mac, vlan string) []ArpItem {

	arpItemList := []ArpItem{}

	dev.ArpTable.Iterator(func(k int, v interface{}) bool {
		arpItem, ok1 := v.(ArpItem)
		if ok1 {
			if strings.Contains(strings.ToUpper(arpItem.Mac), strings.ToUpper(mac)) && arpItem.Vlan == vlan {
				arpItemList = append(arpItemList, arpItem)
			}
		}
		return true
	})

	return arpItemList
}

func (dev *S5700) PrintDhcpLeaseFile(mac, vlan string) {

	arpList := dev.FindArpListByMacVlan(mac, vlan)

	var sContent string
	for _, arpItem := range arpList {
		var sMac string
		sMac = strings.ReplaceAll(arpItem.Mac, "-", "")
		sContent += fmt.Sprintf("%s:%s:%s:%s:%s:%s %s %s\n",
			sMac[:2], sMac[2:4], sMac[4:6], sMac[6:8], sMac[8:10], sMac[10:12],
			arpItem.IP, "2025-10-24T19:27:51+08:00")
	}

	fileName := "leases.txt"
	var err error
	if err = ioutil.WriteFile(fileName, []byte(sContent), 0666); err != nil {
		fmt.Println("Writefile Error =", err)
		return
	}
}

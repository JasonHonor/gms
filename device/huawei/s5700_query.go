package huawei

import "strings"

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

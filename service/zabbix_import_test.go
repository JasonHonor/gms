package service

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/AlekSi/zabbix"

	"github.com/gogf/gf/os/gfile"
)

var ch chan bool

func processText(str string) {
	//fmt.Printf("processText %s\n", str)

	sHosts := strings.Split(str, ",")
	sIP := sHosts[12]

	sIPParts := strings.Split(sIP, ".")

	if len(sIPParts) < 4 {
		return
	}

	sDept := sHosts[1]
	sName := sHosts[4] + sIPParts[2] + "." + sIPParts[3]
	sStaffNo := sHosts[5]
	mac := sHosts[11]

	api := getAPI(nil)

	host, err := api.HostGetByHost(sIP)
	if err != nil {
		fmt.Printf("GetHostError %v  %v\n", err, sIP)
		return
	} else {
		fmt.Printf("Host %v found <%s %s %s %s %s>!\n", host, sDept, sName, sStaffNo, sIP, mac)

		if host != nil {
			//found host.
			inv := host.Inventory
			inv.MacAddressA = mac
			inv.Model = GetModelName(mac)
			inv.Alias = sName

			//fmt.Printf("Host found %v\n", host)
			hostUpdate := zabbix.HostUpdate{
				HostId:    host.HostId,
				Name:      sName,
				Inventory: inv,
			}
			err := api.HostsUpdate(hostUpdate)

			if err != nil {
				log.Printf("HostUpdateError %v\n", err)
			}
			//return host
		} /*else {
			err := api.HostsCreate(hosts)
			if err != nil {
				t.Fatal(err)
			}
			return &hosts[0]
		}*/
	}
	//
}

func TestImportCSV(t *testing.T) {

	ch = make(chan bool)

	err := gfile.ReadLines("2020.7-W.csv", processText)

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}
}

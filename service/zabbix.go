package service

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	//	"encoding/json"
	//	"errors"
	"fmt"

	"github.com/AlekSi/zabbix"

	"math/rand"
	"testing"
)

var (
	_host string
	_api  *zabbix.API
)

func getAPI(t *testing.T) *zabbix.API {
	if _api != nil {
		return _api
	}

	url, user, password := "http://172.20.83.99/ui/api_jsonrpc.php", "Admin", "zabbix"
	_api = zabbix.NewAPI(url)
	_api.SetClient(http.DefaultClient)
	v := "1"
	if v != "" && v != "0" {
		_api.Logger = log.New(os.Stderr, "[zabbix] ", 0)
	}

	if user != "" {
		auth, err := _api.Login(user, password)
		if err != nil {
			t.Fatal(err)
		}
		if auth == "" {
			t.Fatal("Login failed")
		}
	}

	return _api
}

func getHost() string {
	return _host
}

func CreateHost(group *zabbix.HostGroup, t *testing.T) *zabbix.Host {
	name := fmt.Sprintf("%s-%d", getHost(), rand.Int())
	iface := zabbix.HostInterface{DNS: name, Port: "42", Type: zabbix.Agent, UseIP: 0, Main: 1}
	hosts := zabbix.Hosts{{
		Host:       name,
		Name:       "Name for " + name,
		GroupIds:   zabbix.HostGroupIds{{group.GroupId}},
		Interfaces: zabbix.HostInterfaces{iface},
	}}

	err := getAPI(t).HostsCreate(hosts)
	if err != nil {
		t.Fatal(err)
	}
	return &hosts[0]
}

func filterByMacPrefixFile(mac, file string) bool {
	var sMac string = ""
	sMac = strings.ReplaceAll(mac, ":", "")
	sMac = strings.ToUpper(sMac)

	//filter by file lines.
	var sFileContent string = ""
	f, err := ioutil.ReadFile(file)
	if err == nil {
		sFileContent = string(f)
	} else {
		return false
	}

	sFileLines := strings.Split(sFileContent, "\n")
	var bFound bool = false
	for _, sFileLine := range sFileLines {
		//fmt.Printf("Filter 1. %s 2. %s\n", sFileLine, sMac)
		if strings.HasPrefix(sMac, sFileLine) {
			bFound = true
			break
		}
	}

	return bFound
}

func GetModelName(mac string) string {

	var model string
	if filterByMacPrefixFile(mac, "avaya.mac.txt") {
		model = "AVAYA 1608-I"
	}

	return model
}

func CreateHostByGroupId(groupId, ip, mac string, t *testing.T) *zabbix.Host {

	api := getAPI(t)

	iface := zabbix.HostInterface{IP: ip, Port: "10050", Type: zabbix.Agent, UseIP: 1, Main: 1}

	var inv zabbix.HostInventory

	if len(mac) > 0 && mac != "Incomplate" {
		inv = zabbix.HostInventory{MacAddressA: mac, Model: GetModelName(mac)}
	}

	hosts := zabbix.Hosts{{
		Host:          ip,
		Name:          ip,
		GroupIds:      zabbix.HostGroupIds{{groupId}},
		Interfaces:    zabbix.HostInterfaces{iface},
		Inventory:     inv,
		InventoryMode: 0,
	}}

	host, _ := api.HostGetByHost(ip)
	if host != nil {
		inv := host.Inventory
		if len(mac) > 0 && mac != "Incomplete" {
			inv.MacAddressA = mac
			inv.Model = GetModelName(mac)

			fmt.Printf("Host found %v\n", host)

			hostUpdate := zabbix.HostUpdate{
				HostId:    host.HostId,
				Inventory: inv,
			}
			err := api.HostsUpdate(hostUpdate)

			if err != nil {
				log.Fatal(err)
			}
		}

		return host
	} else {
		err := api.HostsCreate(hosts)
		if err != nil {
			t.Fatal(err)
		}
		return &hosts[0]
	}
}

func DeleteHost(host *zabbix.Host, t *testing.T) {
	err := getAPI(t).HostsDelete(zabbix.Hosts{*host})
	if err != nil {
		t.Fatal(err)
	}
}

func CreateHostGroup(t *testing.T) *zabbix.HostGroup {
	hostGroups := zabbix.HostGroups{{Name: fmt.Sprintf("zabbix-testing-%d", rand.Int())}}
	err := getAPI(t).HostGroupsCreate(hostGroups)
	if err != nil {
		t.Fatal(err)
	}
	return &hostGroups[0]
}

func DeleteHostGroup(hostGroup *zabbix.HostGroup, t *testing.T) {
	err := getAPI(t).HostGroupsDelete(zabbix.HostGroups{*hostGroup})
	if err != nil {
		t.Fatal(err)
	}
}

func GetHostGroupIdByName(t *testing.T, vlan string) string {
	groups, err := getAPI(t).HostGroupsGet(zabbix.Params{"filter": map[string][]string{"name": []string{fmt.Sprintf("办公网%s", vlan)}}})
	if err != nil {
		t.Fatal(err)
	}

	if len(groups) > 0 {
		return groups[0].GroupId
	} else {
		return ""
	}
}

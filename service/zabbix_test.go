package service

import (
	//	"encoding/json"
	//	"errors"
	"fmt"

	"github.com/AlekSi/zabbix"

	"reflect"
	"testing"
)

func TestHosts(t *testing.T) {
	return

	api := getAPI(t)

	group := CreateHostGroup(t)
	defer DeleteHostGroup(group, t)

	hosts, err := api.HostsGetByHostGroups(zabbix.HostGroups{*group})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 0 {
		t.Errorf("Bad hosts: %#v", hosts)
	}

	host := CreateHost(group, t)
	if host.HostId == "" || host.Host == "" {
		t.Errorf("Something is empty: %#v", host)
	}
	host.GroupIds = nil
	host.Interfaces = nil

	host2, err := api.HostGetByHost(host.Host)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(host, host2) {
		t.Errorf("Hosts are not equal:\n%#v\n%#v", host, host2)
	}

	host2, err = api.HostGetById(host.HostId)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(host, host2) {
		t.Errorf("Hosts are not equal:\n%#v\n%#v", host, host2)
	}

	hosts, err = api.HostsGetByHostGroups(zabbix.HostGroups{*group})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 1 {
		t.Errorf("Bad hosts: %#v", hosts)
	}

	DeleteHost(host, t)

	hosts, err = api.HostsGetByHostGroups(zabbix.HostGroups{*group})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 0 {
		t.Errorf("Bad hosts: %#v", hosts)
	}
}

func TestListHosts(t *testing.T) {
	//getAPI(t)

	group := GetHostGroupIdByName(t, 65)

	CreateHostByGroupId(group, "192.168.9.9", "mac", t)

	fmt.Printf("Group %v\n", group)
}

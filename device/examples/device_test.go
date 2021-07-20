package examples

import (
	"fmt"
	"testing"

	. "gms/cluster"
	"gms/utils"

	"github.com/gogf/gf/frame/g"
)

func TestSnmpWalkConfigFile(t *testing.T) {

	g.Cfg().SetFileName("config.json")

	//str-array-map by mac addr
	var cluster = Cluster{
		MacMap:  utils.NewArrayMap(),
		IpMap:   utils.NewArrayMap(),
		PortMap: utils.NewArrayMap(),
	}

	cluster.Init()
	cluster.Discovery()
	//cluster.MacMap.Dump()
	/*sMac := utils.MacFormat("00e0-4c68-02ea")
	sRet := cluster.FindByMac(sMac)
	fmt.Printf("FindByMac %s Return\n%v\n", sMac, sRet)

	sMacB := utils.MacFormat("00:e0:4c:68:02:ea")
	sRetB := cluster.FindByMac(sMacB)
	fmt.Printf("FindByMac %s Return\n%v\n", sMacB, sRetB)*/

	sIp := "172.20.68.17"
	sRet2 := cluster.FindByIP(sIp)
	fmt.Printf("FindByIP %s Return %v\n", sIp, sRet2)
}

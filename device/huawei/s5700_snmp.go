package huawei

var oidEnterprise = "1.3.6.1.4.1"
var oidSystem = "1.3.6.1.2.1.1"

type SnmpPublic struct {
	SysName        string
	SysDescription string
	SysUptime      string
	SysObjectID    string

	//Table need to walk
	PhysicalNameTable string
	IfDescrTable      string

	//ArpTable 	PortIndex.IP->Mac
	ArpTable string

	//IpPortTable ifIP->PortIndex
	IpPortTable string
}

type SnmpOids struct {
	SnmpPublic
}

func newSnmpPublic() SnmpPublic {
	return SnmpPublic{
		//system name
		SysName: oidSystem + ".5.0",
		//system discription of software
		SysDescription: oidSystem + ".1.0",
		//system update distance
		SysUptime: oidSystem + ".3.0",
		//object identifer of custom objects.
		SysObjectID: oidSystem + ".2.0",
		//Physical interface table
		PhysicalNameTable: "1.3.6.1.2.1.47.1.1.1.1.7",
		IfDescrTable:      "1.3.6.1.2.1.2.2.1.2",
		ArpTable:          "1.3.6.1.2.1.4.22.1.2",
		IpPortTable:       "1.3.6.1.2.1.4.20.1.2",
	}
}

//Huawei device oids
func NewHuaweiSnmpOids() SnmpOids {
	return SnmpOids{
		SnmpPublic: newSnmpPublic(),
	}
}

//Cisco device oids
func NewCiscoSnmpOids() SnmpOids {
	return SnmpOids{
		SnmpPublic: newSnmpPublic(),
	}
}

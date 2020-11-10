package utils

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	g "github.com/soniah/gosnmp"
)

var SnmpCount int = 0
var Sfilters []FilterItem
var SPrefix string

type FilterItem struct {
	Name  string
	Type  string
	Value string
}

func GetSNMPInfo(host, community string, oids []string) {
	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = host
	g.Default.Community = community
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s %v", i, variable.Name, variable.Type)

		// the Value of each variable returned by Get() implements
		// interface{}. You could do a type switch...
		switch variable.Type {
		case g.ObjectIdentifier:
			fmt.Printf("string: %s\n", variable.Value.(string))
		case g.OctetString:
			fmt.Printf("string: %s\n", variable.Value.(string))
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
}

func WalkSnmpOid(host, community string, oid string, isBinary bool, filters []FilterItem) {

	SnmpCount = 0
	SPrefix = oid
	Sfilters = filters

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = host
	g.Default.Community = community
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	if isBinary {
		err = g.Default.BulkWalk(oid, printBinaryValue)
		if err != nil {
			fmt.Printf("Walk Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		err = g.Default.BulkWalk(oid, printValue)
		if err != nil {
			fmt.Printf("Walk Error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Count=%d\n", SnmpCount)
}

func printValue(pdu g.SnmpPDU) error {

	SnmpCount++
	fmt.Printf("%s(%v)= ", pdu.Name, pdu.Type)

	switch pdu.Type {
	case g.OctetString:
		b := pdu.Value.(string)
		fmt.Printf("STRING: %s\n", string(b))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, g.ToBigInt(pdu.Value))
	}
	return nil
}

func printBinaryValue(pdu g.SnmpPDU) error {

	SnmpCount++
	//fmt.Printf("%s(%v)= ", pdu.Name, pdu.Type)
	//fmt.Printf("%s\n", pdu.Name)

	switch pdu.Type {
	case g.OctetString:
		s := pdu.Value.(string)
		b := []byte(s)
		if len(b) == 6 { //MAC Address

			var str string = ""
			for idx, c := range b {
				str += hex.EncodeToString([]byte{c})
				if idx%2 == 1 && idx < len(b)-1 {
					str += "."
				}
			}

			for _, item := range Sfilters {
				sMac := MacFormat(item.Value)
				if sMac != "" && str == sMac {
					fmt.Printf("%s %s MACADDR: %s\n", item.Name, strings.Replace(pdu.Name, SPrefix, "", -1), str)
				}
			}

			if len(Sfilters) == 0 {
				fmt.Printf("%s MACADDR: %s\n", strings.Replace(pdu.Name, SPrefix, "", -1), str)
			}
		}
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, g.ToBigInt(pdu.Value))
	}
	return nil
}

package utils

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"

	"github.com/gogf/gf/os/gfile"

	"encoding/json"
	//	"reflect"
)

type ArrayMapItem interface {
	Dump() string
}

type ArpItem struct {
	Sys       string
	IP        string
	Mac       string
	Vlan      string
	Interface string
}

func (ai *ArpItem) Dump() string {
	/*if strings.Contains(ai.Interface, "Vlan") || strings.Contains(ai.Interface, "0/47") ||
		strings.Contains(ai.Interface, "Eth-Trunk") || strings.Contains(ai.Interface, "CPU") {
		return ""
	} else*/{
		return "\t" + ai.Sys + "/" + ai.IP + "/" + ai.Mac + "/" + ai.Interface + "\n"
	}
}

func (ai *ArpItem) DumpByMac() string {
	if strings.Contains(ai.Interface, "Vlan") || strings.Contains(ai.Interface, "0/47") ||
		strings.Contains(ai.Interface, "Eth-Trunk") || strings.Contains(ai.Interface, "CPU") {
		return ""
	} else {
		return "\t" + ai.Sys + "/" + ai.IP + "/" + ai.Interface + "\n"
	}
}

func (ai *ArpItem) DumpByIP() string {
	if strings.Contains(ai.Interface, "Vlan") || strings.Contains(ai.Interface, "0/47") ||
		strings.Contains(ai.Interface, "CPU") {
		return "DumpByIPNotFound"
	} else {
		return "\t" + ai.Sys + "/" + ai.Mac + "/" + ai.Interface + "\n"
	}
}

type ArrayMap struct {
	data *gmap.StrAnyMap
}

func NewArrayMap() *ArrayMap {
	return &ArrayMap{
		data: gmap.NewStrAnyMap(),
	}
}

func (am *ArrayMap) Set(key string, value interface{}) {
	var gArray *garray.Array
	if am.data == nil {
		fmt.Printf("am.Data not created!")
		return
	}
	if am.data.Contains(key) {
		oData := am.data.Get(key)
		if oData != nil {
			gArray = oData.(*garray.Array)
		}
		gArray.Append(value)
	} else {
		gArray = garray.NewArray()
		gArray.Append(value)
		if gArray != nil {
			am.data.Set(key, gArray)
		}
	}
}

func (am *ArrayMap) Get(key string) (*garray.Array, error) {
	if am.data == nil {
		return nil, fmt.Errorf("container am.data not created!")
	}
	oData := am.data.Get(key)
	if oData != nil {
		return oData.(*garray.Array), nil
	} else {
		return nil, fmt.Errorf("key %v not found!", key)
	}
}

func (am *ArrayMap) Dump() {
	fmt.Printf("ItemCount=%v\n", am.data.Size())
	am.data.Iterator(func(key string, value interface{}) bool {
		if value == nil {
			return false
		}

		var ret string
		oData, ok := value.(*garray.Array)
		if ok {
			if oData.Len() > 1 {
				oData.Iterator(func(index int, value interface{}) bool {
					item, ok2 := value.(ArpItem)
					if ok2 {
						ret += item.Dump()
					} else {
						fmt.Printf("Dump:ConvertError2 type mismatch.")
					}
					return true
				})

				if len(ret) > 0 {
					fmt.Printf("Key=%s Ret=%s\n", key, ret)
				}
			}
		} else {
			oValue, ok2 := value.([]interface{})
			//_, ok2 := value.([]map[string]interface{})
			if ok2 {
				for _, mValue := range oValue {
					itemMap, _ := mValue.(map[string]interface{})
					fmt.Printf("Key=%s Ret=%s\n", itemMap["IP"], itemMap["Mac"])
				}

			} else {
				fmt.Printf("Dump:ConvertError type mismatch2.\n")
			}
			fmt.Printf("Dump:ConvertError type mismatch.\n")
		}

		return true
	})
}

func (am *ArrayMap) MarshalJSON() ([]byte, error) {
	var pos int = 0
	var ret string = `{`
	am.data.Iterator(func(k string, v interface{}) bool {

		var vRet string
		gArr := v.(*garray.Array)

		var vPos = 0
		gArr.Iterator(func(ki int, vi interface{}) bool {
			//item
			arpItem := vi.(ArpItem)
			btItem, _ := json.Marshal(arpItem)
			vRet += fmt.Sprintf(`%s`, string(btItem))
			vPos++
			if vPos < gArr.Len() {
				vRet += ","
			}

			return true
		})

		ret += fmt.Sprintf(`"%s":[%s]`, k, vRet)
		pos++
		if pos < am.data.Size() {
			ret += ","
		}
		return true
	})
	ret += `}`

	//fmt.Printf("Mashaled:%v\n", ret)
	return []byte(ret), nil
}

func (am *ArrayMap) UnmarshalJSON(b []byte) error {
	var tmp map[string][]ArpItem

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		fmt.Printf("UnmarshalError %v\n", err)
		return err
	}

	for key, value := range tmp {
		for _, item := range value {
			am.Set(key, item)
		}
	}

	//fmt.Printf("MapData %v\n", tmp)
	return nil
}

func (am *ArrayMap) Save(file string) {
	s, err := json.Marshal(am)
	if err == nil {
		gfile.PutBytes(file, s)
	} else {
		fmt.Printf("SaveErr=%v\n", err)
	}
}

func (am *ArrayMap) Load(file string) {
	am.data = gmap.NewStrAnyMap()
	intBytes := gfile.GetBytes(file)
	err := json.Unmarshal(intBytes, am)
	if err == nil {
		fmt.Printf("LoadOK Size=%v\n", am.data.Size())
	} else {
		fmt.Printf("LoadErr=%v\n", err)
	}
}

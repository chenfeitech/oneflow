package lua_helper

import (
	"encoding/xml"
	"fmt"
)

type XmlNode struct {
	XMLName xml.Name
	Data    string    `xml:",chardata"`
	Nodes   []XmlNode `xml:",any"`
}

var _ = fmt.Println

func (l *iState) Lua_xml_encode(root XmlNode) string {
	xml_bytes, err := xml.Marshal(root)
	if err != nil {
		panic(err)
	}
	return xml.Header + string(xml_bytes)
}

func (l *iState) Lua_xml_decode(xml_str string) (result interface{}) {
	var node XmlNode
	err := xml.Unmarshal([]byte(xml_str), &node)
	if err != nil {
		panic(err)
	}
	_, result = flat(&node)
	return
}

func flat(node *XmlNode) (key string, val interface{}) {
	key = node.XMLName.Local
	if len(node.Nodes) == 0 {
		val = node.Data
	} else {
		sub_nodes := make(map[string]interface{}, 0)
		for _, sn := range node.Nodes {
			k, v := flat(&sn)
			if prev_val, contans := sub_nodes[k]; contans {
				if arr_val, ok := prev_val.([]interface{}); ok {
					sub_nodes[k] = append(arr_val, v)
				} else {
					sub_nodes[k] = []interface{}{prev_val, v}
				}
			} else {
				sub_nodes[k] = v
			}
		}
		val = sub_nodes
	}
	return
}

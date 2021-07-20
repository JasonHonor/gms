package utils

import (
	"bytes"
	"runtime"
	"strings"
)

func SplitColumns(data, seperator string) []string {
	var ret []string
	cols := strings.Split(data, seperator)
	for _, col := range cols {
		if col != seperator && col != "" {
			ret = append(ret, col)
		}
	}
	return ret
}

// RemovePart remove string by begin/end tag.
func RemoveStringByTag(buf []byte, beginTag, endTag string, strPrepend string) []byte {

	data := string(buf)

	var startPos, endPos int = 0, 0
	startPos = strings.Index(data, beginTag)

	if startPos > 0 {

		// remove by last tag.
		endPos = strings.LastIndex(data, endTag)
		distance := endPos - startPos
		total := distance + len(endTag)

		buf2 := buf[startPos+total:]

		buf = append(buf[:startPos], []byte(strPrepend)...)
		startPos += len(strPrepend)

		buf = append(buf[:startPos], buf2...)

		return buf
	} else {
		return buf
	}
}

//PathSep get path seperator for current os.
func PathSep() string {
	if runtime.GOOS == "windows" {
		return "\\"
	} else {
		return "/"
	}
}

func RemoveStrByTagBytes(buf []byte, beginTag, endTag []byte, strPrepend string) []byte {

	var startPos, endPos int = 0, 0
	startPos = bytes.Index(buf, beginTag)

	if startPos < 0 {
		return buf
	} else {

		// remove by last tag.
		endPos = bytes.LastIndex(buf, endTag)
		distance := endPos - startPos
		total := distance + len(endTag)

		buf2 := buf[startPos+total:]

		buf = append(buf[:startPos], []byte(strPrepend)...)
		startPos += len(strPrepend)

		buf = append(buf[:startPos], buf2...)

		return buf
	}
}

func RemoveMacSeperator(mac string) string {
	var sMacIn string
	sMacIn = mac

	sMacIn = strings.TrimLeft(sMacIn, " ")
	sMacIn = strings.TrimRight(sMacIn, " ")

	var sNoTagStr string
	if strings.Index(sMacIn, ":") > 0 {
		sNoTagStr = strings.ReplaceAll(sMacIn, ":", "")
	} else if strings.Index(sMacIn, "-") > 0 {
		sNoTagStr = strings.ReplaceAll(sMacIn, "-", "")
	} else if strings.Index(sMacIn, ".") > 0 {
		sNoTagStr = strings.ReplaceAll(sMacIn, ".", "")
	}
	return sNoTagStr
}

func FormatMacAddrWord(mac string) string {
	if len(mac) < 12 {
		return ""
	} else {
		return strings.ToUpper(mac[:4] + "." + mac[4:8] + "." + mac[8:12])
	}
}

func MacFormat(mac string) string {
	return FormatMacAddrWord(RemoveMacSeperator(mac))
}

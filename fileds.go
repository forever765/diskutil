package diskutil

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// 解析MegaCli64返回的字段
func parseFiled(line, filed string, targetType int) (interface{}, error) {
	fileds := strings.SplitN(line, ":", 2)
	if len(fileds) != 2 {
		return nil, errors.New("format illegal: " + line)
	}

	// data为全量vd字段
	data := strings.TrimSpace(fileds[1])
	// Raw Size 需要特别处理
	if filed == "Raw Size" {
		reg := regexp.MustCompile(` \[0x.*Sectors\]`)
		value := reg.ReplaceAllString(line, "")
		value = strings.TrimLeft(value, "Raw Size: ")
		return value, nil
	}

	if targetType == typeString {
		return data, nil
	} else if targetType == typeInt {
		// 很多杂牌机是认不出Enclosure Device ID的
		if data == "N/A" {
			return 999, nil
		}
		value, err := strconv.ParseInt(data, 10, 0)
		if err != nil {
			return nil, err
		}
		return int(value), nil
	} else if targetType == typeUint64 {
		value, err := strconv.ParseUint(data, 10, 0)
		if err != nil {
			return nil, err
		}
		return value, nil
	}
	return nil, errors.New("type not supported")
}

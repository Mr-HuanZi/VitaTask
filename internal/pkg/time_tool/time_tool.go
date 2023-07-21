package time_tool

import (
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"time"
)

// SelectUnixTimeMode 根据 mode 选择返回的Unix格式
func SelectUnixTimeMode(in time.Time, mode string) int64 {
	// 模式选择
	switch mode {
	case "milli": // 毫秒
		return in.UnixMilli()
	case "micro": // 微秒
		return in.UnixMicro()
	case "nano": // 纳秒
		return in.UnixNano()
	case "sec": // 秒
	default:
		return in.Unix()
	}
	return 0
}

// ParseTimeRangeToUnix 解析多个字符串时间到Unix数字时间戳
func ParseTimeRangeToUnix(t []string, format string, mode string) ([]int64, error) {
	// 数组默认长度为Slice长度,后面append时,不需要重新申请内存和拷贝,效率很高
	unix := make([]int64, len(t))
	j := 0
	for _, v := range t {
		instance, err := time.ParseInLocation(format, v, time.Local)
		if err != nil {
			return nil, err
		}
		unix[j] = SelectUnixTimeMode(instance, mode)
		j++
	}
	return unix, nil
}

// ParseStartEndTimeToUnix 解析字符串格式的开始时间与结束时间到数字时间戳
// 结束时间会被调整为当天的最后1秒
func ParseStartEndTimeToUnix(t []string, inFormat string, mode string) ([]int64, error) {
	tLen := len(t)
	if tLen <= 0 {
		return nil, exception.NewException(response.ElementQuantityTooLittle)
	}
	if tLen > 2 {
		return nil, exception.NewException(response.ElementQuantityTooMany)
	}
	// 指定长度为 2，因为只需要2个
	unix := make([]int64, 2)
	j := 0
	for _, v := range t {
		if j > 1 {
			// 数量大于2就退出循环
			break
		}
		instance, err := time.ParseInLocation(inFormat, v, time.Local)
		if err != nil {
			return nil, err
		}

		// 第二个元素的时间调整为当天最后1秒
		if j == 1 {
			instance = time.Date(instance.Year(), instance.Month(), instance.Day(), 23, 59, 59, 0, instance.Location())
		}
		unix[j] = SelectUnixTimeMode(instance, mode)

		j++
	}
	return unix, nil
}

// ChangeFormat 改变时间格式
func ChangeFormat(inFormat, outFormat, v string) (string, error) {
	instance, err := time.ParseInLocation(inFormat, v, time.Local)
	if err != nil {
		return "", err
	}

	return instance.Format(outFormat), nil
}

package pkg

import (
	dto2 "VitaTaskGo/internal/api/model/dto"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/sirupsen/logrus"
	"regexp"
	"sort"
	"strconv"
)

func Encryption(s string) string {
	var appKey = "K9jTVxRMFoUzAzgbaG3h1vrCKbWFYUZ3"
	md5Str := cryptor.Md5String(s)
	rs := []rune(md5Str)                      // 把md5字符串转换成切片
	start := string(rs[0:6])                  // 开头截取6位
	end := string(rs[len(rs)-10:])            // 结尾截取10位
	md5Str = cryptor.Md5String(start + end)   // 头尾拼接MD5
	return cryptor.Md5String(md5Str + appKey) // 加盐再MD5
}

func ParseStringToUi64(number string) uint64 {
	if i64, err := strconv.ParseUint(number, 10, 64); err != nil {
		return 0
	} else {
		return i64
	}
}

func ParseStringToI64(number string) int64 {
	if i64, err := strconv.ParseInt(number, 10, 64); err != nil {
		return 0
	} else {
		return i64
	}
}

// PassFormat 校验密码格式
func PassFormat(s string) bool {
	match, err := regexp.MatchString("^([a-zA-Z\\d.]){8,16}$", s)
	if err != nil {
		logrus.Errorln("校验密码正则错误：", err)
		return false
	}
	return match
}

func SliceOperator[T dto2.Integer](slice []T, in T, operator string) []T {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	j := 0
	r := make([]T, len(slice))
	for _, v := range slice {
		switch operator {
		case "+":
			r[j] = v + in
		case "-":
			r[j] = v - in
		case "*":
			r[j] = v * in
		case "/":
			r[j] = v / in
		case "|":
			r[j] = v | in
		case "&":
			r[j] = v & in
		case "^":
			r[j] = v ^ in
		case "<<":
			r[j] = v << in
		case ">>":
			r[j] = v >> in
		}
		j++
	}
	return r
}

// SliceUnique 移除数切片中重复的值
// github.com/duke-git/lancet/v2/slice 包提供了功能相同的函数，但此函数性能会更优一些(ChatGPT说的)
// 原理如下(ChatGPT分析的，大致差不多)
// 1.使用Go的sort包对原始切片进行升序排序，这里使用的是快速排序算法，其复杂度为O(nlogn)。排序后，相同的元素现在在一起并且直接相邻。
// 2.声明一个变量j，作为下一个非重复元素要放置的位置。接着，循环遍历排序后的切片，并将当前元素与下一个元素进行比较。
// 如果它们相同，则意味着找到了重复的元素，跳过它；否则，将它移到位置j，并将j加1，以便后续的非重复元素放置在正确的位置。
// 3.返回切片的前j个元素，因为它们是不重复的元素。
//
// 该算法的时间复杂度是O(nlogn)，因为排序的复杂度是O(nlogn)，而遍历并删除重复元素的复杂度是O(n)。
// 该算法的空间复杂度也是O(n)，因为需要创建一个排序后的切片和处理后的切片。
func SliceUnique[T dto2.NumberAndString](slice []T) []T {
	// 如果有0或1个元素，则返回切片本身。
	if len(slice) < 2 {
		return slice
	}

	//  使切片升序排序
	sort.SliceStable(slice, func(i, j int) bool { return slice[i] < slice[j] })

	uniqPointer := 0

	for i := 1; i < len(slice); i++ {
		// 比较当前元素和唯一指针指向的元素
		// 如果它们不相同，则将项写入唯一指针的右侧。
		if slice[uniqPointer] != slice[i] {
			uniqPointer++
			slice[uniqPointer] = slice[i]
		}
	}

	return slice[:uniqPointer+1]
}

// PagedResult 返回通用分页实例
func PagedResult[T any](t []T, total, page int64) *dto2.PagedResult[T] {
	return &dto2.PagedResult[T]{
		Items: t,
		Total: total,
		Page:  page,
	}
}

// EmptyPagedResult 空的通用分页实例
func EmptyPagedResult[T any]() *dto2.PagedResult[T] {
	return &dto2.PagedResult[T]{
		Items: nil,
		Total: 0,
		Page:  1,
	}
}

package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var removeColorCodesRegexp = regexp.MustCompile(`\033\[[0-9]*m`)

// RemoveColorCodes 移除字符串中的 ANSI 颜色代码
func RemoveColorCodes(line string) string {
	return removeColorCodesRegexp.ReplaceAllString(line, "")
}

// 将字符串居中
func Center(s string, width int, fill string) string {
	if len(s) >= width {
		return s
	}
	padding := width - len(s)
	leftPadding := padding / 2
	rightPadding := padding - leftPadding
	return strings.Repeat(fill, leftPadding) + s + strings.Repeat(fill, rightPadding)
}

// 分离主机域名（IP）和端口号
//
// 示例：
// "example.com:8096" 					=> 	"example.com", "8096"
// "[240e:da8:a801:5a47::316]:8096" 	=> 	"240e:da8:a801:5a47::316" "8096"
// "192.168.1.1:8096" 					=> 	"192.168.1.1" "8096"
func SplitHostPort(hostPort string) (host string, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		host, port = host[:colon], host[colon+1:]
	}

	// 对地址是IPv6地址进行处理
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

// MD5编码
//
// 对字符串进行MD5哈希运算, 返回十六进制
func MD5Hash(raw string) string {
	hash := md5.New()
	hash.Write([]byte(raw))
	return hex.EncodeToString(hash.Sum(nil))
}

// 获取服务器入口
//
// 包含协议、主机域名（IP）、端口号（标准端口号可省略）
// Example1: https://example1.com:8920
// Example2: http://example2.com:5224
func GetEndpoint(addr string) string {
	if !strings.HasPrefix(addr, "http") {
		addr = "http://" + addr
	}
	return strings.TrimSuffix(addr, "/")
}

// api_key 为 Emby / Jellyfin 10.x 使用；Jellyfin 12 改为大写 ApiKey；X-Emby-Token 为旧版认证头
var embyAPIKeys = []string{"api_key", "ApiKey", "X-Emby-Token"}

// 从 URL 中查询参数中解析 Emby 的 API 键值对
//
// 以 xxx=xxx 的字符串形式返回
func ResolveEmbyAPIKVPairs(urlString *string) (string, error) {
	if urlString == nil {
		return "", fmt.Errorf("urlString is nil")
	}
	url, err := url.Parse(*urlString)
	if err != nil {
		return "", err
	}
	for quryKey, queryValue := range url.Query() {
		for _, key := range embyAPIKeys {
			if strings.EqualFold(quryKey, key) {
				return fmt.Sprintf("%s=%s", quryKey, queryValue[0]), nil
			}
		}
	}
	return "", nil
}

// 判断字符串是否为整型数字
func isInt[T ~[]byte | ~[]rune | ~string](s T) bool {
	_, err := strconv.Atoi(string(s))
	return err == nil
}

// 在 []string 中找到某个字符串的索引
// 如果未找到，返回 -1
//
// slice: 目标切片
// target: 目标字符串
// caseInsensitive: 是否忽略大小写
// trim: 是否去除空白字符
func FindStringIndex(slice []string, target string, caseInsensitive bool, trim bool) int {
	if trim {
		target = strings.TrimSpace(target)
	}
	for i, str := range slice {
		if trim {
			str = strings.TrimSpace(str)
		}
		if caseInsensitive {
			if strings.EqualFold(str, target) {
				return i
			}
		} else {
			if str == target {
				return i
			}
		}
	}
	return -1
}

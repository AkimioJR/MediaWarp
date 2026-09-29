package utils

import (
	"bytes"
	"regexp"
	"strings"
)

const (
	ASSHeader1 = `[Script Info]
; This is an Advanced Sub Station Alpha v4+ script.
Title:
ScriptType: v4.00+
Collisions: Normal
PlayDepth: 0

[V4+ Styles]`
	ASSHeader2 = `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text`
)

var (
	srtSubtitlesPattern      = regexp.MustCompile(`@\d+@\d{2}:\d{2}:\d{2},\d{3} --> \d{2}:\d{2}:\d{2},\d{3}@`) // 用于文本是否为 SRT 字幕的正则表达式
	srtTimePattern           = regexp.MustCompile(`-?\d\d:\d\d:\d\d`)
	timeFormatPattern        = regexp.MustCompile(`\d(\d:\d{2}:\d{2}),(\d{2})\d`)
	arrowPattern             = regexp.MustCompile(`\s+-->\s+`)
	styleTagStartPattern     = regexp.MustCompile(`<([ubi])>`)
	styleTagEndPattern       = regexp.MustCompile(`</([ubi])>`)
	fontColorTagStartPattern = regexp.MustCompile(`<font\s+color="?#[\w]{2}([\w]{2})([\w]{2})([\w]{2})"?\s*>`)
	fontColorTagEndPattern   = regexp.MustCompile(`</font>`)
)

// 判断字幕是否为 SRT 格式
func IsSRT(content []byte) bool {
	content = bytes.ReplaceAll(content, []byte{'\r'}, []byte{})    // 去除 \r 保证多系统兼容
	content = bytes.ReplaceAll(content, []byte{'\n'}, []byte{'@'}) // 将 \n 替换为 @
	return srtSubtitlesPattern.Match(content)                      // 查找第一个匹配项
}

// 将 SRT 字幕转换成 ASS 字幕
//
// srtText: SRT 格式字幕文本
// style: ASS 字幕样式
func SRT2ASS(srtText []byte, style []string) []byte {
	// 预定义常量
	var (
		newLine        = []byte("\n")
		dialogueStart  = []byte("Dialogue: 0,")
		dialogueSuffix = []byte(",Default,,0,0,0,,")
	)

	srtText = bytes.ReplaceAll(srtText, []byte("\r"), []byte(""))
	var lines [][]byte
	for line := range bytes.SplitSeq(srtText, []byte("\n")) {
		line = bytes.TrimSpace(line)
		if len(line) != 0 {
			lines = append(lines, line)
		}
	}
	var (
		subtitleBuffer         bytes.Buffer     // 字幕缓存区（某一行字幕未完成先存取到此处）
		currentSubtitleContent uint8        = 0 // 一个时间下字幕的行数（2表示这一时间有2行字幕）
		subtitleContent        bytes.Buffer     // 字幕内容（预分配 16K 大小）
		result                 bytes.Buffer
	)
	subtitleBuffer.Grow(1024)
	subtitleContent.Grow(16 * 1024)
	result.Grow(16 * 1024)

	for index, line := range lines {
		if isInt(line) && srtTimePattern.Match(lines[index+1]) { // 这一行是 SRT 字幕的序列数且下一行是时间
			if subtitleBuffer.Len() > 0 {
				subtitleContent.Write(subtitleBuffer.Bytes())
				subtitleContent.Write(newLine)
				subtitleBuffer.Reset() // 清空缓存区
			}
			currentSubtitleContent = 0
			continue
		}

		if srtTimePattern.Match(line) { // 这一行是时间行
			subtitleBuffer.Write(dialogueStart)
			subtitleBuffer.Write(bytes.ReplaceAll(line, []byte("-0"), []byte("0"))) // 替换时间中的负号
			subtitleBuffer.Write(dialogueSuffix)
		} else {
			if currentSubtitleContent == 0 {
				subtitleBuffer.Write(newLine) // 同一时间多行字幕需要在一行中使用字面量 \n 表示换行
			}
			subtitleBuffer.Write(line)
			currentSubtitleContent += 1
		}
	}
	// 最后一行字幕
	subtitleContent.Write(subtitleBuffer.Bytes())
	subtitleContent.Write(newLine)

	content := subtitleContent.Bytes()
	content = timeFormatPattern.ReplaceAll(content, []byte("$1.$2"))                 // 替换时间格式
	content = arrowPattern.ReplaceAll(content, []byte(","))                          // 替换箭头符号
	content = styleTagStartPattern.ReplaceAll(content, []byte(`{\\$11}`))            // 替换样式标签
	content = styleTagEndPattern.ReplaceAll(content, []byte(`{\\$10}`))              // 替换字体颜色标签
	content = fontColorTagStartPattern.ReplaceAll(content, []byte(`{\\c&H$3$2$1&}`)) // 替换字体颜色标签
	content = fontColorTagEndPattern.ReplaceAll(content, []byte(""))                 // 删除字体结束标签

	result.WriteString(ASSHeader1 + "\n")
	result.Write(newLine)
	result.WriteString(strings.Join(style, "\n"))
	result.Write(newLine)
	result.Write(newLine)
	result.WriteString(ASSHeader2)
	result.Write(newLine)
	result.Write(newLine)
	result.Write(content)
	return result.Bytes()
}

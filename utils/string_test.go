package utils_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/AkimioJR/MediaWarp/utils"
)

func TestResolveEmbyAPIKVPairs(t *testing.T) {
	type TestCase struct {
		URI    string
		Result string
	}
	testCases := map[string]TestCase{
		"api_key": {
			"/emby/videos/41/master.m3u8?DeviceId=B5F997A7-B2EA-448A-A610-A21872B0B4DD&MediaSourceId=mediasource_41&PlaySessionId=d69e971d45fc45dda567cc60834813ab&api_key=e12acc0815f74e9da6a86c9e8c2d45d8&VideoCodec=hevc,h264,mpeg4&VideoBitrate=42000000&TranscodingMaxAudioChannels=6&SegmentContainer=ts&MinSegments=2&BreakOnNonKeyFrames=True&TranscodeReasons=ContainerNotSupported",
			"api_key=e12acc0815f74e9da6a86c9e8c2d45d8",
		},
		"api_KEY": {
			"/emby/videos/41/master.m3u8?DeviceId=B5F997A7-B2EA-448A-A610-A21872B0B4DD&MediaSourceId=mediasource_41&PlaySessionId=d69e971d45fc45dda567cc60834813ab&api_KEY=e12acc0815f74e9da6a86c9e8c2d45d8&VideoCodec=hevc,h264,mpeg4&VideoBitrate=42000000&TranscodingMaxAudioChannels=6&SegmentContainer=ts&MinSegments=2&BreakOnNonKeyFrames=True&TranscodeReasons=ContainerNotSupported",
			"api_KEY=e12acc0815f74e9da6a86c9e8c2d45d8",
		},
		"X-Emby-Token": {
			"/emby/Videos/8/stream.strm?mediasourceid=mediasource_31&playsessionid=8b1ce7461411479cbfbf14a9c63b41e6&static=true&X-Emby-Token=539c76dc33fc4935857e027698d685c7",
			"X-Emby-Token=539c76dc33fc4935857e027698d685c7",
		},
		"x-emby-token": {
			"/emby/Videos/8/stream.strm?mediasourceid=mediasource_31&playsessionid=8b1ce7461411479cbfbf14a9c63b41e6&static=true&x-emby-token=539c76dc33fc4935857e027698d685c7",
			"x-emby-token=539c76dc33fc4935857e027698d685c7",
		},
		"blank": {
			"/emby/Videos/8/stream.strm",
			"",
		},
	}

	for caseName, testCase := range testCases {
		t.Run(caseName, func(t *testing.T) {
			result, err := utils.ResolveEmbyAPIKVPairs(&testCase.URI)
			if err != nil {
				t.Errorf("%s 解析发生错误。错误: ", err)
			}
			if testCase.Result != result {
				t.Errorf("%s 解析错误。期望: %s, 实际: %s", caseName, testCase.Result, result)
			}
		})
	}
}

// goos: darwin
// goarch: arm64
// pkg: MediaWarp/utils
// cpu: Apple M1
// BenchmarkStringConcat
// BenchmarkStringConcat-8   	  542830	     25741 ns/op	  275330 B/op	       1 allocs/op
// 测试字符串拼接
func BenchmarkStringConcat(b *testing.B) {
	s := ""
	for i := 0; i < b.N; i++ {
		s += "a"
	}
}

// goos: darwin
// goarch: arm64
// pkg: MediaWarp/utils
// cpu: Apple M1
// BenchmarkByteAppend
// BenchmarkByteAppend-8   	819340761	         1.597 ns/op	       6 B/op	       0 allocs/op
// 测试 []byte 的 append
func BenchmarkByteAppend(b *testing.B) {
	buf := make([]byte, 0)
	for i := 0; i < b.N; i++ {
		buf = append(buf, 'a')
	}
	_ = string(buf)
}

// goos: darwin
// goarch: arm64
// pkg: MediaWarp/utils
// cpu: Apple M1
// BenchmarkStringBuilder
// BenchmarkStringBuilder-8   	430995778	         2.531 ns/op	       5 B/op	       0 allocs/op
// PASS
// ok  	MediaWarp/utils	1.829s
// // 测试 strings.Builder
func BenchmarkStringBuilder(b *testing.B) {
	var sb strings.Builder
	for i := 0; i < b.N; i++ {
		sb.WriteByte('a')
	}
	_ = sb.String()
}

// goos: darwin
// goarch: arm64
// pkg: MediaWarp/utils
// cpu: Apple M1
// BenchmarkMethod1
// BenchmarkMethod1-8   	11925727	        86.93 ns/op	      80 B/op	       2 allocs/op
func BenchmarkMethod1(b *testing.B) {
	header := "Header"
	style := []string{"Style1", "Style2"}
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		buf.WriteString(header + "\n")
		buf.WriteString(strings.Join(style, "\n") + "\n\n")
		buf.WriteString("Footer\n\n")
	}
}

// goos: darwin
// goarch: arm64
// pkg: MediaWarp/utils
// cpu: Apple M1
// BenchmarkMethod2
// BenchmarkMethod2-8   	14530182	        82.06 ns/op	      80 B/op	       2 allocs/op
func BenchmarkMethod2(b *testing.B) {
	header := "Header"
	style := []string{"Style1", "Style2"}
	newLine := []byte{'\n'}
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		buf.WriteString(header)
		buf.Write(newLine)
		buf.WriteString(strings.Join(style, "\n"))
		buf.Write(newLine)
		buf.Write(newLine)
		buf.WriteString("Footer")
		buf.Write(newLine)
		buf.Write(newLine)
	}
}

package handler

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/service"
	"MediaWarp/internal/service/alist"
	"MediaWarp/utils"
	"fmt"
	"path"
	"strings"
	"time"
)

func processHTTPStrmPlaybackInfo(jsonChain *utils.JsonChain, bsePath string, itemId string, id string, directStreamURL *string) {
	startTime := time.Now()
	defer func() {
		logging.Debugf("处理 HTTPStrm %s PlaybackInfo 耗时：%s", id, time.Since(startTime))
	}()

	if !config.HTTPStrm.Proxy {
		jsonChain.Set(
			bsePath+"SupportsDirectStream",
			false,
		).Set(
			bsePath+"SupportsTranscoding",
			false,
		).Delete(
			bsePath + "TranscodingUrl",
		).Delete(
			bsePath + "TranscodingContainer",
		).Delete(
			bsePath + "TranscodingSubProtocol",
		).Delete(
			bsePath + "TrancodeLiveStartIndex",
		)

		var msgs []string
		if directStreamURL != nil {
			msgs = append(msgs, fmt.Sprintf("原直链播放链接: %s", *directStreamURL))
			apikeypair, err := utils.ResolveEmbyAPIKVPairs(directStreamURL)
			if err != nil {
				logging.Warning("解析API键值对失败：", err)
			}
			directStreamURL := fmt.Sprintf("/Videos/%s/stream?MediaSourceId=%s&Static=true&%s", itemId, id, apikeypair)
			jsonChain.Set(
				bsePath+"DirectStreamUrl",
				directStreamURL,
			)
			msgs = append(msgs, fmt.Sprintf("修改直链播放链接为: %s", directStreamURL))
		}
		logging.Infof("%s 强制禁止串流/转码行为，%s", id, strings.Join(msgs, ", "))
	} else {
		logging.Infof("%s 保持原有串流/转码行为", id)
	}
}

func processAlistStrmPlaybackInfo(jsonChain *utils.JsonChain, bsePath string, id string, itemId string, alistAddr string, directStreamURL *string, filepath string, size *int64) {
	startTime := time.Now()
	defer func() {
		logging.Debugf("处理 AlistStrm %s PlaybackInfo 耗时：%s", id, time.Since(startTime))
	}()

	if !config.AlistStrm.Proxy {
		jsonChain.Set(
			bsePath+"SupportsDirectStream",
			false,
		).Set(
			bsePath+"SupportsTranscoding",
			false,
		).Delete(
			bsePath + "TranscodingUrl",
		).Delete(
			bsePath + "TranscodingContainer",
		).Delete(
			bsePath + "TranscodingSubProtocol",
		).Delete(
			bsePath + "TrancodeLiveStartIndex",
		)

		var msgs []string
		if directStreamURL != nil {
			msgs = append(msgs, fmt.Sprintf("原直链播放链接: %s", *directStreamURL))
			url := fmt.Sprintf("/Videos/%s/stream?MediaSourceId=%s&Static=true", id, itemId)

			apikeypair, err := utils.ResolveEmbyAPIKVPairs(&url)
			if err != nil {
				logging.Warning("解析API键值对失败：", err)
			} else {
				url += "&" + apikeypair
			}
			jsonChain.Set(
				bsePath+"DirectStreamUrl",
				url,
			)
			msgs = append(msgs, fmt.Sprintf("修改直链播放链接为: %s", url))
		}
		container := strings.TrimPrefix(path.Ext(filepath), ".")
		jsonChain.Set(
			bsePath+"Container",
			container,
		)
		msgs = append(msgs, fmt.Sprintf("容器为： %s", container))
		logging.Infof("%s 强制禁止串流/转码行为，%s", id, strings.Join(msgs, ", "))
	} else {
		logging.Infof("%s 保持原有串流/转码行为", id)
	}

	if size == nil {
		alistClient, err := service.GetAlistClient(alistAddr)
		if err != nil {
			logging.Warning("获取 AlistClient 失败：", err)
			return
		}
		fsGetData, err := alistClient.FsGet(&alist.FsGetRequest{Path: filepath, Page: 1})
		if err != nil {
			logging.Warning("请求 FsGet 失败：", err)
			return
		}
		jsonChain.Set(
			bsePath+"Size",
			fsGetData.Size,
		)
		logging.Infof("%s 设置文件大小为：%d", id, fsGetData.Size)
	}
}

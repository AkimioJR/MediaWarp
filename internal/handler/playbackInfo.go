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

	var msgs []string

	jsonChain.Set(
		bsePath+"SupportsDirectPlay",
		true,
	)

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
		msgs = append(msgs, "禁止串流/转码行为")
	} else {
		msgs = append(msgs, "保持原有串流/转码行为")
	}

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
	logging.Infof("Media(id: %s) %s", id, strings.Join(msgs, ", "))
}

func processAlistStrmPlaybackInfo(jsonChain *utils.JsonChain, bsePath string, itemId string, id string, alistAddr string, directStreamURL *string, filepath string, size *int64) {
	startTime := time.Now()
	defer func() {
		logging.Debugf("处理 AlistStrm %s PlaybackInfo 耗时：%s", id, time.Since(startTime))
	}()

	jsonChain.Set(
		bsePath+"SupportsDirectPlay",
		true,
	)

	var msgs []string

	container := strings.TrimPrefix(path.Ext(filepath), ".")
	jsonChain.Set(
		bsePath+"Container",
		container,
	).Set(
		bsePath+"SupportsDirectStream",
		true,
	)

	msgs = append(msgs, fmt.Sprintf("容器为： %s", container))

	if !config.AlistStrm.Proxy {
		jsonChain.Set(
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
		msgs = append(msgs, "禁止转码行为")
	} else {
		msgs = append(msgs, "保持原有转码行为")
	}

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

	if size == nil {
		alistClient, err := service.GetAlistClient(alistAddr)
		if err != nil {
			logging.Warning("获取 AlistClient 失败：", err)
		} else {
			fsGetData, err := alistClient.FsGet(&alist.FsGetRequest{Path: filepath, Page: 1})
			if err != nil {
				logging.Warning("请求 FsGet 失败：", err)
			} else {
				jsonChain.Set(
					bsePath+"Size",
					fsGetData.Size,
				)
				msgs = append(msgs, fmt.Sprintf("设置文件大小为： %d", fsGetData.Size))
			}
		}
	}

	logging.Infof("Media(id: %s) %s", id, strings.Join(msgs, ", "))
}

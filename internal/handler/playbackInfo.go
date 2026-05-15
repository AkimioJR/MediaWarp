package handler

import (
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/AkimioJR/MediaWarp/internal/config"
	"github.com/AkimioJR/MediaWarp/internal/logging"
	"github.com/AkimioJR/MediaWarp/internal/service"
	"github.com/AkimioJR/MediaWarp/internal/service/alist"
	"github.com/AkimioJR/MediaWarp/utils"
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

		// 尝试从原始 URL 中提取 MediaSourceId，以解决 Emby 4.9+ 的 mediasource_ 前缀问题
		actualMediaSourceId := id
		if u, err := url.Parse(*directStreamURL); err == nil {
			if qid := u.Query().Get("MediaSourceId"); qid != "" {
				actualMediaSourceId = qid
			}
		}

		apikeypair, err := utils.ResolveEmbyAPIKVPairs(directStreamURL)
		if err != nil {
			logging.Warning("解析API键值对失败：", err)
		}
		newDirectStreamURL := fmt.Sprintf("/Videos/%s/stream?MediaSourceId=%s&Static=true&%s", itemId, actualMediaSourceId, apikeypair)
		jsonChain.Set(
			bsePath+"DirectStreamUrl",
			newDirectStreamURL,
		)
		msgs = append(msgs, fmt.Sprintf("修改直链播放链接为: %s", newDirectStreamURL))
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

		// 尝试从原始 URL 中提取 MediaSourceId，以解决 Emby 4.9+ 的 mediasource_ 前缀问题
		actualMediaSourceId := id
		if u, err := url.Parse(*directStreamURL); err == nil {
			if qid := u.Query().Get("MediaSourceId"); qid != "" {
				actualMediaSourceId = qid
			}
		}

		apikeypair, err := utils.ResolveEmbyAPIKVPairs(directStreamURL)
		if err != nil {
			logging.Warning("解析API键值对失败：", err)
		}
		newDirectStreamURL := fmt.Sprintf("/Videos/%s/stream?MediaSourceId=%s&Static=true&%s", itemId, actualMediaSourceId, apikeypair)
		jsonChain.Set(
			bsePath+"DirectStreamUrl",
			newDirectStreamURL,
		)
		msgs = append(msgs, fmt.Sprintf("修改直链播放链接为: %s", newDirectStreamURL))
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

# MediaSourceId 修复与优化说明文档 (Branch: fix/mediasource-id)

## 1. 问题背景
在 Emby 4.9.x 及更高版本中，服务器对媒体源 ID（MediaSourceId）的处理引入了破坏性变更：
*   **旧版本 (<= 4.8):** `MediaSourceId` 通常与 `ItemID` 相同，且多为纯数字（如 `3983`）。
*   **新版本 (>= 4.9):** 官方强制为媒体源 ID 增加了 `mediasource_` 前缀（如 `mediasource_3983`）。

### 导致现象
MediaWarp 在重写 `DirectStreamUrl` 时，直接使用了从 PlaybackInfo 响应体中获取的原始 ID 字段。在 Emby 4.9+ 环境下，这会导致生成的播放链接中 `MediaSourceId` 丢失前缀或指向错误的 ID，从而触发客户端报错“没有兼容的流”或“404 Not Found”。

---

## 2. 修复方案
我们采用了**动态提取逻辑**代替硬编码判断。

### 核心逻辑
不再尝试猜测 Emby 的版本或手动拼接前缀，而是直接从上游（Emby/Jellyfin）原始生成的 `DirectStreamUrl` 中“窃取”最准确的参数。

**修改位置：** `internal/handler/playbackInfo.go`

**逻辑实现：**
1.  解析上游返回的原生播放链接。
2.  使用 `url.Parse` 提取查询参数中的 `MediaSourceId`。
3.  如果上游链接中存在该参数（无论带不带前缀），MediaWarp 将完整沿用该值进行重定向链接的拼接。
4.  如果提取失败，则退化使用传入的默认 `id` 作为兜底。

### 优点
*   **零配置兼容：** 自动适配 Emby 4.8、4.9 甚至未来的版本。
*   **高鲁棒性：** 只要上游服务器生成的原始链接能播，MediaWarp 生成的重定向链接就一定能播。
*   **无感性能：** 由于重用了解析逻辑，对 CPU 和内存的额外消耗几乎为零。

---

## 3. 其他修复：Jellyfin 参数对齐
在排查过程中，发现 `internal/handler/jellyfin.go` 存在一个严重的参数传递笔误。

**问题：**
原本在调用 `ModifyPlaybackInfo` 处理器时，将 `ItemID` 和 `MediaSourceId` 两个位置都传入了同一个 `*mediasource.ID`。

**修复：**
现已将第一个参数修正为真正的 `*mediasource.ItemID`。这确保了在 Jellyfin 这种 `ItemID` 与 `MediaSourceId` 经常不一致的系统中，生成的 `/Videos/{itemId}/stream` 路径是完全正确的。

---

## 4. 涉及文件
*   `internal/handler/playbackInfo.go`: 引入 `net/url`，增加动态 ID 提取逻辑。
*   `internal/handler/jellyfin.go`: 修正函数调用时的参数对应关系。

---

## 5. 测试建议
1.  **Emby 4.9+ 测试：** 确认点击播放后，跳转地址中的 `MediaSourceId=mediasource_xxxx` 是否完整。
2.  **多版本媒体测试：** 针对同一个 Item 下有多个版本（如 4K 和 1080P）的情况，确认不同版本跳转的 ID 是否各不相同。
3.  **Jellyfin 测试：** 确认 `/Videos/` 后的 ID 与 `MediaSourceId` 参数是否能正确区分。
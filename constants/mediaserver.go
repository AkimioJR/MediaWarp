package constants

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type MediaServerType uint8 // 媒体服务器类型

var InvalidMediaServerErr = errors.New("invalid  MediaServerType")

const (
	EMBY     MediaServerType = iota // 媒体服务器类型：EmbyServer
	JELLYFIN                        // 媒体服务器类型：Jellyfin
	PLEX                            // 媒体服务器类型：Plex
	FNTV                            // 媒体服务器类型：飞牛影视
)

func (m MediaServerType) str() (string, error) {
	switch m {
	case EMBY:
		return "Emby", nil
	case JELLYFIN:
		return "Jellyfin", nil
	case PLEX:
		return "Plex", nil
	case FNTV:
		return "FNTV", nil
	default:
		return "", InvalidMediaServerErr
	}
}

func (m MediaServerType) String() string {
	s, err := m.str()
	if err != nil {
		return "Unknown"
	}
	return s
}

func parseMediaServerTypeStr(s string) (MediaServerType, error) {
	switch s {
	case "Emby":
		return EMBY, nil
	case "Jellyfin":
		return JELLYFIN, nil
	case "Plex":
		return PLEX, nil
	case "FNTV":
		return FNTV, nil
	default:
		return 0, InvalidMediaServerErr
	}
}

func (m *MediaServerType) UnMarshalJSON(data []byte) error {
	var s string
	if err := yaml.Unmarshal(data, &s); err != nil {
		return err
	}
	mt, err := parseMediaServerTypeStr(s)
	if err != nil {
		return fmt.Errorf("%w: %s", err, s)
	}
	*m = mt
	return nil
}

func (m *MediaServerType) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	mt, err := parseMediaServerTypeStr(s)
	if err != nil {
		return fmt.Errorf("%w: %s", err, s)
	}
	*m = mt
	return nil
}

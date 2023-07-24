package minecraft

import (
	"fmt"
	"time"

	"github.com/layou233/ZBProxy/common/mcprotocol"
	"github.com/layou233/ZBProxy/config"
)

func generateKickMessage(s *config.ConfigProxyService, name string) mcprotocol.Message {
	return mcprotocol.Message{
		Color: mcprotocol.White,
		Extra: []mcprotocol.Message{
			{Bold: true, Color: mcprotocol.Red, Text: "ZB"},
			{Bold: true, Text: "Proxy"},
			{Text: " - "},
			{Bold: true, Color: mcprotocol.Gold, Text: "连接被拒绝\n"},

			{Text: "服务器拒绝了您的连接请求。\n"},
			{Text: "原因: "},
			{Color: mcprotocol.LightPurple, Text: "您无权访问此服务。\n"},
			{Text: "请加入QQ群或者下方链接购买。QQ群:850447104\n\n"},

			{
				Color: mcprotocol.Gray,
				Text: fmt.Sprintf("时间: %d | 玩家名称: %s | 路线: %s\n",
					time.Now().UnixMilli(), name, s.Name),
			},
			{Text: "购买地址"},
			{
				Color: mcprotocol.Aqua, UnderLined: true,
				Text: "https://www.xinhai.asia/",
				// ClickEvent: chat.OpenURL("https://github.com/layou233/ZBProxy"),
			},
		},
	}
}

func generatePlayerNumberLimitExceededMessage(s *config.ConfigProxyService, name string) mcprotocol.Message {
	return mcprotocol.Message{
		Color: mcprotocol.White,
		Extra: []mcprotocol.Message{
			{Bold: true, Color: mcprotocol.Red, Text: "ZB"},
			{Bold: true, Text: "Proxy"},
			{Text: " - "},
			{Bold: true, Color: mcprotocol.Gold, Text: "连接被拒绝\n"},

			{Text: "您的连接请求被服务器拒绝.\n"},
			{Text: "原因: "},
			{Color: mcprotocol.LightPurple, Text: "超出服务在线玩家数量限制。\n"},
			{Text: "请等一等\n\n"},

			{
				Color: mcprotocol.Gray,
				Text: fmt.Sprintf("时间: %d | 玩家名称: %s | 路线名称: %s\n",
					time.Now().UnixMilli(), name, s.Name),
			},
			{Text: "官方QQ群: "},
			{
				Color: mcprotocol.Aqua, UnderLined: true,
				Text: "850447104",
				// ClickEvent: chat.OpenURL("https://github.com/layou233/ZBProxy"),
			},
		},
	}
}

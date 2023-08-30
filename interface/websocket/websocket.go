package websocket

import (
	"container/list"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

type DashboardWebSocket interface {
	Update()
	Socket(c echo.Context) error
}

// シングルトンな構造体
var instance *dashboardWebSocket

type dashboardWebSocket struct {
	// 新規通知有無
	isUpdated bool
	// クライアントのコネクションリスト
	conns list.List
}

func NewDashboardWebSocket() DashboardWebSocket {
	// 単一インスタンスとなるように設定
	if instance == nil {
		instance = new(dashboardWebSocket)
		instance.isUpdated = false
	}
	return instance
}

func (socket *dashboardWebSocket) Update() {
	socket.isUpdated = true
}

func (socket *dashboardWebSocket) Socket(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		var listElm *list.Element

		// すでに登録済みのコネクションか判定
		contain := false
		for con := socket.conns.Front(); con != nil; con = con.Next() {
			if con.Value == ws {
				contain = true
				break
			}
		}
		// 未登録の場合、connectionリストに追加
		if !contain {
			listElm = socket.conns.PushBack(ws)
			// conns.PushBack(ws)
			log.Print("INFO: 新たなクライアントが接続されました。クライアント数: ", socket.conns.Len())
		}

		for {
			// 5sおきに接続確認
			// 接続できない場合はコネクションを切断
			time.Sleep(time.Second * 5)
			err := websocket.Message.Send(ws, "CONNECTION CHECK")
			if err != nil {
				c.Logger().Error(err)
				socket.conns.Remove(listElm)
				ws.Close()
				log.Print("INFO: 受信できないクライアントを削除しました。クライアント数: ", socket.conns.Len())
				break
			}

			// 更新がない場合
			if !socket.isUpdated {
				continue
			}

			// 更新がある場合、全てのクライアントへ更新を通知
			log.Print("INFO: 実行履歴が更新されたためクライアントへ通知します。クライアント数: ", socket.conns.Len())
			for con := socket.conns.Front(); con != nil; con = con.Next() {
				err := websocket.Message.Send(con.Value.(*websocket.Conn), "UPDATED")
				if err != nil {
					c.Logger().Error(err)
					log.Print("INFO: 正常に送信できないクライアントが存在しました。")
				}
			}
			socket.isUpdated = false
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

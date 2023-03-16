package gate

type UserConn struct {
}

type WServer struct {
	wsUserToConn map[string]map[int][]*UserConn
}

func (ws *WServer) onInit(wsPort int) {

}

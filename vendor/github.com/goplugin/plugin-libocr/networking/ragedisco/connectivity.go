package ragedisco

import ragetypes "github.com/goplugin/plugin-libocr/ragep2p/types"

type connectivityMsgType int

const (
	_ connectivityMsgType = iota
	connectivityAdd
	connectivityRemove
)

type connectivityMsg struct {
	msgType connectivityMsgType
	peerID  ragetypes.PeerID
}

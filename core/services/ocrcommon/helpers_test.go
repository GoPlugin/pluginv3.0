package ocrcommon

import ocrnetworking "github.com/goplugin/plugin-libocr/networking"

func (p *SingletonPeerWrapper) PeerConfig() (ocrnetworking.PeerConfig, error) {
	return p.peerConfig()
}

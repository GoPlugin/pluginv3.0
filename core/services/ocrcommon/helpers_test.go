package ocrcommon

import ocrnetworking "github.com/goplugin/libocr/networking"

func (p *SingletonPeerWrapper) PeerConfig() (ocrnetworking.PeerConfig, error) {
	return p.peerConfig()
}

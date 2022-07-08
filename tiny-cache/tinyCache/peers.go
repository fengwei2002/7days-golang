package tinyCache

// PeerPicker 接口拥有 PickPeer 方法， 用于根据传入的 key 选择相应节点 PeerPicker
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 接口的 Get 方法用于从对应的 group 中查找缓存 value ， PeerGetter 对应于 http 客户端
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}

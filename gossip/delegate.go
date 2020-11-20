package gossip

import "encoding/json"

type MyDelegate struct {
	metaData MyMetaData
}

func (d *MyDelegate) NodeMeta(limit int) []byte {
	return d.metaData.Bytes()
}
func (d *MyDelegate) LocalState(join bool) []byte {
	// not use, noop
	return []byte("")
}
func (d *MyDelegate) NotifyMsg(msg []byte) {
	// not use
}
func (d *MyDelegate) GetBroadcasts(overhead, limit int) [][]byte {
	// not use, noop
	return nil
}
func (d *MyDelegate) MergeRemoteState(buf []byte, join bool) {
	// not use
}

type MyMetaData struct {
	ip   string
	port int
}

func (m *MyMetaData) Bytes() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		return []byte("")
	}
	return data
}
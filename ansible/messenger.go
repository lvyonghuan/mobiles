package ansible

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/lvyonghuan/Ubik-Util/uerr"
)

// read a message from a stream
func readFromStream(s network.Stream, v any) error {
	decoder := json.NewDecoder(s)
	err := decoder.Decode(v)
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

package ansible

import (
	"context"
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/lvyonghuan/Ubik-Util/uerr"
	"github.com/lvyonghuan/mobiles/hainish"
)

func (p *peerManager) sendHeartbeat(peerID peer.ID) error {
	stream, err := p.ansible.host().NewStream(context.Background(), peerID, heartbeatProtocol)
	if err != nil {
		if stream != nil {
			stream.Close()
		}
		return uerr.NewError(err)
	}
	defer stream.Close()

	_, err = stream.Write([]byte(""))
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

func (p *peerManager) sendPluginInfoToLeader(leaderID peer.ID) error {
	stream, err := p.ansible.host().NewStream(context.Background(), leaderID, identityConfirmationProtocol)
	if err != nil {
		if stream != nil {
			stream.Close()
		}
		return uerr.NewError(err)
	}
	defer stream.Close()

	// Encode the plugin metadata to JSON and send it
	jsonData, err := json.Marshal(p.ansible)
	if err != nil {
		return uerr.NewError(err)
	}

	// Send the JSON data over the stream
	_, err = stream.Write(jsonData)
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

func (p *peerManager) sendProcessDataToFollower(data hainish.Edge) error {
	stream, err := p.ansible.host().NewStream(context.Background(), data.Destination, passingDataProtocol)
	if err != nil {
		if stream != nil {
			stream.Close()
		}
		return uerr.NewError(err)
	}

	defer stream.Close()
	// Encode the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return uerr.NewError(err)
	}
	//Send the data
	_, err = stream.Write(jsonData)
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

func (p *peerManager) sendLogToLeader(level int, message string) error {
	stream, err := p.ansible.host().NewStream(context.Background(), p.ansible.getLeader(), logUploadProtocol)
	if err != nil {
		if stream != nil {
			stream.Close()
		}
		return uerr.NewError(err)
	}

	defer stream.Close()
	logMsg := logMessage{
		Level:   level,
		Message: message,
	}

	// Encode the data to JSON
	jsonData, err := json.Marshal(logMsg)
	if err != nil {
		return uerr.NewError(err)
	}
	//Send the
	_, err = stream.Write(jsonData)
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

func (p *peerManager) sendResultToLeader(result any) error {
	stream, err := p.ansible.host().NewStream(context.Background(), p.ansible.getLeader(), resultUploadProtocol)
	if err != nil {
		if stream != nil {
			stream.Close()
		}
		return uerr.NewError(err)
	}

	defer stream.Close()
	// Encode the data to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return uerr.NewError(err)
	}
	//Send the data
	_, err = stream.Write(jsonData)
	if err != nil {
		return uerr.NewError(err)
	}

	return nil
}

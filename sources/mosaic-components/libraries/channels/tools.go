

package channels


import "net"

import "vgl/transcript"


func CreateFromConnection (_callbacks Callbacks, _connection net.Conn) (Channel, error) {
	return Create (_callbacks, _connection, _connection, _connection)
}

func CreateAndDial (_callbacks Callbacks, _network string, _address string) (Channel, error) {
	if _connection, _error := net.Dial (_network, _address); _error != nil {
		return nil, _error
	} else {
		return CreateFromConnection (_callbacks, _connection)
	}
	panic ("fallthrough")
}


var _packageTranscript = transcript.NewPackageTranscript ()

var useTranscript = false

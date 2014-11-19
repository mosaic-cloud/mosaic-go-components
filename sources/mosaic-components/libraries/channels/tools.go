

package channels


import "net"

import "vgl/transcript"


func CreateFromConnection (_callbacks Callbacks, _connection_ net.Conn) (Channel, error) {
	_connection := & connection { connection : _connection_ }
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


type connection struct {
	connection net.Conn
}

func (_connection *connection) Read (data []byte) (int, error) {
	return _connection.connection.Read (data)
}

func (_connection *connection) Write (data []byte) (int, error) {
	return _connection.connection.Write (data)
}

func (_connection *connection) Sync () (error) {
	return nil
}

func (_connection *connection) Close () (error) {
	return _connection.connection.Close ()
}


var useTranscript = false

var _packageTranscript = transcript.NewPackageTranscript (transcript.InformationLevel)

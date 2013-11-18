

package channels


import "encoding/binary"
import "fmt"
import "io"


func (_channel *channel) pullInboundPacket () (*Packet, error) {
	
	if useTranscript { _channel.transcript.TraceDebugging ("inputing an inbound packet...") }
	
	if useTranscript { _channel.transcript.TraceDebugging ("inputing the packet size...") }
	var _size uint32
	if _error := binary.Read (_channel.inboundStream, binary.BigEndian, &_size); _error != nil {
		return nil, _error
	}
	if _size == 0 {
		return nil, fmt.Errorf ("invalid inbound packet (zero payload)")
	}
	
	if useTranscript { _channel.transcript.TraceDebugging ("inputing the packet payload (%d)...", _size) }
	var _payload []byte = make ([]byte, _size)
	if _readSize, _error := io.ReadFull (_channel.inboundStream, _payload); _error != nil {
		return nil, _error
	} else if _readSize != int (_size) {
		panic ("assertion")
	}
	
	if useTranscript { _channel.transcript.TraceDebugging ("decoding the packet payload...") }
	var _packet *Packet
	if _packet_1, _error := DecodePacket (_payload); _error != nil {
		return nil, _error
	} else {
		_packet = _packet_1
	}
	
	if useTranscript { _channel.transcript.TraceDebugging ("completed inputing the inbound packet `%#v`.", _packet) }
	return _packet, nil
}


func (_channel *channel) pushOutboundPacket (_packet *Packet) (error) {
	
	if useTranscript { _channel.transcript.TraceDebugging ("outputing an outbound packet...") }
	
	if useTranscript { _channel.transcript.TraceDebugging ("encoding the packet payload...") }
	var _payload []byte
	if _payload_1, _error := EncodePacket (_packet); _error != nil {
		return _error
	} else {
		_payload = _payload_1
	}
	var _size uint32 = uint32 (len (_payload))
	
	if useTranscript { _channel.transcript.TraceDebugging ("outputing the packet size and payload (%d)...", _size) }
	if _error := binary.Write (_channel.outboundStream, binary.BigEndian, _size); _error != nil {
		return _error
	}
	if _writeSize, _error := _channel.outboundStream.Write (_payload); _error != nil {
		return _error
	} else if _writeSize != int (_size) {
		panic ("assertion")
	}
	
	if useTranscript { _channel.transcript.TraceDebugging ("completed outputing the outbound packet `%#v`.", _packet) }
	return nil
}

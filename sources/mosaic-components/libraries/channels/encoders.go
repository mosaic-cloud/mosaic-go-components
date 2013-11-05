
package channels


import "encoding/json"


func EncodePacket (_packet *Packet) ([]byte, error) {
	
	var _rawData []byte
	if _rawData_1, _error := json.Marshal (_packet.Data); _error != nil {
		return nil, _error
	} else {
		_rawData = _rawData_1
	}
	
	var _rawAttachment []byte
	if _packet.Attachment != nil {
		_rawAttachment = ([]byte) (_packet.Attachment)
	} else {
		_rawAttachment = emptyPacketData
	}
	
	return reframePacket (_rawData, _rawAttachment)
}


func reframePacket (_data []byte, _attachment []byte) ([]byte, error) {
	_payload := make ([]byte, len (_data) + 1 + len (_attachment))
	copy (_payload, _data)
	_payload[len (_data)] = 0
	copy (_payload[len (_data) + 1:], _attachment)
	return _payload, nil
}


var emptyPacketData = []byte {}

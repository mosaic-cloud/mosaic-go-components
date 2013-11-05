

package channels


import "bytes"
import "encoding/json"
import "fmt"


func DecodePacket (_payload []byte) (*Packet, error) {
	
	var _rawData []byte
	var _rawAttachment []byte
	if _rawData_1, _rawAttachment_1, _error := deframePacket (_payload); _error != nil {
		return nil, _error
	} else {
		_rawData = _rawData_1
		_rawAttachment = _rawAttachment_1
	}
	
	var _data PacketData
	if _error := json.Unmarshal (_rawData, &_data); _error != nil {
		return nil, _error
	}
	
	var _attachment PacketAttachment
	if len (_rawAttachment) != 0 {
		_attachment = PacketAttachment (_rawAttachment)
	} else {
		_attachment = nil
	}
	
	_packet := & Packet {
			Data : _data,
			Attachment : _attachment,
	}
	
	return _packet, nil
}


func deframePacket (_payload []byte) ([]byte, []byte, error) {
	
	if _dataLimit := bytes.IndexByte (_payload, 0); _dataLimit == -1 {
		return nil, nil, fmt.Errorf ("invalid inbound packet (missing payload delimiter)")
	} else {
		_data := _payload[:_dataLimit]
		_attachment := _payload[_dataLimit + 1:]
		return _data, _attachment, nil
	}
	panic ("fallthrough")
}

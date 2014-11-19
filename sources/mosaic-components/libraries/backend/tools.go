

package backend


import "fmt"
import "net"

import "vgl/transcript"

import . "mosaic-components/libraries/messages"


func TcpSocketAcquireSync (_backend Controller, _identifier ResourceIdentifier) (net.IP, uint16, string, error) {
	_specification := & TcpSocketSpecification {
			Identifier : _identifier,
	}
	var _descriptor *TcpSocketDescriptor
	if _descriptor_1, _error := _backend.ResourceAcquireSync (_specification); _error != nil {
		return nil, 0, "", _error
	} else if _descriptor_2, _typeOk := _descriptor_1.(*TcpSocketDescriptor); !_typeOk {
		return nil, 0, "", fmt.Errorf ("unexpected resource descriptor `%#v`", _descriptor_1)
	} else {
		_descriptor = _descriptor_2
	}
	return _descriptor.Ip, _descriptor.Port, _descriptor.Fqdn, nil
}


var packageTranscript = transcript.NewPackageTranscript (transcript.InformationLevel)

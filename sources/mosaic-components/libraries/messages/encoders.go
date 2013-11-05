

package messages


import "encoding/json"
import "fmt"

import . "mosaic-components/libraries/channels"


func Encode (_message_0 Message) (*Packet, error) {
	var _data PacketData
	var _attachment PacketAttachment
	var _error error
	switch _message := _message_0.(type) {
		case *ComponentCall :
			_data, _attachment, _error = EncodeComponentCall (_message)
		case *ComponentCallReturn :
			_data, _attachment, _error = EncodeComponentCallReturn (_message)
		case *ComponentCast :
			_data, _attachment, _error = EncodeComponentCast (_message)
		case *ComponentRegister :
			_data, _attachment, _error = EncodeComponentRegister (_message)
		case *ComponentRegisterReturn :
			_data, _attachment, _error = EncodeComponentRegisterReturn (_message)
		case *ResourceAcquire :
			_data, _attachment, _error = EncodeResourceAcquire (_message)
		case *ResourceAcquireReturn :
			_data, _attachment, _error = EncodeResourceAcquireReturn (_message)
		case *TranscriptPush :
			_data, _attachment, _error = EncodeTranscriptPush (_message)
		default :
			_error = fmt.Errorf ("unknown message type `%#v`", _message_0)
	}
	if _error != nil {
		return nil, _error
	}
	return & Packet {
			Data : _data,
			Attachment : _attachment,
	}, nil
}


func EncodeComponentCall (_message *ComponentCall) (PacketData, PacketAttachment, error) {
	_data := map[string]interface{} {
			"__type__" : "exchange",
			"action" : "call",
			"component" : _message.Component,
			"operation" : _message.Operation,
			"inputs" : nil,
			"correlation" : _message.Correlation,
	}
	if _encodedInputs, _error := EncodeObject (_message.Inputs); _error != nil {
		return nil, nil, _error
	} else {
		_data["inputs"] = _encodedInputs
	}
	return PacketData (_data), PacketAttachment (_message.Attachment), nil
}

func EncodeComponentCallReturn (_message *ComponentCallReturn) (PacketData, PacketAttachment, error) {
	_data := map[string]interface{} {
			"__type__" : "exchange",
			"action" : "call-return",
			"ok" : _message.Ok,
			"outputs" : nil,
			"error" : nil,
			"correlation" : _message.Correlation,
	}
	if _message.Ok {
		if _encodedOutputs, _error := EncodeObject (_message.Outputs); _error != nil {
			return nil, nil, _error
		} else {
			_data["outputs"] = _encodedOutputs
			delete (_data, "error")
		}
	} else {
		if _encodedError, _error := EncodeObject (_message.Error); _error != nil {
			return nil, nil, _error
		} else {
			_data["error"] = _encodedError
			delete (_data, "outputs")
		}
	}
	return PacketData (_data), PacketAttachment (_message.Attachment), nil
}

func EncodeComponentCast (_message *ComponentCast) (PacketData, PacketAttachment, error) {
	_data := map[string]interface{} {
			"__type__" : "exchange",
			"action" : "cast",
			"component" : _message.Component,
			"operation" : _message.Operation,
			"inputs" : nil,
	}
	if _encodedInputs, _error := EncodeObject (_message.Inputs); _error != nil {
		return nil, nil, _error
	} else {
		_data["inputs"] = _encodedInputs
	}
	return PacketData (_data), PacketAttachment (_message.Attachment), nil
}

func EncodeComponentRegister (_message *ComponentRegister) (PacketData, PacketAttachment, error) {
	_data := map[string]interface{} {
			"__type__" : "exchange",
			"action" : "register",
			"group" : _message.Group,
			"correlation" : _message.Correlation,
	}
	return PacketData (_data), nil, nil
}

func EncodeComponentRegisterReturn (_message *ComponentRegisterReturn) (PacketData, PacketAttachment, error) {
	_data := map[string]interface{} {
			"__type__" : "exchange",
			"action" : "register-return",
			"ok" : _message.Ok,
			"error" : nil,
			"correlation" : _message.Correlation,
	}
	if _message.Ok {
		delete (_data, "error")
	} else {
		if _encodedError, _error := EncodeObject (_message.Error); _error != nil {
			return nil, nil, _error
		} else {
			_data["error"] = _encodedError
		}
	}
	return PacketData (_data), nil, nil
}


func EncodeResourceAcquire (_message *ResourceAcquire) (PacketData, PacketAttachment, error) {
	_data := map[string]interface{} {
			"__type__" : "resources",
			"action" : "acquire",
			"specifications" : nil,
			"correlation" : _message.Correlation,
	}
	_encodedSpecifications := make (map[string]interface{}, len (_message.Specifications))
	for _, _specification := range _message.Specifications {
		if _encodedSpecification, _identifier, _error := EncodeResourceSpecification (_specification); _error != nil {
			return nil, nil, _error
		} else {
			_encodedSpecifications[string (_identifier)] = _encodedSpecification
		}
	}
	_data["specifications"] = _encodedSpecifications
	return PacketData (_data), nil, nil
}

func EncodeResourceAcquireReturn (_message *ResourceAcquireReturn) (PacketData, PacketAttachment, error) {
	_data := map[string]interface{} {
			"__type__" : "resources",
			"action" : "acquire-return",
			"descriptors" : nil,
			"correlation" : _message.Correlation,
	}
	_encodedDescriptors := make (map[string]interface{}, len (_message.Descriptors))
	for _, _descriptor := range _message.Descriptors {
		if _encodedDescriptor, _identifier, _error := EncodeResourceDescriptor (_descriptor); _error != nil {
			return nil, nil, _error
		} else {
			_encodedDescriptors[string (_identifier)] = _encodedDescriptor
		}
	}
	_data["descriptors"] = _encodedDescriptors
	return PacketData (_data), nil, nil
}

func EncodeResourceSpecification (_specification_0 ResourceSpecification) (interface{}, ResourceIdentifier, error) {
	switch _specification := _specification_0.(type) {
		case *TcpSocketSpecification :
			return "socket:ipv4:tcp", _specification.Identifier, nil
		default :
			return nil, NilResourceIdentifier, fmt.Errorf ("unkunown resource specification type `%#v`", _specification_0)
	}
	panic ("fallthrough")
}

func EncodeResourceDescriptor (_descriptor_0 ResourceDescriptor) (interface{}, ResourceIdentifier, error) {
	switch _descriptor := _descriptor_0.(type) {
		case *TcpSocketDescriptor :
			return map[string]interface{} {
					"type" : "socket:ipv4:tcp",
					"ip" : _descriptor.Ip,
					"port" : _descriptor.Port,
					"fqdn" : _descriptor.Fqdn,
			}, _descriptor.Identifier, nil
		default :
			return nil, NilResourceIdentifier, fmt.Errorf ("unknown resource descriptor type `%#v`", _descriptor_0)
	}
}


func EncodeTranscriptPush (_message *TranscriptPush) (PacketData, PacketAttachment, error) {
	_data := map[string]interface{} {
			"__type__" : "transcript",
			"action" : "push",
	}
	return PacketData (_data), PacketAttachment (_message.Data), nil
}


func EncodeObject (_object interface{}) (*json.RawMessage, error) {
	if _data, _error := json.Marshal (_object); _error != nil {
		return nil, _error
	} else {
		return (*json.RawMessage) (&_data), nil
	}
	panic ("fallthrough")
}

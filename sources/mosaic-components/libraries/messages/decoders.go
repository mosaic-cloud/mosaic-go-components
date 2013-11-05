

package messages


import "fmt"
import "net"

import . "mosaic-components/libraries/channels"


func Decode (_rawMessage *Packet) (Message, error) {
	_data := _rawMessage.Data
	_attachment := _rawMessage.Attachment
	var _type string
	if _type_1, _error := extractString (_data, "__type__"); _error != nil {
		return nil, _error
	} else {
		_type = _type_1
	}
	switch _type {
		case "exchange" :
			return decodeExchange (_data, _attachment)
		case "resources" :
			return decodeResources (_data, _attachment)
		case "transcript" :
			return decodeTranscript (_data, _attachment)
		default :
			return nil, fmt.Errorf ("unknown message type `%s`", _type)
	}
	panic ("fallthrough")
}


func decodeExchange (_data PacketData, _attachment PacketAttachment) (Message, error) {
	var _action string
	if _action_1, _error := extractString (_data, "action"); _error != nil {
		return nil, _error
	} else {
		_action = _action_1
	}
	switch _action {
		case "call" :
			return decodeComponentCall (_data, _attachment)
		case "call-return" :
			return decodeComponentCallReturn (_data, _attachment)
		case "cast" :
			return decodeComponentCast (_data, _attachment)
		case "register" :
			return decodeComponentRegister (_data, _attachment)
		case "register-return" :
			return decodeComponentRegisterReturn (_data, _attachment)
		default :
			return nil, fmt.Errorf ("unknown exchange message action `%s`", _action)
	}
	panic ("fallthrough")
}

func decodeComponentCall (_data map[string]interface{}, _attachment PacketAttachment) (*ComponentCall, error) {
	var _component ComponentIdentifier
	if containsKey (_data, "component") {
		if _component_1, _error := extractComponentIdentifier (_data, "component"); _error != nil {
			return nil, _error
		} else {
			_component = _component_1
		}
	} else {
		_component = ""
	}
	var _operation ComponentOperation
	if _operation_1, _error := extractComponentOperation (_data, "operation"); _error != nil {
		return nil, _error
	} else {
		_operation = _operation_1
	}
	var _inputs interface{}
	if _inputs_1, _error := extractObject (_data, "inputs"); _error != nil {
		return nil, _error
	} else {
		_inputs = _inputs_1
	}
	var _correlation Correlation
	if _correlation_1, _error := extractCorrelation (_data, "correlation"); _error != nil {
		return nil, _error
	} else {
		_correlation = _correlation_1
	}
	_message := & ComponentCall {
			Component : _component,
			Operation : _operation,
			Inputs : _inputs,
			Correlation : _correlation,
			Attachment : Attachment (_attachment),
	}
	return _message, nil
}

func decodeComponentCallReturn (_data map[string]interface{}, _attachment PacketAttachment) (*ComponentCallReturn, error) {
	var _ok bool
	if _ok_1, _error := extractBool (_data, "ok"); _error != nil {
		return nil, _error
	} else {
		_ok = _ok_1
	}
	var _outputs interface{}
	var _merror interface{}
	if _ok {
		if _outputs_1, _error := extractObject (_data, "outputs"); _error != nil {
			return nil, _error
		} else {
			_outputs = _outputs_1
		}
	} else {
		if _merror_1, _error := extractObject (_data, "outputs"); _error != nil {
			return nil, _error
		} else {
			_merror = _merror_1
		}
	}
	var _correlation Correlation
	if _correlation_1, _error := extractCorrelation (_data, "correlation"); _error != nil {
		return nil, _error
	} else {
		_correlation = _correlation_1
	}
	return & ComponentCallReturn {
			Ok : _ok,
			Outputs : _outputs,
			Error : _merror,
			Correlation : _correlation,
			Attachment : Attachment (_attachment),
	}, nil
}

func decodeComponentCast (_data map[string]interface{}, _attachment PacketAttachment) (*ComponentCast, error) {
	var _component ComponentIdentifier
	if containsKey (_data, "component") {
		if _component_1, _error := extractComponentIdentifier (_data, "component"); _error != nil {
			return nil, _error
		} else {
			_component = _component_1
		}
	} else {
		_component = ""
	}
	var _operation ComponentOperation
	if _operation_1, _error := extractComponentOperation (_data, "operation"); _error != nil {
		return nil, _error
	} else {
		_operation = _operation_1
	}
	var _inputs interface{}
	if _inputs_1, _error := extractObject (_data, "inputs"); _error != nil {
		return nil, _error
	} else {
		_inputs = _inputs_1
	}
	_message := & ComponentCast {
			Component : _component,
			Operation : _operation,
			Inputs : _inputs,
			Attachment : Attachment (_attachment),
	}
	return _message, nil
}

func decodeComponentRegister (_data map[string]interface{}, _attachment PacketAttachment) (*ComponentRegister, error) {
	panic ("not-implemented")
}

func decodeComponentRegisterReturn (_data map[string]interface{}, _attachment PacketAttachment) (*ComponentRegisterReturn, error) {
	var _ok bool
	if _ok_1, _error := extractBool (_data, "ok"); _error != nil {
		return nil, _error
	} else {
		_ok = _ok_1
	}
	var _merror interface{}
	if !_ok {
		if _merror_1, _error := extractObject (_data, "error"); _error != nil {
			return nil, _error
		} else {
			_merror = _merror_1
		}
	}
	var _correlation Correlation
	if _correlation_1, _error := extractCorrelation (_data, "correlation"); _error != nil {
		return nil, _error
	} else {
		_correlation = _correlation_1
	}
	return & ComponentRegisterReturn {
			Ok : _ok,
			Error : _merror,
			Correlation : _correlation,
	}, nil
}


func decodeResources (_data PacketData, _attachment PacketAttachment) (Message, error) {
	var _action string
	if _action_1, _error := extractString (_data, "action"); _error != nil {
		return nil, _error
	} else {
		_action = _action_1
	}
	switch _action {
		case "acquire" :
			return decodeResourceAcquire (_data, _attachment)
		case "acquire-return" :
			return decodeResourceAcquireReturn (_data, _attachment)
		default :
			return nil, fmt.Errorf ("unknown resources message action `%s`", _action)
	}
	panic ("fallthrough")
}

func decodeResourceAcquire (_data PacketData, _attachment PacketAttachment) (*ResourceAcquire, error) {
	panic ("not-implemented")
}

func decodeResourceAcquireReturn (_data PacketData, _attachment PacketAttachment) (*ResourceAcquireReturn, error) {
	var _ok bool
	if _ok_1, _error := extractBool (_data, "ok"); _error != nil {
		return nil, _error
	} else {
		_ok = _ok_1
	}
	var _descriptors []ResourceDescriptor
	var _merror interface{}
	if _ok {
		if _descriptors_1, _error := extractMap (_data, "descriptors"); _error != nil {
			return nil, _error
		} else if _descriptors_2, _error := decodeResourceDescriptors (_descriptors_1); _error != nil {
			return nil, _error
		} else {
			_descriptors = _descriptors_2
		}
	} else {
		if _merror_1, _error := extractObject (_data, "error"); _error != nil {
			return nil, _error
		} else {
			_merror = _merror_1
		}
	}
	var _correlation Correlation
	if _correlation_1, _error := extractCorrelation (_data, "correlation"); _error != nil {
		return nil, _error
	} else {
		_correlation = _correlation_1
	}
	return & ResourceAcquireReturn {
			Ok : _ok,
			Descriptors : _descriptors,
			Error : _merror,
			Correlation : _correlation,
	}, nil
}

func decodeResourceSpecifications (_data map[string]interface{}) ([]ResourceSpecification, error) {
	_specifications := make ([]ResourceSpecification, 0, len (_data))
	for _identifier, _specificationData := range _data {
		if _specification, _error := decodeResourceSpecification (_identifier, _specificationData); _error != nil {
			return nil, _error
		} else {
			_specifications = append (_specifications, _specification)
		}
	}
	return _specifications, nil
}

func decodeResourceSpecification (_identifier string, _data interface{}) (ResourceSpecification, error) {
	if _type, _ok := _data.(string); _ok {
		switch _type {
			case "socket:ipv4:tcp" :
				return & TcpSocketSpecification {
						Identifier : ResourceIdentifier (_identifier),
				}, nil
			default :
				return nil, fmt.Errorf ("unexpected resource specification `%#v`", _data)
		}
	} else {
		return nil, fmt.Errorf ("unexpected resource specification `%#v`", _data)
	}
	panic ("fallthrough")
}

func decodeResourceDescriptors (_data map[string]interface{}) ([]ResourceDescriptor, error) {
	_descriptors := make ([]ResourceDescriptor, 0, len (_data))
	for _identifier, _descriptorData := range _data {
		if _descriptor, _error := decodeResourceDescriptor (_identifier, _descriptorData); _error != nil {
			return nil, _error
		} else {
			_descriptors = append (_descriptors, _descriptor)
		}
	}
	return _descriptors, nil
}

func decodeResourceDescriptor (_identifier string, _data_0 interface{}) (ResourceDescriptor, error) {
	if _data, _ok := _data_0.(map[string]interface{}); _ok {
		var _type string
		if _type_1, _error := extractString (_data, "type"); _error != nil {
			return nil, _error
		} else {
			_type = _type_1
		}
		switch _type {
			case "socket:ipv4:tcp" :
				var _ip net.IP
				if _ip_1, _error := extractString (_data, "ip"); _error != nil {
					return nil, _error
				} else if _ip_2 := net.ParseIP (_ip_1); _ip_2 == nil {
					return nil, fmt.Errorf ("invalid tcp socket resource descriptor IP `%s`", _ip_1)
				} else {
					_ip = _ip_2
				}
				var _port uint16
				if _port_1, _error := extractNumber (_data, "port"); _error != nil {
					return nil, _error
				} else {
					// FIXME: Conversion issue!
					_port = uint16 (_port_1)
				}
				var _fqdn string
				if _fqdn_1, _error := extractString (_data, "fqdn"); _error != nil {
					return nil, _error
				} else {
					_fqdn = _fqdn_1
				}
				return & TcpSocketDescriptor {
						Identifier : ResourceIdentifier (_identifier),
						Ip : _ip,
						Port : _port,
						Fqdn : _fqdn,
				}, nil
			default :
				return nil, fmt.Errorf ("unexpected resource descriptor `%#v`", _data_0)
		}
	} else {
		return nil, fmt.Errorf ("unexpected resource descriptor `%#v`", _data_0)
	}
}


func decodeTranscript (_data PacketData, _attachment PacketAttachment) (Message, error) {
	var _action string
	if _action_1, _error := extractString (_data, "action"); _error != nil {
		return nil, _error
	} else {
		_action = _action_1
	}
	switch _action {
		case "push" :
			return decodeTranscriptPush (_data, _attachment)
		default :
			return nil, fmt.Errorf ("unknown transcript message action `%s`", _action)
	}
	panic ("fallthrough")
}

func decodeTranscriptPush (_data PacketData, _attachment PacketAttachment) (*TranscriptPush, error) {
	// FIXME: Enforce extraneous data entries
	return & TranscriptPush {
		Data : Attachment (_attachment),
	}, nil
}


func containsKey (_data_0 PacketData, _key string) (bool) {
	_data := map[string]interface{} (_data_0)
	_, _exists := _data[_key]
	return _exists
}

func extractString (_data_0 PacketData, _key string) (string, error) {
	_data := map[string]interface{} (_data_0)
	if _value_1, _exists := _data[_key]; !_exists {
		return "", fmt.Errorf ("missing message data `%s`", _key)
	} else if _value_2, _ok := _value_1.(string); !_ok {
		return "", fmt.Errorf ("invalid message data `%s`: `%#v`", _key, _value_1)
	} else {
		return _value_2, nil
	}
	panic ("fallthrough")
}

func extractNumber (_data_0 PacketData, _key string) (float64, error) {
	_data := map[string]interface{} (_data_0)
	if _value_1, _exists := _data[_key]; !_exists {
		return 0, fmt.Errorf ("missing message data `%s`", _key)
	} else if _value_2, _ok := _value_1.(float64); !_ok {
		return 0, fmt.Errorf ("invalid message data `%s`: `%#v`", _key, _value_1)
	} else {
		return _value_2, nil
	}
	panic ("fallthrough")
}

func extractBool (_data_0 PacketData, _key string) (bool, error) {
	_data := map[string]interface{} (_data_0)
	if _value_1, _exists := _data[_key]; !_exists {
		return false, fmt.Errorf ("missing message data `%s`", _key)
	} else if _value_2, _ok := _value_1.(bool); !_ok {
		return false, fmt.Errorf ("invalid message data `%s`: `%#v`", _key, _value_1)
	} else {
		return _value_2, nil
	}
	panic ("fallthrough")
}

func extractObject (_data_0 PacketData, _key string) (interface{}, error) {
	_data := map[string]interface{} (_data_0)
	if _value_1, _exists := _data[_key]; !_exists {
		return "", fmt.Errorf ("missing message data `%s`", _key)
	} else {
		return _value_1, nil
	}
	panic ("fallthrough")
}

func extractArray (_data_0 PacketData, _key string) ([]interface{}, error) {
	_data := map[string]interface{} (_data_0)
	if _value_1, _exists := _data[_key]; !_exists {
		return nil, fmt.Errorf ("missing message data `%s`", _key)
	} else if _value_2, _ok := _value_1.([]interface{}); !_ok {
		return nil, fmt.Errorf ("invalid message data `%s`: `%#v`", _key, _value_1)
	} else {
		return _value_2, nil
	}
	panic ("fallthrough")
}

func extractMap (_data_0 PacketData, _key string) (map[string]interface{}, error) {
	_data := map[string]interface{} (_data_0)
	if _value_1, _exists := _data[_key]; !_exists {
		return nil, fmt.Errorf ("missing message data `%s`", _key)
	} else if _value_2, _ok := _value_1.(map[string]interface{}); !_ok {
		return nil, fmt.Errorf ("invalid message data `%s`: `%#v`", _key, _value_1)
	} else {
		return _value_2, nil
	}
	panic ("fallthrough")
}


func extractComponentIdentifier (_data PacketData, _key string) (ComponentIdentifier, error) {
	if _value, _error := extractString (_data, _key); _error != nil {
		return "", _error
	} else {
		// FIXME: Validate syntax!
		return ComponentIdentifier (_value), nil
	}
}

func extractComponentOperation (_data PacketData, _key string) (ComponentOperation, error) {
	if _value, _error := extractString (_data, _key); _error != nil {
		return "", _error
	} else {
		// FIXME: Validate syntax!
		return ComponentOperation (_value), nil
	}
}

func extractCorrelation (_data PacketData, _key string) (Correlation, error) {
	if _value, _error := extractString (_data, _key); _error != nil {
		return "", _error
	} else {
		// FIXME: Validate syntax!
		return Correlation (_value), nil
	}
}

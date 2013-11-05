

package messages


import "net"

import "crypto/rand"
import "encoding/hex"


func NewCorrelation () (Correlation) {
	// FIXME: Validate syntax!
	_buffer := make ([]byte, 16)
	if _, _error := rand.Read (_buffer); _error != nil {
		panic (_error)
	}
	_hex := hex.EncodeToString (_buffer)
	_correlation := Correlation (_hex)
	return _correlation
}


func ComponentCallInvoke (_component ComponentIdentifier, _operation ComponentOperation, _inputs interface{}, _attachment Attachment) (*ComponentCall, Correlation) {
	_correlation := NewCorrelation ()
	_call := & ComponentCall {
			Component : _component,
			Operation : _operation,
			Inputs : _inputs,
			Correlation : _correlation,
			Attachment : _attachment,
	}
	return _call, _correlation
}

func ComponentCallSucceeded (_correlation Correlation, _outputs interface{}, _attachment Attachment) (*ComponentCallReturn) {
	_return := & ComponentCallReturn {
			Ok : true,
			Outputs : _outputs,
			Correlation : _correlation,
			Attachment : _attachment,
	}
	return _return
}

func ComponentCallFailed (_correlation Correlation, _error interface{}, _attachment Attachment) (*ComponentCallReturn) {
	_return := & ComponentCallReturn {
			Ok : false,
			Error : _error,
			Correlation : _correlation,
			Attachment : _attachment,
	}
	return _return
}

func ComponentCallCompleted (_correlation Correlation, _ok bool, _outputsOrError interface{}, _attachment Attachment) (*ComponentCallReturn) {
	if _ok {
		return ComponentCallSucceeded (_correlation, _outputsOrError, _attachment)
	} else {
		return ComponentCallFailed (_correlation, _outputsOrError, _attachment)
	}
}

func (_call *ComponentCall) Succeeded (_outputs interface{}, _attachment Attachment) (*ComponentCallReturn) {
	return ComponentCallSucceeded (_call.Correlation, _outputs, _attachment)
}

func (_call *ComponentCall) ReturnFailure (_error interface{}, _attachment Attachment) (*ComponentCallReturn) {
	return ComponentCallFailed (_call.Correlation, _error, _attachment)
}

func (_call *ComponentCall) Return (_ok bool, _outputsOrError interface{}, _attachment Attachment) (*ComponentCallReturn) {
	return ComponentCallCompleted (_call.Correlation, _ok, _outputsOrError, _attachment)
}


func ComponentCastInvoke (_component ComponentIdentifier, _operation ComponentOperation, _inputs interface{}, _attachment Attachment) (*ComponentCast) {
	_cast := & ComponentCast {
			Component : _component,
			Operation : _operation,
			Inputs : _inputs,
			Attachment : _attachment,
	}
	return _cast
}


func ComponentRegisterInvoke (_group ComponentGroup) (*ComponentRegister, Correlation) {
	_correlation := NewCorrelation ()
	_call := & ComponentRegister {
			Component : NilComponentIdentifier,
			Group : _group,
			Correlation : _correlation,
	}
	return _call, _correlation
}

func ComponentRegisterSucceeded (_correlation Correlation) (*ComponentRegisterReturn) {
	_return := & ComponentRegisterReturn {
			Ok : true,
			Correlation : _correlation,
	}
	return _return
}

func ComponentRegisterFailed (_correlation Correlation, _error interface{}) (*ComponentRegisterReturn) {
	_return := & ComponentRegisterReturn {
			Ok : false,
			Error : _error,
			Correlation : _correlation,
	}
	return _return
}

func (_call *ComponentRegister) Succeeded () (*ComponentRegisterReturn) {
	return ComponentRegisterSucceeded (_call.Correlation)
}

func (_call *ComponentRegister) Failed (_error interface{}) (*ComponentRegisterReturn) {
	return ComponentRegisterFailed (_call.Correlation, _error)
}


func ResourceAcquireInvoke (_specification ResourceSpecification) (*ResourceAcquire, Correlation) {
	_correlation := NewCorrelation ()
	_call := & ResourceAcquire {
			Specifications : []ResourceSpecification { _specification },
			Correlation : _correlation,
	}
	return _call, _correlation
}

func TcpSocketAcquireInvoke (_identifier ResourceIdentifier) (*ResourceAcquire, Correlation) {
	_specification := & TcpSocketSpecification {
			Identifier : ResourceIdentifier (_identifier),
	}
	return ResourceAcquireInvoke (_specification)
}

func ResourceAcquireSucceeded (_correlation Correlation, _descriptor ResourceDescriptor) (*ResourceAcquireReturn) {
	_return := & ResourceAcquireReturn {
			Ok : true,
			Descriptors : []ResourceDescriptor { _descriptor },
			Correlation : _correlation,
	}
	return _return
}

func TcpSocketAcquireSucceeded (_correlation Correlation, _identifier ResourceIdentifier, _ip net.IP, _port uint16, _fqdn string) (*ResourceAcquireReturn) {
	_descriptor := & TcpSocketDescriptor {
			Identifier : _identifier,
			Ip : _ip,
			Port : _port,
			Fqdn : _fqdn,
	}
	return ResourceAcquireSucceeded (_correlation, _descriptor)
}

func ResourceAcquireFailed (_correlation Correlation, _error interface{}) (*ResourceAcquireReturn) {
	_return := & ResourceAcquireReturn {
			Ok : false,
			Error : _error,
			Correlation : _correlation,
	}
	return _return
}

func (_call *ResourceAcquire) Failed (_error interface{}) (*ResourceAcquireReturn) {
	return ResourceAcquireFailed (_call.Correlation, _error)
}


func TranscriptPushInvoke (_data Attachment) (*TranscriptPush) {
	_call := & TranscriptPush {
			Data : _data,
	}
	return _call
}



package backend


import "fmt"

import . "mch/lib/messages"


func DispatchComponentInvoke (_message_0 Message, _callbacks ComponentInvokeCallbacks) (error) {
	switch _message := _message_0.(type) {
		case *ComponentCall :
			if _message.Component != NilComponentIdentifier {
				return fmt.Errorf ("unexpected component identifier `%s` (none expected)", _message.Component)
			}
			return _callbacks.ComponentCallInvoked (_message.Operation, _message.Inputs, _message.Correlation, _message.Attachment)
		case *ComponentCast :
			if _message.Component != NilComponentIdentifier {
				return fmt.Errorf ("unexpected component identifier `%s` (none expected)", _message.Component)
			}
			return _callbacks.ComponentCastInvoked (_message.Operation, _message.Inputs, _message.Attachment)
		default :
			return NotDispatchedError
	}
	panic ("fallthrough")
}


func DispatchComponentReturn (_message_0 Message, _callbacks ComponentReturnCallbacks) (error) {
	switch _message := _message_0.(type) {
		case *ComponentCallReturn :
			if _message.Ok {
				return _callbacks.ComponentCallSucceeded (_message.Correlation, _message.Outputs, _message.Attachment)
			} else {
				return _callbacks.ComponentCallFailed (_message.Correlation, _message.Error, _message.Attachment)
			}
		case *ComponentRegisterReturn :
			if _message.Ok {
				return _callbacks.ComponentRegisterSucceeded (_message.Correlation)
			} else {
				return _callbacks.ComponentRegisterFailed (_message.Correlation, _message.Error)
			}
		default :
			return NotDispatchedError
	}
	panic ("fallthrough")
}


func DispatchResourceReturn (_message_0 Message, _callbacks ResourceReturnCallbacks) (error) {
	switch _message := _message_0.(type) {
		case *ResourceAcquireReturn :
			if _message.Ok {
				var _error error
				for _, _descriptor := range _message.Descriptors {
					if _error_1 := _callbacks.ResourceAcquireSucceeded (_message.Correlation, _descriptor); _error_1 != nil {
						if _error == nil {
							_error = _error_1
						} else {
							// FIXME: Handle this error!
						}
					}
				}
				return _error
			} else {
				return _callbacks.ResourceAcquireFailed (_message.Correlation, _message.Error)
			}
		default :
			return NotDispatchedError
	}
	panic ("fallthrough")
}


func DispatchMessage (_message Message, _callbacks MessageCallbacks) (error) {
	if _error := DispatchComponentInvoke (_message, _callbacks); _error != NotDispatchedError {
		return _error
	}
	if _error := DispatchComponentReturn (_message, _callbacks); _error != NotDispatchedError {
		return _error
	}
	if _error := DispatchResourceReturn (_message, _callbacks); _error != NotDispatchedError {
		return _error
	}
	return NotDispatchedError
}


var NotDispatchedError = fmt.Errorf ("message-not-dispatched")

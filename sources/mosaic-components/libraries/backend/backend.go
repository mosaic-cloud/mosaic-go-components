

package backend


import "fmt"
import "time"

import "mosaic-components/libraries/channels"
import "vgl/transcript"

import . "mosaic-components/libraries/messages"


type backend struct {
	state backendState
	callbacks Callbacks
	controllerIsolates chan func () ()
	callbacksIsolates chan func () ()
	pendingCompletions map[Correlation]chan Message
	channel channels.Controller
	transcript transcript.Transcript
}

type backendController backend
type backendMessageCallbacks backend
type backendChannelCallbacks backend

type backendState uint
const (
	invalidBackendStateMin backendState = iota
	Initializing
	Active
	Terminating
	Terminated
	invalidBackendStateMax
)


const isolateChannelBuffer = 128


func Create (_callbacks Callbacks) (*backend, channels.Callbacks, error) {
	
	_backend := & backend {
			state : Initializing,
			callbacks : _callbacks,
			controllerIsolates : make (chan func () (), isolateChannelBuffer),
			callbacksIsolates : make (chan func () (), isolateChannelBuffer),
			pendingCompletions : make (map[Correlation]chan Message),
			channel : nil,
	}
	
	_backend.transcript = transcript.NewTranscript (_backend, packageTranscript)
	_backend.transcript.TraceDebugging ("creating the backend...")
	
	go _backend.executeIsolateLoop (_backend.controllerIsolates)
	go _backend.executeIsolateLoop (_backend.callbacksIsolates)
	
	return _backend, (*backendChannelCallbacks) (_backend), nil
}

// NOTE: non-isolated
func (_backend *backend) Terminate () (error) {
	_backend.controllerIsolates <- func () () {
		_backend.initiateTerminate (nil)
	}
	return nil
}


// NOTE: non-isolated
func (_backend *backend) WaitTerminated () (error) {
	// FIXME: Implement this properly!
	for {
		if _backend.state != Terminated {
			time.Sleep (1 * time.Second)
			continue
		} else {
			break
		}
	}
	return nil
}


// NOTE: non-isolated
func (_backend_0 *backendController) ComponentCallSync (_component ComponentIdentifier, _operation ComponentOperation, _inputs interface{}, _attachment Attachment) (interface{}, Attachment, error) {
	_backend := (*backend) (_backend_0)
	_invoke, _correlation := ComponentCallInvoke (_component, _operation, _inputs, _attachment)
	if _return_1, _error := _backend.handleSync (_invoke, _correlation); _error != nil {
		return nil, nil, _error
	} else if _return, _typeOk := _return_1.(*ComponentCallReturn); !_typeOk {
		return nil, nil, fmt.Errorf ("unexpected return message `%#v`", _return_1)
	} else {
		if _return.Ok {
			return _return.Outputs, _return.Attachment, nil
		} else {
			return nil, _return.Attachment, fmt.Errorf ("call failed `%#v`", _return.Error)
		}
	}
	panic ("fallthrough")
}

// NOTE: non-isolated
func (_backend_0 *backendController) ComponentRegisterSync (_group ComponentGroup) (error) {
	_backend := (*backend) (_backend_0)
	_invoke, _correlation := ComponentRegisterInvoke (_group)
	if _return_1, _error := _backend.handleSync (_invoke, _correlation); _error != nil {
		return _error
	} else if _return, _typeOk := _return_1.(*ComponentRegisterReturn); !_typeOk {
		return fmt.Errorf ("unexpected return message `%#v`", _return_1)
	} else {
		if _return.Ok {
			return nil
		} else {
			return fmt.Errorf ("call failed `%#v`", _return.Error)
		}
	}
	panic ("fallthrough")
}

// NOTE: non-isolated
func (_backend_0 *backendController) ResourceAcquireSync (_specification ResourceSpecification) (ResourceDescriptor, error) {
	_backend := (*backend) (_backend_0)
	_invoke, _correlation := ResourceAcquireInvoke (_specification)
	if _return_1, _error := _backend.handleSync (_invoke, _correlation); _error != nil {
		return nil, _error
	} else if _return, _typeOk := _return_1.(*ResourceAcquireReturn); !_typeOk {
		return nil, fmt.Errorf ("unexpected return message `%#v`", _return_1)
	} else {
		if _return.Ok {
			if len (_return.Descriptors) == 1 {
				return _return.Descriptors[0], nil
			} else {
				return nil, fmt.Errorf ("unexpected descriptors `%#v`", _return.Descriptors)
			}
		} else {
			return nil, fmt.Errorf ("call failed `%#v`", _return.Error)
		}
	}
	panic ("fallthrough")
}


// NOTE: non-isolated
func (_backend *backend) handleSync (_invoke Message, _correlation Correlation) (Message, error) {
	_completion := make (chan Message, 1)
	defer close (_completion)
	if _error := _backend.handleOutboundMessage1 (_invoke, _correlation, _completion); _error != nil {
		return nil, _error
	}
	_return := <- _completion
	delete (_backend.pendingCompletions, _correlation)
	if _return == nil {
		return nil, fmt.Errorf ("sync-aborted")
	}
	return _return, nil
}


// NOTE: non-isolated
func (_backend_0 *backendController) ComponentCallInvoke (_component ComponentIdentifier, _operation ComponentOperation, _inputs interface{}, _attachment Attachment) (Correlation, error) {
	_backend := (*backend) (_backend_0)
	_message, _correlation := ComponentCallInvoke (_component, _operation, _inputs, _attachment)
	if _error := _backend.handleOutboundMessage (_message); _error != nil {
		return NilCorrelation, _error
	}
	return _correlation, nil
}

// NOTE: non-isolated
func (_backend_0 *backendController) ComponentCallSucceeded (_correlation Correlation, _outputs interface{}, _attachment Attachment) (error) {
	_backend := (*backend) (_backend_0)
	_message := ComponentCallSucceeded (_correlation, _outputs, _attachment)
	return _backend.handleOutboundMessage (_message)
}

// NOTE: non-isolated
func (_backend_0 *backendController) ComponentCallFailed (_correlation Correlation, _error interface{}, _attachment Attachment) (error) {
	_backend := (*backend) (_backend_0)
	_message := ComponentCallFailed (_correlation, _error, _attachment)
	return _backend.handleOutboundMessage (_message)
}

// NOTE: non-isolated
func (_backend_0 *backendController) ComponentCastInvoke (_component ComponentIdentifier, _operation ComponentOperation, _inputs interface{}, _attachment Attachment) (error) {
	_backend := (*backend) (_backend_0)
	_message := ComponentCastInvoke (_component, _operation, _inputs, _attachment)
	return _backend.handleOutboundMessage (_message)
}

// NOTE: non-isolated
func (_backend_0 *backendController) ComponentRegisterInvoke (_group ComponentGroup) (Correlation, error) {
	_backend := (*backend) (_backend_0)
	_message, _correlation := ComponentRegisterInvoke (_group)
	if _error := _backend.handleOutboundMessage (_message); _error != nil {
		return NilCorrelation, _error
	}
	return _correlation, nil
}

// NOTE: non-isolated
func (_backend_0 *backendController) ResourceAcquireInvoke (_specification ResourceSpecification) (Correlation, error) {
	_backend := (*backend) (_backend_0)
	_message, _correlation := ResourceAcquireInvoke (_specification)
	if _error := _backend.handleOutboundMessage (_message); _error != nil {
		return NilCorrelation, _error
	}
	return _correlation, nil
}


// NOTE: non-isolated
func (_backend_0 *backendController) TranscriptPushInvoke (_data Attachment) (error) {
	_backend := (*backend) (_backend_0)
	_message := TranscriptPushInvoke (_data)
	return _backend.handleOutboundMessage (_message)
}

// NOTE: non-isolated
func (_backend_0 *backendController) Terminate () (error) {
	_backend := (*backend) (_backend_0)
	_backend.controllerIsolates <- func () () {
		_backend.initiateTerminate (nil)
	}
	return nil
}


// NOTE: isolated
func (_backend_0 *backendMessageCallbacks) ComponentCallInvoked (_operation ComponentOperation, _inputs interface{}, _correlation Correlation, _attachment Attachment) (error) {
	_backend := (*backend) (_backend_0)
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.ComponentCallInvoked (_operation, _inputs, _correlation, _attachment); _error != nil {
			_backend.handleCallbacksError (_error)
			return
		}
	}
	return nil
}

// NOTE: isolated
func (_backend_0 *backendMessageCallbacks) ComponentCastInvoked (_operation ComponentOperation, _inputs interface{}, _attachment Attachment) (error) {
	_backend := (*backend) (_backend_0)
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.ComponentCastInvoked (_operation, _inputs, _attachment); _error != nil {
			_backend.handleCallbacksError (_error)
			return
		}
	}
	return nil
}

// NOTE: isolated
func (_backend_0 *backendMessageCallbacks) ComponentCallSucceeded (_correlation Correlation, _outputs interface{}, _attachment Attachment) (error) {
	_backend := (*backend) (_backend_0)
	if _completion, _exists := _backend.pendingCompletions[_correlation]; _exists {
		_completion <- ComponentCallSucceeded (_correlation, _outputs, _attachment)
		return nil
	}
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.ComponentCallSucceeded (_correlation, _outputs, _attachment); _error != nil {
			_backend.handleCallbacksError (_error)
			return
		}
	}
	return nil
}

// NOTE: isolated
func (_backend_0 *backendMessageCallbacks) ComponentCallFailed (_correlation Correlation, _error interface{}, _attachment Attachment) (error) {
	_backend := (*backend) (_backend_0)
	if _completion, _exists := _backend.pendingCompletions[_correlation]; _exists {
		_completion <- ComponentCallFailed (_correlation, _error, _attachment)
		return nil
	}
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.ComponentCallFailed (_correlation, _error, _attachment); _error != nil {
			_backend.handleCallbacksError (_error)
			return
		}
	}
	return nil
}

// NOTE: isolated
func (_backend_0 *backendMessageCallbacks) ComponentRegisterSucceeded (_correlation Correlation) (error) {
	_backend := (*backend) (_backend_0)
	if _completion, _exists := _backend.pendingCompletions[_correlation]; _exists {
		_completion <- ComponentRegisterSucceeded (_correlation)
		return nil
	}
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.ComponentRegisterSucceeded (_correlation); _error != nil {
			_backend.handleCallbacksError (_error)
			return
		}
	}
	return nil
}

// NOTE: isolated
func (_backend_0 *backendMessageCallbacks) ComponentRegisterFailed (_correlation Correlation, _error interface{}) (error) {
	_backend := (*backend) (_backend_0)
	if _completion, _exists := _backend.pendingCompletions[_correlation]; _exists {
		_completion <- ComponentRegisterFailed (_correlation, _error)
		return nil
	}
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.ComponentRegisterFailed (_correlation, _error); _error != nil {
			_backend.handleCallbacksError (_error)
			return
		}
	}
	return nil
}

// NOTE: isolated
func (_backend_0 *backendMessageCallbacks) ResourceAcquireSucceeded (_correlation Correlation, _descriptor ResourceDescriptor) (error) {
	_backend := (*backend) (_backend_0)
	if _completion, _exists := _backend.pendingCompletions[_correlation]; _exists {
		_completion <- ResourceAcquireSucceeded (_correlation, _descriptor)
		return nil
	}
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.ResourceAcquireSucceeded (_correlation, _descriptor); _error != nil {
			_backend.handleCallbacksError (_error)
			return
		}
	}
	return nil
}

// NOTE: isolated
func (_backend_0 *backendMessageCallbacks) ResourceAcquireFailed (_correlation Correlation, _error interface{}) (error) {
	_backend := (*backend) (_backend_0)
	if _completion, _exists := _backend.pendingCompletions[_correlation]; _exists {
		_completion <- ResourceAcquireFailed (_correlation, _error)
		return nil
	}
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.ResourceAcquireFailed (_correlation, _error); _error != nil {
			_backend.handleCallbacksError (_error)
			return
		}
	}
	return nil
}


// NOTE: non-isolated
func (_backend *backend) handleInboundMessage (_message Message) () {
	_backend.controllerIsolates <- func () () {
		if _backend.state != Active {
			panic ("invalid-state")
		}
		_backend.transcript.TraceDebugging ("dispatching message `%#v`...", _message)
		if _error := DispatchMessage (_message, (*backendMessageCallbacks) (_backend)); _error != nil {
			panic (_error)
		}
	}
}


// NOTE: non-isolated
func (_backend *backend) handleOutboundMessage (_message Message) (error) {
	return _backend.handleOutboundMessage1 (_message, NilCorrelation, nil)
}

// NOTE: non-isolated
func (_backend *backend) handleOutboundMessage1 (_message Message, _pendingCorrelation Correlation, _pendingCompletion chan Message) (error) {
	var _packet *channels.Packet
	if _packet_1, _error := Encode (_message); _error != nil {
		return _error
	} else {
		_packet = _packet_1
	}
	_completion := make (chan error, 1)
	defer close (_completion)
	_backend.controllerIsolates <- func () () {
		if _backend.state != Active {
			_completion <- fmt.Errorf ("invalid-state")
			return
		}
		if _error := _backend.channel.Push (_packet); _error != nil {
			_completion <- _error
			return
		}
		if _pendingCorrelation != NilCorrelation {
			_backend.pendingCompletions[_pendingCorrelation] = _pendingCompletion
		}
		_completion <- nil
	}
	return <- _completion
}


// NOTE: isolated
func (_backend *backend) initialize (_channel channels.Controller) () {
	if _backend.state != Initializing {
		panic ("illegal-state")
	}
	_backend.state = Active
	_backend.channel = _channel
	_backend.callbacksIsolates <- func () () {
		if _error := _backend.callbacks.Initialized ((*backendController) (_backend)); _error != nil {
			_backend.handleCallbacksError (_error)
		}
	}
}

// NOTE: isolated
func (_backend *backend) initiateTerminate (_error error) () {
	if _backend.state == Terminating {
		// FIXME: Better handle this case?
		return
	} else if _backend.state != Active {
		panic ("illegal-state")
	}
	
	for _, _pendingCompletion := range _backend.pendingCompletions {
		_pendingCompletion <- nil
	}
	
	_backend.state = Terminating
	_backend.channel.Close (channels.InboundFlow)
	_backend.channel.Close (channels.OutboundFlow)
}

// NOTE: isolated
func (_backend *backend) concludeTerminate (_error error) () {
	if _backend.state == Terminated {
		panic ("illegal-state")
	} else if _backend.state != Terminating {
		panic ("illegal-state")
	}
	_backend.state = Terminated
	_backend.callbacksIsolates <- func () () {
		_backend.callbacks.Terminated (_error)
		close (_backend.callbacksIsolates)
	}
	close (_backend.controllerIsolates)
}



// NOTE: non-isolated
func (_backend_0 *backendChannelCallbacks) Initialized (_channel channels.Controller) (error) {
	_backend := (*backend) (_backend_0)
	_backend.controllerIsolates <- func () () {
		_backend.initialize (_channel)
	}
	return nil
}

// NOTE: non-isolated
func (_backend_0 *backendChannelCallbacks) Pushed (_packet *channels.Packet) (error) {
	_backend := (*backend) (_backend_0)
	var _message Message
	if _message_1, _error := Decode (_packet); _error != nil {
		panic (_error)
	} else {
		_message = _message_1
	}
	_backend.handleInboundMessage (_message)
	return nil
}

// NOTE: non-isolated
func (_backend_0 *backendChannelCallbacks) Closed (_flow channels.Flow, _error error) (error) {
	_backend := (*backend) (_backend_0)
	_backend.controllerIsolates <- func () () {
		_backend.initiateTerminate (_error)
	}
	return nil
}

// NOTE: non-isolated
func (_backend_0 *backendChannelCallbacks) Terminated (_error error) (error) {
	_backend := (*backend) (_backend_0)
	_backend.controllerIsolates <- func () () {
		_backend.concludeTerminate (_error)
	}
	return nil
}


func (_backend *backend) executeIsolateLoop (_isolates chan func () ()) () {
	for {
		_isolate := <- _isolates
		if _isolate == nil {
			break
		}
		_isolate ()
	}
}


func (_backend *backend) handleCallbacksError (_error error) () {
	panic (_error)
}

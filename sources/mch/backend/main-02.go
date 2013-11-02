

package main


import "os"

import "mch/lib/backend"
import "mch/lib/channels"
import "vgl/transcript"

import . "mch/lib/messages"


func main () () {
	
	_callbacks := & callbacks {
			backend : nil,
			terminated : make (chan error),
			transcript : nil,
	}
	_callbacks.transcript = transcript.NewTranscript (_callbacks, packageTranscript)
	_transcript := packageTranscript
	var _error error
	
	_transcript.TraceInformation ("creating the backend...")
	var _backend backend.Backend
	var _backendChannelCallbacks channels.Callbacks
	if _backend, _backendChannelCallbacks, _error = backend.Create (_callbacks); _error != nil {
		panic (_error)
	}
	
	_transcript.TraceInformation ("creating the channel...")
	var _channel channels.Channel
	if false {
		_inboundStream := os.Stdin
		_outboundStream := os.Stdout
		if _channel, _error = channels.Create (_backendChannelCallbacks, _inboundStream, _outboundStream, nil); _error != nil {
			panic (_error)
		}
	} else {
		if _channel, _error = channels.CreateAndDial (_backendChannelCallbacks, "tcp", "127.0.0.1:24704"); _error != nil {
			panic (_error)
		}
	}
	
	if _error = <- _callbacks.terminated; _error != nil {
		panic (_error)
	}
	
	_backend.Terminate ()
	_channel.Terminate ()
	
	_transcript.TraceInformation ("done.")
}


type callbacks struct {
	backend backend.Controller
	terminated chan error
	transcript transcript.Transcript
}


func (_callbacks *callbacks) Initialized (_backend backend.Controller) (error) {
	
	_callbacks.transcript.TraceInformation ("initialized")
	_callbacks.backend = _backend
	
	_callbacks.transcript.TraceInformation ("acquiring TCP endpoint")
	if _ip, _port, _fqdn, _error := backend.TcpSocketAcquireSync (_callbacks.backend, ResourceIdentifier ("tests")); _error != nil {
		panic (_error)
	} else {
		_callbacks.transcript.TraceInformation ("aquired `%s:%d` / `%s`.", _ip, _port, _fqdn)
	}
	
	_callbacks.transcript.TraceInformation ("registering in group...")
	if _error := _callbacks.backend.ComponentRegisterSync (testsGroup); _error != nil {
		panic (_error)
	} else {
		_callbacks.transcript.TraceInformation ("registered.")
	}
	
	return nil
}

func (_callbacks *callbacks) Terminated (_error error) (error) {
	_callbacks.transcript.TraceInformation ("terminated")
	_callbacks.terminated <- _error
	return nil
}

func (_callbacks *callbacks) ComponentCallInvoked (_operation ComponentOperation, _inputs interface{}, _correlation Correlation, _attachment Attachment) (error) {
	_callbacks.transcript.TraceInformation ("call invoked for operation `%s` with inputs `%#v`...", _operation, _inputs)
	return _callbacks.backend.ComponentCallSucceeded (_correlation, _inputs, nil)
}

func (_callbacks *callbacks) ComponentCastInvoked (_operation ComponentOperation, _inputs interface{}, _attachment Attachment) (error) {
	panic ("unexpected")
}

func (_callbacks *callbacks) ComponentCallSucceeded (_correlation Correlation, _outputs interface{}, _attachment Attachment) (error) {
	panic ("unexpected")
}

func (_callbacks *callbacks) ComponentCallFailed (_correlation Correlation, _error interface{}, _attachment Attachment) (error) {
	panic ("unexpected")
}

func (_callbacks *callbacks) ComponentRegisterSucceeded (_correlation Correlation) (error) {
	panic ("unexpected")
}

func (_callbacks *callbacks) ComponentRegisterFailed (_correlation Correlation, _error interface{}) (error) {
	panic ("unexpected")
}

func (_callbacks *callbacks) ResourceAcquireSucceeded (_correlation Correlation, _descriptor ResourceDescriptor) (error) {
	panic ("unexpected")
}

func (_callbacks *callbacks) ResourceAcquireFailed (_correlation Correlation, _error interface{}) (error) {
	panic ("unexpected")
}


var packageTranscript = transcript.NewPackageTranscript ()
var testsGroup = ComponentGroup ("85aa675f0f3af10789e2ef4bf07665217fd91bc6")

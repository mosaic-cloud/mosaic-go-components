

package main


import "os"

import "mosaic-components/libraries/backend"
import "mosaic-components/libraries/channels"
import "vgl/transcript"

import . "mosaic-components/libraries/messages"


func main () () {
	
	_callbacks := & callbacks {
			backend : nil,
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
	if true {
		_inboundStream := os.Stdin
		_outboundStream := os.Stdout
		if _, _error = channels.Create (_backendChannelCallbacks, _inboundStream, _outboundStream, nil); _error != nil {
			panic (_error)
		}
	} else {
		if _, _error = channels.CreateAndDial (_backendChannelCallbacks, "tcp", "127.0.0.1:24704"); _error != nil {
			panic (_error)
		}
	}
	
	_backend.WaitTerminated ()
	
	_transcript.TraceInformation ("done.")
}


type callbacks struct {
	backend backend.Controller
	transcript transcript.Transcript
}


func (_callbacks *callbacks) Initialized (_backend backend.Controller) (error) {
	_callbacks.transcript.TraceInformation ("initialized.")
	_callbacks.backend = _backend
	return nil
}

func (_callbacks *callbacks) Terminated (_error error) (error) {
	_callbacks.transcript.TraceInformation ("terminated.")
	return nil
}

func (_callbacks *callbacks) ComponentCallInvoked (_operation ComponentOperation, _inputs interface{}, _correlation Correlation, _attachment Attachment) (error) {
	panic ("unexpected")
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

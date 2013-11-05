

package channels


import "fmt"
import "io"
import "sync"

import "vgl/transcript"


type channel struct {
	
	callbacks Callbacks
	inboundStream io.Reader
	outboundStream io.Writer
	closer io.Closer
	
	controllerActive bool
	inboundActive bool
	outboundActive bool
	
	controllerIsolates chan func () ()
	callbacksIsolates chan func () ()
	inboundPackets chan *Packet
	outboundPackets chan *Packet
	inboundSignal chan bool
	outboundSignal chan bool
	terminateAcknowledgments sync.WaitGroup
	
	transcript transcript.Transcript
}


const packetChannelBuffer = 128
const isolateChannelBuffer = 16


func Create (_callbacks Callbacks, _inbound io.Reader, _outbound io.Writer, _closer io.Closer) (Channel, error) {
	
	_channel := & channel {
			
			callbacks : _callbacks,
			inboundStream : _inbound,
			outboundStream : _outbound,
			closer : _closer,
			
			controllerActive : true,
			inboundActive : true,
			outboundActive : true,
			
			controllerIsolates : make (chan func () (), isolateChannelBuffer),
			callbacksIsolates : make (chan func () (), isolateChannelBuffer),
			inboundPackets : make (chan *Packet, packetChannelBuffer),
			outboundPackets : make (chan *Packet, packetChannelBuffer),
			inboundSignal : make (chan bool, 1),
			outboundSignal : make (chan bool, 1),
			terminateAcknowledgments : sync.WaitGroup {},
	}
	
	_channel.transcript = transcript.NewTranscript (_channel, _packageTranscript)
	_channel.transcript.TraceDebugging ("creating the channel...")
	
	_channel.terminateAcknowledgments.Add (2)
	go _channel.executeControllerLoop ()
	go _channel.executeCallbacksLoop ()
	go _channel.executeInboundLoop ()
	go _channel.executeOutboundLoop ()
	
	return _channel, nil
}


func (_channel *channel) Push (_packet *Packet) (error) {
	_completion := make (chan error, 1)
	defer close (_completion)
	
	_channel.controllerIsolates <- func () () {
		if !_channel.outboundActive {
			_completion <- fmt.Errorf ("outbound channel flow is closed!")
			return
		}
		_channel.outboundPackets <- _packet
		_completion <- nil
	}
	
	return <- _completion
}


func (_channel *channel) Close (_flow Flow) (error) {
	_completion := make (chan error, 1)
	defer close (_completion)
	
	_channel.controllerIsolates <- func () () {
		switch _flow {
			
			case InboundFlow :
				_channel.transcript.TraceDebugging ("closing the channel inbound flow...")
				if !_channel.inboundActive {
					_completion <- fmt.Errorf ("inbound channel flow already closed!")
					break
				}
				_channel.inboundActive = false
				select {
					case _channel.inboundSignal <- true :
						// NOTE: <nop>
					default :
						// NOTE: <nop>
				}
				_completion <- nil
			
			case OutboundFlow :
				_channel.transcript.TraceDebugging ("closing the channel outbound flow...")
				if !_channel.outboundActive {
					_completion <- fmt.Errorf ("outbound channel flow already closed!")
					break
				}
				_channel.outboundActive = false
				select {
					case _channel.outboundSignal <- true :
						// NOTE: <nop>
					default :
						// NOTE: <nop>
				}
				_completion <- nil
			
			default :
				_completion <- fmt.Errorf ("invalid channel flow")
		}
	}
	
	return <- _completion
}


func (_channel *channel) Terminate () (error) {
	
	_completion := make (chan error, 1)
	defer close (_completion)
	
	_channel.controllerIsolates <- func () () {
		
		_channel.transcript.TraceDebugging ("terminating the channel...")
		if ! _channel.controllerActive {
			_completion <- fmt.Errorf ("channel is already terminated!")
			return
		}
		
		_channel.controllerActive = false
		_channel.inboundActive = false
		_channel.outboundActive = false
		
		select {
			case _channel.inboundSignal <- true :
				// NOTE: <nop>
			default :
				// NOTE: <nop>
		}
		select {
			case _channel.outboundSignal <- true :
				// NOTE: <nop>
			default :
				// NOTE: <nop>
		}
		
		_channel.transcript.TraceDebugging ("waiting for the channel background tasks (phase 1)...")
		_channel.terminateAcknowledgments.Wait ()
		
		var _error error
		
		if _channel.closer != nil {
			_channel.transcript.TraceDebugging ("invoking the channel closer...")
			if _closerError := _channel.closer.Close (); _closerError != nil {
				if _error == nil {
					_error = _closerError
				} else {
					// FIXME: Handle this error!
				}
			}
		}
		
		_channel.terminateAcknowledgments.Add (1)
		_channel.callbacksIsolates <- func () () {
			if _error := _channel.callbacks.Terminated (_error); _error != nil {
				_channel.handleCallbacksError (_error)
			}
		}
		
		_channel.transcript.TraceDebugging ("waiting for the channel background tasks (phase 2)...")
		_channel.terminateAcknowledgments.Wait ()
		
		_channel.transcript.TraceDebugging ("terminated the channel.")
		close (_channel.controllerIsolates)
		close (_channel.callbacksIsolates)
		
		_completion <- nil
	}
	
	return <- _completion
}


func (_channel *channel) executeControllerLoop () {
	_channel.transcript.TraceDebugging ("started the channel control background task.")
	for {
		_isolate, _ok := <- _channel.controllerIsolates
		if !_ok {
			_channel.controllerIsolates = nil
			break
		}
		_isolate ()
	}
	if _channel.controllerIsolates != nil {
		close (_channel.controllerIsolates)
		_channel.controllerIsolates = nil
	}
	_channel.transcript.TraceDebugging ("terminated the channel control background task.")
}


func (_channel *channel) executeCallbacksLoop () {
	_channel.transcript.TraceDebugging ("started the channel callbacks background task.")
	if _error := _channel.callbacks.Initialized (_channel); _error != nil {
		_channel.handleCallbacksError (_error)
	}
	for {
		select {
			case _packet, _ok := <- _channel.inboundPackets :
				if !_ok {
					_channel.inboundPackets = nil
					continue
				}
				if _error := _channel.callbacks.Pushed (_packet); _error != nil {
					_channel.handleCallbacksError (_error)
				}
			case _isolate, _ok := <- _channel.callbacksIsolates :
				if !_ok {
					_channel.callbacksIsolates = nil
					break
				}
				_isolate ()
		}
	}
	if _channel.inboundPackets != nil {
		close (_channel.inboundPackets)
		_channel.inboundPackets = nil
	}
	if _channel.callbacksIsolates != nil {
		close (_channel.callbacksIsolates)
		_channel.callbacksIsolates = nil
	}
	_channel.transcript.TraceDebugging ("terminated the channel callbacks background task.")
	_channel.terminateAcknowledgments.Done ()
}


func (_channel *channel) executeInboundLoop () {
	_channel.transcript.TraceDebugging ("started the channel inbound background task.")
	
	loop : for {
		select {
			case _ = <- _channel.inboundSignal :
				goto signalled
			default :
				// NOTE: <nop>
		}
		if _packet, _error := _channel.pullInboundPacket (); _error != nil {
			_channel.inboundActive = false
			if _error == io.EOF {
				_error = nil
			}
			_channel.callbacksIsolates <- func () () {
				if _error := _channel.callbacks.Closed (InboundFlow, _error); _error != nil {
					_channel.handleCallbacksError (_error)
				}
			}
			goto terminate
		} else if _packet != nil {
			select {
				case _channel.inboundPackets <- _packet :
					continue loop
				case _ = <- _channel.inboundSignal :
					goto signalled
			}
		} else {
			panic ("assertion")
		}
		panic ("fallthrough")
	}
	panic ("fallthrough")
	
	signalled :
	if !_channel.inboundActive {
		_channel.callbacksIsolates <- func () () {
			if _error := _channel.callbacks.Closed (InboundFlow, nil); _error != nil {
				_channel.handleCallbacksError (_error)
			}
		}
	} else {
		panic ("assertion")
	}
	
	terminate :
	if _channel.inboundPackets != nil {
		close (_channel.inboundPackets)
		_channel.inboundPackets = nil
	}
	_channel.transcript.TraceDebugging ("terminated the channel inbound background task.")
	_channel.terminateAcknowledgments.Done ()
}


func (_channel *channel) executeOutboundLoop () {
	_channel.transcript.TraceDebugging ("started the channel outbound background task.")
	
	loop : for {
		select {
			case _packet, _ok := <- _channel.outboundPackets :
				if !_ok {
					_channel.outboundPackets = nil
					goto signaled
				}
				if _error := _channel.pushOutboundPacket (_packet); _error != nil {
					_channel.outboundActive = false
					_channel.callbacksIsolates <- func () () {
						if _error := _channel.callbacks.Closed (OutboundFlow, _error); _error != nil {
							_channel.handleCallbacksError (_error)
						}
					}
					goto terminate
				} else {
					continue loop
				}
			case _ = <- _channel.outboundSignal :
				goto signaled
		}
		panic ("fallthrough")
	}
	panic ("fallthrough")
	
	signaled:
	if !_channel.outboundActive {
		_channel.callbacksIsolates <- func () () {
			if _error := _channel.callbacks.Closed (OutboundFlow, nil); _error != nil {
				_channel.handleCallbacksError (_error)
			}
		}
	} else {
		panic ("assertion")
	}
	
	terminate :
	if _channel.outboundPackets != nil {
		close (_channel.outboundPackets)
		_channel.outboundPackets = nil
	}
	_channel.transcript.TraceDebugging ("terminated the channel outbound background task.")
	_channel.terminateAcknowledgments.Done ()
}


func (_channel *channel) handleCallbacksError (_error error) () {
	panic (_error)
}

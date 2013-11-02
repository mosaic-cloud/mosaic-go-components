

package main


import "os"

import "mch/lib/channels"
import "mch/lib/messages"
import "vgl/transcript"


func main () () {
	
	_transcript.TraceInformation ("creating the channel...")
	_queue := newChannelQueue ()
	
	var _channel channels.Channel
	if false {
		_inboundStream := os.Stdin
		_outboundStream := os.Stdout
		if _channel_, _error := channels.Create (_queue, _inboundStream, _outboundStream, nil); _error != nil {
			panic (_error)
		} else {
			_channel = _channel_
		}
	} else {
		if _channel_, _error := channels.CreateAndDial (_queue, "tcp", "127.0.0.1:24704"); _error != nil {
			panic (_error)
		} else {
			_channel = _channel_
		}
	}
	
	for {
		var _error error
		
		_transcript.TraceInformation ("pulling a message from the channel...")
		var _rawInboundMessage *channels.Message
		_rawInboundMessage, _error = _queue.wait ()
		if _error != nil {
			panic (_error)
		}
		if _rawInboundMessage == nil {
			break
		}
		
		_transcript.TraceInformation ("decoding the message...")
		var _inboundMessage messages.Message
		_inboundMessage, _error = messages.Decode (_rawInboundMessage)
		if _error != nil {
			panic (_error)
		}
		
		switch _message := _message_1.(type) {
			case messages.ComponentCall :
				_transcript.TraceInformation ("replying...")
				_channel.Push (_message.ReturnSuccess (_message.Inputs, nil))
			default :
				_transcript.TraceError ("unknown message `%#v`", _message)
		}
	}
	
	_transcript.TraceInformation ("terminating the channel...")
	if _error := _channel.Terminate (); _error != nil {
		panic (_error)
	}
	
	_transcript.TraceInformation ("done.")
}


type channelQueue struct {
	messages chan *channels.Message
	errors chan error
}

func newChannelQueue () (*channelQueue) {
	return & channelQueue {
		messages : make (chan *channels.Message, 1024),
		errors : make (chan error, 16),
	}
}

func (_queue *channelQueue) wait () (*channels.Message, error) {
	select {
		case _message := <- _queue.messages :
			return _message, nil
		case _error := <- _queue.errors :
			return nil, _error
	}
}

func (_queue *channelQueue) Pull (_message *channels.Message) () {
	_queue.messages <- _message
}

func (_queue *channelQueue) Closed (_flow channels.Flow, _error error) () {
	switch _flow {
		case channels.InboundFlow :
			if _error != nil {
				_queue.errors <- _error
			} else {
				_queue.messages <- nil
			}
		case channels.OutboundFlow :
			if _error != nil {
				_queue.errors <- _error
			}
		default :
			panic ("assertion")
	}
}

func (_queue *channelQueue) Terminated (_error error) () {
	if _error != nil {
		_queue.errors <- _error
	}
}


var _transcript = transcript.NewPackageTranscript ()

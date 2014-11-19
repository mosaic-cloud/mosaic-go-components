

package main


import "errors"

import "vgl/transcript"

import . "mosaic-components/examples/simple-server"
import . "mosaic-components/libraries/messages"


type callbacks struct {
}


func (_callbacks *callbacks) Initialize (_server *SimpleServer) (error) {
	
	_server.ProcessExecutable = "/bin/sleep"
	_server.ProcessArguments = []string { "1h" }
	_server.ProcessEnvironment = map[string]string {}
	_server.SelfGroup = ""
	
	return nil
}


func (_callbacks *callbacks) Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error) {
	
	return nil, errors.New ("invalid-operation")
}


func main () () {
	PreMain (& callbacks {}, packageTranscript)
}


var packageTranscript = transcript.NewPackageTranscript (transcript.DebuggingLevel)

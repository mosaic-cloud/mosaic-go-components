

package simple_server


import "bufio"
import "errors"
import "net"
import "os"
import "syscall"
import "time"

import "mosaic-components/libraries/backend"
import "vgl/transcript"

import . "mosaic-components/libraries/messages"


func PreMain (_callbacks SimpleServerCallbacks, _packageTranscript transcript.Transcript) () {
	
	backend.PreMain (
			func (_componentIdentifier string, _channelEndpoint string, _configuration map[string]interface{}) (error) {
				return preMain (_componentIdentifier, _channelEndpoint, _configuration, _callbacks, _packageTranscript)
			})
	panic ("fallthrough")
}


func preMain (_componentIdentifier string, _channelEndpoint string, _configuration map[string]interface{}, _callbacks SimpleServerCallbacks, _packageTranscript transcript.Transcript) (error) {
	
	_server := & SimpleServer {}
	_server.callbacks = _callbacks
	_server.Identifier = ComponentIdentifier (_componentIdentifier)
	_server.Configuration = _configuration
	_server.Transcript = transcript.NewTranscript (_callbacks, _packageTranscript)
	
	_server.Temporary = os.Getenv ("mosaic_component_temporary")
	
	return backend.Execute (_server, _componentIdentifier, _channelEndpoint)
}


type SimpleServerCallbacks interface {
	Initialize (_server *SimpleServer) (_error error)
	Called (_server *SimpleServer, _operation ComponentOperation, _inputs interface{}) (_outputs interface{}, _error error)
}


type SimpleServer struct {
	backend backend.Controller
	callbacks SimpleServerCallbacks
	process *os.Process
	ProcessExecutable string
	ProcessArguments []string
	ProcessEnvironment map[string]string
	SelfGroup ComponentGroup
	Identifier ComponentIdentifier
	Configuration map[string]interface{}
	Transcript transcript.Transcript
	Temporary string
}


func (_server *SimpleServer) TcpSocketAcquire (_identifier ResourceIdentifier) (_ip net.IP, _port uint16, _fqdn string, _error error) {
	return backend.TcpSocketAcquireSync (_server.backend, _identifier)
}

func (_server *SimpleServer) TcpSocketResolve (_group ComponentGroup, _operation ComponentOperation) (_ip net.IP, _port uint16, _fqdn string, _error error) {
	
	var _ip_1 string
	
	if _outputs_1, _, _error := _server.backend.ComponentCallSync (ComponentIdentifier (_group), _operation, nil, nil); _error != nil {
		return nil, 0, "", _error
	} else {
		_outputs := _outputs_1.(map[string]interface{})
		_ip_1 = _outputs["ip"].(string)
		_port = uint16 (_outputs["port"].(float64))
		_fqdn = _outputs["fqdn"].(string)
	}
	
	_ip = net.ParseIP (_ip_1)
	if _ip == nil {
		return nil, 0, "", errors.New ("invalid IP address")
	}
	
	return _ip, _port, _fqdn, nil
}


func (_server *SimpleServer) Initialized (_backend backend.Controller) (error) {
	
	_server.Transcript.TraceInformation ("initializing the component...")
	_server.backend = _backend
	
	_server.Transcript.TraceInformation ("initializing the server...")
	if _error := _server.callbacks.Initialize (_server); _error != nil {
		panic (_error)
	}
	
	_server.Transcript.TraceInformation ("starting the server...")
	if _error := _server.startProcess (); _error != nil {
		panic (_error)
	}
	
	if _server.SelfGroup != "" {
		_server.Transcript.TraceInformation ("registering the component...")
		if _error := _server.backend.ComponentRegisterSync (_server.SelfGroup); _error != nil {
			panic (_error)
		}
	}
	
	_server.Transcript.TraceInformation ("initialized the component.")
	
	return nil
}


func (_server *SimpleServer) Terminated (_error error) (error) {
	
	_server.Transcript.TraceInformation ("terminating the component...")
	
	_server.Transcript.TraceInformation ("signaling the server...")
	if _error := _server.process.Signal (syscall.SIGTERM); _error != nil {
		panic (_error)
	}
	if _error := _server.process.Signal (syscall.SIGINT); _error != nil {
		panic (_error)
	}
	// FIXME: Find a better way to handle this!
	go func () () {
		time.Sleep (3 * time.Second)
		_server.process.Signal (syscall.SIGKILL)
	} ()
	
	_server.Transcript.TraceInformation ("waiting the server...")
	if _, _error := _server.process.Wait (); _error != nil {
		panic (_error)
	}
	
	_server.Transcript.TraceInformation ("terminated the component.")
	return nil
}


func (_server *SimpleServer) ComponentCallInvoked (_operation ComponentOperation, _inputs interface{}, _correlation Correlation, _attachment Attachment) (error) {
	
	_outputs, _error := _server.callbacks.Called (_server, _operation, _inputs)
	
	if _error == nil {
		if _error := _server.backend.ComponentCallSucceeded (_correlation, _outputs, nil); _error != nil {
			panic (_error)
		}
		return nil
	} else {
		if _error := _server.backend.ComponentCallFailed (_correlation, _error, nil); _error != nil {
			panic (_error)
		}
		return nil
	}
}


func (_server *SimpleServer) ComponentCastInvoked (_operation ComponentOperation, _inputs interface{}, _attachment Attachment) (error) {
	_server.Transcript.TraceError ("invoked unexpected component cast operation `%s`; ignoring!", _operation)
	return nil
}

func (_server *SimpleServer) ComponentCallSucceeded (_correlation Correlation, _outputs interface{}, _attachment Attachment) (error) {
	_server.Transcript.TraceError ("returned unexpected component call `%s`; ignoring!", _correlation)
	return nil
}

func (_server *SimpleServer) ComponentCallFailed (_correlation Correlation, _error interface{}, _attachment Attachment) (error) {
	_server.Transcript.TraceError ("returned unexpected component call `%s`; ignoring!", _correlation)
	return nil
}

func (_server *SimpleServer) ComponentRegisterSucceeded (_correlation Correlation) (error) {
	_server.Transcript.TraceError ("returned unexpected component register `%s`; ignoring!", _correlation)
	return nil
}

func (_server *SimpleServer) ComponentRegisterFailed (_correlation Correlation, _error interface{}) (error) {
	_server.Transcript.TraceError ("returned unexpected component register `%s`; ignoring!", _correlation)
	return nil
}

func (_server *SimpleServer) ResourceAcquireSucceeded (_correlation Correlation, _descriptor ResourceDescriptor) (error) {
	_server.Transcript.TraceError ("returned unexpected resource acquire `%s`; ignoring!", _correlation)
	return nil
}

func (_server *SimpleServer) ResourceAcquireFailed (_correlation Correlation, _error interface{}) (error) {
	_server.Transcript.TraceError ("returned unexpected resource acquire `%s`; ignoring!", _correlation)
	return nil
}


func (_server *SimpleServer) startProcess () (error) {
	
	_server.Transcript.TraceDebugging ("staring the server process...")
	
	var _reader, _writer *os.File
	if _reader_1, _writer_1, _error := os.Pipe (); _error != nil {
		panic (_error)
	} else {
		_reader = _reader_1
		_writer = _writer_1
	}
	
	go func () () {
		// FIXME: Handle errors!
		_scanner := bufio.NewScanner (_reader)
		for _scanner.Scan () {
			_server.Transcript.TraceInformation (">>  %s", _scanner.Text ())
		}
		_reader.Close ()
	} ()
	
	_executable := _server.ProcessExecutable
	_server.Transcript.TraceDebugging ("  * using the executable `%s`;", _executable)
	
	_arguments := []string {
			_executable,
	}
	for _, _argument := range _server.ProcessArguments {
		_arguments = append (_arguments, _argument)
	}
	_server.Transcript.TraceDebugging ("  * using the arguments `%#v`;", _arguments)
	
	_environment := []string {
	}
	for _environmentKey, _environmentValue := range _server.ProcessEnvironment {
		_environment = append (_environment, _environmentKey + "=" + _environmentValue)
	}
	_server.Transcript.TraceDebugging ("  * using the environment `%#v`;", _environment)
	
	_attributes := & os.ProcAttr {
			Env : _environment,
			Dir : "",
			Files : []*os.File {
					nil,
					_writer,
					_writer,
			},
	}
	
	if usePdeathSignal {
		_attributes.Sys = & syscall.SysProcAttr {
				Pdeathsig : syscall.SIGTERM,
		}
	}
	
	if _process_1, _error := os.StartProcess (_executable, _arguments, _attributes); _error != nil {
		_server.Transcript.TraceDebugging ("staring failed (while starting the process): `%s`!", _error.Error ())
		panic (_error)
	} else {
		_server.process = _process_1
	}
	
	_server.Transcript.TraceDebugging ("started the server process.")
	
	return nil
}


const usePdeathSignal = false

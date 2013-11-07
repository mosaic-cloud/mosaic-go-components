

package backend


import "encoding/json"
import "io/ioutil"
import "os"


func PreMain (_delegate func (string, string) (error)) () {
	
	_transcript := packageTranscript
	
	var _componentIdentifier string
	var _channelEndpoint string
	
	_arguments := os.Args
	if len (_arguments) < 1 {
		_transcript.TraceError ("invalid arguments (expected at least one)")
		os.Exit (1)
	}
	switch _arguments[1] {
		
		case "component" :
			if len (_arguments) != 3 {
				_transcript.TraceError ("invalid component arguments (expected only the identifier)")
				os.Exit (1)
			}
			_componentIdentifier = _arguments[2]
			_channelEndpoint = "stdio"
		
		case "component-init" :
			if len (_arguments) != 3 {
				_transcript.TraceError ("invalid component arguments (expected only the configuration file)")
				os.Exit (1)
			}
			var _configuration map[string]interface{}
			if _configurationData, _error := ioutil.ReadFile (_arguments[2]); _error != nil {
				_transcript.TraceError ("failed reading the configuration `%s`; aborting!", _error.Error ())
				os.Exit (1)
			} else if _error := json.Unmarshal (_configurationData, &_configuration); _error != nil {
				_transcript.TraceError ("failed parsing the configuration `%s`; aborting!", _error.Error ())
				os.Exit (1)
			} else {
				{
					_componentIdentifierValue, _ok := _configuration["component-identifier"]
					if !_ok {
						_transcript.TraceError ("invalid configuration: component identifier missing; aborting!")
						os.Exit (1)
					}
					_componentIdentifier, _ok = _componentIdentifierValue.(string)
					if !_ok {
						_transcript.TraceError ("invalid configuration: component identifier invalid; aborting!")
						os.Exit (1)
					}
				}
				{
					_channelEndpointValue, _ok := _configuration["channel-endpoint"]
					if !_ok {
						_transcript.TraceError ("invalid configuration: channel endpoint missing; aborting!")
						os.Exit (1)
					}
					_channelEndpoint, _ok = _channelEndpointValue.(string)
					if !_ok {
						_transcript.TraceError ("invalid configuration: channel endpoint invalid; aborting!")
						os.Exit (1)
					}
				}
			}
		
		case "standalone" :
			if len (_arguments) != 2 {
				_transcript.TraceError ("invalid standalone arguments (expected no others)")
				os.Exit (1)
			}
			_transcript.TraceError ("standalone is not implemented; aborting!")
			os.Exit (1)
		
		default :
			_transcript.TraceError ("invalid mode `%s`", _arguments[1])
			os.Exit (1)
	}
	
	if _error := _delegate (_componentIdentifier, _channelEndpoint); _error != nil {
		_transcript.TraceError ("delegate failed: `%s`; aborting!", _error.Error ())
		os.Exit (1)
	}
	
	os.Exit (0)
}

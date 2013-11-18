
package backend


import . "mosaic-components/libraries/messages"


type Backend interface {
	Terminate () (error)
	WaitTerminated () (error)
	// FIXME: Find a more elegant solution!
	TranscriptPush (_data Attachment) (error)
}


type Controller interface {
	Terminate () (error)
	ComponentCallInvoke (_component ComponentIdentifier, _operation ComponentOperation, _inputs interface{}, _attachment Attachment) (Correlation, error)
	ComponentCallSucceeded (_correlation Correlation, _outputs interface{}, _attachment Attachment) (error)
	ComponentCallFailed (_correlation Correlation, _error interface{}, _attachment Attachment) (error)
	ComponentCastInvoke (_component ComponentIdentifier, _operation ComponentOperation, _inputs interface{}, _attachment Attachment) (error)
	ComponentRegisterInvoke (_group ComponentGroup) (Correlation, error)
	ResourceAcquireInvoke (_specification ResourceSpecification) (Correlation, error)
	TranscriptPushInvoke (_data Attachment) (error)
	ComponentCallSync (_component ComponentIdentifier, _operation ComponentOperation, _inputs interface{}, _attachment Attachment) (interface{}, Attachment, error)
	ComponentRegisterSync (ComponentGroup) (error)
	ResourceAcquireSync (_specification ResourceSpecification) (ResourceDescriptor, error)
}


type Callbacks interface {
	Initialized (Controller) (error)
	Terminated (error) (error)
	MessageCallbacks
}


type ComponentInvokeCallbacks interface {
	ComponentCallInvoked (_operation ComponentOperation, _inputs interface{}, _correlation Correlation, _attachment Attachment) (error)
	ComponentCastInvoked (_operation ComponentOperation, _inputs interface{}, _attachment Attachment) (error)
}

type ComponentReturnCallbacks interface {
	ComponentCallSucceeded (_correlation Correlation, _outputs interface{}, _attachment Attachment) (error)
	ComponentCallFailed (_correlation Correlation, _error interface{}, _attachment Attachment) (error)
	ComponentRegisterSucceeded (_correlation Correlation) (error)
	ComponentRegisterFailed (_correlation Correlation, _error interface{}) (error)
}

type ResourceReturnCallbacks interface {
	ResourceAcquireSucceeded (_correlation Correlation, _descriptor ResourceDescriptor) (error)
	ResourceAcquireFailed (_correlation Correlation, _error interface{}) (error)
}

type MessageCallbacks interface {
	ComponentInvokeCallbacks
	ComponentReturnCallbacks
	ResourceReturnCallbacks
}

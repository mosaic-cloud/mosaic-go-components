

package messages


import "net"


type Message interface{}
type Correlation string
type Attachment []byte


type ComponentCall struct {
	Message
	Component ComponentIdentifier
	Operation ComponentOperation
	Inputs interface{}
	Correlation Correlation
	Attachment Attachment
}

type ComponentCallReturn struct {
	Message
	Ok bool
	Outputs interface{}
	Error interface{}
	Correlation Correlation
	Attachment Attachment
}

type ComponentCast struct {
	Message
	Component ComponentIdentifier
	Operation ComponentOperation
	Inputs interface{}
	Attachment Attachment
}

type ComponentRegister struct {
	Message
	Component ComponentIdentifier
	Group ComponentGroup
	Correlation Correlation
}

type ComponentRegisterReturn struct {
	Message
	Ok bool
	Error interface{}
	Correlation Correlation
}

type ComponentIdentifier string
type ComponentOperation string
type ComponentGroup string


type ResourceAcquire struct {
	Message
	Specifications []ResourceSpecification
	Correlation Correlation
}

type ResourceAcquireReturn struct {
	Message
	Ok bool
	Descriptors []ResourceDescriptor
	Error interface{}
	Correlation Correlation
}

type ResourceSpecification interface{}
type ResourceDescriptor interface{}
type ResourceIdentifier string

type TcpSocketSpecification struct {
	ResourceSpecification
	Identifier ResourceIdentifier
}

type TcpSocketDescriptor struct {
	ResourceDescriptor
	Identifier ResourceIdentifier
	Ip net.IP
	Port uint16
	Fqdn string
}


type TranscriptPush struct {
	Message
	Data Attachment
}


const NilComponentIdentifier = ""
const NilComponentGroup = ""
const NilComponentOperation = ""
const NilResourceIdentifier = ""
const NilCorrelation = ""

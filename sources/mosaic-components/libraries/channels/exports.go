

package channels


type Channel interface {
	Terminate () (error)
}


type Controller interface {
	Push (*Packet) (error)
	Close (Flow) (error)
	Terminate () (error)
}


type Callbacks interface {
	Initialized (Controller) (error)
	Pushed (*Packet) (error)
	Closed (Flow, error) (error)
	Terminated (error) (error)
}


type Flow int
const (
	InvalidFlowMin Flow = iota
	InboundFlow
	OutboundFlow
	InvalidFlowMax
)


type Packet struct {
	Data PacketData
	Attachment PacketAttachment
}

type PacketData map[string]interface{}
type PacketAttachment []byte

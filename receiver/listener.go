package receiver

import (
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type RespondFunc func(msg *dns.Msg) (*dns.Msg, error)

type Receiver struct{
	IfaceName string
	cnn *net.UDPConn

	RespondNameMatch map[string][]RespondFunc
}

func NewReceiver() *Receiver{
	r := new(Receiver)

	r.RespondNameMatch = make(map[string][]RespondFunc)

	return r
}

func (r *Receiver) BeginListen() error{
	ifi, err := net.InterfaceByName(r.IfaceName)

	if(err != nil){return err}

	listener, err := net.ListenMulticastUDP("udp", ifi, &net.UDPAddr{IP: net.ParseIP("224.0.0.251"), Port: 5353})

	if(err != nil){return err}
	
	r.cnn = listener

	go r.recv()

	return nil
}

func (r *Receiver) recv(){
	buf := make([]byte, 65536)

	for {
		n, from, err := r.cnn.ReadFrom(buf)

		if(err != nil){
			log.Printf("[ERR] ReadFrom %v\n", err)
			continue
		}

		var msg dns.Msg

		err = msg.Unpack(buf[:n])

		if(err != nil){
			log.Printf("[ERR] Unable to parse DNS Message from %v %v\n", from, err)
			continue
		}
		
		err = r.process_msg(msg, from)

		if(err != nil){
			log.Printf("[ERR] Error proocessing %v", err)
			continue
		}
	}
}

func (r *Receiver) Send(msg *dns.Msg, target net.Addr) error{
	buf, err := msg.Pack()

	if(err != nil){return err}

	_, err = r.cnn.WriteToUDP(buf, target.(*net.UDPAddr))

	return err
}

func (r *Receiver) SpoofFromSend(msg *dns.Msg, target, sender *net.UDPAddr) error{
	c, err := net.DialUDP("udp", sender, target)

	if(err != nil){return err}

	buf, err := msg.Pack()

	if(err != nil){return err}

	_, err = c.WriteToUDP(buf, target)

	return err
}

func (r *Receiver) process_msg(msg dns.Msg, from net.Addr) (error){
	//mdns only uses opcode 0 for both queries and responses
	if(msg.Opcode != dns.OpcodeQuery){
		return nil
	}

	for _, q := range msg.Question{
		//log.Printf("Qry for %v %v from %v\n", q.Name, dns.Type(q.Qtype).String(), from)

		//look for a name match
		for n, fncs := range r.RespondNameMatch{
			nl := strings.ToLower(q.Name)

			if(strings.Contains(nl, n) || n=="*"){
				//this can respond
				for _, f := range fncs{
					repl, err := f(&msg)

					if(err != nil){
						log.Printf("[ERR] Resp function error for %v %v", nl, err)
						continue
					}
					
					if(repl != nil){
						//log.Printf("[RESP] %v %v to %v", q.Name, dns.Type(q.Qtype).String(), from)
						r.Send(repl, from)
					}
				}
			}
		}
	}

	return nil
}
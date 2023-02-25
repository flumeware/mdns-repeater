package targetproxy

import (
	"github.com/miekg/dns"
)

type Proxy struct{
	TargetName string
}

func (p *Proxy) Respond(qry *dns.Msg) (*dns.Msg, error){
	cl := dns.Client{}

	cl.UDPSize = 4096

	resp,_, err := cl.Exchange(qry, p.TargetName)

	return resp, err
}
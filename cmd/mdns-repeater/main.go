package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/flumeware/mdns-repeater/receiver"
	"github.com/flumeware/mdns-repeater/targetproxy"
)

type config struct{
	Interfaces []string `json:"ifaces"`
	SingleTargets []single_target `json:"single_targets"`
}

type single_target struct{
	IP string `json:"ip"`
    Port int `json:"port"`
	Names []string `json:"names"`
}

func main(){
    var cfg config

    var cfg_filename = flag.String("c", "config.json", "Config file")

    flag.Parse()

    cfg_data, err := os.ReadFile(*cfg_filename)

    if(err != nil){
            log.Fatal(err)
    }

    err = json.Unmarshal(cfg_data, &cfg)

    if(err != nil){
            log.Fatal(err)
    }

    var rxs []*receiver.Receiver

    for _, i := range cfg.Interfaces{
        r := receiver.NewReceiver()

        r.IfaceName = i

        rxs = append(rxs, r)
    }

    for _, t := range cfg.SingleTargets{
        tp := new(targetproxy.Proxy)

        if(t.Port == 0){
            t.Port = 5353
        }

        tp.TargetName = fmt.Sprintf("%v:%v", t.IP, t.Port)

        for _, r := range rxs{
            for _, n := range t.Names{
                r.RespondNameMatch[n] = append(r.RespondNameMatch[n], tp.Respond)
            }
        }
    }

    for _, r := range rxs{
        go r.BeginListen()
    }

    select {}
}
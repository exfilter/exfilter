package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/exfilter/exfilter/pkg/exfilterlogger"
	openssltracer "github.com/exfilter/exfilter/pkg/openssl-tracer"
	ruleparser "github.com/exfilter/exfilter/pkg/rule-parser"
	tcpegresstracer "github.com/exfilter/exfilter/pkg/tcpegress-tracer"
	bpf "github.com/iovisor/gobpf/bcc"
	"github.com/mitchellh/go-ps"
)

func main() {
	prMap := ruleparser.ParseRuleFile("example.rules")

	m, err := tcpegresstracer.InitTCPTracer1(0)
	defer tcpegresstracer.DeInitTCPTracer(m)
	if err != nil {
		log.Fatal("error initializing tcp egress tracer:", err)
		return
	}

	m_tls, err := openssltracer.InitTLSTracer()
	defer openssltracer.DeInitTLSTracer(m_tls)
	if err != nil {
		log.Fatal("error initializing tls tracer:", err)
		return
	}

	table, err := tcpegresstracer.LoadBPFTable(m)
	if err != nil {
		log.Fatal("error loading bpf table for tcp egress tracer:", err)
		return
	}

	table_tls, err := openssltracer.LoadBPFTable(m_tls)
	if err != nil {
		log.Fatal("error loading bpf table for tls tracer:", err)
	}

	channel := make(chan []byte)
	tls_channel := make(chan []byte)

	perfMap, err := bpf.InitPerfMap(table, channel, nil)
	if err != nil {
		log.Fatal("failed to init perf map, ", err)
		return
	}

	perfMapTls, err := bpf.InitPerfMap(table_tls, tls_channel, nil)
	if err != nil {
		log.Fatal("failed to init perf map, ", err)
		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	f, err := exfilterlogger.InitLogger("exfilter.log")
	if err != nil {
		return
	}
	defer exfilterlogger.DeinitLogger(f)
	// fmt.Printf("%10s\t%10s\t%30s\t%30s\t%50s\n", "PID", "PROCESSNAME", "LADDR", "RADDR", "DATA")
	go func() {
		var event tcpegresstracer.TCPEgressEvent
		for {
			data := <-channel
			event.Pid = binary.LittleEndian.Uint32(data[0:4])
			event.Saddr = binary.LittleEndian.Uint32(data[4:8])
			event.Daddr = binary.LittleEndian.Uint32(data[8:12])
			event.Lport = binary.LittleEndian.Uint16(data[12:14])
			event.Dport = binary.LittleEndian.Uint16(data[14:16])
			event.DataLen = binary.LittleEndian.Uint32(data[16:20])
			event.Data = data[20:]

			// port match
			if prMap["dstPortRules"][int(event.Dport)] == nil {
				continue
			}

			// payload match
			isMatch := false
			var eventmsg string
			for _, ruleOpt := range prMap["dstPortRules"][int(event.Dport)] {
				if strings.Contains(strings.ToLower(string(event.Data)), strings.ToLower(ruleOpt.Content)) {
					isMatch = true
					eventmsg = ruleOpt.Message
					break
				}
			}

			if !isMatch {
				continue
			}

			p, _ := ps.FindProcess(int(event.Pid))
			logevent := exfilterlogger.EgressEvent{}
			logevent.Pid = event.Pid
			logevent.Saddr = tcpegresstracer.Inet_ntoa(event.Saddr) + ":" + strconv.Itoa(int(event.Lport))
			logevent.Daddr = tcpegresstracer.Inet_ntoa(event.Daddr) + ":" + strconv.Itoa(int(event.Dport))
			logevent.Data = string(event.Data)
			logevent.Msg = eventmsg

			exfilterlogger.LogEvent(logevent)

			fmt.Printf("%-10d\t%-10s\t%-30s\t%-30s\t%-50s\n", event.Pid, p.Executable(), tcpegresstracer.Inet_ntoa(event.Saddr)+":"+strconv.Itoa(int(event.Lport)), tcpegresstracer.Inet_ntoa(event.Daddr)+":"+strconv.Itoa(int(event.Dport)), string(event.Data))
		}
	}()

	go func() {
		var tls_event openssltracer.SSLDataEvent
		for {
			tls_data := <-tls_channel

			err := binary.Read(bytes.NewBuffer(tls_data), binary.LittleEndian, &tls_event)
			if err != nil {
				fmt.Printf("failed to decode received data: %s\n", err)
				continue
			}

			comm := string(tls_event.Data[:])
			var eventType string
			if openssltracer.AttachType(tls_event.EventType) != openssltracer.PROBE_ENTRY {
				continue
			}
			p, _ := ps.FindProcess(int(tls_event.Pid))
			fmt.Printf("%10d\t%10s\t%30s\t%8s\n", tls_event.Pid, p.Executable(), comm, eventType)
		}
	}()

	perfMap.Start()
	perfMapTls.Start()
	<-sig
	perfMap.Stop()
	perfMapTls.Stop()
}

// Code generated by bpf2go; DO NOT EDIT.
//go:build arm64be || armbe || mips || mips64 || mips64p32 || ppc64 || s390 || s390x || sparc || sparc64
// +build arm64be armbe mips mips64 mips64p32 ppc64 s390 s390x sparc sparc64

package probe

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadBpf returns the embedded CollectionSpec for bpf.
func loadBpf() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_BpfBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load bpf: %w", err)
	}

	return spec, err
}

// loadBpfObjects loads bpf and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//     *bpfObjects
//     *bpfPrograms
//     *bpfMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadBpfObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadBpf()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// bpfSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfSpecs struct {
	bpfProgramSpecs
	bpfMapSpecs
}

// bpfSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfProgramSpecs struct {
	KprobeSend     *ebpf.ProgramSpec `ebpf:"kprobe_send"`
	KprobeSendfile *ebpf.ProgramSpec `ebpf:"kprobe_sendfile"`
	KprobeSendmmsg *ebpf.ProgramSpec `ebpf:"kprobe_sendmmsg"`
	KprobeSendmsg  *ebpf.ProgramSpec `ebpf:"kprobe_sendmsg"`
	KprobeSendto   *ebpf.ProgramSpec `ebpf:"kprobe_sendto"`
	KprobeWrite    *ebpf.ProgramSpec `ebpf:"kprobe_write"`
	KprobeWritev   *ebpf.ProgramSpec `ebpf:"kprobe_writev"`
}

// bpfMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfMapSpecs struct {
	SendEvents     *ebpf.MapSpec `ebpf:"send_events"`
	SendfileEvents *ebpf.MapSpec `ebpf:"sendfile_events"`
	SendmmsgEvents *ebpf.MapSpec `ebpf:"sendmmsg_events"`
	SendmsgEvents  *ebpf.MapSpec `ebpf:"sendmsg_events"`
	SendtoEvents   *ebpf.MapSpec `ebpf:"sendto_events"`
	WriteEvents    *ebpf.MapSpec `ebpf:"write_events"`
	WritevEvents   *ebpf.MapSpec `ebpf:"writev_events"`
}

// bpfObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfObjects struct {
	bpfPrograms
	bpfMaps
}

func (o *bpfObjects) Close() error {
	return _BpfClose(
		&o.bpfPrograms,
		&o.bpfMaps,
	)
}

// bpfMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfMaps struct {
	SendEvents     *ebpf.Map `ebpf:"send_events"`
	SendfileEvents *ebpf.Map `ebpf:"sendfile_events"`
	SendmmsgEvents *ebpf.Map `ebpf:"sendmmsg_events"`
	SendmsgEvents  *ebpf.Map `ebpf:"sendmsg_events"`
	SendtoEvents   *ebpf.Map `ebpf:"sendto_events"`
	WriteEvents    *ebpf.Map `ebpf:"write_events"`
	WritevEvents   *ebpf.Map `ebpf:"writev_events"`
}

func (m *bpfMaps) Close() error {
	return _BpfClose(
		m.SendEvents,
		m.SendfileEvents,
		m.SendmmsgEvents,
		m.SendmsgEvents,
		m.SendtoEvents,
		m.WriteEvents,
		m.WritevEvents,
	)
}

// bpfPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfPrograms struct {
	KprobeSend     *ebpf.Program `ebpf:"kprobe_send"`
	KprobeSendfile *ebpf.Program `ebpf:"kprobe_sendfile"`
	KprobeSendmmsg *ebpf.Program `ebpf:"kprobe_sendmmsg"`
	KprobeSendmsg  *ebpf.Program `ebpf:"kprobe_sendmsg"`
	KprobeSendto   *ebpf.Program `ebpf:"kprobe_sendto"`
	KprobeWrite    *ebpf.Program `ebpf:"kprobe_write"`
	KprobeWritev   *ebpf.Program `ebpf:"kprobe_writev"`
}

func (p *bpfPrograms) Close() error {
	return _BpfClose(
		p.KprobeSend,
		p.KprobeSendfile,
		p.KprobeSendmmsg,
		p.KprobeSendmsg,
		p.KprobeSendto,
		p.KprobeWrite,
		p.KprobeWritev,
	)
}

func _BpfClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//go:embed bpf_bpfeb.o
var _BpfBytes []byte

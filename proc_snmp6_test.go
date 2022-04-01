package procfs

import "testing"

func TestProcSnmp6(t *testing.T) {
	p, err := getProcFixtures(t).Proc(26231)
	if err != nil {
		t.Fatal(err)
	}

	procSnmp6, err := p.Snmp6()
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		name string
		want float64
		have float64
	}{
		{name: "pid", want: 26231, have: float64(procSnmp6.PID)},
		{name: "Ip6InReceives", want: 92166, have: procSnmp6.Ip6.InReceives},
		{name: "Ip6InDelivers", want: 92053, have: procSnmp6.Ip6.InDelivers},
		{name: "Ip6OutNoRoutes", want: 169, have: procSnmp6.Ip6.OutNoRoutes},
		{name: "Ip6InOctets", want: 113479132, have: procSnmp6.Ip6.InOctets},
		{name: "Icmp6InMsgs", want: 142, have: procSnmp6.Icmp6.InMsgs},
		{name: "Udp6InDatagrams", want: 2016, have: procSnmp6.Udp6.InDatagrams},
		{name: "UdpLite6InDatagrams", want: 0, have: procSnmp6.UdpLite6.InDatagrams},
	} {
		if test.want != test.have {
			t.Errorf("want %s %f, have %f", test.name, test.want, test.have)
		}
	}

}

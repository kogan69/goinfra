package pnet

import (
	"bytes"
	"net"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func Test_PNet_Parse(t *testing.T) {

	pNet, err := ParseIp("194.90.1.5")
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(pNet)

	pNet, err = ParseIp("194.90.1.5/32")
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(pNet)

	pNet, err = ParseIp("::ffff:85.238.101.249/32")
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(pNet)

}

var parseIPTests = []struct {
	in  string
	out net.IP
}{
	{"127.0.1.2", net.IPv4(127, 0, 1, 2)},
	{"127.0.0.1", net.IPv4(127, 0, 0, 1)},
	{"::ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	{"::ffff:7f01:0203", net.IPv4(127, 1, 2, 3)},
	{"0:0:0:0:0000:ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	{"0:0:0:0:000000:ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},
	{"0:0:0:0::ffff:127.1.2.3", net.IPv4(127, 1, 2, 3)},

	{"2001:4860:0:2001::68", net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}},
	{"2001:4860:0000:2001:0000:0000:0000:0068", net.IP{0x20, 0x01, 0x48, 0x60, 0, 0, 0x20, 0x01, 0, 0, 0, 0, 0, 0, 0x00, 0x68}},

	{"-0.0.0.0", nil},
	{"0.-1.0.0", nil},
	{"0.0.-2.0", nil},
	{"0.0.0.-3", nil},
	{"127.0.0.256", nil},
	{"abc", nil},
	{"123:", nil},
	{"fe80::1%lo0", nil},
	{"fe80::1%911", nil},
	{"", nil},
	{"a1:a2:a3:a4::b1:b2:b3:b4", nil}, // Issue 6628
	{"127.001.002.003", nil},
	{"::ffff:127.001.002.003", nil},
	{"123.000.000.000", nil},
	{"1.2..4", nil},
	{"0123.0.0.1", nil},
}

func TestParseIP(t *testing.T) {
	for _, tt := range parseIPTests {
		pNet, err := ParseIp(tt.in)
		if err != nil {
			if tt.out != nil {
				t.Fatal(err)
			} else {
				continue
			}

		}
		if !pNet.IPNet.IP.Equal(tt.out) {
			t.Fatal("failed")
		}
	}
}

var parseCIDRTests = []struct {
	in  string
	ip  net.IP
	net *net.IPNet
	err error
}{
	{"135.104.0.0/32", net.IPv4(135, 104, 0, 0), &net.IPNet{IP: net.IPv4(135, 104, 0, 0), Mask: net.IPv4Mask(255,
		255,
		255,
		255)}, nil},
	{"0.0.0.0/24", net.IPv4(0, 0, 0, 0), &net.IPNet{IP: net.IPv4(0, 0, 0, 0), Mask: net.IPv4Mask(255,
		255,
		255,
		0)}, nil},
	{"135.104.0.0/24", net.IPv4(135, 104, 0, 0), &net.IPNet{IP: net.IPv4(135, 104, 0, 0), Mask: net.IPv4Mask(255,
		255,
		255,
		0)}, nil},
	{"135.104.0.1/32", net.IPv4(135, 104, 0, 1), &net.IPNet{IP: net.IPv4(135, 104, 0, 1), Mask: net.IPv4Mask(255,
		255,
		255,
		255)}, nil},
	{"135.104.0.1/24", net.IPv4(135, 104, 0, 1), &net.IPNet{IP: net.IPv4(135, 104, 0, 0), Mask: net.IPv4Mask(255,
		255,
		255,
		0)}, nil},
	{"::1/128", net.ParseIP("::1"), &net.IPNet{IP: net.ParseIP("::1"), Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"))}, nil},
	{"abcd:2345::/127", net.ParseIP("abcd:2345::"), &net.IPNet{IP: net.ParseIP("abcd:2345::"), Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff:ffff:ffff:ffff:fffe"))}, nil},
	{"abcd:2345::/65", net.ParseIP("abcd:2345::"), &net.IPNet{IP: net.ParseIP("abcd:2345::"), Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff:8000::"))}, nil},
	{"abcd:2345::/64", net.ParseIP("abcd:2345::"), &net.IPNet{IP: net.ParseIP("abcd:2345::"), Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff:ffff::"))}, nil},
	{"abcd:2345::/63", net.ParseIP("abcd:2345::"), &net.IPNet{IP: net.ParseIP("abcd:2345::"), Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff:fffe::"))}, nil},
	{"abcd:2345::/33", net.ParseIP("abcd:2345::"), &net.IPNet{IP: net.ParseIP("abcd:2345::"), Mask: net.IPMask(net.ParseIP("ffff:ffff:8000::"))}, nil},
	{"abcd:2345::/32", net.ParseIP("abcd:2345::"), &net.IPNet{IP: net.ParseIP("abcd:2345::"), Mask: net.IPMask(net.ParseIP("ffff:ffff::"))}, nil},
	{"abcd:2344::/31", net.ParseIP("abcd:2344::"), &net.IPNet{IP: net.ParseIP("abcd:2344::"), Mask: net.IPMask(net.ParseIP("ffff:fffe::"))}, nil},
	{"abcd:2300::/24", net.ParseIP("abcd:2300::"), &net.IPNet{IP: net.ParseIP("abcd:2300::"), Mask: net.IPMask(net.ParseIP("ffff:ff00::"))}, nil},
	{"abcd:2345::/24", net.ParseIP("abcd:2345::"), &net.IPNet{IP: net.ParseIP("abcd:2300::"), Mask: net.IPMask(net.ParseIP("ffff:ff00::"))}, nil},
	{"2001:DB8::/48", net.ParseIP("2001:DB8::"), &net.IPNet{IP: net.ParseIP("2001:DB8::"), Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff::"))}, nil},
	{"2001:DB8::1/48", net.ParseIP("2001:DB8::1"), &net.IPNet{IP: net.ParseIP("2001:DB8::"), Mask: net.IPMask(net.ParseIP("ffff:ffff:ffff::"))}, nil},
}

func TestParseCIDR(t *testing.T) {
	for _, tt := range parseCIDRTests {
		ip, net, err := net.ParseCIDR(tt.in)
		if !reflect.DeepEqual(err, tt.err) {
			t.Errorf("ParseCIDR(%q) = %v, %v; want %v, %v", tt.in, ip, net, tt.ip, tt.net)
		}
		if err == nil && (!tt.ip.Equal(ip) || !tt.net.IP.Equal(net.IP) || !reflect.DeepEqual(net.Mask, tt.net.Mask)) {
			t.Errorf("ParseCIDR(%q) = %v, {%v, %v}; want %v, {%v, %v}",
				tt.in,
				ip,
				net.IP,
				net.Mask,
				tt.ip,
				tt.net.IP,
				tt.net.Mask)
		}
	}
}

func TestPNetParse(t *testing.T) {
	for _, tt := range parseCIDRTests {
		//ip, net, err := net.ParseCIDR(tt.in)

		pNet, err := ParseIp(tt.in)
		if err != nil {
			t.Fatal(err)
		}
		if !pNet.IPNet.IP.Equal(tt.ip) {
			t.Fatal()
		}
		if !bytes.Equal(pNet.IPNet.Mask, tt.net.Mask) {
			t.Fatal()
		}
	}
}

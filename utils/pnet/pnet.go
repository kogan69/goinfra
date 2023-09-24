package pnet

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net"

	"github.com/jackc/pgtype"
)

type CIDR pgtype.Inet

func ParseIp(s string) (pNet CIDR, err error) {
	i, n, e := net.ParseCIDR(s)
	if e == nil {
		pNet.IPNet = &net.IPNet{IP: i, Mask: n.Mask}
		pNet.Status = pgtype.Present
		return
	} else {
		err = nil
		addr := net.ParseIP(s)
		if addr == nil {
			err = errors.New(fmt.Sprintf("failed to parse: %s", s))
			return
		}
		pNet.IPNet = &net.IPNet{IP: addr}
		if len(addr) == net.IPv4len {
			pNet.IPNet.Mask = net.CIDRMask(32, 32)
		} else if len(addr) == net.IPv6len {
			pNet.IPNet.Mask = net.CIDRMask(128, 128)
		}
		pNet.Status = pgtype.Present
	}
	return
}
func (p CIDR) String() string {
	return p.IPNet.IP.String()
}

func (p CIDR) IsUnspecified() bool {
	return p.IPNet.IP.IsUnspecified()
}

func UnspecifiedV4() CIDR {
	return CIDR{
		IPNet:  &net.IPNet{IP: net.IPv4zero, Mask: net.IPv4Mask(0, 0, 0, 0)},
		Status: pgtype.Present,
	}
}
func (dst *CIDR) Set(src interface{}) error {
	return (*pgtype.Inet)(dst).Set(src)
}

func (dst CIDR) Get() interface{} {
	return (pgtype.Inet)(dst).Get()
}

func (src *CIDR) AssignTo(dst interface{}) error {
	return (*pgtype.Inet)(src).AssignTo(dst)
}

func (dst *CIDR) DecodeText(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.Inet)(dst).DecodeText(ci, src)
}

func (dst *CIDR) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	return (*pgtype.Inet)(dst).DecodeBinary(ci, src)
}

func (src CIDR) EncodeText(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (pgtype.Inet)(src).EncodeText(ci, buf)
}

func (src CIDR) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	return (pgtype.Inet)(src).EncodeBinary(ci, buf)
}

// Scan implements the database/sql Scanner interface.
func (dst *CIDR) Scan(src interface{}) error {
	return (*pgtype.Inet)(dst).Scan(src)
}

// Value implements the database/sql/driver Valuer interface.
func (src CIDR) Value() (driver.Value, error) {
	return (pgtype.Inet)(src).Value()
}
func (p CIDR) ToIp() net.IP {
	return p.IPNet.IP
}

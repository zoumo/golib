// Copyright 2023 jim.zoumo@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package netutil

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	// For connection setup and write operations.
	errMissingAddress = errors.New("missing address")
	// For connection setup operations.
	errNoSuitableAddress = errors.New("no suitable address found")
)

// BalancedDialer
type BalancedDialer interface {
	// DialContext connects to the address on the named network with
	// client side load balance.
	//
	// Balanceable networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only),
	// "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), other networks will not be
	// balanced.
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// A Resolver looks up names and numbers.
// It is an interface that represents net.Resolver
type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) (addrs []net.IPAddr, err error)
	LookupPort(ctx context.Context, network, service string) (port int, err error)
}

// BalancerBuilder build a new balancer instance
type BalancerBuilder interface {
	Build(host string, addrList []net.Addr) Balancer
}

// Balancer resort the addresses coming from resolver and
// put the first priority address to the head
type Balancer interface {
	Balance(ctx context.Context, addrList []net.Addr) []net.Addr
}

type Options struct {
	// BalancerBuilder build a client side load balancer
	BalancerBuilder BalancerBuilder
	// custom resolver, If not set, net.DefaultResolver will be used
	Resolver Resolver
	// custom dail function, If not set, net.DailContext will be used
	dialer func(ctx context.Context, network, address string) (net.Conn, error)
}

var _ BalancedDialer = &baseBalancedDialer{}

type baseBalancedDialer struct {
	resolver        Resolver
	dial            func(ctx context.Context, network, address string) (net.Conn, error)
	balancerbuilder BalancerBuilder
	balancers       sync.Map
}

func NewBalancedDialer(opt Options) BalancedDialer {
	d := &baseBalancedDialer{}
	if opt.Resolver != nil {
		d.resolver = opt.Resolver
	} else {
		d.resolver = net.DefaultResolver
	}
	if opt.dialer != nil {
		d.dial = opt.dialer
	} else {
		d.dial = (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext
	}
	if opt.BalancerBuilder != nil {
		d.balancerbuilder = opt.BalancerBuilder
	} else {
		d.balancerbuilder = &rrBalancerBuilder{}
	}
	return d
}

func (d *baseBalancedDialer) DialContext(ctx context.Context, network, host string) (net.Conn, error) {
	if ctx == nil {
		panic("nil context")
	}

	if !isBalanceableNetwork(network) {
		return d.dial(ctx, network, host)
	}

	addrs, err := d.lookupAddrs(ctx, network, host)
	if err != nil {
		return nil, err
	}
	return d.dialSerial(ctx, network, host, addrs)
}

func (d *baseBalancedDialer) lookupAddrs(ctx context.Context, network, addr string) (AddrList, error) {
	if !isBalanceableNetwork(network) {
		return nil, fmt.Errorf("unsupport network %v", network)
	}
	if addr == "" {
		return nil, errMissingAddress
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	portnum, err := d.resolver.LookupPort(ctx, network, port)
	if err != nil {
		return nil, err
	}

	inetaddr := func(ip net.IPAddr) net.Addr {
		switch network {
		case "tcp", "tcp4", "tcp6":
			return &net.TCPAddr{IP: ip.IP, Port: portnum, Zone: ip.Zone}
		case "udp", "udp4", "udp6":
			return &net.UDPAddr{IP: ip.IP, Port: portnum, Zone: ip.Zone}
		default:
			panic("unexpected network: " + network)
		}
	}

	ips, err := d.resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	if len(ips) == 1 && ips[0].IP.Equal(net.IPv6unspecified) {
		ips = append(ips, net.IPAddr{IP: net.IPv4zero})
	}

	var filter func(net.IPAddr) bool
	if network != "" && network[len(network)-1] == '4' {
		filter = ipv4only
	}
	if network != "" && network[len(network)-1] == '6' {
		filter = ipv6only
	}
	return filterAddrList(filter, ips, inetaddr, host)
}

func (d *baseBalancedDialer) dialSerial(ctx context.Context, network, host string, addrList AddrList) (net.Conn, error) {
	if len(addrList) > 1 {
		// TODO: consider cgo dns resolver
		// purego dns will always return the same dns list
		// sort.Sort(addrList)

		// get balancer to resort addresses
		b, ok := d.balancers.Load(host)
		if !ok {
			b, _ = d.balancers.LoadOrStore(host, d.balancerbuilder.Build(host, addrList))
		}
		balancer := b.(Balancer)
		addrList = balancer.Balance(ctx, addrList)
	}
	var firstErr error
	for _, addr := range addrList {
		addrstr := addr.String()
		c, err := d.dial(ctx, network, addrstr)
		if err == nil {
			return c, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	if firstErr == nil {
		firstErr = &net.OpError{Op: "dial", Net: network, Source: nil, Addr: nil, Err: errMissingAddress}
	}
	return nil, firstErr
}

type AddrList []net.Addr

func (s AddrList) Len() int {
	return len(s)
}

func (s AddrList) Less(i, j int) bool {
	return s[i].String() < s[j].String()
}

func (s AddrList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func isBalanceableNetwork(network string) bool {
	switch network {
	case "tcp", "tcp4", "tcp6", "udp", "udp4", "udp6":
		return true
	}
	return false
}

// filterAddrList applies a filter to a list of IP addresses,
// yielding a list of Addr objects. Known filters are nil, ipv4only,
// and ipv6only. It returns every address when the filter is nil.
// The result contains at least one address when error is nil.
func filterAddrList(filter func(net.IPAddr) bool, ips []net.IPAddr, inetaddr func(net.IPAddr) net.Addr, originalAddr string) (AddrList, error) {
	var addrs []net.Addr
	for _, ip := range ips {
		if filter == nil || filter(ip) {
			addrs = append(addrs, inetaddr(ip))
		}
	}
	if len(addrs) == 0 {
		return nil, &net.AddrError{Err: "", Addr: originalAddr}
	}
	return addrs, nil
}

// ipv4only reports whether addr is an IPv4 address.
func ipv4only(addr net.IPAddr) bool {
	return addr.IP.To4() != nil
}

// ipv6only reports whether addr is an IPv6 address except IPv4-mapped IPv6 address.
func ipv6only(addr net.IPAddr) bool {
	return len(addr.IP) == net.IPv6len && addr.IP.To4() == nil
}

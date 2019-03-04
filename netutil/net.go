/*
Copyright 2018 Jim Zhang (jim.zoumo@gmail.com). All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package netutil

import (
	"fmt"
	"net"
)

// Interface represents the local network interface
type Interface struct {
	net.Interface
}

// Addrs returns a list of unicast interface addresses for a specific
// interface.
func (dev *Interface) Addrs() (AddrSlice, error) {
	addrs, err := dev.Interface.Addrs()
	if err != nil {
		return nil, err
	}

	ret := []Addr{}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok {
			ret = append(ret, Addr{ipnet})
		}
	}
	return ret, nil
}

// IsLoopback returns true if the net interface is lookback
func (dev *Interface) IsLoopback() bool {
	return (dev.Flags & net.FlagLoopback) > 0
}

// Addr represents a network end point address.
type Addr struct {
	*net.IPNet
}

// IPAddr returns the string form of the IP address ip.
// It returns one of 4 forms:
//   - "<nil>", if ip has length 0
//   - dotted decimal ("192.0.2.1"), if ip is an IPv4 or IP4-mapped IPv6 address
//   - IPv6 ("2001:db8::1"), if ip is a valid IPv6 address
//   - the hexadecimal form of ip, without punctuation, if no other cases appl
func (addr Addr) IPAddr() string {
	return addr.IP.String()
}

// MaskSize returns the number of leading ones
func (addr Addr) MaskSize() int {
	mask, _ := addr.Mask.Size()
	return mask
}

// IsIPv4 returns true if the ip is not an IPv4 address
func (addr Addr) IsIPv4() bool {
	return addr.IP.To4() != nil
}

// IsIPv6 returns true if the ip is not an IPv6 address
func (addr Addr) IsIPv6() bool {
	return len(addr.IP) == net.IPv6len
}

// IsLoopback reports whether ip is a loopback address.
func (addr Addr) IsLoopback() bool {
	return addr.IP.IsLoopback()
}

// AddrSlice reprecents a list of ip addresses
type AddrSlice []Addr

// Contains checks if the ip is in the collection
func (addrs AddrSlice) Contains(ip string) bool {
	for _, addr := range addrs {
		if addr.IPAddr() == ip {
			return true
		}
	}
	return false
}

// InterfaceSlice reprecents a list of net interfaces
type InterfaceSlice []Interface

// Get returns net interface device if it exists in the collection
func (ifaces InterfaceSlice) Get(name string) *Interface {
	for _, iface := range ifaces {
		if iface.Name == name {
			return &iface
		}
	}
	return nil
}

// Contains checks if the net interface is in the collection
func (ifaces InterfaceSlice) Contains(name string) bool {
	for _, iface := range ifaces {
		if iface.Name == name {
			return true
		}
	}
	return false
}

// Filter filters some interfaces by filsterFunc if it returns true
func (ifaces InterfaceSlice) Filter(filterFunc func(iface Interface) bool) InterfaceSlice {
	ret := []Interface{}
	for _, iface := range ifaces {
		if filterFunc(iface) {
			continue
		}
		ret = append(ret, iface)
	}
	return ret
}

// One returns the first net interface in the list
// It returns nil if there is no element in the slice
func (ifaces InterfaceSlice) One() *Interface {
	if len(ifaces) == 0 {
		return nil
	}
	return &ifaces[0]
}

// Interfaces returns a slice containing the local network interfaces
func Interfaces() (InterfaceSlice, error) {
	ret := []Interface{}
	ifaces, err := net.Interfaces()
	if err != nil {
		return ret, err
	}
	for _, iface := range ifaces {
		ret = append(ret, Interface{iface})
	}
	return ret, nil
}

// InterfacesByLoopback returns a list of loopback network interfaces.
func InterfacesByLoopback() (InterfaceSlice, error) {
	ifaces, err := Interfaces()
	if err != nil {
		return nil, err
	}
	ret := []Interface{}
	for _, iface := range ifaces {
		if iface.IsLoopback() {
			ret = append(ret, iface)
		}
	}
	return ret, nil
}

// InterfacesByIP returns the local network interfaces that is using the
// specified IP address.
func InterfacesByIP(ip string) (InterfaceSlice, error) {
	ifaces, err := Interfaces()
	if err != nil {
		return nil, err
	}
	ret := []Interface{}
	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		if addrs.Contains(ip) {
			ret = append(ret, iface)
		}
	}
	return ret, nil
}

// InterfaceByName returns the local network interface specified by name.
func InterfaceByName(name string) (*Interface, error) {
	slice, err := Interfaces()
	if err != nil {
		return nil, err
	}

	one := slice.Get(name)
	if one == nil {
		return nil, fmt.Errorf("no such network interface named %v", name)
	}
	return one, nil
}

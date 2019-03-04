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
	"net"
	"reflect"
	"testing"
)

var (
	ifaceLo = Interface{
		Interface: net.Interface{
			Name: "lo",
		},
	}
	ifaceEth0 = Interface{
		Interface: net.Interface{
			Name: "eth0",
		},
	}

	ifaceCases = InterfaceSlice{
		ifaceLo,
		ifaceEth0,
	}
)

func TestInterfacesAndSlice(t *testing.T) {
	ifaces, _ := InterfacesByLoopback()
	if len(ifaces) != 1 {
		t.Errorf("InterfacesByLoopback() got = %v, want = 1", len(ifaces))
	}

	_, err := InterfaceByName("noSuchInterface")
	if err == nil {
		t.Errorf("InterfaceByName() err = %v, wantErr = %v", err, false)
	}

	for _, iface := range ifaces {
		ifaceByName, _ := InterfaceByName(iface.Name)
		if reflect.DeepEqual(iface, ifaceByName) {
			t.Errorf("InterfaceByName() got = %v, want %v", ifaceByName, iface)
		}
		if !ifaces.Contains(iface.Name) {
			t.Errorf("InterfaceSlice.Contains() got = %v, want %v", ifaces.Contains(iface.Name), true)
		}

		addrs, _ := iface.Addrs()
		if !addrs.Contains("127.0.0.1") {
			t.Errorf("AddrSlice.Contains() got = %v, want %v", addrs.Contains(iface.Name), true)
		}
	}

	slice, _ := InterfacesByIP("127.0.0.1")
	if !reflect.DeepEqual(ifaces, slice) {
		t.Errorf("InterfacesByIP got = %v, want = %v", slice, ifaces)
	}
}

func TestInterfaceSlice_Filter(t *testing.T) {
	tests := []struct {
		name       string
		ifaces     InterfaceSlice
		filterFunc func(iface Interface) bool
		want       InterfaceSlice
	}{
		{
			"",
			ifaceCases,
			func(iface Interface) bool {
				if iface.Name == "lo" {
					return true
				}
				return false
			},
			InterfaceSlice{
				ifaceEth0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ifaces.Filter(tt.filterFunc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceSlice.Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceSlice_Get(t *testing.T) {
	tests := []struct {
		name   string
		ifaces InterfaceSlice
		iface  string
		want   *Interface
	}{
		{"", ifaceCases, "lo", &ifaceLo},
		{"", ifaceCases, "eth0", &ifaceEth0},
		{"", ifaceCases, "xxx", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.ifaces.Get(tt.iface)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceSlice.Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceSlice_Contains(t *testing.T) {
	tests := []struct {
		name   string
		ifaces InterfaceSlice
		iface  string
		want   bool
	}{
		{"", ifaceCases, "lo", true},
		{"", ifaceCases, "xxx", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ifaces.Contains(tt.iface); got != tt.want {
				t.Errorf("InterfaceSlice.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterfaceSlice_One(t *testing.T) {
	tests := []struct {
		name   string
		ifaces InterfaceSlice
		want   *Interface
	}{
		{"", ifaceCases, &ifaceLo},
		{"", nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ifaces.One(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InterfaceSlice.One() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddr_IsIPv4(t *testing.T) {
	tests := []struct {
		name string
		addr Addr
		want bool
	}{
		{"", Addr{&net.IPNet{IP: net.IPv4(127, 0, 0, 1)}}, true},
		{"", Addr{&net.IPNet{IP: []byte{1, 1, 1, 1}}}, true},
		{"", Addr{&net.IPNet{IP: []byte{1, 1, 1, 1, 1}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.addr.IsIPv4(); got != tt.want {
				t.Errorf("Addr.IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddr_IsIPv6(t *testing.T) {
	tests := []struct {
		name string
		addr Addr
		want bool
	}{
		{"", Addr{&net.IPNet{IP: net.IPv4(127, 0, 0, 1)}}, true},
		{"", Addr{&net.IPNet{IP: net.ParseIP("fe80::1ba8:6946:ab33:da50")}}, true},
		{"", Addr{&net.IPNet{IP: []byte{1, 1, 1, 1, 1}}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.addr.IsIPv6(); got != tt.want {
				t.Errorf("Addr.IsIPv6() = %v, want %v", got, tt.want)
			}
		})
	}
}

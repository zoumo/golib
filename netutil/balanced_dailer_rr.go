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
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	r  = rand.New(rand.NewSource(time.Now().UnixNano()))
	mu sync.Mutex
)

// Intn implements rand.Intn on the global source.
func randIntn(n int) int {
	mu.Lock()
	defer mu.Unlock()
	return r.Intn(n)
}

// rrBalancerBuilder create a RoundRobin Balancer
type rrBalancerBuilder struct{}

func (b *rrBalancerBuilder) Build(host string, addrs []net.Addr) Balancer {
	return &rrBalancer{
		host: host,
		// Start at a random index
		next: uint64(randIntn(len(addrs))),
	}
}

type rrBalancer struct {
	host string
	next uint64
}

func (b *rrBalancer) Balance(ctx context.Context, addrs []net.Addr) []net.Addr {
	addrsLen := len(addrs)
	if addrsLen <= 1 {
		return addrs
	}
	nextIndex := atomic.AddUint64(&b.next, 1)
	nextIndex = nextIndex % uint64(addrsLen)
	newAddr := addrs[nextIndex:]
	newAddr = append(newAddr, addrs[0:nextIndex]...)
	return newAddr
}

/**
 * Copyright 2018-2023 Wargaming Group Limited
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
**/
package befw

import (
	"bytes"
	"net"
	"testing"
)

func TestNet2Strings(t *testing.T) {
	table := map[*net.IPNet]string{
		{
			IP:   net.IPv4(192, 168, 0, 5),
			Mask: net.IPv4Mask(255, 255, 255, 0),
		}: "192.168.0.0/24",
		{
			IP:   net.IPv4(192, 168, 0, 1),
			Mask: net.IPv4Mask(255, 255, 255, 255),
		}: "192.168.0.1/32",
	}
	keys := make([]*net.IPNet, 0)
	values := make([]string, 0)
	for k, v := range table {
		keys = append(keys, k)
		values = append(values, v)
	}
	for i, x := range nets2string(keys) {
		if values[i] != x {
			t.Errorf("Value for %d doesn't match: %s != %s", i, x, values[i])
		}
	}
}

func TestPath2ipnet(t *testing.T) {
	table := map[string]*net.IPNet{
		"befw/$alias$/$test$/192.168.0.5/24": {
			IP:   net.IPv4(192, 168, 0, 0),
			Mask: net.IPv4Mask(255, 255, 255, 0),
		},
		"befw/sercvice_tcp_2200/192.168.0.5": {
			IP:   net.IPv4(192, 168, 0, 5),
			Mask: net.IPv4Mask(255, 255, 255, 255),
		},
		"befw/sercvice_tcp_2200/10.0.0.5/8": {
			IP:   net.IPv4(10, 0, 0, 0),
			Mask: net.IPv4Mask(255, 0, 0, 0),
		},
		"befw/sercvice_tcp_2200/::1:5ee:bad:c0de/96": {
			IP:   net.IP{0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0x01, 0x05, 0xee, 0, 0, 0, 0},
			Mask: net.CIDRMask(96, 128),
		},
		"befw/sercvice_tcp_2200/cafe:feed::/127": {
			IP:   net.IP{0xca, 0xfe, 0xfe, 0xed, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Mask: net.CIDRMask(127, 128),
		},
		"befw/sercvice_tcp_2200/::/0": {
			IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			Mask: net.CIDRMask(0, 128),
		},
	}
	equals := func(a, b *net.IPNet) bool {
		return a.IP.Equal(b.IP) && bytes.Equal([]byte(a.Mask), []byte(b.Mask))
	}
	for i, x := range table {
		if n := path2ipnet(i); n == nil || !equals(n, x) {
			t.Errorf("Value for %s doesn't match: %s != %s", i, x.String(), n)
		}
	}

}

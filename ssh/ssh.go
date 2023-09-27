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

package ssh

import (
	"context"
	"io"
	"net"
	"time"

	"github.com/spf13/afero"
	"golang.org/x/crypto/ssh"

	"github.com/zoumo/golib/ssh/scp"
	"github.com/zoumo/golib/ssh/shell"
)

func SetClientConfigDefaults(cfg *ssh.ClientConfig) {
	cfg.SetDefaults()
	if cfg.HostKeyCallback == nil {
		cfg.HostKeyCallback = ssh.InsecureIgnoreHostKey()
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	if cfg.User == "" {
		cfg.User = "root"
	}
}

func DialTCP(host, port string, cfg *ssh.ClientConfig) (*Client, error) {
	return Dial("tcp", host, port, cfg)
}

func Dial(network, host, port string, cfg *ssh.ClientConfig) (*Client, error) {
	SetClientConfigDefaults(cfg)

	c, err := ssh.Dial(network, net.JoinHostPort(host, port), cfg)

	if err != nil {
		return nil, err
	}
	return &Client{c}, err
}

type Client struct {
	*ssh.Client
}

func (c *Client) Dial(network, host, port string, cfg *ssh.ClientConfig) (*Client, error) {
	SetClientConfigDefaults(cfg)

	conn, err := c.Client.Dial(network, net.JoinHostPort(host, port))
	if err != nil {
		return nil, err
	}
	ncc, chans, reqs, err := ssh.NewClientConn(conn, net.JoinHostPort(host, port), cfg)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(ncc, chans, reqs)
	return &Client{client}, nil
}

func (c *Client) Shell(stdin io.Reader, stdout, stderr io.Writer) error {
	shell := shell.New(c.Client)
	return shell.Run(stdin, stdout, stderr)
}

func (c *Client) Upload(ctx context.Context, local, remote string) error {
	scp := scp.New(c.Client, afero.NewOsFs(), nil)
	return scp.Upload(ctx, local, remote)
}

func (c *Client) Download(ctx context.Context, remote, local string) error {
	scp := scp.New(c.Client, afero.NewOsFs(), nil)
	return scp.Download(ctx, remote, local)
}

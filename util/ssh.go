package util

import (
	"golang.org/x/crypto/ssh"
	"strconv"
)

type Client struct {
	hostAlias string
	user      string
	pass      string
	port      int
	Client    *ssh.Client
}

func (c *Client) NewClient() *ssh.Client {
	hostFullPath := c.GetFullHostPath()
	pass := []ssh.AuthMethod{ssh.Password(c.pass)}
	conf := &ssh.ClientConfig{User: c.user, Auth: pass}
	client, err := ssh.Dial("tcp", hostFullPath, conf)
	if err != nil {
		panic(err)
	}
	c.Client = client
	return client
}

// 获取远程服务器完整路径
func (c *Client) GetFullHostPath() string {
	host, ok := HostsConfig[c.hostAlias]["host"]
	if !ok {
		panic("未找到该服务器别名：" + c.hostAlias)
	}

	return host.(string) + ":" + strconv.Itoa(c.port)
}

func (c *Client) SetAlias(alias string) {
	c.hostAlias = alias
}

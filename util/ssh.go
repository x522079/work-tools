package util

import (
	"bufio"
	"cmds/middleware"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Client struct {
	HostAlias string
	User      string
	Pass      string
	Port      int
	KeyPath   string
	Client    *ssh.Client
}

func (c *Client) NewClient() *ssh.Client {
	hostFullPath := c.GetFullHostPath(true)
	//pass := []ssh.AuthMethod{ssh.Password(c.Pass)}
	pass := []ssh.AuthMethod{c.getPublicKeyMethod()}
	conf := &ssh.ClientConfig{User: c.User, Auth: pass, HostKeyCallback: ssh.FixedHostKey(c.getHostKey())}
	client, err := ssh.Dial("tcp", hostFullPath, conf)
	if err != nil {
		panic(err)
	}
	c.Client = client
	return client
}

func (c *Client) NewSftpClient() *sftp.Client {
	hostFullPath := c.GetFullHostPath(true)
	pass := []ssh.AuthMethod{c.getPublicKeyMethod()}
	conf := &ssh.ClientConfig{User: c.User, Auth: pass, HostKeyCallback: ssh.FixedHostKey(c.getHostKey())}
	client, err := ssh.Dial("tcp", hostFullPath, conf)
	if err != nil {
		panic(err)
	}
	c.Client = client

	sftpClient, err := sftp.NewClient(client, func(client *sftp.Client) error {
		return nil
	})

	return sftpClient
}

// 获取远程服务器完整路径
func (c *Client) GetFullHostPath(hasPort bool) string {
	host, ok := middleware.HostsConfig[c.HostAlias]["host"]
	if !ok {
		panic("未找到该服务器别名：" + c.HostAlias)
	}

	if hasPort {
		return host.(string) + ":" + strconv.Itoa(c.Port)
	}
	return host.(string)
}

func (c *Client) SetAlias(alias string) {
	c.HostAlias = alias
}

// 获取主机key
func (c *Client) getHostKey() ssh.PublicKey {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts"))
	CheckError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fileds := strings.Split(scanner.Text(), " ")
		if len(fileds) != 3 {
			continue
		}

		if strings.Contains(fileds[0], middleware.HostsConfig[c.HostAlias]["host"].(string)) {
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			CheckError(err)
			break
		}
	}

	return hostKey
}

func (c *Client) getPublicKeyMethod() ssh.AuthMethod {
	file, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	CheckError(err)

	signer, err := ssh.ParsePrivateKey(file)
	CheckError(err)

	return ssh.PublicKeys(signer)
}

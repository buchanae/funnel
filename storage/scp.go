package storage

import (
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"github.com/pkg/sftp"
)

type SCPClient struct {
	Host   string
	Config *ssh.ClientConfig
	Conn   *ssh.Client
	Client *sftp.Client
}

func NewSCPClient(host string) *SCPClient {
	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}
	username := os.Getenv("USER")
	config := &ssh.ClientConfig{
		User: username,
		Auth: auths,
	}

 	// Create a new SCP client
	return &SCPClient{
		Host:   host,
		Config: config,
	}
}

func (s *SCPClient) Connect() error {
	var err error
	s.Conn, err = ssh.Dial("tcp", s.Host, s.Config)
	if err != nil {
		return err
	}
	return nil
	s.Client, err = sftp.NewClient(s.Conn)
	if err != nil {
		return fmt.Errorf("unable to start sftp subsytem: %v", err)
	}
	return nil
}

func (s *SCPClient) Close() error {
	var err error
	err = s.Conn.Close()
	if err != nil {
		return err
	}
	err = s.Client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *SCPClient) SCPLocalToRemote(source string, dest string) error {
	// Connect to the remote server
	err := s.Connect()
	if err != nil{
		return fmt.Errorf("Couldn't establish a connection to the remote server: %v", err)
	}
	// Close connection after the file has been copied
	defer s.Close()

	sf, _ := os.Open(source)
	defer sf.Close()

	dstD := path.Dir(dest)
	if _, err := s.Client.Stat(dstD); err != nil {
		s.Client.MkdirAll(dstD, 0777)
	}
	df, err := s.Client.Create(dest)
	if err != nil {
		return err
	}
	
	_, err = io.Copy(df, sf)
	cerr := df.Close()
	if err != nil {
		return err
	}
	if cerr != nil {
		return cerr
	}
	return nil
}

func (s *SCPClient) SCPRemoteToLocal(source string, dest string) error {
	// Connect to the remote server
	err := s.Connect()
	if err != nil{
		return fmt.Errorf("Couldn't establish a connection to the remote server: %v ", err)
	}
	// Close connection after the file has been copied
	defer s.Close()

	sf, _ := s.Client.Open(source)
	defer sf.Close()

	dstD := path.Dir(dest)
	if _, err := os.Stat(dstD); err != nil {
		os.MkdirAll(dstD, 0777)
	}
	df, err := os.Create(dest)
	if err != nil {
		return err
	}

	_, err = io.Copy(df, sf)
	cerr := df.Close()
	if err != nil {
		return err
	}
	if cerr != nil {
		return cerr
	}
	return nil
}

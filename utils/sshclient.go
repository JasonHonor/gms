package utils

import (
	"encoding/hex"
	//	"encoding/hex"
	//	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
	//	"encoding/hex"
)

// SSHClient SSH client wrapper
type SSHClient struct {
	Host     string
	Port     int
	Username string
	Password string

	MoreTag    string
	MoreWant   string
	IsMoreLine bool

	// color tag to ignore
	ColorTag string

	ReadOnlyPrompt  string
	SysEnablePrompt string
	LineBreak       string

	ExitCmd string
}

// connect connect to ssh server.
func (sshClient *SSHClient) connect() (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)

	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(sshClient.Password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User:            sshClient.Username,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	clientConfig.Ciphers = append(clientConfig.Ciphers, "aes128-cbc")
	clientConfig.KeyExchanges = append(clientConfig.KeyExchanges, "diffie-hellman-group1-sha1")

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", sshClient.Host, sshClient.Port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

// Execute run ssh command.
func (sshClient *SSHClient) Execute(cmdList []string) []string {

	session, err := sshClient.connect()
	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	log.Println("Connected.")

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,            // 0:disable echo
		ssh.TTY_OP_ISPEED: 14400 * 1024, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400 * 1024, //output speed = 14.4kbaud
	}
	if err1 := session.RequestPty("linux", 64, 200, modes); err1 != nil {
		log.Fatalf("request pty error: %s\n", err1.Error())
	}

	w, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}
	r, err := session.StdoutPipe()
	if err != nil {
		panic(err)
	}
	e, err := session.StderrPipe()
	if err != nil {
		panic(err)
	}

	in, out := sshClient.MuxShell(w, r, e)
	if err := session.Shell(); err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", <-out) //ignore the shell output

	var ret []string

	for _, cmd := range cmdList {
		log.Printf("Exec %v\n", cmd)

		in <- cmd

		sOut := <-out
		//log.Printf("%s\n", sOut)
		ret = append(ret, sOut)
	}

	in <- sshClient.ExitCmd
	_ = <-out
	session.Wait()

	return ret
}

// MuxShell interaction with shell.
func (sshClient *SSHClient) MuxShell(w io.Writer, r, e io.Reader) (chan<- string, <-chan string) {
	in := make(chan string, 5)
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1) //for the shell itself
	go func() {
		for cmd := range in {
			wg.Add(1)
			w.Write([]byte(cmd + "\n"))
			wg.Wait()
		}
	}()

	go func() {
		var (
			buf [2 * 1024 * 1024]byte
			t   int
		)
		for {

			//read next buf n.
			n, err := r.Read(buf[t:])
			if err != nil {
				fmt.Println("ReadError:" + err.Error())
				close(in)
				close(out)
				return
			}

			t += n
			result := string(buf[:t])

			line := string(buf[t-n : t])
			//fmt.Printf("Line:=>%v\n", line)

			if strings.Contains(line, sshClient.MoreTag) {
				if sshClient.IsMoreLine {
					t -= n
				} else {
					t -= len(sshClient.MoreTag)
				}
				w.Write([]byte(sshClient.MoreWant))
			} else if len(sshClient.ColorTag) > 0 {

				//invisible char
				if strings.HasSuffix(sshClient.ColorTag, "H") {
					colorTag := strings.Replace(sshClient.ColorTag, "H", "", -1)
					tag, err := hex.DecodeString(colorTag)
					if err == nil {
						// remove colortag
						var newBuf []byte
						newBuf = RemoveStrByTagBytes([]byte(line), tag, tag, "\r\n")

						t -= n
						copy(buf[t:], newBuf[:])
						t += len(newBuf)
					}
				} else {
					// remove colortag
					var newBuf []byte
					newBuf = RemoveStringByTag([]byte(line), sshClient.ColorTag, sshClient.ColorTag, "\r\n")

					t -= n
					copy(buf[t:], newBuf[:])
					t += len(newBuf)
				}
			}

			if strings.Contains(result, "username:") ||
				strings.Contains(result, "password:") ||
				strings.Contains(result, sshClient.ReadOnlyPrompt) ||
				strings.Contains(result, sshClient.SysEnablePrompt) {

				//sOut := string(buf[:t])
				//fmt.Printf("DataOut:%v\n", sOut)

				out <- string(buf[:t])
				t = 0
				wg.Done()
			}
		}
	}()
	return in, out
}

// Shell interaction shell.
func (sshClient *SSHClient) Shell(session *ssh.Session) {
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
	}
	if err1 := session.RequestPty("linux", 32, 160, modes); err1 != nil {
		log.Fatalf("request pty error: %s\n", err1.Error())
	}
	if err2 := session.Shell(); err2 != nil {
		log.Fatalf("start shell error: %s\n", err2.Error())
	}
	if err3 := session.Wait(); err3 != nil {
		log.Fatalf("return error: %s\n", err3.Error())
	}
}

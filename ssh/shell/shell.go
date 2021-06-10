package shell

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type Shell struct {
	client *ssh.Client

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New returns a new Shell
func New(c *ssh.Client) *Shell {
	return &Shell{
		client: c,
	}
}

func (s *Shell) Run(stdin io.Reader, stdout, stderr io.Writer) error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("ssh: failed to create new session")
	}
	defer session.Close()

	s.stdin = stdin
	s.stdout = stdout
	s.stderr = stderr

	return s.runShell(session)
}

func (s *Shell) runShell(session *ssh.Session) error {
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("failed to make raw terminal: %v", err)
	}
	defer terminal.Restore(fd, state) //nolint

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	err = session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
	if err != nil {
		return err
	}

	go s.syncWindowChange(fd, session)

	session.Stdin = s.stdin
	session.Stdout = s.stdout
	session.Stderr = s.stderr

	err = session.Shell()
	if err != nil {
		return err
	}
	err = session.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (s *Shell) syncWindowChange(fd int, session *ssh.Session) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGWINCH)
	width, height, err := terminal.GetSize(fd)
	if err != nil {
		// TODO:
		return
	}

	for sig := range signalCh {
		if sig == nil {
			return
		}
		curWidth, curHeight, _ := terminal.GetSize(fd)

		if curWidth == width && curHeight == height {
			continue
		}
		err := session.WindowChange(curHeight, curWidth)
		if err != nil {
			// closed
			return
		}
		width, height = curWidth, curHeight
	}
}

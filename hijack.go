package dcgi

import (
	"fmt"
	"io"
	"log"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/engine-api/types"
)

// holdHijackedConnection handles copying input to and output from streams to
// the connection. Copied from github.com/docker/docker/api/client.
func holdHijackedConnection(inputStream io.ReadCloser, outputStream, errorStream io.Writer, resp types.HijackedResponse) error {
	var err error

	receiveStdout := make(chan error, 1)
	if outputStream != nil || errorStream != nil {
		go func() {
			_, err = stdcopy.StdCopy(outputStream, errorStream, resp.Reader)
			// log.Printf("[hijack] End of stdout")
			receiveStdout <- err
		}()
	}

	stdinDone := make(chan struct{})
	go func() {
		if inputStream != nil {
			io.Copy(resp.Conn, inputStream)
			// log.Printf("[hijack] End of stdin")
		}

		if err := resp.CloseWrite(); err != nil {
			log.Printf("cgi: couldn't send EOF: %s", err)
		}
		close(stdinDone)
	}()

	select {
	case err := <-receiveStdout:
		if err != nil {
			return fmt.Errorf("Error receiveStdout: %s", err)
		}
	case <-stdinDone:
		if outputStream != nil || errorStream != nil {
			err := <-receiveStdout
			if err != nil {
				return fmt.Errorf("Error receiveStdout: %s", err)
			}
		}
	}

	return nil
}

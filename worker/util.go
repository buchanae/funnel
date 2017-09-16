package worker

import (
	"context"
	"net"
	"os/exec"
	"syscall"
  "github.com/ohsu-comp-bio/funnel/proto/tes"
  "time"
)

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down

		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface

		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err

		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", nil
}

// GetExitCode gets the exit status (i.e. exit code) from the result of an executed command.
// The exit code is zero if the command completed without error.
func GetExitCode(err error) int {
	if err != nil {
		if exiterr, exitOk := err.(*exec.ExitError); exitOk {
			if status, statusOk := exiterr.Sys().(syscall.WaitStatus); statusOk {
				return status.ExitStatus()
			}
		}
	}
	// The error is nil, the command returned successfully, so exit status is 0.
	return 0
}

func PollForCancel(ctx context.Context, get func() tes.State, rate time.Duration) context.Context {
	taskctx, cancel := context.WithCancel(ctx)

	// Start a goroutine that polls the server to watch for a canceled state.
	// If a cancel state is found, "taskctx" is canceled.
	go func() {
		ticker := time.NewTicker(rate)
		defer ticker.Stop()

		for {
			select {
			case <-taskctx.Done():
				return
			case <-ticker.C:
				state := get()
				if tes.TerminalState(state) {
					cancel()
				}
			}
		}
	}()
	return taskctx
}

func LogHostIP(tl TaskLogger, i int) {
	// Grab the IP address of this host. Used to send task metadata updates.
	ip, err := externalIP()
  if err == nil {
    tl.ExecutorHostIP(i, ip)
  }
}

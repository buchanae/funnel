package worker

import (
	"github.com/ohsu-comp-bio/funnel/config"
	pbf "github.com/ohsu-comp-bio/funnel/proto/funnel"
	pscpu "github.com/shirou/gopsutil/cpu"
	psmem "github.com/shirou/gopsutil/mem"
	"net"
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

// detectResources helps determine the amount of resources to report.
// Resources are determined by inspecting the host, but they
// can be overridden by config.
func detectResources(conf config.Resources) config.Resources {
	res := config.Resources{
		Cpus:   conf.Cpus,
		RamGb:  conf.RamGb,
		DiskGb: conf.DiskGb,
	}
	cpuinfo, _ := pscpu.Info()
	vmeminfo, _ := psmem.VirtualMemory()

	if conf.Cpus == 0 {
		// TODO is cores the best metric? with hyperthreading,
		//      runtime.NumCPU() and pscpu.Counts() return 8
		//      on my 4-core mac laptop
		for _, cpu := range cpuinfo {
			res.Cpus += uint32(cpu.Cores)
		}
	}

	if conf.RamGb == 0.0 {
		res.RamGb = float64(vmeminfo.Total) /
			float64(1024) / float64(1024) / float64(1024)
	}

	return res
}

// NoopTaskRunner is useful during testing for creating a worker with a TaskRunner
// that doesn't do anything.
func NoopTaskRunner(l TaskControl, c config.Worker, j *pbf.TaskWrapper) {
}

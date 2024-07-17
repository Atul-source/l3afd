package restart

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/l3af-project/l3afd/v2/apis"
	"github.com/l3af-project/l3afd/v2/bpfprogs"
	"github.com/l3af-project/l3afd/v2/config"
	"github.com/l3af-project/l3afd/v2/pidfile"
	"github.com/l3af-project/l3afd/v2/stats"
)

func CheckKernelVersion(conf *config.Config) error {
	const minVerLen = 2

	kernelVersion, err := GetKernelVersion()
	if err != nil {
		return fmt.Errorf("failed to find kernel version: %v", err)
	}

	//validate version
	ver := strings.Split(kernelVersion, ".")
	if len(ver) < minVerLen {
		return fmt.Errorf("expected minimum kernel version length %d and got %d, ver %+q", minVerLen, len(ver), ver)
	}
	major_ver, err := strconv.Atoi(ver[0])
	if err != nil {
		return fmt.Errorf("failed to find kernel major version: %v", err)
	}
	minor_ver, err := strconv.Atoi(ver[1])
	if err != nil {
		return fmt.Errorf("failed to find kernel minor version: %v", err)
	}

	if major_ver > conf.MinKernelMajorVer {
		return nil
	}
	if major_ver == conf.MinKernelMajorVer && minor_ver >= conf.MinKernelMinorVer {
		return nil
	}

	return fmt.Errorf("expected Kernel version >=  %d.%d", conf.MinKernelMajorVer, conf.MinKernelMinorVer)
}

func GetKernelVersion() (string, error) {
	osVersion, err := os.ReadFile("/proc/version")
	if err != nil {
		return "", fmt.Errorf("failed to read procfs: %v", err)
	}
	var u1, u2, kernelVersion string
	_, err = fmt.Sscanf(string(osVersion), "%s %s %s", &u1, &u2, &kernelVersion)
	if err != nil {
		return "", fmt.Errorf("failed to scan procfs version: %v", err)
	}

	return kernelVersion, nil
}

func DoServerWork() {

	// need to setup a goroutine with
	// who is serving (l3afd state structure) & server file discriptors
	//
	// TODO is setting up servers which are res
}

func DoClientWork(conf *config.Config) {
	// Create a process for l3afd
	if err := pidfile.CreatePID(conf.PIDFilename + "-child"); err != nil {
		log.Fatal().Err(err).Msgf("The PID file: %s, could not be created", conf.PIDFilename)
	}
	if runtime.GOOS == "linux" {
		if err := CheckKernelVersion(conf); err != nil {
			log.Fatal().Err(err).Msg("The unsupported kernel version please upgrade")
		}
	}
	// Get Hostname
	machineHostname, err := os.Hostname()
	if err != nil {
		log.Error().Err(err).Msg("Could not get hostname from OS")
	}

	// setup Metrics endpoint
	// We need to implement a way so that
	stats.SetupMetrics(machineHostname, daemonName, conf.MetricsAddr)

	pMon := bpfprogs.NewpCheck(conf.MaxEBPFReStartCount, conf.BpfChainingEnabled, conf.EBPFPollInterval)
	bpfM := bpfprogs.NewpBPFMetrics(conf.BpfChainingEnabled, conf.NMetricSamples)

	nfConfigs, err := bpfprogs.NewNFConfigs(ctx, machineHostname, conf, pMon, bpfM)
	if err != nil {
		return nil, fmt.Errorf("error in NewNFConfigs setup: %v", err)
	}

	if err := apis.StartConfigWatcher(ctx, machineHostname, daemonName, conf, nfConfigs); err != nil {
		return nil, fmt.Errorf("error in version announcer: %v", err)
	}

	return nfConfigs, nil
}

// For restart purpose we just stop l3afd server
// we have 7080 8898 8899 we need to stop this go routines
// step 1 : we will send restart API request
// step 2 : we will kill all the http server so that our l3afd state is fixed  we have -> then we will serve l3afd state
// step 3 : Invoke new l3afd process with child argument
// step 4 : Get Full l3afd state
// step 5 : send Ready to old process
// step 6 :

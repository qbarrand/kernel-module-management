package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/kubernetes-sigs/kernel-module-management/internal/cmd"
	"github.com/kubernetes-sigs/kernel-module-management/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/klog/v2/klogr"
	ctrl "sigs.k8s.io/controller-runtime"
)

func uname() (string, error) {
	var utsBuf syscall.Utsname
	if err := syscall.Uname(&utsBuf); err != nil {
		return "", fmt.Errorf("error calling the uname syscall: %v", err)
	}

	var sb strings.Builder

	for _, c := range utsBuf.Release {
		// the buffer is padded with zeroes
		if c == 0 {
			break
		}

		sb.WriteByte(byte(c))
	}

	return sb.String(), nil
}

func setFirmwareLookupPath(path string) error {
	const sysFile = "/sys/module/firmware_class/parameters/path"

	fd, err := os.Create(sysFile)
	if err != nil {
		return fmt.Errorf("could not open %s: %v", sysFile, err)
	}
	defer fd.Close()

	_, err = fd.WriteString(path)

	return err
}

func main() {
	logger := klogr.New().WithName("kmm-daemon")

	serverAddr := cmd.GetEnvOrFatalError("GRPC_SERVER_ADDR", logger)
	hostName := cmd.GetEnvOrFatalError("HOSTNAME", logger)

	kernelVersion, err := uname()
	if err != nil {
		cmd.FatalError(logger, err, "Could not get the kernel version")
	}

	if firmwareLookupPath := os.Getenv("FIRMWARE_LOOKUP_PATH"); firmwareLookupPath != "" {
		logger.Info("Setting the firmware lookup path", "path", firmwareLookupPath)

		if err = setFirmwareLookupPath(firmwareLookupPath); err != nil {
			cmd.FatalError(logger, err, "Could not set the firmware lookup path")
		}
	}

	creds := grpc.WithTransportCredentials(insecure.NewCredentials())

	ctx := ctrl.SetupSignalHandler()

	conn, err := grpc.DialContext(ctx, serverAddr, creds)
	if err != nil {
		cmd.FatalError(logger, err, "Could not connect to the server", "address", serverAddr)
	}
	defer conn.Close()

	client := proto.NewKMMClient(conn)

	dcr := &proto.DaemonConfigRequest{
		KernelVersion: kernelVersion,
		HostName:      hostName,
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)

	cfg, err := client.GetDaemonConfig(ctx, dcr)
	if err != nil {
		cmd.FatalError(logger, err, "could not get the daemon configuration")
	}

	cancel()

	_ = cfg
}

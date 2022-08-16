/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// partitionCmd represents the partition command
var partitionCmd = &cobra.Command{
	Use:   "partition",
	Short: "Partitions and formats a disk",
	Run: func(cmd *cobra.Command, args []string) {
		device, _ := cmd.Flags().GetString("device")
		size, _ := cmd.Flags().GetInt("size")
		partition(device, size)
	},
}

func init() {
	partitionCmd.Flags().StringP("device", "d", "", "Device to be partitioned")
	partitionCmd.MarkFlagRequired("device")
	partitionCmd.Flags().IntP("size", "s", 100, "Partition size in GB")
	rootCmd.AddCommand(partitionCmd)

}

func isPartitionSizeTooBig(deviceSize, desiredSize float64) bool {
	return desiredSize > deviceSize
}

func generateGetDeviceSizeCommand(device string) *exec.Cmd {
	return exec.Command("lsblk", device, "-osize", "-dn")
}

func generatePartitionCommand(device string, size int) *exec.Cmd {
	return exec.Command("sgdisk", "-n", fmt.Sprintf("1:-%dGiB:0", size), device, "-g", "-c:1:data")
}

func generateFormatCommand(device string) *exec.Cmd {
	return exec.Command("mkfs.xfs", "-f", device+"1")
}

func partition(device string, size int) {
	cmd := generateGetDeviceSizeCommand(device)
	stdout, err := executeCommand(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to get devize size: %s\n", string(stdout))
		os.Exit(1)
	}

	deviceSizeStr := strings.TrimSpace(string(stdout))
	deviceSizeStr = deviceSizeStr[:len(deviceSizeStr)-1]
	deviceSize, err := strconv.ParseFloat(deviceSizeStr, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to parse device size: %v\n", err)
		os.Exit(1)
	}
	if isPartitionSizeTooBig(deviceSize, float64(size)) {
		fmt.Fprintf(os.Stderr, "error: partition size is too big")
		os.Exit(1)
	}

	cmd = generatePartitionCommand(device, size)
	stdout, err = executeCommand(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to partition device: %s\n", string(stdout))
		os.Exit(1)
	}

	cmd = generateFormatCommand(device)
	stdout, err = executeCommand(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: unable to format device: %s\n", string(stdout))
		os.Exit(1)
	}

}

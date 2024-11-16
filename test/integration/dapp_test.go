package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}

type AppSuite struct {
	suite.Suite
	processes []*exec.Cmd
}

func (s *AppSuite) SetupTest() {
	s.startProcess("nonodo", []string{}, "INF nonodo: ready", "nonodo server")
	s.startProcess("cartesi-machine", []string{
		"--network",
		"--flash-drive=label:root,filename:.cartesi/image.ext2",
		"--env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004",
		"--",
		"/var/opt/cartesi-app/app",
	}, "", "Cartesi machine")
}

func (s *AppSuite) TearDownTest() {
	for _, cmd := range s.processes {
		s.terminateProcessGroup(cmd)
	}
	s.terminateProcessesByName("anvil")
}

func (s *AppSuite) TestItCreateCrowdfundingAndFinishCrowdfundingWithoutPartialSellingAndPayingAllInvestor() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	fmt.Printf("Admin Address: %s\n", admin.Hex())
	fmt.Printf("Creator Address: %s\n", creator.Hex())
}

func (s *AppSuite) startProcess(name string, args []string, readyLog, description string) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Stderr = os.Stderr
	s.processes = append(s.processes, cmd)

	var stdoutPipe io.ReadCloser
	if readyLog != "" {
		var err error
		stdoutPipe, err = cmd.StdoutPipe()
		s.Require().NoError(err, "Failed to create stdout pipe for "+description)
		go s.waitForLog(stdoutPipe, readyLog, description)
	} else {
		cmd.Stdout = os.Stdout
	}

	err := cmd.Start()
	s.Require().NoError(err, "Failed to start "+description)
}

func (s *AppSuite) waitForLog(pipe io.ReadCloser, readyLog, description string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		log := scanner.Text()
		slog.Info(log)
		if strings.Contains(log, readyLog) {
			fmt.Println(description + " is ready!")
			return
		}
	}
}

func (s *AppSuite) terminateProcessGroup(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			s.Require().NoError(err, "Failed to terminate process group for "+cmd.Path)
		}
	}
}

func (s *AppSuite) terminateProcessesByName(name string) {
	output, err := exec.Command("pgrep", "-f", name).Output()
	if err == nil {
		for _, pid := range strings.Fields(string(output)) {
			fmt.Printf("Killing %s process with PID: %s\n", name, pid)
			exec.Command("kill", "-9", pid).Run()
		}
	}
}

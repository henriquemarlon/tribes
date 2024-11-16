package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"
)

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppSuite))
}

type AppSuite struct {
	suite.Suite
	dappCmd   *exec.Cmd
	nonodoCmd *exec.Cmd
}

func (s *AppSuite) SetupTest() {
	// Start the nonodo server
	s.nonodoCmd = exec.Command("nonodo")
	stdoutPipe, err := s.nonodoCmd.StdoutPipe()
	s.Require().NoError(err, "Failed to create stdout pipe for nonodo")
	s.nonodoCmd.Stderr = os.Stderr

	err = s.nonodoCmd.Start()
	s.Require().NoError(err, "Failed to start nonodo server")

	ready := make(chan bool)
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			out := scanner.Text()
			slog.Info(out)
			if strings.Contains(out, "INF nonodo: ready") {
				ready <- true
				return
			}
		}
	}()

	select {
	case <-ready:
		fmt.Println("Nonodo is ready!")
	case <-time.After(5 * time.Second):
		s.Require().Fail("Timed out waiting for nonodo to be ready")
	}

	s.dappCmd = exec.Command(
		"cartesi-machine",
		"--network",
		"--flash-drive=label:root,filename:.cartesi/image.ext2",
		"--env=ROLLUP_HTTP_SERVER_URL=http://10.0.2.2:5004",
		"--",
		"/var/opt/cartesi-app/app",
	)
	s.dappCmd.Stdout = os.Stdout
	s.dappCmd.Stderr = os.Stderr
	err = s.dappCmd.Start()
	s.Require().NoError(err, "Failed to start cartesi-machine")
}

func (s *AppSuite) TearDownTest() {
	if s.dappCmd != nil && s.dappCmd.Process != nil {
		err := s.dappCmd.Process.Kill()
		s.Require().NoError(err, "Failed to kill Cartesi machine process")
	}
	if s.nonodoCmd != nil && s.nonodoCmd.Process != nil {
		err := s.nonodoCmd.Process.Kill()
		s.Require().NoError(err, "Failed to kill nonodo server process")
	}
}

func (s *AppSuite) TestItCreateCrowdfundingAndFinishCrowdfundingWithoutPartialSellingAndPayingAllInvestor() {
	admin := common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	creator := common.HexToAddress("0x0000000000000000000000000000000000000007")
	fmt.Printf("Admin Address: %s\n", admin.Hex())
	fmt.Printf("Creator Address: %s\n", creator.Hex())
}

package service

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/MoTeam-org/OpenBMCLAPI-API-Go/utils"
)

type CommonService struct{}

func NewCommon() *CommonService {
	return &CommonService{}
}

// ClearScreen 清屏函数
func (s *CommonService) ClearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// WaitForEnter 等待用户按回车继续，然后清屏
func (s *CommonService) WaitForEnter() {
	fmt.Print(utils.ColorText(utils.Yellow, "\n按回车键继续..."))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	s.ClearScreen()
}

// WaitForEnterWithoutClear 等待用户按回车继续，但不清屏
func (s *CommonService) WaitForEnterWithoutClear() {
	fmt.Print(utils.ColorText(utils.Yellow, "\n按回车键继续..."))
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

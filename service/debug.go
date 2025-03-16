package service

import "github.com/MoTeam-org/OpenBMCLAPI-API-Go/utils"

// SetDebugLevel 设置调试级别
func SetDebugLevel(level int) {
	utils.SetDebugLevel(level)
}

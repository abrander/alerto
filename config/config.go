package config

import (
	"os"

	"github.com/abrander/alerto/logger"
)

const (
	ConfigDir = "/etc/alerto"
)

func init() {
	_, err := os.Stat("/etc/alerto")
	if err != nil {
		uid := os.Getuid()
		gid := os.Getgid()
		logger.Error("config", "Please run:\nsudo mkdir -p %s && sudo chown %d.%d %s\n", ConfigDir, uid, gid, ConfigDir)
	}
}

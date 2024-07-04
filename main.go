package main

import (
	"github.com/sirupsen/logrus"
	"github.com/xrcuo/api_boot/psutil"
)

func main() {
	logrus.Info("CPU: ", psutil.CpuPercent())
	logrus.Info("Memory: ", psutil.DiskPercent())
	logrus.Info("Disk: ", psutil.MemPercent())

}

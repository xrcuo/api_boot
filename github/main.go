package main

import (
	"time"
	con "wso/config"
	"wso/downloaddir"
	"wso/unzip"

	"github.com/nyancatda/AyaLog"
)

var conf *con.Config

func main() {
	// 初始化本地最新版本为空
	localLatestVersion := ""

	for {
		latestRelease, err := downloaddir.GetLatestRelease(conf.Repoowner, conf.Reponame)
		if err != nil {
			AyaLog.Info("获取最新版本信息失败:", err)
		} else {
			latestVersion := latestRelease.GetTagName()
			AyaLog.Info("最新版本:", latestVersion)

			// 比较本地版本和远程版本
			if localLatestVersion != latestVersion {
				AyaLog.Info("发现新版本，开始下载...")

				// 下载版本文件
				//downloadRelease(latestRelease)
				downloaddir.DownloadRelease(latestRelease)
				AyaLog.Info("下载完成，开始解压...")

				unzip.ExtractRelease()

				// 更新本地最新版本
				localLatestVersion = latestVersion
			} else {
				AyaLog.Info("已经是最新版本")
			}
		}

		time.Sleep(time.Minute * time.Duration(con.Conf.Checkintervainin))
	}
}

func init() {
	// 加载配置文件
	con.Parse()
	AyaLog.Info("System", "config.yaml 加载成功")
	conf = con.Conf

}

package downloaddir

import (
	"context"
	"fmt"
	"os"
	con "wso/config"

	"github.com/google/go-github/v64/github"
	"github.com/nyancatda/AyaLog"
	"github.com/xrcuo/Server-Motd-Webui/pb_download"
)

var (
	ws = "https://mirror.ghproxy.com/"
)

// 获取最新版本信息
func GetLatestRelease(owner, repo string) (*github.RepositoryRelease, error) {
	client := github.NewClient(nil)
	release, _, err := client.Repositories.GetLatestRelease(context.Background(), owner, repo)
	return release, err
}

// 下载版本文件
func DownloadRelease(release *github.RepositoryRelease) error {
	for _, asset := range release.Assets {
		if asset.GetName() == con.Conf.Releaseasset {
			downloadURL := asset.GetBrowserDownloadURL()

			// 创建下载目录
			os.MkdirAll(con.Conf.Downloaddir, os.ModePerm)

			AyaLog.Info("开始下载文件")

			downloader := pb_download.NewDownloader(con.Conf.Downloaddir)
			downloader.AppendResource(con.Conf.Releaseasset, ws+downloadURL)
			// 可自主调整协程数量，默认为CPU核数
			downloader.Concurrent = 1
			err := downloader.Start()
			if err != nil {
				panic(err)
			}

			return nil
		}
	}
	return fmt.Errorf("未找到名为 %s 的发布资源", con.Conf.Releaseasset)
}

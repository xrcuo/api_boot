package unzip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	con "wso/config"

	"github.com/nyancatda/AyaLog"
)

func ExtractRelease() {

	errs := os.RemoveAll(con.Conf.Extractdir) // 删除目录
	if errs != nil {
		AyaLog.Info("删除失败")
	} else {
		AyaLog.Info("删除成功")
	}

	archive, err := zip.OpenReader(con.Conf.Downloaddir + "/" + con.Conf.Releaseasset)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	// 2、循环访问 zip 中的文件 zip.File 切片
	for _, f := range archive.File {
		filePath := filepath.Join(con.Conf.Extractdir, f.Name)
		AyaLog.Info("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(con.Conf.Extractdir)+string(os.PathSeparator)) {
			AyaLog.Info("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			AyaLog.Info("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		// 3、使用 zip.File.Open 方法读取 zip 中文件的内容
		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		// 4、使用 io.Copy 或 io.Writer.Write 保存解压后的文件内容
		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		// 5、使用 zip.Reader.Close 关闭 zip 文件
		dstFile.Close()
		fileInArchive.Close()
	}
}

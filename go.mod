module github.com/xrcuo/api_boot

go 1.22

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.29.10-simp

replace modernc.org/libc => github.com/fumiama/libc v0.0.0-20240530081950-6f6d8586b5c5

require github.com/sirupsen/logrus v1.9.3

require golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect

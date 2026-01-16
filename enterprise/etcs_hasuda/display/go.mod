module enterprise/etcs/display

go 1.20

replace etc-system.jp/iniread => ./iniread

replace etc-system.jp/tcpclient => ./tcpclient

require (
	etc-system.jp/iniread v0.0.0-00010101000000-000000000000
	github.com/pkg/sftp v1.13.5
	github.com/zserge/lorca v0.1.10
	golang.org/x/crypto v0.9.0
)

require (
	github.com/kr/fs v0.1.0 // indirect
	github.com/stretchr/testify v1.8.3 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

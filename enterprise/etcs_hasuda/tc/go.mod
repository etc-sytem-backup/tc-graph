module localhost.com/traffic_counter

go 1.23.0

toolchain go1.24.4

replace localhost.com/conretry => ./conretry

replace localhost.com/tcpclient => ./tcpclient

replace localhost.com/makecmd => ./makecmd

require (
	etc-system.jp/iniread v0.0.0-00010101000000-000000000000
	localhost.com/conretry v0.0.0-00010101000000-000000000000
	localhost.com/findcmd v0.0.0-00010101000000-000000000000
	localhost.com/makecmd v0.0.0-00010101000000-000000000000
	localhost.com/tcpclient v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace localhost.com/findcmd => ./findcmd

replace etc-system.jp/iniread => ./iniread

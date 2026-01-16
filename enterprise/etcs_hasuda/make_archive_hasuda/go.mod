module enterprise/etcs/etcs_hasuda/make_archive_hasuda

go 1.21.1

replace localhost.com/iniread => ./iniread

require localhost.com/iniread v0.0.0-00010101000000-000000000000

require (
	github.com/stretchr/testify v1.8.4 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

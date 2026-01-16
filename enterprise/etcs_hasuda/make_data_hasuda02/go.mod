module etc-system.jp/make_data_hasuda02

go 1.20

replace localhost.com/iniread => ./iniread

require (
	localhost.com/iniread v0.0.0-00010101000000-000000000000
	localhost.com/readcsv v0.0.0-00010101000000-000000000000
)

require (
	github.com/stretchr/testify v1.7.1 // indirect
	gopkg.in/ini.v1 v1.66.4 // indirect
)

replace localhost.com/readcsv => ./readcsv

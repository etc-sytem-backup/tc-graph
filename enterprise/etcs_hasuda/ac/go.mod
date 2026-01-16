module localhost.com/application_counter

go 1.17

require (
	gopkg.in/ini.v1 v1.66.4 // indirect
	localhost.com/iniread v0.0.0-00010101000000-000000000000
)

require github.com/stretchr/testify v1.7.1 // indirect

replace localhost.com/iniread => ./iniread

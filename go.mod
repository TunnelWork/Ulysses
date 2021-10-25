module github.com/TunnelWork/Ulysses

go 1.16

replace github.com/TunnelWork/Ulysses/src => ./src

require (
	github.com/TunnelWork/Ulysses.Lib v0.0.1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-sql-driver/mysql v1.6.0
	gopkg.in/yaml.v2 v2.4.0
)

module gms

go 1.15

require (
	github.com/AlekSi/zabbix v0.0.0-00010101000000-000000000000
	github.com/gogf/gf v1.13.7
	github.com/soniah/gosnmp v1.27.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
)

replace github.com/AlekSi/zabbix => ../zabbix

replace gonuts.io/aleksi/reflector v0.4.1 => github.com/AlekSi/reflector v0.4.1

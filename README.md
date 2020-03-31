# lseed

Utility to seed an LDAP instance with data.

## Help

```
go run lseed.go --help
  -groups int
    	the number of groups (default 2)
  -help
    	show help
  -host string
    	the bind host (default "0.0.0.0")
  -members int
    	the number of members per group (default 10)
  -ou string
    	the organizational unit that will contain the seeded data (default "ou=loadtest,dc=mm,dc=test,dc=com")
  -password string
    	the bind password (default "mostest")
  -photo string
    	the path to the profile photo
  -port int
    	the bind port (default 389)
  -user string
    	the bind user (default "cn=admin,dc=mm,dc=test,dc=com")
```

## Examples
```
go run lseed.go
```
```
go run lseed.go -photo ~/Pictures/test.jpeg
```

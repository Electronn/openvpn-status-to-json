It is a converter from standard OpenVPN status log to json structure.

Compilation :

1. clone repo
2. go build converter.go

Usage :
```
go run converter.go -ovpn.log /var/log/status.log
```

or if you compile it 

```
converter -ovpn.log /var/log/status.log
```


Enjoy!

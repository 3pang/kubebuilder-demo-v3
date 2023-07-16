@echo off
go env -w GOOS=linux
go build -o manager cmd/main.go
scp manager root@192.168.56.101:/root/wy/kubebuilder-demo-v3/

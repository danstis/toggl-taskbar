go build -ldflags "-H=windowsgui -X main.Version=$(gitversion /output json /showvariable SemVer)"

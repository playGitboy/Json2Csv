CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Json2Csv.win.amd64.exe main.go
./upx -9 Json2Csv.win.amd64.exe
tar -czvf Json2Csv.win.amd64.tar.gz Json2Csv.win.amd64.exe
rm  Json2Csv.win.amd64.exe

CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Json2Csv.win.386.exe main.go
./upx -9 Json2Csv.win.386.exe
tar -czvf Json2Csv.win.386.tar.gz Json2Csv.win.386.exe
rm  Json2Csv.win.386.exe

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Json2Csv.darwin.amd64 main.go
./upx -9 Json2Csv.darwin.amd64
tar -czvf Json2Csv.darwin.amd64.tar.gz Json2Csv.darwin.amd64
rm  Json2Csv.darwin.amd64

go build -gcflags=-trimpath=$GOPATH -asmflags=-trimpath=$GOPATH -ldflags "-w -s" -o Json2Csv.linux.amd64 main.go
./upx -9 Json2Csv.linux.amd64
tar -czvf Json2Csv.linux.amd64.tar.gz Json2Csv.linux.amd64
rm  Json2Csv.linux.amd64

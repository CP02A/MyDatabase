$env:GOOS = "darwin"
go build -o .\builds\MyDatabase-MacOS
$env:GOOS = "freebsd"
go build -o .\builds\MyDatabase-FreeBSD
$env:GOOS = "linux"
go build -o .\builds\MyDatabase-Linux
$env:GOOS = "windows"
go build -o .\builds\MyDatabase-Windows.exe
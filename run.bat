@echo off
setlocal
cd /d "%~dp0"

set "ADDR=%FUTU_ADDR%"
if "%ADDR%"=="" set "ADDR=127.0.0.1:11111"

echo Running futuapi4go-demo (OpenD: %ADDR%)...
echo.

go run ./cmd/demo/main.go

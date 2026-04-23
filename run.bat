@echo off
setlocal
cd /d "%~dp0"

set "ADDR=%FUTU_ADDR%"
if "%ADDR%"=="" set "ADDR=127.0.0.1:11111"

set "EXAMPLE=%1"
if "%EXAMPLE%"=="" set "EXAMPLE=00_connect"

echo Running futuapi4go-demo/examples/%EXAMPLE% (OpenD: %ADDR%)...
echo.

go run ./examples/%EXAMPLE%

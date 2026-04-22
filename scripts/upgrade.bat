@echo off
echo Upgrading dependencies...

go get -u github.com/shing1211/futuapi4go@latest
go mod tidy

if %ERRORLEVEL% EQU 0 (
    echo Upgrade successful!
) else (
    echo Upgrade failed!
    exit /b %ERRORLEVEL%
)

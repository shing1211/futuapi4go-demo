@echo off
setlocal
cd /d "%~dp0"

echo Cleaning build artifacts...

if exist cmd\demo\futuapi4go-demo.exe (
    del /q cmd\demo\futuapi4go-demo.exe
    echo Deleted cmd\demo\futuapi4go-demo.exe
)

echo Done.

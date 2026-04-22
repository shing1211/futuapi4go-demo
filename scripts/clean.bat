@echo off
echo Cleaning build artifacts...

if exist cmd\demo\futuapi4go-demo.exe (
    del cmd\demo\futuapi4go-demo.exe
    echo Deleted cmd\demo\futuapi4go-demo.exe
)

if exist go.work (
    del go.work
    echo Deleted go.work
)

if exist go.work.sum (
    del go.work.sum
    echo Deleted go.work.sum
)

echo Done.

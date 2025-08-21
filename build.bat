@echo off
rem ------------------------- 配置部分 -------------------------
rem 定义应用程序名称
set APP_NAME=Crontab

rem 定义输出目录
set OUTPUT_DIR=bin

rem 自定义 main.go 的路径
set MAIN_PATH=cmd\main.go

rem ------------------------- 检查并创建输出目录 -------------------------
rem 检查并创建输出目录
if not exist %OUTPUT_DIR% (
    mkdir %OUTPUT_DIR%
)

rem ------------------------- 编译不同平台 -------------------------

rem 编译为 Linux 可执行文件
echo 编译 Linux 可执行文件...
set GOOS=linux
set GOARCH=amd64
go build -o %OUTPUT_DIR%\%APP_NAME%-linux-amd64 %MAIN_PATH%

rem 完成
echo 所有平台的可执行文件已经生成。
pause

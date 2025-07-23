@echo off
chcp 65001

REM 初始化参数标志
set WEB_FLAG=0
set PUSH_FLAG=0
set COMMENT_FLAG=0

REM 解析参数
:parse_args
if "%~1"=="" goto end_args

if "%~1"=="-web" (
    set WEB_FLAG=1
) else if "%~1"=="-p" (
    set PUSH_FLAG=1
) else if "%~1"=="-m" (
    set COMMENT_FLAG=1
    set USER_COMMENT="%~2"
    shift
)
shift
goto parse_args

:end_args

REM 执行-web操作
if %WEB_FLAG%==1 (
    cd web_ui
    call build.bat
    cd ../
)
set VERSION=""
REM 从core/version.go中读取版本号
for /f "tokens=2 delims==" %%v in ('type core\version.go ^| find "VERSION"') do set VERSION=%%~v
if %VERSION%=="" (
    echo 错误：无法从core/version.go读取版本号
    exit /b 1
)
set tag=v%VERSION%
echo 当前版本: %VERSION% TAG: %tag%

REM 设置comment
echo %COMMENT_FLAG%
if %COMMENT_FLAG%==1 (
    set comment=%USER_COMMENT%
) else (
    set comment=%VERSION%
)

echo %comment%
git add .
git commit -m "%VERSION% %comment%"

REM 执行git操作
if %PUSH_FLAG%==1 (
    git push -u github main 
    git push -u gitee master
)


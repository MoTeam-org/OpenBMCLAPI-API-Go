@echo off
REM 构建前端
cd web
npm run build
cd ..

REM 复制构建文件
rmdir /s /q service\web\dist
mkdir service\web
xcopy /s /e /i web\dist service\web\dist

REM 构建 Go 程序
go build 
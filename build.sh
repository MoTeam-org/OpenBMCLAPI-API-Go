#!/bin/bash
# 构建前端
cd web
npm run build
cd ..

# 复制构建文件
rm -rf service/web/dist
mkdir -p service/web
cp -r web/dist service/web/

# 构建 Go 程序
go build 
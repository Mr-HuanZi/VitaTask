FROM golang:1.20.2-alpine

# 设置工作目录
WORKDIR $GOPATH/src

COPY ../.. ./service

# 编译并安装项目文件
RUN cd $GOPATH/src/service && go build -o $GOPATH/bin/service main.go
# 覆盖默认指令
ENTRYPOINT ["/go/bin/service"]

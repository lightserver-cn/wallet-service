# 第一个阶段：指定基础镜像
FROM golang:1.23.3-alpine3.20 AS build

# 设置环境信息：GOPROXY 和 GO111MODULE
ENV GOPROXY=https://goproxy.cn,direct
ENV GO111MODULE=on

# 设置go包缓存（防止出现 $GOPATH/go.mod exists but should not）
WORKDIR /go/cache

# 将 go.mod 和 go.sum 文件复制到容器的工作目录
# 下载 Go 模块依赖
COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build as release

# 设置工作目录为 /go/release
WORKDIR /go/release

# 将当前目录下的所有文件复制到容器的工作目录
ADD . .

# 编译 Go 项目，并将可执行文件命名为 app，并将其复制到 /usr/local/bin 目录
# -o：指定输出文件名。默认情况下，输出文件名与包名相同。
# -v：显示详细的编译过程信息
# -race：启用数据竞争检测。这将在编译时插入额外的代码以检测潜在的数据竞争问题。
# -ldflags：传递链接器标志。这可以用于设置版本信息、优化等。
# go build -ldflags "-s -w"
# -s: 省略符号表和调试信息
# -w: 省略DWARF符号表
# -gcflags：传递编译器标志。这可以用于设置编译器优化、调试信息等。
# go build -gcflags "-N -l"
# -work：显示编译过程中的临时工作目录。
# -x：显示执行的每个命令，但不执行它们。
# -a：强制重新编译所有依赖包，而不仅仅是修改过的包。
# -i：安装与目标包相关的依赖包，但不编译它们。
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -installsuffix cgo -v -o server cmd/main.go

#!/bin/sh
# 当前用于构建golang服务，打包为二进制文件
# 考虑到目录问题，仅能在项目根目录下执行

# 检查是否传递了参数
if [ $# -eq 0 ]; then
  echo "请输入你将构建的 golang 二进制包包名"
  exit 1
fi

# 临时的镜像文件
image_name="tmp-$1:latest"
# 临时的容器
container_name="tmp-$1-container"
# golang 二进制包名
binary_filename="$1"

echo "您将构建的 golang 二进制包，包名为：$binary_filename"

# 使用 dockerfile 构建临时镜像
docker build -f dockerfile/Dockerfile -t "$image_name" .

# 使用构建好的临时镜像创建一个临时容器
docker create --name "$container_name" "$image_name"

# 将容器中构建的二进制文件复制到宿主机上
docker cp "$container_name:/go/release/server" "./$binary_filename"

# 删除临时容器
docker rm "$container_name"

# 删除镜像（可以不用删除，以便二次构建镜像时可以依赖镜像缓存层）
docker rmi "$image_name"

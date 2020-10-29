#openjdk-8为基础镜像dev
FROM reg.c5game.cn/library/c5game-base-image:latest

#暴露容器端口
EXPOSE 17001

#拷贝执行文件
ADD iris-seed /opt/iris-seed
ADD
#可选jvm参数
ENTRYPOINT exec sh /opt/run.sh
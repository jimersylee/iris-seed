#!/bin/bash
#用来构建jar包,构建docker镜像,推送docker镜像,更新docker镜像
set -e

COMMAND=${1}
IMAGE_NAME=registry.cn-hongkong.aliyuncs.com/c5g_prod/go-steam-proxy:latest
echo "$(date) build..."
#编译包,本地可以使用
function compile(){
  echo "start to compile"
  go build
  echo "compiled"
}
#构建docker镜像,本地可以使用
function build_image(){
  echo "building image"
  docker build . -t ${IMAGE_NAME} -f Dockerfile
  echo "build image done"
}

#将构建好的docker镜像推送到仓库,本地可以使用
function push_image(){
  echo "pushing image";
  docker push ${IMAGE_NAME}
  echo "yes baby"
}



case ${COMMAND} in
  compile)
    compile
    exit
    ;;
  build_image)
    build_image
    exit
    ;;
  push_image)
    push_image
    exit
    ;;
  update_docker)
    update_docker
    exit
    ;;
  all)
    compile 
    build_image 
    push_image 
    exit
    ;;
   deploy)
    update_code
    compile
    build_image
    push_image
    update_docker
    exit
    ;;
  *)
  echo "
  all:编译,打jar包,构建docker镜像,推送docker镜像
  deploy:服务器使用,包含更新代码,编译,打镜像,更新docker全套
  compile:编译jar包
  build_image:构建docker镜像
  push_image:推送docker镜像
  "
  exit
  ;;

esac

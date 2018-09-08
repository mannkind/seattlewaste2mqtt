#!/bin/bash -e
TAG=${1:-latest}
ARCHS="amd64 arm32v6 arm64v8"
IMAGE="mannkind/seattlewaste2mqtt"

for arch in $ARCHS; do
  case ${arch} in
    amd64   ) qemu_arch="x86_64"; golang_arch="amd64";;
    arm32v6 ) qemu_arch="arm"; golang_arch="arm";;
    arm64v8 ) qemu_arch="aarch64"; golang_arch="arm64";;
  esac
  cp Dockerfile.cross Dockerfile.${arch}
  sed -i "" "s|__BASEIMAGE_ARCH__|${arch}|g" Dockerfile.${arch}
  sed -i "" "s|__QEMU_ARCH__|${qemu_arch}|g" Dockerfile.${arch}
  sed -i "" "s|__GOLANG_ARCH__|${golang_arch}|g" Dockerfile.${arch}
  if [ ${arch} == 'amd64' ]; then
    sed -i "" "/__CROSS_/d" Dockerfile.${arch}
  else
    sed -i "" "s/__CROSS_//g" Dockerfile.${arch}
  fi
done


for arch in $ARCHS; do
  docker build -f Dockerfile.${arch} -t ${IMAGE}:${arch}-${TAG} .
  docker push ${IMAGE}:${arch}-${TAG}
done

docker manifest create ${IMAGE}:${TAG} ${IMAGE}:amd64-${TAG} ${IMAGE}:arm32v6-${TAG} ${IMAGE}:arm64v8-${TAG}
docker manifest annotate ${IMAGE}:${TAG} ${IMAGE}:arm32v6-${TAG} --os linux --arch arm
docker manifest annotate ${IMAGE}:${TAG} ${IMAGE}:arm64v8-${TAG} --os linux --arch arm64 --variant armv8
docker manifest push --purge ${IMAGE}:${TAG}

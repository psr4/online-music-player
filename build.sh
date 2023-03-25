#/bin/bash

mkdir build

#docker build . -t gobuild
docker run --rm -v $PWD:/service -it gobuild "go mod download && go build -v -o player"

cp -r ./static ./build/
cp -r ./templates ./build/
cp -r ./music ./build/
mv player ./build/
cp NasDockerfile ./build/Dockerfile
cd build
docker build . -t a5108/online-music-player
docker push a5108/online-music-player

rm -rf build



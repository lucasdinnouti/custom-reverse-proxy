docker build  -t processor -f ./components/processor/Dockerfile .
docker build  -t proxy -f ./components/proxy/Dockerfile .
docker build  -t runner -f ./components/runner/Dockerfile .
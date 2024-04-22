cd components/processor
docker build  -t processor .

cd ../proxy
docker build  -t proxy .

cd ../runner
docker build  -t runner .

cd ../..
kind create cluster --name crp --config ./infra/manifests/cluster/kind.yaml

kubectl cluster-info --context kind-crp

kind load docker-image processor:latest --name crp
kind load docker-image runner:latest --name crp
kind load docker-image proxy:latest --name crp

export ALGORITHM="round_robin"

for f in ./infra/manifests/components/*.yaml; do envsubst < $f | kubectl apply -f -; done

kubectl apply -f ./infra/manifests/metrics
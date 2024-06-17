kind create cluster --name crp --config ./infra/manifests/cluster/kind.yaml

kubectl cluster-info --context kind-crp

kind load docker-image processor:latest --name crp
kind load docker-image runner:latest --name crp
kind load docker-image proxy:latest --name crp

export ALGORITHM="round_robin"
envsubst < infra/manifests/components/proxy.yaml | kubectl apply -f -

kubectl apply -f ./infra/manifests/components
kubectl apply -f ./infra/manifests/metrics
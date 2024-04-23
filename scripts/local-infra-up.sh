kind create cluster --name crp

kubectl cluster-info --context kind-crp

kind load docker-image processor:latest --name crp
kind load docker-image runner:latest --name crp
kind load docker-image proxy:latest --name crp

kubectl apply -f ./infra/manifests
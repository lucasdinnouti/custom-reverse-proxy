## now loop through the above array
for alg in "round_robin" "weighted_round_robin" "metadata" "machine_learning" "machine_learning_weight"
do
    kubectl delete pod runner
    kubectl delete pod proxy

    export ALGORITHM=$alg
    envsubst < infra/manifests/components/proxy.yaml | kubectl apply -f -
    kubectl apply -f infra/manifests/components/runner.yaml

    sleep 30m
done


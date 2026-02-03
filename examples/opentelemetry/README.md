# Example - OpenTelemetry

As the provider supports OpenTelemetry, this example shows how to technically use it in a demo environment.

## Try it out!

- Run an OpenTelemetry Collector to capture signals (here we focus on traces), Jaeger for distributed traces visualization, an OCI registry for distributing scenarios, our [instrumented and repackaged CTFd](https://github.com/ctfer-io/ctfd-packaged) for OpenTelemetry support with the [`ctfd-chall-manager`](https://github.com/ctfer-io/ctfd-chal-manager) plugin, [`ctfd-setup`](https://github.com/ctfer-io/ctfd-setup) for bootstrapping, and [Chall-Manager](https://github.com/ctfer-io/chall-manager) for on-demand challenge instances (quite a lot actually! :partying_face:).
    ```bash
    docker compose up -d
    ```

- Initialize the Terraform setup.
    ```bash
    terraform init
    ```

- Setup the environment variables for the provider to pick up its configuration.
    ```bash
    export OTEL_EXPORTER_OTLP_ENDPOINT=dns://localhost:4317
    export OTEL_EXPORTER_OTLP_INSECURE=true
    export CTFD_URL=http://localhost:8000
    export CTFD_ADMIN_USERNAME=ctfer
    export CTFD_ADMIN_PASSWORD=ctfer
    ```

- Build the scenario
    ```bash
    go install github.com/ctfer-io/chall-manager/cmd/chall-manager-cli@latest
    chall-manager-cli --url localhost:8000 scenario \
        --scenario localhost:5000/some/scenario:v0.1.0 \
        --directory ../../provider/scenario \
        --insecure
    ```

- Use it! :smile:
    ```bash
    terraform apply -auto-approve
    ```

- You can also run the acceptance tests of the TF provider :thinking:
    ```bash
    (
        cd ../..
        export REGISTRY="$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps --format="{{.Names}}" | grep registry)):5000"
        TF_ACC=1 go test ./provider/ -v -run=^TestAcc_
    )
    ```

- Open Jaeger and visualize traces: [`http://localhost:16686`](http://localhost:16686)
<div align="center">
    <img src="jaeger.png" width="1000px">
</div>

- You can delete the infra :wink:
    ```bash
    terraform destroy -auto-approve
    docker compose down -v
    ```

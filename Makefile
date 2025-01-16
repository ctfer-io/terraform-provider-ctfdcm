.PHONY: test-acc
test-acc:
	echo "name: scenario" >> Pulumi.yaml
	echo "runtime:" >> Pulumi.yaml
	echo "  name: go" >> Pulumi.yaml
	echo "  options:" >> Pulumi.yaml
	echo "    binary: ./main" >> Pulumi.yaml
    echo "description: An example scenario." >> Pulumi.yaml
	go build -o main scenario/main.go
	zip -r scenario.zip Pulumi.yaml main
	rm main Pulumi.yaml

	TF_ACC=1 \
	go test ./provider/ -v -run=^TestAcc_ -count=1 -coverprofile=cov.out -coverpkg "github.com/ctfer-io/terraform-provider-ctfdcm/provider"

	rm scenario.zip

.PHONY: docs
docs:
	go generate ./...

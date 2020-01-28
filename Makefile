build:
	go get gopkg.in/yaml.v2
	go build hookexec.go
example:
	go get gopkg.in/yaml.v2
	go run hookexec.go ./examples/config.yml &> ./logs/example.log


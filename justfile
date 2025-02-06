[private]
default:
    @just --list

# install substrate to ~/.local/bin
install:
    go build -o ~/.local/bin/substrate cmd/substrate/main.go

# generate documentation
gen:
    protoc --go_out=. proto/*.proto
    terraform-docs markdown modules/artifacts > modules/artifacts/README.md
    terraform-docs markdown modules/feature > modules/feature/README.md
    terraform-docs markdown modules/role > modules/role/README.md
    terraform-docs markdown modules/substrate > modules/substrate/README.md

# initialize terraform
init:
    AWS_PROFILE=prod.kaixo.io tofu -chdir=spec/prod init
    AWS_PROFILE=dev.kaixo.io tofu -chdir=spec/dev init

# apply terraform
apply:
    AWS_PROFILE=prod.kaixo.io tofu -chdir=spec/prod apply
    AWS_PROFILE=dev.kaixo.io tofu -chdir=spec/dev apply

# destroy terraform
destroy:
    AWS_PROFILE=prod.kaixo.io tofu -chdir=spec/prod destroy
    AWS_PROFILE=dev.kaixo.io tofu -chdir=spec/dev destroy

# test infrastructure
test:
    cd spec && bash test.sh

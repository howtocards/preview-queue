package generated

//go:generate rm -rf models restapi
//go:generate swagger generate server -f ./../../../../swagger.yml --exclude-main --include-buildapi --strict

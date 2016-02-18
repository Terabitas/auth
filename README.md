# Development

## Running integration tests

Execute tests

Provisioning is happening inside `TestMain`
```
go test -v -tags integration
```

## Testing locally

Build project container, it will produce docker image:
```
nildev build github.com/nildev/auth
```

Create `local.env` file in root dir with content:
```
ND_SIGN_KEY=__YOUR_VAL__
ND_ENV=dev
ND_GITHUB_CLIENT_ID=__YOUR_VAL__
ND_GITHUB_SECRETE=__YOUR_VAL__
```

Run containers:
```
docker-compose -f docker-compose.yml up -d
```

Execute HTTP request:
```
curl -X GET  http://$(docker-machine ip nildev):8080/api/v1/auth/github/authCode -v
```

## Project Details

### Release Notes

See the [releases tab](https://github.com/nildev/auth/releases) for more information on each release.

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and contacting developers via IRC and mailing lists.

### License

Project is released under the MIT license. See the [LICENSE](LICENSE) file for details.

Specific components of project use code derivative from software distributed under other licenses; in those cases the appropriate licenses are stipulated alongside the code.
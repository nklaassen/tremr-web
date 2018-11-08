# tremr-web
Front-end and back-end code for the web interface of the Tremr app developed for CMPT 275

## building and running locally
You probably don't want to clone this repo manually, instead:
1. Install [Go](https://github.com/golang/go)
2. `go get github.com/nklaassen/tremr-web`
3. `cd "${GOPATH}/github.com/nklaassen/tremr-web"`
4. `go build`
4. `./tremr-web` will run the server locally on port 8080. This needs to be run from the root directory of the repo so it can find the `www` directory.

At this point, try opening your browser to `localhost:8080`.

If you just want to test the API, try the following:

    curl -X POST -i http://localhost:8080/api/tremors --data '
        {
            "resting": 38,
            "postural": 47
        }'
    curl -X GET localhost:8080/api/tremors

    curl -X POST localhost:8080/api/meds --data '
        {
            "name": "testmed",
            "dosage": "10 mL",
            "schedule": {
              "mo": true,
              "we": true
            },
            "reminder": false
        }'
    curl -X GET localhost:8080/api/meds


## contributing
### front-end
Static html, css, and js files will be served from the `www` directory, add and edit what you need there. Make sure you are running the webserver (see above) if you need access to the api.

### back-end
I hope you like Go! After making changes, make sure to re-run `go build` and run the resulting binary from the root directory of this repo (where this README lives).

`go test` is also really useful for testing things, see `database/tremor_test.go` for an example. Note that this command needs to be run from within the package directory you are testing (`cd database; go test`).

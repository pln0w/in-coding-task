# Coding task

## Commands 

Run Docker container to start HTTP server  
```
make run
```
or directly `docker-compose up`  
  

Check how it works
```
make demo
```  
or directly make request for _0.0.0.0:8080_
  

  
Run tests
```
make test
```
or directly `go test`


## Description

### Web server routes
* GET `/` - health check
* GET `/routes` - retrieve details about best nearest pickup point from source location  
Query params:  
    + `src` - ie. `13.388860,52.517037`
    + `dst` - ie. `13.397634,52.529407`

Example call:  
```
0.0.0.0:8080/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219
```

### Files  
```
.
├── app.yaml                
├── controller.go           # functions for handling HTTP request/response
├── docker-compose.yaml      
├── Dockerfile              
├── main.go                 # entry point - HTTP server run
├── Makefile                
├── parser.go               # core functions for deserializing incoming data
├── README.md               
├── router.go               # router for HTTP server
├── sort.go                 # sorting functions
└── sort_test.go            # sorting test

```

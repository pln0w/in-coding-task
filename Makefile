run:
	docker-compose up

build:
	go build -o main .

demo:
	curl -i -H "Accept: application/json" -H "Content-Type: application/json" -X GET 0.0.0.0:8080/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219&dst=13.397634,52.529598&dst=13.397632,52.529488

test:
	go test

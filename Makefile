# note: call scripts from /scripts

# make new_migration MESSAGE_NAME="message name"
new_migration:
	migrate create -ext sql -dir scripts/migration/ -seq $(MESSAGE_NAME)

# run db migration (change your connection string if your config differs)
up_migration:
	migrate -path scripts/migration/ -database "mysql://root:123456@tcp(127.0.0.1:3306)/engineerpro?charset=utf8mb4&parseTime=True&loc=Local" -verbose up

# reverse db migration (change your connection string if your config differs)
down_migration:
	migrate -path scripts/migration/ -database "mysql://root:123456@tcp(127.0.0.1:3306)/engineerpro?charset=utf8mb4&parseTime=True&loc=Local" -verbose down

# generate grpc Go files from protobuf
proto_aap:
	protoc --proto_path=pkg/types/proto/ --go_out=pkg/types/proto/pb/authen_and_post --go_opt=paths=source_relative \
        --go-grpc_out=pkg/types/proto/pb/authen_and_post --go-grpc_opt=paths=source_relative \
        authen_and_post.proto

# tidy go.mod/go.sum files
tidy:
	go mod tidy

# remove docker volumes & images
docker_clear:
	docker volume rm $(docker volume ls -qf dangling=true) & docker rmi $(docker images -f "dangling=true" -q)

# up all services with forced  build and re-create options
compose_up_rebuild:
	docker compose up --build --force-recreate

# docker compose up all services without rebuilding images
compose_up:
	docker compose up

# docker compose down
compose_down:
	docker compose down

# generate api docs
gen_swagger:
	swag init -g cmd/web_app/main.go
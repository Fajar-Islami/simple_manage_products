# Load environment variables from .env file
include .env
export 


registry:=ghcr.io
username:=fajar-islami
image:=simple_manage_products
tags:=latest


git-add:
	git add .
	git commit -am '${cmt}'

test:
	echo ${cmt}

entermysql:
	docker exec -it mysql_simple_manage_products mysql -u ${mysql_username} -p${mysql_password} ${mysql_dbname}

entermysqlroot:
	docker exec -it mysql_simple_manage_products mysql -u root -p${mysql_root_password} ${mysql_dbname}

enterredis:
	docker exec -it redis_simple_manage_products redis-cli 
	#  AUTH 1234

runenv:
	docker compose up -d

run-redis:
	docker compose up -d redis_simple_manage_products

run:
	docker compose up -d
	go run main.go

commit:
	git add .
	git commit -am '${cmt}'

struct:
	gomodifytags -file ${file} -struct ${struct} -add-tags ${tags}

stop:
	docker compose stop

down:
	docker compose down -v

logs:
	docker compose logs -f

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o  ./dist/example ./main.go

dockerbuild:
	docker build --rm -t ${appName} .
	docker image prune --filter label=stage=dockerbuilder -f

dockerun:
	docker run --name ${appName}  -p 8080:8080 ${appName} 

dockerrm:
	docker rm ${appName} -f
	docker rmi ${appName}

dockeenter:
	docker exec -it ${appName} bash

push-image:
	docker build -t ${registry}/${username}/${image}:${tags} .
	export CR_PAT=${CR_PAT}
	echo ${CR_PAT} | docker login ${registry} -u ${USERNAME} --password-stdin -interactive --tty
	docker push ${registry}/${username}/${image}:${tags}


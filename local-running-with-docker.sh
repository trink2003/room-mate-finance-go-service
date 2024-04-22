#!/bin/bash

handle_error() {
    echo -e "\n\n >> An error occurred on line $1 \n\n"
    exit 1
}

trap 'handle_error $LINENO' ERR

if [ "$#" -lt 1 ]; then
cat <<EOF | cat -


>> Script usage: $s0 images_name

Where:
    - images_name: specify the name of the image will be built

EOF
exit 1
fi

lastest_git_commit_hash_id=$(git log -n 1 --pretty=format:'%h')


images_name="$1"
current_time=$(date -d "$b 0 min" "+%Y%m%d%H%M%S")
images_tag="${current_time}_${lastest_git_commit_hash_id}"

cat <<EOF | cat -


>> Deloying new version of service with images tag: ${images_tag}

EOF

go_mod_command="go mod download"
cat <<EOF | cat -


>> Downloading library with go mod
>> Command: ${go_mod_command}

EOF
eval ${go_mod_command}

go_build_command="CGO_ENABLED=0 GOOS=linux go build -o ./go_app"
cat <<EOF | cat -


>> Building go project to executable program
>> Command: ${go_build_command}

EOF
eval ${go_build_command}

# -----------------
docker_build_command="docker build -f ./Dockerfile -t ${images_name}:${images_tag} ."
cat <<EOF | cat -


>> Building image with Docker
>> Command: ${docker_build_command}

EOF
eval ${docker_build_command}

cat <<EOF | cat -


>> Docker compose

EOF

cat <<EOF | docker compose -f - down
services:
  room-mate-finance-go-service:
    container_name: room-mate-finance-service
    image: ${images_name}:${images_tag}
    environment:
      - "DATABASE_USERNAME=postgres"
      - "DATABASE_PASSWORD=postgres"
      - "DATABASE_HOST=10.0.2.10"
      - "DATABASE_PORT=5432"
      - "DATABASE_NAME=room-mate-finance"
      - "GIN_MODE=release"
      - "JWT_SECRET_KEY=Q8OzIHRo4buDIGfhu41pIGFuaCBsw6AgxJHhurlwIHRyYWkgbmjhuqV0IFZp4buHdCBOYW0="
      - "JWT_EXPIRE_TIME=1440"
      - "DATABASE_MIGRATION=false"
      - "DATABASE_INITIALIZATION_DATA=false"
    ports:
      - "8080:8080"
EOF

old_images=$(docker images | grep room | awk '{print $3}')

for VAR in $old_images
do
    if docker rmi $VAR
    then
        echo -e "\n\n >> Remove successfully"
    else
        echo -e "\n\n >> Remove failed"
    fi
done

cat <<EOF | docker compose -f - up -d
services:
  room-mate-finance-go-service:
    container_name: room-mate-finance-service
    image: ${images_name}:${images_tag}
    environment:
      - "DATABASE_USERNAME=postgres"
      - "DATABASE_PASSWORD=postgres"
      - "DATABASE_HOST=10.0.2.10"
      - "DATABASE_PORT=5432"
      - "DATABASE_NAME=room-mate-finance"
      - "GIN_MODE=release"
      - "JWT_SECRET_KEY=Q8OzIHRo4buDIGfhu41pIGFuaCBsw6AgxJHhurlwIHRyYWkgbmjhuqV0IFZp4buHdCBOYW0="
      - "JWT_EXPIRE_TIME=1440"
      - "DATABASE_MIGRATION=false"
      - "DATABASE_INITIALIZATION_DATA=false"
    ports:
      - "8080:8080"
EOF

remove_built_go_executed_program="rm -f ./go_app"
cat <<EOF | cat -


>> Removing built go executed program
>> Command: ${remove_built_go_executed_program}

EOF
if eval ${remove_built_go_executed_program}
then
cat <<EOF | cat -


>> Successfully

EOF
else
cat <<EOF | cat -


>> Failed

EOF
fi

cat <<EOF | cat -


>> Deployment process has been done

EOF

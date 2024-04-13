#!/bin/bash

#!/bin/bash

cd ./room-mate-finance-go-service/

ssh_user="ae403"
ssh_host="0.tcp.ap.ngrok.io"
ssh_port="17742"
target_dir='/home/ae403/service/room-mate-finance'

images_name='tuanloc/room-mate-finance'
app_name='room-mate-finance'
app_namespace='service'
images_tag=$(date -d "$b 0 min" "+%Y%m%d_%H%M%S")
replica=2

echo -e "\n\n >> go mod download"
go mod download

echo -e "\n\n >> go build -o ./go-app"
CGO_ENABLED=0 GOOS=linux go build -o ./go-app

echo -e "\n\n >> Docker build"
docker build -f ./Dockerfile2 -t ${images_name}:${images_tag} .

echo -e "\n\n >> Push to Docker registry"
docker push ${images_name}:${images_tag}
# kind load docker-image ${images_name}:${images_tag}

echo -e "\n\n >> Copy file to target host"
cat <<EOF | cat - | tee ./helm/Chart.yaml
apiVersion: v2
name: room-mate-finance-go-service
description: A Helm chart for Kubernetes to deploy the room-mate-finance-go-service service

type: application

version: ${images_tag}

appVersion: latest
EOF
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -r ./helm/ $ssh_user@$ssh_host:$target_dir

echo -e "\n\n >> Upgrade application"
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; helm upgrade -i --force --set image.name=${images_name},image.tag=${images_tag},replica=${replica},port=6060 ${app_name} -n ${app_namespace} --create-namespace ./helm"

echo -e "\n\n >> Remove images"
docker rmi ${images_name}:${images_tag} || echo -e "\n\n >> No images"
rm -f ./go-app

echo -e "\n\n >> Revert chart info \n\n"

cat <<EOF | cat - | tee ./helm/Chart.yaml
apiVersion: v2
name: room-mate-finance-go-service
description: A Helm chart for Kubernetes to deploy the room-mate-finance-go-service service

type: application

version: 0.1.0

appVersion: latest
EOF

echo -e "\n\n >> Done \n\n"

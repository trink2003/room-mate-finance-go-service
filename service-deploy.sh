#!/bin/bash

handle_error() {
    echo -e "\n\n >> An error occurred on line $1 \n\n"
    exit 1
}

trap 'handle_error $LINENO' ERR

if [ "$#" -lt 8 ]; then
cat <<EOF | cat -


>> Script usage: $0 images_name app_name app_namespace replica ssh_user ssh_host ssh_port target_dir

Where:
    - images_name: specify the name of the image will be built
    - app_name: the app name when upgrading helm chart
    - app_namespace: K8S namespace to build on
    - replica: number of application instance will be run
    - ssh_user: username of target host to run
    - ssh_host: ip or domain of target host to run
    - ssh_port: port of target host to run
    - target_dir: the directory that command will be execute each time as well as the directory that will be store chart folder

EOF
exit 1
fi

# lastest_git_commit_hash_id=$(git log --branches --format="%H" -n 1) # this command will get full hash commit id
lastest_git_commit_hash_id=$(git log --branches --format="%h" -n 1) # this command will get short hash commit id


images_name="$1"
app_name="$2"
app_namespace="$3"
replica="$4"
current_time=$(date -d "$b 0 min" "+%Y%m%d%H%M%S")
images_tag="${current_time}_${lastest_git_commit_hash_id}"

cat <<EOF | cat -


>> Deloying new version of service with images tag: ${images_tag}

EOF

ssh_user="$5"
ssh_host="$6"
ssh_port="$7"
target_dir="$8"

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
docker_push_command="docker push ${images_name}:${images_tag}"
cat <<EOF | cat -


>> Pushing images to image registry
>> Command: ${docker_push_command}

EOF
eval ${docker_push_command}
# -----------------

cat <<EOF | cat -


>> Updating chart information

EOF
cat <<EOF | cat - | tee ./helm/Chart.yaml
apiVersion: v2
name: room-mate-finance-go-service
description: A Helm chart for Kubernetes to deploy the room-mate-finance-go-service service

type: application

version: ${current_time}

appVersion: latest
EOF

cat <<EOF | cat -


>> Uploading necessary file to target host $ssh_host

EOF
# ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "eval \"mkdir -p ${target_dir}\""
scp -P $ssh_port -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -pr ./helm/ $ssh_user@$ssh_host:$target_dir

helm_upgrade_command="helm upgrade -i --force --set image.name=${images_name},image.tag=${images_tag},replica=${replica},port=6060 ${app_name} -n ${app_namespace} --create-namespace ./helm"
cat <<EOF | cat -


>> Upgrading helm chart of application
>> Command: ${helm_upgrade_command}

EOF
ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no $ssh_user@$ssh_host -p $ssh_port "cd ${target_dir} ; source ~/.bash_profile ; eval ${helm_upgrade_command}"

echo -e "\n\n >> Remove images \n\n"
docker_remove_image_command="docker rmi ${images_name}:${images_tag}"
cat <<EOF | cat -


>> Removing built images
>> Command: ${docker_remove_image_command}

EOF

if eval ${docker_remove_image_command}
then
cat <<EOF | cat -


>> Removing images successfully

EOF
else
cat <<EOF | cat -


>> Can not remove images

EOF
fi

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


>> Revert chart info

EOF

cat <<EOF | cat - | tee ./helm/Chart.yaml
apiVersion: v2
name: room-mate-finance-go-service
description: A Helm chart for Kubernetes to deploy the room-mate-finance-go-service service

type: application

version: 0.1.0

appVersion: latest
EOF

ecat <<EOF | cat -


>> Deployment process has been done

EOF

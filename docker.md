## mysql
  docker run -itd --name mysql-dev -p 3301:3306 --restart=always -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7
  docker run -itd --name mysql-prd -p 3302:3306 --restart=always -v /local_path:/var/lib/mysql -u 1000:1000 --privileged=true -e MYSQL_ROOT_PASSWORD=Vista@devops1 mysql:5.7
  grant all on *.* to 'root'@'%' identified by '123456' with grant option;  # grant root authens
  docker exec -it container_id /bin/bash
  docker run -itd --name xxx  -v /local_path:/dest_path --privileged=true mysql:5.7

## jenkins docker
docker run -itd --name myjenkins -p 8080:8080 -p 50000:50000 --restart=always --env JAVA_OPTS="-Dhudson.util.ProcessTree.disable=true" -v /home/artifactory/jenkins_home:/var/jenkins_home -u 1000 --privileged=true jenkins/jenkins:latest
docker run -itd --name testjenkins -p 7070:8080 -p 50001:50000 --restart=always --env JAVA_OPTS="-Dhudson.util.ProcessTree.disable=true" -v /home/artifactory/test_jenkins_home:/var/jenkins_home -u 1000 --privileged=true jenkins/jenkins:latest

## get container log path
sudo docker inspect --format '{{.LogPath}}' container_id

## transfer docker image
docker pull nginx
docker save -o nginx.docker nginx
docker load -i nginx.docker 

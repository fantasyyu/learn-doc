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

# django

  python manage.py createsuperuser
  python manage.py shell -> Group.objects.create(name='GUEST')

  1. runserver 0.0.0.0:8002 PYTHONUNBUFFERED=1;DJANGO_SETTINGS_MODULE=BOSS.settings.dev  # use to determine if dev env
  2. os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'BOSS.settings.prd')  in prd_wsgi.py
  3. AUTH_USER_MODEL = 'users.Users' # customize user model
  4. #Customize authentication  in settings.py
      AUTHENTICATION_BACKENDS = (
      'utils.auth.MyLdapBackend',
      'django.contrib.auth.backends.ModelBackend',
  )
  5. migrate authtoken  # generate auth token table
  6. queryset = Q() || queryset &= Q(name__icontains=request.GET.get('name'))
  7. apt-get install libsasl2-dev python-dev libldap2-dev libssl-dev && pip install python-ldap
  8. related_name:
    Class: TestCase
    Class: TestCasePlan
      testcase = models.Foriegn(related_name='testcaseplan')
    # to find all special TestCase_id be use in which TestCasePlan
    TestCase.objects.get(id=TestCase_id).testcaseplan.all()
  9. URL and id:
      url: url(r'^api/v1/testplan/(.*)', TestPlanView.as_view()),
      view: def get(self, request, id)
  10. rest call
      curl --header "Authorization: Token d0a90ed9a769616b344353578df93d3e6577bf42" "http://127.0.0.1:8002/api/v1/testtemplate?name=1&type=2"
  11. Change primarky key start from non-1 value
      operations = [
          migrations.RunSQL(
              'ALTER SEQUENCE APPNAME_USER_id_seq RESTART WITH 10000000;'
          )
      ]
  12. reset migrations
      python manage.py showmigrations 应用名
      python manage.py migrate --fake 应用名 zero
  13. 查看迁移
      先将应用模型类恢复到没有修改之前的状态
      python manage.py showmigrations 应用名
      删除app migrations下除init.py之外的所有文件
      删除数据库migrations表中该应用的所有迁移记录
      python manage.py makemigrations 应用名 # 会再次生成0001_initial.py 之类的文件
      python manage.py migrate 应用名 --fake-initial #  --fake-inital 会在数据库migrations表中记录当前这个app 执行到 0001_initial.py ，但不会去执行迁移文件
      修改应用模型类，并执行迁移命令
      python manage.py makemigrations 应用名
      python manage.py migrate 应用名
  14. obj.update will update update_time, filter().update won't
  15. log file:
      http://www.360doc.com/content/14/0708/10/16044571_392797799.shtml

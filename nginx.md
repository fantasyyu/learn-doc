## Install nginx
  apt-get install nginx
  
## enable nginx
  service nginx start
  
## config content
  vim /etc/nginx/conf.d/static.conf
  example:
  server {
        listen      90;
        server_name  localhost;
        location / {
            root   html;
            index  index.html index.htm;
        }
         location /static/ {
          root /root/;
        }
        location /api/v1/ {
                proxy_pass http://10.67.91.245:8002/;
        }
  }

## doc
https://prometheus.ac.cn/docs/introduction/overview/

## installation node-exporter
    wget https://github.com/prometheus/node_exporter/releases/download/v1.10.2/node_exporter-1.10.2.linux-amd64.tar.gz 
    tar -xvf 

## get prom docker
    docker run --name prometheus -d -p 0.0.0.0:9090:9090 -v ./prometheus.yml:/etc/prometheus/prometheus.yml ncr.nioint.com/docker.io/prom/prometheus:latest
    yml内容示例：
    '''
    scrape_configs:

    - job_name: node
    
      static_configs:
    
      - targets:
    
        - 10.64.40.40:9100
    
        - 10.64.40.39:9100
        - 10.64.40.36:9100
    '''
    参数：--web.enable-lifecycle 允许curl -X POST http://10.64.40.35:9090/-/reload

## grafana docker
    docker run -d --name grafana -p 3000:3000 ncr.nioint.com/docker.io/grafana/grafana:latest
    初始登录（admin：admin）
    Data Sources -> Add Prometheus -> URL -> http://你的PrometheusIP:9090
    Dashboards -> 1860



#Run prometheus
```
docker run -p 9090:9090 -d --name prometheus --network monitoring -v $(pwd)/prometheus:/etc/config prom/prometheus --config.file=/etc/config/prometheus.yml
```

```
docker run -p 9090:9090 -it --name prometheus --network monitoring -v $(pwd)/prometheus:/etc/config prom/prometheus --config.file=/etc/config/prometheus.yml
```

# Run alert manager
```
docker run -d -p 9093:9093 -v $(pwd)/alertmanager:/etc/config --network monitoring --name alertmanager prom/alertmanager --config.file=/etc/config/alertmanager.yml
```

# Telegram bot for alerting
```
docker run -d \
	-e 'ALERTMANAGER_URL=http://alertmanager:9093' \
	-e 'BOLT_PATH=/data/bot.db' \
	-e 'STORE=bolt' \
	-e 'TEMPLATE_PATHS=/templates/default.tmpl' \
	-e 'TELEGRAM_ADMIN=6936842' \
	-e 'TELEGRAM_TOKEN=686312965:AAFf0goR98-8uVxGM8tRVxhafHyyQ_KIynk' \
	-v $(pwd)/alertmanager-bot:/data \
	-v $(pwd)/alertmanager-bot:/templates \
  --network monitoring \
	--name alertmanager-bot \
	metalmatze/alertmanager-bot:0.4.0 --listen.addr 0.0.0.0:8080
```

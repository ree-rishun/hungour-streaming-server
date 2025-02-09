
## ローカル開発の流れ

```
docker-compose up -d

```


## ローカルデプロイ
```
docker build -t app -f .Dockerfile .
docker tag app asia-northeast1-docker.pkg.dev/hunger-gourmet/streaming-server/app:latest
docker push asia-northeast1-docker.pkg.dev/hunger-gourmet/streaming-server/app:latest
```

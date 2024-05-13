# GO-RATELIMITER

**IP**: O rate limiter restringe o numéro de requisições recebidas por endereço IP dentro de determinado intervalo de tempo.

**Token**: Da mesma forma que por IP porém, com diferentes tempos de expirações de acordo com o token usado. Por exemplo: `API_KEY: <TOKEN>`

*Opções:*

- FIVE_REQUEST_TOKEN=5
- SIX_REQUEST_TOKEN=6
- SEVEN_REQUEST_TOKEN=7


### Configuração padrão:

- MAX_IP_REQUESTS_PER_SECOND=3 `# Quantidade de requisições por segundo`
- BLOCK_TIME_SECONDS=10 `# Tempo de bloqueio para requisições`


### Como usar:

1. Execute o comando `make docker-compose` para subir o redis - que pode ser acesso na porta 8001;
2. Execute o comando na raiz do projeto: `go run main.go`;
3. Realize as requisições usando o arquivo localizado em: `api/api.http`;

4. Testes estão localizados nos seguintes caminhos: 

- `internal/infra/database/redis/limiter_info_repository_test.go`
- `internal/rate-limiter/rate-limiter_test.go`.
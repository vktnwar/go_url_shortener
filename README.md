# Go URL Shortener

Uma API simples e eficiente para encurtamento de URLs construída em Go, com Postgres, Redis e Rate Limiting por IP.

## ✨ Características

* 🔗 Geração de URLs encurtadas únicas
* 🏎 Redirecionamento rápido usando Redis como cache
* 📊 Métricas de cliques armazenadas no Postgres
* 🚦 Limite de requisições por IP
* 🛠 Healthcheck disponível em `/health`
* 🐳 Dockerizado para desenvolvimento e produção

## 🚀 Instalação

### Opção 1: Build local

```bash
git clone https://github.com/vktnwar/go_url_shortener.git
cd go_url_shortener
go build -o url_shortener ./cmd/main.go
```

### Opção 2: Executar diretamente

```bash
git clone https://github.com/vktnwar/go_url_shortener.git
cd go_url_shortener
go run ./cmd/main.go
```

### Opção 3: Docker Compose

```bash
docker-compose up --build
```

## 📖 Uso

### Healthcheck

```bash
curl http://localhost:8080/health
```

### Encurtar URL

```bash
curl -X POST http://localhost:8080/shorten \
 -H "Content-Type: application/json" \
 -d '{"url":"https://example.com"}'
```

**Resposta:**

```json
{
  "short_url": "http://localhost:8080/Ijzx-XLH"
}
```

### Redirecionar URL

```bash
curl -v http://localhost:8080/Ijzx-XLH
```

Redireciona para a URL original e incrementa o contador de cliques.

## ⚙️ Configuração

Crie um arquivo `.env` na raiz do projeto:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=url_shortener
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
RATE_LIMIT=5
RATE_WINDOW=1m
```

## 🧩 Estrutura do projeto

```
go_url_shortener/
├─ cmd/
│  └─ main.go
├─ config/
│  └─ config.go
├─ repository/
│  ├─ postgres_repository.go
│  └─ redis_repository.go
├─ server/
│  └─ router.go
├─ service/
│  └─ url_service.go
├─ middleware/
│  └─ rate_limiter.go
├─ .env
├─ Dockerfile
└─ docker-compose.yml
```

## 🛠️ Desenvolvimento

```bash
# Clonar repositório
git clone https://github.com/vktnwar/go_url_shortener.git
cd go_url_shortener

# Executar em modo desenvolvimento
go run ./cmd/main.go
```

## 📋 Requisitos

* Go 1.23 ou superior
* Docker e Docker Compose (opcional)
* PostgreSQL
* Redis

## 📄 Licença

MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🤝 Contribuição

Contribuições são bem-vindas! Por favor, abra uma issue ou envie um pull request.

## 🚧 Próximas Funcionalidades

* [ ] Dashboard de métricas de cliques
* [ ] Autenticação de usuários para gerenciamento de URLs
* [ ] Expiração automática de URLs encurtadas
* [ ] Relatórios de uso por IP
* [ ] Melhor tratamento de erros e logs detalhados

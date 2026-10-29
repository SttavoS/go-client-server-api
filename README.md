# Go Client Server API

## Configuração e Inicialização

### Pré-requisitos

- Go 1.27
- Não precisa de gcc/CGO por usar o pacote `moderc.org/sqlite`

### Passos para Configuração

1. **Clone o repositório**

```bash
git clone https://github.com/SttavoS/go-client-server-api.git
cd go-client-server-api
```

2. Rode o servidor primeiro a partir da raíz do projeto

```bash
go run ./cmd/server
```

A única URL do projeto é `http://127.0.0.1:8080/cotacao`

3. Depois rode o client, também a partir da raíz do projeto

```bash
go run ./cmd/client
```

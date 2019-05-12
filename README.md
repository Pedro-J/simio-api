# simio-api

Implementação da simio-api em GO


## Configuração da Aplicação:

**Atenção**: Para utilizar os comandos `go` deste projeto, tenha antes o go instalado na sua máquina.

## 1 - Instalação do GO 

Siga os passos do tutorial oficial da linguagem em: https://golang.org/doc/install. A versão utilizada no projeto foi a 1.12.4

## 2 - Instalação das dependencias do projeto 

```
$   cd {DIRETORIO_DO_PROJETO}
$   dep ensure
```

## 3 - Como Compilar

```
$   cd {DIRETORIO_DO_PROJETO}
$   go build
```

um binário com o nome de simio-api é gerado após esse comando

## 4 - Como Executar

```
$   cd {DIRETORIO_DO_PROJETO}
$   ./simio-api
```

## 5 - Informações da API

A api posseui dois endpoints que são:

http://localhost:5000/simian e http://localhost:5000/stats (ambiente local)

ou 

http://simio-api.us-east-2.elasticbeanstalk.com/simian e http://simio-api.us-east-2.elasticbeanstalk.com/stats (ambiente AWS)


OBS: A porta padrão da aplicação é a 5000 e arquivos com dados relacionados a aplicação serão salvos na pasta "{DIRETORIO_DO_BINARIO}/database/data/simios/" 

## 6 - Teste se a aplicação está rodando

```
$   curl -d '{"dna": ["ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"]}' -X POST http://localhost:5000/simian -w '\n'
```

ou

```
$   curl -d '{"dna": ["ATCGAT","ATCGAT","ATCGAT","TGAACC","TGGTTG","GACGGA"]}' -X POST http://localhost:5000/simian -w '\n'
```

ou


```
$   curl -X GET http://localhost:5000/stats -w '\n'
```


## 7 - Para ver a cobertura dos testes

```
$   cd {DIRETORIO_DO_PROJETO}
$   go test -cover
```

## 8 - Para ver a cobertura dos testes de forma detalhada

```
$   cd {DIRETORIO_DO_PROJETO}
$   go test ./... -coverprofile report_cover_tests.out
$   go tool cover -html=report_cover_tests.out
```
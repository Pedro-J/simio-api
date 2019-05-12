# simio-api

Implementação da simio-api em GO


## Configuração da Aplicação:

**Atenção**: Para utilizar os comandos `go` deste projeto, tenha antes o go instalado na sua máquina.

## 1 - Instalação do GO 

Siga os passos do tutorial oficial da linguagem em: https://golang.org/doc/install

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

OBS: A porta padrão da aplicação é a 8000 e arquivos com dados relacionados a aplicação serão salvos na pasta "{DIRETORIO_DO_BINARIO}/database/data/simios/" 

## 5 - Teste se a aplicação está rodando

```
$   curl -d '{"dna": ["ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"]}' -X POST http://localhost:8000/simian -w '\n'
```

ou

```
$   curl -d '{"dna": ["ATCGAT","ATCGAT","ATCGAT","TGAACC","TGGTTG","GACGGA"]}' -X POST http://localhost:8000/simian -w '\n'
```

ou


```
$   curl -X GET http://localhost:8000/stats -w '\n'
```
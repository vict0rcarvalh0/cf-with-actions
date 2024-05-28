#!/bin/sh

# Arquivo para armazenar o relatório de testes e cobertura
REPORT_FILE="test-report.txt"

# Limpa o conteúdo do arquivo de relatório anterior
> "$REPORT_FILE"

# Executa os testes e gera o arquivo de cobertura
echo "Executando testes e gerando arquivo de cobertura..." | tee -a "$REPORT_FILE"
go test -coverprofile=coverage.out ./... | tee -a "$REPORT_FILE"

# Verifica se o arquivo de cobertura foi gerado
if [ -f coverage.out ]; then
    # Gera o relatório de cobertura em HTML
    go tool cover -html=coverage.out -o coverage.html
    echo "Relatório de cobertura HTML gerado: coverage.html" | tee -a "$REPORT_FILE"
else
    echo "Falha ao gerar o arquivo de cobertura" | tee -a "$REPORT_FILE"
fi

# Gera o relatório dos testes com detalhes
echo "Gerando relatório detalhado dos testes..." | tee -a "$REPORT_FILE"
go test -v ./... | tee -a "$REPORT_FILE"

# Exibe a cobertura no terminal e adiciona ao arquivo de relatório
echo "Cobertura dos testes:" | tee -a "$REPORT_FILE"
go tool cover -func=coverage.out | tee -a "$REPORT_FILE"

echo "Relatório de testes gerado: $REPORT_FILE" | tee -a "$REPORT_FILE"

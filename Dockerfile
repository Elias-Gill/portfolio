ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

# Configurar el directorio de trabajo y copiar los archivos necesarios
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /run-app .

# Crear la imagen final
FROM debian:bookworm

# Copiar el binario generado
COPY --from=builder /run-app /usr/local/bin/

# Copiar todo el contenido del proyecto (templates, estáticos, etc.)
WORKDIR /app
COPY --from=builder /usr/src/app ./

# Comando de ejecución
CMD ["run-app"]

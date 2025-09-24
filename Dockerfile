# ---- etapa de compilación ----
FROM golang:1.23-alpine AS builder

# Instala dependencias necesarias para compilar
RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 genera binario estático
RUN CGO_ENABLED=0 go build -o /run-app .

# ---- imagen final ----
FROM alpine:3.20

# Instala git en la imagen final (lo que pediste)
RUN apk add --no-cache git

# Copia el binario
COPY --from=builder /run-app /usr/local/bin/run-app

# Copia assets (templates, estáticos…)
WORKDIR /app
COPY --from=builder /src ./

CMD ["run-app"]

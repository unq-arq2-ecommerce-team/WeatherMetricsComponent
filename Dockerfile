FROM golang:1.20.2-alpine3.17

# Haciendo set up
RUN apk update && apk upgrade && apk add ca-certificates bash git openssh gcc g++ pkgconfig build-base curl \
    && rm -rf /var/cache/apk/*

# Definir directorio de trabajo
WORKDIR /app

# Copiar archivos de configuración de módulos
COPY go.mod .
COPY go.sum .

# Descargar las dependencias del proyecto
RUN go mod download

# Copiar el código fuente del proyecto y el archivo .env
COPY . .

# update swagger docs/
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.10
RUN swag init -g internal/infrastructure/app.go

# Compilar el proyecto teniendo en cuenta la ubicación de main.go en la carpeta 'cmd'
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main ./cmd

# Ejecutar la aplicación
CMD ["./main"]
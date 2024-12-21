# Gunakan base image untuk Golang versi 1.23
FROM golang:1.23

# Set working directory
WORKDIR /app

# Copy kode ke dalam container
COPY . .

# Build aplikasi
RUN go mod tidy
RUN go build -o server .

# Ekspose port 8080
EXPOSE 8080

# Jalankan aplikasi
CMD ["./server"]

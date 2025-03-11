FROM node:23-alpine AS frontend
WORKDIR /frontend
COPY frontend/package*.json /frontend
RUN npm ci
RUN npm i @tailwindcss/oxide-linux-x64-musl
COPY frontend/ /frontend
RUN npm run build

FROM golang:1.24-alpine AS backend
WORKDIR /backend
COPY backend/go.* /backend/
RUN go mod download
COPY backend/ /backend
COPY --from=frontend /frontend/dist/client/ /backend/public/
RUN go build -o HexDeck

FROM scratch
WORKDIR /app
COPY --from=backend /backend/HexDeck .
CMD ["/app/HexDeck"]
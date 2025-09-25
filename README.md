# Тестовое задание для **АО «КАЛУГА АСТРАЛ»**

## 🚀 Запуск проекта

### 🔧 Требования
- [Docker](https://docs.docker.com/get-docker/)  
- [Docker Compose](https://docs.docker.com/compose/)

### 📋 Инструкция по запуску
1. Создайте файл `.env`, скопировав и заполнив значения из примера `.env.example`:
   ```bash 
   cp .env.example .env
   ```
2. Запустите приложение:
    ```bash 
    docker compose --env-file .env up
    ```
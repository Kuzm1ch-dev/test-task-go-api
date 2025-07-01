# Task API

REST API для управления задачами с асинхронной обработкой.

## Запуск

### Локально
```bash
go run cmd/api/main.go
```

### Docker
```bash
docker build -t task-api .
docker run -p 8080:8080 task-api
```

## API Endpoints

- `POST /tasks` - создать задачу
- `GET /tasks` - получить все задачи  
- `GET /tasks/{id}` - получить задачу по ID
- `DELETE /tasks/{id}` - удалить задачу



## Примеры использования

### Создать задачу
```bash
curl -X POST http://localhost:8080/tasks
```

### Получить задачу
```bash
curl http://localhost:8080/tasks/{task-id}
```

### Удалить задачу
```bash
curl -X DELETE http://localhost:8080/tasks/{task-id}
```

## Особенности
- Задачи обрабатываются 3-5 минут (случайно)
- Прогресс обновляется каждую секунду
- При удалении задачи горутина корректно завершается
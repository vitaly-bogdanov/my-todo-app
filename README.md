# my-todo-app

## Деплой на Render

[![Deploy to Render](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy?repo=https://github.com/vitaly-bogdanov/my-todo-app)

Проект содержит Blueprint `render.yaml`: бесплатный Docker web service и бесплатную базу Postgres.

1. Откройте ссылку https://render.com/deploy?repo=https://github.com/vitaly-bogdanov/my-todo-app и войдите в Render.
2. Проверьте имена сервиса `my-todo-app` и базы `my-todo-app-db` и нажмите **Apply** (план Free для обоих).
3. Дождитесь сборки: при старте приложение само применит миграции из `migrations/` (повторный запуск ничего не меняет), переменные `POSTGRES_*` подставляются из базы автоматически.
4. Проверьте сервис: `curl https://<имя-сервиса>.onrender.com/v1/users`.

Ограничения free-плана: бесплатная Postgres удаляется через 30 дней и в workspace может быть только одна такая база; бесплатный web service засыпает без трафика, и первый запрос после простоя идёт медленно.

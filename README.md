# app-watchdog

<img width="1903" height="616" alt="image" src="https://github.com/user-attachments/assets/2b22ccc9-2355-4549-8920-bbbefb25a9f2" />

1. [Структура репозитория](#structure)
2. [Веб‑приложение](#web-app)
3. [Watchdog‑скрипт](#watchdog)
4. [Деплой через Ansible](#ansible)
5. [Release workflow](#release)

<a id="structure"></a>
## Структура репозитория

- `app/` – golang веб-приложение
- `monitor/` –  bash‑скрипт, для мониторинга веб-приложения
- `ansible/` – инвентарь и роль `web-app`, которая ставит/удаляет приложение и мониторинг

<a id="web-app"></a>
## Веб‑приложение

Переменные окружения:

- `APP_PORT` (по умолчанию `8080`) – порт HTTP.
- `HEALTH_FAIL_RATE` (0…1, по умолчанию `0`) – вероятность искусственно «уронить» `/healthz`.
- `HEALTH_MAX_DELAY_MS` (по умолчанию `0`) – максимальная задержка ответа `/healthz` в миллисекундах.

<a id="watchdog"></a>
## Watchdog‑скрипт

Переменные окружения:

- `SERVICE_NAME` – systemd‑юнит, который нужно перезапускать, например `web-app.service`.
- `CHECK_URL` – URL health‑чека, например `http://127.0.0.1:8080/healthz`.
- `CHECK_TIMEOUT`, `RETRY_COUNT`, `RETRY_DELAY`, `CHECK_HTTP_CODE` – таймаут curl, число попыток, пауза между ними и ожидаемый HTTP‑код.

Скрипт прерывается после успешного ответа; если все попытки исчерпаны, он рестартует указанный юнит и логирует результат.

<a id="ansible"></a>
## Деплой Ansible

Роль `ansible/roles/web-app` управляет установкой:

1. Обновите `ansible/inventory/hosts.yml`, указав свои хосты.
2. Настройте `ansible/inventory/group_vars/all.yml`. Основные параметры:
   - `web_app_*`, `web_app_monitoring_*`
3. Запустите плейбук:

```bash
cd ansible
ansible-playbook -i inventory/hosts.yml site.yml
```

При состоянии `present` роль:

1. Скачивает архив релиза с `https://github.com/jeorji/app-watchdog/releases`.
2. Распаковывает его в `/opt/<имя сервиса>` и рендерит `.env` с параметрами здоровья.
3. Устанавливает systemd‑юнит для приложения (`app.service.j2`).
4. Разворачивает службу и таймер (`monitor.service.j2`, `monitor.timer.j2`), которые гоняют `monitor.sh`.

Переключение в `absent` удаляет юниты и каталоги, возвращая хост в исходное состояние.

<a id="release"></a>
## Release workflow

Файл `.github/workflows/release.yml` автоматизирует публикацию релизов. Он запускается при пуше тега `v*` и выполняет:

1. `build` job — сборка бинарника `web-app`, упаковка в `app-<os>-<arch>.tar.gz`.
2. `release` job — загрузка артефактов, генерация SHA256‑сумм, формирование GitHub Release с тарболами, `monitor.sh` и `checksums.txt`.

Чтобы создать релиз:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```

Через пару минут в разделе Releases появятся свежие артефакты, которые использует Ansible.

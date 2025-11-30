# web-app Ansible Role

Разворачивает веб‑сервис `web-app` и watchdog (`monitor.sh`) из GitHub Releases, настраивает systemd‑юниты и .env файлы.

## Role Variables

### Приложение

| Переменная | Описание | Значение по умолчанию |
| ---------- | -------- | --------------------- |
| `web_app_state` | Управление жизненным циклом (`present`/`absent`) | `absent` |
| `web_app_release_version` | Тег релиза либо `latest` для последнего артефакта | `latest` |
| `web_app_service_name` | Имя systemd‑юнита приложения | `web-app` |
| `web_app_health_max_delay_ms` | Максимальная задержка ответа `/healthz` (мс) | `1000` |
| `web_app_health_fail_rate` | Вероятность искусственной ошибки `/healthz` | `0.1` |
| `web_app_listen_port` | HTTP‑порт приложения | `8080` |
| `web_app_install_path` | Каталог, куда распаковывается архив | `/opt` |
| `web_app_repo_url` | Базовый URL GitHub Releases | `https://github.com/jeorji/app-watchdog/releases` |

### Watchdog

| Переменная | Описание | Значение по умолчанию |
| ---------- | -------- | --------------------- |
| `web_app_monitoring_state` | Жизненный цикл watchdog (`present`/`absent`) | `absent` |
| `web_app_monitoring_release_version` | Тег релиза `monitor.sh` либо `latest` | `latest` |
| `web_app_monitoring_service_name` | Имя таймера/сервиса watchdog | `web-app-monitoring` |
| `web_app_monitoring_watch_service` | Целевой systemd‑юнит, который перезапускает watchdog | **обязательно указать** |
| `web_app_monitoring_check_url` | URL health‑чека для мониторинга | **обязательно указать** |
| `web_app_monitoring_check_timeout` | Таймаут curl в секундах | `3` |
| `web_app_monitoring_retry_count` | Количество повторов проверки | `3` |
| `web_app_monitoring_retry_delay` | Пауза между попытками (секунды) | `1` |
| `web_app_monitoring_timer_interval` | Интервал systemd‑таймера | `1min` |
| `web_app_monitoring_install_path` | Каталог установки watchdog | `/opt` |
| `web_app_monitoring_repo_url` | Источник релизов `monitor.sh` | `https://github.com/jeorji/app-watchdog/releases` |

## Handlers

- `reload systemd`
- `restart web-app service`

## Example Playbook

```yaml
- hosts: app_servers
  vars:
    web_app_state: present
    web_app_monitoring_state: present
    web_app_monitoring_watch_service: web-app
    web_app_monitoring_check_url: http://127.0.0.1:8080/healthz
  roles:
    - web-app
```

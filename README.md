# git-fame

Считает, сколько строк, коммитов и файлов приходится на каждого автора в git-репозитории. Строки привязываются к авторам через `git blame`, поэтому цифры показывают состояние кода на выбранной ревизии, а не всю историю.

## Установка

Нужны Go 1.26+ и `git` в `PATH`.

```sh
git clone https://github.com/mhgffqwoer/git-fame.git
cd git-fame
go install ./cmd/git-fame
```

## Запуск

```sh
git-fame
git-fame --repository ../other-repo --revision v1.2.0
git-fame --extensions .go,.md --order-by commits
git-fame --languages go --exclude 'vendor/*,*_test.go' --format json
```

| Флаг | По умолчанию | |
| --- | --- | --- |
| `--repository` | `.` | путь к репозиторию |
| `--revision` | `HEAD` | коммит, ветка или тег |
| `--order-by` | `lines` | `lines`, `commits` или `files` |
| `--use-committer` | `false` | считать коммитеров вместо авторов |
| `--format` | `tabular` | `tabular`, `csv`, `json` или `json-lines` |
| `--extensions` | | только эти расширения, например `.go,.md` |
| `--languages` | | только эти языки, например `go,markdown` |
| `--exclude` | | пропустить файлы, подходящие под любой из glob-паттернов |
| `--restrict-to` | | оставить только файлы, подходящие хотя бы под один glob-паттерн |

Соответствие языков и расширений лежит в `configs/language_extensions.json`.

## Как считается

Список файлов берется из `git ls-tree` на нужной ревизии, проходит через фильтры, и по каждому файлу запускается `git blame --porcelain`. Для каждого автора считаются:

- строки - строки, которые он менял последним;
- коммиты - разные коммиты, к которым относятся эти строки;
- файлы - файлы, где есть хотя бы одна его строка.

У пустых файлов blame ничего не возвращает, поэтому они засчитываются автору последнего коммита, который их трогал.

## Тесты

Интеграционные тесты собирают бинарник и гоняют его на git-бандлах из `test/integration/testdata`:

```sh
go test ./...
```

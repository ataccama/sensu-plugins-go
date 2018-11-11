# Golang Sensu plugins

> Buildable assets for Sensu 5.0

Sensu [[ https://docs.sensu.io/sensu-core/2.0/reference/assets/ | asset ]] is a self-contained executable that Sensu is able to distribute to the agents who need it and then use it in checks, handlers etc. Old Ruby-based Sensu checks can't be used, because they require Ruby and shitload of gems to work properly.

The preferred way to implement checks now is to use Golang, because it fullfills all the requirements from the previous paragraph.

## Checks

- [x] Load average - `check-load`
- [x] Memory - `check-memory`
- [x] Filesystem usage `check-fs`
- [x] generic SystemD unit check - `check-systemd-unit`
- [x] TCP check - `check-tcp`
- [x] Kernel process (aka `ps ...`) - `check-ps`
- [] S.M.A.R.T. - `check-smart` //started//
- [] mdamd consistency check - `check-mdraid`
- [x] HTTP check - `check-http` //just return codes so far//
- [] ICMP check - `check-icmp`
- [] Prometheus metric - `check-promql`
- [] Docker container status - `check-docker-process`
- [] Nomad
  - `check-nomad-task`
  - `check-nomad-job`

## Handlers

- [] [[ https://github.com/sensu/sensu-slack-handler/blob/master/main.go | Slack ]]

# Building

1. Clone this repo
2. Have `go-1.11` installed and `make` too
3. Run `make`
4. Use built `.tar.gz` with your `sensuctl asset create ...` command
5. Profit :)

# Notes

- Sensu guys have made [[ https://github.com/sensu-plugins/sensu-plugins-go/ | some libs ]] to help with commonalities, those are included
- A guy who's name sounds like a random japanese Samurai already created [[ https://github.com/hico-horiuchi/sensu-plugins-go/ | some checks of his ]], maybe some of the code might be reused


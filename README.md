# Golang Sensu plugins

> Buildable assets for Sensu 5.0

Sensu [asset](https://docs.sensu.io/sensu-go/5.0/reference/assets/) is a self-contained executable that Sensu is able to distribute to the agents who need it and then use it in checks, handlers etc. Old Ruby-based Sensu checks can't be used, because they require Ruby and shitload of gems to work properly.

The preferred way to implement checks now is to use Golang, because it fullfills all the requirements from the previous paragraph.

## Checks

- [x] Load average - `check-load`
- [x] Memory - `check-memory`
- [x] Filesystem usage `check-fs`
- [x] generic SystemD unit check - `check-systemd-unit`
- [x] TCP check - `check-tcp`
- [x] Kernel process (aka `ps ...`) - `check-ps`
- [ ] S.M.A.R.T. - `check-smart` _started_
- [ ] mdamd consistency check - `check-mdraid`
- [x] HTTP check - `check-http` _just return codes so far_
- [ ] ICMP check - `check-icmp`
- [ ] Prometheus metric - `check-promql`
- [ ] Docker container status - `check-docker-process`
- [ ] Nomad
  - `check-nomad-task`
  - `check-nomad-job`

## Handlers

- [ ] [Slack](https://github.com/sensu/sensu-slack-handler/)

# Building

1. Clone this repo
2. Have `go-1.11` installed and `make` too
3. Run `make`
4. Use built `.tar.gz` with your `sensuctl asset create ...` command
5. Profit :)

# Notes

- Sensu guys have made [some libs](https://github.com/sensu-plugins/sensu-plugins-go/) to help with commonalities, those are included
- A guy who's name sounds like a random japanese Samurai already created [some checks of his](https://github.com/hico-horiuchi/sensu-plugins-go/) - it might be worth exploring


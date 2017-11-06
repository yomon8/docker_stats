
### Install

```sh
chmod +x ./docker_stats
sudo mv ./docker_stats /usr/share/munin/plugins/
sudo ln -s /usr/share/munin/plugins/docker_stats /etc/munin/plugins/docker_stats
```

### Adjust user permission 

```sh
cat <<EOF > /etc/munin/plugin-conf.d/docker_stats
[docker_stats]
user root
EOF
```

```
sudo /etc/init.d/munin-node restart
```

### Remove Container from graph

If you want to remove old container from munin graph.

```
rm ${MUNIN_PLUGSTATE}/containerlist
```



### Test

```sh
sudo -u munin munin-run docker_stats config
sudo -u munin munin-run docker_stats
```

```sh
$ nc localhost 4949
# munin node at yourhostname
config docker_stats
...
...
fetch docker_stats
...
...
```

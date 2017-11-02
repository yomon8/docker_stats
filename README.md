
### Install

```sh
chmod +x ./docker_stats
sudo mv ./docker_stats /usr/share/munin/plugins/
sudo ln -s /usr/share/munin/plugins/docker_stats /etc/munin/plugins/docker_stats
```

### Adjust user permission 

```
cat <<EOF > /etc/munin/plugin-conf.d/docker_stats
[docker_stats]
user root
EOF
```


### Test

```
sudo -u munin munin-run docker_stats config
sudo -u munin munin-run docker_stats
```

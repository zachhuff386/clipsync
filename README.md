# clipsync: Clipboard sync

[Clipsync](https://github.com/zachhuff386/clipsync)
is a high performance clipboard sharing application for linux. The clipboard is
shared over the network with support for string clipboard content up to several
megabytes in size. Works with both Wayland and X servers.

## Building

Requires Golang, X11 and XFixes libraries.

```bash
sudo dnf install libX11-devel libXfixes-devel xorg-x11-proto-devel
CGO_ENABLED=1 go install github.com/zachhuff386/clipsync@latest
```

## Computer 111.111.111.111

```bash
./clipsync generate-key
public_key=oDKlgt9NsfJObrQu+Xp2GTLY80EpGkW0Hr09bBwsUTI

tee ./clipsync.conf << EOF
{
  "bind": "0.0.0.0:9774",
  "private_key": "3zyjOTKaXO0zFOIx2cOaeYBmQ8bSsQjTr9dLGBHTNto",
  "public_key": "cjZ5vR4R2t3QI8xzMz0Jw2lGvnim3nBlsmiViyM0iWo",
  "clients": [
    {
      "address": "222.222.222.222:9774",
      "public_key": "C4HZ1DkOIbG3u2zqC4mL8JPhliOfjex0h3E3XoKfJhw"
    }
  ]
}
EOF

./clipsync start
```

## Computer 222.222.222.222

```bash
./clipsync generate-key
public_key=C4HZ1DkOIbG3u2zqC4mL8JPhliOfjex0h3E3XoKfJhw

tee ./clipsync.conf << EOF
{
  "bind": "0.0.0.0:9774",
  "private_key": "Vc7BVAyVFdtmvtOv5uhm/2/EoAZlOXvsL/QgCUVlVAg",
  "public_key": "C4HZ1DkOIbG3u2zqC4mL8JPhliOfjex0h3E3XoKfJhw",
  "clients": [
    {
      "address": "111.111.111.111:9774",
      "public_key": "oDKlgt9NsfJObrQu+Xp2GTLY80EpGkW0Hr09bBwsUTI"
    }
  ]
}
EOF

./clipsync start
```

## Systemd User Service

```bash
sudo cp ./clipsync /usr/local/bin/clipsync
cp ./clipsync.conf ~/.config/clipsync.conf

mkdir -p ~/.config/systemd/user
tee ~/.config/systemd/user/clipsync.service << EOF
[Unit]
Description=Clipsync
Wants=gnome-session.target
After=gnome-session.target

[Service]
ExecStart=/usr/local/bin/clipsync start $HOME/.config/clipsync.conf

[Install]
WantedBy=gnome-session.target
EOF

systemctl --user daemon-reload
systemctl --user enable clipsync
systemctl --user restart clipsync
```

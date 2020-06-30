# 通过 readlink 获取绝对路径，再取出目录
work_path=$(dirname $(readlink -f $0))

systemctl stop saynice-api
systemctl stop saynice-web

mv $work_path/saynice-api-amd64-linux-v1 /usr/local/bin/saynice-api
mv $work_path/saynice-web-amd64-linux /usr/local/bin/saynice-web
mv $work_path/saynice-api.service /lib/systemd/system
mv $work_path/saynice-web.service /lib/systemd/system

rm -rf $work_path

systemctl daemon-reload
systemctl enable saynice-api
systemctl enable saynice-web

systemctl start saynice-api
systemctl start saynice-web
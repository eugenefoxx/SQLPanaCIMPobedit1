# SQLPanaCIMPobedit1

# SQLPanacimP1 query to DB PanaCIM P1

Нужно перед выпуском сверить состояние ео в sap и panacim

systemd
/etc/systemd/system/panasap.service
[Unit]
Description=App

[Service]
ExecStart=/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/panasap
WorkingDirectory=/home/a20272/Code/github.com/eugenefoxx/SQLPanaCIMPobedit1/

[Install]
Wantedby=multi-user.target

/etc/systemd/system/panasap.timer
[Unit]
Description=timer for panasap

[Timer]
OnCalendar=02:50 \*:0/15

[Install]
WantedBy=timers.target

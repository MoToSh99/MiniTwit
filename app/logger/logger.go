package logger

import (
        "github.com/bshuster-repo/logrus-logstash-hook"
        "github.com/sirupsen/logrus"
        "net"
)

func Send(msg string) {
        log := logrus.New()
        conn, err := net.Dial("tcp", "localhost:5000")
        if err != nil {
                log.Fatal(err)
        }
        hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{}))

        log.Hooks.Add(hook)
        log.Info(msg)
}
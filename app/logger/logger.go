package logger

import (
        "github.com/bshuster-repo/logrus-logstash-hook"
        "github.com/sirupsen/logrus"
        "net"
)

func Send(msg string) {
        log := logrus.New()
        conn, err := net.Dial("tcp", "165.22.76.211:5000")
        if err != nil {
                log.Warn(err)
                return
        }
        hook := logrustash.New(conn, logrustash.DefaultFormatter(logrus.Fields{}))

        log.Hooks.Add(hook)
        log.Info(msg)
}
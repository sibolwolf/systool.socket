package audioctrller

import (
	"bytes"
	"log"
	"testing"
	"time"

	"smartconn.cc/tosone/logstash"
)

func TestRecord(t *testing.T) {
	reader, err := Record()
	if err != nil {
		logstash.Error(err.Error())
	}
	var byteBuf = new(bytes.Buffer)
	go func() {
		byteBuf.ReadFrom(reader)
		log.Println("over read")
	}()
	<-time.After(time.Second * 6)
	if err := PlaySE(byteBuf.Bytes()); err != nil {
		logstash.Error(err.Error())
	}
}

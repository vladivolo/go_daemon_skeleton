package redis

import (
	"errors"
	log "github.com/vladivolo/lumber"
	"github.com/vladivolo/redismq"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RedisConn struct {
	proto     string
	ipport    string
	queuename string
	conn      rmq.Connection
	queue     rmq.Queue
	db        string
	lock      *sync.Mutex
}

type LConsumer struct {
	callback func(string)
}

//tcp@127.0.0.1:6379/1/new_data_block
func ParceRedisQueueAddr(Params string) (proto, ipaddr, db, queue string, err error) {
	arr := strings.Split(Params, "@")
	if len(arr) != 2 {
		log.Error("ParceRedisQueueAddr: (%s) Failed format", Params)
		return "", "", "", "", errors.New("Failed full redis queue format")
	}
	in := strings.Split(arr[1], "/")
	if len(in) != 3 {
		log.Error("ParceRedisQueueAddr: (%s) Failed format", Params)
		return "", "", "", "", errors.New("Failed full redis queue format")
	}

	return arr[0], in[0], in[1], in[2], nil
}

func NewRedisConn(proto, ipaddr, db string) (*RedisConn, error) {
	rc := &RedisConn{
		proto:  proto,
		ipport: ipaddr,
		conn:   nil,
		queue:  nil,
		db:     db,
		lock:   &sync.Mutex{},
	}

	if err := rc.Connect(); err != nil {
		log.Error("RedisConn error %s", err.Error())
		return nil, err
	}

	return rc, nil
}

func NewRedisConnection(Params string) (*RedisConn, error) {
	proto, ipaddr, db, queue, err := ParceRedisQueueAddr(Params)
	if err != nil {
		return nil, err
	}

	log.Info("AddQueue: proto %s ipaddr %s db %s queue %s", proto, ipaddr, db, queue)

	rc, err := NewRedisConn(proto, ipaddr, db)
	if err != nil {
		log.Error("Redis connect error: %s", err)
		return nil, err
	}

	rc.queuename = queue

	return rc, nil
}

func (rc *RedisConn) Connect() (err error) {
	rc.conn, err = rmq.OpenConnection("producer", rc.Proto(), rc.IpPort(), rc.Db())
	if err != nil {
		return errors.New("Connect failed")
	}

	return
}

func (rc *RedisConn) OpenQueue() {
	rc.queue = rc.Cli().OpenQueue(rc.QueueName())
	/*
		if r := rc.queue.PurgeReady(); r != 0 {
			log.Debug("PurgeReady: %d", r)
		}
	*/
}

func (rc *RedisConn) QueueName() string {
	return rc.queuename
}

func (rc *RedisConn) Proto() string {
	return rc.proto
}

func (rc *RedisConn) IpPort() string {
	return rc.ipport
}

func (rc *RedisConn) Db() int {
	db, err := strconv.Atoi(rc.db)
	if err != nil {
		log.Error("Unknown convert Db %s", err)
		return -1
	}
	return db
}

func (rc *RedisConn) Cli() rmq.Connection {
	return rc.conn
}

func (rc *RedisConn) Queue() rmq.Queue {
	return rc.queue
}

func (rc *RedisConn) Publish(payload string) error {
	q := rc.Queue()
	if q == nil {
		log.Warn("Publish failed: %s", payload)
		return errors.New("Publish failed")
	}
	log.Info("Publish: %s", payload)

	return q.Publish(payload)
}

func (rc *RedisConn) StartConsuming(callback func(string)) error {
	q := rc.Queue()
	if q == nil {
		return errors.New("Queue is nil")
	}

	InputConsumer := LConsumer{callback: callback}
	if err := q.StartConsuming(1*time.Second, &InputConsumer); err != nil {
		log.Warn("StartConsuming() return false")
		return err
	}

	/*
		if err := q.StartConsuming(1, 1*time.Second); err != nil {
			log.Warn("StartConsuming() return false")
			return err
		}

		InputConsumer := LConsumer{callback: callback}
		rc.Queue().AddConsumer("Consumer", &InputConsumer)
	*/

	return nil
}

func (c *LConsumer) Consume(delivery rmq.Delivery) {
	log.Info("Consume: Queue [%s] data %s", delivery.QueueName(), delivery.Payload())

	c.callback(delivery.Payload())
	delivery.Ack()
}

package redis

import (
	"fmt"
	"testing"
	"time"
)

var (
	FFParams = "tcp@127.0.0.1:6370/1/test_new_data_block"
	FQParams = "tcp@127.0.0.1:6370/1/test_new_data_block"

	IQParams = "tcp@127.0.0.1:6379/1/test_new_data_block_2"
	OQParams = "tcp@127.0.0.1:6379/1/test_new_data_block_2"

	MsgCnt int
)

func TestRedisConnection(t *testing.T) {
	/*
		// Connect with failed params
		FFQueue, err := NewRedisConnection(FFParams)
		if err == nil {
			t.Error("Connect with failed server success")
		}

		if FFQueue != nil {
			t.Error("FFQueue != nil")
		}

		// Connect to unknown redis server
		FailedQueue, err := NewRedisConnection(FQParams)
		if err == nil {
			t.Error("Connect to failed server success")
		}

		if FailedQueue != nil {
			t.Error("FailedQueue != nil")
		}
	*/
	// Up input queue
	InputQueue, err := NewRedisConnection(IQParams)
	if err != nil {
		t.Error("NewRedisConnection failed")
	}

	InputQueue.OpenQueue()
	if err = InputQueue.StartConsuming(InputQueueHandle); err != nil {
		t.Error("StartConsuming failed")
	}

	// Up output queue
	OutputQueue, err := NewRedisConnection(OQParams)
	if err != nil {
		t.Error("NewRedisConnection failed")
	}

	OutputQueue.OpenQueue()

	for i := 0; i < 10; i++ {
		if err = OutputQueue.Publish(fmt.Sprintf("%d-%x", i, i)); err != nil {
			fmt.Printf("%s", err)
		}
		//fmt.Printf("%s: %d", OutputQueue.QueueName(), OutputQueue.Queue().ReadyCount())
	}

	time.Sleep(5 * time.Second)

	if MsgCnt != 10 {
		t.Error("Send & Recieve message failed!")
	}
}

func InputQueueHandle(payload string) {
	fmt.Printf("InputQueueHandle %s\n", payload)

	MsgCnt++
}

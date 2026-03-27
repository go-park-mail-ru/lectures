package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"image/jpeg"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/nfnt/resize"
)

type ImgResizeTask struct {
	Name string
	MD5  string
}

const (
	ImageResizeTopicName = "cool_topic"
)

var (
	kafkaAddr = flag.String("addr", "localhost:9092", "kafka broker address")
	groupID   = flag.String("group", "resizer-group", "kafka consumer group")
)

var (
	sizes = []uint{80, 160, 320}
)

type Consumer struct {
	ready chan bool
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("incoming task %s\n", string(msg.Value))

		task := &ImgResizeTask{}
		err := json.Unmarshal(msg.Value, task)
		if err != nil {
			fmt.Println("cant unpack json", err)
			sess.MarkMessage(msg, "")
			continue
		}

		originalPath := fmt.Sprintf("./images/%s.jpg", task.MD5)
		for _, size := range sizes {
			resizedPath := fmt.Sprintf("./images/%s_%d.jpg", task.MD5, size)
			err := ResizeImage(originalPath, resizedPath, size)
			if err != nil {
				fmt.Println("resize failed", err)
				continue
			}
			fmt.Println("resize success")
			time.Sleep(3 * time.Second)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func main() {
	flag.Parse()

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := Consumer{
		ready: make(chan bool),
	}

	client, err := sarama.NewConsumerGroup([]string{*kafkaAddr}, *groupID, config)
	panicOnError("cant create consumer group", err)
	defer client.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	go func() {
		defer wg.Done()

		for {
			if err := client.Consume(ctx, []string{ImageResizeTopicName}, &consumer); err != nil {
				fmt.Println("Error from consumer:", err)
				time.Sleep(time.Second)
			}

			consumer.ready = make(chan bool)

			// Check if context was cancelled, signaling that the consumer should stop
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	<-consumer.ready
	fmt.Println("worker started")
	<-ctx.Done()
	fmt.Println("worker stopped")
	wg.Wait()
}

func ResizeImage(originalPath string, resizedPath string, size uint) error {
	file, err := os.Open(originalPath)
	if err != nil {
		return fmt.Errorf("cant open file %s: %s", originalPath, err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		return fmt.Errorf("cant jpeg decode file %s", err)
	}

	resizeImage := resize.Resize(size, 0, img, resize.Lanczos3)

	out, err := os.Create(resizedPath)
	if err != nil {
		return fmt.Errorf("cant create file %s: %s", resizedPath, err)
	}
	defer out.Close()

	jpeg.Encode(out, resizeImage, nil)

	return nil
}

// Никогда так не делайте!
func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + " " + err.Error())
	}
}

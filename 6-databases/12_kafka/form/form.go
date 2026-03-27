package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/IBM/sarama"
)

var uploadFormTmpl = []byte(`
<html>
	<body>
	<form action="/upload" method="post" enctype="multipart/form-data">
		Image: <input type="file" name="my_file">
		<input type="submit" value="Upload">
	</form>
	</body>
</html>
`)

type ImgResizeTask struct {
	Name string
	MD5  string
}

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write(uploadFormTmpl)
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	uploadData, handler, err := r.FormFile("my_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer uploadData.Close()

	fmt.Fprintf(w, "handler.Filename %v\n", handler.Filename)
	fmt.Fprintf(w, "handler.Header %#v\n", handler.Header)

	tmpName := RandStringRunes(32)

	tmpFile := "./images/" + tmpName + ".jpg"
	newFile, err := os.Create(tmpFile)
	if err != nil {
		http.Error(w, "cant open file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	hasher := md5.New()
	writtenBytes, err := io.Copy(newFile, io.TeeReader(uploadData, hasher))
	if err != nil {
		http.Error(w, "cant save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	newFile.Sync()
	newFile.Close()

	md5Sum := hex.EncodeToString(hasher.Sum(nil))

	realFile := "./images/" + md5Sum + ".jpg"
	err = os.Rename(tmpFile, realFile)
	if err != nil {
		http.Error(w, "cant rename file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data, _ := json.Marshal(ImgResizeTask{handler.Filename, md5Sum})

	fmt.Println("put task ", string(data))

	// Kafka publish via sarama
	msg := &sarama.ProducerMessage{
		Topic: ImageResizeTopicName,
		Key:   sarama.StringEncoder(md5Sum),
		Value: sarama.ByteEncoder(data),
	}
	kafkaAsyncP.Input() <- msg

	fmt.Fprintf(w, "Upload %d bytes successful\n", writtenBytes)
}

const (
	ImageResizeTopicName = "cool_topic"
)

var (
	kafkaAddr   = flag.String("addr", "localhost:9092", "kafka broker address")
	kafkaAsyncP sarama.AsyncProducer
)

func main() {
	err := os.Mkdir("./images", 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Println(err)
		return
	}

	flag.Parse()

	config := sarama.NewConfig()
	config.Producer.Return.Successes = false
	config.Producer.Return.Errors = true

	kafkaAsyncP, err = sarama.NewAsyncProducer([]string{*kafkaAddr}, config)
	panicOnError("cant create kafka producer", err)
	defer kafkaAsyncP.Close()

	// Отдельная горутина для ошибок продюсера
	go func() {
		for err := range kafkaAsyncP.Errors() {
			fmt.Println("Failed to produce message", err)
		}
	}()

	http.HandleFunc("/", mainPage)
	http.HandleFunc("/upload", uploadPage)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

// Никогда так не делайте!
func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + " " + err.Error())
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

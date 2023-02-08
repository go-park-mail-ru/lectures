package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
это тест на проверку того что у нас это действительно конвейер
неправильное поведение: накапливать результаты выполнения одной функции, а потом слать их в следующую.

это не похволяет запускать на конвейере бесконечные задачи
правильное поведение: обеспечить беспрепятственный поток
*/
func TestPipeline(t *testing.T) {

	var ok = true
	var recieved uint32
	freeFlowCmds := []cmd{
		cmd(func(in, out chan interface{}) {
			out <- 1
			time.Sleep(10 * time.Millisecond)
			currRecieved := atomic.LoadUint32(&recieved)
			// в чем тут суть
			// если вы накапливаете значения, то пока вся функция не отрабоатет - дальше они не пойдут
			// тут я проверяю, что счетчик увеличился в следующей функции
			// это значит что туда дошло значение прежде чем текущая функция отработала
			if currRecieved == 0 {
				ok = false
			}
		}),
		cmd(func(in, out chan interface{}) {
			for range in {
				atomic.AddUint32(&recieved, 1)
			}
		}),
	}
	stat = Stat{}
	RunPipeline(freeFlowCmds...)

	assert.True(t, ok,
		"во второй джобе не увеличислся счетчик, а в первой уже дошли до следующего действия")
	assert.NotEqual(t, 0, recieved,
		"счетчик recieved в итоге не увилился, а должен был")
}

/*
этот тест проверяет то, что все функции действительно выполнились
и дает представление о влиянии time.Sleep в одном из звеньев конвейера на время работы

возможно кому-то будет легче с ним
при правильной реализации ваш код конечно же должен его проходить
*/
func TestPipeline2(t *testing.T) {

	var recieved uint32
	freeFlowCmds := []cmd{
		cmd(func(in, out chan interface{}) {
			out <- uint32(1)
			out <- uint32(3)
			out <- uint32(4)
		}),
		cmd(func(in, out chan interface{}) {
			for val := range in {
				out <- val.(uint32) * 3
				time.Sleep(time.Millisecond * 100)
			}
		}),
		cmd(func(in, out chan interface{}) {
			for val := range in {
				fmt.Println("collected", val)
				atomic.AddUint32(&recieved, val.(uint32))
			}
		}),
	}

	timeStart := time.Now()

	stat = Stat{}
	RunPipeline(freeFlowCmds...)

	expectedTime := time.Millisecond * 350
	timeEnd := time.Since(timeStart)
	assert.Less(t, timeEnd, expectedTime,
		"execition too long. Got: %s. Expected: <%s", timeEnd.String(), expectedTime.String())
	assert.Equal(t, uint32((1+3+4)*3), recieved,
		"f3 have not collected inputs, recieved = %d", recieved)
}

// инициализация джобы, которая просто выплюнет в out подряд все из слайса строк strs
func newCatStrings(strs []string, pauses time.Duration) func(in, out chan interface{}) {
	return func(in, out chan interface{}) {
		for _, email := range strs {
			out <- email
			if pauses != 0 {
				time.Sleep(pauses)
			}
		}
	}
}

// инициализация джобы, которая считает из in все строки, пока канал не закроется. и положит все в strs
func newCollectStrings(strs *[]string) func(in, out chan interface{}) {
	return func(in, out chan interface{}) {
		for dataRaw := range in {
			data := fmt.Sprintf("%v", dataRaw)
			*strs = append(*strs, data)
		}
	}
}

// проверяем, что SelectUsers корректно обрабатывает алиасы и не повторяет одних и тех же юзеров
func TestAlias(t *testing.T) {
	inputData := []string{
		"batman@mail.ru", //is an alias for bruce.wayne@mail.ru
		"bruce.wayne@mail.ru",
	}
	expectedOutput := []string{
		"{12499983457589032104 bruce.wayne@mail.ru}",
	}

	testResult := []string{}
	stat = Stat{}
	RunPipeline(
		cmd(newCatStrings(inputData, 0)),
		cmd(SelectUsers),
		cmd(newCollectStrings(&testResult)),
	)

	assert.Equal(t, expectedOutput, testResult,
		"итоговый результат отличается от ожидаемого")
}

// проверяем, что запуски функций SelectUsers,SelectMessages,CheckSpam в параллельных RunPipeline не влияют друг на друга
func TestParallelPiplines(t *testing.T) {
	inputData := []string{
		"1000@mail.ru",
		"1001@mail.ru",
		"1002@mail.ru",
	}

	stat = Stat{}
	cntFirst, cntSecond := 0, 0

	timeStart := time.Now()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunPipeline(
			cmd(newCatStrings(inputData, 150*time.Millisecond)),
			cmd(SelectUsers),
			cmd(SelectMessages),
			cmd(CheckSpam),
			cmd(func(in, out chan interface{}) {
				for _ = range in {
					cntFirst++
				}
			}),
		)

	}()

	RunPipeline(
		cmd(newCatStrings(inputData, 100*time.Millisecond)),
		cmd(SelectUsers),
		cmd(SelectMessages),
		cmd(CheckSpam),
		cmd(func(in, out chan interface{}) {
			for _ = range in {
				cntSecond++
			}
		}),
	)

	wg.Wait()

	expectedTime := 2700 * time.Millisecond
	timeEnd := time.Since(timeStart)
	assert.Less(t, timeEnd, expectedTime,
		"параллельные пайплайны не должны влиять друг на друга. скорость их выполнения зависит от самого медленного")
	assert.Equal(t, 22, cntFirst)
	assert.Equal(t, 22, cntSecond)
}

func TestTotal(t *testing.T) {
	inputData := []string{
		"harry.dubois@mail.ru",
		"k.kitsuragi@mail.ru",
		"d.vader@mail.ru",
		"noname@mail.ru",
		"e.musk@mail.ru",
		"spiderman@mail.ru", //is an alias for peter.parker@mail.ru
		"red_prince@mail.ru",
		"tomasangelo@mail.ru",
		"batman@mail.ru", //is an alias for bruce.wayne@mail.ru
		"bruce.wayne@mail.ru",
	}
	expectedOutput := []string{
		"true 221945221381252775",
		"true 357347175551886490",
		"true 1595319133252549342",
		"true 1877225754447839300",
		"true 4652873815360231330",
		"true 5108368734614700369",
		"true 7829088386935944034",
		"true 8065084208075053255",
		"true 9323185346293974544",
		"true 10463884548348336960",
		"true 11204847394727393252",
		"true 12026159364158506481",
		"true 12386730660396758454",
		"true 12556782602004681106",
		"true 12728377754914798838",
		"true 13245035231559086127",
		"true 14107154567229229487",
		"true 16476037061321929257",
		"true 16728486308265447483",
		"true 17087986564527251681",
		"true 17259218828069106373",
		"true 17696166526272393238",
		"false 26236336874602209",
		"false 59892029605752939",
		"false 221962074543525747",
		"false 378045830174189628",
		"false 2803967521226628027",
		"false 6652443725402098015",
		"false 7594744397141820297",
		"false 9656111811170476016",
		"false 10167774218733491071",
		"false 10462184946173556768",
		"false 10493933060383355848",
		"false 10523043777071802347",
		"false 11512743696420569029",
		"false 12792092352287413255",
		"false 12975933273041759035",
		"false 14498495926778052146",
		"false 15161554273155698590",
		"false 15262116397886015961",
		"false 15728889559763622673",
		"false 15784986543485231004",
	}

	timeStart := time.Now()
	testResult := []string{}
	stat = Stat{}
	RunPipeline(
		cmd(newCatStrings(inputData, 0)),
		cmd(SelectUsers),
		cmd(SelectMessages),
		cmd(CheckSpam),
		cmd(CombineResults),
		cmd(newCollectStrings(&testResult)),
	)

	expectedTime := 3000 * time.Millisecond
	timeEnd := time.Since(timeStart)
	assert.Less(t, timeEnd, expectedTime,
		"слишком долгоe выполнение. что-то где-то нераспараллелено. должно быть не больше, чем %s, а было %s", timeEnd, expectedTime)
	assert.Equal(t, expectedOutput, testResult,
		"итоговый результат отличается от ожидаемого")
	expectedStat := Stat{
		RunGetUser:            uint32(10),
		RunGetMessages:        uint32(5),
		GetMessagesTotalUsers: uint32(9),
		RunHasSpam:            uint32(42),
	}
	assert.Equal(t, stat, expectedStat, "количество вызово функций не совпадает с ожидаемым")
}

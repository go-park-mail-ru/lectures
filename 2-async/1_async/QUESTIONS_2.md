# Мини-квиз: пул, гонки, select, sync, context

По коду: `3_workerpool` и всё после него.

---

### 1. Воркеры «не финишат». Почему?

```go
workerInput := make(chan string)
for i := 0; i < 3; i++ {
	go startWorker(i, workerInput) // внутри: for input := range in { ... }
}
for _, month := range months {
	workerInput <- month
}
// close(workerInput) // попробуйте закомментировать
time.Sleep(time.Second)
```

- A. Без `close` воркеры сами выйдут из `range`, когда main перестанет писать
- B. Без `close` `range` в воркерах блокируется навсегда; `finished` не напечатается (процесс потом убьёт `Sleep`/выход из main)
- C. Паника: нельзя `range` по каналу из нескольких горутин
- D. Нужен буфер размером с число воркеров, иначе дедлок на send

<details>
<summary>Ответ</summary>

**B.**

Несколько горутин **можно** читать из одного канала — это и есть пул. Останавливает их `close`: `range` дочитает хвост и выйдет.

Без `close` воркеры вечно висят на `<-in`. Программа кажется живой из-за `time.Sleep` в `main`, но это не join. Как только main закончится — процесс умрёт вместе с ними.
</details>

---

### 2. Гонка, которой «не видно глазами»

```go
var counters = map[int]int{}
for i := 0; i < 5; i++ {
	go func(counters map[int]int, th int) {
		for j := 0; j < 5; j++ {
			counters[th*10+j]++
		}
	}(counters, i)
}
fmt.Scanln()
fmt.Println(counters)
```

- A. Всё ок: у каждой горутины свои ключи, пересечений нет
- B. Тихий data race / `fatal error: concurrent map writes`; ловится `go run -race`
- C. Паника только если ключи совпадут
- D. Map в Go потокобезопасен, как `sync.Map`

<details>
<summary>Ответ</summary>

**B.**

Даже **разные ключи** не спасают: внутренности map нельзя трогать параллельно. Писатель+писатель — часто сразу `concurrent map writes`. Писатель+читатель — data race, `-race` орёт, без флага результат неопределённый.

`++` у обычного `int` тоже гонка: это read-modify-write, не одна инструкция. Для счётчика — `atomic.AddInt32` / mutex. Для map — `sync.Mutex` / `RWMutex` (RLock **только** на чтение; инкремент ключа — это запись, нужен `Lock`).
</details>

---

### 3. `select` != `if/else if`

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)
ch1 <- 1
ch2 <- 1

select {
case val := <-ch1:
	fmt.Println("ch1", val)
case ch2 <- 2:
	fmt.Println("put to ch2")
default:
	fmt.Println("default")
}
```

- A. Всегда `ch1 1` — кейсы проверяются сверху вниз
- B. Всегда `default`
- C. Случайно либо `ch1 1`, либо `put to ch2` — оба кейса готовы
- D. Всегда `ch1 1`: у `ch2` буфер уже полный, запись **не** готова

<details>
<summary>Ответ</summary>

**D.** Напечатается `ch1 1`.

В `ch1` лежит значение — чтение готово. В `ch2` тоже лежит `1`, буфер размера 1 **заполнен**, поэтому `ch2 <- 2` ждать некого: этот case **не** участвует.

Это не приоритет «сверху вниз». Если убрать `ch2 <- 1`, буфер пуст — запись станет готовой, и runtime выберет случайно между двумя case. `default` срабатывает только когда **никто** не готов.
</details>

---

### 4. Трюк: закрытый канал в `select` крутится вечно

```go
ch1 := make(chan int)
ch2 := make(chan int)
go func() {
	close(ch2)
	ch1 <- 1
	close(ch1)
}()

for {
	select {
	case v1, ok := <-ch1:
		if !ok {
			fmt.Println("ch1 closed")
			// ch1 = nil  // забыли
		}
		fmt.Println("chan1", v1)
	case v2, ok := <-ch2:
		if !ok {
			fmt.Println("ch2 closed")
			ch2 = nil
		}
		fmt.Println("chan2", v2)
	}
}
```

- A. После `close` кейс больше не выбирается — цикл сам остановится
- B. Закрытый `ch1` **всегда готов**, отдаёт `0, false` → busy-loop; лечится `ch1 = nil`
- C. Паника на чтении закрытого канала
- D. `ok == false` только один раз, дальше кейс исчезает сам

<details>
<summary>Ответ</summary>

**B.**

Закрытый канал в `select` готов **всегда**. Без `ch = nil` этот case выигрывает снова и снова: `v=0`, `ok=false`, CPU в потолок.

`nil`-канал в `select` наоборот: case **выключается**, никогда не выбирается. Так ждут «пока оба не закроются».

Мелочь из лекции: `fmt.Println("chan1", v1)` стоит **после** `if !ok` — на закрытии ещё и нулевое значение напечатается.
</details>

---

### 5. `WaitGroup`: где вызывать `Add`?

```go
wg := &sync.WaitGroup{}
for i := 0; i < 5; i++ {
	go func() {
		wg.Add(1)
		defer wg.Done()
		doWork()
	}()
}
wg.Wait()
```

- A. Нормально: `Add` рядом с работой
- B. Гонка: `Wait` может увидеть счётчик `0` до того, как горутины успеют сделать `Add`
- C. Паника: `Add` нельзя вызывать из другой горутины никогда
- D. Нужен `Add(5)` внутри каждой горутины

<details>
<summary>Ответ</summary>

**B.**

`Add` вызывают **в той горутине, которая порождает воркеров**, до `go ...` (или `wg.Add(n)` один раз перед циклом). Иначе `Wait` в main может пробежаться, пока ни одна горутина ещё не стартовала — счётчик ноль, `Wait` сразу вернётся.

`Done` обычно в `defer`. `Add` после старта `Wait` или отрицательный счётчик — паника.

С Go 1.25 есть `wg.Go(func() { ... })` — сам делает `Add`+`go`+`Done`. Ошибки он всё равно не собирает: для этого `errgroup`. У `errgroup.WithContext` первая ошибка ещё и отменяет остальных через `ctx`.
</details>

---

### 6. Остановить пачку горутин: send или `close`?

```go
cancelCh := make(chan struct{})
for i := 1; i <= 3; i++ {
	go worker(i, cancelCh) // select { case <-cancelCh: return ... }
}
// вариант 1:
cancelCh <- struct{}{}
// вариант 2:
close(cancelCh)
```

- A. Оба будят всех воркеров одинаково
- B. `send` разбудит **одного**; `close` — **всех** (и все последующие `<-cancelCh` тоже сразу готовы)
- C. `close` разбудит одного, send — всех
- D. `close` здесь паника: из канала читают несколько горутин

<details>
<summary>Ответ</summary>

**B.**

Запись в канал — один получатель. Остальные так и будут крутиться. Если в этот момент никто не читает, а канал небуферизованный — ещё и дедлок в main.

`close` — рассылка «сигнал всем»: каждый `<-cancelCh` просыпается. Читать из закрытого можно из многих горутин. Писать после `close` / закрывать дважды — паника. Поэтому закрывает **один** владелец, обычно тот, кто создал канал.

Тот же приём: `signal.Notify` + `close(done)` по Ctrl+C.
</details>

---

### 7. Кто тут вообще отменяется?

```go
parent, cancelParent := context.WithCancel(context.Background())
child, _ := context.WithCancel(parent)

go worker(child) // select { case <-ctx.Done(): ... }

time.Sleep(time.Second)
cancelParent()
```

- A. Воркер не остановится: ему передали `child`, а `cancel` ребёнка мы выбросили (`_`)
- B. `cancelParent` закроет `Done()` и у parent, и у child — воркер проснётся
- C. Отмена родителя отменяет только тех, кому передали сам `parent`
- D. Паника: нельзя выбрасывать `cancel`

<details>
<summary>Ответ</summary>

**B.**

Контексты — дерево. `WithCancel(parent)` подписывает ребёнка на родителя: `cancelParent()` закрывает `Done()` у parent **и** у всех потомков. Воркеру всё равно, что в руках `child`, а не `parent`.

`child, _ := ...` значит только: **точечно** этого ребёнка вы уже не отмените, не трогая остальных. На отмену сверху это не влияет. Обратное неверно: `cancel` ребёнка родителя не трогает.

`cancel` от `WithCancel`/`WithTimeout` всё равно принято звать (`defer cancel()`), чтобы отпустить таймер и не ждать, пока контекст сам истечёт.
</details>

---

### 8. `init()` и `Once.Do` — это одно и то же?

```go
func Init() {
	fmt.Println("Init once")
}

func init() {
	fmt.Println("Init at start of program")
}

func main() {
	once := &sync.Once{}
	for i := 0; i < 10; i++ {
		go func() { once.Do(Init) }()
	}
}
```

- A. Да: оба один раз на пакет, просто разный синтаксис
- B. `init()` — сам, до `main`; `Once.Do` — лениво, один раз на этот `once`, даже из 10 горутин
- C. `Once.Do` вызовет `Init` по разу в каждой горутине — у них разные стеки
- D. Если `Do` никто не вызовет, `init()` тоже не выполнится

<details>
<summary>Ответ</summary>

**B.**

`init()` пакета рантайм зовёт сам, один раз при старте, **до** `main`. Вы его не вызываете.

`sync.Once` — лениво, при первом `Do`. Десять горутин, один `Init`: остальные `Do` сразу вернутся. Если `Do` так и не вызвали — `Init` не бежит, а `init()` пакета уже отработал.

Поэтому в примере сначала `"Init at start of program"`, потом `"Init once"`.
</details>

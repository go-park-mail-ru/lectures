curl http://localhost:8080/debug/pprof/goroutine?debug=2 -o goroutines.txt
curl http://127.0.0.1:8080/debug/pprof/heap -o mem_out.txt

# go tool pprof -svg -inuse_objects pprof_2.exe mem_out.txt > mem_io.svg
# go tool pprof -svg -inuse_space pprof_2.exe mem_out.txt > mem_is.svg

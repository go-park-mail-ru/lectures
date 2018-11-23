curl http://127.0.0.1:8080/debug/pprof/heap -o mem_out.txt
curl http://127.0.0.1:8080/debug/pprof/profile?seconds=5 -o cpu_out.txt

go tool pprof -svg -alloc_objects pprof_1.exe mem_out.txt > mem_ao.svg
go tool pprof -svg -alloc_space pprof_1.exe mem_out.txt > mem_as.svg
go tool pprof -svg -inuse_objects pprof_1.exe mem_out.txt > mem_io.svg
go tool pprof -svg -inuse_space pprof_1.exe mem_out.txt > mem_is.svg
go tool pprof -svg pprof_1.exe cpu_out.txt > cpu.svg

# go tool pprof pprof_1.exe mem_out.txt
# go tool pprof pprof_1.exe cpu_out.txt

# go tool pprof -http=:8081 -alloc_objects pprof_1.exe mem_out.txt
# go tool pprof -http=:8081 pprof_1.exe cpu_out.txt

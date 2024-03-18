FROM tarantool/tarantool
COPY app.lua /opt/tarantool
CMD ["tarantool", "/opt/tarantool/app.lua"]

# docker build --tag=tnt:latest . && docker run --rm -p 3301:3301 tnt
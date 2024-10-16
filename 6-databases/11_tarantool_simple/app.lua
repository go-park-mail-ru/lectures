#!/usr/bin/env tarantool

-- Запустить с помощью ./app.lua при заранее установленным tarantool на устрйотсво
-- Зайти с помощью tarantoolctl connect 127.0.0.1:3301

-- Настроить базу данных
box.cfg {
    listen = 3301
}

-- При поднятии БД создаем спейсы и индексы
box.once('init', function()
	s = box.schema.space.create('users')
	s:format({{name = 'id', type = 'unsigned'},{name = 'name', type = 'string'},{name = 'age', type = 'unsigned'}})
	s:create_index('primary', {type = 'hash', parts = {'id'}})


    box.schema.space.create('sessions')
    box.space.sessions:create_index('primary',
        { type = 'HASH', parts = {1, 'string'}})

     print('Hello, world!')
end)

-- Даем доступ для юзера guest для подключения по guest пользователю
box.once('access:v1', function()
    box.schema.user.grant('guest', 'read,write,execute', 'universe')
end)

-- Можем определять свои функции и вызывать их из кода
function test()
    print('test')
    return 'test'
end

-- Например ID сессии генерируется здесь
function new_session(user_data)
    print('received data', user_data)
    local random_number
    local session_id
    session_id = ""
    for x = 1,64,1 do
        random_number = math.random(65, 90)
        session_id = session_id .. string.char(random_number)
    end

    box.space.sessions:insert{session_id, user_data}

    return session_id
end

function check_session(session_id)
    local session_id = box.space.sessions:select{session_id}[1]
    print('found session', session_id)
    return session_id
end

How to run:
- generate *.pb.go:\
  '*.\protoc-gen.cmd*'
- build project:\
  '*go build -v ./cmd/server*'
- run postgres DB p 5432:\
  from root project dir execute - '*docker-compose up*'\
  tables will be created once server start
- start:\
  *server.exe*\
  or just run by GoLang IDE

Notes:
- configuration file is not implemented
- mock db for tests is not implemented
- sometimes some tests have additional spaces in message response :(
- didn't read go project structure
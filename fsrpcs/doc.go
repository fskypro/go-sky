package fsrpcs

/*
service.go
	表示服务器的一个服务，客户端与服务器连接后，一个连接可以对应很多个 service
client.go
	与 service 对应，一个 service 对应一个 client

serverproxy.go
	与一个连接对应，是客户端与服务器交流的一个代理，一个 serverproxy 可以对应多个 client
*/

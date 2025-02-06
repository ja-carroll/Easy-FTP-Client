.PHONY: run-ftp-container
run-ftp-container:
	docker run -d -v ./ftpserver:/home/vsftpd -p 20:20 -p 21:21 -p 47400-47470:47400-47470 -e FTP_USER=test -e FTP_PASS=test -e PASV_ADDRESS=127.0.0.1 --name ftp --restart=always bogem/ftp

.PHONY: stop
stop:
	docker stop ftp

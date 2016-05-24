all: binary

binary: 
	rm -f ./bin/*
	GOARCH=amd64 GOOS=linux go build -ldflags "-X main.buildstamp=`date '+%Y-%m-%d_%H:%M:%S'` -X main.githash=`git rev-parse HEAD`" -o ./bin/gocf

cloudpush: binary
	cf push  

localpush: binary  
	docker run -v ${PWD}/bin:/opt/bin --link mariadb:mariadb  --env-file ./cf.env -p 4000:4000  \
	   	-it cloudfoundry/cflinuxfs2 /opt/bin/gocf

mariadb-start:
	docker run -d --name mariadb --env-file ./mariadb.env  -p 3306:3306/tcp mariadb  2>/dev/null || echo "MariaDB is already running (make db-stop to start from scratch)"
		sleep 10

mariadb-stop:
	docker rm -f  mariadb || exit 0

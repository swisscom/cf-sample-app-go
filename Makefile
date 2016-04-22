all: binary

binary: 
	rm -f ./bin/*
	GOARCH=amd64 GOOS=linux go build -ldflags "-X main.buildstamp=`date '+%Y-%m-%d_%H:%M:%S'` -X main.githash=`git rev-parse HEAD`" -o ./bin/gocf

cloudpush: binary
	cf push  

localpush: binary 
	docker run -v ${PWD}/bin:/opt/bin  --env-file ./cf.env -p 4000:4000  \
	   	-it cloudfoundry/cflinuxfs2 /opt/bin/gocf

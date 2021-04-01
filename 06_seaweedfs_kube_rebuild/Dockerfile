FROM scratch
ADD weed /
ENTRYPOINT ["/weed"]

#go get github.com/chrislusf/seaweedfs/weed/...
#CGO_ENABLED=0 GOOS=linux go build github.com/chrislusf/seaweedfs/weed
#docker build -t weed:latest -f ./Dockerfile.
#docker tag weed:latest repo/weed:latest
#docker push repo/weed:latest
#rm -f weed
#docker rmi $(docker images -qa -f 'dangling=true')2>/dev/null
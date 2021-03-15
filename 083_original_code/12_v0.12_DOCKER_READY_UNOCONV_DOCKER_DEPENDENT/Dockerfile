#docker image prune --filter label=stage=builder -f
#docker-compose up --build -d
#CGO_ENABLED=0 GOOS=linux go build -a
# 1 choose a compiler OS
FROM golang:alpine AS builder
# Set necessary environmet variables needed for our image

# 2 (optional) label the compiler image
LABEL stage=builder
# 3 (optional) install any compiler-only dependencies
RUN apk add --no-cache gcc libc-dev
WORKDIR /workspace
# 4 copy all the source files
COPY . .
# 5 build the GO program
RUN CGO_ENABLED=0 GOOS=linux go build -a
# 6 choose a runtime OS
FROM alpine AS final
# 7
ARG ENV
WORKDIR /
# 8 copy from builder the GO executable file
COPY --from=builder /workspace/12_v0.12_SERVER_READY_direct_file_serving .
#COPY --from=builder /workspace/_envs/env_$ENV.yaml ./_envs/
# 9 execute the program upon start
CMD [ "./12_v0.12_SERVER_READY_direct_file_serving" ]




#############################################################################################

#ENV PYTHONUNBUFFERED=1
#RUN apk add --update --no-cache python3 && ln -sf python3 /usr/bin/python
#RUN python3 -m ensurepip
#RUN pip3 install --no-cache --upgrade pip setuptools
#RUN pip install unoconv
#RUN apk update
#RUN apk add --upgrade py3-unoconv
RUN apk add --upgrade --no-cache mupdf-tools
#RUN apk add --upgrade --no-cache libreoffice


# Export necessary port
EXPOSE 8080
EXPOSE 8100
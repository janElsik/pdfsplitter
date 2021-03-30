# 1 Setup and start Seaweed (File System)
 - git clone https://github.com/chrislusf/seaweedfs
 - cd seaweedfs/docker
 - docker-compose -f seaweedfs-compose.yml -p seaweedfs up
# 2 Setup and start Nats (Messaging)
 - docker network create nats
 - docker run --name nats --network nats --rm -p 4222:4222 -p 8222:8222 nats
# 3 Setup and start Gotenberg (Office doc conversion)
 - docker run --rm -p 3000:3000 thecodingmachine/gotenberg:6
# 4 Run Go services
 - inside 02_convert_to_pdf run ./02_convert_to_pdf
 - inside 03_create_thumbs run ./03_create_thumbs
 - inside 01_entry_point run ./01_entry_point
 
  

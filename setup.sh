#启动redis-server
redis-server ./conf/redis.conf

#启动fastdfs - tracker
fdfs_trackerd ./conf/tracker.conf restart
#启动fastdfs - storage
fdfs_storaged /root/workspace/go/src/iHome/conf/storage.conf restart

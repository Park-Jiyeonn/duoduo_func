## 项目框架：Hertz
### 极简版抖音
Data中的视频我没上传，因为没什么意义  
#### 目前的结构
.  
├── README.md  
├── biz  
│&emsp;&emsp;├── dal  
│&emsp;&emsp;│&emsp;&emsp;├── comment.go  
│&emsp;&emsp;│&emsp;&emsp;├── init.go  
│&emsp;&emsp;│&emsp;&emsp;├── user.go  
│&emsp;&emsp;│&emsp;&emsp;└── video.go  
│&emsp;&emsp;├── handler  
│&emsp;&emsp;│&emsp;&emsp;├── comment  
│&emsp;&emsp;│&emsp;&emsp;│&emsp;&emsp;└── comment_service.go  
│&emsp;&emsp;│&emsp;&emsp;├── favorite  
│&emsp;&emsp;│&emsp;&emsp;│&emsp;&emsp;└── favorite_service.go  
│&emsp;&emsp;│&emsp;&emsp;├── feed  
│&emsp;&emsp;│&emsp;&emsp;│&emsp;&emsp;└── feed_service.go  
│&emsp;&emsp;│&emsp;&emsp;├── ping.go  
│&emsp;&emsp;│&emsp;&emsp;├── publish  
│&emsp;&emsp;│&emsp;&emsp;│&emsp;&emsp;└── publish_service.go  
│&emsp;&emsp;│&emsp;&emsp;└── user  
│&emsp;&emsp;│&emsp;&emsp;    └── user_service.go  
│&emsp;&emsp;├── model  
│&emsp;&emsp;│&emsp;&emsp;├── comment  
│&emsp;&emsp;│&emsp;&emsp;│&emsp;&emsp;└── comment.go  
│&emsp;&emsp;│&emsp;&emsp;├── favorite  
│&emsp;&emsp;│&emsp;&emsp;│&emsp;&emsp;└── favourite.go  
│&emsp;&emsp;│&emsp;&emsp;├── feed  
│&emsp;&emsp;│&emsp;&emsp;│&emsp;&emsp;└── feed.go  
│&emsp;&emsp;│&emsp;&emsp;├── publish  
│&emsp;&emsp;│&emsp;&emsp;│&emsp;&emsp;└── publish.go  
│&emsp;&emsp;│&emsp;&emsp;└── user  
│&emsp;&emsp;│&emsp;&emsp;    └── user.go  
│&emsp;&emsp;├── mw  
│&emsp;&emsp;│&emsp;&emsp;└── jwt.go  
│&emsp;&emsp;├── redis  
│&emsp;&emsp;│&emsp;&emsp;├── init.go  
│&emsp;&emsp;│&emsp;&emsp;└── redis.go  
│&emsp;&emsp;└── router  
│&emsp;&emsp;    ├── comment  
│&emsp;&emsp;    │&emsp;&emsp;├── comment.go  
│&emsp;&emsp;    │&emsp;&emsp;└── middleware.go  
│&emsp;&emsp;    ├── favorite  
│&emsp;&emsp;    │&emsp;&emsp;├── favourite.go  
│&emsp;&emsp;    │&emsp;&emsp;└── middleware.go  
│&emsp;&emsp;    ├── feed  
│&emsp;&emsp;    │&emsp;&emsp;├── feed.go  
│&emsp;&emsp;    │&emsp;&emsp;└── middleware.go  
│&emsp;&emsp;    ├── publish  
│&emsp;&emsp;    │&emsp;&emsp;├── middleware.go  
│&emsp;&emsp;    │&emsp;&emsp;└── publish.go  
│&emsp;&emsp;    ├── register.go  
│&emsp;&emsp;    └── user  
│&emsp;&emsp;        ├── middleware.go  
│&emsp;&emsp;        └── user.go  
├── data  
│&emsp;&emsp;├── 第一个视频.jpg  
│&emsp;&emsp;├── 第一个视频.mp4  
│&emsp;&emsp;├── 第二个视频.jpg  
│&emsp;&emsp;├── 第二个视频.mp4  
│&emsp;&emsp;├── 高铁上.jpg  
│&emsp;&emsp;└── 高铁上.mp4  
├── go.mod  
├── go.sum  
├── idl  
│&emsp;&emsp;├── comment.thrift  
│&emsp;&emsp;├── favourite.thrift  
│&emsp;&emsp;├── feed.thrift  
│&emsp;&emsp;├── publish.thrift  
│&emsp;&emsp;└── user.thrift  
├── main.go  
├── main_go.exe  
├── pojo  
│&emsp;&emsp;├── comment.go  
│&emsp;&emsp;├── user.go  
│&emsp;&emsp;└── video.go  
├── router.go  
├── router_gen.go  
├── test  
│&emsp;&emsp;├── test.go  
│&emsp;&emsp;└── test_go.exe  
└── util  
├── ffmpeg.go  
└── jwt.go  
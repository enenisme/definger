package utils

// ProbesContent 探针服务内容
var ProbesContent = `[probes]
  [probes.093561eda8a835f5a01738826c77dbf6]
    data = "GET / HTTP/1.1\r\nSec-Ch-Ua: \"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"\r\nSec-Ch-Ua-Mobile: ?0\r\nSec-Fetch-Site: none\r\nSec-Fetch-Mode: navigate\r\nSec-Fetch-User: ?1\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36\r\nAccept-Encoding: gzip, deflate, br\r\nPriority: u=0, i\r\nSec-Ch-Ua-Platform: \"macOS\"\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7\r\nSec-Fetch-Dest: document\r\nAccept-Language: zh-CN,zh;q=0.9,ga;q=0.8,en;q=0.7\r\nConnection: close\r\n\r\n"
    desc = "请补充描述"
    format = "HTTP-RESPONSE"
    timeout = 30
    write_expression = "http(data,3)"
  [probes.6e2c4c3c59397d0806c6490a3368face]
    data = "GET /index HTTP/1.1\r\nUser-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36\r\nContent-Type: application/json\r\n\r\n{\"a\":\""
    desc = "请补充描述"
    format = "HTTP-RESPONSE"
    timeout = 30
    write_expression = "http(data,3)"
  [probes.91c52f3fc61db0797ab5212014b3cf9f]
    data = "GET /login HTTP/1.1\r\nAccept-Encoding: gzip, deflate, br\r\nPriority: u=0, i\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36\r\nSec-Fetch-Dest: document\r\nAccept-Language: zh-CN,zh;q=0.9,ga;q=0.8,en;q=0.7\r\nConnection: close\r\nSec-Ch-Ua-Platform: \"macOS\"\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7\r\nSec-Fetch-Site: none\r\nSec-Ch-Ua: \"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"\r\nSec-Ch-Ua-Mobile: ?0\r\nSec-Fetch-Mode: navigate\r\nSec-Fetch-User: ?1\r\n\r\n"
    desc = "请补充描述"
    format = "HTTP-RESPONSE"
    timeout = 30
    write_expression = "http(data,3)"
`

// probesForGetTitle 获取标题探针服务
var ProbesForGetTitle = `[probes.093561eda8a835f5a01738826c77dbf6]
    data = "GET / HTTP/1.1\r\nSec-Ch-Ua: \"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"\r\nSec-Ch-Ua-Mobile: ?0\r\nSec-Fetch-Site: none\r\nSec-Fetch-Mode: navigate\r\nSec-Fetch-User: ?1\r\nUpgrade-Insecure-Requests: 1\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36\r\nAccept-Encoding: gzip, deflate, br\r\nPriority: u=0, i\r\nSec-Ch-Ua-Platform: \"macOS\"\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7\r\nSec-Fetch-Dest: document\r\nAccept-Language: zh-CN,zh;q=0.9,ga;q=0.8,en;q=0.7\r\nConnection: close\r\n\r\n"
    desc = "请补充描述"
    format = "HTTP-RESPONSE"
    timeout = 30
    write_expression = "http(data,3)"`

// ProbesForGetFavicon 获取favicon探针服务
var ProbesForGetFavicon = `[probes.favicon]
    data = "GET /favicon.ico HTTP/1.1\r\nSec-Ch-Ua: \"Google Chrome\";v=\"125\", \"Chromium\";v=\"125\", \"Not.A/Brand\";v=\"24\"\r\nSec-Ch-Ua-Platform: \"macOS\"\r\nAccept-Encoding: gzip, deflate, br\r\nConnection: close\r\nUpgrade-Insecure-Requests: 1\r\nSec-Fetch-Site: none\r\nSec-Ch-Ua-Mobile: ?0\r\nUser-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36\r\nSec-Fetch-User: ?1\r\nSec-Fetch-Dest: document\r\nAccept-Language: zh-CN,zh;q=0.9,ga;q=0.8,en;q=0.7\r\nPriority: u=0, i\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7\r\nSec-Fetch-Mode: navigate\r\n\r\n"
    desc = "favicon"
    format = "HTTP-RESPONSE"
    timeout = 30
    write_expression = "http(data,3)"`

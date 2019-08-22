数据库操作服务
---

将数据操作封装成 RESTful API, 方便远程访问数据.

1. 认证配置

    config/application.json
    
    ```json
    {
     "web":{
       "port": "127.0.0.1:9001",
       "isDebug": true
     },
     "db":{
       "dbName": {
         "type": "mysql",
         "dsn": "dbuser:dbpassword@tcp(host:port)/dbname?charset=utf8",
         "maxOpenConnections": 1,
         "maxIdleConnections": 1
       }
     }
    }
    ```

    config/auth.json:
    
    ```json
    {
      "$appKey": {
        "secret": "$appSecret",
        "database": "dbname",
        "actions": ["SELECT", "INSERT", "UPDATE", "DELETE"]
      }
    }    
    ```
    $appKey, $appSecret: 请使用32位字符串替换

2. token生成
    
    使用 AES-128-CBC 加解密码
    
    ```
    // PHP 加密方法
    
    $nonce = md5(time());
    $appKey = '...'; // config/auth.json中的$appKey
    $dbName = '...'; // config/application.json中的dbName
 
    // 待加密的原始字符串
    $msg = $nonce .'|'. $appKey .'|'. time();
    
    $key = '***'; // 16位字符串,  config/auth.json中$appSecret的前18位
    $iv = '***'; // 16位字符串,  config/auth.json中$appSecret的后18位
        
    $cipher = 'AES-128-CBC';
    $ivlen = openssl_cipher_iv_length($cipher);
    $token = openssl_encrypt($msg, $cipher, $key, $options=0, $iv);
    ```
  
3. API

    * 3.1. 查询
        ```
        POST /query
        Authorization: token
        ContentType: application/json
        {
            "sql":"INSERT user(username,email) VALUES(?,?)",
            "params":["test1","test1@test.com"1]
        }
        ```
        返回
        ```json
        [
         {"email":"test1@test.com","id":"1","username":"test1"}
        ]
        ```
    
    * 3.2. 更新
        ```
        POST /update
        Authorization: token
        ContentType: application/json
        {
            "sql":"UPDATE user SET username=? WHERE id=?",
            "params":["test", 1]
        }
        ```
       返回
       ```json
       {"affectedRows":1}
       ```

    * 3.3. 删除
        ```
        POST /delete
        Authorization: token
        ContentType: application/json
        {
            "sql":"DELETE FROM user WHERE id=?",
            "params":[1]
        }
        ```
        返回
        ```json
        {"affectedRows":1}
        ```

    * 3.4. 添加
        ```
        POST /insert HTTP/1.1
        Authorization: token
        Content-Type: application/json
        {
            "sql":"INSERT user SET id=?,username=?",
            "params":[1, "test"]
        }
        ```
        返回(需要数据支持)
        ```json
        {"lastInsertId":1}
        ```

    * 3.5. 批量查询 (可以是不同的SQL,不同的查询参数)
        ```
        POST /batch/query
        Authorization: token
        ContentType: application/json
        [
            {
                "sql":"SELECT * FROM user WHERE id=?",
                "params":[1]
            },
            {
                "sql":"SELECT * FROM user WHERE username=?",
                "params":["test2"]
            } 
        ]
        ```
        返回(多维数组)
        ```json
        [
         [{"email":"test1@test.com","id":"1","username":"test1"}],
         [{"email":"test2@test.com","id":"2","username":"test2"}]
        ]
        ```

    * 3.6. 事务 v1 (可以是不同的SQL,不同的查询参数)
        ```
        POST /v1/transaction
        Authorization: token
        ContentType: application/json
        [
            {
                "sql":"INSERT INTO user(username,email) VALUES(?,?)",
                "params":["test1","test1@test.com"]
            },
            {
                "sql":"INSERT INTO user_info(username,email) VALUES(?,?)",
                "params":["test2","test2@test.com"]
            } 
        ]
        ```
        返回
        ```json
        {"result":true}
        ```

    * 3.7. 事务 v2 (同一个SQL, 不同的参数)
        ```
        POST /v2/transaction
        Authorization: token
        ContentType: application/json
        {
            "sql":"INSERT INTO user(username,email) VALUES(?,?)",
            "params":[
                ["test2","test2@test.com"],
                ["test2","test2@test.com"]
            ]
        }
        ```
        返回
        ```json
        {"result":true}
        ```

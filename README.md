# Rss-Aggregator

> 什么是RSS?

[RSS是一种消息来源格式规范，用于聚合多个网站更新的内容并自动通知网站订阅者。使用RSS后，网站订阅者便无需再手动查看网站是否有新的内容，同时RSS可将多个网站更新的内容进行整合，以摘要的形式呈现，有助于订阅者快速获取重要信息，并选择性地点阅查看。](https://zh.wikipedia.org/zh-tw/RSS)[1](https://zh.wikipedia.org/zh-tw/RSS)[2](https://zh.wikipedia.org/wiki/RSS)[3](https://baike.baidu.com/item/RSS)

> RSS Aggregator能干什么?

[RSS Aggregator是一种工具，它可以自动从您喜欢的博客和网站中筛选出相关内容，并将其呈现在一个地方。这个地方可以是网站、移动应用程序或桌面应用程序。如果您想轻松地筛选内容、保持组织和节省时间，那么RSS聚合器是一个很好的工具。如果您是网站发布者，RSS聚合器可以帮助您将内容进行同步并扩大受众。如果您是读者，RSS聚合器可以帮助您将来自您喜欢的网站的所有内容整合到一个地方。](https://www.wprssaggregator.com/what-is-an-rss-aggregator/)

接下来将使用`go`语言来构建一个`Rss-Aggregator`。

## 1、初始化项目

在初始化项目中，我们将编写一个简单的（不带任何`handler`）的http服务器。

- 使用`go mod`命令

`注:`[`go mod`是Go语言从1.11版本之后官方推出的版本管理工具，用于解决之前没有地方记录依赖包具体版本的问题，相比vendor、dep等包管理工具，更加便于依赖包的管理。](https://juejin.cn/post/7022280943608004644)[1](https://juejin.cn/post/7022280943608004644)[2](https://zhuanlan.zhihu.com/p/413040181)[3](http://c.biancheng.net/view/5712.html)

```bash
go mod init github.com/blkcor/Rss-Aggregator
```

- 安装依赖

```bash
go get github.com/joho/godotenv    -- 用于从.env文件中读取环境变量
go get github.com/go-chi/chi       -- api router
go get github.com/go-chi/cors      -- 用于进行跨域配置
```

- 配置文件

```properties
PORT=8000
```

- `http`服务

```go
package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	//load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("an error happen when loading the .env file: %v", err)
	}
	//now we can get the PORT attr in the current environment
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT is not found in the environment")
	}
	//router
	router := chi.NewRouter()
	//register handler
	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{"http://*", "https://*"},
			AllowedMethods: []string{"POST", "GET", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
			ExposedHeaders: []string{"*"},
      MaxAge:        300
		}),
	)
	//http server
	serve := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
	log.Println("server is running on port:", port)
	err = serve.ListenAndServe()
	if err != nil {
		log.Fatalf("an error happen when starting the http server: %v", err)
	}
}
```

- Json工具类(`json.go`)

使用该工具类使`response`以`json`格式进行响应。

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithJson(w http.ResponseWriter, code int, payload interface{}) {
	//get the json string
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to Marshal json response: %v", payload)
		w.WriteHeader(500)
		return
	}
	//write the response data
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err = w.Write(dat)
	if err != nil {
		log.Print("An error happen when responsing the data: %v", err)
		return
	}
}
```

- `HandlerFunc`

我们使用handler来处理用户的每一次请求。(`handler_readiness`)

```go
package main

import "net/http"

/*
*
the function signature is the specific function signature if you want to define the http handler
*/
func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	responseWithJson(w, 200, struct{}{})
}
```

- 挂载（注册）handler

```go
.....
//v1 router
	v1Router := chi.NewRouter()
	v1Router.HandleFunc("/ready", handlerReadiness) //we can use Get Post Delete Put to replace the handleFunc(all can)
	router.Mount("/v1", v1Router)

	//http server
	serve := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}
.....
```

使用`http client`工具，如thunder, postman,insominia等等进行验证。

- 错误处理

错误处理可以让无论开发者还是用户在发生错误的时候都能清楚的知道错误发生的原因。

在`json.go`继续添加如下的函数:

```go
func responseWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX err: %v", msg)
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	responseWithJson(w, code, ErrorResponse{
		Error: msg,
	})
}
```

错误处理`handler`(`handler_error.go`)

```go
package main

import "net/http"

func handlerError(w http.ResponseWriter, r *http.Request) {
	responseWithError(w, 500, "something went wrong!")
}
```

注册`handler`

```go
v1Router.Get("/err", handlerError) 
```

## 2、数据库

本项目使用`PostgreSql`作为数据库，使用`PG-admin`作为客户端工具连接数据库。

为了通过`go`操作`database`，安装下面的依赖：

```shell
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest --将sql映射成安全的go代码
go install github.com/pressly/goose/v3/cmd/goose@latest --管理数据表的工具
```

### 2.1、使用goose创建table

- 根目录下创建`sql/schema`目录，并在下面创建`001_users.sql`

由于`goose`使基于`sql comments`，因此我们这样编写我们的sql文件:(注意逗号)

```sql
-- +goose Up
create TABLE users(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
```

- 在`.env`文件中新增数据库配置

```properties
DB_URL=postgre://chenzilong:Woshiguaiwu123@localhost:5432/rssagg
```

- 进入`sql/schema`目录下执行下面的命令：

```bash
DB_URL=postgre://chenzilong:Woshiguaiwu123@localhost:5432/rssagg up(or down)
```

goose会自动帮我们执行所有sql中up(or down)部分的sql语句。

### 2.2、使用sqlc完成从sql到go code的转换

- sqlc.yaml

该配置文件对`shema(how to define table)`和`query(how to query/modify the data)`的位置进行了配置，让sqlc进行读取生成。

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "sql/queries"
    schema: "sql/schema"
    gen:
      go:
        out: "internal/database"
```

- 根目录下创建`sql/queries`目录，并创建`user.sql`文件

```sql
-- name: CreateUser :one
Insert Into users(id,created_at,updated_at,name)
values ($1,$2,$3,$4)
RETURNING *;
```

`其中`:name指定生成的func的名称，one表示返回单个User对象。

- 执行下面的命令，sqlc会在根目录下新建`internal/database`，并将生成的`.go`文件放入其中:

```bash
sqlc generate
```

### 2.3、如何使用生成的go code

由于要随机给`user`设置id，我们需要生成uuid，添加如下依赖：

```shell
go get github.com/google/uuid
```

新增`CreateUser`的`handler`:(`handler_user.go`)

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (a *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}
	user, err := a.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
    CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error creating user:%v", err))
		return
	}
	responseWithJson(w, 200, user)
}

```

在`main.go`中新增如下代码：

- 新增一个结构体，用来存放数据库连接（`Queries`）

```go
type apiConfig struct {
	DB *database.Queries
}
```

- 新建创建用户的路由：

```go
v1Router.Post("/users", apiCfg.handlerCreateUser)
```

使用客户端工具测试。

## 3、鉴权

为了保证数据的安全，很多操作只对登录鉴权的用户开放，为了简化，我们使用ApiKey来模拟这个过程（一般通过`token`进行鉴权）。

- 在`/sql/schema`目录下新增`002_users_apikey.sql`文件，添加内容如下：

```sql
-- +goose Up
ALTER TABLE users
    ADD COLUMN api_key VARCHAR(64) UNIQUE NOT NULL DEFAULT (
        encode(sha256(random()::text::bytea), 'hex')
        );

-- +goose Down
ALTER TABLE users DROP COLUMN api_key
```

这个`sql`的主要功能是给`users`表添加一个新列`api_key`

执行下面命令：

```bash
goose postgres postgres://chenzilong:Woshiguaiwu123@localhost:5432/rssagg up
```

- 修改`sql/queries/users.sql`文件如下：

```sql
-- name: CreateUser :one
Insert Into users(id, created_at, updated_at, name, api_key)
values ($1, $2, $3, $4, encode(sha256(random()::text::bytea), 'hex')) RETURNING *;

-- name: GetUserByApiKey :one
SELECT * FROM users WHERE api_key = $1;
```

这个`sql`主要是在创建用户的时候随机设置`api_key`的值。并且新增`GetUserByApiKey`这个方法。

- 通过`sqlc`生成代码：

```shell
sqlc generate
```

- 新建文件`/internal/auth/auth.go`，内容如下：

```go
package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*
GetApiKey function could extract ApiKey in the request header
example like this:
Authorization: ApiKey {insert your api key here}
*/
func GetApiKey(header http.Header) (string, error) {
	val := header.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth header")
	}
	return vals[1], nil
}
```

这个方法从`Request Header`中拿到`ApiKey`，并规定了`ApiKey`的格式。

- 修改`handler_user.go`文件，新增如下内容：

```go
func (a *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		responseWithError(w, 403, fmt.Sprintf("Auth err: %v", err))
		return
	}
	user, err := a.DB.GetUserByApiKey(r.Context(), apiKey)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error Getting User: %v", err))
		return
	}
	responseWithJson(w, 200, dbUserToUser(user))
}
```

除了登录和注册逻辑，其他几乎所有的`api`调用都要进行身份验证。

## 4、获取资源

### 4.1、编写中间件

和`handlerGetUser`一样，在每次获取资源之前，也会进行身份验证，我们考虑将这个功能抽象成一个中间件。

- 新建`middleware_auth.go`文件，内容如下：

```go
package main

import (
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/auth"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"net/http"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (a *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			responseWithError(w, 403, fmt.Sprintf("Auth err: %v", err))
			return
		}
		user, err := a.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			responseWithError(w, 401, fmt.Sprintf("Error Getting User: %v", err))
			return
		}
		handler(w, r, user)
	}
}
```

这样我们的每一个需要进行身份验证的`handler`函数都会进行`ApiKey`的校验。

- 修改`handler_user.go`，内容如下：

```go
func (a *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	responseWithJson(w, 200, dbUserToUser(user))
}
```

- 修改`main.go`，内容如下：用`middlewareAuth`进行包裹：

```go
...
v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
...
```

### 4.2、获取资源api

下面编写获取资源内容的api，流程基本和之前一致。

- 新建`/sql/schema/003_feeds.sql`，内容如下：

```sql
-- +goose Up
create TABLE feeds
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP   NOT NULL,
    updated_at TIMESTAMP   NOT NULL,
    name       TEXT        NOT NULL,
    url        TEXT UNIQUE NOT NULL,
    user_id    UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
```

- 执行下面命令，生成数据库表：

```bash
goose postgres postgres://chenzilong:Woshiguaiwu123@localhost:5432/rssagg up
```

- 新建`/sql/queries/feeds.sql`，内容如下：

```sql
-- name: CreateFeed :one
Insert Into feeds(id, created_at, updated_at, name, url, user_id)
values ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;
```

- 执行下面命令，生成`CreateFeed`函数：

```shell
sqlc generate
```

接着编写`feed`的`handler`:

- 新建`handler_feeds.go`，内容如下：

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (a *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}
	feed, err := a.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error creating feed:%v", err))
		return
	}
	responseWithJson(w, 200, dbFeedToFeed(feed))
}

func (a *apiConfig) handlerGetFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := a.DB.GetFeeds(r.Context())
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error getting feeds:%v", err))
		return
	}
	feedsResponse := make([]Feed, 0)
	for _, feed := range feeds {
		feedsResponse = append(feedsResponse, dbFeedToFeed(feed))
	}
	responseWithJson(w, 200, feedsResponse)
}
```

- 将`handler`注册到`router`中：在`main.go`中新增内容如下：

```go
...
//feeds
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeed)
...
```

### 4.3、订阅资源api

我们需要把`user`和`feed`进行关联。

- 新建`sql/schema/feed_follows.sql`，内容如下：

```sql
-- name: CreateFeedFollows :one
Insert Into feed_follows(id, created_at, updated_at, user_id, feed_id)
values ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetFeedFollows :many
SELECT *
FROM feed_follows
WHERE user_id = $1;

-- name: DeleteFeedFollows :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
```

- 执行下面命令，生成相应的方法。

```shell
sqlc generate
```

- 新建`handler_feed_follows.go`，编写处理逻辑：

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func (a *apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	param := parameter{}
	err := decoder.Decode(&param)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}
	feedFollows, err := a.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    param.FeedID,
	})
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Fail to create FeedsFollows: %v", err))
		return
	}
	responseWithJson(w, 200, dbFeedFollowToFeedFollow(feedFollows))
}

func (a *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := a.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Fail to get FeedsFollows: %v", err))
		return
	}
	feedFollowResponse := make([]database.FeedFollow, 0)
	for _, feedFollow := range feedFollows {
		feedFollowResponse = append(feedFollowResponse, feedFollow)
	}
	responseWithJson(w, 200, feedFollowResponse)
}

func (a *apiConfig) handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		FeedID uuid.UUID
	}
	param := parameter{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing JSON:%v", err))
		return
	}
	err = a.DB.DeleteFeedFollows(r.Context(), database.DeleteFeedFollowsParams{
		UserID: user.ID,
		FeedID: param.FeedID,
	})
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error deleting feed_follows:%v", err))
	}
	responseWithJson(w, 200, "delete successfully!")
}
```

- 添加`handler`到`router`

```go
//feed_follows
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))
```

### 4.4、更新资源api

在这一步我们要通过`feeds`表中的`url`字段进行`Rss`资源的抓取。

- 新建`sql/schema/006_posts.sql`内容如下：

```sql
-- +goose Up
create TABLE posts
(
    id           UUID PRIMARY KEY,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL,
    title        TEXT      NOT NULL,
    description  TEXT,
    published_at TIMESTAMP NOT NULL,
    url          TEXT      NOT NULL UNIQUE,
    feed_id      UUID      NOT NULL REFERENCES feeds (id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE posts;
```

- 执行下面命令，生成数据库`table`：

```bash
goose postgres postgres://chenzilong:Woshiguaiwu123@localhost:5432/rssagg up
```

- 新建`/sql/queries/posts.sql`，内容如下：

```sql
-- name: CreatePost :one
Insert Into posts(id, created_at, updated_at, title, description, published_at, url, feed_id)
values ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: GetPostForUser :many
SELECT posts.*
FROM posts
         JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC LIMIT $2;
```

- 执行下面命令，生成对应操作数据库的方法：

```shell
sqlc generate
```

首先，我们需要在服务启动的时候开一个`go routine`去跑一个定时任务：

- 新建`scraper.go`，内容如下：

```go
package main

import (
	"context"
	"database/sql"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

/*
Run a long task to auto update feed resources
*/
func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scraping on %v go routine every %v duration", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	//waiting for ticker
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Fail to fetch feeds: %v", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapFeed(db, wg, feed)
		}
		wg.Wait()
	}
}
func scrapFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Fail to mark feed as fetched: %v", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		if !(err == io.EOF) {
			log.Printf("Error fetching feed: %v", err)
			return
		}
	}
	for _, item := range rssFeed.Channel.Item {
		log.Printf("Found post:%v,On feed %s\n", item.Title, feed.Name)
		//save to the database
		description := sql.NullString{String: item.Description}
		parsedTime, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("Fail to parse the time %v, err:%v", item.PubDate, err)
			continue
		}
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: parsedTime,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if !strings.Contains(err.Error(), "duplicate key") {
				log.Printf("Fail to create post: %v", err)
			}
			continue
		}
	}
	log.Printf("feed %s collected, %v post found", feed.Name, len(rssFeed.Channel.Item))
}
```

- 在`main.go`中新增下面内容：

```go
//start scraping the post
	go startScraping(queries, 10, 10*time.Second)
```

接下来编写获取`posts`的api：

- 新建`handler_post.go`，内容如下：

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/blkcor/Rss-aggregator/internal/database"
	"net/http"
	"strconv"
)

func (a *apiConfig) handlerGetUserPosts(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameter struct {
		Limit string `json:"limit"`
	}
	param := parameter{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&param)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error parsing Json Body:%v", err))
		return
	}
	convertedLimit, err := strconv.Atoi(param.Limit)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error parsing Json parameter:%v,please provide the correct parameter", err))
		return
	}
	posts, err := a.DB.GetPostForUser(r.Context(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(convertedLimit),
	})
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error getting posts for user:%v", err))
		return
	}
	postsResponse := make([]Post, 0)
	for _, post := range posts {
		postsResponse = append(postsResponse, dbPostToPost(post))
	}
	responseWithJson(w, 200, postsResponse)
}
```

- 将`handler`注册到`router`:

```go
//posts
v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetUserPosts))
```


# Twitter Scraper

[![Go Reference](https://pkg.go.dev/badge/github.com/blacheinc/twitter-scraper.svg)](https://pkg.go.dev/github.com/blacheinc/twitter-scraper)

This is an extended fork of [n0madic/twitter-scraper](https://github.com/n0madic/twitter-scraper)

## Installation

```shell
go get -u github.com/blacheinc/twitter-scraper
```

## Usage


### Get user followers

```golang
package main

import (
    "log"
    twitterscraper "github.com/blacheinc/twitter-scraper"
)

func main() {
 scraper := twitterscraper.New()

	if err = scraper.Login("username", "password"); err != nil {
		return false, err
	}
	
    // get the logged in user cookie
	cookie := scraper.GetCookies()

    // set cookie for subsequent
	scraper.SetCookies(cookie)

    // Get the user profile to extract the followers count 
    profile, err := scraper.GetProfile(twitterUsername)
	if err != nil {
		log.Fatal(err)
	}

	followers := scraper.GetFollowers(context.Background(), twitterUserID, profile.FollowersCount)
	
	for follower := range followers {
        // you will get the userIds in this format "user-947425510262562817"
        // when checking if a userid is among the return Ids use this
        // `formattedUserId := "user-" + user ` then compare.
	    fmt.Println(follower.UserID)
	}
}
```


### Get favorite tweets

```golang
package main

import (
    "log"
    twitterscraper "github.com/blacheinc/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New()

	if err = scraper.Login("username", "password"); err != nil {
		return false, err
	}

	// get the logged in user cookie
	cookie := scraper.GetCookies()

	// set cookie for subsequent
	scraper.SetCookies(cookie)

	tweets := scraper.FavoriteTweets(context.Background(), twitterUsername, 10)

    for tweet := range tweets {
        log.Println(tweet.Text)
    }
}
```

### Get user tweets

```golang
package main

import (
    "context"
    "fmt"
    twitterscraper "github.com/blacheinc/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New()

    for tweet := range scraper.GetTweets(context.Background(), "Twitter", 50) {
        if tweet.Error != nil {
            panic(tweet.Error)
        }
        fmt.Println(tweet.Text)
    }
}
```

It appears you can ask for up to 50 tweets (limit ~3200 tweets).

### Get single tweet

```golang
package main

import (
    "fmt"

    twitterscraper "github.com/blacheinc/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New()
    tweet, err := scraper.GetTweet("1328684389388185600")
    if err != nil {
        panic(err)
    }
    fmt.Println(tweet.Text)
}
```

### Search tweets by query standard operators

Now the search only works for authenticated users!

Tweets containing “twitter” and “scraper” and “data“, filtering out retweets:

```golang
package main

import (
    "context"
    "fmt"
    twitterscraper "github.com/blacheinc/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New()
    err := scraper.LoginOpenAccount()
    if err != nil {
        panic(err)
    }
    for tweet := range scraper.SearchTweets(context.Background(),
        "twitter scraper data -filter:retweets", 50) {
        if tweet.Error != nil {
            panic(tweet.Error)
        }
        fmt.Println(tweet.Text)
    }
}
```

The search ends if we have 50 tweets.

See [Rules and filtering](https://developer.twitter.com/en/docs/tweets/rules-and-filtering/overview/standard-operators) for build standard queries.

#### Set search mode

```golang
scraper.SetSearchMode(twitterscraper.SearchLatest)
```

Options:

- `twitterscraper.SearchTop` - default mode
- `twitterscraper.SearchLatest` - live mode
- `twitterscraper.SearchPhotos` - image mode
- `twitterscraper.SearchVideos` - video mode
- `twitterscraper.SearchUsers` - user mode

### Get profile

```golang
package main

import (
    "fmt"
    twitterscraper "github.com/blacheinc/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New()
    profile, err := scraper.GetProfile("Twitter")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", profile)
}
```

### Search profiles by query

```golang
package main

import (
    "context"
    "fmt"
    twitterscraper "github.com/blacheinc/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New().SetSearchMode(twitterscraper.SearchUsers)
    err := scraper.Login(username, password)
    if err !== nil {
        panic(err)
    }
    for profile := range scraper.SearchProfiles(context.Background(), "Twitter", 50) {
        if profile.Error != nil {
            panic(profile.Error)
        }
        fmt.Println(profile.Name)
    }
}
```

### Get trends

```golang
package main

import (
    "fmt"
    twitterscraper "github.com/blacheinc/twitter-scraper"
)

func main() {
    scraper := twitterscraper.New()
    trends, err := scraper.GetTrends()
    if err != nil {
        panic(err)
    }
    fmt.Println(trends)
}
```

### Use authentication

Some specified user tweets are protected that you must login and follow.
It is also required to search.

#### Login

```golang
err := scraper.Login("username", "password")
```

Use username to login, not email!
But if you have email confirmation, use email address in addition:

```golang
err := scraper.Login("username", "password", "email")
```

If you have two-factor authentication, use code:

```golang
err := scraper.Login("username", "password", "code")
```

Status of login can be checked with:

```golang
scraper.IsLoggedIn()
```

Logout (clear session):

```golang
scraper.Logout()
```

If you want save session between restarts, you can save cookies with `scraper.GetCookies()` and restore with `scraper.SetCookies()`.

For example, save cookies:

```golang
cookies := scraper.GetCookies()
// serialize to JSON
js, _ := json.Marshal(cookies)
// save to file
f, _ = os.Create("cookies.json")
f.Write(js)
```

and load cookies:

```golang
f, _ := os.Open("cookies.json")
// deserialize from JSON
var cookies []*http.Cookie
json.NewDecoder(f).Decode(&cookies)
// load cookies
scraper.SetCookies(cookies)
// check login status
scraper.IsLoggedIn()
```

#### Open account

If you don't want to use your account, you can login as a Twitter app:

```golang
err := scraper.LoginOpenAccount()
```

### Use Proxy

Support HTTP(s) and SOCKS5 proxy

#### with HTTP

```golang
err := scraper.SetProxy("http://localhost:3128")
if err != nil {
    panic(err)
}
```

#### with SOCKS5

```golang
err := scraper.SetProxy("socks5://localhost:1080")
if err != nil {
    panic(err)
}
```

### Delay requests

Add delay between API requests (in seconds)

```golang
scraper.WithDelay(5)
```

### Load timeline with tweet replies

```golang
scraper.WithReplies(true)
```

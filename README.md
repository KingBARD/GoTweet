# Go Twitter Library

A quick Twitter API library that I made while attempting to learn Go.

### Version
1.0.0

This library uses the following open source libraries to function:

* [Jason] - Jason
* [OAuth] - Go OAuth Library

### Examples:

Tweeting
```sh
    T := TwtterAPI.Account{"KEY", "KEY")
    
    //Must be called
    T.Auth()
    
    T.Tweet("This is my text","","", false, false)
    //This could also be done with
    resp, err := T.Tweet("This is my text","","", false, false)
    
    if err != nil {
        log.Fatal(err)
    }
    
    
    //This would just be the json response of the tweet been created successfully
    fmt.Println(resp)
    
```

Retweeting:
```sh

    T := TwtterAPI.Account{"KEY", "KEY")
    
    //Must be called
    T.Auth()
    
    T.Retweet("TweetID")
    
    //or
    
    resp, err := T.Retweet("TweetID")
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resp)
    

```

License
----

MIT


[Jason]:https://github.com/antonholmquist/jason
[OAuth]:https://github.com/garyburd/go-oauth
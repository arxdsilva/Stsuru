# Stsuru

This is a **simple** example of a link shortener that displays the MongoDB **fields** written in Go. Made by **[@arxdsilva](https://twitter.com/arxdsilva)**

### 1. Main external packages used on this version:
1. [Iris Web Framework](http://iris-go.com)
2. [MongoDB](https://gopkg.in/mgo.v2)

### 2. Functionalities:
1. Input your link inside MongoDB;
2. Creates & Displays the Hash of each URL;
3. Redirect access by using the hashed URL;
4. Removes the element (one by one) from Mongo by clicking on the element's "X".

### 3. Observations:
* It isn't necessarily shortening, because the Hash used is not the best for this usage.

### 4. [How HTML it is right now](https://github.com/ArxdSilva/Stsuru/blob/master/templates/mypage.html):
![now](https://ia601501.us.archive.org/13/items/ScreenShot20160822At4.44.59PM/Screen%20Shot%202016-08-22%20at%204.46.17%20PM.png)

### 5. [HTML goal](https://github.com/ArxdSilva/Stsuru/blob/master/index.html):
![goal](https://ia601501.us.archive.org/13/items/ScreenShot20160822At4.44.59PM/Screen%20Shot%202016-08-22%20at%204.44.59%20PM.png)

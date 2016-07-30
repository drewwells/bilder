# bilder - web app to host photo albums.

**bilder** is a standalone webserver on port 8173. It monitors a given directory for albums and serves them dynamically, including thumbnail generation. There's a live demo available [here](https://geller.io/bilder/b/kitties).

## Configuration

You can configure **bilder** via an optional JSON file. You can pass its location to **bilder** on startup:

```
$ bilder -config /path/to/your/bilder.json
```

It currently supports the following options:

 + `url-path-prefix` *default:* `""`: This is a prefix that can be added to the assets' paths that are loaded from the browser. This allows running **bilder** behind a proxy like nginx that can terminate the HTTPS connection. Consider the path of the demo linked above: [https://geller.io/bilder/b/kitties](https://geller.io/bilder/b/kitties). In this case nginx proxy passes to the **bilder** process under the `/bilder` location:
```
location /bilder/ {
    proxy_pass http://localhost:8173/;
}
```
 + `bilder-dir` *default:* `"bilder"`: This is the path of the folder that **bilder** scans for album directories. For example, this directory would contain a single album `kitties` (Please note that `index.html` and `*_thumb.jpg` are generated automatically by **bilder**. For details about the contained `bilder.json` consider the next section on [Albums](https://github.com/fgeller/bilder#albums)):
```
$ find bilder
bilder
bilder/kitties
bilder/kitties/400.jpeg
bilder/kitties/400_thumb.jpeg
bilder/kitties/bilder.json
bilder/kitties/index.html
```
 + `access-log` *default:* `""`: When set to a file name, **bilder** logs requests against the `/b` path in combined log format to the set file.

This is the JSON file that is used for the [demo](https://geller.io/bilder/b/kitties):
```
{ "bilder-dir": "/home/fgeller/var/bilder", "url-path-prefix": "/bilder" }
```

### Albums

Each sub-directory of the `bilder-dir` directory is considered an album if it contains JPG images. You can add more information about the album by adding a `bilder.json` to the directory. It currently supports the following options:

 + `user` *default:* `""`, `pass` *default:* `""`: If both are non-empty strings, **bilder** will use them as credentials to enable basic authentication for this album.
 + `title` *default:* `""`: Title that should be set for the album, defaults to the directory name.
 + `captions` *default:* `null`: Map object from filename to caption string (consider the demo example below).
 
This is the `bilder.json` file in the `kitties` directory of the [demo](https://geller.io/bilder/b/kitties):
```
{
  "title": "Kitties",
  "captions": {
    "cat-eyes.jpg": "looking",
    "mini-monster.jpg": "rooooar!",
    "yawning.jpg": "Boring!"
  }
}
```

## Credits

All images in the [demo](https://geller.io/bilder/b/kitties) are free images from [pixabay](https://pixabay.com/).

**bilder** uses the following libraries:

 + @dimsemenov's [PhotoSwipe](https://github.com/dimsemenov/PhotoSwipe) for rendering the album.
 + @nfnt's [resize](https://github.com/nfnt/resize) to generate thumbnails.
 + @oliamb's [cutter](https://github.com/oliamb/cutter) to crop thumbnails to a centered square.
 + @satori's [go.uuid](https://github.com/satori/go.uuid) to generate a random session ID.
 + @gorilla's [handlers](https://github.com/gorilla/handlers) for logging requests.

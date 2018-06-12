# Big Brother

_bigbro_ is a website interaction logging service. 

## Client

To use _bigbro_ on a website, include the following Javascript snippet (this repository hosts `bigbro.js` in the [js folder](js)):

```html
<script type="text/javascript" src="bigbro.js"></script>
<script type="text/javascript">
    BigBro.init("username", "localhost:1984");
</script>
```

This will allow the page to capture most events that occur by interacting with the page.

Custom logging events can also be added, see the following example:

```html
<script type="text/javascript" src="bigbro.js"></script>
<script type="text/javascript">
    let bb = BigBro.init("username", "localhost:1984");
    window.addEventListener("click", function (e) {
        bb.log(e, "custom_event");
    })
</script>
```

The arguments to `BigBro.init` are as follows:

 - `actor`: A unique identifier of the current user.
 - `server`: The address _bigbro_ is running on (please omit protocol; this will be determined automatically).
 - (optional) `events`: A list of events that will be listened on globally (at the window level); e.g. "click", "mousemove".

## Server

_bigbro_ is written in Go. To install locally, please use:

```bash
go get -u install github.com/hscells/bigbro
```

Alternatively, download a [prebuilt binary](https://github.com/hscells/bigbro/releases).

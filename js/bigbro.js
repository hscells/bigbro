let BigBro = {
    // This should not be modified outside of the init method.
    data: {
        user: "",
        server: "",
        events: ["click", "dblclick", "mousedown", "mouseup",
            "mouseenter", "mouseout", "wheel", "loadstart", "loadend", "load",
            "unload", "reset", "submit", "scroll", "resize",
            "cut", "copy", "paste", "select", "keydown", "keyup",
            "ontouchstart", "ontouchmove", "ontouchend", "ontouchcancel"
        ],
    },
    eventQueue: [],
    captureQueue: [],
    // State associated with recording user.
    startedRecording: false,
    captureStream: null,
    captureInterval: 1000,

    // init must be called with the user and the server, and optionally a list of
    // events to listen to globally.
    init: function (user, server, events) {
        this.data.user = user;
        this.data.server = server;
        this.data.events = events || this.data.events;

        let protocol = 'ws://';
        if (window.location.protocol === 'https:') {
            protocol = 'wss://';
        }

        this.wsEvent = new WebSocket(protocol + this.data.server + "/event");
        this.wsRecord = new WebSocket(protocol + this.data.server + "/capture");
        let self = this;
        this.wsEvent.onopen = function () {
            for (let i = 0; i < self.data.events.length; i++) {
                window.addEventListener(self.data.events[i], function (e) {
                    self.log(e, self.data.events[i]);
                })
            }
        };
        this.wsRecord.onopen = function () {
            window.setInterval(self.capture, self.captureInterval, self)
        };
        return this
    },
    // log logs an event with a specified method name (normally the actual event name).
    log: function (e, method, comment) {
        let event = {
            target: e.target.tagName,
            name: e.target.name,
            id: e.target.id,
            method: method,
            location: window.location.href,
            time: new Date().toISOString(),
            x: e.x,
            y: e.y,
            screenWidth: window.innerWidth,
            screenHeight: window.innerHeight,
            actor: this.data.user
        };
        if (method === "keydown" || method === "keyup") {
            // Which key was actually pressed?
            event.comment = e.code;
        }
        if (method === "paste" || method === "cut" || method === "copy") {
            // Seems like we can only get data for paste events.
            event.comment = e.clipboardData.getData("text/plain")
        }
        if (method === "wheel") {
            // Strength of the wheel rotation.
            event.comment = e.deltaY.toString();
        }
        if (comment != null) {
            event.comment = comment;
        }

        if (this.wsEvent.readyState !== 1) {
            this.eventQueue.push(event);
            return false;
        }

        while (this.eventQueue.length > 0) {
            this.wsEvent.send(JSON.stringify(this.eventQueue.pop()))
        }

        this.wsEvent.send(JSON.stringify(event));
    },
    startCaptureWhenClicked: async function (elementId, interval) {
        let self = this;
        document.getElementById(elementId).addEventListener("click", function () {
            try {
                if (!self.startedRecording) {
                    self.startedRecording = true;
                    self.captureStream = navigator.mediaDevices.getDisplayMedia({video: {cursor: "always"}, displaySurface: "browser", browserWindow: false});
                }
            } catch (err) {
                console.error("Error: " + err);
            }
        });
    },
    capture: function (self) {
        if (self.captureStream === null || self.captureStream === undefined) {
            return
        }
        if (!self.startedRecording) {
            return
        }
        const canvas = document.createElement("canvas");

        const context = canvas.getContext("2d");
        const video = document.createElement("video");

        try {

            self.captureStream.then((stream) => {
                video.srcObject = stream;
                video.play().then(() => {
                    canvas.setAttribute("height", window.innerHeight);
                    canvas.setAttribute("width", window.innerWidth);
                    context.drawImage(video, 0, 0, window.innerWidth, window.innerHeight);
                    let blob = canvas.toDataURL("image/png");

                    console.log(blob);

                    let capture = {
                        actor: self.data.user,
                        time: new Date().toISOString(),
                        data: blob,
                    };

                    if (self.wsRecord.readyState !== 1) {
                        self.captureQueue.push(capture);
                        return false;
                    }

                    while (self.captureQueue.length > 0) {
                        self.captureQueue.send(JSON.stringify(self.captureQueue.pop()))
                    }

                    self.wsRecord.send(JSON.stringify(capture))
                });


            })

        } catch (err) {
            console.error("Error: " + err);
        }
    }
};
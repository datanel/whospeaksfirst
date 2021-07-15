var app = new Vue({
    el: '#app',
    data() {
        return {
            ws: null,
            speakers: {},
        }
    },
    computed: {
        present: {
            get() {
                return this.speakers["present"]
            },
            set(newPresent) {
                this.speakers["present"] = newPresent
            }
        },
        absent: {
            get() {
                return this.speakers["absent"]
            },
            set(newAbsent) {
                this.speakers["absent"] = newAbsent
            }
        }
    },
    created: function () {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + trimTrailingSlash(window.location.pathname) + '/ws');
        this.ws.addEventListener('message', function (e) {
            var msg = JSON.parse(e.data);
            self.$data.speakers = msg;
        });
    },
    methods: {
        send() {
            this.present = shuffle(this.present)
            this.ws.send(
                JSON.stringify(this.speakers)
            );
        },
        reset() {
            const present = [...this.present, ...this.absent]
            present.sort()
            this.present = present
            this.absent = []
            this.ws.send(
                JSON.stringify(this.speakers)
            );
        },
        startDrag: function(event, item) {
            event.dataTransfer.dropEffect = "move";
            event.dataTransfer.effectAllowed = "move";
            event.dataTransfer.setData("speaker", item)
        },
        onDrop: function(event, presence) {
            const speaker = event.dataTransfer.getData("speaker");
            this.speakers[presence].push(speaker)
            this.present = this.present.filter((s, index) => s !== speaker);
            this.ws.send(
                JSON.stringify(this.speakers)
            );
        }
    }
});

function shuffle(arr) {
    return arr
    .map(a => [Math.random(), a])
    .sort((a, b) => a[0] - b[0])
    .map(a => a[1]);
}

function trimTrailingSlash(val) {
    return val.replace(/\/$/, "");
}
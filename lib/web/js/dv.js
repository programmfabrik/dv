
window.addEventListener("load", function(evt) {

    let ws = new WebSocket("ws://"+document.location.host+"/stream");
    ws.onopen = function(evt) {
        // print("OPEN");
    }
    ws.onclose = function(evt) {
        // print("CLOSE");
        ws = null;
    }
    var currentHeader = {}
    var currentType = ""
    ws.onmessage = function(evt) {
        // print("RESPONSE: " + evt.data);
        console.debug("message received", evt)
        if (evt.data instanceof Blob) {
            console.debug("received body: "+evt.data.size+" bytes")
            console.info(currentType+": "+evt.data.size+" bytes")
            if (!evt.data.size) {
                addText(currentType+": 0 bytes")
                return
            }
            switch (currentType) {
            case "text/plain":
            case "text/html":
                evt.data.text().then(function(text) {
                    addText(text)
                })
                break
            case "text/csv":
                evt.data.text().then(function(text) {
                    var results = Papa.parse(text.trim())
                    console.debug("results", results)

                    var tb = document.createElement("table")
                    tb.classList.add("csv")
                    var tbody = document.createElement("tbody")

                    for (row of results.data) {
                        var tr = document.createElement("tr")
                        for (cell of row) {
                            var td = document.createElement("td")
                            td.innerText = cell
                            tr.appendChild(td)
                        }
                        tbody.appendChild(tr)
                        console.debug("row", row)
                    }
                    tb.appendChild(tbody)
                    document.body.appendChild(tb)
                })
                break
            case "application/json":
                evt.data.text().then(function(jsonText) {
                    var treeEl = document.createElement("code")
                    treeEl.classList.add("json")
                    const tree = JsonView.renderJSON(jsonText, treeEl)
                    JsonView.expandChildrenDepth(tree, 1)
                    document.body.appendChild(treeEl)
                })
                break
            case "image/jpeg":
            case "image/png":
                var imgURL = URL.createObjectURL(evt.data)
                var img = document.createElement("img")
                img.classList.add("image")
                img.src = imgURL
                document.body.appendChild(img)
                console.debug("img", img)
                // URL.revokeObjectURL(imgURL)
                break
            default:
                let warning = currentType+": unsupported content-type"
                addText(warning).classList.add("error")
                console.warn(warning)
                break
            }
        } else if (typeof(evt.data) === "string") {
            try {
                currentHeader = JSON.parse(evt.data)
                currentType = currentHeader["Content-Type"][0].split(";")[0]
            } catch (ex) {
                currentHeader = null
                console.error(ex)
            }
            console.debug("received header: ["+currentType+"] "+evt.data.length+" bytes")
        } else {
            print("received: ["+typeof(evt.data)+"] unsupported type")
        }
    }
    ws.onerror = function(evt) {
        print("ERROR: " + evt.data);
    }
})

function addText(text) {
    var textDiv = document.createElement("div")
    textDiv.classList.add("text")
    textDiv.innerText = text
    document.body.appendChild(textDiv)
    return textDiv
}
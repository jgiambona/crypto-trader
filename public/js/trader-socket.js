
define(() => {
  if (document.readyState == "complete") {
    let loc = window.location
    let uri = 'ws:'
    if (loc.protocol === 'https:') {
      uri = 'wss:'
    }
    uri += '//' + loc.host
    uri += loc.pathname + 'ws'

    let ws = new WebSocket(uri)
    ws.onopen = () => {
      console.log('Connected')
    }

    ws.onmessage = (e) => {
      let out = document.getElementById('output')
      out.innerHTML += e.data + '<br>'
    }

    setInterval(() => {
      ws.send('Hello, Server!')
    }, 1000)
  }
})

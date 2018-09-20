// WS abstracts and easy to sue api over a raw WebSocket connection.
class WS {
    constructor(addr) {
        this.addr = addr
        this.filters = new Map()
        this._init()
    }
    send(name, data) {
        this.conn.send(JSON.stringify({name, data}))
    }
    on(name, cb) {
        this.filters.set(name, cb)
    }
    catch(cb) {
        this.conn.onerror = cb
        this._onErr = cb
    }
    reset() {
        this._init()
    }
    _init() {
        this.conn = new WebSocket(this.addr)
        this.conn.onmessage = (msg) => {
            let data = JSON.parse(msg.data)
            let handler = this.filters.get(data.name)
            if (handler != null) {
                handler(data.data)
            }
        }
        this.conn.onerror = this._onErr
    }
}
export default WS
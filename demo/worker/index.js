importScripts("wasm_exec.js");

(async () => {
  const go = new Go()
  const { instance } = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
  go.run(instance)
  global.postMessage({
    type: "ready",
  })
  global.addEventListener("message", messageListener)
})();

const messageListener = function ({ data = {} }) {
  const { type, message } = data;
  return ({
    "md5": md5ResponseHandler,
  }[type] || function () { })(message)
}

const md5ResponseHandler = function (data) {
  global.postMessage({
    type: "md5",
    message: getMd5(data)
  })
}




importScripts("wasm_exec.js");

(async () => {
  const go = new Go()
  const { instance } = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
  go.run(instance)
  global.addEventListener("message", function ({ data = {} }) {
    const { type, message } = data;
    global.postMessage({
      type: type,
      message: (global[type] || function () {
        return
      })(...message)
    })
  });
})();
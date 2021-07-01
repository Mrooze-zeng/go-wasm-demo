importScripts("wasm_exec.js");

(async () => {
  const go = new Go()
  const { instance } = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
  go.run(instance)
  global.addEventListener("message",async function ({ data = {} }) {
    const { type, message } = data;
   await global.postMessage({
      type: type,
      message: (global[type] || function () {
        return
      })(...message)
    })
  });
})();
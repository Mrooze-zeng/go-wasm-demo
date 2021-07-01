importScripts("wasm_exec.js");

(async () => {
  const go = new Go()
  const app = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
  const { instance } =app;
  console.log(app,go)
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
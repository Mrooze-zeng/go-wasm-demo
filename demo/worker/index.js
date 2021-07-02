importScripts("wasm_exec.js");

(async () => {
  const go = new Go()
  let app
  if (WebAssembly.instantiateStreaming) {
    app = await WebAssembly.instantiateStreaming(fetch("app.wasm"), go.importObject)
  } else {
    const res = await fetch("app.wasm");
    app = await WebAssembly.instantiate(await res.arrayBuffer(), go.importObject)
  }
  const { instance } = app;
  go.run(instance)
  global.addEventListener("message", async function ({ data = {} }) {
    const { type = "", message = [] } = data;
    const sendMessage = type.split('.').reduce(function (global, key) {
      return global[key]||function(){};
    }, global);
    await global.postMessage({
      type: type,
      message: await sendMessage(...message)
    })
  });
})();
const initWorker = function ({ file = "worker/index.js", handler = function () { } }) {
  return new Promise(function (resolve, reject) {
    const myWorker = new Worker(file)
    myWorker.addEventListener("message", function ({ data = {} }) {
      switch (data.type) {
        case "ready":
          resolve(myWorker)
          break;
        default:
          handler.call(myWorker, data)
      }
    })
  })
}

const getFileBuffer = function (callback = function () { }) {
  const fileInput = document.getElementById("j-upload")
  fileInput.addEventListener("change", function (e) {
    const fileReader = new FileReader()
    fileReader.onload = function () {
      callback(new Uint8Array(this.result))
    }
    fileReader.readAsArrayBuffer(this.files[0])
  })
}

const getTextBuffer = function (callback = function () { }) {
  const $input = document.getElementById("j-text");
  const $btn = document.getElementById("j-md5-text")

  $btn.addEventListener("click", function () {
    callback(new TextEncoder().encode($input.value))
  })
}


const md5ResponseHandler = function (md5Code) {
  document.getElementById("j-result").innerHTML = `md5: ${md5Code}`
  console.log("md5:" + md5Code)
}

const workerHandler = function ({ type, message }) {
  return ({
    "md5": md5ResponseHandler,
  }[type] || function () { }).call(this, message)
};

(async () => {
  const worker = await initWorker({
    handler: workerHandler
  })
  getFileBuffer(function (buffer) {
    worker.postMessage({ type: "md5", message: buffer })
  })

  getTextBuffer(function (buffer) {
    worker.postMessage({ type: "md5", message: buffer })
  })


})();
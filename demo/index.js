const getFileBuffer = function (callback = function () { }) {
  const fileInput = document.getElementById("j-upload")
  fileInput.addEventListener("change", function (e) {
    const fileReader = new FileReader()
    fileReader.onload = function () {
      const worker = callback(new Uint8Array(this.result));
      const listener = function ({ data = {} }) {
        const { type, message } = data;
        if (type === "getMd5") {
          document.getElementById("j-result").innerHTML = `md5: ${message}`
          console.log("md5:" + message)
          worker.removeEventListener("message", listener)
        }
      }
      worker.addEventListener("message", listener)
    }
    fileReader.readAsArrayBuffer(this.files[0])
  })
}

const getTextBuffer = function (callback = function () { }) {
  const $input = document.getElementById("j-text");
  const $btn = document.getElementById("j-md5-text")
  $btn.addEventListener("click", function () {
    const worker = callback(new TextEncoder().encode($input.value));
    const listener = function ({ data = {} }) {
      const { type, message } = data;
      if (type === "getMd5") {
        document.getElementById("j-result").innerHTML = `md5: ${message}`
        console.log("md5:" + message)
        worker.removeEventListener("message", listener)
      }
    }
    worker.addEventListener("message", listener)
  })
}

const rotateImage = function (callback = function () { }) {
  const $img = document.getElementById("j-img")
  const $btn = document.getElementById("j-img-control")
  const $preview = document.createElement("img")

  $preview.width = 250

  const getBufferCache = function () {
    let buffer;
    return async function (url) {
      if (buffer) {
        return buffer;
      }
      const res = await fetch(url)
      buffer = await res.arrayBuffer()
      return buffer;
    }
  }


  const getBuffer = getBufferCache()

  $btn.addEventListener("click", async function () {
    const self = this;
    const buffer = await getBuffer($img.src)
    const directions = [1, 3, 6, 8, 9]
    const direction = directions[Math.floor(Math.random() * directions.length)]
    const worker = callback(new Uint8Array(buffer), direction);
    this.setAttribute("disabled", true)
    const listener = function ({ data = {} }) {
      const { type, message } = data;
      if (type === "imageRotate" && message) {
        URL.revokeObjectURL($preview.src)
        $preview.src = URL.createObjectURL(new Blob([message.buffer.buffer], { type: message.type }))
        $preview.onload = () => {
          self.removeAttribute("disabled")
          worker.removeEventListener("message", listener)
        }
      }
    }
    worker.addEventListener("message", listener)
    document.body.appendChild($preview)
  })
}

const worker = new Worker("worker/index.js");

getFileBuffer(function (buffer) {
  worker.postMessage({ type: "getMd5", message: [buffer] })
  return worker;
})

getTextBuffer(function (buffer) {
  worker.postMessage({ type: "getMd5", message: [buffer] })
  return worker;
})

rotateImage(function (buffer, direction) {
  worker.postMessage({ type: "imageRotate", message: [buffer, direction] })
  return worker
})
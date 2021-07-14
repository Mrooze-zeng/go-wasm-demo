const getFileBuffer = function (callback = function () {}) {
  const $fileInput = document.getElementById("j-upload");
  const $btn = document.getElementById("j-md5-file");
  const fileMd5Handler = function (file) {
    const start = window.performance.now();
    const fileReader = new FileReader();
    fileReader.onload = function () {
      const worker = callback(new Uint8Array(this.result));
      const listener = function ({ data = {} }) {
        const { type, message } = data;
        if (type === "getMd5" && message) {
          document.getElementById(
            "j-result",
          ).innerHTML = `md5: ${message.data}`;
          console.log(message);
          console.log("耗时:", window.performance.now() - start, "毫秒");
        }
      };
      worker.addEventListener("message", listener, { once: true });
    };
    fileReader.readAsArrayBuffer(file);
  };
  $fileInput.addEventListener("change", function (e) {
    fileMd5Handler(this.files[0]);
  });
  $btn.addEventListener("click", function () {
    $fileInput.files[0] && fileMd5Handler($fileInput.files[0]);
  });
};

const getTextBuffer = function (callback = function () {}) {
  const $input = document.getElementById("j-text");
  const $btn = document.getElementById("j-md5-text");
  $btn.addEventListener("click", function () {
    const start = window.performance.now();
    const worker = callback(new TextEncoder().encode($input.value));
    const listener = function ({ data = {} }) {
      const { type, message } = data;
      if (type === "getMd5" && message) {
        document.getElementById("j-result").innerHTML = `md5: ${message.data}`;
        console.log(message);
        console.log("耗时:", window.performance.now() - start, "毫秒");
      }
    };
    worker.addEventListener("message", listener, { once: true });
  });
};

const rotateImage = function (
  callback = function () {},
  release = function () {},
) {
  const $img = document.getElementById("j-img");
  const $btn = document.getElementById("j-img-control");
  const $release = document.getElementById("j-img-release");
  const $preview = document.createElement("img");

  $preview.width = 250;

  const getBufferCache = function () {
    let buffer;
    return async function (url) {
      if (buffer) {
        return buffer;
      }
      const res = await fetch(url);
      buffer = await res.arrayBuffer();
      return buffer;
    };
  };

  const getBuffer = getBufferCache();

  $btn.addEventListener("click", async function () {
    const self = this;
    const start = window.performance.now();
    const buffer = await getBuffer($img.src);
    const directions = [1, 3, 6, 8, 9];
    const direction = directions[Math.floor(Math.random() * directions.length)];
    const worker = callback(new Uint8Array(buffer), direction);
    this.setAttribute("disabled", true);
    const listener = function ({ data = {} }) {
      const { type, message } = data;
      if (type === "imageRotate.run" && message) {
        console.log("耗时:", window.performance.now() - start, "毫秒");
        URL.revokeObjectURL($preview.src);
        $preview.src = URL.createObjectURL(
          new Blob([message.data.buffer], { type: message.type }),
        );
        $preview.onload = () => {
          self.removeAttribute("disabled");
        };
      }
    };
    worker.addEventListener("message", listener, { once: true });
    document.body.appendChild($preview);
  });
  $release.addEventListener("click", release);
};

const createDownloadLink = function (url = "", filename = "") {
  const link = document.createElement("a");
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
};

const getExcel = function (callback = function () {}) {
  $btn = document.getElementById("j-excel");
  $btn.addEventListener("click", callback);
  const listener = function ({ data = {} }) {
    const { type, message } = data;
    if (type === "getExcel" && message) {
      const url = URL.createObjectURL(
        new Blob([message.data.buffer], {
          type: message.type,
        }),
      );
      createDownloadLink(url, "test.xlsx");
      URL.revokeObjectURL(url);
    }
  };
  worker.addEventListener("message", listener, { once: true });
};

const getCSV = function (callback = function () {}) {
  let $btn = document.getElementById("j-csv");
  $btn.addEventListener("click", callback);
  const listener = function ({ data = {} }) {
    const { type, message } = data;
    if (type === "setCSV" && message) {
      console.log(message);
      const url = URL.createObjectURL(
        new Blob([message.data.buffer], {
          type: message.type,
        }),
      );
      createDownloadLink(url, "test.csv");
      URL.revokeObjectURL(url);
    }
  };
  worker.addEventListener("message", listener, { once: true });
};

const worker = new Worker("worker/index.js");

getFileBuffer(function (buffer) {
  worker.postMessage({ type: "getMd5", message: [buffer] });
  return worker;
});

getTextBuffer(function (buffer) {
  worker.postMessage({ type: "getMd5", message: [buffer] });
  return worker;
});

rotateImage(
  function (buffer, direction) {
    worker.postMessage({
      type: "imageRotate.run",
      message: [buffer, direction],
    });
    return worker;
  },
  function () {
    worker.postMessage({ type: "imageRotate.release" });
  },
);

getExcel(function () {
  worker.postMessage({ type: "getExcel" });
  return worker;
});

getCSV(function () {
  worker.postMessage({ type: "setCSV" });
  return worker;
});

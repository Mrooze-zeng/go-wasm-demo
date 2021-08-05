const getFileBuffer = function (callback = function () {}) {
  const $fileInput = document.getElementById("j-upload");
  const $btn = document.getElementById("j-md5-file");
  const fileMd5Handler = async function (file) {
    const start = window.performance.now();
    const result = await fileReader(file);
    const worker = callback(new Uint8Array(result));
    const listener = function ({ data = {} }) {
      const { type, message } = data;
      if (type === "getMd5" && message) {
        document.getElementById("j-result").innerHTML = `md5: ${message.data}`;
        console.log(message);
        console.log("耗时:", window.performance.now() - start, "毫秒");
      }
    };
    worker.addEventListener("message", listener, { once: true });
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

const fileReader = function (file) {
  return new Promise((resolve, reject) => {
    if (!file) {
      reject();
    }
    const fileReader = new FileReader();
    fileReader.onload = function () {
      resolve(this.result);
    };
    fileReader.readAsArrayBuffer(file);
  });
};

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

const rotateImage = function (
  callback = function () {},
  release = function () {},
) {
  const $img = document.getElementById("j-img");
  const $btn = document.getElementById("j-img-control");
  const $release = document.getElementById("j-img-release");
  const $preview = document.createElement("img");

  const getBuffer = getBufferCache();

  $preview.width = 250;

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
  $btn.addEventListener("click", function () {
    const worker = callback();
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
  });
};

const getCSV = function (callback = function () {}) {
  let $btn = document.getElementById("j-csv");
  $btn.addEventListener("click", function () {
    const worker = callback();
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
  });
};

const getImageThumbnail = function (callback = function () {}) {
  const $btn = document.getElementById("j-img-thumbnail");
  const $img = document.getElementById("j-img");
  const getBuffer = getBufferCache();

  $btn.addEventListener("click", async function () {
    const start = window.performance.now();
    const buffer = await getBuffer($img.src);
    const worker = callback(new Uint8Array(buffer));
    const listener = function ({ data = {} }) {
      const { type, message } = data;
      if (type === "getImageThumbnail" && message) {
        const url = URL.createObjectURL(
          new Blob([message.data.buffer], {
            type: message.type,
          }),
        );
        const image = document.createElement("img");
        image.src = url;
        document.body.appendChild(image);
        image.onload = function () {
          URL.revokeObjectURL(url);
          console.log("耗时:", window.performance.now() - start, "毫秒");
        };
      }
    };
    worker.addEventListener("message", listener, { once: true });
  });
};

const parseVideo = function (callback = function () {}) {
  const $btn = document.getElementById("j-video-btn");
  const $video = document.getElementById("j-video");
  $btn.addEventListener("click", async function () {
    if ($video.files[0]) {
      const result = await fileReader($video.files[0]);
      const worker = callback(
        URL.createObjectURL(
          new Blob([result], {
            type: $video.files[0].type,
          }),
        ),
      );
      const listener = function ({ data = {} }) {
        const { type, message } = data;
        if (type === "parseVideo" && message) {
          console.log(message);
        }
      };
      worker.addEventListener("message", listener, { once: true });
    }
  });
};

const gzipFile = function (callback = function () {}) {
  const $btn = document.getElementById("j-gzip-file");
  const $fileInput = document.getElementById("j-upload");

  $btn.addEventListener("click", async function () {
    const file = $fileInput.files[0];
    if (file) {
      const result = await fileReader(file);
      const worker = callback(new Uint8Array(result), file.name);
      const listener = function ({ data = {} }) {
        const { type, message } = data;
        if (type === "compress.gzip" && message) {
          const url = URL.createObjectURL(
            new Blob([message.data.buffer], {
              type: message.type,
            }),
          );
          createDownloadLink(url, `${file.name}.gz`);
          URL.revokeObjectURL(url);
        }
      };
      worker.addEventListener("message", listener, { once: true });
    }
  });
};

const ungzipFile = function (callback = function () {}) {
  const $btn = document.getElementById("j-ungzip-file");
  const $fileInput = document.getElementById("j-upload");

  $btn.addEventListener("click", async function () {
    const file = $fileInput.files[0];
    console.log(file);
    if (file) {
      const result = await fileReader(file);
      const worker = callback(new Uint8Array(result), file.name);
      const listener = function ({ data = {} }) {
        const { type, message } = data;
        if (type === "compress.ungzip" && message) {
          console.log(message);
          const url = URL.createObjectURL(new Blob([message.data.buffer]));
          createDownloadLink(url, message.name);
        }
      };
      worker.addEventListener("message", listener, { once: true });
    }
  });
};

const sliceUpload = function(callback=function(){}){
  const $fileInput = document.getElementById("j-file");
  const $btn = document.getElementById("j-file-btn");
  $btn.addEventListener("click", async function () {
    const file = $fileInput.files[0];
    console.log(file);
    if (file) {
      const result = await fileReader(file);
      const worker = callback(new Uint8Array(result),{
        buffer:new Uint8Array(result),
        mintype:file.type,
        name:file.name,
        chunkSize:1024*1024
      },"http://127.0.0.1:8080/upload");
      const listener = function ({ data = {} }) {
        const { type, message } = data;
        if (type === "sliceUpload" && message) {
          console.log(message);
        }
      };
      worker.addEventListener("message", listener, { once: true });
    }
  })
}

const sliceDownload = function(callback=function(){}){
  const $btn = document.getElementById("j-down-btn");
  $btn.addEventListener("click", async function () {
    // const worker = callback(location.origin+'/cow.jpg',1024*1024)
    const worker = callback("http://127.0.0.1:5000/tmp/dump.zip",1024*1024)
    const listener = function ({ data = {} }) {
      const { type, message } = data;
      if (type === "sliceDownload" && message) {
        console.log(message);
      }
    };
    worker.addEventListener("message", listener, { once: true });
  })
}

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

getImageThumbnail(function (buffer) {
  worker.postMessage({ type: "getImageThumbnail", message: [buffer] });
  return worker;
});

parseVideo(function (url) {
  console.log(url);
  const video = document.createElement("video");
  video.src = url;
  video.controls = true;
  document.body.appendChild(video);
  worker.postMessage({ type: "parseVideo", message: [url] });
  return worker;
});

gzipFile(function (buffer, name) {
  worker.postMessage({ type: "compress.gzip", message: [buffer, name] });
  return worker;
});

ungzipFile(function (buffer, name,) {
  worker.postMessage({ type: "compress.ungzip", message: [buffer, name] });
  return worker;
});

sliceUpload(function(buffer,options,apiUrl){
  worker.postMessage({ type: "sliceUpload", message: [buffer, options,apiUrl] });
  return worker;
})

sliceDownload(function(url,size){
  worker.postMessage({ type: "sliceDownload", message: [url,size] });
  return worker;
})
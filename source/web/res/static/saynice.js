var themes = [
  "primary",
  "secondary",
  "green",
  "blue",
  "orange",
  "red",
  "black",
  "accent primary",
  "accent secondary",
  "accent green",
  "accent blue",
  "accent blue",
  "accent orange",
  "accent black",
]

function initTheme() {
  var theme = Cookies.get("theme")

  if (theme) {
    setTheme(theme)
  } else {
    randTheme()
  }
}

function randTheme() {
  var i = Math.random() * 10000 % themes.length
  var x = themes[Math.floor(i)]
  document.body.className = document.body.className + " " + x

  Cookies.set("theme", x)
}

function setTheme(theme) {
  document.body.className = document.body.className + " " + theme

  Cookies.set("theme", theme)
}

initTheme()

function httpGet(url, callback) {
  var request = new XMLHttpRequest();
  request.open("GET", url);
  request.responseType = "json";
  request.send();
  request.onload = function () {
    if (200 != request.status) {
      Snackbar.show({ text: "网络异常，请稍后再试。", });
    } else {
      callback(request.response)
    }
  }
}

function httpPost(url, data, callback) {
  var request = new XMLHttpRequest();
  request.open("POST", url);
  request.send(JSON.stringify(data));
  request.onload = function () {
    if (200 != request.status) {
      Snackbar.show({ text: "网络异常，请稍后再试。", });
    } else {
      callback(request.response)
    }
  }
}
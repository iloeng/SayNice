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

function alert(message, a, b, aCallback, bCallback) {
  var modal = document.createElement("div")
  modal.className = "modal is-active"

  var bg = document.createElement("div")
  bg.className = "modal-background"
  modal.appendChild(bg)

  var content = document.createElement("div")
  content.className = "modal-content"

  var dialog = document.createElement("div")
  dialog.className = "dialog"
  dialog.style = "text-align: center;"

  var title = document.createElement("h5")
  title.innerText = message

  var yesButton = document.createElement("button")
  yesButton.innerText = a
  yesButton.className = "text"
  yesButton.onclick = function () {
    document.body.removeChild(modal)
    if (undefined != aCallback) {
      aCallback()
    }
  }

  var noButton = document.createElement("button")
  noButton.innerText = b
  noButton.onclick = function () {
    document.body.removeChild(modal)
    if (undefined != bCallback) {
      bCallback()
    }
  }

  dialog.appendChild(title)
  dialog.appendChild(yesButton)
  dialog.appendChild(noButton)

  content.appendChild(dialog)

  modal.appendChild(content)
  document.body.appendChild(modal)
}

function report(url, title, articles, a, b) {
  var reportModal = document.createElement("div")
  reportModal.className = "modal is-active"

  var bg = document.createElement("div")
  bg.className = "modal-background"
  bg.onclick = function () {
    document.body.removeChild(reportModal)
  }

  var content = document.createElement("div")
  content.className = "modal-content"

  var dialog = document.createElement("div")
  dialog.className = "dialog"
  dialog.style = "text-align: center;"

  var t = document.createElement("h4")
  t.innerText = title

  var checked = []
  var checkboxs = document.createElement("div")
  checkboxs.style = "text-align: start"

  for (let index = 0; index < articles.length; index++) {
    const article = articles[index];
    checked[index] = false

    var p = document.createElement("p")
    var checkbox = document.createElement("input")
    checkbox.type = "checkbox"
    checkbox.name = "article" + index
    checkbox.id = "article" + index
    checkbox.onclick = function () {
      var cb = document.getElementById("article" + index)
      checked[index] = cb.checked
    }
    var label = document.createElement("label")
    label.style = "margin-left: 10px;"
    label.htmlFor = checkbox.id
    label.innerText = article
    p.appendChild(checkbox)
    p.appendChild(label)
    checkboxs.appendChild(p)
  }

  var yesButton = document.createElement("button")
  yesButton.innerText = a
  yesButton.className = "text"
  yesButton.onclick = function () {
    document.body.removeChild(reportModal)
  }

  var noButton = document.createElement("button")
  noButton.innerText = b
  noButton.onclick = function () {
    document.body.removeChild(reportModal)

    var remark = ""
    for (let index = 0; index < checked.length; index++) {
      if (checked[index]) {
        remark += articles[index] + ";"
      }
    }

    apiPost(url, { "remark": remark }, function (resp) {
      Snackbar.show({ text: "SUCCESS", });
    }, function (status, resp) {
      if (200 != status) {
        Snackbar.show({ text: "网络异常，请稍后再试。", });
      } else {
        Snackbar.show({ text: "提交失败：" + resp.erro, });
      }
    })
  }

  dialog.appendChild(t)
  dialog.appendChild(checkboxs)
  dialog.appendChild(yesButton)
  dialog.appendChild(noButton)

  content.appendChild(dialog)

  var close = document.createElement("close")
  close.className = "modal-close"
  close.onclick = function () {
    document.body.removeChild(reportModal)
  }

  reportModal.appendChild(bg)
  reportModal.appendChild(content)
  reportModal.appendChild(close)

  document.body.appendChild(reportModal)
}

function appendArticle(parent, article, yes, no, yesCallback, noCallback) {
  var articleDiv = document.createElement("div")
  articleDiv.className = "article"

  var articleText = document.createElement("p")
  articleText.innerText = article

  var yesButton = document.createElement("button")
  yesButton.innerText = yes
  yesButton.onclick = yesCallback

  var noButton = document.createElement("button")
  noButton.innerText = no
  noButton.onclick = noCallback

  articleDiv.appendChild(articleText)
  articleDiv.appendChild(yesButton)
  articleDiv.appendChild(noButton)

  parent.appendChild(articleDiv)
}

function httpGet(url, success, failure = undefined) {
  var request = new XMLHttpRequest();
  request.open("GET", url);
  request.responseType = "json";
  request.send();
  request.onload = function () {
    httpOnLoad(request, success, failure)
  }
}

function httpPost(url, data, success, failure = undefined) {
  var request = new XMLHttpRequest();
  request.open("POST", url);
  request.responseType = "json";
  request.send(JSON.stringify(data));
  request.onload = function () {
    httpOnLoad(request, success, failure)
  }
}

function httpOnLoad(request, success, failure) {
  if (200 != request.status) {
    if (undefined == failure) {
      Snackbar.show({ text: "网络异常，请稍后再试。", });
    } else {
      failure(request.status, request.response)
    }
  } else {
    success(request.response)
  }
}

function apiGet(url, success, failure = undefined) {
  httpGet(url, function (resp) {
    if (0 == resp.code) {
      success(resp)
    } else if (undefined == failure) {
      Snackbar.show({ text: "请求失败: " + resp.erro, });
    } else {
      failure(200, resp)
    }
  }, failure)
}

function apiPost(url, data, success, failure = undefined) {
  httpPost(url, data, function (resp) {
    if (0 == resp.code) {
      success(resp)
    } else if (undefined == failure) {
      Snackbar.show({ text: "提交失败: " + resp.erro, });
    } else {
      failure(200, resp)
    }
  }, failure)
}

// see: https://blog.csdn.net/qq_33401924/article/details/88898811
let Time = {
  //    获取当前时间戳
  getUnix: function () {
    let date = new Date()
    return date.getTime()
  },
  //    获取今天0点0分0秒的时间戳
  getTodayUnix: function () {
    let date = new Date()
    date.setHours(0)
    date.setMinutes(0)
    date.setSeconds(0)
    date.setMilliseconds(0)
    return date.getTime()
  },
  //    获取今年1月1日0点0分0秒的时间戳
  getYearUnix: function () {
    let date = new Date()
    date.setMonth(0)
    date.setDate(1)
    date.setHours(0)
    date.setMinutes(0)
    date.setSeconds(0)
    date.setMilliseconds(0)
    return date.getTime()
  },
  //    获取标准年月
  getLastDate: function (time) {
    let date = new Date(time)
    let month = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1
    let day = date.getDate() + 1 < 10 ? '0' + (date.getDate() + 1) : date.getDate() + 1
    return date.getFullYear() + '-' + month + '-' + day
  },
  //    转换时间
  getFormatTime: function (timestamp) {
    let now = this.getUnix()
    let today = this.getTodayUnix()
    let year = this.getYearUnix()
    let timer = (now - timestamp) / 1000
    let tip = ''

    if (timer <= 0) {
      tip = '刚刚'
    } else if (Math.floor(timer / 60) <= 0) {
      tip = '刚刚'
    } else if (Math.floor(timer < 3600)) {
      tip = Math.floor(timer / 60) + '分钟前'
    } else if (timer >= 3600 && timer < 86400) {
      tip = Math.floor(timer / 3600) + '小时前'
    } else if (timer / 86400 <= 31) {
      tip = Math.floor(timer / 86400) + '天前'
    } else {
      tip = this.getLastDate(timestamp)
    }
    return tip
  }
}

function autoLoad(offset, callback) {
  //文档内容实际高度（包括超出视窗的溢出部分）
  var scrollHeight = Math.max(document.documentElement.scrollHeight, document.body.scrollHeight);
  //滚动条滚动距离
  var scrollTop = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop;
  //窗口可视范围高度
  var clientHeight = window.innerHeight || Math.min(document.documentElement.clientHeight, document.body.clientHeight);

  if (scrollHeight <= clientHeight + scrollTop + offset) {
    callback()
  }
}

var scrollTimer
var scrollStatus = 0
const scrollTimeout = 800

function onWindowScroll(handler) {
  clearTimeout(scrollTimer)
  if (0 == scrollStatus) {
    scrollStatus = 1
    handler(scrollStatus)
  } else {
    scrollStatus = 2
    handler(scrollStatus)
  }
  scrollTimer = setTimeout(function () {
    scrollStatus = 3
    handler(scrollStatus)
    scrollStatus = 0
  }, scrollTimeout)
}

function onNextPosts(
  url,
  postsDivName,
  reportButtonText,
  onTextAllCallback,
  onEmojiCallback,
  onReportCallback,
  onAttitudeCallback,
  completedCallback,
  noMoreCallback
) {
  apiGet(url, function (resp) {
    var postsDiv = document.getElementById(postsDivName)
    for (let index = 0; index < resp.data.length; index++) {
      const post = resp.data[index];

      var postDiv = document.createElement("div")
      postDiv.className = "post"

      var postTextDiv = document.createElement("div")
      var postTextAllBtn = document.createElement("button")

      postTextDiv.id = "t" + post.id
      postTextDiv.className = "text"
      postTextDiv.innerText = post.text
      postTextAllBtn.innerText = "全文"
      postTextAllBtn.onclick = function () { onTextAllCallback(post.id) }
      postTextDiv.appendChild(postTextAllBtn)

      postDiv.appendChild(postTextDiv)

      var timeSpan = document.createElement("span")
      var createdAt = new Date(post.createdAt.replace(/-/gi, "/")).getTime()
      timeSpan.innerText = Time.getFormatTime(createdAt)

      var toolDiv = document.createElement("div")
      var emojiBtn = document.createElement("button")
      var reportBtn = document.createElement("button")

      emojiBtn.className = "text"
      emojiBtn.innerHTML = "😀"
      emojiBtn.onclick = function () { onEmojiCallback(post.id) }
      reportBtn.className = "none"
      reportBtn.innerHTML = reportButtonText
      reportBtn.onclick = function () { onReportCallback(post.id) }

      toolDiv.appendChild(emojiBtn)
      toolDiv.appendChild(reportBtn)

      var barDiv = document.createElement("div")
      barDiv.className = "bar"
      barDiv.appendChild(timeSpan)
      barDiv.appendChild(toolDiv)

      postDiv.appendChild(barDiv)

      var attitudesDiv = document.createElement("div")
      attitudesDiv.id = "a" + post.id
      attitudesDiv.className = "attitudes"

      if (undefined != post.attitudes && null != post.attitudes) {

        attitudes = Object.entries(post.attitudes)
        attitudes.sort((a, b) => b[1] - a[1]);

        for (let [emojiId, count] of attitudes) {
          emojiId = parseInt(emojiId)

          var attitudeDiv = document.createElement("button")
          attitudeDiv.className = "attitude"
          attitudeDiv.dataset.id = emojiId
          attitudeDiv.dataset.count = count
          attitudeDiv.onclick = function () { onAttitudeCallback(this, post.id) }

          let emojiShow = emojiMap.get(emojiId)

          if (undefined == emojiShow || null == emojiShow) {
            continue
          }
          if (1 == count) {
            attitudeDiv.innerHTML = emojiShow
          } else {
            attitudeDiv.innerHTML = emojiShow + "+" + count
          }

          attitudesDiv.appendChild(attitudeDiv)
        }
      }

      postDiv.appendChild(attitudesDiv)

      var emojisDiv = document.createElement("div")
      emojisDiv.id = "f" + post.id
      emojisDiv.className = "emojis"
      emojisDiv.hidden = true

      postDiv.appendChild(emojisDiv)

      postsDiv.appendChild(postDiv)
    }

    completedCallback(resp.data.length)
  }, function (status, resp) {
    if (200 != status) {
      Snackbar.show({ text: "网络异常，请稍后再试", });
    } else if (10020 == resp.code) { // 没有更多了
      noMoreCallback()
    } else {
      Snackbar.show({ text: "请求失败: " + resp.erro, });
    }
  })
}

function initCreatedAt() {
  var createAtDiv = document.getElementsByClassName("createdAt")
  for (let index = 0; index < createAtDiv.length; index++) {
    const element = createAtDiv[index];
    let createdAt = new Date(element.dataset.time.replace(/-/gi, "/")).getTime()
    element.innerHTML = Time.getFormatTime(createdAt)
    // element.className = ""
  }
}

function initTheme() {
  var theme = Cookies.get("theme")

  if (theme) {
    setTheme(theme)
  } else {
    randTheme()
  }
}

function initArticles(url, callback) {
  // var articles = Cookies.get("articles")
  // if (undefined != articles && null != articles && 0 != articles.length) {
  //   callback(JSON.parse(articles))
  // } else {
  //   apiGet(url, function (resp) {
  //     Cookies.set("articles", JSON.stringify(resp.data))
  //     callback(resp.data)
  //   })
  // }
  apiGet(url, function (resp) {
    Cookies.set("articles", JSON.stringify(resp.data))
    callback(resp.data)
  })
}

initTheme()
initCreatedAt()
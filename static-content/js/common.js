// WebSocket通信を開始する関数
let startWebSocket = () => {
  // uriの組み立て
  var loc = window.location;
  var uri = 'ws:';
  if (loc.protocol === 'https:') {
    uri = 'wss:';
  }
  uri += '//' + loc.host + '/ws';
  ws = new WebSocket(uri);

  // WebSocket通信開始時のコールバック関数
  ws.onopen = function () {
    console.log('Connected');
  }

  // メッセージ受信時のコールバック関数
  ws.onmessage = function (evt) {
    // 更新された場合は検索実行
    if (evt.data == "UPDATED") {
      var event = new CustomEvent('history_update');
      document.dispatchEvent(event);
    }
  }
}

// トーストを表示する関数
let createToast = (msg) => {
  var toastEl = document.getElementById("toast");
  var toast = new bootstrap.Toast(toastEl, {
    animation: true,
    autohide: true,
    delay: 3000
  });
  toastBody = document.getElementById("toast_body");
  toastBody.textContent = msg
  toast.show();
}

// 初期化関数
let initCommon = () => {
  startWebSocket();
  document.addEventListener('history_update', () => { createToast("検索結果が更新されました。") });
}

// イベントに対応する処理の追加
window.addEventListener('load', initCommon);
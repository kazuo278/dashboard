// 当月1日の日付を返す関数
let getFirstDayOfTheMonth = () => {
  var today = new Date();
  var year = today.getFullYear().toString().padStart(4, "0");
  var month = (today.getMonth() + 1).toString().padStart(2, "0");
  var day = "01"
  return year + "-" + month + "-" + day;
}

// 検索URLを組み立てる関数
let createRequestUri = () => {
  var uri = window.location.protocol + "//" + window.location.host + "/actions/count";
  var params = new URLSearchParams();

  var repository_name = document.getElementById("repository_name").value;
  if (repository_name) {
    params.append("repository_name", repository_name);
  }

  var started_at = document.getElementById("started_at").value;
  if (started_at) {
    params.append("started_at", started_at + "T00:00:00+09:00");
  }

  var finished_at = document.getElementById("finished_at").value;
  if (finished_at) {
    params.append("finished_at", finished_at + "T00:00:00+09:00");
  }
  uri += "?" + new URLSearchParams(params).toString();
  return uri;
}

// JSONデータからテーブルレコードを作成する関数
let displayRedcords = (records) => {
  var tbody = document.getElementById("table_body");
  // 既存データを削除
  while (tbody.firstChild) {
    tbody.removeChild(tbody.firstChild);
  }
  // 実行回数の最大値
  var max = 0;
  // 新規データを表示
  records.forEach(record => {
    if (record.count > max) {
      max = record.count;
    }
  });
  // 表示最大値を10の位で切り上げた値に設定
  var maxCount = Math.ceil(max / 10) * 10;

  records.forEach(record => {
    var tr = document.createElement("tr");
    // リポジトリ名
    var td1 = document.createElement("td");
    td1.classList.add("d-block");
    td1.classList.add("w-100");
    td1.textContent = record.repository_name;
    tr.appendChild(td1);
    // 実行回数
    var td2 = document.createElement("td");
    td2.classList.add("d-block");
    td2.classList.add("w-100");
    td2.classList.add("bar-background");
    var span = document.createElement("span");
    var widthRate = record.count / maxCount * 100;
    if (widthRate >= 75) {
      span.classList.add("bar-blue");
    } else if (widthRate >= 50) {
      span.classList.add("bar-green");
    } else {
      span.classList.add("bar-yellow");
    }
    span.classList.add("bar");
    span.style.width = widthRate + "%";
    span.textContent = record.count;
    td2.appendChild(span);
    tr.appendChild(td2);

    tbody.appendChild(tr);
  });
}

// total数を表示する
let displayTotalNum = (totalNum) => {
  var totalRecords = document.getElementById("total");
  totalRecords.textContent = "Total: " + totalNum + "回"
}

// 表示する関数
let display = (data) => {
  displayRedcords(data.repositories);
  displayTotalNum(data.total_count);
}

// 検索する関数
let search = () => {
  console.log("検索します")
  uri = createRequestUri();
  fetch(uri)
    .then((response) => response.json())
    .then((data) => display(data));
}

// 初期化関数
let initCount = () => {
  // 初期値設定
  // 検索開始日を当月1日に設定
  document.getElementById("started_at").value = getFirstDayOfTheMonth();
  // イベントリスナー登録
  // 検索ボタン押下時に検索実行
  searchButton = document.getElementById("search_button");
  searchButton.addEventListener('click', search);
  // 実行履歴更新時に検索実行
  document.addEventListener('history_update', search);
  // 初期化処理
  search();
}

// イベントに対応する処理の追加
window.addEventListener('load', initCount);
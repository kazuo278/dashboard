// 現在のページ数
let currentPageNum = 1;
// 合計ページ数
let totalPageNum = 1;
// 1ページあたりの表示数
const pageSize = 10;
// カラム表示制御チェックボックスのIDプリフィックス
const CHEKC_PREFIX = "check-"
// カラムIDリスト
const columnIdList = [
  "job-id",
  "repository-id",
  "run-id",
  "run-attempt",
  "repository-name",
  "workflow-ref",
  "job-name",
  "status",
  "conclusion",
  "started-at",
  "finished-at"
];

let getToday = () => {
  var today = new Date();
  var year = today.getFullYear().toString().padStart(4, "0");
  var month = (today.getMonth() + 1).toString().padStart(2, "0");
  var day = (today.getDate()).toString().padStart(2, "0");
  return year + "-" + month + "-" + day;
}

// 検索URLを組み立てる関数
let createRequestUri = () => {
  var uri = window.location.protocol + "//" + window.location.host + "/actions";
  var params = new URLSearchParams();

  params.append("limit", pageSize);
  params.append("offset", (currentPageNum - 1) * pageSize);

  var repository_name = document.getElementById("repository_name").value;
  if (repository_name) {
    params.append("repository_name", repository_name);
  }

  var workflow_ref = document.getElementById("workflow_ref").value;
  if (workflow_ref) {
    params.append("workflow_ref", workflow_ref);
  }

  var job_name = document.getElementById("job_name").value;
  if (job_name) {
    params.append("job_name", job_name);
  }

  var started_at = document.getElementById("started_at").value;
  if (started_at) {
    params.append("started_at", started_at + "T00:00:00+09:00");
  }

  var finished_at = document.getElementById("finished_at").value;
  if (finished_at) {
    params.append("finished_at", finished_at + "T00:00:00+09:00");
  }

  var status = document.getElementById("status").value
  if (status !== "ALL") {
    params.append("status", status);
  }

  var conclusion = document.getElementById("conclusion").value
  if (conclusion) {
    params.append("conclusion", conclusion);
  }

  uri += "?" + new URLSearchParams(params).toString();
  return uri;
}

// 検索条件をクリアする関数
let clear = () => {
  document.getElementById("repository_name").value = "";
  document.getElementById("workflow_ref").value = "";
  document.getElementById("job_name").value = "";
  document.getElementById("started_at").value = null;
  document.getElementById("finished_at").value = null;
  document.getElementById("status").value = "ALL";
  document.getElementById("conclusion").value = "";
}

// 単一のテーブルカラムの表示を制御する関数
let updateColumnDisplay = (idName) => {
  if (document.getElementById(CHEKC_PREFIX + idName).checked) {
    Array.from(document.getElementsByClassName(idName)).forEach(function (x) { x.style.display = '' });
  } else {
    Array.from(document.getElementsByClassName(idName)).forEach(function (x) { x.style.display = 'none' });
  }
}

// 日付変換する関数
let formatDate = (dateStr) => {
  if (!dateStr) {
    return "-"
  }
  date = new Date(dateStr);
  var year = date.getFullYear().toString().padStart(4, "0");
  var month = (date.getMonth() + 1).toString().padStart(2, "0");
  var day = date.getDate().toString().padStart(2, "0");
  var hours = date.getHours().toString().padStart(2, "0");
  var minutes = date.getMinutes().toString().padStart(2, "0");
  var secounds = date.getSeconds().toString().padStart(2, "0");

  return year + "/" + month + "/" + day + " " + hours + ":" + minutes + ":" + secounds;
}

// ステータスを変換する関数
let formatStatus = (status) => {
  var span = document.createElement("span");
  span.classList.add("status-icon");
  if (status === "IN_PROGRESS") {
    span.classList.add("in-progress");
    span.textContent = "実行中";
  } else if (status === "QUEUED") {
    span.classList.add("queued");
    span.textContent = "待機中";
  } else if (status === "COMPLETED") {
    span.classList.add("completed");
    span.textContent = "完了";
  } else {
    span.textContent = "不明";
  }
  return span;
}

// 全件数を表示する
let displayTotalRecord = (totalNum) => {
  var totalRecords = document.getElementById("total_records");
  totalRecords.textContent = "全" + totalNum + "件"
}

// JSONデータからテーブルレコードを作成する関数
let displayRedcords = (records) => {
  var tbody = document.getElementById("table_body");
  // 既存データを削除
  while (tbody.firstChild) {
    tbody.removeChild(tbody.firstChild)
  }
  // 新規データを表示
  records.forEach(record => {
    var tr = document.createElement("tr");
    // Job ID
    var tdJobId = document.createElement("td");
    tdJobId.classList.add("job-id");
    // 詳細取得リンクを作成
    var aJobId = document.createElement("a");
    var uriJobId = window.location.protocol + "//" + window.location.host + "/dashboard/detail.html"
    var paramsJobId = new URLSearchParams();
    paramsJobId.append("job_id", record.job_id);
    paramsJobId.append("run_id", record.run_id);
    paramsJobId.append("run_attempt", record.run_attempt);
    aJobId.href = uriJobId + "?" + new URLSearchParams(paramsJobId).toString();
    aJobId.text = record.job_id;
    tdJobId.appendChild(aJobId);
    tr.appendChild(tdJobId);

    // リポジトリID
    var tdRepoId = document.createElement("td");
    tdRepoId.classList.add("repository-id");
    tdRepoId.textContent = record.repository_id;
    tr.appendChild(tdRepoId);

    // RUN　ID
    var tdRunId = document.createElement("td");
    tdRunId.classList.add("run-id");
    var aRunId = document.createElement("a");
    aRunId.href = "https://github.com/" + record.repository_name + "/actions/runs/" + record.run_id;
    aRunId.text = record.run_id;
    aRunId.target = '_blank';
    tdRunId.appendChild(aRunId)
    tr.appendChild(tdRunId);

    // RUN ATTEMPT
    var tdRunAttempt = document.createElement("td");
    tdRunAttempt.classList.add("run-attempt");
    var aRunAttempt = document.createElement("a");
    aRunAttempt.href = "https://github.com/" + record.repository_name + "/actions/runs/" + record.run_id + "/attempts/" + record.run_attempt;
    aRunAttempt.text = record.run_attempt;
    aRunAttempt.target = '_blank';
    tdRunAttempt.appendChild(aRunAttempt)
    tr.appendChild(tdRunAttempt);

    // リポジトリ名
    var tdRepoName = document.createElement("td");
    tdRepoName.classList.add("repository-name");
    var aRepoName = document.createElement("a");
    aRepoName.href = "https://github.com/" + record.repository_name;
    aRepoName.text = record.repository_name;
    aRepoName.target = '_blank';
    tdRepoName.appendChild(aRepoName)
    tr.appendChild(tdRepoName);

    // ワークフローRef
    var tdWorkflowRef = document.createElement("td");
    tdWorkflowRef.textContent = record.workflow_ref;
    tdWorkflowRef.classList.add("workflow-ref");
    tr.appendChild(tdWorkflowRef);

    // JOB名
    var tdJobName = document.createElement("td");
    tdJobName.textContent = record.job_name;
    tdJobName.classList.add("job-name");
    tr.appendChild(tdJobName);

    // 実行ステータス
    var tdStatus = document.createElement("td");
    tdStatus.appendChild(formatStatus(record.status));
    tdStatus.classList.add("status");
    tr.appendChild(tdStatus);

    // 実行結果
    var tdConclusion = document.createElement("td");
    tdConclusion.textContent = record.conclusion;
    tdConclusion.classList.add("conclusion");
    tr.appendChild(tdConclusion);

    // 開始日時
    var tdStartedAt = document.createElement("td");
    tdStartedAt.textContent = formatDate(record.started_at);
    tdStartedAt.classList.add("started-at");
    tr.appendChild(tdStartedAt);

    // 終了日時
    var tdFinishedAt = document.createElement("td");
    tdFinishedAt.textContent = formatDate(record.finished_at);
    tdFinishedAt.classList.add("finished-at");
    tr.appendChild(tdFinishedAt);

    tbody.appendChild(tr);
  });
  // カラム表示項目の初期化
  columnIdList.forEach(idName => updateColumnDisplay(idName));
}

// ページネーターを作成する関数
let displayPageNation = (totalNum) => {
  var navList = document.getElementById("page_list");
  // 既存データを削除
  while (navList.firstChild) {
    navList.removeChild(navList.firstChild)
  }
  // 合計ページ数
  totalPageNum = Math.ceil(totalNum / pageSize);

  // 表示するページリンクの最大数
  const displayPageLinkMaxNum = 20
  // 開始ページ番号
  startPageNum = currentPageNum - Math.floor((displayPageLinkMaxNum - 1) / 2);
  // 終了ページ番号
  endPageNum = currentPageNum + Math.ceil((displayPageLinkMaxNum - 1) / 2);

  // 表示範囲をプラス方向にずらす数を計算(startが１ページよりも小さくなる場合)
  needShiftRightNum = 0;
  if (startPageNum < 1) {
    needShiftRightNum = 1 - startPageNum;
  }
  // 表示範囲をマイナス方向にずらす数を計算(endPageが合計ページ数よりも大きくなる場合)
  needShiftLeftNum = 0;
  if (endPageNum > totalPageNum) {
    needShiftLeftNum = endPageNum - totalPageNum;
  }
  // 表示範囲を調整
  startPageNum = startPageNum + needShiftRightNum - needShiftLeftNum;
  endPageNum = endPageNum + needShiftRightNum - needShiftLeftNum;
  if (startPageNum < 1) {
    startPageNum = 1;
  }
  if (endPageNum > totalPageNum) {
    endPageNum = totalPageNum;
  }
  // ページネータ作成
  // 先頭ページ用ボタン
  var liFirst = document.createElement("li");
  liFirst.classList.add("page-item");
  if (currentPageNum == 1) {
    liFirst.classList.add("disabled");
  }
  var aFirst = document.createElement("a");
  aFirst.href = "#";
  aFirst.id = "page_first";
  aFirst.classList.add("page-link");
  aFirst.textContent = "<<";
  // クリック時に選択したページで再表示させる
  aFirst.addEventListener('click', () => {
    currentPageNum = 1;
    search();
  })
  liFirst.appendChild(aFirst)
  navList.appendChild(liFirst);
  // 個別ページボタン
  for (let page = 1; page <= totalPageNum; page++) {
    var li = document.createElement("li");
    li.classList.add("page-item");
    if (page == currentPageNum) {
      li.classList.add("active");
    }
    var a = document.createElement("a");
    a.href = "#";
    a.id = "page_" + page;
    a.classList.add("page-link");
    a.textContent = page;
    // 表示範囲外は非表示CSSを適用
    if (page < startPageNum || page > endPageNum) {
      a.classList.add("d-none")
    }
    // クリック時に選択したページで再表示させる
    a.addEventListener('click', (event) => {
      currentPageNum = Number(document.getElementById(event.target.id).textContent);
      search();
    })
    li.appendChild(a)
    navList.appendChild(li);
  }
  // 最終ページ用ボタン
  var liLast = document.createElement("li");
  liLast.classList.add("page-item");
  if (currentPageNum == totalPageNum || totalPageNum == 0) {
    liLast.classList.add("disabled");
  }
  var aLast = document.createElement("a");
  aLast.href = "#";
  aLast.id = "page_last";
  aLast.classList.add("page-link");
  aLast.textContent = ">>";
  // クリック時に選択したページで再表示させる
  aLast.addEventListener('click', () => {
    currentPageNum = totalPageNum;
    search();
  })
  liLast.appendChild(aLast)
  navList.appendChild(liLast);
}

// 表示する関数
let display = (data) => {
  displayRedcords(data.jobs);
  displayPageNation(data.count);
  displayTotalRecord(data.count);
}

// 検索する関数
let search = () => {
  console.log("検索します")
  uri = createRequestUri();
  fetch(uri)
    .then((response) => response.json())
    .then((data) => display(data));
}

// 現在のページ数を1に初期化する関数
let initCurrentPageNum = () => {
  currentPageNum = 1;
}

// 初期化関数
let initDashboard = () => {
  // 開始日を当日に変更
  document.getElementById("started_at").value = getToday();
  // イベントリスナー登録
  // 検索ボタン押下時に検索実行
  searchButton = document.getElementById("search_button");
  searchButton.addEventListener('click', search);
  searchButton.addEventListener('click', initCurrentPageNum);
  // クリアボタン押下時に検索条件をクリア
  clearButton = document.getElementById("clear_button");
  clearButton.addEventListener('click', clear);
  // カラム表示制御チェックボックスに表示制御関数を登録
  columnIdList.forEach(idName => {
    document.getElementById(CHEKC_PREFIX + idName).addEventListener('change', {
      arg: idName,
      handleEvent: function () { updateColumnDisplay(this.arg) }
    });
  });

  // 実行履歴更新時に検索実行
  document.addEventListener('history_update', search);
  // 初期化処理
  search();
  // カラム表示項目の初期化
  columnIdList.forEach(idName => updateColumnDisplay(idName));
}

// イベントに対応する処理の追加
window.addEventListener('load', initDashboard);


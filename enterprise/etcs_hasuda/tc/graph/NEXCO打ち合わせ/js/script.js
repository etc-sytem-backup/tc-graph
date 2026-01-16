// 最終受信時刻
const updateReceiveTime = (time, type) => {
  if(!time || !type) return
  document.getElementById("last-receive-time-" + type).innerText = time;
}

/* ========== メイン画面 ========== */
// 大型駐車台数
const updateLCarCnt = (c) => {
  if(!c) return
  document.getElementById("l-car-cnt").innerText = c;
}

// 小型駐車台数
const updateSCarCnt = (c) => {
  if(!c) return
  document.getElementById("s-car-cnt").innerText = c;
}

// 大型駐車室数
const updateLParkCnt = (c) => {
  if(!c) return
  document.getElementById("l-park-cnt").innerText = c;
}

// 小型駐車室数
const updateSParkCnt = (c) => {
  if(!c) return
  document.getElementById("s-park-cnt").innerText = c;
}

// 大型満車フラグ
const updateLMansyaflg = (f) => {
  if(!f) return
  if(f === "1") {
    document.getElementById("l-mansya-flag").innerText = "満"
  } else {
    document.getElementById("l-mansya-flag").innerText = ""
  }
}

// 小型満車フラグ
const updateSMansyaflg = (f) => {
  if(!f) return
  if(f === "1") {
    document.getElementById("s-mansya-flag").innerText = "満"
  } else {
    document.getElementById("s-mansya-flag").innerText = ""
  }
}

// 大型満空率
const updateLMankuRatio = (r, f) => {
  if(!r || !f) return
  document.getElementById("l-manku-ratio").innerText = r;
  if(f === "1") {
    document.getElementById("l-manku-ratio").classList.remove("text-primary")
    document.getElementById("l-manku-ratio").classList.add("text-danger")
  } else {
    document.getElementById("l-manku-ratio").classList.remove("text-danger")
    document.getElementById("l-manku-ratio").classList.add("text-primary")
  }
}

// 小型満空率
const updateSMankuRatio = (r, f) => {
  if(!r || !f) return
  document.getElementById("s-manku-ratio").innerText = r;
  if(f === "1") {
    document.getElementById("s-manku-ratio").classList.remove("text-primary")
    document.getElementById("s-manku-ratio").classList.add("text-danger")
  } else {
    document.getElementById("s-manku-ratio").classList.remove("text-danger")
    document.getElementById("s-manku-ratio").classList.add("text-primary")
  }
}

// 大型超過台数
const updateLTyoukaCnt = (c) => {
  if(!c) return
  document.getElementById("l-tyouka-cnt").innerText = c;
}

// 小型超過台数
const updateSTyoukaCnt = (c) => {
  if(!c) return
  document.getElementById("s-tyouka-cnt").innerText = c;
}

// 大型パス台数
const updateLPassCnt = (c) => {
  if(!c) return
  document.getElementById("l-pass-cnt").innerText = c;
}

// 小型パス台数
const updateSPassCnt = (c) => {
  if(!c) return
  document.getElementById("s-pass-cnt").innerText = c;
}

// 電波ステータス
const updateRadioStatus = (s) => {
  if(!s) return
  if(s === "1") {
    document.getElementById("radio-status").innerText = "発射中";
  } else if(s === "0") {
    document.getElementById("radio-status").innerText = "停止中";
  } else {
    document.getElementById("radio-status").innerText = "-";
  }
}

/* ========== 統計(平均)画面 ========== */

const updateJamFlg = (f) => {
  if(!f) return
  if(f === "1") {
    document.getElementById("jam-flg").innerText = "●"
  } else {
    document.getElementById("jam-flg").innerText = "　"
  }
}
const updateAParkingMinAvgDay = (c) => {
  if(!c) return
  document.getElementById("a-parking-min-avg-day").innerText = c + " min"
}

const updateAParkingMinAvgWeek = (c) => {
  if(!c) return
  document.getElementById("a-parking-min-avg-week").innerText = c + " min"
}

const updateLParkingMinAvgDay = (c) => {
  if(!c) return
  document.getElementById("l-parking-min-avg-day").innerText = c + " min"
}

const updateLParkingMinAvgWeek = (c) => {
  if(!c) return
  document.getElementById("l-parking-min-avg-week").innerText = c + " min"
}

const updateSParkingMinAvgDay = (c) => {
  if(!c) return
  document.getElementById("s-parking-min-avg-day").innerText = c + " min"
}

const updateSParkingMinAvgWeek = (c) => {
  if(!c) return
  document.getElementById("s-parking-min-avg-week").innerText = c + " min"
}

/* ========== 統計(一覧)画面 ========== */

const initTableRows = () => {
  document.querySelectorAll(".disp-table-rows").forEach((e, i) => {
    let tableHTML = "";
    for (var j = 0; j < 200; j++) {
        tableHTML += "<tr>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
        + "</tr>"
    }
    e.innerHTML = tableHTML;
  })
}

const updateTableWeek = (data) => {
  // 全て文字列
  // data[i][0]: CSVファイル更新時刻
  // data[i][1]: 駐車開始時刻
  // data[i][2]: WCN番号
  // data[i][3]: ETCカード番号
  // data[i][4]: 支局
  // data[i][5]: 用途
  // data[i][6]: 種別
  // data[i][7]: 一連番号
  // data[i][8]: 回数
  updateReceiveTime(data[0][0], "week");
  data.forEach((v, i) => {
    if(i >= 200) {
      return
    }
    let tr = document.getElementById("disp-table-week").children[i];
    for(let j = 0; j < 8; j++) {
      tr.children[j].innerText = v[j+1];
    }
  })
}

const updateTableMonth = (data) => {
  // 全て文字列
  // data[i][0]: CSVファイル更新時刻
  // data[i][1]: 駐車開始時刻
  // data[i][2]: WCN番号
  // data[i][3]: ETCカード番号
  // data[i][4]: 支局
  // data[i][5]: 用途
  // data[i][6]: 種別
  // data[i][7]: 一連番号
  // data[i][8]: 回数
  updateReceiveTime(data[0][0], "month");
  data.forEach((v, i) => {
    if(i >= 200) {
      return
    }
    let tr = document.getElementById("disp-table-month").children[i];
    for(let j = 0; j < 8; j++) {
      tr.children[j].innerText = v[j+1];
    }
  })
}

const updateParkingTime = (data) => {
  // 全て文字列
  // data[i][0]: CSVファイル更新時刻
  // data[i][1]: 駐車開始時刻
  // data[i][2]: WCN番号
  // data[i][3]: ETCカード番号
  // data[i][4]: 支局
  // data[i][5]: 用途
  // data[i][6]: 種別
  // data[i][7]: 一連番号
  // data[i][8]: 駐車時間
  updateReceiveTime(data[0][0], "pariking");
  data.forEach((v, i) => {
    if(i >= 200) {
      return
    }
    let tr = document.getElementById("disp-parking-time").children[i];
    for(let j = 0; j < 8; j++) {
      if(j+1 === 8) {
          let parsedInt = parseInt(v[j+1]);
          if(!isNaN(parsedInt)) {
              if(parseInt(v[j+1]) >= 1440) {
                  tr.classList.add("text-danger");
              }
              tr.children[j].innerText = convertMinutes(parseInt(v[j+1]));
          }
      } else {
        tr.classList.remove("text-danger");
        tr.children[j].innerText = v[j+1];
      }
    }
  })
}

function convertMinutes(minutes) {
  // 時間数を計算する
  let hh = Math.floor(minutes / 60);
  // 分数を計算する
  let mm = minutes % 60;

  // 'hh:mm' 形式に変換する
  return `${hh.toString().padStart(2, '0')}:${mm.toString().padStart(2, '0')}`;
} 

/* ========== 統計(逆走検知)画面 ========== */

const initReverseRows = () => {
  document.querySelectorAll(".disp-reverse-rows").forEach((e, i) => {
    let tableHTML = "";
    for (var j = 0; j < 200; j++) {
        tableHTML += "<tr>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
          + "<td>　</td>"
        + "</tr>"
    }
    e.innerHTML = tableHTML;
  })
}

const updateReverse = (data) => {
  // 全て文字列
  // data[i][0]: CSVファイル更新時刻
  // data[i][1]: 駐車開始時刻
  // data[i][2]: WCN番号
  // data[i][3]: ETCカード番号
  // data[i][4]: 支局
  // data[i][5]: 用途
  // data[i][6]: 種別
  // data[i][7]: 一連番号
  // data[i][8]: アンテナ
  updateReceiveTime(data[0][0], "reverse");
  console.log(data[0]);
  data.forEach((v, i) => {
    if(i >= 200) {
      return
    }
    let tr = document.getElementById("disp-reverse").children[i];
    for(let j = 0; j < 8; j++) {
      tr.children[j].innerText = v[j+1];
    }
  })
}

/* ========== 断面交通量画面 ========== */

// アンテナ1
const updatePassageCnt1 = (c) => {
  if(!c) return
  document.getElementById("ant-1-cnt").innerText = c;
}

// アンテナ2
const updatePassageCnt2 = (c) => {
  if(!c) return
  document.getElementById("ant-2-cnt").innerText = c;
}

// アンテナ3
const updatePassageCnt3 = (c) => {
  if(!c) return
  document.getElementById("ant-3-cnt").innerText = c;
}

/* ========== 設定画面 ========== */

// 大型駐車台数オフセット
const updateLParkingOffset = (c) => {
  if(!c) return
  document.getElementById("l-park-offset").value = c;
  document.getElementById("l-park-offset-value").textContent = c;
}

// 小型駐車台数オフセット
const updateSParkingOffset = (c) => {
  if(!c) return
  document.getElementById("s-park-offset").value = c;
  document.getElementById("s-park-offset-value").textContent = c;
}
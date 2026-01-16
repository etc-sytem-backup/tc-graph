package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SessionDescription struct {
	SDP  string `json:"sdp"`
	Type string `json:"type"`
}

type EnterUserInfo struct {
	ETCNumber  string `json:"etcNumber"`
	Department string `json:"department"`
	Name       string `json:"name"`
	LastUpdate string `json:"lastUpdate"`
	IsVisitor  bool   `json:"isVisitor"`
	IsEnter    bool   `json:"isEnter"`
}

type ExitUserInfo struct {
	ETCNumber  string
	Department string
	Name       string
	LastUpdate string
	IsVisitor  bool
	IsExit     bool
}

type CarRecord struct {
	ETCNumber  string
	Name       string
	Department string
	Status     string
	EntryTime  string
	ExitTime   string
	StayTime   string
}

type VehicleInfo struct {
	ETCNumber    string `json:"etcNumber"`
	Name         string `json:"name"`
	Department   string `json:"department"`
	IsEnter      bool   `json:"isEnter"`
	SerialNumber string `json:"serialNumber"`
	VehicleType  string `json:"vehicleType"`
	PassingTime  string `json:"passingTime"`
}

type SystemStats struct {
	TodayEntries      int `json:"today_entries"`
	TodayExits        int `json:"today_exits"`
	CurrentInside     int `json:"current_inside"`
	VisitorCount      int `json:"visitor_count"`
	UnauthorizedCount int `json:"unauthorized_count"`
}

var (
	videoLock sync.Mutex
	videoPath = "/home/mizukami-ryosuke/sk_prj/enterprise/etcs_hasuda/tc/gatesystem/www/hls/output.mp4"
	// 最新検出車両情報を保持
	currentVehicle *VehicleInfo
	vehicleMutex   sync.RWMutex
	// システム統計情報
	systemStats SystemStats
	statsMutex  sync.RWMutex
)

func init() {
	// TODO:カメラ映像を表示するときはコメントアウトを外す
	// go convertRTSPtoMP4() // 非同期で変換開始
	// TODO:カメラ画像を解析するときはコメントアウトを外す
	// go imageRecognition()
}

// TODO:カメラ映像を表示するときはコメントアウトを外す
// func convertRTSPtoMP4() {
// 	videoLock.Lock()
// 	defer videoLock.Unlock()

// 	// 既存ファイル削除
// 	if _, err := os.Stat(videoPath); err == nil {
// 		os.Remove(videoPath)
// 	}

// 	cmd := exec.Command(
// 		"ffmpeg",
// 		"-i", "rtsp://admin:etcs-Mizukami@192.168.110.64:554/Streaming/channels/101",
// 		"-c:v", "copy",
// 		"-c:a", "aac",
// 		"-f", "mp4",
// 		"-movflags", "frag_keyframe+empty_moov",
// 		videoPath,
// 	)

// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr

// 	log.Println("RTSP→MP4変換プロセスを開始します...")
// 	if err := cmd.Run(); err != nil {
// 		log.Printf("変換失敗: %v\n", err)
// 	} else {
// 		log.Println("MP4ファイルが正常に生成されました")
// 	}
// }

func main() {
	// 静的ファイル配信（wwwディレクトリ内の全てのファイルを配信）
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("www"))))
	http.Handle("/hls/", http.StripPrefix("/hls/",
		http.FileServer(http.Dir(filepath.Dir(videoPath)))))

	// ルートパスは待機画面（ホームダッシュボード）を表示
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータで強制的に待機画面を表示する場合
		if r.URL.Path == "/" || r.URL.Path == "/home" {
			handleHomeDashboard(w, r)
			return
		}
		// 他のパスは静的ファイルとして扱う
		http.FileServer(http.Dir("www")).ServeHTTP(w, r)
	})

	// 待機画面（ホームダッシュボード）
	http.HandleFunc("/home", handleHomeDashboard)

	// 許可車両画面
	http.HandleFunc("/authorized", handleAuthorizedVehicle)

	// 来客車両画面
	http.HandleFunc("/visitor", handleVisitorDetection)

	// 許可車両一覧（既存のページ）
	http.HandleFunc("/index", handleIndex)
	http.HandleFunc("/vehicle-list", func(w http.ResponseWriter, r *http.Request) {
		// 既存の許可車両一覧ページへリダイレクト
		http.Redirect(w, r, "/index", http.StatusFound)
	})

	// APIエンドポイント
	http.HandleFunc("/api/user-info", handleUserInfo)
	http.HandleFunc("/api/current-vehicle", handleCurrentVehicle)
	http.HandleFunc("/api/detected-vehicle", handleDetectedVehicle)
	http.HandleFunc("/api/stats", handleStats)
	http.HandleFunc("/api/vehicle-detected", handleVehicleDetected)
	http.HandleFunc("/api/register-visitor", handleRegisterVisitor)
	// http.HandleFunc("/api/latest-image", handleLatestImage)

	// システム制御エンドポイント
	// http.HandleFunc("/webrtc", handleWebRTC)
	http.HandleFunc("/open-gate", handleOpenGate)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/register", handleRegister)

	// 初期統計情報の更新を開始
	go updateStatsPeriodically()

	log.Println("サーバー起動: http://localhost:8080")
	log.Println("待機画面: http://localhost:8080/home")
	log.Println("許可車両画面: http://localhost:8080/authorized")
	log.Fatal(http.ListenAndServe(":8080", nil))
	// var enterUserInfo EnterUserInfo

	// 静的ファイル配信
	http.Handle("/www/", http.StripPrefix("/www/", http.FileServer(http.Dir("www"))))
	http.Handle("/hls/", http.StripPrefix("/hls/",
		http.FileServer(http.Dir(filepath.Dir(videoPath)))))

	// // メイン画面
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl := template.Must(template.ParseFiles("www/main.html"))

	// 	enterUserInfo = EnterUserInfo{
	// 		ETCNumber:  "",
	// 		Department: "",
	// 		Name:       "",
	// 		LastUpdate: time.Now().Format("15:04"),
	// 		IsVisitor:  false,
	// 		IsEnter:    true,
	// 	}

	// 	data := struct {
	// 		UserInfo    EnterUserInfo
	// 		AutoRefresh bool
	// 	}{
	// 		UserInfo:    enterUserInfo,
	// 		AutoRefresh: true,
	// 	}
	// 	tmpl.Execute(w, data)
	// })

	// // 許可車両画面
	// http.HandleFunc("/authorized", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl, err := template.ParseFiles("www/authorized_vehicle.html")
	// 	if err != nil {
	// 		log.Printf("テンプレート読み込みエラー: %v", err)
	// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	// 現在の車両情報を取得
	// 	vehicleMutex.RLock()
	// 	vehicle := currentVehicle
	// 	vehicleMutex.RUnlock()

	// 	// 車両情報がない場合はデフォルト値を設定
	// 	if vehicle == nil {
	// 		vehicle = &VehicleInfo{
	// 			ETCNumber:   "N/A",
	// 			Name:        "N/A",
	// 			Department:  "N/A",
	// 			IsEnter:     true,
	// 			PassingTime: time.Now().Format("15:04:05"),
	// 		}
	// 	}

	// 	tmpl.Execute(w, vehicle)
	// })

	// // ユーザー情報取得
	// http.HandleFunc("/api/user-info", handleUserInfo)

	// // 現在の車両情報取得API
	// http.HandleFunc("/api/current-vehicle", func(w http.ResponseWriter, r *http.Request) {
	// 	vehicleMutex.RLock()
	// 	vehicle := currentVehicle
	// 	vehicleMutex.RUnlock()

	// 	if vehicle == nil {
	// 		http.Error(w, "No vehicle detected", http.StatusNotFound)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	json.NewEncoder(w).Encode(vehicle)
	// })

	// http.HandleFunc("/webrtc", func(w http.ResponseWriter, r *http.Request) {
	// 	// CORS設定（開発用）
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 	if r.Method == "OPTIONS" {
	// 		return
	// 	}

	// 	var offer SessionDescription
	// 	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
	// 		log.Printf("Offer decode error: %v", err)
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 		return
	// 	}

	// 	// MediaMTXへのリクエストを修正
	// 	reqBody := strings.NewReader(fmt.Sprintf(`{"type":"%s","sdp":"%s"}`, offer.Type, strings.ReplaceAll(offer.SDP, "\r\n", "\\r\\n")))

	// 	req, err := http.NewRequest("POST", "http://localhost:8889/cam/whep", reqBody)
	// 	if err != nil {
	// 		log.Printf("Request creation error: %v", err)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	req.Header.Set("Content-Type", "application/json")

	// 	client := &http.Client{}
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		log.Printf("MediaMTX connection error: %v", err)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// 	defer resp.Body.Close()

	// 	if resp.StatusCode != http.StatusCreated {
	// 		body, _ := io.ReadAll(resp.Body)
	// 		log.Printf("MediaMTX error response: %s", string(body))
	// 		http.Error(w, string(body), resp.StatusCode)
	// 		return
	// 	}

	// 	var answer SessionDescription
	// 	if err := json.NewDecoder(resp.Body).Decode(&answer); err != nil {
	// 		log.Printf("Answer decode error: %v", err)
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	if err := json.NewEncoder(w).Encode(answer); err != nil {
	// 		log.Printf("Answer encode error: %v", err)
	// 	}
	// })

	// // 画像取得
	// http.HandleFunc("/api/latest-image", handleLatestImage)

	// // ゲートオープン
	// http.HandleFunc("/open-gate", func(w http.ResponseWriter, r *http.Request) {
	// 	gateOpenCommand()
	// 	registerVisitor(getLatestEnterWCN())

	// 	w.Write([]byte("ゲートを開放しました"))
	// })
	// http.HandleFunc("/index", handleIndex)
	// http.HandleFunc("/list", handleList)
	// http.HandleFunc("/delete", handleDelete)
	// http.HandleFunc("/register", handleRegister)

	// log.Println("サーバー起動: http://localhost:8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}

// 待機画面（ホームダッシュボード）ハンドラ
func handleHomeDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("www/home_dashboard.html")
	if err != nil {
		log.Printf("テンプレート読み込みエラー: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 統計情報を更新
	updateStats()

	statsMutex.RLock()
	stats := systemStats
	statsMutex.RUnlock()

	tmpl.Execute(w, stats)
}

// 許可車両画面ハンドラ
func handleAuthorizedVehicle(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("www/authorized_vehicle.html")
	if err != nil {
		log.Printf("テンプレート読み込みエラー: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 現在の車両情報を取得
	vehicleMutex.RLock()
	vehicle := currentVehicle
	vehicleMutex.RUnlock()

	// 車両情報がない場合はデフォルト値を設定
	if vehicle == nil {
		vehicle = &VehicleInfo{
			ETCNumber:   "N/A",
			Name:        "N/A",
			Department:  "N/A",
			IsEnter:     true,
			PassingTime: time.Now().Format("15:04:05"),
		}
	}

	tmpl.Execute(w, vehicle)
}

// 現在の車両情報取得API
func handleCurrentVehicle(w http.ResponseWriter, r *http.Request) {
	vehicleMutex.RLock()
	vehicle := currentVehicle
	vehicleMutex.RUnlock()

	if vehicle == nil {
		http.Error(w, "No vehicle detected", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

// 来客車両画面ハンドラ
func handleVisitorDetection(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("www/visitor_detection.html")
	if err != nil {
		log.Printf("テンプレート読み込みエラー: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 現在の車両情報を取得
	vehicleMutex.RLock()
	vehicle := currentVehicle
	vehicleMutex.RUnlock()

	// 車両情報がない場合はデフォルト値を設定
	if vehicle == nil {
		vehicle = &VehicleInfo{
			ETCNumber:   getLatestEnterWCN(),
			Name:        "来客車",
			Department:  "来客車",
			IsEnter:     true,
			PassingTime: time.Now().Format("15:04:05"),
		}
	}

	tmpl.Execute(w, vehicle)
}

// 来客車両登録API
func handleRegisterVisitor(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var visitorData struct {
		ETCNumber    string `json:"etcNumber"`
		SerialNumber string `json:"serialNumber"`
		VehicleType  string `json:"vehicleType"`
	}

	if err := json.NewDecoder(r.Body).Decode(&visitorData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 来客車両をCSVに登録
	if visitorData.ETCNumber != "" {
		registerVisitor(visitorData.ETCNumber)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "来客車両を登録しました",
	})
}

// 統計情報取得API
func handleStats(w http.ResponseWriter, r *http.Request) {
	updateStats()

	statsMutex.RLock()
	stats := systemStats
	statsMutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// 車両検出API（フロントエンドから定期的に呼び出される）
func handleVehicleDetected(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 最新のWCNを取得
	etcNumber := getLatestEnterWCN()
	if etcNumber == "" {
		// 車両未検出
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"detected": false,
		})
		return
	}

	// 車両情報を取得
	enterUserInfo := fetchLatestEnterUserInfo()

	// 車両検出処理
	if enterUserInfo.ETCNumber != "" && enterUserInfo.Department != "未検出" {
		// 車両情報を保存
		vehicleMutex.Lock()
		currentVehicle = &VehicleInfo{
			ETCNumber:    enterUserInfo.ETCNumber,
			Name:         enterUserInfo.Name,
			Department:   enterUserInfo.Department,
			IsEnter:      enterUserInfo.IsEnter,
			SerialNumber: extractSerialNumber(enterUserInfo.ETCNumber),
			VehicleType:  getRandomVehicleType(),
			PassingTime:  time.Now().Format("15:04:05"),
		}
		vehicleMutex.Unlock()

		// 許可車両かどうかを判定
		isAuthorized := !enterUserInfo.IsVisitor && enterUserInfo.Department != "来客車"

		// レスポンスデータを作成
		responseData := map[string]interface{}{
			"detected":     true,
			"etcNumber":    enterUserInfo.ETCNumber,
			"department":   enterUserInfo.Department,
			"name":         enterUserInfo.Name,
			"isVisitor":    enterUserInfo.IsVisitor,
			"isEnter":      enterUserInfo.IsEnter,
			"isAuthorized": isAuthorized,
		}

		// 画面遷移の決定
		if isAuthorized {
			// 許可車両の場合
			responseData["redirect"] = "/authorized"
			go gateOpenCommand()
			updateUser(enterUserInfo)
		} else {
			// 来客車両の場合
			responseData["redirect"] = "/visitor"
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)

	} else {
		// 車両未検出
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"detected": false,
		})
	}
}

// 検出車両情報取得API
func handleDetectedVehicle(w http.ResponseWriter, r *http.Request) {
	vehicleMutex.RLock()
	vehicle := currentVehicle
	vehicleMutex.RUnlock()

	if vehicle == nil {
		// 最新のWCNを取得して仮の車両情報を作成
		etcNumber := getLatestEnterWCN()
		if etcNumber != "" {
			vehicle = &VehicleInfo{
				ETCNumber:    etcNumber,
				SerialNumber: extractSerialNumber(etcNumber),
				VehicleType:  "SUV/MPV", // デフォルト値
				PassingTime:  time.Now().Format("15:04:05"),
			}
		}
	}

	if vehicle == nil {
		http.Error(w, "No vehicle detected", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vehicle)
}

// ETC番号からシリアル番号を抽出（仮実装）
func extractSerialNumber(etcNumber string) string {
	if len(etcNumber) >= 3 {
		// 最後の3桁をシリアル番号として使用
		return etcNumber[len(etcNumber)-3:]
	}
	return "000"
}

// 許可車両かどうかで画面を分岐
func handleUserInfo(w http.ResponseWriter, r *http.Request) {
	enterUserInfo := fetchLatestEnterUserInfo()

	// 車両情報を更新
	updateUser(enterUserInfo)

	// gateOpenCommand()
	// w.Header().Set("Content-Type", "application/json")
	// w.Write([]byte(`{"type":"registered","department":"` + enterUserInfo.Department + `","name":"` + enterUserInfo.Name + `"}`))
	// time.Sleep(10 * time.Second)
	// gateCloseCommand()

	// 現在の車両情報を保存
	vehicleMutex.Lock()
	currentVehicle = &VehicleInfo{
		ETCNumber:   enterUserInfo.ETCNumber,
		Name:        enterUserInfo.Name,
		Department:  enterUserInfo.Department,
		IsEnter:     enterUserInfo.IsEnter,
		PassingTime: time.Now().Format("15:04:05"),
	}
	vehicleMutex.Unlock()

	// レスポンスデータを作成
	responseData := map[string]interface{}{
		"etcNumber":     enterUserInfo.ETCNumber,
		"department":    enterUserInfo.Department,
		"name":          enterUserInfo.Name,
		"isVisitor":     enterUserInfo.IsVisitor,
		"isEnter":       enterUserInfo.IsEnter,
		"detectionTime": time.Now().Format("15:04:05"),
	}

	// 許可車両の場合（来客車でない場合）
	if !enterUserInfo.IsVisitor && enterUserInfo.Department != "未検出" && enterUserInfo.Department != "来客車" {
		// 自動的にゲートを開ける
		gateOpenCommand()
		updateUser(enterUserInfo)
		responseData["redirect"] = "/authorized"
	} else {
		// 来客車または未検出の場合は従来のレスポンス
		responseData["type"] = "registered"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func handleLatestImage(w http.ResponseWriter, r *http.Request) {
	imagePath := findClosestImage()
	http.ServeFile(w, r, imagePath)
}

func findClosestImage() string {
	dir := "/home/mizukami-ryosuke/IP CAMERA/01"
	files, _ := os.ReadDir(dir)

	var closestFile string
	var smallestDiff time.Duration = 1<<63 - 1 // 最大値
	currentTime := time.Now()

	for _, file := range files {
		filename := file.Name()
		if strings.Contains(filename, "VEHICLE_PICTURE") {
			// ファイル名から日時部分を抽出（例: _01_20250706233616542_）
			parts := strings.Split(filename, "_")
			if len(parts) < 3 {
				continue
			}

			// 日時文字列を解析（YYYYMMDDhhmmssSSS形式）
			timeStr := parts[2]
			if len(timeStr) != 17 {
				continue
			}

			// タイムゾーンを考慮して解析
			imgTime, err := time.ParseInLocation(
				"20060102150405.000",
				timeStr[:14]+"."+timeStr[14:],
				time.Local,
			)
			if err != nil {
				log.Printf("時間解析エラー: %v (ファイル: %s)", err, filename)
				continue
			}

			// 現在時刻との差分計算
			diff := currentTime.Sub(imgTime).Abs()
			if diff < smallestDiff {
				smallestDiff = diff
				closestFile = filepath.Join(dir, filename)
			}
		}
	}

	if closestFile == "" {
		log.Println("該当する画像ファイルが見つかりません")
	}
	return closestFile
}

// 画像認識処理
func imageRecognition() {
	pythonPath := filepath.Join(os.Getenv("HOME"), "ocr_project/venv/bin/python")
	scriptPath := filepath.Join(os.Getenv("HOME"), "ocr_project/ocr_monitor.py")

	// 画像認識処理から抜き出す値のデフォルト値
	direction := "Unknown"
	timestamp := time.Now().Format("15:04") // 現在時刻をデフォルト

	cmd := exec.Command(pythonPath, scriptPath)
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = os.Stderr

	// プロセス開始
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start OCR monitor: %v", err)
	}

	// 出力をリアルタイムで処理
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		// 方向情報を抽出
		dirRe := regexp.MustCompile(`Direction:\s*(Forward|Reverse|Unknown)`)
		if match := dirRe.FindStringSubmatch(line); len(match) > 1 {
			direction = match[1]
		}

		// タイムスタンプを抽出してフォーマット変換
		timeRe := regexp.MustCompile(`timestamp['"]?:\s*['"]?(\d{4}-\d{2}-\d{2}\s+(\d{2}:\d{2}):\d{2})`)
		if match := timeRe.FindStringSubmatch(line); len(match) > 2 {
			// "hh:mm" 部分だけを取得
			timestamp = match[2]
		}

		// 結果行だけを処理
		if strings.Contains(line, "Result:") {
			fmt.Printf("[%s] 方向: %s\n", timestamp, direction)
		} else {
			// 他のログも表示（オプション）
			fmt.Println(line)
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("OCR monitor exited with error: %v", err)
	}

	// 進行方向、通過時刻を返す
	// return direction, timestamp

	// enterUserInfo := EnterUserInfo{
	// 	ETCNumber:  etcNumber,
	// 	Department: "登録済み",
	// 	Name:       "来客車",
	// 	LastUpdate: time.Now().Format("15:04"),
	// 	IsVisitor:  false,
	// 	IsEnter:    true,
	// }
	// updateUser()
}

// 統計情報を定期的に更新
func updateStatsPeriodically() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		updateStats()
	}
}

// 統計情報を計算して更新
func updateStats() {
	// WCN登録リストから統計を計算
	csvPath := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")
	file, err := os.Open(csvPath)
	if err != nil {
		log.Printf("統計情報更新エラー: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("CSV読み込みエラー: %v", err)
		return
	}

	var todayEntries, todayExits, currentInside, visitorCount, unauthorizedCount int
	// todayDate := time.Now().Format("2006-01-02")

	for i, record := range records {
		if i == 0 {
			continue // ヘッダー行をスキップ
		}

		if len(record) >= 5 {
			entryTime := record[4]

			// 本日の入場数を計算
			if record[3] == "入場" && strings.HasPrefix(entryTime, time.Now().Format("15")) {
				todayEntries++
			}

			// 本日の退場数を計算
			if record[3] == "退場" && len(record) > 5 && strings.HasPrefix(record[5], time.Now().Format("15")) {
				todayExits++
			}

			// 現在構内にいる車両数を計算
			if record[3] == "入場" && (len(record) <= 5 || record[5] == "") {
				currentInside++
			}

			// 来客車両数を計算
			if len(record) >= 2 && record[1] == "来客車" {
				visitorCount++
			}
		}
	}

	statsMutex.Lock()
	systemStats = SystemStats{
		TodayEntries:      todayEntries,
		TodayExits:        todayExits,
		CurrentInside:     currentInside,
		VisitorCount:      visitorCount,
		UnauthorizedCount: unauthorizedCount,
	}
	statsMutex.Unlock()
}

// 最新の入口側ゲート待機中のWCNを取得
func getLatestEnterWCN() string {

	// ファイルオープン処理
	path := filepath.Join(os.Getenv("HOME"), "opt", "aps", "sbox01", "tc_csv_table", "WCN_rireki.csv")
	file, err := os.Open(path)
	if err != nil {
		log.Println("WCN_rireki.csvが開けません:", err)
		return ""
	}
	defer file.Close()

	// ファイルの最後から読み取り（sbox01/tc_csv_table/WCN_rireki.csvは最後の行が最新読み取り行）
	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if lastLine == "" {
		return ""
	}

	// コンマ区切りで分割
	columns := strings.Split(lastLine, ",")
	if len(columns) >= 4 {
		previousFormatWCN := strings.ReplaceAll(columns[3], " ", "") // 4列目のWCNに対して空白を除去
		targetWCN := strings.ReplaceAll(previousFormatWCN, "\"", "") // 更にバックスラッシュを除去→これで数字の羅列のWCNの出来上がり
		return targetWCN
	}

	return ""
}

// 最新の出口側ゲート待機中のWCNを取得
func getLatestExitWCN() string {

	// ファイルオープン処理
	path := filepath.Join(os.Getenv("HOME"), "opt", "aps", "sbox01", "tc_csv_table", "WCN_rireki.csv")
	file, err := os.Open(path)
	if err != nil {
		log.Println("WCN_rireki.csvが開けません:", err)
		return ""
	}
	defer file.Close()

	// ファイルの最後から読み取り
	scanner := bufio.NewScanner(file)
	var lastLine string
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if lastLine == "" {
		return ""
	}

	columns := strings.Split(lastLine, ",")
	if len(columns) >= 4 {
		previousFormatWCN := strings.ReplaceAll(columns[3], " ", "") // 4列目のWCNに対して空白を除去
		targetWCN := strings.ReplaceAll(previousFormatWCN, "\"", "") // 更にバックスラッシュを除去→これで数字の羅列のWCNの出来上がり
		return targetWCN
	}

	return ""
}

// 入口側ゲート待機中のWCNを許可車両リストと照合し、車両情報を返す関数
func fetchLatestEnterUserInfo() EnterUserInfo {
	// WCNリストから検索
	path := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")
	file, err := os.Open(path)
	if err != nil {
		return EnterUserInfo{"FOE", "エラー", "システムエラー", time.Now().Format("15:04"), true, true}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return EnterUserInfo{"DRE", "エラー", "データ読み込み失敗", time.Now().Format("15:04"), true, true}
	}

	etcNumber := getLatestEnterWCN()
	if etcNumber == "" {
		return EnterUserInfo{"", "未検出", "ゲートオープン", time.Now().Format("15:04"), true, true}
	}

	// var etcNumber string
	// OuterLoop: // ラベルを定義
	// 	for {
	// 		etcNumber = getLatestEnterWCN()
	// 		if etcNumber == "" {
	// 			return EnterUserInfo{"", "未検出", "ゲートオープン", time.Now().Format("15:04"), true, false}
	// 		}
	// 		for _, record := range records {
	// 			if len(record) >= 3 && strings.ReplaceAll(record[0], " ", "") == etcNumber && record[3] == "入場" {
	// 				if record[1] == "来客車" {
	// 					return EnterUserInfo{
	// 						ETCNumber:  etcNumber,
	// 						Department: "登録済み",
	// 						Name:       "来客車",
	// 						LastUpdate: time.Now().Format("15:04"),
	// 						IsVisitor:  false,
	// 						IsEnter:    true,
	// 					}
	// 				} else {
	// 					continue OuterLoop
	// 				}
	// 			}
	// 			if len(record) >= 3 && strings.ReplaceAll(record[0], " ", "") == etcNumber && record[3] == "退場" {

	// 				// gateOpenCommand()

	//				return EnterUserInfo{
	//					ETCNumber:  etcNumber,
	//					Department: record[2],
	//					Name:       record[1],
	//					LastUpdate: time.Now().Format("15:04"),
	//					IsVisitor:  false,
	//					IsEnter:    true,
	//				}
	//			}
	//		}
	//		return EnterUserInfo{etcNumber, "来客車", "ゲートオープン", time.Now().Format("15:04"), true, true}
	//	}
	// etcNumber = getLatestEnterWCN()
	// if etcNumber == "" {
	// 	return EnterUserInfo{"", "未検出", "ゲートオープン", time.Now().Format("15:04"), true}
	// }
	// for _, record := range records {
	// 	if len(record) >= 3 && strings.ReplaceAll(record[0], " ", "") == etcNumber {
	// 		return EnterUserInfo{
	// 			ETCNumber:  etcNumber,
	// 			Department: record[2],
	// 			Name:       record[1],
	// 			LastUpdate: time.Now().Format("15:04"),
	// 			IsVisitor:  false,
	// 		}
	// 	}
	// }

	// CSV内で該当するETC番号を検索
	found := false
	for _, record := range records {
		if len(record) >= 1 {
			recordETC := strings.ReplaceAll(strings.TrimSpace(record[0]), " ", "")
			if recordETC == etcNumber {
				found = true
				// 入場か退場かを判断
				isEnter := true
				if len(record) >= 4 && record[3] == "退場" {
					isEnter = false
				}

				// 来客車かどうかを判定
				isVisitor := false
				if len(record) >= 2 && record[1] == "来客車" {
					isVisitor = true
				}

				return EnterUserInfo{
					ETCNumber:  etcNumber,
					Department: record[2],
					Name:       record[1],
					LastUpdate: time.Now().Format("15:04"),
					IsVisitor:  isVisitor,
					IsEnter:    isEnter,
				}
			}
		}
	}

	// CSVに存在しない場合は来客車
	if !found {
		return EnterUserInfo{
			ETCNumber:  etcNumber,
			Department: "来客車",
			Name:       "来客車",
			LastUpdate: time.Now().Format("15:04"),
			IsVisitor:  true,
			IsEnter:    true,
		}
	}

	return EnterUserInfo{"", "未検出", "ゲートオープン", time.Now().Format("15:04"), true, true}
}

func handleOpenGate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ゲートを開ける(この部分はラインアイと接続時のみコメントアウトを外す)
	// gateOpenCommand()

	// 来客車の場合のみ登録
	visitorWCN := getLatestEnterWCN()
	if visitorWCN != "" {
		registerVisitor(visitorWCN)
	}

	w.Write([]byte("ゲートを開放しました"))
}

// ランダムな車種を取得（デモ用）
func getRandomVehicleType() string {
	vehicleTypes := []string{"セダン", "トラック", "軽トラック", "SUV/MPV", "バン", "バス"}
	rand.Seed(time.Now().UnixNano())
	return vehicleTypes[rand.Intn(len(vehicleTypes))]
}

func gateOpenCommand() {

	// 接続先アドレスとポートを指定
	address := "192.168.110.101:10003" // 適切なIPアドレスとポートに変更してください

	// TCP接続を確立
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("接続失敗: %v", err)
	}
	defer conn.Close()

	// 送信するコマンド (F0h,00000001b)は開くトリガーをONにしただけ
	// F0h = 0xF0 (16進数)
	// 00000001b = 0x01 (16進数)
	command := []byte{0xF0, 0x01}

	// コマンドを送信
	_, err = conn.Write(command)
	if err != nil {
		log.Fatalf("送信失敗: %v", err)
	}

	log.Println("コマンド送信完了:", command)
}

func gateCloseCommand() {
	// 接続先アドレスとポートを指定
	address := "192.168.110.101:10003" // 適切なIPアドレスとポートに変更してください

	// TCP接続を確立
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalf("接続失敗: %v", err)
	}
	defer conn.Close()

	// 送信するコマンド (F0h,00000010b)は閉じるトリガーをONにしただけ
	// F0h = 0xF0 (16進数)
	// 00000010b = 0x02 (16進数)
	command := []byte{0xF0, 0x02}

	// コマンドを送信
	_, err = conn.Write(command)
	if err != nil {
		log.Fatalf("送信失敗: %v", err)
	}

	// 閉じるトリガーをONにしたらOFFに戻しておく
	// F0h = 0xF0 (16進数)
	// 00000000b = 0x0 (16進数)
	command = []byte{0xF0, 0x0}

	// コマンドを送信
	_, err = conn.Write(command)
	if err != nil {
		log.Fatalf("送信失敗: %v", err)
	}

	log.Println("コマンド送信完了:", command)
}

func registerVisitor(visitorWCN string) {
	// CSVファイルに追記
	csvPath := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")
	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("CSVファイルオープンエラー:", err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	newRecord := []string{
		visitorWCN,
		"来客車",
		"登録済み",
		"入場",
		time.Now().Format("15:04"), // 現在時刻
		"",                         // 退場時間
		"",                         // 滞在時間
	}
	if err := writer.Write(newRecord); err != nil {
		log.Println("CSV書き込みエラー:", err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		// return
	}
	writer.Flush()
}

func fetchLatestExitUserInfo() ExitUserInfo {
	// var etcNumber string
	// for {
	// 	etcNumber = getLatestExitWCN()
	// 	if etcNumber != readExitWCN {
	// 		readExitWCN = etcNumber
	// 		break
	// 	}
	// }

	// if etcNumber == "" {
	// 	return ExitUserInfo{"", "未検出", "ゲートオープン", time.Now().Format("15:04"), true, false}
	// }

	// WCNリストから検索
	path := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")
	file, err := os.Open(path)
	if err != nil {
		return ExitUserInfo{"FOE", "エラー", "システムエラー", time.Now().Format("15:04"), true, true}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return ExitUserInfo{"DRE", "エラー", "データ読み込み失敗", time.Now().Format("15:04"), true, true}
	}

	// for _, record := range records {
	// 	if len(record) >= 3 && strings.ReplaceAll(record[0], " ", "") == etcNumber {
	// 		command := exec.Command("ping", "-c", "4", "192.168.110.1")
	// 		output, a := command.CombinedOutput()
	// 		if a != nil {
	// 			fmt.Printf("Ping実行エラー: %v\n", err)
	// 		}

	// 		// 本当にこれでpingコマンドが実行されているかコンソールに表示して検証
	// 		fmt.Println(string(output))

	// 		return ExitUserInfo{
	// 			ETCNumber:  etcNumber,
	// 			Department: record[2],
	// 			Name:       record[1],
	// 			LastUpdate: time.Now().Format("15:04"),
	// 			IsVisitor:  false,
	// 			IsExit:     true,
	// 		}
	// 	}
	// }

	// return ExitUserInfo{etcNumber, "来客車", "ゲートオープン", time.Now().Format("15:04"), true, true}

	var etcNumber string
OuterLoop: // ラベルを定義
	for {
		etcNumber = getLatestExitWCN()
		if etcNumber == "" {
			return ExitUserInfo{"", "未検出", "ゲートオープン", time.Now().Format("15:04"), true, false}
		}
		for _, record := range records {
			if len(record) >= 3 && strings.ReplaceAll(record[0], " ", "") == etcNumber && record[3] == "退場" {
				continue OuterLoop
			}
			if len(record) >= 3 && strings.ReplaceAll(record[0], " ", "") == etcNumber && record[3] == "入場" {

				gateOpenCommand()

				return ExitUserInfo{
					ETCNumber:  etcNumber,
					Department: record[2],
					Name:       record[1],
					LastUpdate: time.Now().Format("15:04"),
					IsVisitor:  false,
					IsExit:     true,
				}
			}
		}
	}
}

// 事前登録者の情報を「WCN」と「カメラで捉えた画像の進行方向」によって更新
func updateUser(enterUser EnterUserInfo) {
	csvPath := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")

	file, err := os.Open(csvPath)
	if err != nil {
		log.Println("CSVファイルオープンエラー:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("CSV読み込みエラー:", err)
		return
	}

	// // 画像処理によって入退場を判別
	// direction := imageRecognition()

	// 更新対象を検索
	updated := false
	currentTime := time.Now().Format("15:04")
	for i, record := range records {
		if len(record) > 0 && strings.TrimSpace(record[0]) == enterUser.ETCNumber /*&& strings.Contains(direction, "F")*/ {
			// ステータスと入場時間を更新
			if len(record) >= 4 {
				records[i][3] = "入場"
				records[i][4] = currentTime
				updated = true
			}
			fmt.Println("入場しました")
			break
		} else if len(record) > 0 && strings.TrimSpace(record[0]) == enterUser.ETCNumber /*&& strings.Contains(direction, "R")*/ {
			// ステータスと退場時間を更新
			if len(record) >= 4 {
				records[i][3] = "退場"
				records[i][5] = currentTime
				updated = true
			}
			fmt.Println("退場しました")
			break
		}
	}

	if !updated {
		log.Println("該当ユーザーが見つかりませんでした")
		return
	}

	// ファイルに書き戻し
	file.Close()
	output, err := os.Create(csvPath)
	if err != nil {
		log.Println("CSV書き込みオープンエラー:", err)
		return
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	if err := writer.WriteAll(records); err != nil {
		log.Println("CSV書き込みエラー:", err)
	}
}

// 退場していく事前登録者の情報を更新
func updateCSVExit(exitUser ExitUserInfo) {
	csvPath := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")

	file, err := os.Open(csvPath)
	if err != nil {
		log.Println("CSVファイルオープンエラー:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("CSV読み込みエラー:", err)
		return
	}

	// 更新対象を検索
	updated := false
	currentTime := time.Now().Format("15:04")
	for i, record := range records {
		if len(record) > 0 && strings.TrimSpace(record[0]) == exitUser.ETCNumber {
			// ステータスと退場時間を更新
			if len(record) >= 4 {
				records[i][3] = "退場"
				records[i][5] = currentTime
				updated = true
			}
			break
		}
	}

	if !updated {
		log.Println("該当ユーザーが見つかりませんでした")
		return
	}

	// ファイルに書き戻し
	file.Close()
	output, err := os.Create(csvPath)
	if err != nil {
		log.Println("CSV書き込みオープンエラー:", err)
		return
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	if err := writer.WriteAll(records); err != nil {
		log.Println("CSV書き込みエラー:", err)
	}
	log.Printf("退場記録を更新: %s %s\n", exitUser.ETCNumber, currentTime)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	cars := loadCSVData()
	tmpl := template.Must(template.ParseFiles(filepath.Join("www", "permitted_cars.html")))
	tmpl.Execute(w, cars)
}

func handleList(w http.ResponseWriter, r *http.Request) {
	cars := loadCSVData()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	selectedIDs := strings.Split(r.FormValue("selectedIds"), ",")
	log.Printf("削除対象のETC番号: %v\n", selectedIDs)

	// CSVファイルの読み込み
	csvPath := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")
	file, err := os.Open(csvPath)
	if err != nil {
		log.Println("CSVファイルオープンエラー:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Println("CSV読み込みエラー:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 削除対象を除外した新しいデータ作成
	var newRecords [][]string
	header := records[0] // ヘッダー行を保持
	newRecords = append(newRecords, header)

	for i, row := range records {
		if i == 0 {
			continue // ヘッダー行はスキップ
		}
		if len(row) == 0 {
			continue
		}

		// 選択されたETC番号と一致しない行のみ保持
		shouldKeep := true
		for _, id := range selectedIDs {
			if row[0] == id {
				shouldKeep = false
				break
			}
		}

		if shouldKeep {
			newRecords = append(newRecords, row)
		}
	}

	// CSVファイルを上書き保存
	file.Close() // 読み込み用ファイルを一旦閉じる

	output, err := os.Create(csvPath)
	if err != nil {
		log.Println("CSVファイル作成エラー:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer output.Close()

	writer := csv.NewWriter(output)
	if err := writer.WriteAll(newRecords); err != nil {
		log.Println("CSV書き込みエラー:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Flush()
	log.Println("CSVファイルを更新しました")
	w.Write([]byte("選択した行を削除しました"))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// フォームデータの取得
	etcNumber := r.FormValue("etcNumber")
	name := r.FormValue("name")
	department := r.FormValue("department")

	// バリデーション
	if etcNumber == "" || name == "" || department == "" {
		http.Error(w, "すべてのフィールドが必須です", http.StatusBadRequest)
		return
	}

	// CSVファイルに追記
	csvPath := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")
	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("CSVファイルオープンエラー:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	newRecord := []string{
		etcNumber,
		name,
		department,
		"入場",                       // デフォルト状態
		time.Now().Format("15:04"), // 現在時刻
		"",                         // 退場時間
		"",                         // 滞在時間
	}

	if err := writer.Write(newRecord); err != nil {
		log.Println("CSV書き込みエラー:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	writer.Flush()
	log.Printf("新規登録: %v\n", newRecord)
	w.Write([]byte("登録が完了しました"))
}

func loadCSVData() []CarRecord {
	csvPath := filepath.Join(os.Getenv("HOME"), "sk_prj", "enterprise", "etcs", "etcs_hasuda", "tc", "gatesystem", "WCN_Register_List.csv")
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal("CSVファイルが開けません:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("CSV読み込みエラー:", err)
	}

	var cars []CarRecord
	for i, row := range records {
		if i == 0 || len(row) < 6 {
			continue
		}
		cars = append(cars, CarRecord{
			ETCNumber:  row[0],
			Name:       row[1],
			Department: row[2],
			Status:     row[3],
			EntryTime:  row[4],
			ExitTime:   row[5],
			StayTime:   calculateStayTime(row[4], row[5]),
		})
	}

	return cars
}

func calculateStayTime(entry, exit string) string {
	// 退場時間が未設定の場合は""を返す
	if exit == "" {
		return ""
	}

	// 時間のパース（HH:MM形式）
	parseTime := func(s string) (hour, min int, err error) {
		parts := strings.Split(s, ":")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid time format")
		}
		hour, err1 := strconv.Atoi(parts[0])
		min, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return 0, 0, fmt.Errorf("invalid time value")
		}
		return hour, min, nil
	}

	// 入場時間と退場時間をパース
	entryHour, entryMin, err := parseTime(entry)
	if err != nil {
		return "フォーマットエラー"
	}

	exitHour, exitMin, err := parseTime(exit)
	if err != nil {
		return "フォーマットエラー"
	}

	// 時間差分計算（分単位）
	totalEntryMinutes := entryHour*60 + entryMin
	totalExitMinutes := exitHour*60 + exitMin

	// 24時間を超える場合を考慮
	if totalExitMinutes < totalEntryMinutes {
		totalExitMinutes += 24 * 60
	}

	diffMinutes := totalExitMinutes - totalEntryMinutes

	// 時間と分に変換
	hours := diffMinutes / 60
	minutes := diffMinutes % 60

	return fmt.Sprintf("%d時間%02d分", hours, minutes)
}

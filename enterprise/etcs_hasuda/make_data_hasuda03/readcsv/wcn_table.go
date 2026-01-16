package readcsv

import (
	"bufio"
	"log"
	"os"
)

// Read wcnテーブルファイルを読み込む
func Wcn_table() []string {

    //	実行環境「ac」直下のディレクトリにあるファイルを直接参照する。
    file, err := os.Open("../ac/tc_wcn_table/WCN_table.csv")
	if err != nil {
        log.Printf("../ac/tc_wcn_table/WCN_table.csv Open Error!!\n")
		log.Fatal(err)
	}
	defer file.Close()

	var wcn []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		wcn = append(wcn, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
        log.Printf("../ac/tc_wcn_table/WCN_table.csv Scanner Error!!\n")
		log.Fatal(err)
	}

	return wcn
}

/* 指定ファイル(CSV)の読込  */
func Read(file_name string) []string {

    //	指定のディレクトリにあるCSVファイルを直接参照する。
    file, err := os.Open(file_name)
	if err != nil {
		//log.Fatal(err)
        log.Printf("Nothing!! %s\n",file_name)
	}
	defer file.Close()

	var result []string

    // CSVファイルを１行ずつ読み込んで配列化
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
        //		log.Fatal(err)
        log.Printf("Can't scan !! %s\n",file_name)
	}
	return result
}

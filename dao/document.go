package dao

import (
	"fmt"
)

func GetDocumentId(title string) (int, error) {
	stmt, err := db.Prepare("SELECT id FROM documents WHERE title = ?;")
	if err != nil {
		fmt.Println("failed to get document id, prepare sql err: ", err.Error())
		return 0, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(title).Scan(&id)
	if err != nil {
		fmt.Println("failed to get document id, err: ", err.Error())
		return 0, err
	}
	return id, nil
}

func GetDocumentTitle(id int) (string, error) {
	stmt, err := db.Prepare("SELECT title FROM documents WHERE id = ?;")
	if err != nil {
		fmt.Println("failed to get document title, prepare sql err: ", err.Error())
		return "", err
	}
	defer stmt.Close()

	var title string
	err = stmt.QueryRow(id).Scan(&title)
	if err != nil {
		fmt.Println("failed to get document title, err: ", err.Error())
		return "", err
	}
	return title, nil
}

func InsertDocument(title, body string) bool {
	stmt, err := db.Prepare("INSERT INTO documents (title, body) VALUES (?, ?);")
	if err != nil {
		fmt.Println("failed to insert document, prepare sql err: ", err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(title, body)
	if err != nil {
		fmt.Println("failed to insert document, err: ", err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	fmt.Println("failed to insert document no rows affected, err: ", err)
	return false
}

func UpdateDocument(id int, body string) bool {
	stmt, err := db.Prepare("UPDATE documents set body = ? WHERE id = ?;")
	if err != nil {
		fmt.Println("failed to update document, prepare sql err: ", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(body, id)
	if err != nil {
		fmt.Println("failed to update document, err: ", err.Error())
		return false
	}
	return true
}

func GetTokenId(token string) (int, int, error) {
	stmt, err := db.Prepare("SELECT id, docs_count FROM tokens WHERE token = ?;")
	if err != nil {
		fmt.Println("failed to get token id, prepare sql err: ", err.Error())
		return 0, 0, err
	}
	defer stmt.Close()

	var id, count int
	err = stmt.QueryRow(token).Scan(&id, &count)
	if err != nil {
		panic(err)
		fmt.Println("failed to get token id, err: ", err.Error())
		return 0, 0, err
	}
	return id, count, nil
}

func GetToken(id int) (string, error) {
	stmt, err := db.Prepare("SELECT token FROM tokens WHERE id = ?;")
	if err != nil {
		fmt.Println("failed to get token, prepare sql err: ", err.Error())
		return "", err
	}
	defer stmt.Close()

	var token string
	err = stmt.QueryRow(id).Scan(&token)
	if err != nil {
		fmt.Println("failed to get token, err: ", err.Error())
		return "", err
	}
	return token, nil
}

func StoreToken(token, postings string) bool {
	stmt, err := db.Prepare("INSERT IGNORE INTO tokens (token, docs_count, postings) VALUES (?, 1, ?);")
	if err != nil {
		fmt.Println("failed to store token, prepare sql err: ", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(token, postings)
	if err != nil {
		fmt.Println("failed to store token, err: ", err.Error())
		return false
	}
	return true
}

func GetPostings(id int) (int, string, error) {
	stmt, err := db.Prepare("SELECT docs_count, postings FROM tokens WHERE id = ?;")
	if err != nil {
		fmt.Println("failed to get postings, prepare sql err: ", err.Error())
		return 0, "", err
	}
	defer stmt.Close()

	var count int
	var postings string
	err = stmt.QueryRow(id).Scan(&count, &postings)
	if err != nil {
		fmt.Println("failed to get document count, err: ", err.Error())
		return 0, "", err
	}
	return count, postings, nil
}

func UpdatePostings(id, count int, postings string) bool {
	stmt, err := db.Prepare("UPDATE tokens SET docs_count = ?, postings = ? WHERE id = ?;")
	if err != nil {
		fmt.Println("failed to update postings, prepare sql err: ", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(count, postings, id)
	if err != nil {
		fmt.Println("failed to update postings, err: ", err.Error())
		return false
	}
	return true
}

func GetSettings(key string) (string, error) {
	stmt, err := db.Prepare("SELECT value FROM settings WHERE key = ?;")
	if err != nil {
		fmt.Println("failed to get settings, prepare sql err: ", err.Error())
		return "", err
	}
	defer stmt.Close()

	var value string
	err = stmt.QueryRow(key).Scan(&value)
	if err != nil {
		fmt.Println("failed to get settings, err: ", err.Error())
		return "", err
	}
	return value, nil
}

func ReplaceSettings(key, value string) bool {
	stmt, err := db.Prepare("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?);")
	if err != nil {
		fmt.Println("failed to replace settings, prepare sql err: ", err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, value)
	if err != nil {
		fmt.Println("failed to replace settings, err: ", err.Error())
		return false
	}
	return true
}

func GetDocumentCount() (int, error) {
	stmt, err := db.Prepare("SELECT COUNT(*) FROM documents;")
	if err != nil {
		fmt.Println("failed to get document count, prepare sql err: ", err.Error())
		return 0, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		fmt.Println("failed to get document count, err: ", err.Error())
		return 0, err
	}
	return count, nil
}

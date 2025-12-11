package dbx

import (
	"dbu-api/internal/env"
	"dbu-api/internal/logger"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	// _ "github.com/lib/pq"
	_ "github.com/go-sql-driver/mysql"
)

// Consulta el número y la fecha de notificación para un faltaID
type NotificationInfo struct {
	NumeroNotificacion string `db:"numero_notificacion"`
	Fecha              string `db:"fecha"`
}

func GetNotificationInfo(faltaID string) (*NotificationInfo, error) {
	db := GetConnection()
	var info NotificationInfo
	err := db.Get(&info, "SELECT numero_notificacion, fecha FROM notificacion WHERE falta_id = ?", faltaID)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// Obtiene o crea el número de notificación único para un faltaID y año
func GetOrCreateNotificationNumber(faltaID string) (string, error) {
	db := GetConnection()
	year := time.Now().Year()
	var numero string
	err := db.Get(&numero, "SELECT numero_notificacion FROM notificacion WHERE falta_id = ?", faltaID)
	if err == nil && numero != "" {
		return numero, nil // Ya existe, retorna el mismo
	}
	// Si no existe, genera nuevo
	next, err := GetNextNotificationNumber()
	if err != nil {
		return "", err
	}
	_, err = db.Exec("INSERT INTO notificacion (falta_id, numero_notificacion, anio) VALUES (?, ?, ?)", faltaID, next, year)
	if err != nil {
		return "", err
	}
	return next, nil
}

// Obtiene el siguiente número de notificación único para el año actual, formato 00001
func GetNextNotificationNumber() (string, error) {
	db := GetConnection()
	year := time.Now().Year()
	var numero int
	err := db.Get(&numero, "SELECT ultimo_numero FROM notificacion_secuencia WHERE anio = ?", year)
	if err != nil {
		// Si no existe, iniciar en 1
		numero = 1
		_, err = db.Exec("INSERT INTO notificacion_secuencia (anio, ultimo_numero) VALUES (?, ?)", year, numero)
		if err != nil {
			return "", err
		}
	} else {
		numero++
		_, err = db.Exec("UPDATE notificacion_secuencia SET ultimo_numero = ? WHERE anio = ?", numero, year)
		if err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%05d", numero), nil
}

var (
	dbx      *sqlx.DB
	once     sync.Once
	DBEngine string
)

func init() {
	once.Do(func() {
		setConnection()
	})
}

func setConnection() {
	var err error
	c := env.NewConfiguration()
	DBEngine = c.DB.Engine

	// Check the connection
	dbx, err = sqlx.Open(DBEngine, connectionString("data"))
	if err != nil {
		logger.Error.Printf("no se puede conectar a la base de datos: %v", err)
		panic(err)
	}
	err = dbx.Ping()
	if err != nil {
		logger.Error.Printf("couldn't connect to database: %v", err)
		panic(err)
	}
	dbx.SetMaxIdleConns(5)
	dbx.SetConnMaxLifetime(2 * time.Minute)
	dbx.SetMaxOpenConns(95)
}

func GetConnection() *sqlx.DB {
	if dbx == nil {
		setConnection()
	}
	return dbx
}

func connectionString(t string) string {
	c := env.NewConfiguration()

	var host, database, username, password, instance string
	var port int
	switch t {
	case "data":
		host = c.DB.Server
		database = c.DB.Name
		username = c.DB.User
		password = c.DB.Password
		instance = c.DB.Instance
		port = c.DB.Port
	default:
		logger.Error.Print("El tipo de conexión no correspondea data/logs")
		return ""
	}
	switch strings.ToLower(DBEngine) {
	case "postgres":
		return fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=disable", database, username, password, host, port)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
	case "sqlserver":
		return fmt.Sprintf(
			"server=%s\\%s;User id=%s;database=%s;password=%s;port=%d", host, instance, username, database, password, port)
	}
	logger.Error.Print("el motor de bases de datos solicitado no está configurado aún")

	return ""
}

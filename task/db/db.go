package db

import (
	"encoding/json"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/oklog/ulid/v2"
	bolt "go.etcd.io/bbolt"
)

func SetupDB() (*bolt.DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(homeDir, taskDBFile)

	db, err := bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(taskBucket))
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists([]byte(completedBucket))
		return err
	}); err != nil {
		return nil, err
	}

	return db, nil
}

func AddTask(des string, db *bolt.DB) (*Task, error) {
	var task *Task
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))

		entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
		id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

		task = &Task{
			ULID:        id,
			Description: des,
			CreatedAt:   time.Now(),
		}

		byteTask, err := json.Marshal(task)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), byteTask)
	}); err != nil {
		return nil, err
	}
	return task, nil
}

func ListTasks(db *bolt.DB) ([]Task, error) {
	var tasks []Task
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task Task
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			tasks = append(tasks, task)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return tasks, nil
}

func SetTaskAsCompleted(task Task, db *bolt.DB) (*CompletedTask, error) {
	var ct CompletedTask
	if err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(completedBucket))
		ct = CompletedTask{
			ULID:        task.ULID,
			Description: task.Description,
			CompletedAt: time.Now(),
		}
		byteData, err := json.Marshal(ct)
		if err != nil {
			return err
		}
		if err := b.Put([]byte(ct.ULID), byteData); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &ct, nil
}

func DeleteTask(id string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(taskBucket))
		return b.Delete([]byte(id))
	})
}

func ListCompletedTasks(db *bolt.DB) ([]CompletedTask, error) {
	var tasks []CompletedTask
	if err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(completedBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var task CompletedTask
			if err := json.Unmarshal(v, &task); err != nil {
				return err
			}
			tasks = append(tasks, task)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return tasks, nil
}

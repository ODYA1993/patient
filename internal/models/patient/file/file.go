package file

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
	"patients/internal/models/patient"
	"patients/pkg/logging"
)

var ErrNotFound = errors.New("not found")

const (
	patientDataFilePath = "data/list_patients.json"
)

type repositoryFile struct {
	Patients []*patient.Person `json:"patients"`
	sync.RWMutex
}

func NewPatientRepo(logger *logging.Logger) patient.Storage {
	return &repositoryFile{
		Patients: make([]*patient.Person, 0),
	}
}

func (r *repositoryFile) FindAll(ctx context.Context) ([]*patient.Person, error) {
	if len(r.Patients) == 0 {
		if err := r.loadPatients(); err != nil {
			return nil, err
		}
	}

	return r.Patients, nil
}

func (r *repositoryFile) Create(ctx context.Context, person *patient.Person) error {
	person.Guid = uuid.New().String()

	// Загружаем существующий список пациентов из файла
	err := r.loadPatients()
	if err != nil {
		return err
	}

	// Добавляем нового пациента в список
	r.Patients = append(r.Patients, person)

	// Сохраняем обновленный список пациентов в файл
	return r.savePatients()
}

func (r *repositoryFile) Update(ctx context.Context, person *patient.Person) error {
	for i, pat := range r.Patients {
		if pat.Guid == person.Guid {
			// Изменяем данные пациента в слайсе
			r.Patients[i] = person
			// Сохраняем обновленные данные в файле
			if err := r.savePatients(); err != nil {
				return err
			}
			return nil
		}
	}

	return fmt.Errorf("patient with GUID %s not found", person.Guid)
}

func (r *repositoryFile) Delete(ctx context.Context, guid interface{}) error {
	for i, pat := range r.Patients {
		if pat.Guid == guid.(string) {
			r.Patients = append(r.Patients[:i], r.Patients[i+1:]...)
			return r.savePatients()
		}
	}

	return ErrNotFound
}

func (r *repositoryFile) loadPatients() error {
	r.RLock()
	defer r.RUnlock()

	data, err := os.ReadFile(patientDataFilePath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &r.Patients)
}

func (r *repositoryFile) savePatients() error {
	r.Lock()
	defer r.Unlock()

	// Открываем файл для записи
	file, err := os.Create(patientDataFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Сериализуем весь слайс пациентов в JSON массив
	data, err := json.Marshal(r.Patients)
	if err != nil {
		return err
	}

	// Записываем JSON массив в файл
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	// Добавляем символ новой строки после каждого элемента
	_, err = file.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}

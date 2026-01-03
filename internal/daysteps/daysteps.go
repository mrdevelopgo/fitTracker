package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// Убираем пробелы в начале и конце
	data = strings.TrimSpace(data)

	// parts разделяет слайс на строки по запятой
	parts := strings.Split(data, ",")

	// Проверка количества частей равна 2
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат строки: ожидается 'шаги,время'")
	}

	// Преобразование количества шагов в тип int
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат шагов: %s", stepsStr)
	}

	// Проверка, что количество шагов > 0
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть больше 0: %d", steps)
	}

	// Преобразование времени (вторая часть) в time.Duration
	durationStr := strings.TrimSpace(parts[1])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("неверный формат времени: %s", durationStr)
	}

	// Проверка, что продолжительность > 0
	if duration <= 0 {
		return 0, 0, fmt.Errorf("продолжительность должна быть больше 0: %v", duration)
	}

	// Возвращение успешного результата
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// 1. Получаем данные через parsePackage
	steps, duration, err := parsePackage(data)
	if err != nil {
		// Возвращаем строку с ошибкой
		return fmt.Sprintf("Ошибка обработки данных '%s': %v", data, err)
	}

	// 2. Проверяем, что шаги > 0 (можно удалить, т.к. уже проверено в parsePackage)
	if steps <= 0 {
		return ""
	}

	// 3. Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 4. Переводим в километры
	distanceKm := distanceMeters / float64(mInKm)

	// 5. Вычисляем калории (используем duration!)
	calories := calculateCalories(steps, weight, height, duration)

	// 6. Формируем строку результата
	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Время активности: %v.\n"+ // Добавили информацию о времени
			"Вы сожгли %.2f ккал.",
		steps, distanceKm, duration, calories,
	)

	return result
}

// Вспомогательная функция для расчета калорий
func calculateCalories(steps int, weight, height float64, duration time.Duration) float64 {
	// Простой расчет калорий на основе шагов и времени
	// В реальном приложении здесь будет вызов функции из spentcalories

	if duration <= 0 || weight <= 0 {
		return 0.0
	}

	// Примерная формула: 0.04 ккал на шаг + влияние скорости
	baseCalories := float64(steps) * 0.04

	// Учет скорости: чем быстрее, тем больше калорий
	speed := 0.0
	distanceKm := float64(steps) * stepLength / float64(mInKm)
	if duration.Hours() > 0 {
		speed = distanceKm / duration.Hours()
	}

	// Множитель скорости
	speedMultiplier := 1.0
	if speed > 5.0 {
		speedMultiplier = 1.2
	}
	if speed > 7.0 {
		speedMultiplier = 1.5
	}

	return baseCalories * speedMultiplier
}
